package transform

import (
	"context"
	"go.uber.org/zap"
	"go_logger_reference/log"
	"go_logger_reference/streamprocessor/model"
)

func NewTransformer(config log.ZapConfig) (*Transformer, error) {
	logger, err := log.NewZap(config)
	if err != nil {
		return nil, err
	}

	logger = logger.Named("transformer")

	return &Transformer{logger: logger}, nil
}

type Transformer struct {
	logger *zap.SugaredLogger
}

func (t *Transformer) TransformUnit(ctx context.Context, unit model.DataUnit) (model.TransformedUnit, error) {
	if ctx.Err() != nil {
		return model.TransformedUnit{}, ctx.Err()
	}

	return model.TransformedUnit{
		ID:               unit.ID + 1000,
		AgregatedPayload: "Aggregated Payload => " + unit.Payload,
	}, nil
}
