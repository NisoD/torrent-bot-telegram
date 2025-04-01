package server

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

type Logger struct {
	logFile   *os.File
	logger    *log.Logger
	logPath   string
	debugMode bool
}

func NewLogger(logPath string, debugMode bool) (*Logger, error) {

	if err := os.MkdirAll(logPath, 0755); err != nil {
		return nil, err
	}

	timestamp := time.Now().Format("2006-01-02")
	logFileName := fmt.Sprintf("bot-log-%s.txt", timestamp)
	logFilePath := filepath.Join(logPath, logFileName)

	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	logger := log.New(logFile, "", log.LstdFlags)

	return &Logger{
		logFile:   logFile,
		logger:    logger,
		logPath:   logPath,
		debugMode: debugMode,
	}, nil
}

// Write
func (l *Logger) LogInfo(format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	l.logger.Printf("[INFO] %s", message)

	log.Printf("[INFO] %s", message)
}

func (l *Logger) LogError(format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	l.logger.Printf("[ERROR] %s", message)

	log.Printf("[ERROR] %s", message)
}

//  debug mode
func (l *Logger) LogDebug(format string, v ...interface{}) {
	if !l.debugMode {
		return
	}

	message := fmt.Sprintf(format, v...)
	l.logger.Printf("[DEBUG] %s", message)

	log.Printf("[DEBUG] %s", message)
}

func (l *Logger) Close() error {
	return l.logFile.Close()
}
