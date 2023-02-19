package domain

import "time"

type PrivateMessage struct {
	ID        uint `gorm:"primaryKey"`
	Sender    string
	Receiver  string
	Timestamp time.Time
	Message   string
}

type Group struct {
	ID    uint   `gorm:"primaryKey"`
	Name  string `gorm:"unique"`
	Users []User `gorm:"many2many:group_users;"`
}

type User struct {
	ID     uint    `gorm:"primaryKey"`
	Name   string  `gorm:"unique"`
	Groups []Group `gorm:"many2many:group_users;"`
}

type GroupMessage struct {
	ID        uint32 `gorm:"primaryKey"`
	SenderID  uint32
	Group     string
	Timestamp time.Time
	Content   string
}

type MessageQueue struct {
	Queue chan GroupMessage
}

func NewMessageQueue(size int) *MessageQueue {
	return &MessageQueue{
		Queue: make(chan GroupMessage, size),
	}
}
