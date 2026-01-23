package public

import "github.com/rs/zerolog"

type comparisonController struct {
	logger *zerolog.Logger
}

func NewComparisonController(Logger *zerolog.Logger) ComparisonController {
	return &comparisonController{}
}

type ComparisonController interface {
}
