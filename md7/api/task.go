package api

import (
	"net/http"
	"strconv"

	"md7/model"
	"md7/sql"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TaskHandler struct {
	db *gorm.DB
}

func NewTaskHandler() *TaskHandler {
	db, _ := sql.InitDatabase()
	return &TaskHandler{db: db.DB}
}

type CreateTaskRequest struct {
	Title string `json:"title" binding:"required"`
}
type UpdateTaskRequest struct {
	Title  *string `json:"title"`
	Status *bool   `json:"status"`
}
func (h *TaskHandler) CreateTask(c *gin.Context) {
	var req CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  400,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}
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

	task := model.Task{
		Title:     req.Title,
		Status:    false,
		StudentID: student.ID,
	}

	if err := h.db.Create(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  500,
			"message": "创建任务失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  201,
		"message": "创建任务成功",
		"data":    task,
	})
}
func (h *TaskHandler) GetTasks(c *gin.Context) {
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

	var tasks []model.Task
	if err := h.db.Where("student_id = ?", student.ID).Order("created_at desc").Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  500,
			"message": "获取任务列表失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "获取任务列表成功",
		"data": gin.H{
			"total": len(tasks),
			"tasks": tasks,
		},
	})
}
func (h *TaskHandler) GetTask(c *gin.Context) {
	taskID := c.Param("id")
	if taskID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  400,
			"message": "任务ID不能为空",
		})
		return
	}

	id, err := strconv.ParseUint(taskID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  400,
			"message": "无效的任务ID",
		})
		return
	}

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
	var task model.Task
	if err := h.db.Where("id = ? AND student_id = ?", id, student.ID).First(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  404,
			"message": "任务不存在或无权限访问",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "获取任务成功",
		"data":    task,
	})
}
func (h *TaskHandler) UpdateTask(c *gin.Context) {
	taskID := c.Param("id")
	if taskID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  400,
			"message": "任务ID不能为空",
		})
		return
	}

	id, err := strconv.ParseUint(taskID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  400,
			"message": "无效的任务ID",
		})
		return
	}

	var req UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  400,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

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
	var task model.Task
	if err := h.db.Where("id = ? AND student_id = ?", id, student.ID).First(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  404,
			"message": "任务不存在或无权限访问",
		})
		return
	}
	updates := make(map[string]interface{})
	if req.Title != nil {
		updates["title"] = *req.Title
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}

	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  400,
			"message": "没有要更新的字段",
		})
		return
	}
	if err := h.db.Model(&task).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  500,
			"message": "更新任务失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "更新任务成功",
		"data":    task,
	})
}
func (h *TaskHandler) DeleteTask(c *gin.Context) {
	taskID := c.Param("id")
	if taskID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  400,
			"message": "任务ID不能为空",
		})
		return
	}

	id, err := strconv.ParseUint(taskID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  400,
			"message": "无效的任务ID",
		})
		return
	}

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
	result := h.db.Where("id = ? AND student_id = ?", id, student.ID).Delete(&model.Task{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  500,
			"message": "删除任务失败",
			"error":   result.Error.Error(),
		})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  404,
			"message": "任务不存在或无权限删除",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "删除任务成功",
	})
}
