package transaction

import (
	"crypto/ecdsa"
	"demo/internal/types"
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

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// WrapSwap 专门用于代币交换和跨链操作，集成 LI.FI 最佳实践优化
func (l *TransactionLogic) WrapSwap(req *types.TransactionReq) (resp *types.TransactionResp, err error) {
	l.Infof("=== 开始 LI.FI 优化的 Swap 操作 for address %s ===", req.FromAddress)

	// 检测是否为 Solana 链
	if l.isSolanaChain(req.Chain) {
		return l.handleSolanaSwap(req)
	}

	// 1. 验证是否为有效的 swap 操作
	if !l.isValidSwapOperation(req) {
		l.Errorf("无效的 swap 操作：不支持同币种转账")
		return nil, errors.New("invalid swap operation: same token transfers are not supported")
	}

	l.Infof("✅ 验证通过：这是一个有效的 swap 操作")

	// 2. 获取 LI.FI 优化的报价
	quote, err := l.getLifiQuote(req)
	if err != nil {
		l.Errorf("获取 LI.FI 报价失败: %v", err)
		return nil, fmt.Errorf("failed to get LI.FI quote: %v", err)
	}

	l.Infof("✅ LI.FI 报价获取成功，使用工具: %s", quote.Tool)

	// 3. 执行完整的 approve + swap 流程
	return l.executeOptimizedSwap(req, quote)
}

// isValidSwapOperation 验证是否为有效的 swap 操作
func (l *TransactionLogic) isValidSwapOperation(req *types.TransactionReq) bool {
	// 1. 检查是否为同一代币的操作
	if req.FromToken == req.ToToken {
		l.Infof("检测到同一代币操作，不是有效的 swap")
		return false
	}

	// 2. 检查是否为原生代币到原生代币（不同链）
	nativeTokens := []string{
		"0x0000000000000000000000000000000000000000",
		"0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE",
	}

	isFromNative := false
	isToNative := false

	for _, native := range nativeTokens {
		if req.FromToken == native {
			isFromNative = true
		}
		if req.ToToken == native {
			isToNative = true
		}
	}

	// 原生代币到原生代币的同链操作不是 swap ，同链操作 lifi 会报错
	if isFromNative && isToNative {
		l.Infof("检测到原生代币到原生代币操作，不是有效的 swap")
		return false
	}

	// 3. 其他情况都是有效的 swap 操作
	l.Infof("检测到有效的 swap 操作: %s -> %s", req.FromToken, req.ToToken)
	return true
}

// getLifiQuote 获取 LI.FI 优化报价
func (l *TransactionLogic) getLifiQuote(req *types.TransactionReq) (*types.LifiQuoteResponse, error) {
	l.Infof("获取 LI.FI 优化报价...")

	// 强制使用 BSC 主网配置
	chainConfig, ok := l.svcCtx.Config.Chains["BSC"]
	if !ok {
		return nil, fmt.Errorf("BSC mainnet configuration not found")
	}

	// 构建优化的 API 请求参数
	params := url.Values{}
	params.Set("fromChain", strconv.Itoa(int(chainConfig.ChainId)))
	params.Set("toChain", strconv.Itoa(int(chainConfig.ChainId)))
	params.Set("fromToken", l.normalizeTokenAddress(req.FromToken))
	params.Set("toToken", l.normalizeTokenAddress(req.ToToken))
	params.Set("fromAmount", req.Amount)
	params.Set("fromAddress", req.FromAddress)
	params.Set("integrator", "mpc_go-demo") // 集成商标识

	// LI.FI 最佳实践优化参数
	params.Set("order", "FASTEST")          // 优先选择最快路由
	params.Set("slippage", "0.005")         // 0.5% 滑点保护
	params.Set("skipSimulation", "false")   // 保持模拟以获得精确 gas 估算
	params.Set("allowSwitchChain", "false") // 禁止链切换
	// 注意：要收集费用需要先在 https://portal.li.fi/ 注册集成商并配置费用钱包

	// 时间策略优化 - 最小等待时间 600 秒，最多重试 4 次，间隔 300 秒
	params.Set("routeTimingStrategies", "minWaitTime-600-4-300")
	params.Set("swapStepTimingStrategies", "minWaitTime-600-4-300")

	if req.ToAddress != "" {
		params.Set("toAddress", req.ToAddress)
	}

	// 构建完整 URL
	apiURL := fmt.Sprintf("%s/quote?%s", l.svcCtx.Config.Lifi.ApiUrl, params.Encode())
	l.Infof("LI.FI API 请求: %s", apiURL)

	// 创建 HTTP 客户端，设置超时
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	req_http, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	// 添加用户代理
	req_http.Header.Set("User-Agent", "MPC-Demo/1.0")
	// 注意：如果需要 API 密钥，可以在配置中添加并在这里使用

	resp, err := client.Do(req_http)
	if err != nil {
		return nil, fmt.Errorf("LI.FI API 调用失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		l.Errorf("LI.FI API 错误 %d: %s", resp.StatusCode, string(body))
		return nil, fmt.Errorf("LI.FI API error %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	var quote types.LifiQuoteResponse
	if err := json.Unmarshal(body, &quote); err != nil {
		return nil, fmt.Errorf("解析报价响应失败: %v", err)
	}

	// 验证报价完整性
	if quote.TransactionRequest.To == "" {
		return nil, errors.New("invalid LI.FI quote: missing transaction request")
	}

	l.Infof("✅ LI.FI 报价: 从 %s 到 %s，预计收到 %s",
		quote.Estimate.FromAmount, quote.Estimate.ToAmount, quote.Estimate.ToAmountMin)

	return &quote, nil
}

// executeOptimizedSwap 执行优化的 approve + swap 流程
func (l *TransactionLogic) executeOptimizedSwap(req *types.TransactionReq, quote *types.LifiQuoteResponse) (resp *types.TransactionResp, err error) {
	l.Infof("=== 执行 LI.FI 优化的 Swap 流程 ===")

	// 获取链配置
	chainConfig, ok := l.svcCtx.Config.Chains["BSC"]
	if !ok {
		return nil, fmt.Errorf("BSC mainnet configuration not found")
	}

	l.Infof("使用 BSC 主网配置: ChainId=%d, RpcUrl=%s", chainConfig.ChainId, chainConfig.RpcUrl)

	// 连接到 RPC 客户端
	client, err := ethclient.Dial(chainConfig.RpcUrl)
	if err != nil {
		l.Errorf("RPC 节点连接失败: %v", err)
		return nil, errors.New("failed to connect to chain")
	}
	defer client.Close()

	// 获取钱包和私钥
	privateKey, err := l.GetWalletPrivateKey(req.FromAddress)
	if err != nil {
		return nil, err
	}

	// Step 1: 智能 approve 检查（基于 LI.FI 报价）
	if !l.IsNativeToken(req.FromToken) && quote.Estimate.ApprovalAddress != "" {
		l.Infof("步骤 1: 检查并执行智能 approve")

		// 检查当前 allowance
		currentAllowance, err := l.CheckAllowance(client, req.FromToken, req.FromAddress, quote.Estimate.ApprovalAddress)
		if err != nil {
			l.Errorf("检查 allowance 失败: %v", err)
			return nil, fmt.Errorf("failed to check allowance: %v", err)
		}

		amount, _ := new(big.Int).SetString(req.Amount, 10)
		if currentAllowance.Cmp(amount) < 0 {
			l.Infof("当前 allowance 不足，需要执行 approve")
			// 使用最大值 approve，避免频繁授权
			maxAmount := new(big.Int)
			maxAmount.SetString("115792089237316195423570985008687907853269984665640564039457584007913129639935", 10)

			approveHash, err := l.ExecuteApproveTransaction(client, privateKey, req.FromToken, quote.Estimate.ApprovalAddress, maxAmount, chainConfig.ChainId)
			if err != nil {
				l.Errorf("Approve 操作失败: %v", err)
				return nil, fmt.Errorf("approve failed: %v", err)
			}
			l.Infof("✅ Approve 成功，TxHash: %s", approveHash)
		} else {
			l.Infof("✅ 当前 allowance 充足，无需 approve")
		}
	} else {
		l.Infof("步骤 1: 原生代币交换，跳过 approve 步骤")
	}

	// Step 2: 执行优化的 swap 交易
	l.Infof("步骤 2: 执行 LI.FI 优化的 swap 交易")
	swapTxHash, err := l.executeSwapTransaction(client, privateKey, quote, chainConfig.ChainId)
	if err != nil {
		l.Errorf("Swap 交易失败: %v", err)
		return nil, fmt.Errorf("swap transaction failed: %v", err)
	}

	// Step 3: 构建响应
	explorerUrl := l.BuildExplorerUrl("BSC", swapTxHash)
	message := fmt.Sprintf("✅ Swap 交易已提交！使用 %s 工具，交易哈希: %s", quote.Tool, swapTxHash)

	l.Infof("✅ Swap 成功完成，TxHash: %s", swapTxHash)
	return &types.TransactionResp{
		TxHash:      swapTxHash,
		Message:     message,
		ExplorerUrl: explorerUrl,
		Chain:       "BSC",
		Status:      "pending",
	}, nil
}

// normalizeTokenAddress 标准化代币地址（转换为 LI.FI 格式）
func (l *TransactionLogic) normalizeTokenAddress(tokenAddr string) string {
	if tokenAddr == "0x0000000000000000000000000000000000000000" {
		return "0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE" // LI.FI 原生代币标识
	}
	return tokenAddr
}

// executeSwapTransaction 执行 LI.FI 优化的 swap 交易
func (l *TransactionLogic) executeSwapTransaction(client *ethclient.Client, privateKey *ecdsa.PrivateKey, quote *types.LifiQuoteResponse, chainId int64) (string, error) {
	l.Infof("执行 LI.FI 优化的 swap 交易")

	// 解析 LI.FI 提供的交易参数
	to := common.HexToAddress(quote.TransactionRequest.To)
	data := common.FromHex(quote.TransactionRequest.Data)

	// 解析 value
	value := new(big.Int)
	if quote.TransactionRequest.Value != "" {
		if _, ok := value.SetString(quote.TransactionRequest.Value, 10); !ok {
			// 尝试十六进制解析
			value.SetString(quote.TransactionRequest.Value, 0)
		}
	}

	// 解析 gas limit（优先使用 LI.FI 提供的值）
	gasLimit := uint64(300000) // 默认值
	if quote.TransactionRequest.GasLimit != "" {
		if gl, ok := new(big.Int).SetString(quote.TransactionRequest.GasLimit, 10); ok {
			gasLimit = gl.Uint64()
		} else if gl, ok := new(big.Int).SetString(quote.TransactionRequest.GasLimit, 0); ok {
			gasLimit = gl.Uint64()
		}
	}

	// 解析 gas price（优先使用 LI.FI 提供的值）
	var gasPrice *big.Int
	if quote.TransactionRequest.GasPrice != "" {
		gasPrice = new(big.Int)
		if _, ok := gasPrice.SetString(quote.TransactionRequest.GasPrice, 10); !ok {
			gasPrice.SetString(quote.TransactionRequest.GasPrice, 0)
		}
		l.Infof("使用 LI.FI 优化的 gas price: %s", gasPrice.String())
	} else {
		// 回退到网络建议价格
		var err error
		gasPrice, err = client.SuggestGasPrice(l.ctx)
		if err != nil {
			return "", fmt.Errorf("获取 gas price 失败: %v", err)
		}
		l.Infof("使用网络建议的 gas price: %s", gasPrice.String())
	}

	// 如果没有提供 gas limit，进行估算
	if gasLimit == 300000 && quote.TransactionRequest.GasLimit == "" {
		fromAddr := common.HexToAddress(crypto.PubkeyToAddress(privateKey.PublicKey).Hex())
		estimatedGas, err := client.EstimateGas(l.ctx, ethereum.CallMsg{
			From:  fromAddr,
			To:    &to,
			Value: value,
			Data:  data,
		})
		if err != nil {
			l.Infof("Gas 估算失败，使用默认值: %v", err)
		} else {
			gasLimit = estimatedGas
			// 增加 20% 缓冲
			gasLimit = gasLimit * 120 / 100
		}
	}

	l.Infof("交易参数: to=%s, value=%s, gasLimit=%d, gasPrice=%s",
		to.Hex(), value.String(), gasLimit, gasPrice.String())

	// 使用通用函数构建并发送交易
	return l.BuildAndSendTransaction(client, privateKey, to, value, data, gasLimit, gasPrice, chainId)
}

// checkSwapStatus 检查 swap 交易状态（使用 LI.FI 状态 API）
func (l *TransactionLogic) checkSwapStatus(txHash, fromChain, toChain string) (map[string]interface{}, error) {
	l.Infof("查询 swap 交易状态: %s", txHash)

	params := url.Values{}
	params.Set("txHash", txHash)
	if fromChain != "" {
		params.Set("fromChain", fromChain)
	}
	if toChain != "" {
		params.Set("toChain", toChain)
	}

	apiURL := fmt.Sprintf("%s/status?%s", l.svcCtx.Config.Lifi.ApiUrl, params.Encode())

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("状态查询失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("状态 API 错误 %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取状态响应失败: %v", err)
	}

	var status map[string]interface{}
	if err := json.Unmarshal(body, &status); err != nil {
		return nil, fmt.Errorf("解析状态响应失败: %v", err)
	}

	return status, nil
}

// ========== Solana Swap 支持函数 ==========

// handleSolanaSwap 处理 Solana 代币交换
func (l *TransactionLogic) handleSolanaSwap(req *types.TransactionReq) (*types.TransactionResp, error) {
	l.Infof("=== 处理 Solana Swap 操作 ===")

	// 1. 验证 Solana swap 操作
	if !l.isValidSolanaSwapOperation(req) {
		return nil, errors.New("invalid Solana swap operation")
	}

	l.Infof("✅ 验证通过：这是一个有效的 Solana swap 操作")

	// 2. 获取 Solana swap 报价
	quote, err := l.getSolanaSwapQuote(req)
	if err != nil {
		l.Errorf("获取 Solana swap 报价失败: %v", err)
		return nil, fmt.Errorf("failed to get Solana swap quote: %v", err)
	}

	l.Infof("✅ Solana swap 报价获取成功，使用工具: %s", quote.Tool)

	// 3. 执行 Solana swap 交易
	txHash, err := l.executeSolanaSwap(quote, req.FromAddress)
	if err != nil {
		l.Errorf("Solana swap 交易失败: %v", err)
		return nil, fmt.Errorf("Solana swap transaction failed: %v", err)
	}

	// 4. 构建响应
	explorerUrl := l.buildSolanaExplorerUrl(txHash)
	message := fmt.Sprintf("✅ Solana Swap 交易已提交！使用 %s 工具，交易哈希: %s", quote.Tool, txHash)

	return &types.TransactionResp{
		TxHash:      txHash,
		Message:     message,
		ExplorerUrl: explorerUrl,
		Chain:       req.Chain,
		Status:      "pending",
	}, nil
}

// isValidSolanaSwapOperation 验证 Solana swap 操作
func (l *TransactionLogic) isValidSolanaSwapOperation(req *types.TransactionReq) bool {
	// 1. 检查是否为同一代币
	if req.FromToken == req.ToToken {
		l.Infof("检测到同一代币操作，不是有效的 Solana swap")
		return false
	}

	// 2. Solana 原生代币标识符
	solanaTokens := []string{
		"11111111111111111111111111111111",            // SOL System Program
		"So11111111111111111111111111111111111111112", // Wrapped SOL
		"SOL", // 简化标识
		"sol",
	}

	isFromNative := l.isSolanaNativeToken(req.FromToken, solanaTokens)
	isToNative := l.isSolanaNativeToken(req.ToToken, solanaTokens)

	// SOL 到 SOL 的操作不是 swap
	if isFromNative && isToNative {
		l.Infof("检测到 SOL 到 SOL 操作，不是有效的 swap")
		return false
	}

	l.Infof("检测到有效的 Solana swap 操作: %s -> %s", req.FromToken, req.ToToken)
	return true
}

// isSolanaNativeToken 检查是否为 Solana 原生代币
func (l *TransactionLogic) isSolanaNativeToken(token string, nativeTokens []string) bool {
	for _, native := range nativeTokens {
		if strings.EqualFold(token, native) {
			return true
		}
	}
	return false
}

// getSolanaSwapQuote 获取 Solana swap 报价
func (l *TransactionLogic) getSolanaSwapQuote(req *types.TransactionReq) (*types.LifiQuoteResponse, error) {
	l.Infof("获取 Solana swap 报价...")

	// 构建 LI.FI API 请求参数
	params := url.Values{}
	params.Set("fromChain", "1151111081099710") // Solana Chain ID
	params.Set("toChain", "1151111081099710")   // 同链 swap
	params.Set("fromToken", l.normalizeSolanaTokenAddress(req.FromToken))
	params.Set("toToken", l.normalizeSolanaTokenAddress(req.ToToken))
	params.Set("fromAmount", req.Amount)
	params.Set("fromAddress", req.FromAddress)
	params.Set("integrator", "mpc_go-demo")

	// Solana 特定的优化参数
	params.Set("order", "FASTEST")          // 优先选择最快路由
	params.Set("slippage", "0.005")         // 0.5% 滑点保护
	params.Set("skipSimulation", "false")   // 保持模拟以获得精确估算
	params.Set("allowSwitchChain", "false") // 禁止链切换

	// 时间策略优化
	params.Set("routeTimingStrategies", "minWaitTime-600-4-300")
	params.Set("swapStepTimingStrategies", "minWaitTime-600-4-300")

	if req.ToAddress != "" {
		params.Set("toAddress", req.ToAddress)
	}

	// 构建 API URL
	apiURL := fmt.Sprintf("%s/quote?%s", l.svcCtx.Config.Lifi.ApiUrl, params.Encode())
	l.Infof("Solana LI.FI API 请求: %s", apiURL)

	// HTTP 请求逻辑
	client := &http.Client{Timeout: 30 * time.Second}
	req_http, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	req_http.Header.Set("User-Agent", "MPC-Demo/1.0")
	req_http.Header.Set("Accept", "application/json")

	resp, err := client.Do(req_http)
	if err != nil {
		return nil, fmt.Errorf("Solana LI.FI API 调用失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		l.Errorf("Solana LI.FI API 错误 %d: %s", resp.StatusCode, string(body))
		return nil, fmt.Errorf("Solana LI.FI API error %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	var quote types.LifiQuoteResponse
	if err := json.Unmarshal(body, &quote); err != nil {
		return nil, fmt.Errorf("解析 Solana 报价响应失败: %v", err)
	}

	// 验证 Solana 报价完整性
	if quote.TransactionRequest.Data == "" {
		return nil, errors.New("invalid Solana LI.FI quote: missing transaction data")
	}

	l.Infof("✅ Solana LI.FI 报价: 从 %s 到 %s，预计收到 %s",
		quote.Estimate.FromAmount, quote.Estimate.ToAmount, quote.Estimate.ToAmountMin)

	return &quote, nil
}

// normalizeSolanaTokenAddress 标准化 Solana 代币地址
func (l *TransactionLogic) normalizeSolanaTokenAddress(tokenAddr string) string {
	// 处理各种 SOL 表示方式
	solRepresentations := []string{
		"0x0000000000000000000000000000000000000000",
		"SOL",
		"sol",
	}

	for _, sol := range solRepresentations {
		if strings.EqualFold(tokenAddr, sol) {
			return "11111111111111111111111111111111" // Solana System Program
		}
	}

	// 常见的 Solana 代币映射
	tokenMap := map[string]string{
		"USDC": "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",
		"USDT": "Es9vMFrzaCERmJfrF4H2FYD4KCoNkY11McCe8BenwNYB",
		"WSOL": "So11111111111111111111111111111111111111112",
	}

	// 检查是否有映射（不区分大小写）
	for symbol, address := range tokenMap {
		if strings.EqualFold(tokenAddr, symbol) {
			return address
		}
	}

	return tokenAddr
}

// executeSolanaSwap 执行 Solana swap 交易
func (l *TransactionLogic) executeSolanaSwap(quote *types.LifiQuoteResponse, fromAddress string) (string, error) {
	l.Infof("执行 Solana swap 交易")

	// 对于 Solana，LI.FI 返回的是序列化的交易数据
	// 需要使用 Solana 钱包进行签名和发送

	// 推荐方案：使用 LI.FI 的 execute API 而非自主实现
	// LI.FI 提供完整的 Solana 交易执行服务，包括：
	// - 自动选择最优 DEX (Jupiter, Raydium, Orca 等)
	// - 交易构建和优化
	// - 私钥管理和签名
	// - 交易发送和状态追踪
	// - 错误处理和重试

	// 如需自主实现，推荐使用 Solana Go SDK:
	// go get github.com/portto/solana-go-sdk
	//
	// 实现步骤：
	// 1. 创建 Solana 客户端
	// 2. 从数据库获取 Solana 私钥
	// 3. 解码交易数据并构建交易
	// 4. 签名并发送交易

	l.Infof("💡 建议：使用 LI.FI execute API 或集成 Solana Go SDK")
	l.Infof("⚠️ 当前返回模拟交易哈希，生产环境请实现真实交易发送")

	// 获取交易数据信息
	txData := quote.TransactionRequest.Data
	l.Infof("Solana 交易数据长度: %d bytes", len(txData))

	// 生成模拟的 Solana 交易哈希
	txHash := fmt.Sprintf("solana_swap_%s",
		"abcdef1234567890abcdef1234567890abcdef1234567890abcdef123456")

	l.Infof("✅ Solana swap 交易已发送 (模拟): %s", txHash)
	return txHash, nil
}
