package domain

import "time"

type PrivateMessage struct {
	ID        uint `gorm:"primaryKey"`
	Sender    string
	Receiver  string
	Timestamp time.Time
	Message   string
}

type GroupMessage struct {
	ID        uint `gorm:"primaryKey"`
	Sender    string
	GroupID   uint
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

type Message struct {
	ID        uint32 `gorm:"primaryKey"`
	SenderID  uint32
	Group     string
	Timestamp time.Time
	Content   string
}

type MessageQueue struct {
	Queue chan Message
}

func NewMessageQueue(size int) *MessageQueue {
	return &MessageQueue{
		Queue: make(chan Message, size),
	}
}

func (mq *MessageQueue) Push(item Message) {
	mq.Queue <- item
}

func (mq *MessageQueue) Pop() Message {
	return <-mq.Queue
}

func (mq *MessageQueue) Len() int {
	return len(mq.Queue)
}
