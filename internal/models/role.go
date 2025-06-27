package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Role struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Name        string    `gorm:"type:varchar(50);unique;not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	CreatedBy   uuid.UUID `gorm:"type:uuid" json:"created_by"`
	UpdatedAt   time.Time `json:"updated_at"`
	UpdatedBy   uuid.UUID `gorm:"type:uuid" json:"updated_by"`
	Users       []User    `gorm:"foreignKey:RoleID"`
}

func (r *Role) BeforeCreate(tx *gorm.DB) (err error) {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	r.CreatedAt = time.Now()
	return
}

func (r *Role) BeforeUpdate(tx *gorm.DB) (err error) {
	r.UpdatedAt = time.Now()
	return
}
