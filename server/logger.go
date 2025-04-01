package server

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

// Logger handles logging operations
type Logger struct {
	logFile   *os.File
	logger    *log.Logger
	logPath   string
	debugMode bool
}

// NewLogger creates a new logger instance
func NewLogger(logPath string, debugMode bool) (*Logger, error) {
	// Ensure log directory exists
	if err := os.MkdirAll(logPath, 0755); err != nil {
		return nil, err
	}

	// Create log file with timestamp in name
	timestamp := time.Now().Format("2006-01-02")
	logFileName := fmt.Sprintf("bot-log-%s.txt", timestamp)
	logFilePath := filepath.Join(logPath, logFileName)

	// Open log file
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}

	// Create logger
	logger := log.New(logFile, "", log.LstdFlags)

	return &Logger{
		logFile:   logFile,
		logger:    logger,
		logPath:   logPath,
		debugMode: debugMode,
	}, nil
}

// LogInfo logs an informational message
func (l *Logger) LogInfo(format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	l.logger.Printf("[INFO] %s", message)

	// Also print to console
	log.Printf("[INFO] %s", message)
}

// LogError logs an error message
func (l *Logger) LogError(format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	l.logger.Printf("[ERROR] %s", message)

	// Also print to console
	log.Printf("[ERROR] %s", message)
}

// LogDebug logs a debug message (only if debug mode is enabled)
func (l *Logger) LogDebug(format string, v ...interface{}) {
	if !l.debugMode {
		return
	}

	message := fmt.Sprintf(format, v...)
	l.logger.Printf("[DEBUG] %s", message)

	// Also print to console
	log.Printf("[DEBUG] %s", message)
}

// Close closes the log file
func (l *Logger) Close() error {
	return l.logFile.Close()
}
