package services

import (
	"github.com/jroden2/holmes-go/pkg/utils"
	"github.com/rs/zerolog"
)

type encodeService struct {
	logger *zerolog.Logger
}

func NewEncodeService(logger *zerolog.Logger) EncodeService {
	return &encodeService{
		logger: logger,
	}
}

type EncodeService interface {
	EncodeSha256(content string) string
	ComputeSha256(input, comparisonSha string) bool
}

func (c *encodeService) EncodeSha256(content string) string {
	retVal := utils.Sha256Hex(content)
	c.logger.Info().Msgf("Encoded content: %s", retVal)
	return retVal
}

func (c *encodeService) ComputeSha256(input, comparisonSha string) bool {
	hashed := utils.Sha256Hex(input)
	return hashed == comparisonSha
}
