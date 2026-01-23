package public

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func Routes(route *gin.RouterGroup, logger *zerolog.Logger) {
	baseController := NewBaseController(logger)

	baseControllerGroup := route.Group("")
	{
		baseControllerGroup.GET("/", baseController.Home)
		baseControllerGroup.POST("/compare", baseController.Compare)

	}
}
