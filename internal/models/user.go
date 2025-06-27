package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Name           string    `gorm:"type:varchar(255);not null" json:"name"`
	Email          string    `gorm:"type:varchar(255);unique;not null" json:"email"`
	Password       string    `gorm:"type:varchar(255);not null" json:"-"`
	WhatsappNumber string    `gorm:"type:varchar(20)" json:"whatsapp_number"`
	Position       string    `gorm:"type:varchar(100)" json:"position"`
	RoleID         uuid.UUID `gorm:"type:uuid;not null" json:"role_id"`
	SchoolID       uuid.UUID `gorm:"type:uuid;null" json:"school_id"`
	Role           Role      `gorm:"foreignKey:RoleID" json:"role"`
	CreatedAt      time.Time `json:"created_at"`
	CreatedBy      uuid.UUID `gorm:"type:uuid" json:"created_by"`
	UpdatedAt      time.Time `json:"updated_at"`
	UpdatedBy      uuid.UUID `gorm:"type:uuid" json:"updated_by"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	u.CreatedAt = time.Now()
	return
}

func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
	u.UpdatedAt = time.Now()
	return
}
