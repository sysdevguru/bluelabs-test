package wallet

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/sysdevguru/bluelabs/model"
	"github.com/sysdevguru/bluelabs/pkg"
	"gorm.io/gorm"
)

type Repo interface {
	Create(ctx context.Context, userID int64) (*model.Wallet, error)
	Deposit(ctx context.Context, userID, walletID int64, funds float64) (*model.Wallet, error)
	Withdraw(ctx context.Context, userID, walletID int64, funds float64) (*model.Wallet, error)
	GetWallet(ctx context.Context, userID, walletID int64) (*model.Wallet, error)
}

type UseCase struct {
	taskName string
	repo     Repo
}

func (uc *UseCase) Create(
	ctx context.Context,
	userID int64,
) (*model.Wallet, error) {
	wallet, err := uc.repo.Create(ctx, userID)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return nil, pkg.StatusError{
				Code:   http.StatusConflict,
				ErrMsg: pkg.ErrDuplicated,
			}
		}

		return nil, pkg.StatusError{
			Code:   http.StatusInternalServerError,
			ErrMsg: err.Error(),
		}
	}

	return wallet, nil
}

func (uc *UseCase) GetWallet(
	ctx context.Context,
	userID, walletID int64,
) (*model.Wallet, error) {
	wallet, err := uc.repo.GetWallet(ctx, userID, walletID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pkg.StatusError{
				Code:   http.StatusNotFound,
				ErrMsg: pkg.ErrWalletNotFound,
			}
		}

		return nil, pkg.StatusError{
			Code:   http.StatusInternalServerError,
			ErrMsg: err.Error(),
		}
	}

	return wallet, nil
}

func (uc *UseCase) Deposit(
	ctx context.Context,
	userID, walletID int64,
	funds float64,
) (*model.Wallet, error) {
	wallet, err := uc.repo.Deposit(ctx, userID, walletID, funds)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pkg.StatusError{
				Code:   http.StatusNotFound,
				ErrMsg: pkg.ErrWalletNotFound,
			}
		}

		return nil, pkg.StatusError{
			Code:   http.StatusInternalServerError,
			ErrMsg: err.Error(),
		}
	}

	return wallet, nil
}

func (uc *UseCase) Withdraw(
	ctx context.Context,
	userID, walletID int64,
	funds float64,
) (*model.Wallet, error) {
	wallet, err := uc.repo.Withdraw(ctx, userID, walletID, funds)
	if err != nil {
		if err.Error() == pkg.ErrWalletFunds {
			return nil, pkg.StatusError{
				Code:   http.StatusBadRequest,
				ErrMsg: err.Error(),
			}
		}

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pkg.StatusError{
				Code:   http.StatusNotFound,
				ErrMsg: pkg.ErrWalletNotFound,
			}
		}

		return nil, pkg.StatusError{
			Code:   http.StatusInternalServerError,
			ErrMsg: err.Error(),
		}
	}

	return wallet, nil
}

func New(taskName string, repo Repo) *UseCase {
	return &UseCase{
		taskName,
		repo,
	}
}
