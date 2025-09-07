package logic

import (
	"context"
	"crypto/ecdsa"
	"demo/internal/svc"
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
	evmTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/zeromicro/go-zero/core/logx"
)

// BridgeLogic 跨链逻辑结构体
type BridgeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

// NewBridgeLogic 创建跨链逻辑实例
func NewBridgeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BridgeLogic {
	return &BridgeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetBridgeQuote 获取跨链报价（按照 LI.FI 最佳实践）
func (l *BridgeLogic) GetBridgeQuote(req *types.BridgeQuoteReq) (*types.BridgeQuoteResp, error) {
	l.Infof("--- 开始获取跨链报价 fromChain=%d toChain=%d ---", req.FromChain, req.ToChain)

	// 验证同链同币情况（LI.FI 不支持同链同币的 quote）
	if req.FromChain == req.ToChain && req.FromToken == req.ToToken {
		l.Errorf("同链同币转账不支持使用 LI.FI quote，请使用普通转账")
		return nil, errors.New("same chain same token transfer not supported by LI.FI")
	}

	// 构建 LI.FI API 请求参数
	params := url.Values{}
	params.Set("fromChain", strconv.Itoa(req.FromChain))
	params.Set("toChain", strconv.Itoa(req.ToChain))
	params.Set("fromToken", l.normalizeTokenAddress(req.FromToken))
	params.Set("toToken", l.normalizeTokenAddress(req.ToToken))
	params.Set("fromAmount", req.FromAmount)
	params.Set("fromAddress", req.FromAddress)
	params.Set("toAddress", req.ToAddress)
	params.Set("integrator", "mpc-demo") // 集成商标识

	// LI.FI 最佳实践优化参数
	if req.Order != "" {
		params.Set("order", req.Order)
	} else {
		params.Set("order", "FASTEST") // 优先选择最快路由
	}

	if req.Slippage != "" {
		params.Set("slippage", req.Slippage)
	} else {
		params.Set("slippage", "0.005") // 0.5% 滑点保护
	}

	// 添加 LI.FI 最佳实践参数
	params.Set("skipSimulation", "false")   // 保持模拟以获得精确 gas 估算
	params.Set("allowSwitchChain", "false") // 禁止链切换
	// 注意：要收集费用需要先在 https://portal.li.fi/ 注册集成商并配置费用钱包

	// 时间策略优化 - 最小等待时间 600 秒，最多重试 4 次，间隔 300 秒
	params.Set("routeTimingStrategies", "minWaitTime-600-4-300")
	params.Set("bridgeStepTimingStrategies", "minWaitTime-600-4-300")

	// 调用 LI.FI API
	apiURL := l.svcCtx.Config.Lifi.ApiUrl + "/quote?" + params.Encode()
	l.Infof("调用 LI.FI API: %s", apiURL)

	// 创建 HTTP 客户端，设置超时
	client := &http.Client{Timeout: 30 * time.Second}

	req_http, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	// 添加用户代理
	req_http.Header.Set("User-Agent", "MPC-Demo/1.0")
	// 注意：如果需要 API 密钥，可以在配置中添加并在这里使用

	resp, err := client.Do(req_http)
	if err != nil {
		l.Errorf("LI.FI API 调用失败: %v", err)
		return nil, fmt.Errorf("LI.FI API 调用失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		l.Errorf("LI.FI API 返回错误: %d, %s", resp.StatusCode, string(body))
		return nil, fmt.Errorf("LI.FI API error: %d", resp.StatusCode)
	}

	// 解析响应
	var lifiResp struct {
		TransactionRequest types.BridgeTxRequest `json:"transactionRequest"`
		Estimate           types.BridgeEstimate  `json:"estimate"`
		Tool               string                `json:"tool,omitempty"`        // 使用的桥接工具
		FromChainId        int                   `json:"fromChainId,omitempty"` // 源链ID
		ToChainId          int                   `json:"toChainId,omitempty"`   // 目标链ID
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		l.Errorf("读取响应失败: %v", err)
		return nil, errors.New("failed to read response")
	}

	if err := json.Unmarshal(body, &lifiResp); err != nil {
		l.Errorf("解析响应失败: %v", err)
		return nil, errors.New("failed to parse response")
	}

	// 验证响应完整性
	if lifiResp.TransactionRequest.To == "" {
		l.Errorf("LI.FI 响应缺少交易参数")
		return nil, errors.New("invalid LI.FI response: missing transaction request")
	}

	// 构建成功消息
	message := fmt.Sprintf("✅ 跨链报价获取成功！从链 %d 到链 %d，使用工具: %s，预计 Gas 费用已计算完成。",
		req.FromChain, req.ToChain, lifiResp.Tool)

	l.Infof("✅ 跨链报价获取成功，使用工具: %s", lifiResp.Tool)

	return &types.BridgeQuoteResp{
		TransactionRequest: lifiResp.TransactionRequest,
		Estimate:           lifiResp.Estimate,
		Message:            message,
	}, nil
}

// ExecuteBridge 执行跨链转账
func (l *BridgeLogic) ExecuteBridge(req *types.BridgeExecuteReq) (*types.BridgeExecuteResp, error) {
	l.Infof("--- 开始执行跨链转账 fromChain=%d toChain=%d ---", req.FromChain, req.ToChain)

	// 1. 先获取报价
	quoteReq := &types.BridgeQuoteReq{
		FromChain:   req.FromChain,
		ToChain:     req.ToChain,
		FromToken:   req.FromToken,
		ToToken:     req.ToToken,
		FromAmount:  req.Amount,
		FromAddress: req.FromAddress,
		ToAddress:   req.ToAddress,
		Order:       req.Order,
		Slippage:    req.Slippage,
	}

	quoteResp, err := l.GetBridgeQuote(quoteReq)
	if err != nil {
		l.Errorf("获取跨链报价失败: %v", err)
		return nil, fmt.Errorf("failed to get bridge quote: %v", err)
	}

	// 2. 获取源链配置
	chainName := l.getChainNameByID(req.FromChain)
	chainConfig, ok := l.svcCtx.Config.Chains[chainName]
	if !ok {
		l.Errorf("不支持的源链: %d", req.FromChain)
		return nil, fmt.Errorf("unsupported from chain: %d", req.FromChain)
	}

	// 3. 连接源链 RPC
	client, err := ethclient.Dial(chainConfig.RpcUrl)
	if err != nil {
		l.Errorf("连接源链 RPC 失败: %v", err)
		return nil, errors.New("failed to connect to source chain")
	}
	defer client.Close()

	// 4. 获取钱包私钥
	txLogic := NewTransactionLogic(l.ctx, l.svcCtx)
	privateKey, err := txLogic.GetWalletPrivateKey(req.FromAddress)
	if err != nil {
		return nil, err
	}

	// 5. 检查是否需要 approve（ERC20 代币）
	if !txLogic.IsNativeToken(req.FromToken) && quoteResp.Estimate.ApprovalAddress != "" {
		l.Infof("需要先执行 approve 操作")
		err := l.executeApprove(client, req, quoteResp.Estimate.ApprovalAddress, privateKey, chainConfig.ChainId)
		if err != nil {
			l.Errorf("approve 操作失败: %v", err)
			return nil, fmt.Errorf("approve failed: %v", err)
		}
	}

	// 6. 构建并发送跨链交易
	txHash, err := l.sendBridgeTransaction(client, quoteResp.TransactionRequest, privateKey, chainConfig.ChainId)
	if err != nil {
		l.Errorf("发送跨链交易失败: %v", err)
		return nil, fmt.Errorf("failed to send bridge transaction: %v", err)
	}

	// 7. 构建响应
	explorerUrl := l.buildBridgeExplorerUrl(req.FromChain, txHash)
	message := fmt.Sprintf("✅ 跨链转账已提交！从链 %d 到链 %d，交易哈希: %s", req.FromChain, req.ToChain, txHash)

	return &types.BridgeExecuteResp{
		TxHash:      txHash,
		Message:     message,
		ExplorerUrl: explorerUrl,
		FromChain:   req.FromChain,
		ToChain:     req.ToChain,
		Status:      "pending",
	}, nil
}

// GetBridgeStatus 查询跨链状态
func (l *BridgeLogic) GetBridgeStatus(req *types.BridgeStatusReq) (*types.BridgeStatusResp, error) {
	l.Infof("--- 查询跨链状态 txHash=%s ---", req.TxHash)

	// 调用 LI.FI 状态查询 API
	params := url.Values{}
	params.Set("txHash", req.TxHash)

	apiURL := l.svcCtx.Config.Lifi.ApiUrl + "/status?" + params.Encode()
	l.Infof("调用 LI.FI 状态 API: %s", apiURL)

	// 创建 HTTP 客户端，设置超时
	client := &http.Client{Timeout: 30 * time.Second}

	req_http, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建状态查询请求失败: %v", err)
	}

	// 添加用户代理
	req_http.Header.Set("User-Agent", "MPC-Demo/1.0")

	resp, err := client.Do(req_http)
	if err != nil {
		l.Errorf("LI.FI 状态 API 调用失败: %v", err)
		return nil, fmt.Errorf("LI.FI 状态 API 调用失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		l.Errorf("LI.FI 状态 API 返回错误: %d, %s", resp.StatusCode, string(body))
		return nil, fmt.Errorf("LI.FI status API error: %d", resp.StatusCode)
	}

	// 解析响应
	var statusResp struct {
		Status     string `json:"status"`
		SubStatus  string `json:"substatus,omitempty"`
		TxHash     string `json:"txHash"`
		FromTxLink string `json:"fromTxLink,omitempty"`
		ToTxLink   string `json:"toTxLink,omitempty"`
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		l.Errorf("读取状态响应失败: %v", err)
		return nil, errors.New("failed to read status response")
	}

	if err := json.Unmarshal(body, &statusResp); err != nil {
		l.Errorf("解析状态响应失败: %v", err)
		return nil, errors.New("failed to parse status response")
	}

	// 构建状态消息
	var message string
	switch statusResp.Status {
	case "PENDING":
		message = "⏳ 跨链交易处理中，请耐心等待..."
	case "DONE":
		message = "✅ 跨链交易已完成！"
	case "FAILED":
		message = "❌ 跨链交易失败"
	default:
		message = fmt.Sprintf("📊 跨链交易状态: %s", statusResp.Status)
	}

	return &types.BridgeStatusResp{
		Status:     statusResp.Status,
		SubStatus:  statusResp.SubStatus,
		TxHash:     statusResp.TxHash,
		FromTxLink: statusResp.FromTxLink,
		ToTxLink:   statusResp.ToTxLink,
		Message:    message,
	}, nil
}

// executeApprove 执行 ERC20 approve 操作（按照 LI.FI 最佳实践）
func (l *BridgeLogic) executeApprove(client *ethclient.Client, req *types.BridgeExecuteReq, approvalAddress string, privateKey *ecdsa.PrivateKey, chainId int64) error {
	l.Infof("执行 ERC20 approve 操作，approvalAddress: %s", approvalAddress)

	// 构建 approve 调用数据
	// approve(address spender, uint256 amount)
	approveMethodId := []byte{0x09, 0x5e, 0xa7, 0xb3} // approve(address,uint256)

	spender := common.HexToAddress(approvalAddress)
	paddedSpender := common.LeftPadBytes(spender.Bytes(), 32)

	// 使用最大 uint256 值进行 approve（避免未来重复授权）
	maxAmount := new(big.Int)
	maxAmount.SetString("115792089237316195423570985008687907853269984665640564039457584007913129639935", 10)
	paddedAmount := common.LeftPadBytes(maxAmount.Bytes(), 32)

	data := append(approveMethodId, paddedSpender...)
	data = append(data, paddedAmount...)

	// 获取 nonce 和 gas 参数
	fromAddr := common.HexToAddress(req.FromAddress)
	nonce, err := client.PendingNonceAt(l.ctx, fromAddr)
	if err != nil {
		return fmt.Errorf("failed to get nonce: %v", err)
	}

	gasPrice, err := client.SuggestGasPrice(l.ctx)
	if err != nil {
		return fmt.Errorf("failed to get gas price: %v", err)
	}

	// 估算 gas limit
	tokenAddr := common.HexToAddress(req.FromToken)
	gasLimit, err := client.EstimateGas(l.ctx, ethereum.CallMsg{
		From: fromAddr,
		To:   &tokenAddr,
		Data: data,
	})
	if err != nil {
		l.Infof("Gas 估算失败，使用默认值: %v", err)
		gasLimit = 100000 // 使用默认值
	}

	// 增加 gas limit 缓冲（ERC20 approve 可能消耗更多 gas）
	gasLimit = gasLimit * 120 / 100 // 增加 20% 缓冲

	// 构建 approve 交易
	tx := evmTypes.NewTransaction(nonce, common.HexToAddress(req.FromToken), big.NewInt(0), gasLimit, gasPrice, data)

	// 签名并发送
	signedTx, err := evmTypes.SignTx(tx, evmTypes.NewEIP155Signer(big.NewInt(chainId)), privateKey)
	if err != nil {
		return fmt.Errorf("failed to sign approve transaction: %v", err)
	}

	err = client.SendTransaction(l.ctx, signedTx)
	if err != nil {
		// 检查错误信息中是否包含交易哈希（有些 RPC 节点会在错误信息中返回成功的交易哈希）
		if strings.Contains(err.Error(), "result") && strings.Contains(err.Error(), "0x") {
			l.Infof("⚠️ RPC 返回误导性错误，但交易可能已成功发送: %v", err)
			l.Infof("使用本地计算的交易哈希继续流程: %s", signedTx.Hash().Hex())
		} else {
			return fmt.Errorf("failed to send approve transaction: %v", err)
		}
	}

	l.Infof("✅ Approve 交易已发送: %s", signedTx.Hash().Hex())

	// 等待交易确认（至少 1 个区块确认）
	l.Infof("等待 approve 交易确认...")
	receipt, err := l.waitForTransactionReceipt(client, signedTx.Hash(), 30*time.Second)
	if err != nil {
		l.Errorf("等待 approve 交易确认失败: %v", err)
		return fmt.Errorf("failed to wait for approve transaction confirmation: %v", err)
	}

	if receipt.Status == 0 {
		l.Errorf("Approve 交易失败")
		return errors.New("approve transaction failed")
	}

	l.Infof("✅ Approve 交易确认成功，区块: %d", receipt.BlockNumber.Uint64())
	return nil
}

// waitForTransactionReceipt 等待交易确认
func (l *BridgeLogic) waitForTransactionReceipt(client *ethclient.Client, txHash common.Hash, timeout time.Duration) (*evmTypes.Receipt, error) {
	ctx, cancel := context.WithTimeout(l.ctx, timeout)
	defer cancel()

	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-ticker.C:
			receipt, err := client.TransactionReceipt(ctx, txHash)
			if err != nil {
				if err == ethereum.NotFound {
					l.Infof("交易尚未确认，继续等待...")
					continue
				}
				return nil, err
			}
			return receipt, nil
		}
	}
}

// sendBridgeTransaction 发送跨链交易（按照 LI.FI 最佳实践）
func (l *BridgeLogic) sendBridgeTransaction(client *ethclient.Client, txReq types.BridgeTxRequest, privateKey *ecdsa.PrivateKey, chainId int64) (string, error) {
	l.Infof("发送跨链交易")

	// 解析交易参数
	to := common.HexToAddress(txReq.To)
	data := common.FromHex(txReq.Data)

	value := new(big.Int)
	if txReq.Value != "" {
		_, ok := value.SetString(txReq.Value, 10)
		if !ok {
			// 尝试十六进制解析
			_, _ = value.SetString(txReq.Value, 0)
		}
	}

	// 获取 nonce
	fromAddr := crypto.PubkeyToAddress(privateKey.PublicKey)
	nonce, err := client.PendingNonceAt(l.ctx, fromAddr)
	if err != nil {
		return "", fmt.Errorf("failed to get nonce: %v", err)
	}

	// 处理 gas 参数
	gasLimit := uint64(300000)
	if txReq.GasLimit != "" {
		if gl, ok := new(big.Int).SetString(txReq.GasLimit, 10); ok {
			gasLimit = gl.Uint64()
		}
	}

	gasPrice := new(big.Int)
	if txReq.GasPrice != "" {
		_, _ = gasPrice.SetString(txReq.GasPrice, 10)
	}

	// 如果 LI.FI 没有提供 gas price，使用节点建议的价格
	if gasPrice.Sign() == 0 {
		gp, err := client.SuggestGasPrice(l.ctx)
		if err != nil {
			return "", fmt.Errorf("failed to get gas price: %v", err)
		}
		gasPrice = gp
	}

	// 如果 LI.FI 没有提供 gas limit，进行估算
	if gasLimit == 0 || gasLimit == 300000 {
		estimatedGas, err := client.EstimateGas(l.ctx, ethereum.CallMsg{
			From:  fromAddr,
			To:    &to,
			Value: value,
			Data:  data,
		})
		if err != nil {
			l.Infof("Gas 估算失败，使用默认值: %v", err)
			gasLimit = 300000
		} else {
			gasLimit = estimatedGas
			// 增加 20% 缓冲
			gasLimit = gasLimit * 120 / 100
		}
	}

	l.Infof("交易参数: to=%s, value=%s, gasLimit=%d, gasPrice=%s", to.Hex(), value.String(), gasLimit, gasPrice.String())

	// 构建交易
	tx := evmTypes.NewTransaction(nonce, to, value, gasLimit, gasPrice, data)

	// 签名交易
	signedTx, err := evmTypes.SignTx(tx, evmTypes.NewEIP155Signer(big.NewInt(chainId)), privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign transaction: %v", err)
	}

	// 发送交易
	err = client.SendTransaction(l.ctx, signedTx)
	if err != nil {
		// 检查错误信息中是否包含交易哈希（有些 RPC 节点会在错误信息中返回成功的交易哈希）
		if strings.Contains(err.Error(), "result") && strings.Contains(err.Error(), "0x") {
			l.Infof("⚠️ RPC 返回误导性错误，但交易可能已成功发送: %v", err)
			l.Infof("使用本地计算的交易哈希继续流程: %s", signedTx.Hash().Hex())
		} else {
			return "", fmt.Errorf("failed to send transaction: %v", err)
		}
	}

	l.Infof("✅ 跨链交易已发送: %s", signedTx.Hash().Hex())
	return signedTx.Hash().Hex(), nil
}

// getChainNameByID 根据链ID获取链名称
func (l *BridgeLogic) getChainNameByID(chainId int) string {
	chainMap := map[int]string{
		1:        "ETH",
		5:        "ETH-Goerli",
		11155111: "ETH-Sepolia",
		56:       "BSC",
		97:       "BSC-TestNet",
		137:      "Polygon",
		80001:    "Polygon-Mumbai",
		8453:     "Base",
		10:       "Optimism",
		42161:    "Arbitrum",
	}

	if name, ok := chainMap[chainId]; ok {
		return name
	}
	return "UNKNOWN"
}

// buildBridgeExplorerUrl 构建跨链浏览器链接
func (l *BridgeLogic) buildBridgeExplorerUrl(chainId int, txHash string) string {
	explorerMap := map[int]string{
		1:        "https://etherscan.io/tx/%s",
		5:        "https://goerli.etherscan.io/tx/%s",
		11155111: "https://sepolia.etherscan.io/tx/%s",
		56:       "https://bscscan.com/tx/%s",
		97:       "https://testnet.bscscan.com/tx/%s",
		137:      "https://polygonscan.com/tx/%s",
		80001:    "https://mumbai.polygonscan.com/tx/%s",
		8453:     "https://basescan.org/tx/%s",
		10:       "https://optimistic.etherscan.io/tx/%s",
		42161:    "https://arbiscan.io/tx/%s",
	}

	if template, ok := explorerMap[chainId]; ok {
		return fmt.Sprintf(template, txHash)
	}
	return fmt.Sprintf("https://explorer.example.com/tx/%s", txHash)
}

// WrapBridge 完整的跨链操作流程（按照 LI.FI 最佳实践）
func (l *BridgeLogic) WrapBridge(req *types.BridgeExecuteReq) (*types.BridgeExecuteResp, error) {
	l.Infof("=== 开始完整跨链流程 fromChain=%d toChain=%d ===", req.FromChain, req.ToChain)

	// 步骤1: 获取跨链报价
	l.Infof("步骤1: 获取跨链报价...")
	quoteReq := &types.BridgeQuoteReq{
		FromChain:   req.FromChain,
		ToChain:     req.ToChain,
		FromToken:   req.FromToken,
		ToToken:     req.ToToken,
		FromAmount:  req.Amount,
		FromAddress: req.FromAddress,
		ToAddress:   req.ToAddress,
		Order:       req.Order,
		Slippage:    req.Slippage,
	}

	quoteResp, err := l.GetBridgeQuote(quoteReq)
	if err != nil {
		l.Errorf("获取跨链报价失败: %v", err)
		return nil, fmt.Errorf("failed to get bridge quote: %v", err)
	}
	l.Infof("✅ 跨链报价获取成功")

	// 步骤2: 获取源链配置和连接
	chainName := l.getChainNameByID(req.FromChain)
	chainConfig, ok := l.svcCtx.Config.Chains[chainName]
	if !ok {
		l.Errorf("不支持的源链: %d", req.FromChain)
		return nil, fmt.Errorf("unsupported from chain: %d", req.FromChain)
	}

	client, err := ethclient.Dial(chainConfig.RpcUrl)
	if err != nil {
		l.Errorf("连接源链 RPC 失败: %v", err)
		return nil, errors.New("failed to connect to source chain")
	}
	defer client.Close()

	// 步骤3: 获取钱包私钥
	wallet, err := l.svcCtx.WalletsDao.FindOneByAddress(l.ctx, req.FromAddress)
	if err != nil {
		l.Errorf("查询钱包失败: %v", err)
		return nil, errors.New("wallet not found")
	}

	privateKey, err := crypto.HexToECDSA(wallet.EncryptedPrivateKey)
	if err != nil {
		l.Errorf("私钥解析失败: %v", err)
		return nil, errors.New("invalid private key")
	}

	// 创建 TransactionLogic 实例用于调用通用方法
	txLogic := NewTransactionLogic(l.ctx, l.svcCtx)

	// 步骤4: 检查并执行 ERC20 approve（如果需要）
	if !txLogic.IsNativeToken(req.FromToken) && quoteResp.Estimate.ApprovalAddress != "" {
		l.Infof("步骤4: 检查并执行 ERC20 approve...")

		// 检查当前 allowance
		currentAllowance, err := txLogic.CheckAllowance(client, req.FromToken, req.FromAddress, quoteResp.Estimate.ApprovalAddress)
		if err != nil {
			l.Errorf("检查 allowance 失败: %v", err)
			return nil, fmt.Errorf("failed to check allowance: %v", err)
		}

		amount, _ := new(big.Int).SetString(req.Amount, 10)
		if currentAllowance.Cmp(amount) < 0 {
			l.Infof("当前 allowance 不足，需要执行 approve")
			maxAmount := new(big.Int)
			maxAmount.SetString("115792089237316195423570985008687907853269984665640564039457584007913129639935", 10)

			_, err := txLogic.ExecuteApproveTransaction(client, privateKey, req.FromToken, quoteResp.Estimate.ApprovalAddress, maxAmount, chainConfig.ChainId)
			if err != nil {
				l.Errorf("approve 操作失败: %v", err)
				return nil, fmt.Errorf("approve failed: %v", err)
			}
			l.Infof("✅ ERC20 approve 完成")
		} else {
			l.Infof("✅ 当前 allowance 充足，无需 approve")
		}
	} else {
		l.Infof("步骤4: 原生代币，跳过 approve")
	}

	// 步骤5: 发送主跨链交易
	l.Infof("步骤5: 发送主跨链交易...")
	txHash, err := l.sendBridgeTransactionWithRetry(client, quoteResp.TransactionRequest, privateKey, chainConfig.ChainId)
	if err != nil {
		l.Errorf("发送跨链交易失败: %v", err)
		return nil, fmt.Errorf("failed to send bridge transaction: %v", err)
	}
	l.Infof("✅ 跨链交易已发送: %s", txHash)

	// 步骤6: 构建响应
	explorerUrl := l.buildBridgeExplorerUrl(req.FromChain, txHash)
	message := fmt.Sprintf("✅ 跨链转账已提交！从链 %d 到链 %d，交易哈希: %s。请使用 /bridge/status 查询进度。", req.FromChain, req.ToChain, txHash)

	return &types.BridgeExecuteResp{
		TxHash:      txHash,
		Message:     message,
		ExplorerUrl: explorerUrl,
		FromChain:   req.FromChain,
		ToChain:     req.ToChain,
		Status:      "pending",
	}, nil
}

// checkAllowance 检查 ERC20 代币的 allowance
func (l *BridgeLogic) checkAllowance(client *ethclient.Client, tokenAddress, owner, spender string) (*big.Int, error) {
	l.Infof("检查 allowance: token=%s, owner=%s, spender=%s", tokenAddress, owner, spender)

	// 构建 allowance 调用数据
	// allowance(address owner, address spender) returns (uint256)
	allowanceMethodId := []byte{0xdd, 0x62, 0xed, 0x3e} // allowance(address,address)

	ownerAddr := common.HexToAddress(owner)
	spenderAddr := common.HexToAddress(spender)

	paddedOwner := common.LeftPadBytes(ownerAddr.Bytes(), 32)
	paddedSpender := common.LeftPadBytes(spenderAddr.Bytes(), 32)

	data := append(allowanceMethodId, paddedOwner...)
	data = append(data, paddedSpender...)

	// 调用合约
	tokenAddr := common.HexToAddress(tokenAddress)
	result, err := client.CallContract(l.ctx, ethereum.CallMsg{
		To:   &tokenAddr,
		Data: data,
	}, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to call allowance: %v", err)
	}

	// 解析结果
	allowance := new(big.Int).SetBytes(result)
	l.Infof("当前 allowance: %s", allowance.String())
	return allowance, nil
}

// executeApproveWithRetry 带重试的 approve 操作
func (l *BridgeLogic) executeApproveWithRetry(client *ethclient.Client, req *types.BridgeExecuteReq, approvalAddress string, privateKey *ecdsa.PrivateKey, chainId int64) error {
	maxRetries := 3
	for i := 0; i < maxRetries; i++ {
		err := l.executeApprove(client, req, approvalAddress, privateKey, chainId)
		if err != nil {
			l.Errorf("approve 操作失败 (尝试 %d/%d): %v", i+1, maxRetries, err)
			if i == maxRetries-1 {
				return err
			}
			time.Sleep(time.Duration(i+1) * time.Second)
			continue
		}
		return nil
	}
	return errors.New("approve 操作最终失败")
}

// sendBridgeTransactionWithRetry 带重试的跨链交易发送
func (l *BridgeLogic) sendBridgeTransactionWithRetry(client *ethclient.Client, txReq types.BridgeTxRequest, privateKey *ecdsa.PrivateKey, chainId int64) (string, error) {
	maxRetries := 3
	for i := 0; i < maxRetries; i++ {
		txHash, err := l.sendBridgeTransaction(client, txReq, privateKey, chainId)
		if err != nil {
			l.Errorf("发送跨链交易失败 (尝试 %d/%d): %v", i+1, maxRetries, err)
			if i == maxRetries-1 {
				return "", err
			}
			time.Sleep(time.Duration(i+1) * time.Second)
			continue
		}
		return txHash, nil
	}
	return "", errors.New("发送跨链交易最终失败")
}

// GetBridgeStatusWithPolling 带轮询的跨链状态查询
func (l *BridgeLogic) GetBridgeStatusWithPolling(txHash string, maxAttempts int) (*types.BridgeStatusResp, error) {
	l.Infof("开始轮询跨链状态: %s (最多 %d 次尝试)", txHash, maxAttempts)

	for i := 0; i < maxAttempts; i++ {
		req := &types.BridgeStatusReq{TxHash: txHash}
		resp, err := l.GetBridgeStatus(req)
		if err != nil {
			l.Errorf("查询状态失败 (尝试 %d/%d): %v", i+1, maxAttempts, err)
			if i == maxAttempts-1 {
				return nil, err
			}
			time.Sleep(15 * time.Second)
			continue
		}

		l.Infof("状态查询结果: %s (子状态: %s)", resp.Status, resp.SubStatus)

		// 如果交易完成或失败，停止轮询
		if resp.Status == "DONE" || resp.Status == "FAILED" {
			return resp, nil
		}

		// 继续轮询
		if i < maxAttempts-1 {
			l.Infof("交易仍在处理中，15秒后再次查询...")
			time.Sleep(15 * time.Second)
		}
	}

	// 轮询超时
	return &types.BridgeStatusResp{
		Status:  "TIMEOUT",
		TxHash:  txHash,
		Message: "⏰ 状态查询超时，请手动检查交易状态",
	}, nil
}

// GetSupportedChains 获取支持的链列表
func (l *BridgeLogic) GetSupportedChains() ([]ChainInfo, error) {
	l.Infof("获取支持的链列表")

	apiURL := l.svcCtx.Config.Lifi.ApiUrl + "/chains"

	// 创建 HTTP 客户端，设置超时
	client := &http.Client{Timeout: 30 * time.Second}

	req_http, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	// 添加用户代理
	req_http.Header.Set("User-Agent", "MPC-Demo/1.0")

	resp, err := client.Do(req_http)
	if err != nil {
		l.Errorf("获取链列表失败: %v", err)
		return nil, fmt.Errorf("获取链列表失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: %d", resp.StatusCode)
	}

	var chainsResp struct {
		Chains []ChainInfo `json:"chains"`
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("failed to read response")
	}

	if err := json.Unmarshal(body, &chainsResp); err != nil {
		return nil, errors.New("failed to parse response")
	}

	return chainsResp.Chains, nil
}

// normalizeTokenAddress 标准化代币地址（转换为 LI.FI 格式）
func (l *BridgeLogic) normalizeTokenAddress(tokenAddr string) string {
	if tokenAddr == "0x0000000000000000000000000000000000000000" {
		return "0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE" // LI.FI 原生代币标识
	}
	return tokenAddr
}

// GetSupportedTokens 获取支持的代币列表
func (l *BridgeLogic) GetSupportedTokens(chainIds []int) (map[int][]TokenInfo, error) {
	l.Infof("获取支持的代币列表")

	params := url.Values{}
	if len(chainIds) > 0 {
		chainStrs := make([]string, len(chainIds))
		for i, id := range chainIds {
			chainStrs[i] = strconv.Itoa(id)
		}
		params.Set("chains", strings.Join(chainStrs, ","))
	}

	apiURL := l.svcCtx.Config.Lifi.ApiUrl + "/tokens?" + params.Encode()

	// 创建 HTTP 客户端，设置超时
	client := &http.Client{Timeout: 30 * time.Second}

	req_http, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	// 添加用户代理
	req_http.Header.Set("User-Agent", "MPC-Demo/1.0")

	resp, err := client.Do(req_http)
	if err != nil {
		l.Errorf("获取代币列表失败: %v", err)
		return nil, fmt.Errorf("获取代币列表失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: %d", resp.StatusCode)
	}

	var tokensResp map[int][]TokenInfo
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("failed to read response")
	}

	if err := json.Unmarshal(body, &tokensResp); err != nil {
		return nil, errors.New("failed to parse response")
	}

	return tokensResp, nil
}

// ChainInfo 链信息
type ChainInfo struct {
	Key       string `json:"key"`
	Name      string `json:"name"`
	ChainType string `json:"chainType"`
	Coin      string `json:"coin"`
	Id        int    `json:"id"`
	Mainnet   bool   `json:"mainnet"`
	LogoURI   string `json:"logoURI"`
}

// TokenInfo 代币信息
type TokenInfo struct {
	Address  string `json:"address"`
	Decimals int    `json:"decimals"`
	Symbol   string `json:"symbol"`
	ChainId  int    `json:"chainId"`
	Name     string `json:"name"`
	CoinKey  string `json:"coinKey"`
	LogoURI  string `json:"logoURI"`
	PriceUSD string `json:"priceUSD"`
}
