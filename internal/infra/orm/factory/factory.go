package factory

import (
	ormdomain "github.com/QingYunTasha/go-Chatting/domain/infra/orm"
	orm "github.com/QingYunTasha/go-Chatting/internal/infra/orm"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDb(dsn string) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(dsn), &gorm.Config{
		// logger all sql
		//Logger: logger.Default.LogMode(logger.Info),
	})
}

type OrmRepository struct {
	Group          ormdomain.GroupRepository
	User           ormdomain.UserRepository
	GroupMessage   ormdomain.GroupMessageRepository
	PrivateMessage ormdomain.PrivateMessageRepository
}

func NewOrmRepository(db *gorm.DB) (*OrmRepository, error) {
	if err := db.AutoMigrate(&ormdomain.Group{}, &ormdomain.GroupMessage{}, &ormdomain.User{}, &ormdomain.PrivateMessage{}); err != nil {
		return nil, err
	}

	return &OrmRepository{
		Group:          orm.NewGroupRepository(db),
		User:           orm.NewUserRepository(db),
		GroupMessage:   orm.NewGroupMessageRepository(db),
		PrivateMessage: orm.NewPrivateMessageRepository(db),
	}, nil
}
