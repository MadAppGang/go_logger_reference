package utils

import "github.com/sirupsen/logrus"

type DefaultFieldHook struct {
	AddField func() (string, interface{})
}

func (h DefaultFieldHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h DefaultFieldHook) Fire(e *logrus.Entry) error {
	if field, value := h.AddField(); field != "" {
		e.Data[field] = value
	}
	return nil
}
