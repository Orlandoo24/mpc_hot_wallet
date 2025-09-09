package types

// WalletInitReq defines the request body for initializing a new wallet.
// 钱包将根据系统配置自动为多个链创建地址
type WalletInitReq struct {
	// A user-defined name for the wallet.
	Name string `json:"name"`
	// The user's phone number (optional).
	PhoneNumber string `json:"phone_number,omitempty"`
	// The user's email address (optional).
	Email string `json:"email,omitempty" validate:"omitempty,email"`
}

// WalletAddress 单个链的钱包地址信息
type WalletAddress struct {
	// The blockchain the wallet belongs to.
	Chain string `json:"chain"`
	// The public address of the newly created wallet.
	Address string `json:"address"`
}

// WalletInitResp defines the response body for a successful wallet initialization.
type WalletInitResp struct {
	// 所有创建的钱包地址列表
	Wallets []WalletAddress `json:"wallets"`
	// 钱包总数
	TotalCount int `json:"total_count"`
	// 创建成功的链数量
	SuccessCount int `json:"success_count"`
	// 失败的链（如果有）
	FailedChains []string `json:"failed_chains,omitempty"`
}
