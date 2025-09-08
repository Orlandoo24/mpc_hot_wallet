package model

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

var ErrNotFound = gorm.ErrRecordNotFound

// WalletsDao defines the interface for database operations on the wallets table.
type WalletsDao interface {
	Insert(ctx context.Context, data *Wallets) error
	FindOneByAddress(ctx context.Context, address string) (*Wallets, error)
	FindAll(ctx context.Context) ([]*Wallets, error)
}

type walletsDao struct {
	db *gorm.DB
}

// NewWalletsDao creates a new instance of WalletsDao.
func NewWalletsDao(db *gorm.DB) WalletsDao {
	return &walletsDao{
		db: db,
	}
}

// Insert adds a new record to the wallets table.
func (d *walletsDao) Insert(ctx context.Context, data *Wallets) error {
	return d.db.WithContext(ctx).Create(data).Error
}

// FindOneByAddress retrieves a single wallet record by its address.
func (d *walletsDao) FindOneByAddress(ctx context.Context, address string) (*Wallets, error) {
	var resp Wallets
	err := d.db.WithContext(ctx).Where("address = ?", address).First(&resp).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &resp, nil
}

// FindAll retrieves all wallet records.
func (d *walletsDao) FindAll(ctx context.Context) ([]*Wallets, error) {
	var wallets []*Wallets
	err := d.db.WithContext(ctx).Find(&wallets).Error
	if err != nil {
		return nil, err
	}
	return wallets, nil
}
