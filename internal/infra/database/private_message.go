package orm

import (
	dbdomain "go-Chatting/domain/infra/database"

	"gorm.io/gorm"
)

type PrivateMessageRepository struct {
	db *gorm.DB
}

func NewPrivateMessageRepository(db *gorm.DB) dbdomain.PrivateMessageRepository {
	return &PrivateMessageRepository{
		db: db,
	}
}

func (r PrivateMessageRepository) Create(message *dbdomain.PrivateMessage) error {
	return r.db.Create(message).Error
}

func (r PrivateMessageRepository) Get(ID uint32) (*dbdomain.PrivateMessage, error) {
	var message dbdomain.PrivateMessage
	if err := r.db.Preload("Sender").Preload("Receiver").First(&message, ID).Error; err != nil {
		return nil, err
	}

	return &message, nil
}

func (r PrivateMessageRepository) GetBetweenUsersAfterOrderID(senderID, receiverID, lastOrderID uint32) ([]dbdomain.PrivateMessage, error) {
	var messages []dbdomain.PrivateMessage
	if err := r.db.Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?) AND (order_id > ?)",
		senderID, receiverID, receiverID, senderID, lastOrderID).Find(&messages).Error; err != nil {
		return nil, err
	}

	return messages, nil
}

func (r PrivateMessageRepository) Delete(ID uint32) error {
	return r.db.Delete(&dbdomain.PrivateMessage{}, ID).Error
}

func (r PrivateMessageRepository) GetLastIDBetweenUsers(senderID, receiverID uint32) (uint32, error) {
	var maxOrderID uint32
	err := r.db.Model(&dbdomain.PrivateMessage{}).Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)",
		senderID, receiverID, receiverID, senderID).Select("MAX(order_id)").Row().Scan(&maxOrderID)

	if err != nil {
		return 0, err
	}
	return maxOrderID, nil
}
