package repositories

import (
	"core-backend/internal/database"
	"core-backend/internal/models"
	"core-backend/pkg/logger"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type MessageRepository interface {
	SaveMessage(message *models.EncryptedMessages) error
	GetUndeliveredMessages(receiverID uuid.UUID) ([]models.EncryptedMessages, error)
	MarkAsDelivered(messageIDs []uuid.UUID) error
}

type messageRepository struct {
	db *gorm.DB
}

func NewMessageRepository() MessageRepository {
	return &messageRepository{db: database.DB}
}

func (r *messageRepository) SaveMessage(message *models.EncryptedMessages) error {
	if err := r.db.Create(message).Error; err != nil {
		logger.Log.Error("failed to save encrypted message", zap.Error(err))
		return err
	}
	return nil
}

func (r *messageRepository) GetUndeliveredMessages(receiverID uuid.UUID) ([]models.EncryptedMessages, error) {
	var messages []models.EncryptedMessages

	err := r.db.Preload("Status").
		Joins("JOIN delivery_statues ON delivery_statues.message_id = encrypted_messages.message_id").
		Where("encrypted_messages.receiver_id = ? AND delivery_statues.is_delivered = ?", receiverID, false).
		Find(&messages).Error

	if err != nil {
		logger.Log.Error("failed to fetch undelivered messages", zap.Error(err))
		return nil, err
	}
	return messages, nil
}

func (r *messageRepository) MarkAsDelivered(messageIDs []uuid.UUID) error {
	err := r.db.Model(&models.DeliveryStatus{}).
		Where("message_id IN ?", messageIDs).
		Updates(map[string]interface{}{
			"is_delivered": true,
			"delivered_at": gorm.Expr("NOW()"),
		}).Error

	if err != nil {
		logger.Log.Error("failed to update message delivery status", zap.Error(err))
		return err
	}
	return nil
}
