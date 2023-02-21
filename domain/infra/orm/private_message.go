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

type PrivateMessageRepository interface {
	Get(ID uint32) (*PrivateMessage, error)
	GetBetweenUsersAfterMsgID(senderID, receiverID, lastID uint32) ([]PrivateMessage, error)
	Create(message *PrivateMessage) error
	Delete(ID uint32) error
}
