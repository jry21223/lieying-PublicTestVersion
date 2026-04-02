package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lieying/engine/internal/models"
	"gorm.io/gorm"
)

type POCHandler struct {
	db *gorm.DB
}

func NewPOCHandler(db *gorm.DB) *POCHandler {
	return &POCHandler{db: db}
}

func (h *POCHandler) List(c *gin.Context) {
	var pocs []models.POC
	if err := h.db.Find(&pocs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, pocs)
}

func (h *POCHandler) Create(c *gin.Context) {
	var poc models.POC
	if err := c.ShouldBindJSON(&poc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.db.Create(&poc).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, poc)
}

func (h *POCHandler) Get(c *gin.Context) {
	id := c.Param("id")
	var poc models.POC
	if err := h.db.First(&poc, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "POC not found"})
		return
	}
	c.JSON(http.StatusOK, poc)
}
