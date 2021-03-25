package consumer

import (
	"context"

	"github.com/sirupsen/logrus"
	"go_logger_reference/streamprocessor/model"
	"go_logger_reference/utils"
)

func NewRedisConsumer(config string) (Consumer, error) {
	logger := utils.NewLoggerFromConfig(config)

	logger.AddHook(utils.DefaultFieldHook{func() (string, interface{}) {
		return "who", "redis-consumer"
	}})

	return &redisConsumer{
		logger: logger,
	}, nil
}

type redisConsumer struct {
	logger *logrus.Logger
}

func (r *redisConsumer) Consume(ctx context.Context, unit model.TransformedUnit) error {
	if r.logger.IsLevelEnabled(logrus.DebugLevel) {
		r.logger.WithField("id", unit.ID).WithField("payload", unit.AgregatedPayload).Debug("unit consumed")
	} else {
		r.logger.WithField("id", unit.ID).Info("unit consumed")
	}

	return nil
}
