package api_test

import (
	"os"

	. "github.com/sysdevguru/bluelabs/api"
	"github.com/sysdevguru/bluelabs/pkg"

	. "github.com/onsi/ginkgo"
	"github.com/stretchr/testify/assert"
)

var _ = Describe("Dependency", func() {
	var (
		service *Service
	)

	Context("Dependency", func() {
		It("create service", func() {
			os.Setenv("DATABASE_URL", "postgres://db_user:db_pass@localhost:5434/bluelabs?sslmode=disable")

			cfg, err := pkg.Load()
			assert.NoError(GinkgoT(), err)

			service, err = NewService(cfg)
			assert.NoError(GinkgoT(), err)
		})

		It("stop service", func() {
			err := service.Shutdown()
			assert.NoError(GinkgoT(), err)
		})
	})
})
