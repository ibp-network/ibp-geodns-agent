package logging

import (
	"log"
	"os"
	"strings"
)

type LogLevel int

const (
	LevelDebug LogLevel = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

var (
	currentLevel LogLevel = LevelInfo
	logger       *log.Logger
)

// Init initializes the logger
func Init(logLevelOverride ...string) {
	logger = log.New(os.Stdout, "", log.LstdFlags|log.Lmicroseconds)
	
	if len(logLevelOverride) > 0 && logLevelOverride[0] != "" {
		SetLevel(logLevelOverride[0])
	}
}

// SetLevel sets the log level
func SetLevel(level string) {
	switch strings.ToLower(level) {
	case "debug":
		currentLevel = LevelDebug
	case "info":
		currentLevel = LevelInfo
	case "warn", "warning":
		currentLevel = LevelWarn
	case "error":
		currentLevel = LevelError
	case "fatal":
		currentLevel = LevelFatal
	default:
		currentLevel = LevelInfo
	}
}

// shouldLog returns true if the given level should be logged
func shouldLog(level LogLevel) bool {
	return level >= currentLevel
}

// Debug logs a debug message
func Debug(msg string, args ...interface{}) {
	if shouldLog(LevelDebug) {
		if len(args) > 0 {
			logger.Printf("[DEBUG] "+msg, args...)
		} else {
			logger.Printf("[DEBUG] " + msg)
		}
	}
}

// Info logs an info message
func Info(msg string, args ...interface{}) {
	if shouldLog(LevelInfo) {
		if len(args) > 0 {
			logger.Printf("[INFO] "+msg, args...)
		} else {
			logger.Printf("[INFO] " + msg)
		}
	}
}

// Warn logs a warning message
func Warn(msg string, args ...interface{}) {
	if shouldLog(LevelWarn) {
		if len(args) > 0 {
			logger.Printf("[WARN] "+msg, args...)
		} else {
			logger.Printf("[WARN] " + msg)
		}
	}
}

// Error logs an error message
func Error(msg string, args ...interface{}) {
	if shouldLog(LevelError) {
		if len(args) > 0 {
			logger.Printf("[ERROR] "+msg, args...)
		} else {
			logger.Printf("[ERROR] " + msg)
		}
	}
}

// Fatal logs a fatal message and exits
func Fatal(msg string, args ...interface{}) {
	if len(args) > 0 {
		logger.Printf("[FATAL] "+msg, args...)
	} else {
		logger.Printf("[FATAL] " + msg)
	}
	os.Exit(1)
}
