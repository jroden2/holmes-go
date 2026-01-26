package public

import "github.com/gin-gonic/gin"

type encodeController struct {
}

func NewEncodeController() EncodeController {
	return &encodeController{}
}

type EncodeController interface{}

func (c *encodeController) EncodeSha256(ctx *gin.Context) {

}
