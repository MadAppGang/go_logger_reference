package streamprocessor

import (
	"go_logger_reference/log"
	"go_logger_reference/streamprocessor/consumer"
	"go_logger_reference/streamprocessor/producer"
	"go_logger_reference/streamprocessor/transform"
)

func BuildService(config log.ZapConfig) *Service {
	kafkaProducer, _ := producer.NewKafkaProducer(config)
	redisConsumer, _ := consumer.NewRedisConsumer(config)
	transformer, _ := transform.NewTransformer(config)

	service := NewService(kafkaProducer, transformer, redisConsumer)

	return service
}
