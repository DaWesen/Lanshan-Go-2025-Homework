package api

import (
	"net/http"
	"time"

	"md7/model"
	"md7/sql"
	"md7/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthHandler struct {
	db *gorm.DB
}

func NewAuthHandler() *AuthHandler {
	db, _ := sql.InitDatabase()
	return &AuthHandler{db: db.DB}
}

type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	School   string `json:"school"`
	Password string `json:"password" binding:"required,min=6"`
}
type LoginRequest struct {
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
}
type AuthResponse struct {
	AccessToken  string           `json:"access_token"`
	RefreshToken string           `json:"refresh_token"`
	ExpiresAt    time.Time        `json:"expires_at"`
	Student      *StudentResponse `json:"student"`
}
type StudentResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	School    string    `json:"school"`
	CreatedAt time.Time `json:"created_at"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  400,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}
	var existingStudent model.Student
	if err := h.db.Where("name = ?", req.Name).First(&existingStudent).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"status":  409,
			"message": "用户已存在",
		})
		return
	}
	student := model.Student{
		Name:     req.Name,
		School:   req.School,
		Password: req.Password,
	}
	if err := student.HashPassword(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  500,
			"message": "密码加密失败",
			"error":   err.Error(),
		})
		return
	}
	if err := h.db.Create(&student).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  500,
			"message": "注册失败",
			"error":   err.Error(),
		})
		return
	}
	accessToken, err := utils.GenerateAccessToken(student.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  500,
			"message": "生成token失败",
			"error":   err.Error(),
		})
		return
	}
	refreshToken, err := utils.GenerateRefreshToken(student.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  500,
			"message": "生成refresh token失败",
			"error":   err.Error(),
		})
		return
	}
	response := AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(15 * time.Minute),
		Student: &StudentResponse{
			ID:        student.ID,
			Name:      student.Name,
			School:    student.School,
			CreatedAt: student.CreatedAt,
		},
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  201,
		"message": "注册成功",
		"data":    response,
	})
}
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  400,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}
	var student model.Student
	if err := h.db.Where("name = ?", req.Name).First(&student).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  401,
			"message": "用户名或密码错误",
		})
		return
	}
	if !student.CheckPassword(req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  401,
			"message": "用户名或密码错误",
		})
		return
	}
	accessToken, err := utils.GenerateAccessToken(student.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  500,
			"message": "生成token失败",
			"error":   err.Error(),
		})
		return
	}
	refreshToken, err := utils.GenerateRefreshToken(student.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  500,
			"message": "生成refresh token失败",
			"error":   err.Error(),
		})
		return
	}
	response := AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(15 * time.Minute),
		Student: &StudentResponse{
			ID:        student.ID,
			Name:      student.Name,
			School:    student.School,
			CreatedAt: student.CreatedAt,
		},
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "登录成功",
		"data":    response,
	})
}
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  400,
			"message": "缺少token",
		})
		return
	}
	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}
	claims, err := utils.ParseToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  401,
			"message": "无效的token",
			"error":   err.Error(),
		})
		return
	}

	if claims.Subject != "refresh" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  401,
			"message": "请使用refresh token",
		})
		return
	}

	accessToken, err := utils.GenerateAccessToken(claims.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  500,
			"message": "生成token失败",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "token刷新成功",
		"data": gin.H{
			"access_token": accessToken,
			"expires_at":   time.Now().Add(15 * time.Minute),
		},
	})
}

func (h *AuthHandler) GetProfile(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  401,
			"message": "未认证的用户",
		})
		return
	}
	var student model.Student
	if err := h.db.Where("name = ?", username).First(&student).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  404,
			"message": "用户不存在",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "获取用户信息成功",
		"data": StudentResponse{
			ID:        student.ID,
			Name:      student.Name,
			School:    student.School,
			CreatedAt: student.CreatedAt,
		},
	})
}
