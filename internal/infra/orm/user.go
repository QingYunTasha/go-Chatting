package orm

import (
	ormdomain "github.com/QingYunTasha/go-Chatting/domain/infra/orm"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) ormdomain.UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (u *UserRepository) Get(ID uint32) (*ormdomain.User, error) {
	var user ormdomain.User
	err := u.db.Take(&user, ID).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *UserRepository) Create(user *ormdomain.User) error {
	return u.db.Create(user).Error
}

func (u *UserRepository) Update(ID uint32, user *ormdomain.User) error {
	return u.db.Model(&ormdomain.User{}).Where("id = ?", ID).Updates(user).Error
}

func (u *UserRepository) Delete(ID uint32) error {
	return u.db.Delete(&ormdomain.User{}, ID).Error
}

func (u *UserRepository) GetGroups(ID uint32) ([]ormdomain.Group, error) {
	var user ormdomain.User
	if err := u.db.Preload("Groups").Take(&user, ID).Error; err != nil {
		return nil, err
	}

	return user.Groups, nil
}
