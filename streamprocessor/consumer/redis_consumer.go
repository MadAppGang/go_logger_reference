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
	r.logger.Debugw("consuming unit", "payload", unit.AgregatedPayload, "id", unit.ID)
	r.logger.Infow("unit consumed", "id", unit.ID)

	return nil
}
