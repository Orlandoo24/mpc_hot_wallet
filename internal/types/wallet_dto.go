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
