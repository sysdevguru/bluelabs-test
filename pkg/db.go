package pkg

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewGormWithPostgres initializes a Postgres database connection.
func NewGormWithPostgres(cfg Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.Database.URL), &gorm.Config{
		Logger: nil,
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConnections)

	return db, nil
}

// CloseDatabaseConnection cleans up the connection to the db.
func CloseDatabaseConnection(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	sqlDB.Close()
	return nil
}
