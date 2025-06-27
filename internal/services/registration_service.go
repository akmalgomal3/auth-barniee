package services

import (
	"errors"
	"fmt"
	"time"

	"auth-barniee/internal/config"
	"auth-barniee/internal/models"
	"auth-barniee/internal/repositories"
	"auth-barniee/internal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RegistrationService interface {
	RegisterSchoolInfo(schoolName, educationLevel, status, address string, initialStudentCount int) (*models.School, error)
	RegisterAdminInfo(schoolID uuid.UUID, adminName, adminEmail, whatsappNumber, position string) (*models.User, string, error)
	SelectPackage(schoolID, packageID uuid.UUID) (*models.School, error)
	RequestEmailVerificationOTP(userID uuid.UUID) error
	VerifyEmailOTP(userID uuid.UUID, otp string) error
	CompleteRegistration(schoolID uuid.UUID) (*models.School, error)
	GetSchoolByID(schoolID uuid.UUID) (*models.School, error)
	GetPackageByID(packageID uuid.UUID) (*models.Package, error)
}

type registrationService struct {
	schoolRepo      repositories.SchoolRepository
	userRepo        repositories.UserRepository
	roleRepo        repositories.RoleRepository
	packageRepo     repositories.PackageRepository
	emailVerifyRepo repositories.EmailVerificationRepository
	config          *config.Config
}

func NewRegistrationService(
	schoolRepo repositories.SchoolRepository,
	userRepo repositories.UserRepository,
	roleRepo repositories.RoleRepository,
	packageRepo repositories.PackageRepository,
	emailVerifyRepo repositories.EmailVerificationRepository,
	cfg *config.Config,
) RegistrationService {
	return &registrationService{
		schoolRepo:      schoolRepo,
		userRepo:        userRepo,
		roleRepo:        roleRepo,
		packageRepo:     packageRepo,
		emailVerifyRepo: emailVerifyRepo,
		config:          cfg,
	}
}

func (s *registrationService) RegisterSchoolInfo(schoolName, educationLevel, status, address string, initialStudentCount int) (*models.School, error) {
	freeTrialPkg, err := s.packageRepo.FindByName("Free Trial")
	if err != nil {
		return nil, fmt.Errorf("Free Trial package not found in system: %w", err)
	}

	school := &models.School{
		Name:                schoolName,
		EducationLevel:      educationLevel,
		Status:              status,
		Address:             address,
		InitialStudentCount: initialStudentCount,
		PackageID:           freeTrialPkg.ID,
		MaxStudentsAllowed:  *freeTrialPkg.MaxStudents,
		CreatedBy:           uuid.Nil,
	}

	err = s.schoolRepo.Create(school)
	if err != nil {
		return nil, fmt.Errorf("failed to register school info: %w", err)
	}
	return school, nil
}

func (s *registrationService) RegisterAdminInfo(schoolID uuid.UUID, adminName, adminEmail, whatsappNumber, position string) (*models.User, string, error) {
	school, err := s.schoolRepo.FindByID(schoolID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", errors.New("school not found")
		}
		return nil, "", fmt.Errorf("failed to find school: %w", err)
	}

	existingUser, err := s.userRepo.FindByEmail(adminEmail)
	if err == nil && existingUser != nil {
		return nil, "", errors.New("admin user with this email already exists")
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, "", fmt.Errorf("failed to check existing user: %w", err)
	}

	adminRole, err := s.roleRepo.FindByName("admin")
	if err != nil {
		return nil, "", fmt.Errorf("admin role not found: %w", err)
	}

	generatedPassword := utils.GenerateRandomPassword(12)
	hashedPassword, err := utils.HashPassword(generatedPassword)
	if err != nil {
		return nil, "", fmt.Errorf("failed to hash generated password: %w", err)
	}

	adminUser := &models.User{
		Name:           adminName,
		Email:          adminEmail,
		Password:       hashedPassword,
		WhatsappNumber: whatsappNumber,
		Position:       position,
		RoleID:         adminRole.ID,
		SchoolID:       schoolID,
		CreatedBy:      uuid.Nil,
	}

	err = s.userRepo.Create(adminUser)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create admin user: %w", err)
	}

	school.AdminUserID = adminUser.ID
	school.UpdatedBy = uuid.Nil
	err = s.schoolRepo.Update(school)
	if err != nil {
		return nil, "", fmt.Errorf("failed to update school with admin user ID: %w", err)
	}

	return adminUser, generatedPassword, nil
}

func (s *registrationService) SelectPackage(schoolID, packageID uuid.UUID) (*models.School, error) {
	school, err := s.schoolRepo.FindByID(schoolID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("school not found")
		}
		return nil, fmt.Errorf("failed to find school: %w", err)
	}

	pkg, err := s.packageRepo.FindByID(packageID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("package not found")
		}
		return nil, fmt.Errorf("failed to find package: %w", err)
	}

	school.PackageID = pkg.ID
	if pkg.MaxStudents != nil {
		school.MaxStudentsAllowed = *pkg.MaxStudents
	} else {
		school.MaxStudentsAllowed = 0
	}

	if pkg.Name == "Free Trial" && pkg.DurationDays != nil {
		now := time.Now()
		school.SubscriptionStartDate = &now
		expiry := now.Add(time.Duration(*pkg.DurationDays) * 24 * time.Hour)
		school.SubscriptionEndDate = &expiry
	} else {
		school.SubscriptionStartDate = nil
		school.SubscriptionEndDate = nil
	}

	school.UpdatedBy = uuid.Nil
	err = s.schoolRepo.Update(school)
	if err != nil {
		return nil, fmt.Errorf("failed to select package for school: %w", err)
	}
	return school, nil
}

func (s *registrationService) RequestEmailVerificationOTP(userID uuid.UUID) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found for verification")
		}
		return fmt.Errorf("failed to find user for OTP request: %w", err)
	}

	existingVerification, err := s.emailVerifyRepo.FindByUserID(userID)
	if err == nil && existingVerification != nil && !existingVerification.IsVerified {
		existingVerification.IsVerified = false // Mark as invalidated if not yet verified
		s.emailVerifyRepo.Update(existingVerification)
	}

	otp := utils.GenerateOTP()
	expiresAt := time.Now().Add(time.Duration(s.config.OTPExpiryMinutes) * time.Minute)

	verification := &models.EmailVerification{
		UserID:    userID,
		Email:     user.Email,
		OTP:       otp,
		ExpiresAt: expiresAt,
	}

	err = s.emailVerifyRepo.Create(verification)
	if err != nil {
		return fmt.Errorf("failed to save OTP: %w", err)
	}

	subject := "Barniee: Kode Verifikasi Email Anda"
	body := fmt.Sprintf("Halo %s,\n\nKode verifikasi Anda adalah: %s\nKode ini akan kedaluwarsa dalam %d menit.\n\nTerima kasih,\nTim Barniee", user.Name, otp, s.config.OTPExpiryMinutes)

	err = utils.SendEmail(s.config, user.Email, subject, body)
	if err != nil {
		return fmt.Errorf("failed to send OTP email: %w", err)
	}
	return nil
}

func (s *registrationService) VerifyEmailOTP(userID uuid.UUID, otp string) error {
	verification, err := s.emailVerifyRepo.FindByUserIDAndOTP(userID, otp)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("invalid OTP or user ID")
		}
		return fmt.Errorf("failed to find OTP: %w", err)
	}

	if verification.IsVerified {
		return errors.New("email already verified")
	}

	if time.Now().After(verification.ExpiresAt) {
		return errors.New("OTP has expired")
	}

	verification.IsVerified = true
	err = s.emailVerifyRepo.Update(verification)
	if err != nil {
		return fmt.Errorf("failed to mark OTP as verified: %w", err)
	}

	return nil
}

func (s *registrationService) CompleteRegistration(schoolID uuid.UUID) (*models.School, error) {
	school, err := s.schoolRepo.FindByID(schoolID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("school not found for completion")
		}
		return nil, fmt.Errorf("failed to find school: %w", err)
	}

	if school.AdminUserID != uuid.Nil {
		adminUser, err := s.userRepo.FindByID(school.AdminUserID)
		if err != nil {
			return nil, fmt.Errorf("failed to find admin user for school: %w", err)
		}

		latestVerification, err := s.emailVerifyRepo.FindByUserID(adminUser.ID)
		if err != nil || latestVerification == nil || !latestVerification.IsVerified {
			return nil, errors.New("admin email not verified yet")
		}
	} else {
		return nil, errors.New("admin user not associated with school")
	}

	// If Free Trial, set subscription start/end dates if not already set by SelectPackage
	if school.Package.Name == "Free Trial" && school.SubscriptionStartDate == nil {
		now := time.Now()
		school.SubscriptionStartDate = &now
		expiry := now.Add(time.Duration(*school.Package.DurationDays) * 24 * time.Hour)
		school.SubscriptionEndDate = &expiry
		school.UpdatedBy = uuid.Nil
		err = s.schoolRepo.Update(school)
		if err != nil {
			return nil, fmt.Errorf("failed to finalize free trial subscription dates: %w", err)
		}
	}

	return school, nil
}

func (s *registrationService) GetSchoolByID(schoolID uuid.UUID) (*models.School, error) {
	school, err := s.schoolRepo.FindByID(schoolID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("school not found")
		}
		return nil, fmt.Errorf("failed to retrieve school: %w", err)
	}
	return school, nil
}

func (s *registrationService) GetPackageByID(packageID uuid.UUID) (*models.Package, error) {
	pkg, err := s.packageRepo.FindByID(packageID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("package not found")
		}
		return nil, fmt.Errorf("failed to retrieve package: %w", err)
	}
	return pkg, nil
}
