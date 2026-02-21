package repositories

import (
	"core-backend/internal/database"
	"core-backend/internal/models"
	"core-backend/pkg/logger"
	"errors"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *models.User) error
	GetUserByCoreGuardID(coreGuardID string) (*models.User, error)
	UpdateDevice(userID uuid.UUID, newDevice *models.UserDevice) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository() UserRepository {
	return &userRepository{db: database.DB}
}

func (r *userRepository) CreateUser(user *models.User) error {
	if err := r.db.Create(user).Error; err != nil {
		logger.Log.Error("Failed to create user", zap.Error(err))
		return err
	}
	return nil
}

func (r *userRepository) GetUserByCoreGuardID(coreGuardID string) (*models.User, error) {
	var user models.User

	err := r.db.Preload("Key").Preload("Device").Where("core_guard_id = ?", coreGuardID).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Log.Warn("user not found", zap.String("coreGuardID", coreGuardID))
			return nil, err
		}
		logger.Log.Error("database error while fetching user", zap.Error(err))
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) UpdateDevice(userID uuid.UUID, newDevice *models.UserDevice) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ?", userID).Delete(&models.UserDevice{}).Error; err != nil {
			return err
		}
		return tx.Create(newDevice).Error
	})

	if err != nil {
		logger.Log.Error("failed to update user device", zap.Error(err))
		return err
	}
	return nil
}
