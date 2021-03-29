package producer

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/sirupsen/logrus"
	"go_logger_reference/streamprocessor/model"
	"go_logger_reference/utils"
)

func NewKafkaProducer(config string) (Producer, error) {
	logger := utils.NewLoggerFromConfig(config)
	logger.AddHook(utils.LogDefaultField("who", "kafka-producer"))

	return &kafkaProducer{
		logger: logger,
	}, nil
}

type kafkaProducer struct {
	logger *logrus.Logger
}

func (p *kafkaProducer) ProduceOne(context.Context) (model.DataUnit, error) {
	id := int(rand.Int31n(100))
	unit := model.DataUnit{
		ID:      id,
		Payload: fmt.Sprintf("payload-%d", id),
	}
	p.logger.WithField("id", unit.ID).Info("unit produced")

	return unit, nil
}
