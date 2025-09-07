package types

// TransactionReq defines the request for various transaction operations (send, swap, etc.)
type TransactionReq struct {
	FromAddress string `json:"from_address" validate:"required"`
	ToAddress   string `json:"to_address" validate:"required"`
	Chain       string `json:"chain" validate:"required"`      // e.g., "BSC"
	FromToken   string `json:"from_token" validate:"required"` // e.g., "0x55d398326f99059fF775485246999027B3197955" for USDT
	ToToken     string `json:"to_token" validate:"required"`   // e.g., "0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE" for native BNB
	Amount      string `json:"amount" validate:"required"`     // e.g., "1000000000000000000" for 1 USDT
}

// TransactionResp defines the response for transaction operations.
type TransactionResp struct {
	TxHash      string `json:"tx_hash"`
	Message     string `json:"message"`
	ExplorerUrl string `json:"explorer_url"`
	Chain       string `json:"chain"`
	Status      string `json:"status"`
}

// LifiToken LI.FI API 中的代币信息
type LifiToken struct {
	Address  string `json:"address"`
	Symbol   string `json:"symbol"`
	Decimals int    `json:"decimals"`
	ChainId  int    `json:"chainId"`
	Name     string `json:"name"`
	PriceUSD string `json:"priceUSD,omitempty"`
}

// LifiGasCost LI.FI API 中的 Gas 费用信息
type LifiGasCost struct {
	Type     string    `json:"type"`
	Price    string    `json:"price"`
	Estimate string    `json:"estimate"`
	Limit    string    `json:"limit"`
	Token    LifiToken `json:"token"` // 修改为对象类型
}

type LifiQuoteResponse struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Tool     string `json:"tool"`
	Estimate struct {
		FromAmount      string        `json:"fromAmount"`
		ToAmount        string        `json:"toAmount"`
		ToAmountMin     string        `json:"toAmountMin"`
		ApprovalAddress string        `json:"approvalAddress"`
		GasCosts        []LifiGasCost `json:"gasCosts"` // 使用新的结构体
	} `json:"estimate"`
	TransactionRequest struct {
		Data     string `json:"data"`
		To       string `json:"to"`
		Value    string `json:"value"`
		GasLimit string `json:"gasLimit"`
		GasPrice string `json:"gasPrice"`
	} `json:"transactionRequest"`
}

// BridgeQuoteReq 跨链报价请求
type BridgeQuoteReq struct {
	FromChain   int    `json:"from_chain" validate:"required"`   // 源链ID，如 56 (BSC)
	ToChain     int    `json:"to_chain" validate:"required"`     // 目标链ID，如 8453 (Base)
	FromToken   string `json:"from_token" validate:"required"`   // 源代币地址，原生币用 0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE
	ToToken     string `json:"to_token" validate:"required"`     // 目标代币地址
	FromAmount  string `json:"from_amount" validate:"required"`  // 转账金额
	FromAddress string `json:"from_address" validate:"required"` // 发送地址
	ToAddress   string `json:"to_address" validate:"required"`   // 接收地址
	Order       string `json:"order,omitempty"`                  // 排序方式: FASTEST, CHEAPEST, SAFEST
	Slippage    string `json:"slippage,omitempty"`               // 滑点，如 "0.005" 表示 0.5%
}

// BridgeQuoteResp 跨链报价响应
type BridgeQuoteResp struct {
	TransactionRequest BridgeTxRequest `json:"transaction_request"`
	Estimate           BridgeEstimate  `json:"estimate"`
	Message            string          `json:"message"`
}

// BridgeTxRequest 跨链交易请求
type BridgeTxRequest struct {
	To       string `json:"to"`
	Data     string `json:"data"`
	Value    string `json:"value"`
	GasLimit string `json:"gasLimit"`
	GasPrice string `json:"gasPrice"`
}

// BridgeEstimate 跨链估算信息
type BridgeEstimate struct {
	GasCosts        []BridgeGasCost `json:"gasCosts"`
	ApprovalAddress string          `json:"approvalAddress,omitempty"`
}

// BridgeGasCost Gas费用信息
type BridgeGasCost struct {
	Type     string    `json:"type"`
	Price    string    `json:"price"`
	Amount   string    `json:"amount"`
	Token    LifiToken `json:"token"` // 修改为对象类型，复用 LifiToken 结构体
	Decimals int       `json:"decimals"`
}

// BridgeExecuteReq 执行跨链转账请求
type BridgeExecuteReq struct {
	FromAddress string `json:"from_address" validate:"required"`
	FromChain   int    `json:"from_chain" validate:"required"`
	ToChain     int    `json:"to_chain" validate:"required"`
	FromToken   string `json:"from_token" validate:"required"`
	ToToken     string `json:"to_token" validate:"required"`
	Amount      string `json:"amount" validate:"required"`
	ToAddress   string `json:"to_address" validate:"required"`
	Order       string `json:"order,omitempty"`
	Slippage    string `json:"slippage,omitempty"`
}

// BridgeExecuteResp 执行跨链转账响应
type BridgeExecuteResp struct {
	TxHash      string `json:"tx_hash"`
	Message     string `json:"message"`
	ExplorerUrl string `json:"explorer_url"`
	FromChain   int    `json:"from_chain"`
	ToChain     int    `json:"to_chain"`
	Status      string `json:"status"`
}

// BridgeStatusReq 查询跨链状态请求
type BridgeStatusReq struct {
	TxHash string `json:"tx_hash" validate:"required"`
}

// BridgeStatusResp 查询跨链状态响应
type BridgeStatusResp struct {
	Status     string `json:"status"`
	SubStatus  string `json:"sub_status,omitempty"`
	TxHash     string `json:"tx_hash"`
	FromTxLink string `json:"from_tx_link,omitempty"`
	ToTxLink   string `json:"to_tx_link,omitempty"`
	Message    string `json:"message"`
}

// CheckAllowanceReq 检查授权请求
type CheckAllowanceReq struct {
	TokenAddress   string `json:"token_address" validate:"required"`
	OwnerAddress   string `json:"owner_address" validate:"required"`
	SpenderAddress string `json:"spender_address" validate:"required"`
	Chain          string `json:"chain" validate:"required"`
}

// CheckAllowanceResp 检查授权响应
type CheckAllowanceResp struct {
	TokenAddress   string `json:"token_address"`
	OwnerAddress   string `json:"owner_address"`
	SpenderAddress string `json:"spender_address"`
	Allowance      string `json:"allowance"`
	Chain          string `json:"chain"`
	Message        string `json:"message"`
}

// ApproveTokenReq 授权代币请求
type ApproveTokenReq struct {
	TokenAddress   string `json:"token_address" validate:"required"`
	SpenderAddress string `json:"spender_address" validate:"required"`
	OwnerAddress   string `json:"owner_address" validate:"required"`
	Amount         string `json:"amount,omitempty"` // "max" 表示无限授权，空值也表示无限授权
	Chain          string `json:"chain" validate:"required"`
}

// ApproveTokenResp 授权代币响应
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

// RevokeApprovalReq 取消授权请求
type RevokeApprovalReq struct {
	TokenAddress   string `json:"token_address" validate:"required"`
	SpenderAddress string `json:"spender_address" validate:"required"`
	OwnerAddress   string `json:"owner_address" validate:"required"`
	Chain          string `json:"chain" validate:"required"`
}

// RevokeApprovalResp 取消授权响应
type RevokeApprovalResp struct {
	TxHash         string `json:"tx_hash"`
	TokenAddress   string `json:"token_address"`
	SpenderAddress string `json:"spender_address"`
	Chain          string `json:"chain"`
	ExplorerUrl    string `json:"explorer_url"`
	Message        string `json:"message"`
	Status         string `json:"status"`
}

// GetUserApprovalsReq 获取用户授权请求
type GetUserApprovalsReq struct {
	UserAddress    string   `json:"user_address" validate:"required"`
	TokenAddresses []string `json:"token_addresses" validate:"required"`
	Chain          string   `json:"chain" validate:"required"`
}

// GetUserApprovalsResp 获取用户授权响应
type GetUserApprovalsResp struct {
	UserAddress string         `json:"user_address"`
	Chain       string         `json:"chain"`
	Approvals   []ApprovalInfo `json:"approvals"`
	Message     string         `json:"message"`
}

// ApprovalInfo 授权信息
type ApprovalInfo struct {
	TokenAddress   string `json:"token_address"`
	SpenderAddress string `json:"spender_address"`
	SpenderName    string `json:"spender_name"`
	Allowance      string `json:"allowance"`
	IsUnlimited    bool   `json:"is_unlimited"`
}

// SpenderInfo Spender 信息
type SpenderInfo struct {
	Address string `json:"address"`
	Name    string `json:"name"`
}
