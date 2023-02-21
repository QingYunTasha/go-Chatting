package orm

import (
	ormdomain "github.com/QingYunTasha/go-Chatting/domain/infra/orm"
	"gorm.io/gorm"
)

type PrivateMessageRepository struct {
	db *gorm.DB
}

func NewPrivateMessageRepository(db *gorm.DB) ormdomain.PrivateMessageRepository {
	return &PrivateMessageRepository{
		db: db,
	}
}

func (r PrivateMessageRepository) Create(message *ormdomain.PrivateMessage) error {
	return r.db.Create(message).Error
}

func (r PrivateMessageRepository) Get(ID uint32) (*ormdomain.PrivateMessage, error) {
	var message ormdomain.PrivateMessage
	if err := r.db.Preload("Sender").Preload("Receiver").First(&message, ID).Error; err != nil {
		return nil, err
	}

	return &message, nil
}

func (r PrivateMessageRepository) GetBetweenUsersAfterMsgID(senderID, receiverID, lastID uint32) ([]ormdomain.PrivateMessage, error) {
	var messages []ormdomain.PrivateMessage
	if err := r.db.Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)",
		senderID, receiverID, receiverID, senderID).
		Order("timestamp").Find(&messages).Error; err != nil {
		return nil, err
	}

	return messages, nil
}

func (r PrivateMessageRepository) Delete(ID uint32) error {
	return r.db.Delete(&ormdomain.PrivateMessage{}, ID).Error
}
