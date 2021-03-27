package sample

import (
	"context"
	"go.uber.org/zap"
)

const RequestIDKey = "X-Request-ID"

type Handler struct {
	logger *zap.SugaredLogger
}

func NewSampleHandler(logger *zap.SugaredLogger) *Handler {
	return &Handler{
		logger: logger,
	}
}

// Handle imitates gin http handler function.
func (h *Handler) Handle(ctx context.Context) {
	reqID := ctx.Value(RequestIDKey)
	logger := h.logger.With(RequestIDKey, reqID)

	logger.With("sampleStructure", struct {
		UserName string
	}{
		UserName: "User1",
	}).Info("some test log")
}
