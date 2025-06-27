package handlers

import (
	"net/http"

	"auth-barniee/internal/models"
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

// Request models for multi-step registration
type RegisterSchoolInfoRequest struct {
	Name                string `json:"name" binding:"required" example:"Barniee Academy"`
	EducationLevel      string `json:"education_level" binding:"required,oneof=SD SMP SMA SMK PerguruanTinggi Lainnya" example:"SMA"`
	Status              string `json:"status" binding:"required,oneof=Negeri Swasta" example:"Swasta"`
	Address             string `json:"address" binding:"required" example:"Jl. Inovasi No. 10, Kota Teknologi"`
	InitialStudentCount int    `json:"initial_student_count" binding:"required,min=1" example:"150"`
}

type RegisterAdminInfoRequest struct {
	SchoolID       uuid.UUID `json:"school_id" binding:"required" example:"a1b2c3d4-e5f6-7890-1234-567890abcdef"`
	AdminName      string    `json:"admin_name" binding:"required" example:"Siti Aminah"`
	AdminEmail     string    `json:"admin_email" binding:"required,email" example:"siti.aminah@example.com"`
	WhatsappNumber string    `json:"whatsapp_number" binding:"required" example:"081234567890"`
	Position       string    `json:"position" binding:"required" example:"Direktur"`
}

type SelectPackageRequest struct {
	SchoolID  uuid.UUID `json:"school_id" binding:"required" example:"a1b2c3d4-e5f6-7890-1234-567890abcdef"`
	PackageID uuid.UUID `json:"package_id" binding:"required" example:"package-uuid-for-premium"`
}

type RequestOTPRequest struct {
	UserID uuid.UUID `json:"user_id" binding:"required" example:"f1e2d3c4-b5a6-9876-5432-10fedcba9876"`
}

type VerifyOTPRequest struct {
	UserID uuid.UUID `json:"user_id" binding:"required" example:"f1e2d3c4-b5a6-9876-5432-10fedcba9876"`
	OTP    string    `json:"otp" binding:"required,len=6" example:"123456"`
}

type CompleteRegistrationRequest struct {
	SchoolID uuid.UUID `json:"school_id" binding:"required" example:"a1b2c3d4-e5f6-7890-1234-567890abcdef"`
}

// RegisterSchoolInfoResponseData represents the data returned after registering school info.
type RegisterSchoolInfoResponseData struct {
	SchoolID   uuid.UUID `json:"school_id" example:"a1b2c3d4-e5f6-7890-1234-567890abcdef"`
	SchoolName string    `json:"school_name" example:"Barniee Academy"`
}

// RegisterAdminInfoResponseData represents the data returned after registering admin info.
type RegisterAdminInfoResponseData struct {
	UserID   uuid.UUID `json:"user_id" example:"f1e2d3c4-b5a6-9876-5432-10fedcba9876"`
	Email    string    `json:"email" example:"siti.aminah@example.com"`
	Password string    `json:"password" example:"GeneratedPass123"`
	SchoolID uuid.UUID `json:"school_id" example:"a1b2c3d4-e5f6-7890-1234-567890abcdef"`
}

// SelectPackageResponseData represents the data returned after selecting a package.
type SelectPackageResponseData struct {
	School models.School `json:"school"`
}

// GetAllPackagesResponseData represents the data returned for all available packages.
type GetAllPackagesResponseData struct {
	Packages []models.Package `json:"packages"`
}

// CompleteRegistrationResponseData represents the data returned after completing registration.
type CompleteRegistrationResponseData struct {
	School models.School `json:"school"`
}

// @Summary Register School Information
// @Description Step 1 of school registration: Register basic school details.
// @Tags School Registration
// @Accept json
// @Produce json
// @Param registerSchoolInfoRequest body RegisterSchoolInfoRequest true "School Information"
// @Success 201 {object} CommonResponse{data=RegisterSchoolInfoResponseData} "School info registered successfully"
// @Failure 400 {object} CommonResponse "Bad request"
// @Failure 500 {object} CommonResponse "Internal server error"
// @Router /register/school-info [post]
func (h *RegistrationHandler) RegisterSchoolInfo(c *gin.Context) {
	var req RegisterSchoolInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, CommonResponse{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	school, err := h.regService.RegisterSchoolInfo(req.Name, req.EducationLevel, req.Status, req.Address, req.InitialStudentCount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, CommonResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusCreated, CommonResponse{
		Status:  http.StatusCreated,
		Message: "School info registered successfully",
		Data: RegisterSchoolInfoResponseData{
			SchoolID:   school.ID,
			SchoolName: school.Name,
		},
	})
}

// @Summary Register Admin Information
// @Description Step 2 of school registration: Register the primary admin user for the school.
// @Tags School Registration
// @Accept json
// @Produce json
// @Param registerAdminInfoRequest body RegisterAdminInfoRequest true "Admin Information"
// @Success 201 {object} CommonResponse{data=RegisterAdminInfoResponseData} "Admin user created and linked to school"
// @Failure 400 {object} CommonResponse "Bad request"
// @Failure 500 {object} CommonResponse "Internal server error"
// @Router /register/admin-info [post]
func (h *RegistrationHandler) RegisterAdminInfo(c *gin.Context) {
	var req RegisterAdminInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, CommonResponse{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	adminUser, generatedPassword, err := h.regService.RegisterAdminInfo(req.SchoolID, req.AdminName, req.AdminEmail, req.WhatsappNumber, req.Position)
	if err != nil {
		c.JSON(http.StatusInternalServerError, CommonResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusCreated, CommonResponse{
		Status:  http.StatusCreated,
		Message: "Admin user created and linked to school",
		Data: RegisterAdminInfoResponseData{
			UserID:   adminUser.ID,
			Email:    adminUser.Email,
			Password: generatedPassword,
			SchoolID: adminUser.SchoolID,
		},
	})
}

// @Summary Get All Available Packages
// @Description Retrieves a list of all available subscription packages (Free Trial, Premium, Enterprise).
// @Tags School Registration
// @Produce json
// @Success 200 {object} CommonResponse{data=GetAllPackagesResponseData} "Packages retrieved successfully"
// @Failure 500 {object} CommonResponse "Internal server error"
// @Router /register/packages [get]
func (h *RegistrationHandler) GetAllPackages(c *gin.Context) {
	db, ok := c.MustGet("db").(*gorm.DB)
	if !ok {
		c.JSON(http.StatusInternalServerError, CommonResponse{
			Status:  http.StatusInternalServerError,
			Message: "Database connection not found in context",
			Data:    nil,
		})
		return
	}
	pkgRepo := repositories.NewPackageRepository(db)
	packages, err := pkgRepo.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, CommonResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed to retrieve packages: " + err.Error(),
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, CommonResponse{
		Status:  http.StatusOK,
		Message: "Packages retrieved successfully",
		Data:    GetAllPackagesResponseData{Packages: packages},
	})
}

// @Summary Select Package
// @Description Step 3 of school registration: Selects a subscription package for the school.
// @Tags School Registration
// @Accept json
// @Produce json
// @Param selectPackageRequest body SelectPackageRequest true "Package Selection"
// @Success 200 {object} CommonResponse{data=SelectPackageResponseData} "Package selected successfully"
// @Failure 400 {object} CommonResponse "Bad request"
// @Failure 500 {object} CommonResponse "Internal server error"
// @Router /register/select-package [post]
func (h *RegistrationHandler) SelectPackage(c *gin.Context) {
	var req SelectPackageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, CommonResponse{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	school, err := h.regService.SelectPackage(req.SchoolID, req.PackageID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, CommonResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, CommonResponse{
		Status:  http.StatusOK,
		Message: "Package selected successfully",
		Data:    SelectPackageResponseData{School: *school},
	})
}

// @Summary Request Email Verification OTP
// @Description Step 4 of school registration: Sends an OTP to the user's email for verification.
// @Tags School Registration
// @Accept json
// @Produce json
// @Param requestOTPRequest body RequestOTPRequest true "User ID for OTP request"
// @Success 200 {object} CommonResponse "OTP sent to email successfully"
// @Failure 400 {object} CommonResponse "Bad request"
// @Failure 500 {object} CommonResponse "Internal server error"
// @Router /register/email-verification/request-otp [post]
func (h *RegistrationHandler) RequestEmailVerificationOTP(c *gin.Context) {
	var req RequestOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, CommonResponse{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	err := h.regService.RequestEmailVerificationOTP(req.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, CommonResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, CommonResponse{
		Status:  http.StatusOK,
		Message: "OTP sent to email successfully",
		Data:    nil,
	})
}

// @Summary Verify Email OTP
// @Description Step 4 of school registration: Verifies the OTP sent to the user's email.
// @Tags School Registration
// @Accept json
// @Produce json
// @Param verifyOTPRequest body VerifyOTPRequest true "User ID and OTP for verification"
// @Success 200 {object} CommonResponse "Email verified successfully"
// @Failure 400 {object} CommonResponse "Bad request (invalid/expired OTP)"
// @Failure 500 {object} CommonResponse "Internal server error"
// @Router /register/email-verification/verify-otp [post]
func (h *RegistrationHandler) VerifyEmailOTP(c *gin.Context) {
	var req VerifyOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, CommonResponse{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	err := h.regService.VerifyEmailOTP(req.UserID, req.OTP)
	if err != nil {
		// Differentiate between user error (invalid OTP) and internal error
		statusCode := http.StatusInternalServerError
		if err.Error() == "invalid OTP or user ID" || err.Error() == "email already verified" || err.Error() == "OTP has expired" {
			statusCode = http.StatusBadRequest
		}
		c.JSON(statusCode, CommonResponse{
			Status:  statusCode,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, CommonResponse{
		Status:  http.StatusOK,
		Message: "Email verified successfully",
		Data:    nil,
	})
}

// @Summary Complete School Registration
// @Description Step 6 of school registration: Finalizes the registration process after all previous steps are complete (including payment if applicable).
// @Tags School Registration
// @Accept json
// @Produce json
// @Param completeRegistrationRequest body CompleteRegistrationRequest true "School ID to complete registration for"
// @Success 200 {object} CommonResponse{data=CompleteRegistrationResponseData} "School registration completed successfully"
// @Failure 400 {object} CommonResponse "Bad request"
// @Failure 500 {object} CommonResponse "Internal server error"
// @Router /register/complete [post]
func (h *RegistrationHandler) CompleteRegistration(c *gin.Context) {
	var req CompleteRegistrationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, CommonResponse{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	school, err := h.regService.CompleteRegistration(req.SchoolID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, CommonResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, CommonResponse{
		Status:  http.StatusOK,
		Message: "School registration completed successfully",
		Data:    CompleteRegistrationResponseData{School: *school},
	})
}
