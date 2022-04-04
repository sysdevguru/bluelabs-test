package pkg_test

import (
	"os"
	"time"

	. "github.com/sysdevguru/bluelabs/pkg"

	. "github.com/onsi/ginkgo"
	"github.com/stretchr/testify/assert"
)

var _ = Describe("Config", func() {
	Context("Load config", func() {
		It("without database url", func() {
			_, err := Load()
			assert.Equal(GinkgoT(), "required key DATABASE_URL missing value", err.Error())
		})

		It("with database url", func() {
			os.Setenv("DATABASE_URL", "http://localhost:5432")
			cfg, err := Load()

			assert.NoError(GinkgoT(), err)
			assert.Equal(GinkgoT(), 60*time.Second, cfg.Server.IdleTimeout)
			assert.Equal(GinkgoT(), 8080, cfg.Server.Port)
			assert.Equal(GinkgoT(), 1*time.Second, cfg.Server.ReadTimeout)
			assert.Equal(GinkgoT(), 2*time.Second, cfg.Server.WriteTimeout)
			assert.Equal(GinkgoT(), "http://localhost:5432", cfg.Database.URL)
			assert.Equal(GinkgoT(), "warn", cfg.Database.LogLevel)
			assert.Equal(GinkgoT(), 10, cfg.Database.MaxOpenConnections)
		})
	})
})
