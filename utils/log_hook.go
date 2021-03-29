package utils

import "github.com/sirupsen/logrus"

func LogDefaultField(name string, value interface{}) logrus.Hook {
	return defaultFieldHook{name, value}
}

type defaultFieldHook struct {
	name  string
	value interface{}
}

func (h defaultFieldHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h defaultFieldHook) Fire(e *logrus.Entry) error {
	if h.name != "" {
		e.Data[h.name] = h.value
	}
	return nil
}
