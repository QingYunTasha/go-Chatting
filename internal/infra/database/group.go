package orm

import (
	dbdomain "go-Chatting/domain/infra/database"

	"gorm.io/gorm"
)

type GroupRepository struct {
	db *gorm.DB
}

func NewGroupRepository(db *gorm.DB) dbdomain.GroupRepository {
	return &GroupRepository{
		db: db,
	}
}

func (r *GroupRepository) Get(ID uint32) (*dbdomain.Group, error) {
	var group dbdomain.Group
	err := r.db.Preload("Members").First(&group, ID).Error
	if err != nil {
		return nil, err
	}

	return &group, nil
}

func (r *GroupRepository) Create(group *dbdomain.Group) error {
	return r.db.Create(group).Error
}

func (r *GroupRepository) Update(ID uint32, group *dbdomain.Group) error {
	return r.db.Model(&dbdomain.Group{}).Where("id = ?", ID).Updates(group).Error
}

func (r *GroupRepository) Delete(ID uint32) error {
	return r.db.Delete(&dbdomain.Group{}, ID).Error
}

func (r *GroupRepository) GetUsers(ID uint32) ([]dbdomain.User, error) {
	var users []dbdomain.User
	if err := r.db.Model(&dbdomain.Group{}).Where("id = ?", ID).Association("Members").Find(&users); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *GroupRepository) GetMessages(ID uint32) ([]dbdomain.GroupMessage, error) {
	var messages []dbdomain.GroupMessage
	if err := r.db.Model(&dbdomain.Group{}).Preload("Sender").Preload("Group").Where("group_id = ?", ID).Find(&messages).Error; err != nil {
		return nil, err
	}

	return messages, nil
}
