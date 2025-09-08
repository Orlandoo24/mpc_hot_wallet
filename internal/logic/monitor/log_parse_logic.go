package monitor

import (
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

// 事件签名常量
var (
	// Transfer(address,address,uint256)
	TransferEventSignature = crypto.Keccak256Hash([]byte("Transfer(address,address,uint256)"))
	// Approval(address,address,uint256)
	ApprovalEventSignature = crypto.Keccak256Hash([]byte("Approval(address,address,uint256)"))
	// Swap 事件签名 (PancakeSwap V2)
	SwapEventSignature = crypto.Keccak256Hash([]byte("Swap(address,uint256,uint256,uint256,uint256,address)"))
	// PancakeSwap V3 Swap
	SwapV3EventSignature = crypto.Keccak256Hash([]byte("Swap(address,address,int256,int256,uint160,uint128,int24)"))
	// Bridge 相关事件 (常见跨链桥)
	BridgeEventSignature = crypto.Keccak256Hash([]byte("Bridge(address,address,uint256,uint256,bytes)"))
	// Deposit/Withdraw (Wrapped tokens)
	DepositEventSignature  = crypto.Keccak256Hash([]byte("Deposit(address,uint256)"))
	WithdrawEventSignature = crypto.Keccak256Hash([]byte("Withdrawal(address,uint256)"))
)

// 已知合约地址 (可扩展)
var (
	// PancakeSwap V2 Router
	PancakeV2Router = common.HexToAddress("0x10ED43C718714eb63d5aA57B78B54704E256024E")
	// PancakeSwap V3 Router
	PancakeV3Router = common.HexToAddress("0x13f4EA83D0bd40E75C8222255bc855a974568Dd4")
	// BSC Bridge 合约地址 (示例)
	BSCBridge = common.HexToAddress("0x0000000000000000000000000000000000000000")
)

// LogParser 日志解析器
type LogParser struct {
	knownContracts map[common.Address]string // 已知合约地址映射
}

// NewLogParser 创建日志解析器
func NewLogParser() *LogParser {
	parser := &LogParser{
		knownContracts: make(map[common.Address]string),
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

	// 如果没有解析到任何事件，但是有BNB转账，创建一个基础事件
	if len(events) == 0 && tx.Value().Cmp(big.NewInt(0)) > 0 {
		event := &TokenEvent{
			BlockNumber: blockNumber,
			TxHash:      tx.Hash().Hex(),
			Timestamp:   timestamp,
			EventType:   "send",
			FromAddr:    fromAddr,
			ToAddr:      toAddr,
			TokenAddr:   "0x0000000000000000000000000000000000000000", // BNB
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

	// 判断事件类型
	eventType := "send"
	if p.isReceiveEvent(from, to, txFromAddr) {
		eventType = "receive"
	}

	return &TokenEvent{
		BlockNumber: blockNumber,
		TxHash:      tx.Hash().Hex(),
		Timestamp:   timestamp,
		EventType:   eventType,
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

	// 判断是授权还是撤销
	eventType := "approve"
	if amount.Cmp(big.NewInt(0)) == 0 {
		eventType = "revoke"
	}

	return &TokenEvent{
		BlockNumber: blockNumber,
		TxHash:      tx.Hash().Hex(),
		Timestamp:   timestamp,
		EventType:   eventType,
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
		EventType:   "swap",
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

	return &TokenEvent{
		BlockNumber: blockNumber,
		TxHash:      tx.Hash().Hex(),
		Timestamp:   timestamp,
		EventType:   "bridge",
		FromAddr:    from,
		ToAddr:      to,
		TokenAddr:   vLog.Address.Hex(),
		Amount:      amount.String(),
		ChainId:     chainId,
	}
}

// isReceiveEvent 判断是否为接收事件
func (p *LogParser) isReceiveEvent(from, to common.Address, txFromAddr string) bool {
	// 简单判断：如果transfer的to地址是交易发起者，则为receive
	txFrom := common.HexToAddress(txFromAddr)
	return to == txFrom
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
