package logging

import (
	"fmt"
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

func ensureLogger() {
	if logger == nil {
		logger = log.New(os.Stdout, "", log.LstdFlags|log.Lmicroseconds)
	}
}

// Init initializes the logger
func Init(logLevelOverride ...string) {
	ensureLogger()
	logger = log.New(os.Stdout, "", log.LstdFlags|log.Lmicroseconds)

	if len(logLevelOverride) > 0 && logLevelOverride[0] != "" {
		SetLevel(logLevelOverride[0])
	}
}

func formatMessage(msg string, args ...interface{}) string {
	if len(args) == 0 {
		return msg
	}

	looksLikeKeyValueArgs := len(args)%2 == 0
	if looksLikeKeyValueArgs {
		for i := 0; i < len(args); i += 2 {
			if _, ok := args[i].(string); !ok {
				looksLikeKeyValueArgs = false
				break
			}
		}
	}

	var builder strings.Builder
	builder.WriteString(msg)
	if !looksLikeKeyValueArgs {
		for _, arg := range args {
			builder.WriteByte(' ')
			builder.WriteString(fmt.Sprint(arg))
		}
		return builder.String()
	}

	for i := 0; i < len(args); i += 2 {
		builder.WriteByte(' ')
		builder.WriteString(fmt.Sprint(args[i]))
		if i+1 < len(args) {
			builder.WriteByte('=')
			builder.WriteString(fmt.Sprint(args[i+1]))
		}
	}

	return builder.String()
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
		ensureLogger()
		logger.Print("[DEBUG] " + formatMessage(msg, args...))
	}
}

// Info logs an info message
func Info(msg string, args ...interface{}) {
	if shouldLog(LevelInfo) {
		ensureLogger()
		logger.Print("[INFO] " + formatMessage(msg, args...))
	}
}

// Warn logs a warning message
func Warn(msg string, args ...interface{}) {
	if shouldLog(LevelWarn) {
		ensureLogger()
		logger.Print("[WARN] " + formatMessage(msg, args...))
	}
}

// Error logs an error message
func Error(msg string, args ...interface{}) {
	if shouldLog(LevelError) {
		ensureLogger()
		logger.Print("[ERROR] " + formatMessage(msg, args...))
	}
}

// Fatal logs a fatal message and exits
func Fatal(msg string, args ...interface{}) {
	ensureLogger()
	logger.Print("[FATAL] " + formatMessage(msg, args...))
	os.Exit(1)
}
