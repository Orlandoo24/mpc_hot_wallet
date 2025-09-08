package monitor

import (
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/url"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

// 事件签名常量
var (
	// TransferEventSignature Transfer(address,address,uint256)
	TransferEventSignature = crypto.Keccak256Hash([]byte("Transfer(address,address,uint256)"))
	// ApprovalEventSignature Approval(address,address,uint256)
	ApprovalEventSignature = crypto.Keccak256Hash([]byte("Approval(address,address,uint256)"))
	// SwapEventSignature Swap 事件签名 (PancakeSwap V2)
	SwapEventSignature = crypto.Keccak256Hash([]byte("Swap(address,uint256,uint256,uint256,uint256,address)"))
	// SwapV3EventSignature PancakeSwap V3 Swap
	SwapV3EventSignature = crypto.Keccak256Hash([]byte("Swap(address,address,int256,int256,uint160,uint128,int24)"))
	// BridgeEventSignature Bridge 相关事件 (常见跨链桥)
	BridgeEventSignature = crypto.Keccak256Hash([]byte("Bridge(address,address,uint256,uint256,bytes)"))
	// DepositEventSignature & WithdrawEventSignature Deposit/Withdraw (Wrapped tokens)
	DepositEventSignature  = crypto.Keccak256Hash([]byte("Deposit(address,uint256)"))
	WithdrawEventSignature = crypto.Keccak256Hash([]byte("Withdrawal(address,uint256)"))
)

// 已知合约地址 (可扩展)
var (
	// PancakeV2Router PancakeSwap V2 Router
	PancakeV2Router = common.HexToAddress("0x10ED43C718714eb63d5aA57B78B54704E256024E")
	// PancakeV3Router PancakeSwap V3 Router
	PancakeV3Router = common.HexToAddress("0x13f4EA83D0bd40E75C8222255bc855a974568Dd4")
	// BSCBridge BSC Bridge 合约地址 (示例)
	BSCBridge = common.HexToAddress("0x0000000000000000000000000000000000000000")
)

// LiFiResponse LI.FI API响应结构
type LiFiResponse struct {
	TransactionType string `json:"transactionType"`
	Action          string `json:"action"`
	Tool            string `json:"tool"`
	FromToken       string `json:"fromToken"`
	ToToken         string `json:"toToken"`
	Amount          string `json:"amount"`
}

// LogParser 日志解析器
type LogParser struct {
	knownContracts map[common.Address]string // 已知合约地址映射
	lifiAPIKey     string                    // LI.FI API密钥
}

// NewLogParser 创建日志解析器
func NewLogParser() *LogParser {
	parser := &LogParser{
		knownContracts: make(map[common.Address]string),
		lifiAPIKey:     "", // 可以从环境变量或配置文件读取
	}

	// 初始化已知合约
	parser.knownContracts[PancakeV2Router] = "PancakeSwap-V2"
	parser.knownContracts[PancakeV3Router] = "PancakeSwap-V3"
	parser.knownContracts[BSCBridge] = "BSC-Bridge"

	return parser
}

// ParseLogs 解析交易日志，生成TokenEvent列表
func (p *LogParser) ParseLogs(logs []*types.Log, tx *types.Transaction, blockNumber uint64, timestamp int64, chainId uint64) []*TokenEvent {
	var events []*TokenEvent

	// 获取交易发送方
	var fromAddr string
	signer := types.LatestSignerForChainID(tx.ChainId())
	if from, err := signer.Sender(tx); err == nil {
		fromAddr = from.Hex()
	}

	// 获取交易接收方
	var toAddr string
	if tx.To() != nil {
		toAddr = tx.To().Hex()
	}

	// 尝试使用LI.FI API增强解析
	p.enhanceWithLiFiAPI(tx, chainId)

	// 解析每个日志
	for _, vLog := range logs {
		if len(vLog.Topics) == 0 {
			continue
		}

		eventSig := vLog.Topics[0]

		switch eventSig {
		case TransferEventSignature:
			if event := p.parseTransferEvent(vLog, tx, blockNumber, timestamp, chainId, fromAddr, toAddr); event != nil {
				events = append(events, event)
			}
		case ApprovalEventSignature:
			if event := p.parseApprovalEvent(vLog, tx, blockNumber, timestamp, chainId, fromAddr, toAddr); event != nil {
				events = append(events, event)
			}
		case SwapEventSignature, SwapV3EventSignature:
			if event := p.parseSwapEvent(vLog, tx, blockNumber, timestamp, chainId, fromAddr, toAddr); event != nil {
				events = append(events, event)
			}
		case BridgeEventSignature, DepositEventSignature, WithdrawEventSignature:
			if event := p.parseBridgeEvent(vLog, tx, blockNumber, timestamp, chainId, fromAddr, toAddr); event != nil {
				events = append(events, event)
			}
		}
	}

	// 如果没有解析到任何事件，但是有原生代币转账，创建一个基础事件
	if len(events) == 0 && tx.Value().Cmp(big.NewInt(0)) > 0 {
		// 判断原生代币转账方向
		var direction string
		if p.isNativeReceiveEvent(fromAddr, toAddr, fromAddr) {
			direction = "IN" // 接收原生代币
		} else {
			direction = "OUT" // 发送原生代币
		}

		event := &TokenEvent{
			BlockNumber: blockNumber,
			TxHash:      tx.Hash().Hex(),
			Timestamp:   timestamp,
			EventType:   "NativeTransfer", // 标准事件名称
			Direction:   direction,
			FromAddr:    fromAddr,
			ToAddr:      toAddr,
			TokenAddr:   "0x0000000000000000000000000000000000000000", // 原生代币地址
			Amount:      tx.Value().String(),
			ChainId:     chainId,
		}
		events = append(events, event)
	}

	return events
}

// parseTransferEvent 解析Transfer事件
func (p *LogParser) parseTransferEvent(vLog *types.Log, tx *types.Transaction, blockNumber uint64, timestamp int64, chainId uint64, txFromAddr, txToAddr string) *TokenEvent {
	if len(vLog.Topics) < 3 || len(vLog.Data) != 32 {
		return nil
	}

	from := common.BytesToAddress(vLog.Topics[1].Bytes())
	to := common.BytesToAddress(vLog.Topics[2].Bytes())
	amount := new(big.Int).SetBytes(vLog.Data)

	// 判断Transfer方向
	var direction string
	if p.isReceiveEvent(from, to, txFromAddr) {
		direction = "IN" // 接收代币
	} else {
		direction = "OUT" // 发送代币
	}

	return &TokenEvent{
		BlockNumber: blockNumber,
		TxHash:      tx.Hash().Hex(),
		Timestamp:   timestamp,
		EventType:   "Transfer", // 标准ABI事件名称
		Direction:   direction,
		FromAddr:    from.Hex(),
		ToAddr:      to.Hex(),
		TokenAddr:   vLog.Address.Hex(),
		Amount:      amount.String(),
		ChainId:     chainId,
	}
}

// parseApprovalEvent 解析Approval事件
func (p *LogParser) parseApprovalEvent(vLog *types.Log, tx *types.Transaction, blockNumber uint64, timestamp int64, chainId uint64, txFromAddr, txToAddr string) *TokenEvent {
	if len(vLog.Topics) < 3 || len(vLog.Data) != 32 {
		return nil
	}

	owner := common.BytesToAddress(vLog.Topics[1].Bytes())
	spender := common.BytesToAddress(vLog.Topics[2].Bytes())
	amount := new(big.Int).SetBytes(vLog.Data)

	// 使用标准ABI事件名称
	eventType := "Approval"

	return &TokenEvent{
		BlockNumber: blockNumber,
		TxHash:      tx.Hash().Hex(),
		Timestamp:   timestamp,
		EventType:   eventType,
		Direction:   "NONE", // Approval事件没有方向性
		FromAddr:    owner.Hex(),
		ToAddr:      spender.Hex(),
		TokenAddr:   vLog.Address.Hex(),
		Amount:      amount.String(),
		ChainId:     chainId,
	}
}

// parseSwapEvent 解析Swap事件
func (p *LogParser) parseSwapEvent(vLog *types.Log, tx *types.Transaction, blockNumber uint64, timestamp int64, chainId uint64, txFromAddr, txToAddr string) *TokenEvent {
	// 基础swap事件，具体金额解析需要根据不同DEX的事件结构
	var amount *big.Int

	if len(vLog.Data) >= 32 {
		// 简化处理，取第一个uint256值作为金额
		amount = new(big.Int).SetBytes(vLog.Data[:32])
	} else {
		amount = big.NewInt(0)
	}

	return &TokenEvent{
		BlockNumber: blockNumber,
		TxHash:      tx.Hash().Hex(),
		Timestamp:   timestamp,
		EventType:   "Swap",
		Direction:   "NONE", // Swap事件没有明确方向性
		FromAddr:    txFromAddr,
		ToAddr:      vLog.Address.Hex(), // Swap合约地址
		TokenAddr:   vLog.Address.Hex(),
		Amount:      amount.String(),
		ChainId:     chainId,
	}
}

// parseBridgeEvent 解析Bridge相关事件
func (p *LogParser) parseBridgeEvent(vLog *types.Log, tx *types.Transaction, blockNumber uint64, timestamp int64, chainId uint64, txFromAddr, txToAddr string) *TokenEvent {
	var amount *big.Int
	var from, to string

	// 根据不同的事件类型解析
	if vLog.Topics[0] == BridgeEventSignature && len(vLog.Topics) >= 3 {
		from = common.BytesToAddress(vLog.Topics[1].Bytes()).Hex()
		to = common.BytesToAddress(vLog.Topics[2].Bytes()).Hex()
		if len(vLog.Data) >= 32 {
			amount = new(big.Int).SetBytes(vLog.Data[:32])
		}
	} else if vLog.Topics[0] == DepositEventSignature || vLog.Topics[0] == WithdrawEventSignature {
		from = txFromAddr
		to = vLog.Address.Hex()
		if len(vLog.Data) >= 32 {
			amount = new(big.Int).SetBytes(vLog.Data[:32])
		}
	}

	if amount == nil {
		amount = big.NewInt(0)
	}

	// 根据事件类型确定ABI标准名称
	var eventType string
	switch vLog.Topics[0] {
	case BridgeEventSignature:
		eventType = "Bridge"
	case DepositEventSignature:
		eventType = "Deposit"
	case WithdrawEventSignature:
		eventType = "Withdrawal"
	default:
		eventType = "Bridge"
	}

	return &TokenEvent{
		BlockNumber: blockNumber,
		TxHash:      tx.Hash().Hex(),
		Timestamp:   timestamp,
		EventType:   eventType,
		Direction:   "NONE", // Bridge/Deposit/Withdrawal事件暂时不区分方向
		FromAddr:    from,
		ToAddr:      to,
		TokenAddr:   vLog.Address.Hex(),
		Amount:      amount.String(),
		ChainId:     chainId,
	}
}

// isReceiveEvent 判断Transfer事件是否为接收事件
func (p *LogParser) isReceiveEvent(from, to common.Address, txFromAddr string) bool {
	// 判断：如果Transfer事件的to地址是交易发起者，则为接收
	// 这种情况通常发生在：别人向你转账，或者你从合约中提取代币
	txFrom := common.HexToAddress(txFromAddr)
	return to == txFrom
}

// isNativeReceiveEvent 判断原生代币转账是否为接收事件
func (p *LogParser) isNativeReceiveEvent(fromAddr, toAddr, txFromAddr string) bool {
	// 对于原生代币转账，交易发起者就是发送方
	// 所以如果txFromAddr == fromAddr，则为发送(OUT)
	// 这里的逻辑与Transfer事件不同，因为原生转账没有合约中介
	return false // 原生转账的交易发起者总是发送方，所以总是OUT
}

// AddKnownContract 添加已知合约
func (p *LogParser) AddKnownContract(address common.Address, name string) {
	p.knownContracts[address] = name
}

// GetContractName 获取合约名称
func (p *LogParser) GetContractName(address common.Address) string {
	if name, exists := p.knownContracts[address]; exists {
		return name
	}
	return "Unknown"
}

// IsSwapContract 判断是否为交换合约
func (p *LogParser) IsSwapContract(address common.Address) bool {
	name := p.GetContractName(address)
	return strings.Contains(strings.ToLower(name), "swap") ||
		strings.Contains(strings.ToLower(name), "pancake") ||
		strings.Contains(strings.ToLower(name), "uniswap")
}

// IsBridgeContract 判断是否为跨链桥合约
func (p *LogParser) IsBridgeContract(address common.Address) bool {
	name := p.GetContractName(address)
	return strings.Contains(strings.ToLower(name), "bridge") ||
		strings.Contains(strings.ToLower(name), "cross")
}

// FormatTokenAmount 格式化代币金额为可读格式
func FormatTokenAmount(amountStr string, tokenAddr string, chainId uint64) string {
	amount := new(big.Int)
	amount.SetString(amountStr, 10)

	// 检查是否为原生代币 (BNB/ETH等)
	if tokenAddr == "0x0000000000000000000000000000000000000000" {
		return formatNativeToken(amount, chainId)
	}

	// 对于ERC20代币，默认使用18位小数
	return formatERC20Token(amount, tokenAddr, 18)
}

// formatNativeToken 格式化原生代币
func formatNativeToken(amount *big.Int, chainId uint64) string {
	// 1 ETH/BNB = 10^18 wei
	eth := new(big.Float).SetInt(amount)
	eth.Quo(eth, new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)))

	var symbol string
	switch chainId {
	case 1, 5, 11155111: // Ethereum 主网, Goerli, Sepolia
		symbol = "ETH"
	case 56, 97: // BSC 主网, BSC 测试网
		symbol = "BNB"
	case 137, 80001: // Polygon 主网, Mumbai
		symbol = "MATIC"
	default:
		symbol = "Native"
	}

	// 格式化为合适的精度
	floatVal, _ := eth.Float64()
	if floatVal >= 1 {
		return fmt.Sprintf("%.6f %s", floatVal, symbol)
	} else if floatVal >= 0.001 {
		return fmt.Sprintf("%.8f %s", floatVal, symbol)
	} else {
		return fmt.Sprintf("%.12f %s", floatVal, symbol)
	}
}

// formatERC20Token 格式化ERC20代币
func formatERC20Token(amount *big.Int, tokenAddr string, decimals int) string {
	if decimals == 0 {
		decimals = 18 // 默认18位小数
	}

	token := new(big.Float).SetInt(amount)
	token.Quo(token, new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)))

	floatVal, _ := token.Float64()

	// 获取代币符号（简化版本，实际项目中可以从合约查询）
	symbol := getTokenSymbol(tokenAddr)

	if floatVal >= 1 {
		return fmt.Sprintf("%.6f %s", floatVal, symbol)
	} else if floatVal >= 0.001 {
		return fmt.Sprintf("%.8f %s", floatVal, symbol)
	} else {
		return fmt.Sprintf("%.12f %s", floatVal, symbol)
	}
}

// getTokenSymbol 获取代币符号（简化版本）
func getTokenSymbol(tokenAddr string) string {
	// 常见代币地址映射（BSC测试网）
	knownTokens := map[string]string{
		"0x55d398326f99059ff775485246999027b3197955": "USDT",
		"0x8ac76a51cc950d9822d68b83fe1ad97b32cd580d": "USDC",
		"0xe9e7cea3dedca5984780bafc599bd69add087d56": "BUSD",
		"0x2170ed0880ac9a755fd29b2688956bd959f933f8": "ETH",
		"0x7130d2a12b9bcbfae4f2634d864a1ee1ce3ead9c": "BTCB",
	}

	if symbol, exists := knownTokens[strings.ToLower(tokenAddr)]; exists {
		return symbol
	}

	// 返回地址的简短形式
	if len(tokenAddr) >= 10 {
		return fmt.Sprintf("Token(%s...%s)", tokenAddr[:6], tokenAddr[len(tokenAddr)-4:])
	}
	return "Unknown"
}

// enhanceWithLiFiAPI 使用LI.FI API增强交易解析
func (p *LogParser) enhanceWithLiFiAPI(tx *types.Transaction, chainId uint64) {
	// 异步调用LI.FI API，不阻塞主流程
	go func() {
		if len(tx.Data()) == 0 {
			return // 没有调用数据
		}

		// 调用LI.FI解析API
		if lifiData := p.callLiFiParseAPI(chainId, tx.Data()); lifiData != nil {
			// 可以将解析结果存储到缓存中，供后续使用
			p.cacheLiFiResult(tx.Hash().Hex(), lifiData)
		}

		// 调用LI.FI状态API
		if statusData := p.callLiFiStatusAPI(tx.Hash().Hex(), chainId); statusData != nil {
			// 存储状态信息
			p.cacheLiFiStatus(tx.Hash().Hex(), statusData)
		}
	}()
}

// callLiFiParseAPI 调用LI.FI解析API
func (p *LogParser) callLiFiParseAPI(chainId uint64, callData []byte) *LiFiResponse {
	baseURL := "https://li.quest/v1/calldata/parse"
	params := url.Values{}
	params.Add("chainId", fmt.Sprintf("%d", chainId))
	params.Add("callData", fmt.Sprintf("0x%x", callData))

	resp, err := http.Get(baseURL + "?" + params.Encode())
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil
	}

	var lifiResp LiFiResponse
	if err := json.Unmarshal(body, &lifiResp); err != nil {
		return nil
	}

	return &lifiResp
}

// callLiFiStatusAPI 调用LI.FI状态API
func (p *LogParser) callLiFiStatusAPI(txHash string, chainId uint64) *LiFiResponse {
	baseURL := "https://li.quest/v1/status"
	params := url.Values{}
	params.Add("txHash", txHash)
	params.Add("fromChain", fmt.Sprintf("%d", chainId))

	resp, err := http.Get(baseURL + "?" + params.Encode())
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil
	}

	var lifiResp LiFiResponse
	if err := json.Unmarshal(body, &lifiResp); err != nil {
		return nil
	}

	return &lifiResp
}

// cacheLiFiResult 缓存LI.FI解析结果
func (p *LogParser) cacheLiFiResult(txHash string, data *LiFiResponse) {
	// 简单实现，实际可以使用Redis或内存缓存
	// 这里仅作示例
}

// cacheLiFiStatus 缓存LI.FI状态信息
func (p *LogParser) cacheLiFiStatus(txHash string, data *LiFiResponse) {
	// 简单实现，实际可以使用Redis或内存缓存
	// 这里仅作示例
}
