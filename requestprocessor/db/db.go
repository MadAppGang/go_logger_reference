package db

import "github.com/sirupsen/logrus"

func NewDBSource() *DBSource {
	return &DBSource{}
}

type DBSource struct {
}

func (s *DBSource) SelectSomething(logger *logrus.Logger, tableName string) []string {
	logger.WithField("table", tableName).Debug("selecting from table")

	extracted := []string{}

	switch tableName {
	case "system":
		extracted = []string{"Create", "Format", "Delete", "Grant"}
	case "months":
		extracted = []string{"Jan", "Feb", "Mar"}
	}
	logger.WithField("table", tableName).Infof("Extracted data from table: %v", extracted)
	return extracted
}
