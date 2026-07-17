package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/jordisetiawan/insurance-auth-service/internal/dto"
	"github.com/jordisetiawan/insurance-auth-service/internal/service"
	"github.com/jordisetiawan/insurance-auth-service/internal/utils"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(
	authService *service.AuthService,
) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// @Summary Register a new user
// @Tags Auth
// @Accept json
// @Produce json
// @Param register body dto.RegisterRequest true "Register Data"
// @Success 201 {object} utils.APIResponse{data=dto.UserResponse}
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	user, err := h.authService.Register(req)
	if err != nil {
		if errors.Is(err, service.ErrEmailAlreadyExists) {
			utils.ErrorResponse(c, http.StatusConflict, err.Error(), nil)
		} else {
			utils.Log.Error("Registration failed", zap.Error(err))
			utils.ErrorResponse(c, http.StatusInternalServerError, "Internal Server Error", nil)
		}
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "User registered successfully", dto.UserResponse{
		ID:       user.ID.String(),
		Email:    user.Email,
		FullName: user.FullName,
		Role:     user.Role,
	})
}

// @Summary Login user
// @Tags Auth
// @Accept json
// @Produce json
// @Param login body dto.LoginRequest true "Login Credentials"
// @Success 200 {object} dto.LoginResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	user, permissions, accessToken, refreshToken, err := h.authService.Login(req)
	if err != nil {
		status := http.StatusInternalServerError

		switch {
		case errors.Is(err, service.ErrInvalidCredentials):
			status = http.StatusUnauthorized
		case errors.Is(err, service.ErrUserInactive):
			status = http.StatusForbidden
		default:
			utils.Log.Error("Login error", zap.Error(err))
		}

		utils.ErrorResponse(c, status, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Login successful", dto.LoginResponse{
		User: dto.UserResponse{
			ID:          user.ID.String(),
			Email:       user.Email,
			FullName:    user.FullName,
			Role:        user.Role,
			Permissions: permissions,
		},
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

// @Summary Refresh access token
// @Tags Auth
// @Accept json
// @Produce json
// @Param refresh body dto.RefreshRequest true "Refresh Token"
// @Success 200 {object} map[string]string
// @Router /auth/refresh [post]
func (h *AuthHandler) Refresh(c *gin.Context) {
	var req dto.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	newAccessToken, err := h.authService.Refresh(req.RefreshToken)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid refresh token", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Token refreshed", gin.H{"access_token": newAccessToken})
}

// @Summary User logout
// @Tags Auth
// @Produce json
// @Success 200 {object} map[string]string
// @Security BearerAuth
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	// Pada JWT stateless, logout biasanya dilakukan di client-side dengan menghapus token.
	// Namun, di sini kita bisa menambahkan logika blacklist di Redis jika diperlukan di masa depan.
	utils.SuccessResponse(c, http.StatusOK, "Successfully logged out", nil)
}

// @Summary Get current user info
// @Tags Auth
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /auth/me [get]
func (h *AuthHandler) GetMe(c *gin.Context) {
	userID, okID := c.Get("user_id")
	email, okEmail := c.Get("email")
	role, okRole := c.Get("role")
	perms, okPerms := c.Get("permissions")
	fullName, okFullName := c.Get("full_name")

	if !okID || !okEmail || !okRole || !okPerms || !okFullName {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to resolve current user context", nil)
		return
	}

	userIDStr, ok1 := userID.(string)
	emailStr, ok2 := email.(string)
	roleStr, ok3 := role.(string)
	fullNameStr, ok4 := fullName.(string)
	permsSlice, ok5 := perms.([]string)
	if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to resolve current user context", nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Current user info fetched", dto.UserResponse{
		ID:          userIDStr,
		Email:       emailStr,
		FullName:    fullNameStr,
		Role:        roleStr,
		Permissions: permsSlice,
	})
}

// @Summary Admin dashboard data
// @Tags Admin
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /admin/dashboard [get]
func (h *AuthHandler) AdminDashboard(c *gin.Context) {
	utils.SuccessResponse(c, http.StatusOK, "Welcome to Admin Dashboard", gin.H{
		"stats":       "Summary statistics for admin",
		"permissions": []string{"manage_users", "view_all_quotations", "system_audit"},
	})
}

// @Summary List all users
// @Tags User Management
// @Security BearerAuth
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Param search query string false "Search by name or email"
// @Success 200 {object} utils.APIResponse
// @Router /users [get]
func (h *AuthHandler) GetUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")

	users, total, err := h.authService.GetAllUsers(page, limit, search)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch users", err.Error())
		return
	}

	var res []dto.UserResponse
	for _, u := range users {
		res = append(res, dto.UserResponse{
			ID:       u.ID.String(),
			Email:    u.Email,
			FullName: u.FullName,
			Role:     u.Role,
		})
	}

	utils.SuccessResponse(c, http.StatusOK, "Users fetched successfully", gin.H{
		"users": res,
		"meta": gin.H{
			"total": total,
			"page":  page,
			"limit": limit,
		},
	})
}

// @Summary Get user detail
// @Tags User Management
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} utils.APIResponse
// @Router /users/{id} [get]
func (h *AuthHandler) GetUserByID(c *gin.Context) {
	id := c.Param("id")
	user, err := h.authService.GetUserByID(id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "User not found", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "User detail fetched", dto.UserResponse{
		ID:       user.ID.String(),
		Email:    user.Email,
		FullName: user.FullName,
		Role:     user.Role,
	})
}

// @Summary Update user
// @Tags User Management
// @Security BearerAuth
// @Param id path string true "User ID"
// @Param body body dto.RegisterRequest true "Update data"
// @Success 200 {object} utils.APIResponse
// @Router /users/{id} [put]
func (h *AuthHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	user, err := h.authService.UpdateUser(id, req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Update failed", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "User updated", user)
}

// @Summary Assign role to user
// @Tags User Management
// @Security BearerAuth
// @Param id path string true "User ID"
// @Param role query string true "Role Code"
// @Success 200 {object} utils.APIResponse
// @Router /users/{id}/role [patch]
func (h *AuthHandler) AssignRole(c *gin.Context) {
	id := c.Param("id")
	role := c.Query("role")
	if err := h.authService.AssignRole(id, role); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Assignment failed", err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Role assigned successfully", nil)
}

// @Summary User dashboard data
// @Tags User
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /user/dashboard [get]
func (h *AuthHandler) UserDashboard(c *gin.Context) {
	userID, _ := c.Get("user_id")
	utils.SuccessResponse(c, http.StatusOK, "Welcome to User Dashboard", gin.H{
		"user_id":  userID,
		"features": []string{"my_policies", "submit_claim", "profile_settings"},
	})
}

// @Summary Assign permission to role
// @Tags Role Management
// @Security BearerAuth
// @Param code path string true "Role Code"
// @Param permission query string true "Permission Code"
// @Success 200 {object} utils.APIResponse
// @Router /roles/{code}/permissions [post]
func (h *AuthHandler) AssignPermission(c *gin.Context) {
	roleCode := c.Param("code")
	permCode := c.Query("permission")
	if err := h.authService.AssignPermissionToRole(roleCode, permCode); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Permission assignment failed", err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Permission assigned to role successfully", nil)
}
