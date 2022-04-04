package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"

	. "github.com/sysdevguru/bluelabs/api"
	"github.com/sysdevguru/bluelabs/model"
	"github.com/sysdevguru/bluelabs/pkg"

	"github.com/gorilla/mux"
	. "github.com/onsi/ginkgo"
	"github.com/stretchr/testify/assert"
)

var _ = Describe("Router", func() {
	var (
		walletID   int64
		router     *mux.Router
		runRequest func(srv http.Handler, r *http.Request) *httptest.ResponseRecorder
		getPayload func(data interface{}) (io.Reader, error)
	)

	BeforeEach(func() {
		os.Setenv("DATABASE_URL", "postgres://db_user:db_pass@localhost:5434/bluelabs?sslmode=disable")
		cfg, err := pkg.Load()
		assert.NoError(GinkgoT(), err)

		service, err := NewService(cfg)
		assert.NoError(GinkgoT(), err)

		router = NewRouter(service)

		runRequest = func(srv http.Handler, r *http.Request) *httptest.ResponseRecorder {
			response := httptest.NewRecorder()
			r.Header.Set("content-type", "application/json")
			srv.ServeHTTP(response, r)
			return response
		}

		getPayload = func(data interface{}) (io.Reader, error) {
			var buf bytes.Buffer
			err := json.NewEncoder(&buf).Encode(data)
			return &buf, err
		}
	})

	Context("Create", func() {
		It("with invalid user id", func() {
			createWalletReq := httptest.NewRequest("POST", fmt.Sprintf("/users/%s/wallets", "invalid_user"), nil)
			createWalletResp := runRequest(router, createWalletReq)

			assert.Equal(GinkgoT(), 400, createWalletResp.Code)
			assert.Equal(GinkgoT(), "invalid user id\n", createWalletResp.Body.String())
		})

		It("as expected", func() {
			createWalletReq := httptest.NewRequest("POST", fmt.Sprintf("/users/%d/wallets", 2), nil)
			createWalletResp := runRequest(router, createWalletReq)

			var wallet model.Wallet
			decoder := json.NewDecoder(createWalletResp.Body)
			err := decoder.Decode(&wallet)
			assert.NoError(GinkgoT(), err)

			assert.Equal(GinkgoT(), 2, int(wallet.UserID))
			assert.Equal(GinkgoT(), 0.00, wallet.Balance)

			walletID = wallet.ID
		})

		It("another wallet for user", func() {
			createWalletReq := httptest.NewRequest("POST", fmt.Sprintf("/users/%d/wallets", 2), nil)
			createWalletResp := runRequest(router, createWalletReq)

			assert.Equal(GinkgoT(), 409, createWalletResp.Code)
			assert.Equal(GinkgoT(), "user already has a wallet\n", createWalletResp.Body.String())
		})
	})

	Context("Get wallet", func() {
		It("with invalid user id", func() {
			getWalletReq := httptest.NewRequest("GET", fmt.Sprintf("/users/%s/wallets/%d", "invalid_user", walletID), nil)
			getWalletResp := runRequest(router, getWalletReq)

			assert.Equal(GinkgoT(), 400, getWalletResp.Code)
			assert.Equal(GinkgoT(), "invalid user id\n", getWalletResp.Body.String())
		})

		It("with invalid wallet id", func() {
			getWalletReq := httptest.NewRequest("GET", fmt.Sprintf("/users/%d/wallets/%s", 2, "invalid_wallet"), nil)
			getWalletResp := runRequest(router, getWalletReq)

			assert.Equal(GinkgoT(), 400, getWalletResp.Code)
			assert.Equal(GinkgoT(), "invalid wallet id\n", getWalletResp.Body.String())
		})

		It("with mismatching user/wallet", func() {
			getWalletReq := httptest.NewRequest("GET", fmt.Sprintf("/users/%d/wallets/%d", 3, walletID), nil)
			getWalletResp := runRequest(router, getWalletReq)

			assert.Equal(GinkgoT(), 404, getWalletResp.Code)
			assert.Equal(GinkgoT(), "wallet not found\n", getWalletResp.Body.String())
		})

		It("as expect", func() {
			getWalletReq := httptest.NewRequest("GET", fmt.Sprintf("/users/%d/wallets/%d", 2, walletID), nil)
			getWalletResp := runRequest(router, getWalletReq)

			var wallet model.Wallet
			decoder := json.NewDecoder(getWalletResp.Body)
			err := decoder.Decode(&wallet)
			assert.NoError(GinkgoT(), err)

			assert.Equal(GinkgoT(), 2, int(wallet.UserID))
			assert.Equal(GinkgoT(), walletID, wallet.ID)
			assert.Equal(GinkgoT(), 0.00, wallet.Balance)
		})
	})

	Context("Deposit", func() {
		It("with invalid user id", func() {
			updateWalletReq := httptest.NewRequest("PUT", fmt.Sprintf("/users/%s/wallets/%d", "invalid_user", walletID), nil)
			updateWalletResp := runRequest(router, updateWalletReq)

			assert.Equal(GinkgoT(), 400, updateWalletResp.Code)
			assert.Equal(GinkgoT(), "invalid user id\n", updateWalletResp.Body.String())
		})

		It("with invalid wallet id", func() {
			updateWalletReq := httptest.NewRequest("PUT", fmt.Sprintf("/users/%d/wallets/%s", 2, "invalid_wallet"), nil)
			updateWalletResp := runRequest(router, updateWalletReq)

			assert.Equal(GinkgoT(), 400, updateWalletResp.Code)
			assert.Equal(GinkgoT(), "invalid wallet id\n", updateWalletResp.Body.String())
		})

		It("with mismatching user/wallet", func() {
			reqData := model.Transaction{
				Action: "deposit",
				Fund:   100.00,
			}
			payload, err := getPayload(reqData)
			assert.NoError(GinkgoT(), err)

			updateWalletReq := httptest.NewRequest("PUT", fmt.Sprintf("/users/%d/wallets/%d", 3, walletID), payload)
			updateWalletResp := runRequest(router, updateWalletReq)

			assert.Equal(GinkgoT(), 404, updateWalletResp.Code)
			assert.Equal(GinkgoT(), "wallet not found\n", updateWalletResp.Body.String())
		})

		It("with invalid funds", func() {
			reqData := model.Transaction{
				Action: "deposit",
				Fund:   -100.00,
			}
			payload, err := getPayload(reqData)
			assert.NoError(GinkgoT(), err)

			updateWalletReq := httptest.NewRequest("PUT", fmt.Sprintf("/users/%d/wallets/%d", 2, walletID), payload)
			updateWalletResp := runRequest(router, updateWalletReq)

			assert.Equal(GinkgoT(), 400, updateWalletResp.Code)
			assert.Equal(GinkgoT(), "cannot update balance with nagetive fund\n", updateWalletResp.Body.String())
		})

		It("as expected", func() {
			reqData := model.Transaction{
				Action: "deposit",
				Fund:   100.00,
			}
			payload, err := getPayload(reqData)
			assert.NoError(GinkgoT(), err)

			updateWalletReq := httptest.NewRequest("PUT", fmt.Sprintf("/users/%d/wallets/%d", 2, walletID), payload)
			updateWalletResp := runRequest(router, updateWalletReq)

			var wallet model.Wallet
			decoder := json.NewDecoder(updateWalletResp.Body)
			err = decoder.Decode(&wallet)
			assert.NoError(GinkgoT(), err)
			assert.Equal(GinkgoT(), 200, updateWalletResp.Code)
			assert.Equal(GinkgoT(), 2, int(wallet.UserID))
			assert.Equal(GinkgoT(), walletID, wallet.ID)
			assert.Equal(GinkgoT(), 100.00, wallet.Balance)
		})

		It("confirm the update", func() {
			getWalletReq := httptest.NewRequest("GET", fmt.Sprintf("/users/%d/wallets/%d", 2, walletID), nil)
			getWalletResp := runRequest(router, getWalletReq)

			var wallet model.Wallet
			decoder := json.NewDecoder(getWalletResp.Body)
			err := decoder.Decode(&wallet)
			assert.NoError(GinkgoT(), err)

			assert.Equal(GinkgoT(), 2, int(wallet.UserID))
			assert.Equal(GinkgoT(), walletID, wallet.ID)
			assert.Equal(GinkgoT(), 100.00, wallet.Balance)
		})
	})

	Context("Withdraw", func() {
		It("with invalid user id", func() {
			updateWalletReq := httptest.NewRequest("PUT", fmt.Sprintf("/users/%s/wallets/%d", "invalid_user", walletID), nil)
			updateWalletResp := runRequest(router, updateWalletReq)

			assert.Equal(GinkgoT(), 400, updateWalletResp.Code)
			assert.Equal(GinkgoT(), "invalid user id\n", updateWalletResp.Body.String())
		})

		It("with invalid wallet id", func() {
			updateWalletReq := httptest.NewRequest("PUT", fmt.Sprintf("/users/%d/wallets/%s", 2, "invalid_wallet"), nil)
			updateWalletResp := runRequest(router, updateWalletReq)

			assert.Equal(GinkgoT(), 400, updateWalletResp.Code)
			assert.Equal(GinkgoT(), "invalid wallet id\n", updateWalletResp.Body.String())
		})

		It("with mismatching user/wallet", func() {
			reqData := model.Transaction{
				Action: "withdraw",
				Fund:   100.00,
			}
			payload, err := getPayload(reqData)
			assert.NoError(GinkgoT(), err)

			updateWalletReq := httptest.NewRequest("PUT", fmt.Sprintf("/users/%d/wallets/%d", 3, walletID), payload)
			updateWalletResp := runRequest(router, updateWalletReq)

			assert.Equal(GinkgoT(), 404, updateWalletResp.Code)
			assert.Equal(GinkgoT(), "wallet not found\n", updateWalletResp.Body.String())
		})

		It("with invalid funds", func() {
			reqData := model.Transaction{
				Action: "withdraw",
				Fund:   -100.00,
			}
			payload, err := getPayload(reqData)
			assert.NoError(GinkgoT(), err)

			updateWalletReq := httptest.NewRequest("PUT", fmt.Sprintf("/users/%d/wallets/%d", 2, walletID), payload)
			updateWalletResp := runRequest(router, updateWalletReq)

			assert.Equal(GinkgoT(), 400, updateWalletResp.Code)
			assert.Equal(GinkgoT(), "cannot update balance with nagetive fund\n", updateWalletResp.Body.String())
		})

		It("as expected", func() {
			reqData := model.Transaction{
				Action: "withdraw",
				Fund:   45.00,
			}
			payload, err := getPayload(reqData)
			assert.NoError(GinkgoT(), err)

			updateWalletReq := httptest.NewRequest("PUT", fmt.Sprintf("/users/%d/wallets/%d", 2, walletID), payload)
			updateWalletResp := runRequest(router, updateWalletReq)

			var wallet model.Wallet
			decoder := json.NewDecoder(updateWalletResp.Body)
			err = decoder.Decode(&wallet)
			assert.NoError(GinkgoT(), err)
			assert.Equal(GinkgoT(), 200, updateWalletResp.Code)
			assert.Equal(GinkgoT(), 2, int(wallet.UserID))
			assert.Equal(GinkgoT(), walletID, wallet.ID)
			assert.Equal(GinkgoT(), 55.00, wallet.Balance)
		})

		It("confirm the update", func() {
			getWalletReq := httptest.NewRequest("GET", fmt.Sprintf("/users/%d/wallets/%d", 2, walletID), nil)
			getWalletResp := runRequest(router, getWalletReq)

			var wallet model.Wallet
			decoder := json.NewDecoder(getWalletResp.Body)
			err := decoder.Decode(&wallet)
			assert.NoError(GinkgoT(), err)

			assert.Equal(GinkgoT(), 2, int(wallet.UserID))
			assert.Equal(GinkgoT(), walletID, wallet.ID)
			assert.Equal(GinkgoT(), 55.00, wallet.Balance)
		})
	})
})
