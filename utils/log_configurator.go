package utils

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"

	"github.com/sirupsen/logrus"
)

func NewLoggerFromConfig(config string) *logrus.Logger {
	_ = config // not used in this reference implementation, but will be used in real implementation
	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.InfoLevel) // set log level based on config
	//logger.SetFormatter() // set formatter based on config
	//logger.SetOutput() // set output based on config

	return logger
}

type ZapConfig struct {
	Format OutputFormat
	MinLogLevel zapcore.Level
}

type OutputFormat string

const (
	JSON OutputFormat = "json"
	PlainText OutputFormat = "console"
)

func NewZap(config ZapConfig) (*zap.SugaredLogger, error) {
	cfg := zap.NewProductionConfig()
	cfg.Level = zap.NewAtomicLevelAt(config.MinLogLevel)
	cfg.Encoding = string(config.Format)

	logger, err := zap.NewProductionConfig().Build()
	if err != nil {
		return nil, err
	}

	return logger.Sugar(), nil
}
