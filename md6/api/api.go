package api

import (
	"md6/dao"
	"md6/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func register(c *gin.Context) {
	name := c.PostForm("name")
	password := c.PostForm("password")
	school := c.PostForm("school")

	if name == "" || password == "" || school == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  400,
			"message": "姓名、密码和学校不能为空",
		})
		return
	}

	success := dao.Addstudent(name, password, school)
	if !success {
		c.JSON(http.StatusConflict, gin.H{
			"status":  409,
			"message": "用户已存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "用户注册成功",
		"data": gin.H{
			"name":   name,
			"school": school,
		},
	})
}

func login(c *gin.Context) {
	name := c.PostForm("name")
	password := c.PostForm("password")

	if name == "" || password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  400,
			"message": "姓名和密码不能为空",
		})
		return
	}

	exists := dao.Givestudent(name)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  401,
			"message": "用户不存在",
		})
		return
	}

	storedPassword := dao.Givepassword(name)
	if storedPassword != password {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  401,
			"message": "密码错误",
		})
		return
	}
	// 生成access token和refresh token
	accessToken, err := utils.GenerateAccessToken(name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  500,
			"message": "生成访问令牌失败",
		})
		return
	}

	refreshToken, err := utils.GenerateRefreshToken(name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  500,
			"message": "生成刷新令牌失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "登录成功",
		"data": gin.H{
			"name":         name,
			"accesstoken":  accessToken,
			"refreshtoken": refreshToken,
			"expires_in":   900,
		},
	})
}

// 修改密码 - 需要登录且只能修改自己的密码
func changePassword(c *gin.Context) {
	currentUsername, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  401,
			"message": "未授权访问",
		})
		return
	}
	name := c.PostForm("name")
	if name != currentUsername.(string) {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  403,
			"message": "只能修改自己的密码",
		})
		return
	}

	oldPassword := c.PostForm("oldpassword")
	newPassword := c.PostForm("newpassword")

	if name == "" || oldPassword == "" || newPassword == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  400,
			"message": "所有字段都是必需的",
		})
		return
	}

	storedPassword := dao.Givepassword(name)
	if storedPassword != oldPassword {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  401,
			"message": "旧密码错误",
		})
		return
	}

	success := dao.UpdatePassword(name, newPassword)
	if !success {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  404,
			"message": "用户不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "密码修改成功",
		"data": gin.H{
			"name": name,
		},
	})
}

func getstudents(c *gin.Context) {
	students := dao.Getallstudents()
	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"data": gin.H{
			"count": len(students),
			"users": students,
		},
	})
}

func verifyToken(c *gin.Context) {
	token := c.PostForm("token")

	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  400,
			"message": "token不能为空",
		})
		return
	}

	claims, err := utils.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  401,
			"message": "无效或过期的令牌",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "令牌有效",
		"data": gin.H{
			"name":  claims.Name,
			"valid": true,
		},
	})
}

func getstudent(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  401,
			"message": "未授权访问",
		})
		return
	}

	name := username.(string)
	students := dao.Getallstudents()
	student, exists := students[name]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  404,
			"message": "用户不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"data": gin.H{
			"name":   student.Name,
			"school": student.School,
		},
	})
}

// 刷新token
func refresh(c *gin.Context) {
	refreshToken := c.PostForm("refreshtoken")
	if refreshToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  400,
			"message": "refreshtoken不能为空",
		})
		return
	}

	claims, err := utils.ParseToken(refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  401,
			"message": "无效或过期的刷新令牌",
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
	exists := dao.Givestudent(claims.Name)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  401,
			"message": "用户不存在",
		})
		return
	}

	// 生成新的access token
	newAccessToken, err := utils.GenerateAccessToken(claims.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  500,
			"message": "生成新访问令牌失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "令牌刷新成功",
		"data": gin.H{
			"accesstoken": newAccessToken,
			"expiresin":   900,
		},
	})
}
