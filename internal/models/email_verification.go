package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EmailVerification struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserID     uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	Email      string    `gorm:"type:varchar(255);not null" json:"email"`
	OTP        string    `gorm:"type:varchar(6);not null" json:"otp"`
	ExpiresAt  time.Time `gorm:"not null" json:"expires_at"`
	IsVerified bool      `gorm:"default:false" json:"is_verified"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	User       User      `gorm:"foreignKey:UserID"`
}

func (ev *EmailVerification) BeforeCreate(tx *gorm.DB) (err error) {
	if ev.ID == uuid.Nil {
		ev.ID = uuid.New()
	}
	ev.CreatedAt = time.Now()
	return
}

func (ev *EmailVerification) BeforeUpdate(tx *gorm.DB) (err error) {
	ev.UpdatedAt = time.Now()
	return
}
