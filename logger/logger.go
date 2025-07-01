package logger

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/startrek92/kube-admission-webhook/config"
)

func InitLogger(cfg *config.Config) error {
	level := parseLevel(cfg.Logging.LogLevel)

	logPath := filepath.Join(cfg.Logging.LogDir, cfg.Logging.LogFile)
	if err := os.MkdirAll(cfg.Logging.LogDir, 0755); err != nil {
		return fmt.Errorf("failed to create log dir: %w", err)
	}

	logFile, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}

	env := strings.ToLower(cfg.Server.Env)
	var handler slog.Handler

	if env == "production" {
		handler = slog.NewJSONHandler(logFile, &slog.HandlerOptions{Level: level})
	} else {
		multi := io.MultiWriter(os.Stdout, logFile)
		handler = slog.NewJSONHandler(multi, &slog.HandlerOptions{Level: level})
	}

	slog.SetDefault(slog.New(handler))
	return nil
}

func parseLevel(str string) slog.Level {
	switch strings.ToLower(str) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
