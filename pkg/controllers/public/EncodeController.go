package public

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jroden2/holmes-go/pkg/services"
	"github.com/rs/zerolog"
)

type encodeController struct {
	logger *zerolog.Logger
	es     services.EncodeService
}

func NewEncodeController(logger *zerolog.Logger, es *services.EncodeService) EncodeController {
	var newEs services.EncodeService
	if es == nil {
		newEs = services.NewEncodeService(logger)
	}

	return &encodeController{
		logger: logger,
		es:     newEs,
	}
}

type EncodeController interface {
	EncodeSha256(ctx *gin.Context)
	ComputeSha256(ctx *gin.Context)
}

func (c *encodeController) EncodeSha256(ctx *gin.Context) {
	content := ctx.GetString("content")
	ctx.JSON(http.StatusOK, gin.H{
		"content": c.es.EncodeSha256(content),
	})
}

func (c *encodeController) ComputeSha256(ctx *gin.Context) {
	content := ctx.GetString("content")
	comparison := ctx.GetString("comparison")
	ctx.JSON(http.StatusOK, gin.H{
		"result": c.es.ComputeSha256(content, comparison),
	})
}
