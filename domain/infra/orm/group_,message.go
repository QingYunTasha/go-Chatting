package domain

import "time"

type GroupMessage struct {
	ID        uint32    `gorm:"primaryKey;autoIncrement"`
	SenderID  uint32    `gorm:"not null;default:null"`
	GroupID   uint32    `gorm:"not null;default:null"`
	Timestamp time.Time `gorm:"not null;default:null"`
	Content   string    `gorm:"not null;default:null"`
}

type GroupMessageRepository interface {
	Get(ID uint32) (*GroupMessage, error)
	Create(*GroupMessage) error
	Delete(ID uint32) error
	DeleteByGroup(ID uint32) error
}
