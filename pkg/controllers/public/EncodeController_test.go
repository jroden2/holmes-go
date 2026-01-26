package public

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jroden2/holmes-go/pkg/services"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestNewEncodeController(t *testing.T) {
	logger := zerolog.New(os.Stdout)

	t.Run("with nil service", func(t *testing.T) {
		controller := NewEncodeController(logger, nil)
		assert.NotNil(t, controller)
	})

	t.Run("with provided service", func(t *testing.T) {
		es := services.NewEncodeService(&logger)
		controller := NewEncodeController(logger, &es)
		assert.NotNil(t, controller)
	})
}

func TestEncodeSha256(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := zerolog.New(os.Stdout)
	controller := NewEncodeController(logger, nil)

	t.Run("successful encoding", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Set("content", "test string")

		controller.EncodeSha256(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"content"`)
		// sha256 of "test string" is d5579c46dfcc7f18207013e65b44e4cb4e2c2298f4ac457ba8f82743f31e930b
		assert.Contains(t, w.Body.String(), "d5579c46dfcc7f18207013e65b44e4cb4e2c2298f4ac457ba8f82743f31e930b")
	})

	t.Run("empty content", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Set("content", "")

		controller.EncodeSha256(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		// sha256 of "" is e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
		assert.Contains(t, w.Body.String(), "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855")
	})
}

func TestComputeSha256(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := zerolog.New(os.Stdout)
	controller := NewEncodeController(logger, nil)

	t.Run("matching hash", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Set("content", "test string")
		ctx.Set("comparison", "d5579c46dfcc7f18207013e65b44e4cb4e2c2298f4ac457ba8f82743f31e930b")

		controller.ComputeSha256(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"result":true`)
	})

	t.Run("non-matching hash", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Set("content", "test string")
		ctx.Set("comparison", "wrong hash")

		controller.ComputeSha256(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"result":false`)
	})
}
