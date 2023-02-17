package factory

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDb(dsn string) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(dsn), &gorm.Config{
		// logger all sql
		//Logger: logger.Default.LogMode(logger.Info),
	})
}

type OrmRepository struct{}

func InitOrmRepository(db *gorm.DB) (*OrmRepository, error) {
	if err := db.AutoMigrate(); err != nil {
		return nil, err
	}

	return &OrmRepository{}, nil
}
