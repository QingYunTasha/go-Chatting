package domain

type Group struct {
	ID            uint32 `gorm:"primaryKey;autoIncrement"`
	Name          string `gorm:"not null;default:null;unique"`
	Users         []User `gorm:"many2many:group_users;constraint:OnDelete:CASCADE;"`
	OwnerID       uint32
	GroupMessages []GroupMessage
}

type GroupRepository interface {
	Get(ID uint32) (*Group, error)
	GetByName(name string) (*Group, error)
	Create(group *Group) error
	Update(ID uint32, group *Group) error
	Delete(ID uint32) error
	GetUsers(ID uint32) ([]User, error)
	GetMessages(ID uint32) ([]GroupMessage, error)
}
