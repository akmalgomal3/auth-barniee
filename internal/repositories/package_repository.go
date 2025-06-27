package repositories

import (
	"auth-barniee/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PackageRepository interface {
	FindByID(id uuid.UUID) (*models.Package, error)
	FindByName(name string) (*models.Package, error)
	FindAll() ([]models.Package, error)
}

type packageRepository struct {
	db *gorm.DB
}

func NewPackageRepository(db *gorm.DB) PackageRepository {
	return &packageRepository{db: db}
}

func (r *packageRepository) FindByID(id uuid.UUID) (*models.Package, error) {
	var pkg models.Package
	result := r.db.First(&pkg, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &pkg, nil
}

func (r *packageRepository) FindByName(name string) (*models.Package, error) {
	var pkg models.Package
	result := r.db.Where("name = ?", name).First(&pkg)
	if result.Error != nil {
		return nil, result.Error
	}
	return &pkg, nil
}

func (r *packageRepository) FindAll() ([]models.Package, error) {
	var pkgs []models.Package
	result := r.db.Find(&pkgs)
	if result.Error != nil {
		return nil, result.Error
	}
	return pkgs, nil
}
