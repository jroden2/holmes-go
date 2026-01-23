package utils

import (
	"io"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func ReadGinFile(c *gin.Context, field string) (string, string) {
	fh, err := c.FormFile(field)
	if err != nil || fh == nil {
		return "", ""
	}

	f, err := fh.Open()
	if err != nil {
		return "", ""
	}
	defer f.Close()

	const max = 16 << 20 // 16MB
	b, err := io.ReadAll(io.LimitReader(f, max))
	if err != nil {
		return "", ""
	}

	return string(b), filepath.Base(fh.Filename)
}
