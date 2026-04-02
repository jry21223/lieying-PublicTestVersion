package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/kunlun-sec/lunying/internal/models"
)

type ProxyHandler struct {
	db *gorm.DB
}

type ProxyPoolItem struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	Type      string    `json:"type"`
	Host      string    `json:"host"`
	Port      int       `json:"port"`
	Username  string    `json:"username,omitempty"`
	Password  string    `json:"password,omitempty"`
	Enabled   bool      `json:"enabled"`
	CreatedAt time.Time `json:"created_at"`
}

type UpdateProxyRequest struct {
	Type     string `json:"type"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Enabled  bool   `json:"enabled"`
}

func NewProxyHandler(db *gorm.DB) *ProxyHandler {
	return &ProxyHandler{db: db}
}

func (h *ProxyHandler) Get(w http.ResponseWriter, r *http.Request) {
	var configs []models.ProxyConfig
	if err := h.db.Find(&configs).Error; err != nil {
		http.Error(w, "获取代理配置失败", http.StatusInternalServerError)
		return
	}

	if len(configs) == 0 {
		defaultConfig := models.ProxyConfig{
			ID:      uuid.New().String(),
			Enabled: false,
			Type:    "http",
			Host:    "127.0.0.1",
			Port:    8080,
		}
		h.db.Create(&defaultConfig)
		respondJSON(w, http.StatusOK, defaultConfig)
		return
	}

	respondJSON(w, http.StatusOK, configs[0])
}

func (h *ProxyHandler) Update(w http.ResponseWriter, r *http.Request) {
	var req UpdateProxyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "无效的请求体", http.StatusBadRequest)
		return
	}

	var config models.ProxyConfig
	if h.db.First(&config).Error != nil {
		config.ID = uuid.New().String()
	}

	config.Enabled = req.Enabled
	config.Type = req.Type
	config.Host = req.Host
	config.Port = req.Port
	config.Username = req.Username
	config.Password = req.Password

	if config.ID == "" {
		if err := h.db.Create(&config).Error; err != nil {
			http.Error(w, "创建配置失败", http.StatusInternalServerError)
			return
		}
	} else {
		if err := h.db.Save(&config).Error; err != nil {
			http.Error(w, "更新配置失败", http.StatusInternalServerError)
			return
		}
	}

	respondJSON(w, http.StatusOK, config)
}

func (h *ProxyHandler) Test(w http.ResponseWriter, r *http.Request) {
	var testReq struct {
		URL string `json:"url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&testReq); err != nil {
		http.Error(w, "无效的请求体", http.StatusBadRequest)
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "连接测试功能",
		"url":     testReq.URL,
	})
}

func (h *ProxyHandler) GetPool(w http.ResponseWriter, r *http.Request) {
	var proxies []ProxyPoolItem
	if err := h.db.Find(&proxies).Error; err != nil {
		http.Error(w, "获取代理池失败", http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, proxies)
}

func (h *ProxyHandler) AddToPool(w http.ResponseWriter, r *http.Request) {
	var proxy ProxyPoolItem
	if err := json.NewDecoder(r.Body).Decode(&proxy); err != nil {
		http.Error(w, "无效的请求体", http.StatusBadRequest)
		return
	}

	proxy.ID = uuid.New().String()
	proxy.Enabled = true
	proxy.CreatedAt = time.Now()

	if err := h.db.Create(&proxy).Error; err != nil {
		http.Error(w, "添加代理失败", http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusCreated, proxy)
}

func (h *ProxyHandler) DeleteFromPool(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.db.Delete(&ProxyPoolItem{}, "id = ?", id).Error; err != nil {
		http.Error(w, "删除代理失败", http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "代理已删除"})
}

func (h *ProxyHandler) ImportPool(w http.ResponseWriter, r *http.Request) {
	var proxies []ProxyPoolItem
	if err := json.NewDecoder(r.Body).Decode(&proxies); err != nil {
		http.Error(w, "无效的请求体", http.StatusBadRequest)
		return
	}

	for _, proxy := range proxies {
		proxy.ID = uuid.New().String()
		proxy.Enabled = true
		proxy.CreatedAt = time.Now()
		h.db.Create(&proxy)
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "代理已导入"})
}
