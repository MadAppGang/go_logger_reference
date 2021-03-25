package consumer

import (
	"context"

	"go_logger_reference/streamprocessor/model"
)

type Consumer interface {
	Consume(ctx context.Context, unit model.TransformedUnit) error
}
