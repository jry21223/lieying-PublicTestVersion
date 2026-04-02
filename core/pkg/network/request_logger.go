package network

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RequestLog struct {
	ID              string    `json:"id" gorm:"primaryKey"`
	Module          string    `json:"module"`
	Target          string    `json:"target"`
	Method          string    `json:"method"`
	URL             string    `json:"url"`
	RequestHeaders  string    `json:"request_headers"`
	RequestBody    string    `json:"request_body"`
	ResponseStatus  int       `json:"response_status"`
	ResponseHeaders string    `json:"response_headers"`
	ResponseBody   string    `json:"response_body"`
	Duration        int64     `json:"duration"`
	Error           string    `json:"error,omitempty"`
	Timestamp       time.Time `json:"timestamp"`
}

type RequestLogger struct {
	db        *gorm.DB
	logs      []RequestLog
	mu        sync.RWMutex
	maxMemory int
}

func NewRequestLogger(db *gorm.DB, maxMemory int) *RequestLogger {
	return &RequestLogger{
		db:        db,
		logs:      make([]RequestLog, 0, maxMemory),
		maxMemory: maxMemory,
	}
}

func (r *RequestLogger) Log(req *http.Request, resp *http.Response, duration int64, err error) {
	log := RequestLog{
		ID:        uuid.New().String(),
		Module:    getModuleFromRequest(req),
		Target:    extractTarget(req.URL.String()),
		Method:    req.Method,
		URL:       req.URL.String(),
		Duration:  duration,
		Timestamp: time.Now(),
	}

	if err != nil {
		log.Error = err.Error()
	}

	if req != nil {
		log.RequestHeaders = serializeHeaders(req.Header)
		if req.Body != nil {
			body, _ := io.ReadAll(req.Body)
			log.RequestBody = string(body)
			req.Body = io.NopCloser(bytes.NewBuffer(body))
		}
	}

	if resp != nil {
		log.ResponseStatus = resp.StatusCode
		log.ResponseHeaders = serializeHeaders(resp.Header)
		if resp.Body != nil {
			body, _ := io.ReadAll(resp.Body)
			log.ResponseBody = truncateString(string(body), 10000)
			resp.Body = io.NopCloser(bytes.NewBuffer(body))
		}
	}

	r.mu.Lock()
	r.logs = append(r.logs, log)
	if len(r.logs) > r.maxMemory {
		r.flush()
	}
	r.mu.Unlock()
}

func (r *RequestLogger) LogRequest(module, method, url, reqBody, respBody string, status int, duration int64) {
	log := RequestLog{
		ID:             uuid.New().String(),
		Module:         module,
		Method:         method,
		URL:            url,
		RequestBody:    truncateString(reqBody, 5000),
		ResponseBody:   truncateString(respBody, 10000),
		ResponseStatus: status,
		Duration:       duration,
		Timestamp:      time.Now(),
	}

	r.mu.Lock()
	r.logs = append(r.logs, log)
	if len(r.logs) > r.maxMemory {
		r.flush()
	}
	r.mu.Unlock()
}

func (r *RequestLogger) flush() {
	if len(r.logs) == 0 {
		return
	}

	if r.db != nil {
		r.db.Create(&r.logs)
	}
	r.logs = make([]RequestLog, 0, r.maxMemory)
}

func (r *RequestLogger) GetLogs(module, target string, limit, offset int) ([]RequestLog, int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	query := r.db.Model(&RequestLog{})

	if module != "" {
		query = query.Where("module = ?", module)
	}
	if target != "" {
		query = query.Where("target LIKE ?", "%"+target+"%")
	}

	var total int64
	query.Count(&total)

	var logs []RequestLog
	err := query.Order("timestamp DESC").Limit(limit).Offset(offset).Find(&logs).Error

	return logs, total, err
}

func (r *RequestLogger) GetLogByID(id string) (*RequestLog, error) {
	var log RequestLog
	err := r.db.Where("id = ?", id).First(&log).Error
	if err != nil {
		return nil, err
	}
	return &log, nil
}

func (r *RequestLogger) ExportLogs(module, target string, format string) ([]byte, error) {
	logs, _, err := r.GetLogs(module, target, 10000, 0)
	if err != nil {
		return nil, err
	}

	switch format {
	case "json":
		return json.MarshalIndent(logs, "", "  ")
	case "csv":
		return exportLogsToCSV(logs)
	default:
		return json.Marshal(logs)
	}
}

func (r *RequestLogger) ClearLogs() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.db.Delete(&RequestLog{}).Error
}

func (r *RequestLogger) Flush() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.flush()
}

func getModuleFromRequest(req *http.Request) string {
	if req == nil || req.URL == nil {
		return "unknown"
	}
	path := req.URL.Path
	if strings.HasPrefix(path, "/api/recon") {
		return "recon"
	}
	if strings.HasPrefix(path, "/api/scan") {
		return "scan"
	}
	if strings.HasPrefix(path, "/api/active") {
		return "active-scan"
	}
	if strings.HasPrefix(path, "/api/ai") {
		return "ai"
	}
	if strings.HasPrefix(path, "/api/edu") {
		return "edu"
	}
	return "general"
}

func extractTarget(urlStr string) string {
	if urlStr == "" {
		return ""
	}
	if idx := strings.Index(urlStr, "?"); idx != -1 {
		urlStr = urlStr[:idx]
	}
	return strings.TrimPrefix(urlStr, "http://")
}

func serializeHeaders(headers http.Header) string {
	var parts []string
	for key, values := range headers {
		for _, value := range values {
			parts = append(parts, fmt.Sprintf("%s: %s", key, value))
		}
	}
	return strings.Join(parts, "\n")
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

func exportLogsToCSV(logs []RequestLog) ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteString("ID,Module,Target,Method,URL,Status,Duration(ms),Timestamp\n")

	for _, log := range logs {
		buf.WriteString(fmt.Sprintf("%s,%s,%s,%s,%s,%d,%d,%s\n",
			log.ID,
			log.Module,
			log.Target,
			log.Method,
			log.URL,
			log.ResponseStatus,
			log.Duration,
			log.Timestamp.Format(time.RFC3339),
		))
	}

	return buf.Bytes(), nil
}
