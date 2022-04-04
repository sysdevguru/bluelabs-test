package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/sysdevguru/bluelabs/model"
	"github.com/sysdevguru/bluelabs/pkg"

	"github.com/gorilla/mux"
)

func (handler *HTTPHandler) GetWallet(w http.ResponseWriter, r *http.Request) error {
	// validate request path params
	userID, err := strconv.Atoi(mux.Vars(r)["userId"])
	if err != nil {
		return pkg.StatusError{
			Code:   http.StatusBadRequest,
			ErrMsg: pkg.ErrUserID,
		}
	}

	walletID, err := strconv.Atoi(mux.Vars(r)["walletId"])
	if err != nil {
		return pkg.StatusError{
			Code:   http.StatusBadRequest,
			ErrMsg: pkg.ErrWalletID,
		}
	}

	wallet, err := handler.WalletUC.GetWallet(r.Context(), int64(userID), int64(walletID))
	if err != nil {
		return err
	}

	return renderJSON(w, wallet)
}

func (handler *HTTPHandler) CreateWallet(w http.ResponseWriter, r *http.Request) error {
	// validate request path params
	userID, err := strconv.Atoi(mux.Vars(r)["userId"])
	if err != nil {
		return pkg.StatusError{
			Code:   http.StatusBadRequest,
			ErrMsg: pkg.ErrUserID,
		}
	}

	wallet, err := handler.WalletUC.Create(r.Context(), int64(userID))
	if err != nil {
		return err
	}

	return renderJSON(w, wallet)
}

func (handler *HTTPHandler) UpdateWallet(w http.ResponseWriter, r *http.Request) error {
	var wallet *model.Wallet
	var err error

	// validate request path params
	userID, err := strconv.Atoi(mux.Vars(r)["userId"])
	if err != nil {
		return pkg.StatusError{
			Code:   http.StatusBadRequest,
			ErrMsg: pkg.ErrUserID,
		}
	}

	walletID, err := strconv.Atoi(mux.Vars(r)["walletId"])
	if err != nil {
		return pkg.StatusError{
			Code:   http.StatusBadRequest,
			ErrMsg: pkg.ErrWalletID,
		}
	}

	transaction := model.Transaction{}
	err = json.NewDecoder(r.Body).Decode(&transaction)
	if err != nil {
		return pkg.StatusError{
			Code:   http.StatusBadRequest,
			ErrMsg: err.Error(),
		}
	}

	// validate input data
	if transaction.Fund < 0 {
		if transaction.Action == model.ActionDeposit {
			return pkg.StatusError{
				Code:   http.StatusBadRequest,
				ErrMsg: pkg.ErrWalletDeposit,
			}
		}

		if transaction.Action == model.ActionWithdraw {
			return pkg.StatusError{
				Code:   http.StatusBadRequest,
				ErrMsg: pkg.ErrWalletWithdraw,
			}
		}
	}

	switch transaction.Action {
	case model.ActionDeposit:
		if wallet, err = handler.WalletUC.Deposit(r.Context(), int64(userID), int64(walletID), transaction.Fund); err != nil {
			return err
		}

		return renderJSON(w, wallet)
	case model.ActionWithdraw:
		if wallet, err = handler.WalletUC.Withdraw(r.Context(), int64(userID), int64(walletID), transaction.Fund); err != nil {
			return err
		}

		return renderJSON(w, wallet)
	default:
		return pkg.StatusError{
			Code:   http.StatusBadRequest,
			ErrMsg: pkg.ErrInvalidAction,
		}
	}
}
