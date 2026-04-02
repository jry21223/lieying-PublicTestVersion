package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/kunlun-sec/lunying/internal/models"
	"github.com/kunlun-sec/lunying/pkg/network"
)

type ActiveScanHandler struct {
	db     *gorm.DB
	engine *network.NetworkEngine
}

func NewActiveScanHandler(db *gorm.DB, engine *network.NetworkEngine) *ActiveScanHandler {
	return &ActiveScanHandler{
		db:     db,
		engine: engine,
	}
}

type CreateActiveScanRequest struct {
	Name     string            `json:"name"`
	Target   string            `json:"target"`
	Type     string            `json:"type"`
	URL      string            `json:"url"`
	Method   string            `json:"method"`
	Headers  map[string]string `json:"headers"`
	Body     string            `json:"body"`
	Payloads map[string][]string `json:"payloads"`
	Options  map[string]interface{} `json:"options"`
}

type ActiveScanResult struct {
	ID        string            `json:"id"`
	Method    string            `json:"method"`
	URL       string            `json:"url"`
	Payload   string            `json:"payload,omitempty"`
	Status    int               `json:"status"`
	Headers   map[string]string `json:"headers"`
	Body      string            `json:"body"`
	Duration  int64             `json:"duration"`
	Error     string            `json:"error,omitempty"`
	Timestamp time.Time        `json:"timestamp"`
}

func (h *ActiveScanHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateActiveScanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "无效的请求体", http.StatusBadRequest)
		return
	}

	url := req.URL
	if url == "" {
		url = req.Target
	}

	if url == "" {
		http.Error(w, "URL不能为空", http.StatusBadRequest)
		return
	}

	method := req.Method
	if method == "" {
		method = "GET"
	}

	ctx := r.Context()
	results := []*ActiveScanResult{}

	if len(req.Payloads) > 0 {
		for param, values := range req.Payloads {
			for _, value := range values {
				modifiedURL := url
				modifiedBody := req.Body

				if param != "" {
					modifiedURL = replacePayload(url, param, value)
					modifiedBody = replacePayload(req.Body, param, value)
				}

				result := h.executeSingleRequest(ctx, method, modifiedURL, req.Headers, modifiedBody)
				result.Payload = value
				results = append(results, result)
			}
		}
	} else {
		result := h.executeSingleRequest(ctx, method, url, req.Headers, req.Body)
		results = append(results, result)
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"results": results,
	})
}

func (h *ActiveScanHandler) executeSingleRequest(ctx context.Context, method, url string, headers map[string]string, body string) *ActiveScanResult {
	startTime := time.Now()

	netReq := &network.Request{
		Method:  method,
		URL:     url,
		Headers: headers,
		Body:    body,
	}

	resp, err := h.engine.DoRequest(ctx, netReq)
	duration := time.Since(startTime).Milliseconds()

	result := &ActiveScanResult{
		ID:        uuid.New().String(),
		Method:    method,
		URL:       url,
		Timestamp: time.Now(),
	}

	if err != nil {
		result.Error = err.Error()
		return result
	}

	result.Status = resp.Status
	result.Headers = make(map[string]string)
	for k, v := range resp.Headers {
		result.Headers[k] = v
	}
	result.Body = resp.Body
	result.Duration = duration

	h.saveRequestToHistory(method, url, headers, body, result)
	return result
}

func (h *ActiveScanHandler) saveRequestToHistory(method, url string, headers map[string]string, body string, result *ActiveScanResult) {
	requestHeaders := ""
	for k, v := range headers {
		requestHeaders += k + ": " + v + "\n"
	}

	responseHeaders := ""
	for k, v := range result.Headers {
		responseHeaders += k + ": " + v + "\n"
	}

	reqLog := models.NetworkRequest{
		ID:              uuid.New().String(),
		Module:          "active-scan",
		Target:          extractTarget(url),
		Method:          method,
		URL:             url,
		RequestHeaders:  requestHeaders,
		RequestBody:     body,
		ResponseStatus:  result.Status,
		ResponseHeaders: responseHeaders,
		ResponseBody:    result.Body,
		Duration:        result.Duration,
		Timestamp:       time.Now(),
	}

	if result.Error != "" {
		reqLog.Error = result.Error
	}

	h.db.Create(&reqLog)
}

func replacePayload(str, param, value string) string {
	placeholder := "{{" + param + "}}"
	return strings.ReplaceAll(str, placeholder, value)
}

func extractTarget(urlStr string) string {
	if urlStr == "" {
		return ""
	}
	if idx := strings.Index(urlStr, "://"); idx != -1 {
		urlStr = urlStr[idx+3:]
	}
	if idx := strings.Index(urlStr, "/"); idx != -1 {
		urlStr = urlStr[:idx]
	}
	if idx := strings.Index(urlStr, "?"); idx != -1 {
		urlStr = urlStr[:idx]
	}
	return urlStr
}
