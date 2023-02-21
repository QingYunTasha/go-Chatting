package orm

import (
	ormdomain "github.com/QingYunTasha/go-Chatting/domain/infra/orm"
	"gorm.io/gorm"
)

type GroupMessageRepository struct {
	db *gorm.DB
}

func NewGroupMessageRepository(db *gorm.DB) ormdomain.GroupMessageRepository {
	return &GroupMessageRepository{
		db: db,
	}
}

func (r *GroupMessageRepository) Get(ID uint32) (*ormdomain.GroupMessage, error) {
	var message ormdomain.GroupMessage
	err := r.db.First(&message, ID).Error
	if err != nil {
		return nil, err
	}
	return &message, nil
}

func (r *GroupMessageRepository) Create(message *ormdomain.GroupMessage) error {
	return r.db.Create(message).Error
}

func (r *GroupMessageRepository) Delete(ID uint32) error {
	return r.db.Delete(&ormdomain.GroupMessage{}, ID).Error
}

func (r *GroupMessageRepository) DeleteByGroup(ID uint32) error {
	return r.db.Where("group_id = ?", ID).Delete(&ormdomain.GroupMessage{}).Error
}
