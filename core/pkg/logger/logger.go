package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Level int

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
	FATAL
)

func (l Level) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

type Logger struct {
	level      Level
	output     io.Writer
	fileOutput *os.File
	format     string
	mu         sync.Mutex
}

var (
	defaultLogger *Logger
	once          sync.Once
)

func Init(level string, logPath string, format string) error {
	var err error
	once.Do(func() {
		defaultLogger, err = New(level, logPath, format)
	})
	return err
}

func New(level string, logPath string, format string) (*Logger, error) {
	l := &Logger{
		level:  parseLevel(level),
		format: format,
	}

	if logPath != "" {
		if err := os.MkdirAll(logPath, 0755); err != nil {
			return nil, err
		}

		logFile := filepath.Join(logPath, fmt.Sprintf("lieying_%s.log", time.Now().Format("2006-01-02")))
		file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return nil, err
		}
		l.fileOutput = file
		l.output = io.MultiWriter(os.Stdout, file)
	} else {
		l.output = os.Stdout
	}

	return l, nil
}

func parseLevel(level string) Level {
	switch level {
	case "debug":
		return DEBUG
	case "info":
		return INFO
	case "warn":
		return WARN
	case "error":
		return ERROR
	case "fatal":
		return FATAL
	default:
		return INFO
	}
}

func (l *Logger) log(level Level, format string, args ...interface{}) {
	if level < l.level {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	msg := fmt.Sprintf(format, args...)

	var logLine string
	if l.format == "json" {
		logLine = fmt.Sprintf(`{"time":"%s","level":"%s","msg":"%s"}`+"\n", timestamp, level.String(), msg)
	} else {
		logLine = fmt.Sprintf("[%s] %s: %s\n", timestamp, level.String(), msg)
	}

	fmt.Fprint(l.output, logLine)
}

func (l *Logger) Debug(format string, args ...interface{}) {
	l.log(DEBUG, format, args...)
}

func (l *Logger) Info(format string, args ...interface{}) {
	l.log(INFO, format, args...)
}

func (l *Logger) Warn(format string, args ...interface{}) {
	l.log(WARN, format, args...)
}

func (l *Logger) Error(format string, args ...interface{}) {
	l.log(ERROR, format, args...)
}

func (l *Logger) Fatal(format string, args ...interface{}) {
	l.log(FATAL, format, args...)
	os.Exit(1)
}

func (l *Logger) Close() error {
	if l.fileOutput != nil {
		return l.fileOutput.Close()
	}
	return nil
}

func Debug(format string, args ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Debug(format, args...)
	}
}

func Info(format string, args ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Info(format, args...)
	}
}

func Warn(format string, args ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Warn(format, args...)
	}
}

func Error(format string, args ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Error(format, args...)
	}
}

func Fatal(format string, args ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Fatal(format, args...)
	}
}
