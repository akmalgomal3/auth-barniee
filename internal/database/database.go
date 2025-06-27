package database

import (
	"fmt"
	"gorm.io/driver/postgres"
	"log"
	"time"

	"auth-barniee/internal/config"
	"auth-barniee/internal/models"
	"auth-barniee/internal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func InitDB(cfg *config.Config) *gorm.DB {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	err = db.AutoMigrate(
		&models.Role{},
		&models.User{},
		&models.School{},
		&models.Package{},
		&models.EmailVerification{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	seedRoles(db)
	seedPackages(db)
	seedAdminUser(db)

	return db
}

func seedRoles(db *gorm.DB) {
	roles := []models.Role{
		{Name: "admin", Description: "Administrator"},
		{Name: "teacher", Description: "Teacher"},
		{Name: "student", Description: "Student"},
	}

	for _, role := range roles {
		var existingRole models.Role
		db.Where("name = ?", role.Name).First(&existingRole)
		if existingRole.ID == uuid.Nil {
			role.ID = uuid.New()
			db.Create(&role)
		}
	}
}

func seedPackages(db *gorm.DB) {
	freeTrialMaxStudents := 50
	freeTrialDurationDays := 30
	premiumPricePerStudent := 50000.0
	enterprisePricePerYear := 10000000000.0

	packages := []models.Package{
		{
			Name:         "Free Trial",
			DurationDays: &freeTrialDurationDays,
			MaxStudents:  &freeTrialMaxStudents,
			Features:     `["Dashboard dasar", "Laporan bulanan", "Email support", "Data backup"]`,
		},
		{
			Name:            "Premium",
			PricePerStudent: &premiumPricePerStudent,
			Features:        `["Unlimited siswa", "AI Analytics lengkap", "Real-time monitoring", "Priority support 24/7", "Custom reports", "Parent app access"]`,
		},
		{
			Name:         "Enterprise",
			PricePerYear: &enterprisePricePerYear,
			Features:     `["Unlimited siswa", "Multi-campus support", "Custom AI features", "Dedicated account manager", "On-site training", "API integration"]`,
		},
	}

	for _, pkg := range packages {
		var existingPkg models.Package
		db.Where("name = ?", pkg.Name).First(&existingPkg)
		if existingPkg.ID == uuid.Nil {
			pkg.ID = uuid.New()
			db.Create(&pkg)
		}
	}
}

func seedAdminUser(db *gorm.DB) {
	var adminRole models.Role
	db.Where("name = ?", "admin").First(&adminRole)

	if adminRole.ID == uuid.Nil {
		log.Println("Admin role not found, cannot seed master admin user.")
		return
	}

	var existingAdmin models.User
	db.Where("email = ?", "masteradmin@barniee.com").First(&existingAdmin)

	if existingAdmin.ID == uuid.Nil {
		hashedPassword, err := utils.HashPassword("masteradminpassword")
		if err != nil {
			log.Fatalf("Failed to hash master admin password: %v", err)
		}

		masterAdminUser := models.User{
			ID:        uuid.New(),
			Name:      "Barniee Master Admin",
			Email:     "masteradmin@barniee.com",
			Password:  hashedPassword,
			RoleID:    adminRole.ID,
			CreatedAt: time.Now(),
		}
		db.Create(&masterAdminUser)
		log.Println("Default master admin user 'masteradmin@barniee.com' created.")
	}
}
