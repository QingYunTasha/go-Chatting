package domain

import "database/sql"

type UserStatus string

const (
	Online  UserStatus = "online"
	Offline UserStatus = "offline"
	//Typeing UserStatus = "typing"
)

type User struct {
	ID                 uint32         `gorm:"primaryKey;autoIncrement"`
	Name               string         `gorm:"not null;default:null;unique"`
	Email              string         `gorm:"not null;default:null;unique"`
	Password           string         `gorm:"not null;default:null"`
	Status             UserStatus     `gorm:"not null;default:offline"`
	Friends            []User         `gorm:"many2many:user_friends"`
	Groups             []Group        `gorm:"many2many:group_users;constraint:OnDelete:CASCADE;"`
	PasswordResetToken sql.NullString `gorm:"default:null"`
}

type UserRepository interface {
	Get(ID uint32) (*User, error)
	GetByEmail(email string) (*User, error)
	Create(user *User) error
	Update(ID uint32, user *User) error
	Delete(ID uint32) error
	GetGroups(ID uint32) ([]Group, error)
	GetFriends(ID uint32) ([]User, error)
	SavePasswordResetToken(userID uint32, token string) error
	FindByPasswordResetToken(token string) (*User, error)
	ClearPasswordResetToken(userID uint32) error
	ResetPassword(userID uint32, password string) error
	IsGroupMember(userID uint32, groupID uint32) (bool, error)
	AddGroupMember(userID, groupID uint32) error
	RemoveGroupMember(userID uint32, groupID uint32) error
	AreFriends(userID uint32, friendID uint32) (bool, error)
	AddFriend(userID uint32, friendID uint32) error
	RemoveFriend(userID uint32, friendID uint32) error
}
