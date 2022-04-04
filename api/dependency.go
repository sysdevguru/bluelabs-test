package api

import (
	"github.com/sysdevguru/bluelabs/pkg"
	"github.com/sysdevguru/bluelabs/usecase/wallet"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type Service struct {
	cfg      pkg.Config
	db       *gorm.DB
	walletUC *wallet.UseCase
}

func NewService(cfg pkg.Config) (*Service, error) {
	db, err := pkg.NewGormWithPostgres(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect database")
	}

	walletUC := wallet.New(
		"wallet_task",
		pkg.NewRepo(db),
	)

	return &Service{
		cfg,
		db,
		walletUC,
	}, nil
}

func (s *Service) Shutdown() error {
	return pkg.CloseDatabaseConnection(s.db)
}
