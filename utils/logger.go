package utils

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LoggerConfig struct {
	LoggerLevel    int    `default:"-1"`
	ATLoggerLevel  int    `default:"-1"`
	ExtLoggerLevel int    `default:"-1"`
	ATLogPath      string `default:"/var/log/at/event.log"`
	ExtLogPath     string `default:"/var/log/schematics/%s.log"`
}

const (
	loggerLabel = "logger"
)

func newProductionConfig(logConf *LoggerConfig) zap.Config {
	return zap.Config{
		Level:            zap.NewAtomicLevelAt(zapcore.Level(logConf.LoggerLevel)),
		Development:      false,
		Encoding:         "json",
		EncoderConfig:    newProductionEncoderConfig(),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}
}

func newProductionEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "utc_date",
		LevelKey:       "level",
		NameKey:        loggerLabel,
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

// NewLogger ...
func NewLogger() (logger *zap.Logger, err error) {
	return newProductionConfig(&LoggerConfig{
		LoggerLevel:    -1,
		ATLoggerLevel:  -1,
		ExtLoggerLevel: -1,
	}).Build()
}
