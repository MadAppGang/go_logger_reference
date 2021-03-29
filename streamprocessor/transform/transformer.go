package transform

import (
	"context"

	"github.com/sirupsen/logrus"
	"go_logger_reference/streamprocessor/model"
	"go_logger_reference/utils"
)

func NewTransformer(config string) (*Transformer, error) {
	logger := utils.NewLoggerFromConfig(config)

	logger.AddHook(utils.LogDefaultField("who", "transformer"))

	return &Transformer{logger: logger}, nil
}

type Transformer struct {
	logger *logrus.Logger
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
