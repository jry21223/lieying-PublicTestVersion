package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/lieying/engine/internal/config"
)

type AIClient interface {
	Generate(prompt string) (string, error)
	IsAvailable() bool
}

type OllamaClient struct {
	endpoint string
	model    string
	client   *http.Client
}

func NewOllamaClient(cfg config.LocalAI, httpClient *http.Client) *OllamaClient {
	return &OllamaClient{
		endpoint: cfg.Endpoint,
		model:    cfg.Model,
		client:   httpClient,
	}
}

func (c *OllamaClient) Generate(prompt string) (string, error) {
	reqBody := map[string]interface{}{
		"model":  c.model,
		"prompt": prompt,
		"stream": false,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	resp, err := c.client.Post(
		fmt.Sprintf("%s/api/generate", c.endpoint),
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result struct {
		Response string `json:"response"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	return result.Response, nil
}

func (c *OllamaClient) IsAvailable() bool {
	resp, err := c.client.Get(fmt.Sprintf("%s/api/tags", c.endpoint))
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

type AIClientManager struct {
	local  AIClient
	remote AIClient
	mode   config.AIMode
}

func NewAIClientManager(cfg config.AIConfig, httpClient *http.Client) *AIClientManager {
	mgr := &AIClientManager{
		mode: cfg.Mode,
	}

	if cfg.Mode == config.AIModeLocal || cfg.Mode == config.AIModeHybrid {
		mgr.local = NewOllamaClient(cfg.Local, httpClient)
	}

	if cfg.Mode == config.AIModeRemote || cfg.Mode == config.AIModeHybrid {
		mgr.remote = NewRemoteClient(cfg.Remote, httpClient)
	}

	return mgr
}

func (m *AIClientManager) Generate(prompt string) (string, error) {
	switch m.mode {
	case config.AIModeLocal:
		if m.local != nil && m.local.IsAvailable() {
			return m.local.Generate(prompt)
		}
		return "", fmt.Errorf("local AI not available")
	case config.AIModeRemote:
		if m.remote != nil {
			return m.remote.Generate(prompt)
		}
		return "", fmt.Errorf("remote AI not available")
	case config.AIModeHybrid:
		if m.local != nil && m.local.IsAvailable() {
			return m.local.Generate(prompt)
		}
		if m.remote != nil {
			return m.remote.Generate(prompt)
		}
		return "", fmt.Errorf("no AI client available")
	}
	return "", fmt.Errorf("invalid AI mode")
}

func (m *AIClientManager) IsAvailable() bool {
	switch m.mode {
	case config.AIModeLocal:
		return m.local != nil && m.local.IsAvailable()
	case config.AIModeRemote:
		return m.remote != nil
	case config.AIModeHybrid:
		return (m.local != nil && m.local.IsAvailable()) || m.remote != nil
	}
	return false
}

type RemoteClient struct {
	apiKey  string
	baseURL string
	model   string
	client  *http.Client
}

func NewRemoteClient(cfg config.RemoteAI, httpClient *http.Client) *RemoteClient {
	return &RemoteClient{
		apiKey:  cfg.APIKey,
		baseURL: cfg.BaseURL,
		model:   cfg.Model,
		client:  httpClient,
	}
}

func (c *RemoteClient) Generate(prompt string) (string, error) {
	reqBody := map[string]interface{}{
		"model": c.model,
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": prompt,
			},
		},
		"stream": false,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/chat/completions", c.baseURL),
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("no response from AI")
	}

	return result.Choices[0].Message.Content, nil
}

func (c *RemoteClient) IsAvailable() bool {
	return c.apiKey != ""
}
