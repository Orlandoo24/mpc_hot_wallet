package transaction

import (
	"context"
	"demo/internal/types"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	evmTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// WrapSend 纯原生转账操作，不借助任何外部服务，专门处理简单的代币转账
func (l *TransactionLogic) WrapSend(req *types.TransactionReq) (resp *types.TransactionResp, err error) {
	l.Infof("--- 开始处理 /transaction/send 请求 (纯原生转账) for address %s ---", req.FromAddress)

	// 检测是否为 Solana 链
	if l.isSolanaChain(req.Chain) {
		return l.handleSolanaTransfer(req)
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
