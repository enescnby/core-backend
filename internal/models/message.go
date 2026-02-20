package models

import (
	"time"

	"github.com/google/uuid"
)

type EncryptedMessages struct {
	MessageID  uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	SenderID   uuid.UUID `gorm:"type:uuid;index;not null"`
	ReceiverID uuid.UUID `gorm:"type:uuid;index;not null"`
	Ciphertext []byte    `gorm:"type:bytea;not null"`
	Nonce      string    `gorm:"type:varchar;not null"`
	AuthTag    string    `gorm:"type:varchar;not null"`
	KeyVersion int       `gorm:"type:int"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`

	Status []DeliveryStatus `gorm:"foreignKey:MessageID"`
}

type DeliveryStatus struct {
	DeliveryID  int       `gorm:"primaryKey;autoIncrement"`
	MessageID   uuid.UUID `gorm:"type:uuid;index;not null"`
	IsDelivered bool      `gorm:"default:false"`
	DeliveredAt *time.Time
	Details     string `gorm:"type:text"`
}

type SecurityAuditLog struct {
	AuditID    int       `gorm:"primaryKey;autoIncrement"`
	UserID     uuid.UUID `gorm:"type:uuid;index;not null"`
	ActionType string    `gorm:"type:varchar;not null"`
	IPAddress  string    `gorm:"type:varchar"`
	Timestamp  time.Time `gorm:"autoCreateTime"`
}
