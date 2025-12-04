package api

import (
	"md6/utils"

	"github.com/gin-gonic/gin"
)

func InitRouter() {
	router := gin.Default()

	// 添加CORS中间件
	router.Use(corsMiddleware())

	// 静态文件服务 - 提供前端文件
	router.Static("/static", "./static")
	router.StaticFile("/", "./static/index.html")
	router.StaticFile("/index.html", "./static/index.html")
	router.POST("/register", register)
	router.POST("/login", login)
	router.POST("/verifyToken", verifyToken)
	router.POST("/refresh", refresh)

	auth := router.Group("/auth")
	auth.Use(utils.JWTAuthMiddleware())
	{
		auth.POST("/changePassword", changePassword)
		auth.GET("/student", getstudent)
		auth.GET("/students", getstudents)
	}
	router.Run(":8080")
}

// CORS中间件
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Content-Length, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}

		c.Next()
	}
}
