package orm

import (
	ormdomain "github.com/QingYunTasha/go-Chatting/domain/infra/orm"
	"gorm.io/gorm"
)

type GroupRepository struct {
	db *gorm.DB
}

func NewGroupRepository(db *gorm.DB) ormdomain.GroupRepository {
	return &GroupRepository{
		db: db,
	}
}

func (r *GroupRepository) Get(ID uint32) (*ormdomain.Group, error) {
	var group ormdomain.Group
	err := r.db.Preload("Members").First(&group, ID).Error
	if err != nil {
		return nil, err
	}

	return &group, nil
}

func (r *GroupRepository) Create(group *ormdomain.Group) error {
	return r.db.Create(group).Error
}

func (r *GroupRepository) Update(ID uint32, group *ormdomain.Group) error {
	return r.db.Model(&ormdomain.Group{}).Where("id = ?", ID).Updates(group).Error
}

func (r *GroupRepository) Delete(ID uint32) error {
	return r.db.Delete(&ormdomain.Group{}, ID).Error
}

func (r *GroupRepository) GetUsers(ID uint32) ([]ormdomain.User, error) {
	var users []ormdomain.User
	if err := r.db.Model(&ormdomain.Group{}).Where("id = ?", ID).Association("Members").Find(&users); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *GroupRepository) GetMessages(ID uint32) ([]ormdomain.GroupMessage, error) {
	var messages []ormdomain.GroupMessage
	if err := r.db.Model(&ormdomain.Group{}).Preload("Sender").Preload("Group").Where("group_id = ?", ID).Find(&messages).Error; err != nil {
		return nil, err
	}

	return messages, nil
}
