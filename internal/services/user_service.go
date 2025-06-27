package services

import (
	"errors"
	"fmt"

	"auth-barniee/internal/models"
	"auth-barniee/internal/repositories"
	"auth-barniee/internal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserService interface {
	CreateTeacherOrStudent(name, email, password, roleName string, adminID uuid.UUID) (*models.User, error)
	GetAllUsers(roleName string, adminUserID uuid.UUID) ([]models.User, error) // Added adminUserID
	GetUserByID(userID uuid.UUID) (*models.User, error)
	UpdateUser(userID, adminID uuid.UUID, name, email *string, roleName *string) (*models.User, error)
	DeleteUser(userID, adminID uuid.UUID) error
}

type userService struct {
	userRepo repositories.UserRepository
	roleRepo repositories.RoleRepository
}

func NewUserService(userRepo repositories.UserRepository, roleRepo repositories.RoleRepository) UserService {
	return &userService{
		userRepo: userRepo,
		roleRepo: roleRepo,
	}
}

func (s *userService) CreateTeacherOrStudent(name, email, password, roleName string, adminID uuid.UUID) (*models.User, error) {
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

	if role.Name != "teacher" && role.Name != "student" {
		return nil, errors.New("can only create users with 'teacher' or 'student' roles")
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	adminUser, err := s.userRepo.FindByID(adminID)
	if err != nil {
		return nil, fmt.Errorf("admin user not found: %w", err)
	}
	if adminUser.Role.Name != "admin" {
		return nil, errors.New("only administrators can create new teachers or students")
	}

	user := &models.User{
		Name:      name,
		Email:     email,
		Password:  hashedPassword,
		RoleID:    role.ID,
		SchoolID:  adminUser.SchoolID, // Assign to the same school as the admin who created it
		CreatedBy: adminID,
	}

	err = s.userRepo.Create(user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return user, nil
}

// GetAllUsers now filters by the admin's school ID, unless it's a master admin.
func (s *userService) GetAllUsers(roleName string, adminUserID uuid.UUID) ([]models.User, error) {
	adminUser, err := s.userRepo.FindByID(adminUserID)
	if err != nil {
		return nil, fmt.Errorf("admin user not found: %w", err)
	}

	var targetRoleID *uuid.UUID
	if roleName != "" {
		role, err := s.roleRepo.FindByName(roleName)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, fmt.Errorf("role '%s' not found", roleName)
			}
			return nil, fmt.Errorf("failed to find role: %w", err)
		}
		targetRoleID = &role.ID
	}

	var targetSchoolID *uuid.UUID
	// If it's a school admin (not master admin), filter by their school_id
	if adminUser.SchoolID != uuid.Nil {
		targetSchoolID = &adminUser.SchoolID
	}
	// If it's a master admin (SchoolID is nil), they can see all users across all schools
	// This logic is implicitly handled if targetSchoolID remains nil.

	users, err := s.userRepo.FindAll(targetRoleID, targetSchoolID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve users: %w", err)
	}
	return users, nil
}

func (s *userService) GetUserByID(userID uuid.UUID) (*models.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to retrieve user: %w", err)
	}
	return user, nil
}

func (s *userService) UpdateUser(userID, adminID uuid.UUID, name, email *string, roleName *string) (*models.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	adminUser, err := s.userRepo.FindByID(adminID)
	if err != nil {
		return nil, fmt.Errorf("admin user not found: %w", err)
	}

	// Authorization check: master admin can update any user; school admin can only update users within their school.
	if adminUser.Role.Name == "admin" && adminUser.SchoolID != uuid.Nil && user.SchoolID != adminUser.SchoolID {
		return nil, errors.New("unauthorized: school admin cannot update users outside their school")
	}

	if name != nil {
		user.Name = *name
	}
	if email != nil {
		if *email != user.Email {
			existingUser, err := s.userRepo.FindByEmail(*email)
			if err == nil && existingUser != nil && existingUser.ID != user.ID {
				return nil, errors.New("email already taken by another user")
			}
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, fmt.Errorf("failed to check existing email: %w", err)
			}
		}
		user.Email = *email
	}
	if roleName != nil {
		role, err := s.roleRepo.FindByName(*roleName)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, fmt.Errorf("role '%s' not found", *roleName)
			}
			return nil, fmt.Errorf("failed to find role: %w", err)
		}
		user.RoleID = role.ID
	}
	user.UpdatedBy = adminID

	err = s.userRepo.Update(user)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}
	return user, nil
}

func (s *userService) DeleteUser(userID, adminID uuid.UUID) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return fmt.Errorf("failed to find user: %w", err)
	}

	adminUser, err := s.userRepo.FindByID(adminID)
	if err != nil {
		return fmt.Errorf("admin user not found: %w", err)
	}

	if adminUser.Role.Name == "admin" && adminUser.SchoolID != uuid.Nil && user.SchoolID != adminUser.SchoolID {
		return errors.New("unauthorized: school admin cannot delete users outside their school")
	}

	if user.ID == adminID {
		return errors.New("cannot delete your own admin account")
	}

	return s.userRepo.Delete(user.ID)
}
