package api

import (
	"net/http"

	"github.com/sysdevguru/bluelabs/api/handlers"

	"github.com/gorilla/mux"
)

func NewRouter(service *Service) *mux.Router {
	handler := handlers.HTTPHandler{WalletUC: service.walletUC}

	r := mux.NewRouter()
	r.Handle("/users/{userId}/wallets/{walletId}", handlers.HTTPHandler{Handle: handler.GetWallet}).Methods(http.MethodGet)
	r.Handle("/users/{userId}/wallets", handlers.HTTPHandler{Handle: handler.CreateWallet}).Methods(http.MethodPost)
	r.Handle("/users/{userId}/wallets/{walletId}", handlers.HTTPHandler{Handle: handler.UpdateWallet}).Methods(http.MethodPut)

	return r
}
