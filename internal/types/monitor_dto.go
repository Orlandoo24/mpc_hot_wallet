package types

type TokenEvent struct {
	BlockNumber uint64 `json:"blockNumber"`
	TxHash      string `json:"txHash"`
	Timestamp   int64  `json:"timestamp"`
	EventType   string `json:"eventType"` // send/receive/swap/bridge/approve/revoke
	FromAddr    string `json:"fromAddr"`
	ToAddr      string `json:"toAddr"`
	TokenAddr   string `json:"tokenAddr"` // 代币合约地址
	Amount      string `json:"amount"`    // 使用string存储以避免精度问题
	ChainId     uint64 `json:"chainId"`   // 支持跨链场景
}
