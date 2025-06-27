package handlers

import (
	"net/http"

	"auth-barniee/internal/models"
	"auth-barniee/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthHandler struct {
	authService services.AuthService
}

func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// CommonResponse represents the standardized API response structure.
type CommonResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// LoginRequest represents the request body for login.
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"user@example.com"`
	Password string `json:"password" binding:"required" example:"password123"`
}

// LoginResponseData represents the data returned upon successful login.
type LoginResponseData struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// UserProfileResponseData represents the data returned for user profile.
type UserProfileResponseData struct {
	User   models.User    `json:"user"`
	School *models.School `json:"school,omitempty"`
}

// @Summary User Login
// @Description Authenticates a user and returns a JWT token.
// @Tags Auth
// @Accept json
// @Produce json
// @Param loginRequest body LoginRequest true "Login Credentials"
// @Success 200 {object} CommonResponse{data=LoginResponseData} "Login successful"
// @Failure 400 {object} CommonResponse "Bad request"
// @Failure 401 {object} CommonResponse "Unauthorized"
// @Failure 500 {object} CommonResponse "Internal server error"
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, CommonResponse{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	token, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, CommonResponse{
			Status:  http.StatusUnauthorized,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, CommonResponse{
		Status:  http.StatusOK,
		Message: "Login successful",
		Data:    LoginResponseData{Token: token},
	})
}

// @Summary User Logout
// @Description Invalidates the client-side JWT token (no server-side session invalidation for stateless JWT).
// @Tags Auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} CommonResponse "Logout successful"
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, CommonResponse{
		Status:  http.StatusOK,
		Message: "Logout successful (token discarded by client)",
		Data:    nil,
	})
}

// @Summary Get User Profile
// @Description Retrieves the basic profile information of the authenticated user.
// @Tags Auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} CommonResponse{data=UserProfileResponseData} "User profile retrieved successfully"
// @Failure 401 {object} CommonResponse "Unauthorized"
// @Failure 500 {object} CommonResponse "Internal server error"
// @Router /profile [get]
func (h *AuthHandler) GetUserProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, CommonResponse{
			Status:  http.StatusUnauthorized,
			Message: "User ID not found in context",
			Data:    nil,
		})
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, CommonResponse{
			Status:  http.StatusInternalServerError,
			Message: "Invalid user ID type in context",
			Data:    nil,
		})
		return
	}

	user, school, err := h.authService.GetUserProfile(userUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, CommonResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	responseData := UserProfileResponseData{
		User:   *user,
		School: school,
	}

	c.JSON(http.StatusOK, CommonResponse{
		Status:  http.StatusOK,
		Message: "User profile retrieved successfully",
		Data:    responseData,
	})
}
