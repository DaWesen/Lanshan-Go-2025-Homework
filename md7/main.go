package main

import (
	"log"
	"md7/api"
	"md7/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.Static("/static", "./frontend")
	r.StaticFile("/", "./frontend/index.html")
	r.StaticFile("/index.html", "./frontend/index.html")
	r.NoRoute(func(c *gin.Context) {
		if c.Request.Method == "GET" {
			c.File("./frontend/index.html")
		} else {
			c.JSON(404, gin.H{
				"status":  404,
				"message": "API endpoint not found",
			})
		}
	})
	authHandler := api.NewAuthHandler()
	taskHandler := api.NewTaskHandler()
	public := r.Group("/api")
	{
		public.POST("/register", authHandler.Register)
		public.POST("/login", authHandler.Login)
		public.POST("/refresh", authHandler.RefreshToken)
	}
	protected := r.Group("/api")
	protected.Use(utils.JWTAuthMiddleware())
	{
		protected.GET("/profile", authHandler.GetProfile)
		protected.POST("/tasks", taskHandler.CreateTask)
		protected.GET("/tasks", taskHandler.GetTasks)
		protected.GET("/tasks/:id", taskHandler.GetTask)
		protected.PUT("/tasks/:id", taskHandler.UpdateTask)
		protected.DELETE("/tasks/:id", taskHandler.DeleteTask)
	}
	log.Println("✅ 服务器启动在 :8080")
	log.Println("🌐 前端地址: http://localhost:8080")
	log.Println("🔄 健康检查: http://localhost:8080/health")
	log.Println("📝 API文档:")
	log.Println("  注册: POST http://localhost:8080/api/register")
	log.Println("  登录: POST http://localhost:8080/api/login")
	log.Println("  获取任务: GET http://localhost:8080/api/tasks")

	if err := r.Run(":8080"); err != nil {
		log.Fatal("启动服务器失败:", err)
	}
}
