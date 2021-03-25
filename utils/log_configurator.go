package utils

import (
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
