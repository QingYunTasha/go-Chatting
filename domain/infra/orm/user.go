package domain

type UserStatus string

const (
	Online  UserStatus = "online"
	Offline UserStatus = "offline"
	Typeing UserStatus = "typing"
)

type User struct {
	ID     uint32     `gorm:"primaryKey;autoIncrement"`
	Name   string     `gorm:"not null;default:null;unique"`
	Status UserStatus `gorm:"not null;default:offline"`
	Groups []Group    `gorm:"many2many:group_users;constraint:OnDelete:CASCADE;"`
}

type UserRepository interface {
	Get(ID uint32) (*User, error)
	Create(user *User) error
	Update(ID uint32, user *User) error
	Delete(ID uint32) error
	GetGroups(ID uint32) ([]Group, error)
}
