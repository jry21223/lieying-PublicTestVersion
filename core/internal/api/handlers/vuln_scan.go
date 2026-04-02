package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/kunlun-sec/lunying/internal/models"
)

type VulnScanHandler struct {
	db *gorm.DB
}

func NewVulnScanHandler(db *gorm.DB) *VulnScanHandler {
	return &VulnScanHandler{
		db: db,
	}
}

type CreateVulnScanRequest struct {
	Assets     []string `json:"assets"`
	Severity   []string `json:"severity"`
	Tags       []string `json:"tags"`
	CustomPOCs []string `json:"custom_pocs"`
}

func (h *VulnScanHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateVulnScanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "无效的请求体", http.StatusBadRequest)
		return
	}

	if len(req.Assets) == 0 {
		http.Error(w, "请至少选择一个资产", http.StatusBadRequest)
		return
	}

	target := req.Assets[0]
	optionsJSON, _ := json.Marshal(map[string]interface{}{
		"assets":   req.Assets,
		"severity": req.Severity,
		"tags":     req.Tags,
	})

	task := models.ActiveScanTask{
		ID:       uuid.New().String(),
		Name:     "漏洞扫描 - " + target,
		Target:   target,
		Type:     "vuln-scan",
		Status:   "pending",
		Progress: 0,
		Options:  string(optionsJSON),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := h.db.Create(&task).Error; err != nil {
		http.Error(w, "创建任务失败", http.StatusInternalServerError)
		return
	}

	ctx := r.Context()
	go func() {
		h.executeTask(ctx, task.ID, req, target)
	}()

	respondJSON(w, http.StatusCreated, task)
}

func (h *VulnScanHandler) executeTask(ctx context.Context, taskID string, req CreateVulnScanRequest, target string) {
	var task models.ActiveScanTask
	if err := h.db.First(&task, "id = ?", taskID).Error; err != nil {
		return
	}

	task.Status = "running"
	task.StartedAt = ptr(time.Now())
	task.UpdatedAt = time.Now()
	h.db.Save(&task)

	defer func() {
		task.Status = "completed"
		task.FinishedAt = ptr(time.Now())
		task.Progress = 100
		task.UpdatedAt = time.Now()
		h.db.Save(&task)
	}()

	// 模拟扫描进度
	stages := []int{25, 50, 75, 100}
	for _, progress := range stages {
		select {
		case <-ctx.Done():
			task.Status = "cancelled"
			return
		case <-time.After(1 * time.Second):
		}

		task.Progress = progress
		task.UpdatedAt = time.Now()
		h.db.Save(&task)
	}
}

func (h *VulnScanHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "taskId")

	var task models.ActiveScanTask
	if err := h.db.First(&task, "id = ?", taskID).Error; err != nil {
		http.Error(w, "任务不存在", http.StatusNotFound)
		return
	}

	respondJSON(w, http.StatusOK, task)
}

func (h *VulnScanHandler) ListTasks(w http.ResponseWriter, r *http.Request) {
	var tasks []models.ActiveScanTask
	if err := h.db.Where("type = ?", "vuln-scan").Order("created_at DESC").Find(&tasks).Error; err != nil {
		http.Error(w, "获取任务列表失败", http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, tasks)
}
