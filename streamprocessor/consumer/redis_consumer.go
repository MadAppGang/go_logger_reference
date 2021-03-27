package consumer

import (
	"context"
	"go.uber.org/zap"
	"go_logger_reference/log"

	"go_logger_reference/streamprocessor/model"
)

func NewRedisConsumer(config log.ZapConfig) (Consumer, error) {
	logger, err := log.NewZap(config)
	if err != nil {
		return nil, err
	}

	logger = logger.Named("redis-consumer")

	return &redisConsumer{
		logger: logger,
	}, nil
}

type redisConsumer struct {
	logger *zap.SugaredLogger
}

func (r *redisConsumer) Consume(ctx context.Context, unit model.TransformedUnit) error {
	r.logger.With("id", unit.ID).With("payload", unit.AgregatedPayload).Debug("consuming unit")
	r.logger.With("id", unit.ID).Info("unit consumed")

	return nil
}
