package executor

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"sync"
	"syscall"
	"time"

	"github.com/google/uuid"
)

type ToolExecutor struct {
	ID         string
	ToolPath   string
	Args       []string
	Env        []string
	WorkingDir string
	Timeout    time.Duration
	Stdout     bytes.Buffer
	Stderr     bytes.Buffer
	StartedAt  time.Time
	FinishedAt time.Time
	ExitCode   int
	CancelFunc context.CancelFunc
}

type ExecutionResult struct {
	ID         string
	Command    string
	Args       []string
	Stdout     string
	Stderr     string
	ExitCode   int
	Duration   time.Duration
	Success    bool
	Error      error
}

func NewToolExecutor(toolPath string, args []string) *ToolExecutor {
	return &ToolExecutor{
		ID:         uuid.New().String(),
		ToolPath:   toolPath,
		Args:       args,
		Env:        []string{},
		Timeout:    300 * time.Second,
	}
}

func (e *ToolExecutor) SetWorkingDir(dir string) {
	e.WorkingDir = dir
}

func (e *ToolExecutor) SetEnv(env []string) {
	e.Env = env
}

func (e *ToolExecutor) SetTimeout(timeout time.Duration) {
	e.Timeout = timeout
}

func (e *ToolExecutor) AddEnv(key, value string) {
	e.Env = append(e.Env, fmt.Sprintf("%s=%s", key, value))
}

func (e *ToolExecutor) SetProxy(proxyConfig *ProxyConfig) {
	if proxyConfig != nil && proxyConfig.Enabled {
		proxyURL := fmt.Sprintf("%s://%s:%d", proxyConfig.Type, proxyConfig.Host, proxyConfig.Port)
		e.AddEnv("HTTP_PROXY", proxyURL)
		e.AddEnv("HTTPS_PROXY", proxyURL)
		e.AddEnv("http_proxy", proxyURL)
		e.AddEnv("https_proxy", proxyURL)

		if proxyConfig.Username != "" && proxyConfig.Password != "" {
			e.AddEnv("PROXY_USER", proxyConfig.Username)
			e.AddEnv("PROXY_PASS", proxyConfig.Password)
		}
	}
}

func (e *ToolExecutor) Run(ctx context.Context) (*ExecutionResult, error) {
	ctx, cancel := context.WithTimeout(ctx, e.Timeout)
	defer cancel()
	e.CancelFunc = cancel

	e.StartedAt = time.Now()
	e.Stdout.Reset()
	e.Stderr.Reset()

	cmd := exec.CommandContext(ctx, e.ToolPath, e.Args...)
	cmd.Stdout = &e.Stdout
	cmd.Stderr = &e.Stderr

	if e.WorkingDir != "" {
		cmd.Dir = e.WorkingDir
	}

	if len(e.Env) > 0 {
		cmd.Env = append(os.Environ(), e.Env...)
	}

	err := cmd.Start()
	if err != nil {
		return e.createResult(err), err
	}

	err = cmd.Wait()
	e.FinishedAt = time.Now()
	e.ExitCode = 0

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
				e.ExitCode = status.ExitStatus()
			}
		}
	}

	return e.createResult(err), nil
}

func (e *ToolExecutor) RunAsync(ctx context.Context, callback func(*ExecutionResult)) {
	go func() {
		result, _ := e.Run(ctx)
		if callback != nil {
			callback(result)
		}
	}()
}

func (e *ToolExecutor) createResult(err error) *ExecutionResult {
	duration := time.Duration(0)
	if !e.FinishedAt.IsZero() {
		duration = e.FinishedAt.Sub(e.StartedAt)
	}

	success := err == nil && e.ExitCode == 0

	return &ExecutionResult{
		ID:         e.ID,
		Command:    e.ToolPath,
		Args:       e.Args,
		Stdout:     e.Stdout.String(),
		Stderr:     e.Stderr.String(),
		ExitCode:   e.ExitCode,
		Duration:   duration,
		Success:    success,
		Error:      err,
	}
}

func (e *ToolExecutor) Stop() error {
	if e.CancelFunc != nil {
		e.CancelFunc()
		return nil
	}
	return fmt.Errorf("executor not running")
}

func (e *ToolExecutor) GetStdout() string {
	return e.Stdout.String()
}

func (e *ToolExecutor) GetStderr() string {
	return e.Stderr.String()
}

func (e *ToolExecutor) IsRunning() bool {
	return !e.StartedAt.IsZero() && e.FinishedAt.IsZero()
}

func (e *ToolExecutor) GetDuration() time.Duration {
	if e.StartedAt.IsZero() {
		return 0
	}
	if e.FinishedAt.IsZero() {
		return time.Since(e.StartedAt)
	}
	return e.FinishedAt.Sub(e.StartedAt)
}

type ProxyConfig struct {
	Enabled  bool   `json:"enabled"`
	Type     string `json:"type"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

func CheckToolExists(toolPath string) bool {
	_, err := os.Stat(toolPath)
	if err == nil {
		return true
	}

	_, err = exec.LookPath(toolPath)
	return err == nil
}

func FindTool(name string) (string, error) {
	path, err := exec.LookPath(name)
	if err == nil {
		return path, nil
	}

	commonPaths := []string{
		"./bin/" + name,
		"./tools/" + name,
		"/usr/local/bin/" + name,
		"/usr/bin/" + name,
		"/opt/" + name + "/" + name,
	}

	for _, p := range commonPaths {
		if CheckToolExists(p) {
			return p, nil
		}
	}

	return "", fmt.Errorf("tool %s not found", name)
}

func ExecuteCommand(name string, args []string, timeout time.Duration, env []string) (*ExecutionResult, error) {
	toolPath, err := FindTool(name)
	if err != nil {
		return nil, err
	}

	executor := NewToolExecutor(toolPath, args)
	if timeout > 0 {
		executor.SetTimeout(timeout)
	}
	if env != nil {
		executor.SetEnv(env)
	}

	return executor.Run(context.Background())
}

type ToolRegistry struct {
	tools map[string]string
	mu    sync.RWMutex
}

func NewToolRegistry() *ToolRegistry {
	return &ToolRegistry{
		tools: make(map[string]string),
	}
}

func (r *ToolRegistry) Register(name, path string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tools[name] = path
}

func (r *ToolRegistry) Get(name string) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if path, ok := r.tools[name]; ok {
		return path, nil
	}

	path, err := FindTool(name)
	if err != nil {
		return "", err
	}

	r.mu.Lock()
	r.tools[name] = path
	r.mu.Unlock()

	return path, nil
}

func (r *ToolRegistry) List() map[string]string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make(map[string]string)
	for k, v := range r.tools {
		result[k] = v
	}
	return result
}

var DefaultToolRegistry = NewToolRegistry()

func init() {
	DefaultToolRegistry.Register("nuclei", "nuclei")
	DefaultToolRegistry.Register("amass", "amass")
	DefaultToolRegistry.Register("subfinder", "subfinder")
	DefaultToolRegistry.Register("nmap", "nmap")
	DefaultToolRegistry.Register("ffuf", "ffuf")
	DefaultToolRegistry.Register("masscan", "masscan")
}
