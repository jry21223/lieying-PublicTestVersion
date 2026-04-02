package handlers

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lieying/engine/pkg/ai"
)

type AIHandler struct {
	aiManager *ai.AIClientManager
}

func NewAIHandler(aiManager *ai.AIClientManager) *AIHandler {
	return &AIHandler{aiManager: aiManager}
}

func (h *AIHandler) GeneratePOC(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read request body"})
		return
	}
	defer c.Request.Body.Close()

	prompt := string(body)
	result, err := h.aiManager.Generate(prompt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": result})
}

func (h *AIHandler) GenerateReport(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read request body"})
		return
	}
	defer c.Request.Body.Close()

	prompt := string(body)
	result, err := h.aiManager.Generate(prompt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": result})
}

func (h *AIHandler) AnalyzeAttackPath(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read request body"})
		return
	}
	defer c.Request.Body.Close()

	prompt := string(body)
	result, err := h.aiManager.Generate(prompt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": result})
}

func (h *AIHandler) Status(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"available": h.aiManager.IsAvailable(),
	})
}
