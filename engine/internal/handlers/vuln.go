package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type VulnHandler struct {
	db *gorm.DB
}

func NewVulnHandler(db *gorm.DB) *VulnHandler {
	return &VulnHandler{db: db}
}

func (h *VulnHandler) List(c *gin.Context) {
	var vulns []interface{}
	if err := h.db.Raw("SELECT * FROM vulnerabilities LIMIT 100").Scan(&vulns).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, vulns)
}

func (h *VulnHandler) Get(c *gin.Context) {
	id := c.Param("id")
	var vuln interface{}
	if err := h.db.Raw("SELECT * FROM vulnerabilities WHERE id = ?", id).Scan(&vuln).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Vulnerability not found"})
		return
	}
	c.JSON(http.StatusOK, vuln)
}
