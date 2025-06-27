package handlers

import (
	"net/http"

	"auth-barniee/internal/services"
	"auth-barniee/internal/utils" // Import new utils package

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthHandler struct {
	authService services.AuthService
}

func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// Login godoc
// @Summary User login
// @Description Authenticate user and return JWT token
// @Tags Auth
// @Accept json
// @Produce json
// @Param loginRequest body LoginRequest true "Login Credentials"
// @Success 200 {object} utils.APIResponse{data=object{token=string}} "Login successful"
// @Failure 400 {object} utils.APIResponse "Bad Request"
// @Failure 401 {object} utils.APIResponse "Unauthorized"
// @Failure 500 {object} utils.APIResponse "Internal Server Error"
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.NewErrorResponse(c, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	token, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		utils.NewErrorResponse(c, http.StatusUnauthorized, "Authentication failed", err)
		return
	}

	utils.NewSuccessResponse(c, http.StatusOK, "Login successful", gin.H{"token": token})
}

// Logout godoc
// @Summary User logout
// @Description Simulate user logout by instructing client to discard token
// @Tags Auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} utils.APIResponse "Logout successful"
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	utils.NewSuccessResponse(c, http.StatusOK, "Logout successful (token discarded by client)", nil)
}

// GetUserProfile godoc
// @Summary Get user profile
// @Description Retrieve authenticated user's basic profile information
// @Tags Users
// @Security BearerAuth
// @Produce json
// @Success 200 {object} utils.APIResponse{data=object{user=models.User,school=models.School}} "User profile retrieved successfully"
// @Failure 401 {object} utils.APIResponse "Unauthorized"
// @Failure 500 {object} utils.APIResponse "Internal Server Error"
// @Router /profile [get]
func (h *AuthHandler) GetUserProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.NewErrorResponse(c, http.StatusUnauthorized, "User ID not found in context", nil)
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		utils.NewErrorResponse(c, http.StatusInternalServerError, "Invalid user ID type in context", nil)
		return
	}

	user, school, err := h.authService.GetUserProfile(userUUID)
	if err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve user profile", err)
		return
	}

	responseData := gin.H{"user": user}
	if school != nil {
		responseData["school"] = school
	}
	utils.NewSuccessResponse(c, http.StatusOK, "User profile retrieved successfully", responseData)
}
