package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapConfig struct {
	Format      OutputFormat
	MinLogLevel zapcore.Level
}

type OutputFormat string

const (
	JSON      OutputFormat = "json"
	PlainText OutputFormat = "console"
)

func NewZap(config ZapConfig) (*zap.SugaredLogger, error) {
	cfg := zap.NewProductionConfig()
	cfg.Level = zap.NewAtomicLevelAt(config.MinLogLevel)
	cfg.Encoding = string(config.Format)

	logger, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	return logger.Sugar(), nil
}
