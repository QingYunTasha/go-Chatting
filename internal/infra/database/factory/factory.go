package factory

import (
	dbdomain "go-Chatting/domain/infra/database"

	database "go-Chatting/internal/infra/database"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDb(dsn string) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(dsn), &gorm.Config{
		// logger all sql
		//Logger: logger.Default.LogMode(logger.Info),
	})
}

type DatabaseRepository struct {
	Group          dbdomain.GroupRepository
	User           dbdomain.UserRepository
	GroupMessage   dbdomain.GroupMessageRepository
	PrivateMessage dbdomain.PrivateMessageRepository
}

func NewDbRepository(db *gorm.DB) (*DatabaseRepository, error) {
	if err := db.AutoMigrate(&dbdomain.Group{}, &dbdomain.GroupMessage{}, &dbdomain.User{}, &dbdomain.PrivateMessage{}); err != nil {
		return nil, err
	}

	return &DatabaseRepository{
		Group:          database.NewGroupRepository(db),
		User:           database.NewUserRepository(db),
		GroupMessage:   database.NewGroupMessageRepository(db),
		PrivateMessage: database.NewPrivateMessageRepository(db),
	}, nil
}
