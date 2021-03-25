package streamprocessor

import (
	"context"
	"log"
	"time"

	"go_logger_reference/streamprocessor/consumer"
	"go_logger_reference/streamprocessor/producer"
	"go_logger_reference/streamprocessor/transform"
)

func NewService(config string, producer producer.Producer, transformer *transform.Transformer, consumer consumer.Consumer) *Service {
	_ = config
	return &Service{
		producer:    producer,
		transformer: transformer,
		consumer:    consumer,
	}
}

type Service struct {
	producer    producer.Producer
	transformer *transform.Transformer
	consumer    consumer.Consumer
}

func (s Service) Run(ctx context.Context) {
	for true {
		select {
		case <-ctx.Done():
			log.Printf("stop")
			return
		default:
		}

		dataUnit, err := s.producer.ProduceOne(ctx)
		if err != nil {
			// log
			continue
		}
		transformedUnit, err := s.transformer.TransformUnit(ctx, dataUnit)
		if err != nil {

		}

		err = s.consumer.Consume(ctx, transformedUnit)
		if err != nil {

		}

		time.Sleep(500 * time.Millisecond)
	}
}
