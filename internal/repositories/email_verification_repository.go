package repositories

import (
	"auth-barniee/internal/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EmailVerificationRepository interface {
	Create(verification *models.EmailVerification) error
	FindByUserIDAndOTP(userID uuid.UUID, otp string) (*models.EmailVerification, error)
	Update(verification *models.EmailVerification) error
	DeleteExpired() error
	FindByUserID(userID uuid.UUID) (*models.EmailVerification, error)
}

type emailVerificationRepository struct {
	db *gorm.DB
}

func NewEmailVerificationRepository(db *gorm.DB) EmailVerificationRepository {
	return &emailVerificationRepository{db: db}
}

func (r *emailVerificationRepository) Create(verification *models.EmailVerification) error {
	return r.db.Create(verification).Error
}

func (r *emailVerificationRepository) FindByUserIDAndOTP(userID uuid.UUID, otp string) (*models.EmailVerification, error) {
	var verification models.EmailVerification
	result := r.db.Where("user_id = ? AND otp = ?", userID, otp).First(&verification)
	if result.Error != nil {
		return nil, result.Error
	}
	return &verification, nil
}

func (r *emailVerificationRepository) Update(verification *models.EmailVerification) error {
	return r.db.Save(verification).Error
}

func (r *emailVerificationRepository) DeleteExpired() error {
	return r.db.Where("expires_at < ?", time.Now()).Delete(&models.EmailVerification{}).Error
}

func (r *emailVerificationRepository) FindByUserID(userID uuid.UUID) (*models.EmailVerification, error) {
	var verification models.EmailVerification
	result := r.db.Where("user_id = ?", userID).Order("created_at DESC").First(&verification)
	if result.Error != nil {
		return nil, result.Error
	}
	return &verification, nil
}
