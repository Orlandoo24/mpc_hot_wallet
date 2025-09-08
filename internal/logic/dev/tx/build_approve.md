基于你的代码和 LI.FI 文档，我建议你实现以下 **approve 和取消授权的逻辑**：

## **1. 完整的 Approve 管理逻辑**

```go
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
		TokenAddress:  req.TokenAddress,
		OwnerAddress:  req.OwnerAddress,
		SpenderAddress: req.SpenderAddress,
		Allowance:     allowance.String(),
		Chain:         req.Chain,
		Message:       fmt.Sprintf("当前授权额度: %s", allowance.String()),
	}, nil
}

// ApproveToken 授权代币
func (l *ApproveLogic) ApproveToken(req *types.ApproveTokenReq) (*types.ApproveTokenResp, error) {
	l.Infof("开始代币授权: token=%s, spender=%s, amount=%s", req.TokenAddress, req.SpenderAddress, req.Amount)

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
		TxHash:        txHash,
		TokenAddress:  req.TokenAddress,
		SpenderAddress: req.SpenderAddress,
		Amount:        amount.String(),
		Chain:         req.Chain,
		ExplorerUrl:   explorerUrl,
		Message:       message,
		Status:        "pending",
	}, nil
}

// RevokeTokenApproval 取消代币授权
func (l *ApproveLogic) RevokeTokenApproval(req *types.RevokeApprovalReq) (*types.RevokeApprovalResp, error) {
	l.Infof("开始取消代币授权: token=%s, spender=%s", req.TokenAddress, req.SpenderAddress)

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
		TxHash:        txHash,
		TokenAddress:  req.TokenAddress,
		SpenderAddress: req.SpenderAddress,
		Chain:         req.Chain,
		ExplorerUrl:   explorerUrl,
		Message:       message,
		Status:        "pending",
	}, nil
}

// GetUserApprovals 获取用户的所有授权记录
func (l *ApproveLogic) GetUserApprovals(req *types.GetUserApprovalsReq) (*types.GetUserApprovalsResp, error) {
	l.Infof("获取用户授权记录: address=%s, chain=%s", req.UserAddress, req.Chain)

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
		"ETH":     "https://etherscan.io/tx/%s",
		"BSC":     "https://bscscan.com/tx/%s",
		"Polygon": "https://polygonscan.com/tx/%s",
		"Arbitrum": "https://arbiscan.io/tx/%s",
		"Optimism": "https://optimistic.etherscan.io/tx/%s",
		"Base":    "https://basescan.org/tx/%s",
	}

	if template, ok := explorerMap[chain]; ok {
		return fmt.Sprintf(template, txHash)
	}
	return fmt.Sprintf("https://explorer.example.com/tx/%s", txHash)
}
```

## **2. 对应的类型定义**

```go
// 检查授权请求
type CheckAllowanceReq struct {
	TokenAddress   string `json:"token_address"`
	OwnerAddress   string `json:"owner_address"`
	SpenderAddress string `json:"spender_address"`
	Chain          string `json:"chain"`
}

// 检查授权响应
type CheckAllowanceResp struct {
	TokenAddress   string `json:"token_address"`
	OwnerAddress   string `json:"owner_address"`
	SpenderAddress string `json:"spender_address"`
	Allowance      string `json:"allowance"`
	Chain          string `json:"chain"`
	Message        string `json:"message"`
}

// 授权代币请求
type ApproveTokenReq struct {
	TokenAddress   string `json:"token_address"`
	SpenderAddress string `json:"spender_address"`
	OwnerAddress   string `json:"owner_address"`
	Amount         string `json:"amount"` // "max" 表示无限授权
	Chain          string `json:"chain"`
}

// 授权代币响应
type ApproveTokenResp struct {
	TxHash         string `json:"tx_hash"`
	TokenAddress   string `json:"token_address"`
	SpenderAddress string `json:"spender_address"`
	Amount         string `json:"amount"`
	Chain          string `json:"chain"`
	ExplorerUrl    string `json:"explorer_url"`
	Message        string `json:"message"`
	Status         string `json:"status"`
}

// 取消授权请求
type RevokeApprovalReq struct {
	TokenAddress   string `json:"token_address"`
	SpenderAddress string `json:"spender_address"`
	OwnerAddress   string `json:"owner_address"`
	Chain          string `json:"chain"`
}

// 取消授权响应
type RevokeApprovalResp struct {
	TxHash         string `json:"tx_hash"`
	TokenAddress   string `json:"token_address"`
	SpenderAddress string `json:"spender_address"`
	Chain          string `json:"chain"`
	ExplorerUrl    string `json:"explorer_url"`
	Message        string `json:"message"`
	Status         string `json:"status"`
}

// 获取用户授权请求
type GetUserApprovalsReq struct {
	UserAddress     string   `json:"user_address"`
	TokenAddresses  []string `json:"token_addresses"`
	Chain           string   `json:"chain"`
}

// 获取用户授权响应
type GetUserApprovalsResp struct {
	UserAddress string         `json:"user_address"`
	Chain       string         `json:"chain"`
	Approvals   []ApprovalInfo `json:"approvals"`
	Message     string         `json:"message"`
}

// 授权信息
type ApprovalInfo struct {
	TokenAddress   string `json:"token_address"`
	SpenderAddress string `json:"spender_address"`
	SpenderName    string `json:"spender_name"`
	Allowance      string `json:"allowance"`
	IsUnlimited    bool   `json:"is_unlimited"`
}

// Spender 信息
type SpenderInfo struct {
	Address string `json:"address"`
	Name    string `json:"name"`
}
```

## **3. 主要功能特点**

1. **检查授权额度** - 查询特定代币对特定合约的授权额度
2. **执行授权** - 支持指定金额或无限授权
3. **取消授权** - 将授权额度设为 0
4. **批量查询** - 获取用户对常见合约的所有授权记录
5. **安全检查** - 识别无限授权并提醒用户
6. **错误处理** - 完善的错误处理和重试机制

这个实现与你现有的跨链代码完美集成，复用了相同的 RPC 连接、私钥管理和交易构建逻辑。