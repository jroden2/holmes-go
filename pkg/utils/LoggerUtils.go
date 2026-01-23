package utils

import (
	"os"

	"github.com/rs/zerolog"
)

func InitLogger() *zerolog.Logger {
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
	return &logger
}
