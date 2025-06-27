package handlers

import (
	"net/http"

	"auth-barniee/internal/models"
	"auth-barniee/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	userService services.UserService
}

func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// CreateUserRequest represents the request body for creating a user.
type CreateUserRequest struct {
	Name     string `json:"name" binding:"required" example:"Teacher John"`
	Email    string `json:"email" binding:"required,email" example:"john@example.com"`
	Password string `json:"password" binding:"required,min=6" example:"securepassword"`
	RoleName string `json:"role_name" binding:"required,oneof=teacher student" example:"teacher"`
}

// UpdateUserRequest represents the request body for updating a user.
type UpdateUserRequest struct {
	Name     *string `json:"name" example:"John Doe"`
	Email    *string `json:"email" example:"john.doe@example.com"`
	RoleName *string `json:"role_name,omitempty" binding:"omitempty,oneof=teacher student admin" example:"student"`
}

// UserDataResponse represents a single user's data for API response.
type UserDataResponse struct {
	User models.User `json:"user"`
}

// UserListResponse represents a list of users' data for API response.
type UserListResponse struct {
	Users []models.User `json:"users"`
}

// @Summary Create Teacher or Student
// @Description Allows an admin to create a new teacher or student account within their school.
// @Tags Admin - User Management
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param createUserRequest body CreateUserRequest true "User details to create"
// @Success 201 {object} CommonResponse{data=UserDataResponse} "User created successfully"
// @Failure 400 {object} CommonResponse "Bad request"
// @Failure 401 {object} CommonResponse "Unauthorized"
// @Failure 403 {object} CommonResponse "Forbidden"
// @Failure 500 {object} CommonResponse "Internal server error"
// @Router /admin/users [post]
func (h *UserHandler) CreateTeacherOrStudent(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, CommonResponse{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	adminID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, CommonResponse{
			Status:  http.StatusUnauthorized,
			Message: "Admin ID not found in context",
			Data:    nil,
		})
		return
	}
	adminUUID, ok := adminID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, CommonResponse{
			Status:  http.StatusInternalServerError,
			Message: "Invalid admin ID type in context",
			Data:    nil,
		})
		return
	}

	user, err := h.userService.CreateTeacherOrStudent(req.Name, req.Email, req.Password, req.RoleName, adminUUID)
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
		Message: "User created successfully",
		Data:    UserDataResponse{User: *user},
	})
}

// @Summary Get All Users
// @Description Retrieves a list of all users, with optional filtering by role. Accessible by admins.
// @Tags Admin - User Management
// @Security BearerAuth
// @Produce json
// @Param role query string false "Filter users by role (teacher, student, admin)" example:"teacher"
// @Success 200 {object} CommonResponse{data=UserListResponse} "Users retrieved successfully"
// @Failure 401 {object} CommonResponse "Unauthorized"
// @Failure 403 {object} CommonResponse "Forbidden"
// @Failure 500 {object} CommonResponse "Internal server error"
// @Router /admin/users [get]
func (h *UserHandler) GetAllUsers(c *gin.Context) {
	roleName := c.Query("role")

	adminID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, CommonResponse{
			Status:  http.StatusUnauthorized,
			Message: "Admin ID not found in context",
			Data:    nil,
		})
		return
	}
	adminUUID, ok := adminID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, CommonResponse{
			Status:  http.StatusInternalServerError,
			Message: "Invalid admin ID type in context",
			Data:    nil,
		})
		return
	}

	users, err := h.userService.GetAllUsers(roleName, adminUUID)
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
		Message: "Users retrieved successfully",
		Data:    UserListResponse{Users: users},
	})
}

// @Summary Get User By ID
// @Description Retrieves a specific user's details by their ID. Accessible by admins.
// @Tags Admin - User Management
// @Security BearerAuth
// @Produce json
// @Param id path string true "User ID" format:"uuid" example:"f1e2d3c4-b5a6-9876-5432-10fedcba9876"
// @Success 200 {object} CommonResponse{data=UserDataResponse} "User retrieved successfully"
// @Failure 400 {object} CommonResponse "Bad request"
// @Failure 401 {object} CommonResponse "Unauthorized"
// @Failure 403 {object} CommonResponse "Forbidden"
// @Failure 404 {object} CommonResponse "User not found"
// @Failure 500 {object} CommonResponse "Internal server error"
// @Router /admin/users/{id} [get]
func (h *UserHandler) GetUserByID(c *gin.Context) {
	userIDParam := c.Param("id")
	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, CommonResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid user ID format",
			Data:    nil,
		})
		return
	}

	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "user not found" {
			statusCode = http.StatusNotFound
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
		Message: "User retrieved successfully",
		Data:    UserDataResponse{User: *user},
	})
}

// @Summary Update User
// @Description Updates the details of an existing user. Accessible by admins.
// @Tags Admin - User Management
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "User ID" format:"uuid" example:"f1e2d3c4-b5a6-9876-5432-10fedcba9876"
// @Param updateUserRequest body UpdateUserRequest true "User details to update"
// @Success 200 {object} CommonResponse{data=UserDataResponse} "User updated successfully"
// @Failure 400 {object} CommonResponse "Bad request"
// @Failure 401 {object} CommonResponse "Unauthorized"
// @Failure 403 {object} CommonResponse "Forbidden"
// @Failure 404 {object} CommonResponse "User not found"
// @Failure 500 {object} CommonResponse "Internal server error"
// @Router /admin/users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	userIDParam := c.Param("id")
	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, CommonResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid user ID format",
			Data:    nil,
		})
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, CommonResponse{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	adminID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, CommonResponse{
			Status:  http.StatusUnauthorized,
			Message: "Admin ID not found in context",
			Data:    nil,
		})
		return
	}
	adminUUID, ok := adminID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, CommonResponse{
			Status:  http.StatusInternalServerError,
			Message: "Invalid admin ID type in context",
			Data:    nil,
		})
		return
	}

	updatedUser, err := h.userService.UpdateUser(userID, adminUUID, req.Name, req.Email, req.RoleName)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "user not found" || err.Error() == "email already taken by another user" || err.Error() == "role '...' not found" || err.Error() == "unauthorized: school admin cannot update users outside their school" {
			statusCode = http.StatusBadRequest // Or 403 for forbidden
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
		Message: "User updated successfully",
		Data:    UserDataResponse{User: *updatedUser},
	})
}

// @Summary Delete User
// @Description Deletes a user by their ID. Accessible by admins.
// @Tags Admin - User Management
// @Security BearerAuth
// @Produce json
// @Param id path string true "User ID" format:"uuid" example:"f1e2d3c4-b5a6-9876-5432-10fedcba9876"
// @Success 200 {object} CommonResponse "User deleted successfully"
// @Failure 400 {object} CommonResponse "Bad request"
// @Failure 401 {object} CommonResponse "Unauthorized"
// @Failure 403 {object} CommonResponse "Forbidden"
// @Failure 404 {object} CommonResponse "User not found"
// @Failure 500 {object} CommonResponse "Internal server error"
// @Router /admin/users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	userIDParam := c.Param("id")
	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, CommonResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid user ID format",
			Data:    nil,
		})
		return
	}

	adminID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, CommonResponse{
			Status:  http.StatusUnauthorized,
			Message: "Admin ID not found in context",
			Data:    nil,
		})
		return
	}
	adminUUID, ok := adminID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, CommonResponse{
			Status:  http.StatusInternalServerError,
			Message: "Invalid admin ID type in context",
			Data:    nil,
		})
		return
	}

	err = h.userService.DeleteUser(userID, adminUUID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "user not found" || err.Error() == "cannot delete your own admin account" || err.Error() == "unauthorized: school admin cannot delete users outside their school" {
			statusCode = http.StatusBadRequest // Or 403 for forbidden
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
		Message: "User deleted successfully",
		Data:    nil,
	})
}
