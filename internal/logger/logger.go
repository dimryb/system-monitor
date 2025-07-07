package logger

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
)

type Logger struct {
	level slog.Level
}

func New(level string) *Logger {
	lvl := parseLevel(level)
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: lvl})
	logger := slog.New(handler)
	slog.SetDefault(logger)

	return &Logger{
		level: lvl,
	}
}

func parseLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn", "warning":
		return slog.LevelWarn
	case "error", "err":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func (l *Logger) Debugf(format string, args ...any) {
	slog.Debug(fmt.Sprintf(format, args...))
}

func (l *Logger) Infof(format string, args ...any) {
	slog.Info(fmt.Sprintf(format, args...))
}

func (l *Logger) Warnf(format string, args ...any) {
	slog.Warn(fmt.Sprintf(format, args...))
}

func (l *Logger) Errorf(format string, args ...any) {
	slog.Error(fmt.Sprintf(format, args...))
}

func (l *Logger) Fatalf(format string, args ...any) {
	slog.Error(fmt.Sprintf(format, args...))
	os.Exit(1)
}

func (l *Logger) Debug(format string, args ...any) {
	slog.Debug(format, args...)
}

func (l *Logger) Info(format string, args ...any) {
	slog.Info(format, args...)
}

func (l *Logger) Warn(format string, args ...any) {
	slog.Warn(format, args...)
}

func (l *Logger) Error(format string, args ...any) {
	slog.Error(format, args...)
}

func (l *Logger) Fatal(format string, args ...any) {
	slog.Error(format, args...)
	os.Exit(1)
}
