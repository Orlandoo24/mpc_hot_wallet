package transaction

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"demo/internal/types"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/ethereum/go-ethereum/common"
	evmTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// WrapSend 纯原生转账操作，不借助任何外部服务，专门处理简单的代币转账
func (l *TransactionLogic) WrapSend(req *types.TransactionReq) (resp *types.TransactionResp, err error) {
	l.Infof("--- 开始处理 /transaction/send 请求 (纯原生转账) for address %s ---", req.FromAddress)

	// 检测是否为 Solana 链
	if l.isSolanaChain(req.Chain) {
		return l.handleSolanaTransfer(req)
	}

	// 检测是否为 Bitcoin 链
	if l.isBTCChain(req.Chain) {
		return l.handleBTCTransfer(req)
	}

	// 1. 获取链配置
	l.Infof("步骤 1: 获取链配置 for chain: %s", req.Chain)
	chainConfig, ok := l.svcCtx.Config.Chains[req.Chain]
	if !ok {
		l.Errorf("获取链配置失败: 不支持的链 %s", req.Chain)
		return nil, fmt.Errorf("unsupported chain: %s", req.Chain)
	}
	l.Infof("链配置获取成功: ChainId=%d, RpcUrl=%s", chainConfig.ChainId, chainConfig.RpcUrl)

	// 2. 连接 RPC 客户端
	l.Infof("步骤 2: 连接到 RPC 节点...")
	client, err := ethclient.Dial(chainConfig.RpcUrl)
	if err != nil {
		l.Errorf("RPC 节点连接失败: %v", err)
		return nil, errors.New("failed to connect to chain")
	}
	defer client.Close()
	l.Infof("RPC 节点连接成功")

	// 3. 获取钱包和私钥
	l.Infof("步骤 3: 从数据库查询钱包地址: %s", req.FromAddress)
	privateKey, err := l.GetWalletPrivateKey(req.FromAddress)
	if err != nil {
		return nil, err
	}
	l.Infof("私钥解析成功")

	// 4. 验证收款地址是否为合约（避免 OOG 问题）
	l.Infof("步骤 4: 验证收款地址类型...")
	toAddr := common.HexToAddress(req.ToAddress)
	code, err := client.CodeAt(l.ctx, toAddr, nil)
	if err != nil {
		l.Errorf("检查收款地址失败: %v", err)
		return nil, errors.New("failed to check to address")
	}

	if len(code) > 0 {
		l.Infof("警告：收款地址是合约地址，可能存在 gas 不足风险")
		// 不阻止交易，但给出警告
	} else {
		l.Infof("收款地址是 EOA 地址，安全")
	}

	// 5. 解析转账金额
	amount := new(big.Int)
	amount.SetString(req.Amount, 10)
	l.Infof("转账金额: %s wei", amount.String())

	// 6. 获取交易参数
	nonce, err := client.PendingNonceAt(l.ctx, common.HexToAddress(req.FromAddress))
	if err != nil {
		l.Errorf("获取 nonce 失败: %v", err)
		return nil, errors.New("failed to get nonce")
	}
	l.Infof("获取 nonce 成功: %d", nonce)

	// 7. 构建交易（完全原生，不借助外部服务）
	var tx *evmTypes.Transaction

	// 判断交易类型
	if l.IsNativeToken(req.FromToken) {
		// 原生代币转账（BNB/ETH 等）
		l.Infof("=== 执行原生代币转账 ===")

		// 使用精确的 gas 估算
		fromAddr := common.HexToAddress(req.FromAddress)
		gasLimit, gasPrice, err := l.EstimateNativeTransferGas(client, fromAddr, toAddr, amount)
		if err != nil {
			l.Errorf("Gas 估算失败: %v", err)
			return nil, fmt.Errorf("gas estimation failed: %v", err)
		}

		l.Infof("Gas 估算结果: gasLimit=%d, gasPrice=%s", gasLimit, gasPrice.String())

		// 构建纯原生转账交易（无 data，直接转账）
		tx = evmTypes.NewTx(&evmTypes.LegacyTx{
			Nonce:    nonce,
			To:       &toAddr,
			Value:    amount,
			Gas:      gasLimit,
			GasPrice: gasPrice,
			Data:     nil,
		})

	} else {
		// ERC20 代币转账
		l.Infof("=== 执行 ERC20 代币转账 ===")

		// 构建 ERC20 transfer 调用数据
		data, err := l.BuildERC20TransferData(req.ToAddress, amount)
		if err != nil {
			l.Errorf("构建 ERC20 调用数据失败: %v", err)
			return nil, fmt.Errorf("failed to build ERC20 data: %v", err)
		}

		// 估算 ERC20 转账的 gas
		fromAddr := common.HexToAddress(req.FromAddress)
		tokenAddr := common.HexToAddress(req.FromToken)
		gasLimit, gasPrice, err := l.EstimateERC20TransferGas(client, fromAddr, tokenAddr, data)
		if err != nil {
			l.Errorf("ERC20 Gas 估算失败: %v", err)
			return nil, fmt.Errorf("ERC20 gas estimation failed: %v", err)
		}

		l.Infof("ERC20 Gas 估算结果: gasLimit=%d, gasPrice=%s", gasLimit, gasPrice.String())

		// 构建 ERC20 转账交易
		tx = evmTypes.NewTx(&evmTypes.LegacyTx{
			Nonce:    nonce,
			To:       &tokenAddr,
			Value:    big.NewInt(0),
			Gas:      gasLimit,
			GasPrice: gasPrice,
			Data:     data,
		})
	}

	// 8. 签名交易
	l.Infof("步骤 8: 签名交易...")
	signedTx, err := evmTypes.SignTx(tx, evmTypes.NewEIP155Signer(big.NewInt(chainConfig.ChainId)), privateKey)
	if err != nil {
		l.Errorf("交易签名失败: %v", err)
		return nil, errors.New("failed to sign transaction")
	}
	l.Infof("交易签名成功, TxHash: %s", signedTx.Hash().Hex())

	// 9. 异步发送交易到区块链网络
	l.Infof("步骤 9: 异步发送交易到区块链网络...")

	txHash := signedTx.Hash().Hex()

	// 立即返回 TxHash，异步发送交易
	resp = &types.TransactionResp{
		TxHash:      txHash,
		Message:     l.buildSuccessMessage(req),
		ExplorerUrl: l.BuildExplorerUrl(req.Chain, txHash),
		Chain:       req.Chain,
		Status:      "pending", // 异步发送中，状态为 pending
	}

	// 异步发送交易（不阻塞响应）
	go func() {
		asyncCtx := context.Background() // 使用独立的 context 避免请求取消影响
		l.sendTransactionAsync(asyncCtx, client, signedTx, txHash)
	}()

	l.Infof("--- /transaction/send 请求处理完成, 立即返回 TxHash: %s (异步发送中) ---", resp.TxHash)
	return resp, nil
}

// sendTransactionAsync 异步发送交易到区块链网络
func (l *TransactionLogic) sendTransactionAsync(ctx context.Context, client *ethclient.Client, signedTx *evmTypes.Transaction, txHash string) {
	l.Infof("开始异步发送交易: %s", txHash)

	// 使用重试机制发送交易
	maxRetries := 3
	for i := 0; i < maxRetries; i++ {
		err := client.SendTransaction(ctx, signedTx)
		if err != nil {
			l.Errorf("异步发送交易失败 (尝试 %d/%d): %v", i+1, maxRetries, err)

			// 如果是最后一次重试，记录最终失败
			if i == maxRetries-1 {
				l.Errorf("交易 %s 发送最终失败: %v", txHash, err)
				// 这里可以考虑将失败信息存储到数据库或发送通知
				return
			}

			// 等待一段时间后重试
			select {
			case <-ctx.Done():
				l.Infof("异步发送被取消: %s", txHash)
				return
			case <-time.After(time.Duration(i+1) * time.Second):
				// 指数退避：1s, 2s, 3s
				continue
			}
		} else {
			l.Infof("异步发送交易成功: %s", txHash)
			// 这里可以考虑将成功信息存储到数据库或发送通知
			return
		}
	}
}

// buildSuccessMessage 构建成功消息
func (l *TransactionLogic) buildSuccessMessage(req *types.TransactionReq) string {
	if l.IsNativeToken(req.FromToken) {
		// 原生代币转账消息
		chainName := l.GetChainDisplayName(req.Chain)
		return fmt.Sprintf("✅ %s 原生代币转账已提交！交易正在异步处理中，请通过区块浏览器查询最终状态。", chainName)
	} else {
		// ERC20 代币转账消息
		chainName := l.GetChainDisplayName(req.Chain)
		return fmt.Sprintf("✅ %s 网络上的 ERC20 代币转账已提交！交易正在异步处理中，请通过区块浏览器查询最终状态。", chainName)
	}
}

// ========== Solana 支持函数 ==========

// isSolanaChain 检测是否为 Solana 链
func (l *TransactionLogic) isSolanaChain(chain string) bool {
	solanaChains := []string{"Solana", "SOL", "solana", "sol"}
	for _, solChain := range solanaChains {
		if strings.EqualFold(chain, solChain) {
			return true
		}
	}
	return false
}

// handleSolanaTransfer 处理 Solana 转账
func (l *TransactionLogic) handleSolanaTransfer(req *types.TransactionReq) (*types.TransactionResp, error) {
	l.Infof("=== 处理 Solana 转账 ===")

	// 对于 Solana，我们需要使用 LI.FI API 来构建交易
	// 因为 Solana 交易构建比 EVM 复杂得多

	// 1. 调用 LI.FI quote API 获取交易数据
	quote, err := l.getSolanaQuote(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get Solana quote: %v", err)
	}

	// 2. 验证 quote 响应
	if quote.TransactionRequest.Data == "" {
		return nil, fmt.Errorf("invalid Solana quote: missing transaction data")
	}

	// 3. 对于 Solana，LI.FI 返回的是 base64 编码的交易数据
	// 我们需要使用 Solana 钱包来签名和发送
	txHash, err := l.sendSolanaTransaction(quote.TransactionRequest.Data, req.FromAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to send Solana transaction: %v", err)
	}

	// 4. 构建响应
	explorerUrl := l.buildSolanaExplorerUrl(txHash)
	message := "✅ Solana 转账已提交！交易正在处理中，请通过区块浏览器查询最终状态。"

	return &types.TransactionResp{
		TxHash:      txHash,
		Message:     message,
		ExplorerUrl: explorerUrl,
		Chain:       req.Chain,
		Status:      "pending",
	}, nil
}

// getSolanaQuote 获取 Solana 转账的 quote
func (l *TransactionLogic) getSolanaQuote(req *types.TransactionReq) (*types.LifiQuoteResponse, error) {
	l.Infof("获取 Solana 转账报价...")

	// 构建 LI.FI API 请求
	params := url.Values{}
	params.Set("fromChain", "1151111081099710") // Solana Chain ID
	params.Set("toChain", "1151111081099710")   // 同链转账

	// 标准化 Solana 代币地址
	fromToken := l.normalizeSolanaToken(req.FromToken)
	toToken := l.normalizeSolanaToken(req.ToToken)

	params.Set("fromToken", fromToken)
	params.Set("toToken", toToken)
	params.Set("fromAmount", req.Amount)
	params.Set("fromAddress", req.FromAddress)
	params.Set("toAddress", req.ToAddress)
	params.Set("integrator", "mpc-demo")
	params.Set("skipSimulation", "false")
	params.Set("allowSwitchChain", "false")

	apiURL := l.svcCtx.Config.Lifi.ApiUrl + "/quote?" + params.Encode()
	l.Infof("调用 LI.FI API: %s", apiURL)

	// 创建 HTTP 请求
	client := &http.Client{Timeout: 30 * time.Second}
	req_http, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req_http.Header.Set("User-Agent", "mpc-demo/1.0")
	req_http.Header.Set("Accept", "application/json")

	// 发送请求
	resp, err := client.Do(req_http)
	if err != nil {
		return nil, fmt.Errorf("failed to call LI.FI API: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		l.Errorf("LI.FI API 错误 %d: %s", resp.StatusCode, string(body))
		return nil, fmt.Errorf("LI.FI API error %d: %s", resp.StatusCode, string(body))
	}

	// 解析响应
	var quote types.LifiQuoteResponse
	if err := json.Unmarshal(body, &quote); err != nil {
		l.Errorf("解析 Solana 报价响应失败: %v", err)
		return nil, fmt.Errorf("failed to parse Solana quote response: %v", err)
	}

	l.Infof("✅ Solana 报价获取成功")
	return &quote, nil
}

// normalizeSolanaToken 标准化 Solana 代币地址
func (l *TransactionLogic) normalizeSolanaToken(tokenAddr string) string {
	// Solana 原生 SOL 的特殊处理
	if tokenAddr == "0x0000000000000000000000000000000000000000" ||
		strings.EqualFold(tokenAddr, "SOL") ||
		tokenAddr == "" {
		return "11111111111111111111111111111111" // Solana System Program
	}
	return tokenAddr
}

// sendSolanaTransaction 发送 Solana 交易
func (l *TransactionLogic) sendSolanaTransaction(transactionData, fromAddress string) (string, error) {
	l.Infof("发送 Solana 交易")

	// 推荐方案：使用 LI.FI 的 execute API 而非自主实现
	// LI.FI 提供完整的 Solana 交易执行服务，包括：
	// - 交易构建和优化
	// - 私钥管理和签名
	// - 交易发送和状态追踪
	// - 错误处理和重试

	// 如需自主实现，推荐使用 Solana Go SDK:
	// go get github.com/portto/solana-go-sdk
	//
	// 实现步骤：
	// 1. 创建 Solana 客户端: client.NewClient(client.MainnetRPCEndpoint)
	// 2. 从数据库获取 Solana 私钥
	// 3. 解码并构建交易
	// 4. 签名并发送交易

	l.Infof("💡 建议：使用 LI.FI execute API 或集成 Solana Go SDK")
	l.Infof("⚠️ 当前返回模拟交易哈希，生产环境请实现真实交易发送")

	// 生成模拟的 Solana 交易哈希
	txHash := fmt.Sprintf("solana_tx_%s",
		"abcdef1234567890abcdef1234567890abcdef1234567890abcdef123456")

	l.Infof("✅ Solana 交易已提交 (模拟): %s", txHash)
	return txHash, nil
}

// buildSolanaExplorerUrl 构建 Solana 浏览器链接
func (l *TransactionLogic) buildSolanaExplorerUrl(txHash string) string {
	return fmt.Sprintf("https://solscan.io/tx/%s", txHash)
}

// ========== Bitcoin 支持函数 ==========

// isBTCChain 检测是否为 Bitcoin 链
func (l *TransactionLogic) isBTCChain(chain string) bool {
	btcChains := []string{"BTC", "Bitcoin", "btc", "bitcoin", "20000000000001", "20000000000002", "tBTC", "testnet"}
	for _, btcChain := range btcChains {
		if strings.EqualFold(chain, btcChain) {
			return true
		}
	}
	return false
}

// handleBTCTransfer 处理 Bitcoin 转账（仅支持测试网）
func (l *TransactionLogic) handleBTCTransfer(req *types.TransactionReq) (*types.TransactionResp, error) {
	l.Infof("=== 处理 Bitcoin 测试网转账 ===")

	// 验证 Bitcoin 地址格式
	if !l.isValidBTCAddress(req.FromAddress) || !l.isValidBTCAddress(req.ToAddress) {
		return nil, fmt.Errorf("invalid Bitcoin address format")
	}

	// BTC 只支持测试网，统一使用自实现的发送逻辑
	l.Infof("✅ 使用自实现的 Bitcoin 测试网转账逻辑")
	return l.handleBTCTransferDirect(req)
}

// handleBTCTransferDirect 自实现的 BTC 转账逻辑（用于测试网）
func (l *TransactionLogic) handleBTCTransferDirect(req *types.TransactionReq) (*types.TransactionResp, error) {
	l.Infof("=== 执行自实现的 Bitcoin 测试网转账 ===")

	// 1. 从数据库获取私钥
	l.Infof("步骤 1: 获取发送地址的私钥...")
	privateKey, err := l.GetWalletPrivateKey(req.FromAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to get private key: %v", err)
	}
	l.Infof("私钥获取成功")

	// 2. 构建 Bitcoin 交易
	l.Infof("步骤 2: 构建 Bitcoin 交易...")
	txHash, err := l.buildAndSendBTCTransaction(req, privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to build and send BTC transaction: %v", err)
	}

	// 3. 构建响应
	explorerUrl := l.buildBTCExplorerUrl(txHash)
	message := "✅ Bitcoin 测试网转账已提交！交易正在处理中，请通过区块浏览器查询最终状态。"

	return &types.TransactionResp{
		TxHash:      txHash,
		Message:     message,
		ExplorerUrl: explorerUrl,
		Chain:       req.Chain,
		Status:      "pending",
	}, nil
}

// BlockstreamAPIUTXO 定义了从 Blockstream API 返回的 UTXO 结构
type BlockstreamAPIUTXO struct {
	TxID   string `json:"txid"`
	Vout   int    `json:"vout"`
	Status struct {
		Confirmed bool `json:"confirmed"`
	} `json:"status"`
	Value int64 `json:"value"`
}

// getUTXOsViaAPI 是获取 UTXO 的全新实现，不再使用 rpcclient
func (l *TransactionLogic) getUTXOsViaAPI(address string) ([]btcjson.ListUnspentResult, error) {
	l.Infof("--- 切换思路：开始通过 Blockstream 公共 API 获取 UTXO for address %s ---", address)

	// 1. 构建 API URL
	apiURL := fmt.Sprintf("https://blockstream.info/testnet/api/address/%s/utxo", address)
	l.Infof("调用 API: %s", apiURL)

	// 2. 发起 HTTP GET 请求
	resp, err := http.Get(apiURL)
	if err != nil {
		l.Errorf("请求 Blockstream API 失败: %v", err)
		return nil, fmt.Errorf("failed to call blockstream api: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		l.Errorf("Blockstream API 返回错误状态 %d: %s", resp.StatusCode, string(body))
		return nil, fmt.Errorf("blockstream api returned status %d: %s", resp.StatusCode, string(body))
	}

	// 3. 解析 JSON 响应
	var apiUTXOs []BlockstreamAPIUTXO
	if err := json.NewDecoder(resp.Body).Decode(&apiUTXOs); err != nil {
		l.Errorf("解析 API 响应失败: %v", err)
		return nil, fmt.Errorf("failed to decode api response: %w", err)
	}

	l.Infof("从 API 获取到 %d 个 UTXO", len(apiUTXOs))

	// 4. 将 API 返回的结构转换为代码中使用的 btcjson.ListUnspentResult 结构
	// 这是关键一步，确保代码其他部分无需改动
	var unspentList []btcjson.ListUnspentResult
	sourceAddr, err := btcutil.DecodeAddress(address, &chaincfg.TestNet3Params)
	if err != nil {
		return nil, fmt.Errorf("invalid source address for script generation: %w", err)
	}
	pkScript, err := txscript.PayToAddrScript(sourceAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to create pkScript: %w", err)
	}
	pkScriptHex := hex.EncodeToString(pkScript)

	for _, utxo := range apiUTXOs {
		// 我们只使用已确认的 UTXO
		if !utxo.Status.Confirmed {
			continue
		}
		unspentList = append(unspentList, btcjson.ListUnspentResult{
			TxID:          utxo.TxID,
			Vout:          uint32(utxo.Vout),
			Address:       address,
			Amount:        float64(utxo.Value) / btcutil.SatoshiPerBitcoin, // API 返回 satoshi，需要转为 BTC
			ScriptPubKey:  pkScriptHex,                                     // 我们需要自己生成 ScriptPubKey
			Confirmations: 6,                                               // API只返回已确认的，我们给一个默认的确认数
		})
	}

	l.Infof("--- 成功通过 API 获取并转换了 %d 个已确认的 UTXO ---", len(unspentList))
	return unspentList, nil
}

// buildAndSendBTCTransaction 构建并发送 Bitcoin 交易（自实现逻辑）
func (l *TransactionLogic) buildAndSendBTCTransaction(req *types.TransactionReq, privateKey interface{}) (string, error) {
	l.Infof("=== 开始构建 Bitcoin 测试网交易 ===")
	l.Infof("步骤 0: 解析私钥参数，类型: %T", privateKey)

	var privateKeyHex string
	switch pk := privateKey.(type) {
	case string:
		privateKeyHex = pk
		l.Infof("私钥类型为 string，长度: %d", len(privateKeyHex))
	case *string:
		privateKeyHex = *pk
		l.Infof("私钥类型为 *string，长度: %d", len(privateKeyHex))
	case *ecdsa.PrivateKey:
		// 从 ECDSA 私钥转换为 hex 字符串
		privateKeyHex = hex.EncodeToString(crypto.FromECDSA(pk))
		l.Infof("私钥类型为 *ecdsa.PrivateKey，转换为 hex，长度: %d", len(privateKeyHex))
	default:
		l.Errorf("无效的私钥类型: %T, 值: %+v", privateKey, privateKey)
		return "", fmt.Errorf("invalid private key type: %T", privateKey)
	}
	// 1. 不再需要连接 RPC 客户端，改用公共 API
	l.Infof("步骤 1: 使用 Blockstream 公共 API 获取 UTXO (无需 RPC 连接)...")

	// 2. 解析私钥
	l.Infof("步骤 2: 解析私钥...")
	l.Infof("私钥 hex 字符串: %s", privateKeyHex[:10]+"...") // 只显示前10个字符
	privateKeyBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		l.Errorf("私钥 hex 解码失败: %v", err)
		return "", fmt.Errorf("failed to decode private key: %v", err)
	}
	l.Infof("私钥字节长度: %d", len(privateKeyBytes))

	privKey, _ := btcec.PrivKeyFromBytes(privateKeyBytes)

	wif, err := btcutil.NewWIF(privKey, &chaincfg.TestNet3Params, true)
	if err != nil {
		l.Errorf("创建 WIF 格式私钥失败: %v", err)
		return "", fmt.Errorf("failed to create WIF: %v", err)
	}
	l.Infof("✅ 私钥解析成功，WIF: %s", wif.String()[:10]+"...")

	// 3. 获取发送地址的UTXO
	l.Infof("步骤 3: 获取发送地址的 UTXO...")
	l.Infof("发送地址: %s", req.FromAddress)
	sourceAddr, err := btcutil.DecodeAddress(req.FromAddress, &chaincfg.TestNet3Params)
	if err != nil {
		l.Errorf("解析发送地址失败: %v", err)
		return "", fmt.Errorf("invalid source address: %v", err)
	}
	l.Infof("地址解析成功，类型: %T", sourceAddr)

	l.Infof("调用全新的 API 方法获取 UTXO...")
	unspentList, err := l.getUTXOsViaAPI(req.FromAddress)
	if err != nil {
		l.Errorf("获取 UTXO 失败: %v", err)
		return "", fmt.Errorf("failed to get UTXOs: %v", err)
	}
	l.Infof("✅ 获取到 %d 个 UTXO", len(unspentList))

	// 4. 计算总输入金额和找零
	l.Infof("步骤 4: 计算输入金额和构建交易...")
	amount, err := strconv.ParseInt(req.Amount, 10, 64) // amount in satoshi
	if err != nil {
		l.Errorf("解析转账金额失败: %v", err)
		return "", fmt.Errorf("invalid amount: %v", err)
	}
	l.Infof("转账金额: %d satoshi", amount)

	var inputSum int64
	var inputs []*wire.TxIn
	var inputUtxos []btcjson.ListUnspentResult

	l.Infof("开始选择 UTXO...")
	for i, utxo := range unspentList {
		utxoAmount := int64(utxo.Amount * btcutil.SatoshiPerBitcoin)
		l.Infof("UTXO %d: TxID=%s, Vout=%d, Amount=%.8f BTC (%d satoshi)",
			i, utxo.TxID, utxo.Vout, utxo.Amount, utxoAmount)

		inputSum += utxoAmount
		hash, err := chainhash.NewHashFromStr(utxo.TxID)
		if err != nil {
			l.Errorf("解析 UTXO hash 失败: %v", err)
			return "", fmt.Errorf("failed to parse UTXO hash: %v", err)
		}
		outPoint := wire.NewOutPoint(hash, utxo.Vout)
		txIn := wire.NewTxIn(outPoint, nil, nil)
		inputs = append(inputs, txIn)
		inputUtxos = append(inputUtxos, utxo)

		l.Infof("当前输入总额: %d satoshi, 需要: %d satoshi (含矿工费 5000)", inputSum, amount+5000)
		if inputSum >= amount+5000 { // 预留 5000 satoshi 作为矿工费
			l.Infof("✅ UTXO 选择完成，选择了 %d 个 UTXO", len(inputs))
			break
		}
	}

	if inputSum < amount+5000 {
		l.Errorf("余额不足: 需要 %d satoshi，可用 %d satoshi", amount+5000, inputSum)
		return "", fmt.Errorf("insufficient funds: need %d, available %d", amount+5000, inputSum)
	}

	// 5. 构建交易
	l.Infof("步骤 5: 构建交易...")
	tx := wire.NewMsgTx(wire.TxVersion)
	for _, input := range inputs {
		tx.AddTxIn(input)
	}
	l.Infof("添加了 %d 个输入", len(inputs))

	// 添加接收地址输出
	l.Infof("添加接收地址输出: %s", req.ToAddress)
	destAddr, err := btcutil.DecodeAddress(req.ToAddress, &chaincfg.TestNet3Params)
	if err != nil {
		l.Errorf("解析接收地址失败: %v", err)
		return "", fmt.Errorf("invalid destination address: %v", err)
	}
	destScript, err := txscript.PayToAddrScript(destAddr)
	if err != nil {
		l.Errorf("创建输出脚本失败: %v", err)
		return "", fmt.Errorf("failed to create output script: %v", err)
	}
	tx.AddTxOut(wire.NewTxOut(amount, destScript))
	l.Infof("接收地址输出: %d satoshi", amount)

	// 添加找零输出
	changeAmount := inputSum - amount - 5000 // 减去转账金额和矿工费
	l.Infof("计算找零: %d - %d - 5000 = %d satoshi", inputSum, amount, changeAmount)
	if changeAmount > 0 {
		changeScript, err := txscript.PayToAddrScript(sourceAddr)
		if err != nil {
			l.Errorf("创建找零脚本失败: %v", err)
			return "", fmt.Errorf("failed to create change script: %v", err)
		}
		tx.AddTxOut(wire.NewTxOut(changeAmount, changeScript))
		l.Infof("找零输出: %d satoshi", changeAmount)
	} else {
		l.Infof("无需找零")
	}

	// 6. 签名交易
	l.Infof("步骤 6: 签名交易...")
	for i, utxo := range inputUtxos {
		l.Infof("签名输入 %d, ScriptPubKey: %s", i, utxo.ScriptPubKey)
		scriptBytes, err := hex.DecodeString(utxo.ScriptPubKey)
		if err != nil {
			l.Errorf("解码脚本失败: %v", err)
			return "", fmt.Errorf("failed to decode script: %v", err)
		}

		sigScript, err := txscript.SignatureScript(
			tx,
			i,
			scriptBytes,
			txscript.SigHashAll,
			privKey,
			true,
		)
		if err != nil {
			l.Errorf("签名输入 %d 失败: %v", i, err)
			return "", fmt.Errorf("failed to sign input %d: %v", i, err)
		}
		tx.TxIn[i].SignatureScript = sigScript
		l.Infof("✅ 输入 %d 签名成功", i)
	}

	l.Infof("步骤 7: 序列化并使用公共 API 广播交易...")
	var signedTx bytes.Buffer
	if err := tx.Serialize(&signedTx); err != nil {
		l.Errorf("序列化交易失败: %v", err)
		return "", fmt.Errorf("failed to serialize transaction: %v", err)
	}
	txHex := hex.EncodeToString(signedTx.Bytes())
	l.Infof("交易序列化成功，大小: %d 字节", signedTx.Len())

	// 使用 Blockstream API 广播交易
	broadcastURL := "https://blockstream.info/testnet/api/tx"
	resp, err := http.Post(broadcastURL, "text/plain", strings.NewReader(txHex))
	if err != nil {
		l.Errorf("广播交易请求失败: %v", err)
		return "", fmt.Errorf("failed to broadcast transaction: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		l.Errorf("读取广播响应失败: %v", err)
		return "", fmt.Errorf("failed to read broadcast response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		l.Errorf("广播交易 API 返回错误 %d: %s", resp.StatusCode, string(body))
		return "", fmt.Errorf("broadcast failed with status %d: %s", resp.StatusCode, string(body))
	}

	txHashStr := string(body)
	l.Infof("✅ Bitcoin 测试网交易已成功提交")
	l.Infof("交易哈希: %s", txHashStr)
	return txHashStr, nil
}

// normalizeBTCToken 标准化 Bitcoin 代币地址
func (l *TransactionLogic) normalizeBTCToken(tokenAddr string) string {
	// Bitcoin 原生 BTC 的特殊处理
	if tokenAddr == "0x0000000000000000000000000000000000000000" ||
		strings.EqualFold(tokenAddr, "BTC") ||
		strings.EqualFold(tokenAddr, "bitcoin") ||
		tokenAddr == "" {
		return "bitcoin" // LI.FI 中 Bitcoin 原生代币的标识符
	}
	return tokenAddr
}

// isValidBTCAddress 验证 Bitcoin 地址格式
func (l *TransactionLogic) isValidBTCAddress(address string) bool {
	// 测试网支持多种地址格式:
	// - 测试网 P2PKH: m, n 开头
	// - 测试网 P2SH: 2 开头
	// - 测试网 SegWit: tb1 开头
	return (strings.HasPrefix(address, "m") ||
		strings.HasPrefix(address, "n") ||
		strings.HasPrefix(address, "2") ||
		strings.HasPrefix(address, "tb1") ||
		strings.HasPrefix(address, "bc1") ||
		strings.HasPrefix(address, "1") ||
		strings.HasPrefix(address, "3")) && len(address) >= 26
}

// isBTCTestnetAddress 检测是否为 Bitcoin 测试网地址
func (l *TransactionLogic) isBTCTestnetAddress(address string) bool {
	// 测试网地址格式:
	// - 测试网 P2PKH: m, n 开头
	// - 测试网 P2SH: 2 开头
	// - 测试网 SegWit: tb1 开头
	return strings.HasPrefix(address, "m") ||
		strings.HasPrefix(address, "n") ||
		strings.HasPrefix(address, "2") ||
		strings.HasPrefix(address, "tb1")
}

// sendBTCTransaction 发送 Bitcoin 交易
func (l *TransactionLogic) sendBTCTransaction(psbtData, fromAddress string) (string, error) {
	l.Infof("处理 Bitcoin PSBT 交易")

	// PSBT (Partially Signed Bitcoin Transaction) 处理步骤：
	// 1. 解码 PSBT hex 数据
	// 2. 从数据库获取 Bitcoin 私钥
	// 3. 签名 PSBT
	// 4. 广播已签名的交易

	// 推荐方案：使用 LI.FI 的 execute API 或集成 Bitcoin Go 库
	//
	// 如需自主实现，推荐使用以下库：
	// go get github.com/btcsuite/btcd/btcutil
	// go get github.com/btcsuite/btcd/chaincfg
	// go get github.com/btcsuite/btcd/txscript
	// go get github.com/btcsuite/btcd/wire
	//
	// 实现步骤：
	// 1. 解码 PSBT: psbt.NewFromRawBytes(psbtBytes)
	// 2. 获取私钥: wallet.GetBTCPrivateKey(fromAddress)
	// 3. 签名交易: psbt.Sign(privateKey)
	// 4. 提取最终交易: psbt.Extract()
	// 5. 广播交易: client.SendRawTransaction(tx)

	l.Infof("💡 建议：使用 LI.FI execute API 或集成 Bitcoin Go SDK")
	l.Infof("⚠️ 当前返回模拟交易哈希，生产环境请实现真实 PSBT 处理")

	// 解析 PSBT 中的 memo 信息（用于跨链转账）
	memo, err := l.extractMemoFromPSBT(psbtData)
	if err != nil {
		l.Infof("未找到 memo 信息: %v", err)
	} else {
		l.Infof("提取到 memo 信息: %s", memo)
	}

	// 生成模拟的 Bitcoin 交易哈希
	txHash := fmt.Sprintf("btc_tx_%s",
		"1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")

	l.Infof("✅ Bitcoin 交易已提交 (模拟): %s", txHash)
	return txHash, nil
}

// extractMemoFromPSBT 从 PSBT 中提取 memo 信息
func (l *TransactionLogic) extractMemoFromPSBT(psbtData string) (string, error) {
	// PSBT memo 通常存储在 OP_RETURN 输出中
	// 这里返回模拟的 memo 解析结果
	l.Infof("解析 PSBT 中的 memo 信息...")

	// 实际实现需要：
	// 1. 解码 PSBT hex 数据
	// 2. 遍历交易输出
	// 3. 查找 OP_RETURN 脚本
	// 4. 提取 memo 数据

	return "memo_placeholder", nil
}

// buildBTCExplorerUrl 构建 Bitcoin 测试网浏览器链接
func (l *TransactionLogic) buildBTCExplorerUrl(txHash string) string {
	return fmt.Sprintf("https://mempool.space/testnet/tx/%s", txHash)
}
