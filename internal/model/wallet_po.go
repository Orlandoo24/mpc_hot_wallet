package model

import (
	"database/sql"
	"time"
)

// Wallets corresponds to the wallets table in the database.
type Wallets struct {
	Id                  int64          `db:"id"`
	UserId              string         `db:"user_id"`
	Address             string         `db:"address"`
	EncryptedPrivateKey string         `db:"encrypted_private_key"`
	PhoneNumber         sql.NullString `db:"phone_number"`
	Email               sql.NullString `db:"email"`
	CreatedAt           time.Time      `db:"created_at"`
	UpdatedAt           time.Time      `db:"updated_at"`
	ChainType           sql.NullString `db:"chain_type"`
}
