package pkg

import (
	"context"
	"errors"

	"github.com/sysdevguru/bluelabs/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type GormRepo struct {
	db *gorm.DB
}

func (g *GormRepo) Create(ctx context.Context, userID int64) (*model.Wallet, error) {
	wallet := &model.Wallet{
		UserID:  userID,
		Balance: 0,
	}

	return wallet, g.db.WithContext(ctx).Create(wallet).Error
}

func (g *GormRepo) Deposit(ctx context.Context, userID, walletID int64, funds float64) (*model.Wallet, error) {
	tx := g.db.WithContext(ctx).Begin()
	defer tx.Commit()

	wallet := &model.Wallet{}
	result := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id=?", walletID).
		Where("user_id=?", userID).
		First(wallet)
	if result.Error != nil {
		return wallet, result.Error
	}

	wallet.Balance += funds
	return wallet, tx.Save(wallet).Error
}

func (g *GormRepo) Withdraw(ctx context.Context, userID, walletID int64, funds float64) (*model.Wallet, error) {
	tx := g.db.WithContext(ctx).Begin()
	defer tx.Commit()

	wallet := &model.Wallet{}
	result := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id=?", walletID).
		Where("user_id=?", userID).
		First(wallet)
	if result.Error != nil {
		return wallet, result.Error
	}

	if wallet.Balance < funds {
		return nil, errors.New(ErrWalletBalance)
	}

	wallet.Balance -= funds
	return wallet, tx.Save(wallet).Error
}

func (g *GormRepo) GetWallet(ctx context.Context, userID, walletID int64) (*model.Wallet, error) {
	tx := g.db.WithContext(ctx).Begin()
	defer tx.Commit()

	wallet := &model.Wallet{}
	result := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id=?", walletID).
		Where("user_id=?", userID).
		First(wallet)

	return wallet, result.Error
}

func NewRepo(db *gorm.DB) *GormRepo {
	return &GormRepo{
		db: db,
	}
}
