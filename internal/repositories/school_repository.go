package repositories

import (
	"auth-barniee/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SchoolRepository interface {
	Create(school *models.School) error
	FindByID(id uuid.UUID) (*models.School, error)
	Update(school *models.School) error
	FindByAdminUserID(adminUserID uuid.UUID) (*models.School, error)
}

type schoolRepository struct {
	db *gorm.DB
}

func NewSchoolRepository(db *gorm.DB) SchoolRepository {
	return &schoolRepository{db: db}
}

func (r *schoolRepository) Create(school *models.School) error {
	return r.db.Create(school).Error
}

func (r *schoolRepository) FindByID(id uuid.UUID) (*models.School, error) {
	var school models.School
	result := r.db.Preload("Package").First(&school, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &school, nil
}

func (r *schoolRepository) Update(school *models.School) error {
	return r.db.Save(school).Error
}

func (r *schoolRepository) FindByAdminUserID(adminUserID uuid.UUID) (*models.School, error) {
	var school models.School
	result := r.db.Preload("Package").Where("admin_user_id = ?", adminUserID).First(&school)
	if result.Error != nil {
		return nil, result.Error
	}
	return &school, nil
}
