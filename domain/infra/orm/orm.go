package domain

import "time"

type PrivateMessage struct {
	ID         uint32    `gorm:"primaryKey;autoIncrement"`
	SenderID   uint32    `gorm:"not null;default:null;index:sender_receiver"`
	ReceiverID uint32    `gorm:"not null;default:null;index:sender_receiver"`
	Timestamp  time.Time `gorm:"not null;default:null"`
	Content    string    `gorm:"not null;default:null"`
	Sender     User      `gorm:"foreignKey:SenderID;not null;default:null"`
	Receiver   User      `gorm:"foreignKey:ReceiverID;not null;default:null"`
}

type Group struct {
	ID            uint32 `gorm:"primaryKey;autoIncrement"`
	Name          string `gorm:"not null;default:null;unique"`
	Users         []User `gorm:"many2many:group_users;constraint:OnDelete:CASCADE;"`
	GroupMessages []GroupMessage
}

type User struct {
	ID     uint32  `gorm:"primaryKey;autoIncrement"`
	Name   string  `gorm:"not null;default:null;unique"`
	Groups []Group `gorm:"many2many:group_users;constraint:OnDelete:CASCADE;"`
}

type GroupMessage struct {
	ID        uint32    `gorm:"primaryKey;autoIncrement"`
	SenderID  uint32    `gorm:"not null;default:null"`
	GroupID   uint32    `gorm:"not null;default:null"`
	Timestamp time.Time `gorm:"not null;default:null"`
	Content   string    `gorm:"not null;default:null"`
}

type MessageQueue struct {
	Queue chan GroupMessage
}

func NewMessageQueue(size int) *MessageQueue {
	return &MessageQueue{
		Queue: make(chan GroupMessage, size),
	}
}
