package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/sysdevguru/bluelabs/pkg"
	"github.com/sysdevguru/bluelabs/usecase/wallet"
)

type HTTPHandler struct {
	Handle   func(w http.ResponseWriter, r *http.Request) error
	WalletUC *wallet.UseCase
}

// ServeHTTP allows custom handler to satisfy http.Handler
func (h HTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h.Handle(w, r)
	if err != nil {
		switch e := err.(type) {
		case pkg.StatusError:
			http.Error(w, e.Error(), e.Status())
		default:
		}
	}
}

func renderJSON(w http.ResponseWriter, value interface{}) error {
	buffer, err := json.Marshal(value)
	if err != nil {
		return err
	}

	w.Header().Add("Content-Type", "application/json")

	_, err = w.Write(buffer)
	if err != nil {
		return err
	}

	return nil
}
