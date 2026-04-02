package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type OllamaClient struct {
	baseURL    string
	model      string
	httpClient *http.Client
}

type OllamaRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OllamaResponse struct {
	Message Message `json:"message"`
	Done    bool    `json:"done"`
}

type OllamaModel struct {
	Name       string `json:"name"`
	Model      string `json:"model"`
	ModifiedAt string `json:"modified_at"`
	Size       int64  `json:"size"`
	Digest     string `json:"digest"`
	Details    struct {
		Format            string   `json:"format"`
		Family            string   `json:"family"`
		Families          []string `json:"families"`
		ParameterSize     string   `json:"parameter_size"`
		QuantizationLevel string   `json:"quantization_level"`
	} `json:"details"`
}

func NewOllamaClient(baseURL string, model string) *OllamaClient {
	if baseURL == "" {
		baseURL = "http://localhost:11434"
	}
	if model == "" {
		model = "qwen2.5:14b"
	}

	return &OllamaClient{
		baseURL: baseURL,
		model:   model,
		httpClient: &http.Client{
			Timeout: 120 * time.Second,
		},
	}
}

func (oc *OllamaClient) IsAvailable() bool {
	resp, err := oc.httpClient.Get(oc.baseURL + "/api/tags")
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == 200
}

func (oc *OllamaClient) ListModels() ([]OllamaModel, error) {
	resp, err := oc.httpClient.Get(oc.baseURL + "/api/tags")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result struct {
		Models []OllamaModel `json:"models"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return result.Models, nil
}

func (oc *OllamaClient) Chat(messages []Message) (string, error) {
	req := OllamaRequest{
		Model:    oc.model,
		Messages: messages,
		Stream:   false,
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	resp, err := oc.httpClient.Post(
		oc.baseURL+"/api/chat",
		"application/json",
		bytes.NewBuffer(reqBody),
	)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result OllamaResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	return result.Message.Content, nil
}

func (oc *OllamaClient) Generate(prompt string) (string, error) {
	messages := []Message{
		{
			Role:    "user",
			Content: prompt,
		},
	}
	return oc.Chat(messages)
}

func (oc *OllamaClient) PullModel(modelName string) error {
	fmt.Printf("📥 正在拉取模型: %s\n", modelName)

	reqBody, err := json.Marshal(map[string]string{
		"name": modelName,
	})
	if err != nil {
		return err
	}

	resp, err := oc.httpClient.Post(
		oc.baseURL+"/api/pull",
		"application/json",
		bytes.NewBuffer(reqBody),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("拉取模型失败: %s", string(body))
	}

	fmt.Printf("✅ 模型 %s 拉取完成\n", modelName)
	return nil
}
