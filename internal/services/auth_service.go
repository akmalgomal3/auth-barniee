package services

import (
	"errors"
	"fmt"

	"auth-barniee/internal/config"
	"auth-barniee/internal/models"
	"auth-barniee/internal/repositories"
	"auth-barniee/internal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthService interface {
	Login(email, password string) (string, error)
	RegisterUser(name, email, password, roleName string, createdBy uuid.UUID) (*models.User, error)
	GetUserProfile(userID uuid.UUID) (*models.User, *models.School, error) // Returns user and its school
}

type authService struct {
	userRepo   repositories.UserRepository
	roleRepo   repositories.RoleRepository
	schoolRepo repositories.SchoolRepository // New: to fetch school details
	config     *config.Config
}

func NewAuthService(userRepo repositories.UserRepository, roleRepo repositories.RoleRepository, schoolRepo repositories.SchoolRepository, cfg *config.Config) AuthService {
	return &authService{
		userRepo:   userRepo,
		roleRepo:   roleRepo,
		schoolRepo: schoolRepo,
		config:     cfg,
	}
}

func (s *authService) Login(email, password string) (string, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("invalid credentials: user not found")
		}
		return "", fmt.Errorf("failed to find user: %w", err)
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		return "", errors.New("invalid credentials: incorrect password")
	}

	token, err := utils.GenerateToken(user, s.config)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}
	return token, nil
}

func (s *authService) RegisterUser(name, email, password, roleName string, createdBy uuid.UUID) (*models.User, error) {
	existingUser, err := s.userRepo.FindByEmail(email)
	if err == nil && existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}

	role, err := s.roleRepo.FindByName(roleName)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("role '%s' not found", roleName)
		}
		return nil, fmt.Errorf("failed to find role: %w", err)
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &models.User{
		Name:      name,
		Email:     email,
		Password:  hashedPassword,
		RoleID:    role.ID,
		CreatedBy: createdBy,
	}

	err = s.userRepo.Create(user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return user, nil
}

func (s *authService) GetUserProfile(userID uuid.UUID) (*models.User, *models.School, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, errors.New("user not found")
		}
		return nil, nil, fmt.Errorf("failed to retrieve user profile: %w", err)
	}

	var school *models.School
	if user.SchoolID != uuid.Nil {
		school, err = s.schoolRepo.FindByID(user.SchoolID)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			// Log this error, but don't fail the user profile retrieval if school isn't found
			// A user might exist without a linked school (e.g., master admin)
			fmt.Printf("Warning: Could not retrieve school for user %s: %v\n", user.ID.String(), err)
		}
	}

	return user, school, nil
}
