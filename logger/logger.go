package logger

import (
	"io"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/startrek92/kube-admission-webhook/config"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Log = logrus.New()

func InitLogger(cfg *config.Config) {
	logDir := cfg.Logging.LogDir
	logFile := cfg.Logging.LogFile

	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		_ = os.MkdirAll(logDir, os.ModePerm)
	}

	logPath := filepath.Join(logDir, logFile)

	rotator := &lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    10,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   true,
	}

	Log.SetOutput(io.MultiWriter(os.Stdout, rotator))
	Log.SetFormatter(&logrus.JSONFormatter{})

	level, err := logrus.ParseLevel(cfg.Logging.LogLevel)
	if err != nil {
		Log.Warnf("Invalid log level in config: %s. Defaulting to info.", cfg.Logging.LogLevel)
		level = logrus.InfoLevel
	}
	Log.SetLevel(level)

	Log.Debugf("Logger initialized with level: %s", level)
}
