package producer

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"go_logger_reference/log"
	"math/rand"

	"go_logger_reference/streamprocessor/model"
)

func NewKafkaProducer(config log.ZapConfig) (Producer, error) {
	logger, err := log.NewZap(config)
	if err != nil {
		return nil, err
	}

	logger = logger.Named("kafka-producer")

	return &kafkaProducer{
		logger: logger,
	}, nil
}

type kafkaProducer struct {
	logger *zap.SugaredLogger
}

func (p *kafkaProducer) ProduceOne(context.Context) (model.DataUnit, error) {
	id := int(rand.Int31n(100))
	unit := model.DataUnit{
		ID:      id,
		Payload: fmt.Sprintf("payload-%d", id),
	}
	p.logger.Infow("unit produced", "id", unit.ID)

	return unit, nil
}
