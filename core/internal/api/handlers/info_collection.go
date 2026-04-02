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

type InfoCollectionHandler struct {
	db *gorm.DB
}

func NewInfoCollectionHandler(db *gorm.DB) *InfoCollectionHandler {
	return &InfoCollectionHandler{
		db: db,
	}
}

type CreateInfoCollectionRequest struct {
	Target  string   `json:"target"`
	Modules []string `json:"modules"`
	Options struct {
		PortRange string `json:"port_range"`
		Threads   int    `json:"threads"`
		Timeout   int    `json:"timeout"`
	} `json:"options"`
}

type InfoCollectionTask struct {
	ID        string                 `json:"id"`
	Target    string                 `json:"target"`
	Modules   []string               `json:"modules"`
	Status    string                 `json:"status"`
	Progress  float64                `json:"progress"`
	Options   map[string]interface{} `json:"options"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}

func (h *InfoCollectionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateInfoCollectionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "无效的请求体", http.StatusBadRequest)
		return
	}

	if req.Target == "" {
		http.Error(w, "目标不能为空", http.StatusBadRequest)
		return
	}

	optionsJSON, _ := json.Marshal(map[string]interface{}{
		"modules": req.Modules,
		"options": req.Options,
	})

	task := models.ActiveScanTask{
		ID:       uuid.New().String(),
		Name:     "信息收集 - " + req.Target,
		Target:   req.Target,
		Type:     "info-collection",
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
		h.executeTask(ctx, task.ID, req)
	}()

	respondJSON(w, http.StatusCreated, task)
}

func (h *InfoCollectionHandler) executeTask(ctx context.Context, taskID string, req CreateInfoCollectionRequest) {
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

	for i, _ := range req.Modules {
		progress := float64(i+1) / float64(len(req.Modules)) * 100
		task.Progress = int(progress)
		task.UpdatedAt = time.Now()
		h.db.Save(&task)

		select {
		case <-ctx.Done():
			task.Status = "cancelled"
			return
		case <-time.After(500 * time.Millisecond):
		}
	}
}

func (h *InfoCollectionHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "taskId")

	var task models.ActiveScanTask
	if err := h.db.First(&task, "id = ?", taskID).Error; err != nil {
		http.Error(w, "任务不存在", http.StatusNotFound)
		return
	}

	respondJSON(w, http.StatusOK, task)
}

func (h *InfoCollectionHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "taskId")

	if err := h.db.Delete(&models.ActiveScanTask{}, "id = ?", taskID).Error; err != nil {
		http.Error(w, "删除任务失败", http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "任务已删除"})
}

func (h *InfoCollectionHandler) GetProgress(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "taskId")

	var task models.ActiveScanTask
	if err := h.db.First(&task, "id = ?", taskID).Error; err != nil {
		http.Error(w, "任务不存在", http.StatusNotFound)
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"status":   task.Status,
		"progress": task.Progress,
	})
}
