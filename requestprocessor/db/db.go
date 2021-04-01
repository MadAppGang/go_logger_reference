package db

import (
	"go.uber.org/zap"
)

func NewDBSource() *DBSource {
	return &DBSource{}
}

type DBSource struct {
}

func (s *DBSource) SelectSomething(logger *zap.SugaredLogger, tableName string) []string {
	log := logger.With("table", tableName)
	log.Debug("selecting from table")

	extracted := []string{}

	switch tableName {
	case "system":
		extracted = []string{"Create", "Format", "Delete", "Grant"}
	case "months":
		extracted = []string{"Jan", "Feb", "Mar"}
	}

	log.Infow("Extracted data from table", "data",extracted)
	return extracted
}
