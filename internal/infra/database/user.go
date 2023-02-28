package orm

import (
	dbdomain "go-Chatting/domain/infra/database"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) dbdomain.UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (u *UserRepository) Get(ID uint32) (*dbdomain.User, error) {
	var user dbdomain.User
	err := u.db.Take(&user, ID).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *UserRepository) Create(user *dbdomain.User) error {
	return u.db.Create(user).Error
}

func (u *UserRepository) Update(ID uint32, user *dbdomain.User) error {
	return u.db.Model(&dbdomain.User{}).Where("id = ?", ID).Updates(user).Error
}

func (u *UserRepository) Delete(ID uint32) error {
	return u.db.Delete(&dbdomain.User{}, ID).Error
}

func (u *UserRepository) GetGroups(ID uint32) ([]dbdomain.Group, error) {
	var user dbdomain.User
	if err := u.db.Preload("Groups").Take(&user, ID).Error; err != nil {
		return nil, err
	}

	return user.Groups, nil
}

func (u *UserRepository) GetFriends(ID uint32) ([]dbdomain.User, error) {
	var user dbdomain.User
	if err := u.db.Preload("Friends").Take(&user, ID).Error; err != nil {
		return nil, err
	}

	return user.Friends, nil
}
