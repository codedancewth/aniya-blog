package handlers

import (
	"net/http"

	"github.com/cworld1/aniya-blog/backend/internal/models"
	"github.com/cworld1/aniya-blog/backend/internal/repository"
	"github.com/cworld1/aniya-blog/backend/pkg/auth"
	"github.com/cworld1/aniya-blog/backend/pkg/response"
	"github.com/cworld1/aniya-blog/backend/pkg/validator"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	userRepo   *repository.UserRepository
	jwtManager *auth.JWTManager
}

// NewAuthHandler 创建认证处理器
func NewAuthHandler(userRepo *repository.UserRepository, jwtManager *auth.JWTManager) *AuthHandler {
	return &AuthHandler{
		userRepo:   userRepo,
		jwtManager: jwtManager,
	}
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=20"`
	Password string `json:"password" binding:"required,min=6"`
	Email    string `json:"email" binding:"required,email"`
	Nickname string `json:"nickname"`
}

// Login 用户登录
// @Summary 用户登录
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body LoginRequest true "登录信息"
// @Success 200 {object} response.Response
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.ERROR, validator.GetErrorMsg(err))
		return
	}

	// 查找用户
	user, err := h.userRepo.FindByUsername(req.Username)
	if err != nil {
		response.Error(c, response.USER_NOT_FOUND, "invalid username or password")
		return
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		response.Error(c, response.INVALID_PASSWORD, "invalid username or password")
		return
	}

	// 检查用户状态
	if user.Status != 1 {
		response.Error(c, response.USER_DISABLED, "user is disabled")
		return
	}

	// 生成令牌
	token, err := h.jwtManager.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		response.Error(c, response.ERROR, "failed to generate token")
		return
	}

	response.Success(c, gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"nickname": user.Nickname,
			"avatar":   user.Avatar,
			"role":     user.Role,
		},
	})
}

// Register 用户注册
// @Summary 用户注册
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "注册信息"
// @Success 200 {object} response.Response
// @Router /api/v1/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.ERROR, validator.GetErrorMsg(err))
		return
	}

	// 检查用户名是否已存在
	if exists, _ := h.userRepo.UsernameExists(req.Username); exists {
		response.Error(c, response.ERROR, "username already exists")
		return
	}

	// 检查邮箱是否已存在
	if exists, _ := h.userRepo.EmailExists(req.Email); exists {
		response.Error(c, response.ERROR, "email already exists")
		return
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		response.Error(c, response.ERROR, "failed to hash password")
		return
	}

	// 创建用户
	user := &models.User{
		Username: req.Username,
		Password: string(hashedPassword),
		Email:    req.Email,
		Nickname: req.Nickname,
		Role:     "user",
		Status:   1,
	}

	if err := h.userRepo.Create(user); err != nil {
		response.Error(c, response.ERROR, "failed to create user")
		return
	}

	// 生成令牌
	token, err := h.jwtManager.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		response.Error(c, response.ERROR, "failed to generate token")
		return
	}

	response.Success(c, gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"nickname": user.Nickname,
			"avatar":   user.Avatar,
			"role":     user.Role,
		},
	})
}

// GetCurrentUser 获取当前用户信息
// @Summary 获取当前用户信息
// @Tags 认证
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response
// @Router /api/v1/auth/me [get]
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	userID, _ := c.Get("user_id")
	user, err := h.userRepo.FindByID(userID.(uint))
	if err != nil {
		response.Error(c, response.USER_NOT_FOUND, "user not found")
		return
	}

	response.Success(c, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
		"nickname": user.Nickname,
		"avatar":   user.Avatar,
		"role":     user.Role,
	})
}

// RefreshToken 刷新令牌
// @Summary 刷新令牌
// @Tags 认证
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response
// @Router /api/v1/auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	tokenString := authHeader[7:] // 去掉 "Bearer " 前缀

	newToken, err := h.jwtManager.RefreshToken(tokenString)
	if err != nil {
		response.Error(c, response.INVALID_TOKEN, "failed to refresh token")
		return
	}

	response.Success(c, gin.H{
		"token": newToken,
	})
}

// ChangePassword 修改密码
// @Summary 修改密码
// @Tags 认证
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body ChangePasswordRequest true "密码信息"
// @Success 200 {object} response.Response
// @Router /api/v1/auth/change-password [post]
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.ERROR, validator.GetErrorMsg(err))
		return
	}

	userID, _ := c.Get("user_id")
	user, err := h.userRepo.FindByID(userID.(uint))
	if err != nil {
		response.Error(c, response.USER_NOT_FOUND, "user not found")
		return
	}

	// 验证旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		response.Error(c, response.INVALID_PASSWORD, "old password is incorrect")
		return
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		response.Error(c, response.ERROR, "failed to hash password")
		return
	}

	user.Password = string(hashedPassword)
	if err := h.userRepo.Update(user); err != nil {
		response.Error(c, response.ERROR, "failed to update password")
		return
	}

	response.Success(c, nil)
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}
