package orm

import (
	dbdomain "go-Chatting/domain/infra/database"

	"gorm.io/gorm"
)

type GroupMessageRepository struct {
	db *gorm.DB
}

func NewGroupMessageRepository(db *gorm.DB) dbdomain.GroupMessageRepository {
	return &GroupMessageRepository{
		db: db,
	}
}

func (r *GroupMessageRepository) Get(ID uint32) (*dbdomain.GroupMessage, error) {
	var message dbdomain.GroupMessage
	err := r.db.First(&message, ID).Error
	if err != nil {
		return nil, err
	}
	return &message, nil
}

func (r *GroupMessageRepository) Create(message *dbdomain.GroupMessage) error {
	return r.db.Create(message).Error
}

func (r *GroupMessageRepository) Delete(ID uint32) error {
	return r.db.Delete(&dbdomain.GroupMessage{}, ID).Error
}

func (r *GroupMessageRepository) DeleteByGroup(ID uint32) error {
	return r.db.Where("group_id = ?", ID).Delete(&dbdomain.GroupMessage{}).Error
}

func (r *GroupMessageRepository) GetByGroupAfterOrderID(groupID uint32, lastOrderID uint32) ([]dbdomain.GroupMessage, error) {
	var groupMessages []dbdomain.GroupMessage
	if err := r.db.Where("(group_id =?) AND (order_id > ?)", groupID, lastOrderID).Find(&groupMessages).Error; err != nil {
		return nil, err
	}
	return groupMessages, nil
}

func (r *GroupMessageRepository) GetLastIDByGroup(groupID uint32) (uint32, error) {
	var maxOrderID uint32
	if err := r.db.Where("group_id =?", groupID).Select("MAX(order_id)").Row().Scan(&maxOrderID); err != nil {
		return 0, err
	}

	return maxOrderID, nil
}
