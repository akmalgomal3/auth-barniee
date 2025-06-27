package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type School struct {
	ID                    uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	Name                  string     `gorm:"type:varchar(255);not null" json:"name"`
	EducationLevel        string     `gorm:"type:varchar(50);not null" json:"education_level"`
	Status                string     `gorm:"type:varchar(50);not null" json:"status"`
	Address               string     `gorm:"type:text;not null" json:"address"`
	InitialStudentCount   int        `gorm:"not null" json:"initial_student_count"`
	AdminUserID           uuid.UUID  `gorm:"type:uuid;null" json:"admin_user_id"`
	PackageID             uuid.UUID  `gorm:"type:uuid;not null" json:"package_id"`
	SubscriptionStartDate *time.Time `json:"subscription_start_date"`
	SubscriptionEndDate   *time.Time `json:"subscription_end_date"`
	MaxStudentsAllowed    int        `gorm:"not null" json:"max_students_allowed"`
	CreatedAt             time.Time  `json:"created_at"`
	CreatedBy             uuid.UUID  `gorm:"type:uuid" json:"created_by"`
	UpdatedAt             time.Time  `json:"updated_at"`
	UpdatedBy             uuid.UUID  `gorm:"type:uuid" json:"updated_by"`
	Package               Package    `gorm:"foreignKey:PackageID"`
}

func (s *School) BeforeCreate(tx *gorm.DB) (err error) {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	s.CreatedAt = time.Now()
	return
}

func (s *School) BeforeUpdate(tx *gorm.DB) (err error) {
	s.UpdatedAt = time.Now()
	return
}
