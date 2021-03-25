package producer

import (
	"context"

	"go_logger_reference/streamprocessor/model"
)

type Producer interface {
	ProduceOne(context.Context) (model.DataUnit, error)
}
