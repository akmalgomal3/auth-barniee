package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Package struct {
	ID              uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Name            string    `gorm:"type:varchar(100);not null;unique" json:"name"`
	PricePerStudent *float64  `gorm:"type:decimal(10,2)" json:"price_per_student,omitempty"`
	PricePerYear    *float64  `gorm:"type:decimal(10,2)" json:"price_per_year,omitempty"`
	DurationDays    *int      `json:"duration_days,omitempty"`
	MaxStudents     *int      `json:"max_students,omitempty"`
	Features        string    `gorm:"type:jsonb" json:"features"`
	CreatedAt       time.Time `json:"created_at"`
	CreatedBy       uuid.UUID `gorm:"type:uuid" json:"created_by"`
	UpdatedAt       time.Time `json:"updated_at"`
	UpdatedBy       uuid.UUID `gorm:"type:uuid" json:"updated_by"`
}

func (p *Package) BeforeCreate(tx *gorm.DB) (err error) {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	p.CreatedAt = time.Now()
	return
}

func (p *Package) BeforeUpdate(tx *gorm.DB) (err error) {
	p.UpdatedAt = time.Now()
	return
}
