package handlers

import (
	"net/http"

	"auth-barniee/internal/repositories"
	"auth-barniee/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RegistrationHandler struct {
	regService services.RegistrationService
}

func NewRegistrationHandler(regService services.RegistrationService) *RegistrationHandler {
	return &RegistrationHandler{regService: regService}
}

type RegisterSchoolInfoRequest struct {
	Name                string `json:"name" binding:"required"`
	EducationLevel      string `json:"education_level" binding:"required,oneof=SD SMP SMA SMK PerguruanTinggi Lainnya"`
	Status              string `json:"status" binding:"required,oneof=Negeri Swasta"`
	Address             string `json:"address" binding:"required"`
	InitialStudentCount int    `json:"initial_student_count" binding:"required,min=1"`
}

type RegisterAdminInfoRequest struct {
	SchoolID       uuid.UUID `json:"school_id" binding:"required"`
	AdminName      string    `json:"admin_name" binding:"required"`
	AdminEmail     string    `json:"admin_email" binding:"required,email"`
	WhatsappNumber string    `json:"whatsapp_number" binding:"required"`
	Position       string    `json:"position" binding:"required"`
}

type SelectPackageRequest struct {
	SchoolID  uuid.UUID `json:"school_id" binding:"required"`
	PackageID uuid.UUID `json:"package_id" binding:"required"`
}

type RequestOTPRequest struct {
	UserID uuid.UUID `json:"user_id" binding:"required"`
}

type VerifyOTPRequest struct {
	UserID uuid.UUID `json:"user_id" binding:"required"`
	OTP    string    `json:"otp" binding:"required,len=6"`
}

type CompleteRegistrationRequest struct {
	SchoolID uuid.UUID `json:"school_id" binding:"required"`
}

func (h *RegistrationHandler) RegisterSchoolInfo(c *gin.Context) {
	var req RegisterSchoolInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	school, err := h.regService.RegisterSchoolInfo(req.Name, req.EducationLevel, req.Status, req.Address, req.InitialStudentCount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "School info registered successfully", "school_id": school.ID, "school_name": school.Name})
}

func (h *RegistrationHandler) RegisterAdminInfo(c *gin.Context) {
	var req RegisterAdminInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	adminUser, generatedPassword, err := h.regService.RegisterAdminInfo(req.SchoolID, req.AdminName, req.AdminEmail, req.WhatsappNumber, req.Position)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message":   "Admin user created and linked to school",
		"user_id":   adminUser.ID,
		"email":     adminUser.Email,
		"password":  generatedPassword,
		"school_id": adminUser.SchoolID,
	})
}

func (h *RegistrationHandler) SelectPackage(c *gin.Context) {
	var req SelectPackageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	school, err := h.regService.SelectPackage(req.SchoolID, req.PackageID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Package selected successfully", "school": school})
}

func (h *RegistrationHandler) GetAllPackages(c *gin.Context) {
	db, ok := c.MustGet("db").(*gorm.DB)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection not found"})
		return
	}
	pkgRepo := repositories.NewPackageRepository(db)
	packages, err := pkgRepo.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve packages"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Packages retrieved successfully", "packages": packages})
}

func (h *RegistrationHandler) RequestEmailVerificationOTP(c *gin.Context) {
	var req RequestOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.regService.RequestEmailVerificationOTP(req.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "OTP sent to email successfully"})
}

func (h *RegistrationHandler) VerifyEmailOTP(c *gin.Context) {
	var req VerifyOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.regService.VerifyEmailOTP(req.UserID, req.OTP)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Email verified successfully"})
}

func (h *RegistrationHandler) CompleteRegistration(c *gin.Context) {
	var req CompleteRegistrationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	school, err := h.regService.CompleteRegistration(req.SchoolID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "School registration completed successfully", "school": school})
}
