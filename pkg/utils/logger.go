package utils

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

func InitLogger() {
	encoderConf := zap.NewProductionEncoderConfig()
	encoderConf.EncodeTime = zapcore.ISO8601TimeEncoder
	config := zap.Config{
		Encoding:      "json",
		Level:         zap.NewAtomicLevelAt(getLogLevel()),
		OutputPaths:   []string{"stdout"},
		EncoderConfig: encoderConf,
	}

	Logger, _ = config.Build()
}

func FlushLogger() {
	if Logger != nil {
		Logger.Sync()
	}
}

func getLogLevel() zapcore.Level {
	log_level := GetEnvVarOrDefault("LOG_LEVEL", "INFO")
	switch log_level {
	case "INFO":
		return zapcore.InfoLevel
	case "WARN":
		return zapcore.WarnLevel
	case "ERROR":
		return zapcore.ErrorLevel
	default:
		return zapcore.DebugLevel
	}
}

