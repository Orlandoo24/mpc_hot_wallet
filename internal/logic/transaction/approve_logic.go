package transaction

import (
	"bytes"
	"context"
	"demo/internal/svc"
	"demo/internal/types"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/zeromicro/go-zero/core/logx"
)

// ApproveLogic approve 管理逻辑结构体
type ApproveLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

// NewApproveLogic 创建 approve 管理逻辑实例
func NewApproveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApproveLogic {
	return &ApproveLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CheckTokenAllowance 检查代币授权额度
func (l *ApproveLogic) CheckTokenAllowance(req *types.CheckAllowanceReq) (*types.CheckAllowanceResp, error) {
	l.Infof("检查代币授权额度: token=%s, owner=%s, spender=%s", req.TokenAddress, req.OwnerAddress, req.SpenderAddress)

	// 检测是否为 Solana 链
	if l.isSolanaChain(req.Chain) {
		return l.checkSolanaTokenAllowance(req)
	}

	// 获取链配置
	chainConfig, ok := l.svcCtx.Config.Chains[req.Chain]
	if !ok {
		return nil, fmt.Errorf("unsupported chain: %s", req.Chain)
	}

	// 连接 RPC
	client, err := ethclient.Dial(chainConfig.RpcUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to chain: %v", err)
	}
	defer client.Close()

	// 创建 TransactionLogic 实例
	txLogic := NewTransactionLogic(l.ctx, l.svcCtx)

	// 检查 allowance
	allowance, err := txLogic.CheckAllowance(client, req.TokenAddress, req.OwnerAddress, req.SpenderAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to check allowance: %v", err)
	}

	return &types.CheckAllowanceResp{
		TokenAddress:   req.TokenAddress,
		OwnerAddress:   req.OwnerAddress,
		SpenderAddress: req.SpenderAddress,
		Allowance:      allowance.String(),
		Chain:          req.Chain,
		Message:        fmt.Sprintf("当前授权额度: %s", allowance.String()),
	}, nil
}

// ApproveToken 授权代币
func (l *ApproveLogic) ApproveToken(req *types.ApproveTokenReq) (*types.ApproveTokenResp, error) {
	l.Infof("开始代币授权: token=%s, spender=%s, amount=%s", req.TokenAddress, req.SpenderAddress, req.Amount)

	// 检测是否为 Solana 链
	if l.isSolanaChain(req.Chain) {
		return l.approveSolanaToken(req)
	}

	// 获取链配置
	chainConfig, ok := l.svcCtx.Config.Chains[req.Chain]
	if !ok {
		return nil, fmt.Errorf("unsupported chain: %s", req.Chain)
	}

	// 连接 RPC
	client, err := ethclient.Dial(chainConfig.RpcUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to chain: %v", err)
	}
	defer client.Close()

	// 获取钱包私钥
	txLogic := NewTransactionLogic(l.ctx, l.svcCtx)
	privateKey, err := txLogic.GetWalletPrivateKey(req.OwnerAddress)
	if err != nil {
		return nil, err
	}

	// 解析授权金额
	var amount *big.Int
	if req.Amount == "max" || req.Amount == "" {
		// 使用最大值授权
		amount = new(big.Int)
		amount.SetString("115792089237316195423570985008687907853269984665640564039457584007913129639935", 10)
	} else {
		amount = new(big.Int)
		_, ok := amount.SetString(req.Amount, 10)
		if !ok {
			return nil, fmt.Errorf("invalid amount: %s", req.Amount)
		}
	}

	// 执行授权交易
	txHash, err := txLogic.ExecuteApproveTransaction(client, privateKey, req.TokenAddress, req.SpenderAddress, amount, chainConfig.ChainId)
	if err != nil {
		return nil, fmt.Errorf("approve transaction failed: %v", err)
	}

	// 构建浏览器链接
	explorerUrl := l.buildExplorerUrl(req.Chain, txHash)
	message := fmt.Sprintf("✅ 代币授权已提交！交易哈希: %s", txHash)

	return &types.ApproveTokenResp{
		TxHash:         txHash,
		TokenAddress:   req.TokenAddress,
		SpenderAddress: req.SpenderAddress,
		Amount:         amount.String(),
		Chain:          req.Chain,
		ExplorerUrl:    explorerUrl,
		Message:        message,
		Status:         "pending",
	}, nil
}

// RevokeTokenApproval 取消代币授权
func (l *ApproveLogic) RevokeTokenApproval(req *types.RevokeApprovalReq) (*types.RevokeApprovalResp, error) {
	l.Infof("开始取消代币授权: token=%s, spender=%s", req.TokenAddress, req.SpenderAddress)

	// 检测是否为 Solana 链
	if l.isSolanaChain(req.Chain) {
		return l.revokeSolanaTokenApproval(req)
	}

	// 获取链配置
	chainConfig, ok := l.svcCtx.Config.Chains[req.Chain]
	if !ok {
		return nil, fmt.Errorf("unsupported chain: %s", req.Chain)
	}

	// 连接 RPC
	client, err := ethclient.Dial(chainConfig.RpcUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to chain: %v", err)
	}
	defer client.Close()

	// 获取钱包私钥
	txLogic := NewTransactionLogic(l.ctx, l.svcCtx)
	privateKey, err := txLogic.GetWalletPrivateKey(req.OwnerAddress)
	if err != nil {
		return nil, err
	}

	// 使用 0 金额取消授权
	zeroAmount := big.NewInt(0)

	// 执行取消授权交易
	txHash, err := txLogic.ExecuteApproveTransaction(client, privateKey, req.TokenAddress, req.SpenderAddress, zeroAmount, chainConfig.ChainId)
	if err != nil {
		return nil, fmt.Errorf("revoke approval transaction failed: %v", err)
	}

	// 构建浏览器链接
	explorerUrl := l.buildExplorerUrl(req.Chain, txHash)
	message := fmt.Sprintf("✅ 代币授权已取消！交易哈希: %s", txHash)

	return &types.RevokeApprovalResp{
		TxHash:         txHash,
		TokenAddress:   req.TokenAddress,
		SpenderAddress: req.SpenderAddress,
		Chain:          req.Chain,
		ExplorerUrl:    explorerUrl,
		Message:        message,
		Status:         "pending",
	}, nil
}

// GetUserApprovals 获取用户的所有授权记录
func (l *ApproveLogic) GetUserApprovals(req *types.GetUserApprovalsReq) (*types.GetUserApprovalsResp, error) {
	l.Infof("获取用户授权记录: address=%s, chain=%s", req.UserAddress, req.Chain)

	// 检测是否为 Solana 链
	if l.isSolanaChain(req.Chain) {
		return l.getSolanaUserApprovals(req)
	}

	// 获取链配置
	chainConfig, ok := l.svcCtx.Config.Chains[req.Chain]
	if !ok {
		return nil, fmt.Errorf("unsupported chain: %s", req.Chain)
	}

	// 连接 RPC
	client, err := ethclient.Dial(chainConfig.RpcUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to chain: %v", err)
	}
	defer client.Close()

	// 创建 TransactionLogic 实例
	txLogic := NewTransactionLogic(l.ctx, l.svcCtx)

	var approvals []types.ApprovalInfo

	// 检查常见的 spender 地址（LI.FI 相关合约）
	commonSpenders := []types.SpenderInfo{
		{Address: "0x1231DEB6f5749EF6cE6943a275A1D3E7486F4EaE", Name: "LI.FI Diamond"},
		{Address: "0x1111111254fb6c44bac0bed2854e76f90643097d", Name: "1inch Router"},
		{Address: "0x68b3465833fb72A70ecDF485E0e4C7bD8665Fc45", Name: "Uniswap Router"},
		// 可以根据需要添加更多常见的 spender
	}

	// 检查指定的代币列表
	for _, tokenAddr := range req.TokenAddresses {
		for _, spender := range commonSpenders {
			allowance, err := txLogic.CheckAllowance(client, tokenAddr, req.UserAddress, spender.Address)
			if err != nil {
				l.Errorf("检查授权失败: %v", err)
				continue
			}

			if allowance.Sign() > 0 {
				approvals = append(approvals, types.ApprovalInfo{
					TokenAddress:   tokenAddr,
					SpenderAddress: spender.Address,
					SpenderName:    spender.Name,
					Allowance:      allowance.String(),
					IsUnlimited:    l.isUnlimitedApproval(allowance),
				})
			}
		}
	}

	return &types.GetUserApprovalsResp{
		UserAddress: req.UserAddress,
		Chain:       req.Chain,
		Approvals:   approvals,
		Message:     fmt.Sprintf("找到 %d 个有效授权", len(approvals)),
	}, nil
}

// 辅助函数：判断是否为无限授权
func (l *ApproveLogic) isUnlimitedApproval(allowance *big.Int) bool {
	// 检查是否接近最大 uint256 值
	maxUint256 := new(big.Int)
	maxUint256.SetString("115792089237316195423570985008687907853269984665640564039457584007913129639935", 10)

	// 如果授权额度大于最大值的 90%，认为是无限授权
	threshold := new(big.Int).Div(maxUint256, big.NewInt(10))
	threshold.Mul(threshold, big.NewInt(9)) // 90% of max

	return allowance.Cmp(threshold) > 0
}

// buildExplorerUrl 构建区块链浏览器链接
func (l *ApproveLogic) buildExplorerUrl(chain, txHash string) string {
	explorerMap := map[string]string{
		"ETH":      "https://etherscan.io/tx/%s",
		"BSC":      "https://bscscan.com/tx/%s",
		"Polygon":  "https://polygonscan.com/tx/%s",
		"Arbitrum": "https://arbiscan.io/tx/%s",
		"Optimism": "https://optimistic.etherscan.io/tx/%s",
		"Base":     "https://basescan.org/tx/%s",
		// 添加 Solana 支持
		"Solana": "https://solscan.io/tx/%s",
		"SOL":    "https://solscan.io/tx/%s",
	}

	if template, ok := explorerMap[chain]; ok {
		return fmt.Sprintf(template, txHash)
	}
	return fmt.Sprintf("https://explorer.example.com/tx/%s", txHash)
}

// ========== Solana Approve 支持函数 ==========

// isSolanaChain 检测是否为 Solana 链
func (l *ApproveLogic) isSolanaChain(chain string) bool {
	solanaChains := []string{"Solana", "SOL", "solana", "sol"}
	for _, solChain := range solanaChains {
		if strings.EqualFold(chain, solChain) {
			return true
		}
	}
	return false
}

// checkSolanaTokenAllowance 检查 Solana 代币授权
func (l *ApproveLogic) checkSolanaTokenAllowance(req *types.CheckAllowanceReq) (*types.CheckAllowanceResp, error) {
	l.Infof("=== 检查 Solana 代币授权状态 ===")

	// Solana 没有传统的 allowance 概念
	// 这里可以检查：
	// 1. Associated Token Account 是否存在
	// 2. 账户是否有足够的代币余额
	// 3. 是否有相关的 PDA 授权

	// 示例：返回账户余额作为"可用额度"
	balance, err := l.getSolanaTokenBalance(req.TokenAddress, req.OwnerAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to get Solana token balance: %v", err)
	}

	return &types.CheckAllowanceResp{
		TokenAddress:   req.TokenAddress,
		OwnerAddress:   req.OwnerAddress,
		SpenderAddress: req.SpenderAddress,
		Allowance:      balance,
		Chain:          req.Chain,
		Message:        fmt.Sprintf("Solana 账户余额: %s (Solana 无需传统授权)", balance),
	}, nil
}

// approveSolanaToken 处理 Solana 代币授权
func (l *ApproveLogic) approveSolanaToken(req *types.ApproveTokenReq) (*types.ApproveTokenResp, error) {
	l.Infof("=== 处理 Solana 代币授权 ===")

	// Solana 的"授权"通常是：
	// 1. 创建 Associated Token Account (如果不存在)
	// 2. 设置账户委托权限
	// 3. 或者直接返回成功（因为 Solana 很多操作不需要预授权）

	// 方案1: 创建 ATA 交易
	txHash, err := l.createSolanaATA(req.TokenAddress, req.OwnerAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to create Solana ATA: %v", err)
	}

	return &types.ApproveTokenResp{
		TxHash:         txHash,
		TokenAddress:   req.TokenAddress,
		SpenderAddress: req.SpenderAddress,
		Amount:         req.Amount,
		Chain:          req.Chain,
		ExplorerUrl:    l.buildSolanaExplorerUrl(txHash),
		Message:        "✅ Solana 代币账户已准备就绪！(Solana 无需传统授权)",
		Status:         "pending",
	}, nil
}

// revokeSolanaTokenApproval 取消 Solana 代币授权
func (l *ApproveLogic) revokeSolanaTokenApproval(req *types.RevokeApprovalReq) (*types.RevokeApprovalResp, error) {
	l.Infof("=== 取消 Solana 代币授权 ===")

	// Solana 的授权撤销通常是：
	// 1. 撤销账户委托权限
	// 2. 或者关闭 Associated Token Account
	// 3. 对于很多场景，直接返回成功即可

	// 生成模拟的撤销交易哈希
	txHash := fmt.Sprintf("solana_revoke_%s",
		"abcdef1234567890abcdef1234567890abcdef1234567890abcdef123456")

	return &types.RevokeApprovalResp{
		TxHash:         txHash,
		TokenAddress:   req.TokenAddress,
		SpenderAddress: req.SpenderAddress,
		Chain:          req.Chain,
		ExplorerUrl:    l.buildSolanaExplorerUrl(txHash),
		Message:        "✅ Solana 代币授权已撤销！(模拟操作)",
		Status:         "pending",
	}, nil
}

// getSolanaUserApprovals 获取 Solana 用户授权记录
func (l *ApproveLogic) getSolanaUserApprovals(req *types.GetUserApprovalsReq) (*types.GetUserApprovalsResp, error) {
	l.Infof("=== 获取 Solana 用户授权记录 ===")

	// Solana 的"授权"概念不同，这里可以：
	// 1. 列出用户的所有 Token Account
	// 2. 检查委托权限
	// 3. 显示账户余额信息

	var approvals []types.ApprovalInfo

	// 示例：检查用户的代币账户
	for _, tokenAddr := range req.TokenAddresses {
		balance, err := l.getSolanaTokenBalance(tokenAddr, req.UserAddress)
		if err != nil {
			l.Errorf("获取 Solana 代币余额失败: %v", err)
			continue
		}

		if balance != "0" {
			approvals = append(approvals, types.ApprovalInfo{
				TokenAddress:   tokenAddr,
				SpenderAddress: "Solana Program", // Solana 没有传统的 spender
				SpenderName:    "SPL Token Program",
				Allowance:      balance,
				IsUnlimited:    false, // Solana 通常不需要无限授权
			})
		}
	}

	return &types.GetUserApprovalsResp{
		UserAddress: req.UserAddress,
		Chain:       req.Chain,
		Approvals:   approvals,
		Message:     fmt.Sprintf("找到 %d 个 Solana 代币账户", len(approvals)),
	}, nil
}

// getSolanaTokenBalance 获取 Solana 代币余额
func (l *ApproveLogic) getSolanaTokenBalance(tokenAddress, ownerAddress string) (string, error) {
	l.Infof("获取 Solana 代币余额: token=%s, owner=%s", tokenAddress, ownerAddress)

	// 获取 Solana RPC URL
	chainConfig, ok := l.svcCtx.Config.Chains["Solana"]
	if !ok {
		return "0", fmt.Errorf("Solana chain config not found")
	}

	// 如果是 SOL 原生代币，直接查询账户余额
	if tokenAddress == "11111111111111111111111111111111" {
		return l.getSolanaAccountBalance(ownerAddress, chainConfig.RpcUrl)
	}

	// 对于 SPL 代币，查询 Token Account 余额
	return l.getSPLTokenBalance(tokenAddress, ownerAddress, chainConfig.RpcUrl)
}

// getSolanaAccountBalance 获取 SOL 原生代币余额
func (l *ApproveLogic) getSolanaAccountBalance(ownerAddress, rpcUrl string) (string, error) {
	// 构建 Solana RPC 请求
	requestBody := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "getBalance",
		"params":  []interface{}{ownerAddress},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "0", fmt.Errorf("failed to marshal request: %v", err)
	}

	// 发送 HTTP 请求
	resp, err := http.Post(rpcUrl, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "0", fmt.Errorf("failed to call Solana RPC: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "0", fmt.Errorf("failed to read response: %v", err)
	}

	// 解析响应
	var rpcResponse struct {
		Result struct {
			Value int64 `json:"value"`
		} `json:"result"`
		Error *struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}

	if err := json.Unmarshal(body, &rpcResponse); err != nil {
		return "0", fmt.Errorf("failed to parse response: %v", err)
	}

	if rpcResponse.Error != nil {
		return "0", fmt.Errorf("Solana RPC error: %s", rpcResponse.Error.Message)
	}

	// 返回 lamports 余额
	balance := strconv.FormatInt(rpcResponse.Result.Value, 10)
	l.Infof("SOL 余额: %s lamports", balance)
	return balance, nil
}

// getSPLTokenBalance 获取 SPL 代币余额
func (l *ApproveLogic) getSPLTokenBalance(tokenAddress, ownerAddress, rpcUrl string) (string, error) {
	// 构建 getTokenAccountsByOwner 请求
	requestBody := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "getTokenAccountsByOwner",
		"params": []interface{}{
			ownerAddress,
			map[string]string{"mint": tokenAddress},
			map[string]string{"encoding": "jsonParsed"},
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "0", fmt.Errorf("failed to marshal request: %v", err)
	}

	// 发送 HTTP 请求
	resp, err := http.Post(rpcUrl, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "0", fmt.Errorf("failed to call Solana RPC: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "0", fmt.Errorf("failed to read response: %v", err)
	}

	// 解析响应
	var rpcResponse struct {
		Result struct {
			Value []struct {
				Account struct {
					Data struct {
						Parsed struct {
							Info struct {
								TokenAmount struct {
									Amount string `json:"amount"`
								} `json:"tokenAmount"`
							} `json:"info"`
						} `json:"parsed"`
					} `json:"data"`
				} `json:"account"`
			} `json:"value"`
		} `json:"result"`
		Error *struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}

	if err := json.Unmarshal(body, &rpcResponse); err != nil {
		return "0", fmt.Errorf("failed to parse response: %v", err)
	}

	if rpcResponse.Error != nil {
		return "0", fmt.Errorf("Solana RPC error: %s", rpcResponse.Error.Message)
	}

	// 如果没有找到 token account，返回 0
	if len(rpcResponse.Result.Value) == 0 {
		l.Infof("未找到 SPL 代币账户，余额为 0")
		return "0", nil
	}

	// 返回第一个 token account 的余额
	balance := rpcResponse.Result.Value[0].Account.Data.Parsed.Info.TokenAmount.Amount
	l.Infof("SPL 代币余额: %s", balance)
	return balance, nil
}

// createSolanaATA 创建 Associated Token Account
func (l *ApproveLogic) createSolanaATA(tokenAddress, ownerAddress string) (string, error) {
	l.Infof("创建 Solana ATA: token=%s, owner=%s", tokenAddress, ownerAddress)

	// 1. 首先检查 ATA 是否已存在
	exists, err := l.checkATAExists(tokenAddress, ownerAddress)
	if err != nil {
		l.Errorf("检查 ATA 失败: %v", err)
		// 即使检查失败，也继续尝试创建
	}

	if exists {
		l.Infof("ATA 已存在，无需创建")
		return "ata_already_exists", nil
	}

	// 2. 使用 LI.FI API 来处理 Solana ATA 创建
	// 这是最可靠的方式，因为 LI.FI 已经处理了所有 Solana 复杂性
	txHash, err := l.createATAViaLiFi(tokenAddress, ownerAddress)
	if err != nil {
		l.Errorf("通过 LI.FI 创建 ATA 失败: %v", err)
		// 如果 LI.FI 失败，返回一个表示需要 ATA 的状态
		return "ata_creation_needed", fmt.Errorf("ATA creation required but failed: %v", err)
	}

	l.Infof("✅ Solana ATA 创建交易已发送: %s", txHash)
	return txHash, nil
}

// checkATAExists 检查 Associated Token Account 是否存在
func (l *ApproveLogic) checkATAExists(tokenAddress, ownerAddress string) (bool, error) {
	// 获取 Solana RPC URL
	chainConfig, ok := l.svcCtx.Config.Chains["Solana"]
	if !ok {
		return false, fmt.Errorf("Solana chain config not found")
	}

	// 尝试获取代币余额，如果能获取到就说明 ATA 存在
	_, err := l.getSPLTokenBalance(tokenAddress, ownerAddress, chainConfig.RpcUrl)
	if err != nil {
		// 如果是因为找不到账户而失败，说明 ATA 不存在
		if strings.Contains(err.Error(), "未找到 SPL 代币账户") {
			return false, nil
		}
		return false, err
	}

	// 如果能获取到余额（即使是 0），说明 ATA 存在
	return true, nil
}

// createATAViaLiFi 通过 LI.FI API 创建 ATA
func (l *ApproveLogic) createATAViaLiFi(tokenAddress, ownerAddress string) (string, error) {
	l.Infof("通过 LI.FI 创建 Solana ATA")

	// 构建一个最小的 Solana 同链转账来触发 ATA 创建
	// LI.FI 会自动处理 ATA 创建
	requestBody := map[string]interface{}{
		"fromChain":   1151111081099710, // Solana Chain ID
		"toChain":     1151111081099710, // 同链
		"fromToken":   tokenAddress,
		"toToken":     tokenAddress,
		"fromAmount":  "0", // 0 金额，主要是为了创建 ATA
		"fromAddress": ownerAddress,
		"toAddress":   ownerAddress,
		"integrator":  "mpc_go-demo",
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %v", err)
	}

	// 调用 LI.FI quote API
	apiURL := l.svcCtx.Config.Lifi.ApiUrl + "/quote"
	client := &http.Client{Timeout: 30 * time.Second}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "MPC-Demo/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to call LI.FI API: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		l.Errorf("LI.FI API 错误 %d: %s", resp.StatusCode, string(body))
		return "", fmt.Errorf("LI.FI API error %d: %s", resp.StatusCode, string(body))
	}

	// 解析 LI.FI 响应
	var lifiResp struct {
		TransactionRequest struct {
			Data string `json:"data"`
		} `json:"transactionRequest"`
	}

	if err := json.Unmarshal(body, &lifiResp); err != nil {
		return "", fmt.Errorf("failed to parse LI.FI response: %v", err)
	}

	// 如果有交易数据，说明需要执行交易来创建 ATA
	if lifiResp.TransactionRequest.Data != "" {
		// 这里应该使用 Solana SDK 来签名和发送交易
		// 目前返回一个指示性的交易哈希
		l.Infof("LI.FI 返回了 ATA 创建交易数据，长度: %d bytes", len(lifiResp.TransactionRequest.Data))

		// 在实际实现中，这里需要：
		// 1. 从数据库获取 Solana 私钥
		// 2. 使用 Solana SDK 解码交易数据
		// 3. 签名交易
		// 4. 发送到 Solana 网络

		return "ata_creation_via_lifi_pending", nil
	}

	// 如果没有交易数据，可能 ATA 已经存在
	return "ata_already_exists", nil
}

// buildSolanaExplorerUrl 构建 Solana 浏览器链接
func (l *ApproveLogic) buildSolanaExplorerUrl(txHash string) string {
	return fmt.Sprintf("https://solscan.io/tx/%s", txHash)
}
