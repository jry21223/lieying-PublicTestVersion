package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lieying/engine/internal/models"
	"github.com/lieying/engine/internal/scheduler"
	"gorm.io/gorm"
)

type TaskHandler struct {
	db        *gorm.DB
	scheduler *scheduler.TaskScheduler
}

func NewTaskHandler(db *gorm.DB, scheduler *scheduler.TaskScheduler) *TaskHandler {
	return &TaskHandler{db: db, scheduler: scheduler}
}

func (h *TaskHandler) Create(c *gin.Context) {
	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.Create(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, task)
}

func (h *TaskHandler) List(c *gin.Context) {
	var tasks []models.Task
	if err := h.db.Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

func (h *TaskHandler) Get(c *gin.Context) {
	id := c.Param("id")
	var task models.Task
	if err := h.db.First(&task, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) Start(c *gin.Context) {
	id := c.Param("id")
	taskID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	var task models.Task
	if err := h.db.First(&task, "id = ?", taskID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	if err := h.scheduler.Submit(&task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task started"})
}

func (h *TaskHandler) Stop(c *gin.Context) {
	id := c.Param("id")
	taskID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	var task models.Task
	if err := h.db.First(&task, "id = ?", taskID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	task.Status = models.TaskStatusPaused
	h.db.Save(&task)

	c.JSON(http.StatusOK, gin.H{"message": "Task stopped"})
}
