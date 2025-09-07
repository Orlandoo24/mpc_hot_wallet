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
	Type     string `json:"type"`
	Price    string `json:"price"`
	Amount   string `json:"amount"`
	Token    string `json:"token"`
	Decimals int    `json:"decimals"`
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
