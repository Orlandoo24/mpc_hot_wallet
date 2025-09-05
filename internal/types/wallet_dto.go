package types

// WalletInitReq defines the request body for initializing a new wallet.
type WalletInitReq struct {
	// The blockchain to initialize the wallet for.
	Chain string `json:"chain" validate:"required"`
	// A user-defined name for the wallet.
	Name string `json:"name"`
	// The user's phone number (optional).
	PhoneNumber string `json:"phone_number,omitempty"`
	// The user's email address (optional).
	Email string `json:"email,omitempty" validate:"omitempty,email"`
}

// WalletInitResp defines the response body for a successful wallet initialization.
type WalletInitResp struct {
	// The public address of the newly created wallet.
	Address string `json:"address"`
	// The blockchain the wallet belongs to.
	Chain string `json:"chain"`
}

// TransactionSendReq defines the request for sending a transaction via LI.FI
type TransactionSendReq struct {
	FromAddress string `json:"from_address" validate:"required"`
	ToAddress   string `json:"to_address" validate:"required"`
	FromChain   string `json:"from_chain" validate:"required"` // e.g., "BSC"
	ToChain     string `json:"to_chain" validate:"required"`   // e.g., "BSC"
	FromToken   string `json:"from_token" validate:"required"` // e.g., "0x55d398326f99059fF775485246999027B3197955" for USDT
	ToToken     string `json:"to_token" validate:"required"`   // e.g., "0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE" for native BNB
	Amount      string `json:"amount" validate:"required"`     // e.g., "1000000000000000000" for 1 USDT
}

// TransactionSendResp defines the response for a send transaction request.
type TransactionSendResp struct {
	TxHash string `json:"tx_hash"`
}
