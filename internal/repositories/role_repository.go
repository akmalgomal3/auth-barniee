package repositories

import (
	"auth-barniee/internal/models"
	"gorm.io/gorm"
)

type RoleRepository interface {
	FindByName(name string) (*models.Role, error)
}

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) FindByName(name string) (*models.Role, error) {
	var role models.Role
	result := r.db.Where("name = ?", name).First(&role)
	if result.Error != nil {
		return nil, result.Error
	}
	return &role, nil
}
