package wallet_test

import (
	"context"
	"os"

	"github.com/sysdevguru/bluelabs/pkg"
	. "github.com/sysdevguru/bluelabs/usecase/wallet"

	. "github.com/onsi/ginkgo"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var _ = Describe("Wallet", func() {
	var (
		uc       *UseCase
		db       *gorm.DB
		ctx      context.Context
		walletID int64
	)

	BeforeEach(func() {
		ctx = context.Background()

		os.Setenv("DATABASE_URL", "postgres://db_user:db_pass@localhost:5434/bluelabs?sslmode=disable")
		cfg, err := pkg.Load()
		assert.NoError(GinkgoT(), err)

		db, err = pkg.NewGormWithPostgres(cfg)
		assert.NoError(GinkgoT(), err)

		uc = New(
			"wallet_task_test",
			pkg.NewRepo(db),
		)
	})

	AfterEach(func() {
		err := pkg.CloseDatabaseConnection(db)
		assert.NoError(GinkgoT(), err)
	})

	Context("Create wallet", func() {
		It("for user 1", func() {
			wallet, err := uc.Create(ctx, 1)
			walletID = wallet.ID
			assert.NoError(GinkgoT(), err)
			assert.Equal(GinkgoT(), 1, int(wallet.UserID))
			assert.Equal(GinkgoT(), 0.00, wallet.Balance)
		})
	})

	Context("Deposit", func() {
		It("from non-existing wallet", func() {
			_, err := uc.Deposit(ctx, 1, 1000, 10)
			assert.Equal(GinkgoT(), "wallet not found", err.Error())
		})

		It("as expected", func() {
			wallet, err := uc.Deposit(ctx, 1, walletID, 100.00)
			assert.NoError(GinkgoT(), err)
			assert.Equal(GinkgoT(), 1, int(wallet.UserID))
			assert.Equal(GinkgoT(), walletID, wallet.ID)
			assert.Equal(GinkgoT(), 100.00, wallet.Balance)
		})
	})

	Context("Withdraw", func() {
		It("from non-existing wallet", func() {
			_, err := uc.Withdraw(ctx, 1, 1000, 10)
			assert.Equal(GinkgoT(), "wallet not found", err.Error())
		})

		It("more than balance", func() {
			_, err := uc.Withdraw(ctx, 1, walletID, 1000.00)
			assert.Equal(GinkgoT(), "wallet funds not enough", err.Error())
		})

		It("as expected", func() {
			wallet, err := uc.Withdraw(ctx, 1, walletID, 35.00)
			assert.NoError(GinkgoT(), err)
			assert.Equal(GinkgoT(), 1, int(wallet.UserID))
			assert.Equal(GinkgoT(), walletID, wallet.ID)
			assert.Equal(GinkgoT(), 65.00, wallet.Balance)
		})
	})

	Context("Get Wallet", func() {
		It("of non-existing user", func() {
			_, err := uc.GetWallet(ctx, 12, walletID)
			assert.Equal(GinkgoT(), "wallet not found", err.Error())
		})

		It("of user 1", func() {
			wallet, err := uc.GetWallet(ctx, 1, walletID)
			assert.NoError(GinkgoT(), err)
			assert.Equal(GinkgoT(), 1, int(wallet.UserID))
			assert.Equal(GinkgoT(), walletID, wallet.ID)
			assert.Equal(GinkgoT(), 65.00, wallet.Balance)
		})
	})
})
