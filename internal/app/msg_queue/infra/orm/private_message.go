package orm

import (
	"errors"

	"gorm.io/gorm"
)

type PrivateMessageRepository struct {
	db *gorm.DB
}

func NewPrivateMessageRepository(db *gorm.DB) PrivateMessageRepository {
	return PrivateMessageRepository{
		db: db,
	}
}

func (r PrivateMessageRepository) Create(message *PrivateMessage) error {
	return r.db.Create(message).Error
}

func (r PrivateMessageRepository) Get(ID uint32) (*PrivateMessage, error) {
	var message PrivateMessage
	if err := r.db.Preload("Sender").Preload("Receiver").First(&message, ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &message, nil
}

func (r PrivateMessageRepository) GetBetweenUsers(senderID, receiverID uint32) ([]PrivateMessage, error) {
	var messages []PrivateMessage
	if err := r.db.Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)",
		senderID, receiverID, receiverID, senderID).
		Order("timestamp").Find(&messages).Error; err != nil {
		return nil, err
	}

	return messages, nil
}

func (r PrivateMessageRepository) Delete(ID uint32) error {
	return r.db.Delete(&PrivateMessage{}, ID).Error
}
