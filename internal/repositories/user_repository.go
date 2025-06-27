package repositories

import (
	"auth-barniee/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *models.User) error
	FindByEmail(email string) (*models.User, error)
	FindByID(id uuid.UUID) (*models.User, error)
	FindAll(roleID *uuid.UUID, schoolID *uuid.UUID) ([]models.User, error) // Added schoolID
	Update(user *models.User) error
	Delete(id uuid.UUID) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	// Removed Preload("School") to prevent circular dependency
	result := r.db.Preload("Role").Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *userRepository) FindByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	// Removed Preload("School") to prevent circular dependency
	result := r.db.Preload("Role").First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *userRepository) FindAll(roleID *uuid.UUID, schoolID *uuid.UUID) ([]models.User, error) {
	var users []models.User
	// Removed Preload("School")
	query := r.db.Preload("Role")
	if roleID != nil && *roleID != uuid.Nil {
		query = query.Where("role_id = ?", *roleID)
	}
	if schoolID != nil && *schoolID != uuid.Nil {
		query = query.Where("school_id = ?", *schoolID)
	}
	result := query.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

func (r *userRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.User{}, id).Error
}
