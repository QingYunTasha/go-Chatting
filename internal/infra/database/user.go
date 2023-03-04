package orm

import (
	"database/sql"
	"errors"
	"fmt"
	dbdomain "go-Chatting/domain/infra/database"

	"golang.org/x/crypto/bcrypt"
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

func (u *UserRepository) GetByEmail(email string) (*dbdomain.User, error) {
	var user dbdomain.User
	if err := u.db.Where("email = ?", email).Take(&user).Error; err != nil {
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

func (u *UserRepository) SavePasswordResetToken(userID uint32, token string) error {
	user := &dbdomain.User{}

	// Find the user by ID
	if err := u.db.Where("id = ?", userID).First(user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// User not found
			return fmt.Errorf("user with ID %d not found", userID)
		}
		// Other errors
		return err
	}

	// Update the user's password reset token
	user.PasswordResetToken = sql.NullString{
		String: token,
		Valid:  true,
	}

	if err := u.db.Save(user).Error; err != nil {
		// Error while updating the user
		return err
	}

	return nil
}

func (r *UserRepository) FindByPasswordResetToken(token string) (*dbdomain.User, error) {
	var user dbdomain.User
	result := r.db.Where("password_reset_token = ?", token).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *UserRepository) ClearPasswordResetToken(userID uint32) error {
	result := r.db.Model(&dbdomain.User{}).Where("id = ?", userID).Update("password_reset_token", nil)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *UserRepository) ResetPassword(userID uint32, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	result := r.db.Model(&dbdomain.User{}).Where("id = ?", userID).Update("password", string(hashedPassword))
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *UserRepository) IsGroupMember(userID uint32, groupName string) (bool, error) {
	groups, err := r.GetGroups(userID)
	if err != nil {
		return false, err
	}

	for _, group := range groups {
		if groupName == group.Name {
			return true, nil
		}
	}
	return false, nil
}

func (r *UserRepository) AddGroupMember(userID uint32, groupName string) error {
	user := dbdomain.User{}
	group := dbdomain.Group{
		Name: groupName,
	}
	err := r.db.Take(&user, userID).Error
	if err != nil {
		return err
	}
	err = r.db.Take(&group).Error
	if err != nil {
		return err
	}
	err = r.db.Model(&user).Association("Groups").Append(&group)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) RemoveGroupMember(userID uint32, groupName string) error {
	user := &dbdomain.User{}
	group := &dbdomain.Group{}

	// Find the user and group
	err := r.db.Preload("Groups", "name = ?", groupName).First(user, userID).Error
	if err != nil {
		return err
	}
	if len(user.Groups) == 0 {
		return errors.New("group not found")
	}
	group = &user.Groups[0]

	// Remove the user from the group
	err = r.db.Model(group).Association("Members").Delete(&dbdomain.User{ID: userID})
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) AreFriends(userID uint32, friendEmail string) (bool, error) {
	user, err := r.Get(userID)
	if err != nil {
		return false, err
	}

	// check if friend exists in user's friend list
	for _, friend := range user.Friends {
		if friend.Email == friendEmail {
			return true, nil
		}
	}

	return false, nil
}

func (r *UserRepository) AddFriend(userID uint32, friendEmail string) error {
	user, err := r.Get(userID)
	if err != nil {
		return err
	}

	friend, err := r.GetByEmail(friendEmail)
	if err != nil {
		return err
	}

	// add friend to user's friend list
	user.Friends = append(user.Friends, *friend)

	// save changes
	return r.db.Save(user).Error
}

func (r *UserRepository) RemoveFriend(userID uint32, friendEmail string) error {
	user, err := r.Get(userID)
	if err != nil {
		return err
	}

	// remove friend from user's friend list
	for i, friend := range user.Friends {
		if friend.Email == friendEmail {
			user.Friends = append(user.Friends[:i], user.Friends[i+1:]...)
			break
		}
	}

	// save changes
	return r.db.Save(user).Error
}
