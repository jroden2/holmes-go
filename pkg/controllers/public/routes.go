package public

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func Routes(route *gin.RouterGroup, logger *zerolog.Logger) {

	baseControllerGroup := route.Group("")
	{
		bc := NewBaseController(logger)
		baseControllerGroup.GET("/", bc.Home)
		baseControllerGroup.POST("/compare", bc.Compare)
	}
	encodeControllerGroup := route.Group("sha")
	{
		ec := NewEncodeController(logger, nil)
		encodeControllerGroup.POST("/encode", ec.EncodeSha256)
		encodeControllerGroup.POST("/compute", ec.ComputeSha256)
	}
}
