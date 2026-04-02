package logger

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	// 测试创建logger
	tempDir := t.TempDir()
	logPath := filepath.Join(tempDir, "logs")

	logger, err := New("info", logPath, "text")
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Close()

	// 验证日志目录是否创建
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		t.Error("Log directory was not created")
	}
}

func TestLogLevels(t *testing.T) {
	tempDir := t.TempDir()
	logPath := filepath.Join(tempDir, "logs")

	logger, err := New("info", logPath, "text")
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Close()

	// 测试各级别日志
	logger.Debug("Debug message") // 应该被过滤
	logger.Info("Info message")
	logger.Warn("Warn message")
	logger.Error("Error message")
}

func TestParseLevel(t *testing.T) {
	tests := []struct {
		input    string
		expected Level
	}{
		{"debug", DEBUG},
		{"info", INFO},
		{"warn", WARN},
		{"error", ERROR},
		{"fatal", FATAL},
		{"unknown", INFO}, // 默认值
		{"", INFO},        // 空字符串默认值
	}

	for _, test := range tests {
		result := parseLevel(test.input)
		if result != test.expected {
			t.Errorf("parseLevel(%q) = %v, expected %v", test.input, result, test.expected)
		}
	}
}

func TestLevelString(t *testing.T) {
	tests := []struct {
		level    Level
		expected string
	}{
		{DEBUG, "DEBUG"},
		{INFO, "INFO"},
		{WARN, "WARN"},
		{ERROR, "ERROR"},
		{FATAL, "FATAL"},
		{Level(99), "UNKNOWN"},
	}

	for _, test := range tests {
		result := test.level.String()
		if result != test.expected {
			t.Errorf("Level(%d).String() = %q, expected %q", test.level, result, test.expected)
		}
	}
}

func TestLogFileCreation(t *testing.T) {
	tempDir := t.TempDir()
	logPath := filepath.Join(tempDir, "logs")

	logger, err := New("info", logPath, "text")
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	// 写入日志
	logger.Info("Test message")
	logger.Close()

	// 检查日志文件是否创建
	entries, err := os.ReadDir(logPath)
	if err != nil {
		t.Fatalf("Failed to read log directory: %v", err)
	}

	if len(entries) == 0 {
		t.Error("No log file was created")
	}

	// 验证文件内容
	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), ".log") {
			content, err := os.ReadFile(filepath.Join(logPath, entry.Name()))
			if err != nil {
				t.Fatalf("Failed to read log file: %v", err)
			}
			if !strings.Contains(string(content), "Test message") {
				t.Error("Log file does not contain expected message")
			}
		}
	}
}

func TestJSONFormat(t *testing.T) {
	tempDir := t.TempDir()
	logPath := filepath.Join(tempDir, "logs")

	logger, err := New("info", logPath, "json")
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Close()

	logger.Info("JSON test message")

	// 读取并验证JSON格式
	entries, _ := os.ReadDir(logPath)
	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), ".log") {
			content, _ := os.ReadFile(filepath.Join(logPath, entry.Name()))
			if !strings.Contains(string(content), `"level":"INFO"`) {
				t.Error("JSON log does not contain expected format")
			}
		}
	}
}

func BenchmarkLogInfo(b *testing.B) {
	logger, _ := New("info", "", "text")
	defer logger.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("Benchmark message %d", i)
	}
}
