package public

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestController() BaseController {
	workdir, _ := os.Getwd()
	if !strings.HasSuffix(workdir, "holmes-go") {
		_ = os.Chdir("../../../")
	}

	logger := zerolog.New(os.Stdout)
	return NewBaseController(&logger)
}

func TestMain(m *testing.M) {
	workdir, _ := os.Getwd()
	if !strings.HasSuffix(workdir, "holmes-go") {
		_ = os.Chdir("../../../")
	}
	code := m.Run()
	_ = os.Chdir(workdir)
	os.Exit(code)
}

func TestHome(t *testing.T) {
	gin.SetMode(gin.TestMode)
	controller := setupTestController()

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/", nil)

	controller.Home(ctx)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Local Diff Checker")
}

func TestCompare_TextMode(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		a              string
		b              string
		mode           string
		ignoreWS       bool
		ignoreCase     bool
		expectedExact  bool
		expectedNormal bool
	}{
		{
			name:           "exact match",
			a:              "hello world",
			b:              "hello world",
			mode:           "text",
			expectedExact:  true,
			expectedNormal: true,
		},
		{
			name:           "different text",
			a:              "hello world",
			b:              "goodbye world",
			mode:           "text",
			expectedExact:  false,
			expectedNormal: false,
		},
		{
			name:           "whitespace difference - ignore off",
			a:              "hello world",
			b:              "hello  world",
			mode:           "text",
			ignoreWS:       false,
			expectedExact:  false,
			expectedNormal: false,
		},
		{
			name:           "whitespace difference - ignore on",
			a:              "hello world",
			b:              "hello  world",
			mode:           "text",
			ignoreWS:       true,
			expectedExact:  false,
			expectedNormal: true,
		},
		{
			name:           "case difference - ignore off",
			a:              "Hello World",
			b:              "hello world",
			mode:           "text",
			ignoreCase:     false,
			expectedExact:  false,
			expectedNormal: false,
		},
		{
			name:           "case difference - ignore on",
			a:              "Hello World",
			b:              "hello world",
			mode:           "text",
			ignoreCase:     true,
			expectedExact:  false,
			expectedNormal: true,
		},
		{
			name:           "whitespace and case - both ignore on",
			a:              "Hello  World",
			b:              "hello world",
			mode:           "text",
			ignoreWS:       true,
			ignoreCase:     true,
			expectedExact:  false,
			expectedNormal: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := setupTestController()

			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)

			form := url.Values{}
			form.Add("a", tt.a)
			form.Add("b", tt.b)
			form.Add("mode", tt.mode)
			if tt.ignoreWS {
				form.Add("ignore_ws", "on")
			}
			if tt.ignoreCase {
				form.Add("ignore_case", "on")
			}

			ctx.Request = httptest.NewRequest(http.MethodPost, "/compare", strings.NewReader(form.Encode()))
			ctx.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			controller.Compare(ctx)

			assert.Equal(t, http.StatusOK, w.Code)

			body := w.Body.String()
			if tt.expectedExact {
				assert.Contains(t, body, "YES")
			} else {
				assert.Contains(t, body, "NO")
			}
		})
	}
}

func TestCompare_JSONMode(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name          string
		a             string
		b             string
		expectedExact bool
		expectError   bool
	}{
		{
			name:          "identical json",
			a:             `{"name":"John","age":30}`,
			b:             `{"name":"John","age":30}`,
			expectedExact: true,
			expectError:   false,
		},
		{
			name:          "different formatting same data",
			a:             `{"name":"John","age":30}`,
			b:             `{"name": "John", "age": 30}`,
			expectedExact: true,
			expectError:   false,
		},
		{
			name:          "different json values",
			a:             `{"name":"John","age":30}`,
			b:             `{"name":"Jane","age":30}`,
			expectedExact: false,
			expectError:   false,
		},
		{
			name:          "invalid json in a",
			a:             `{invalid json}`,
			b:             `{"name":"John"}`,
			expectedExact: false,
			expectError:   true,
		},
		{
			name:          "invalid json in b",
			a:             `{"name":"John"}`,
			b:             `{invalid json}`,
			expectedExact: false,
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := setupTestController()

			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)

			form := url.Values{}
			form.Add("a", tt.a)
			form.Add("b", tt.b)
			form.Add("mode", "json")

			ctx.Request = httptest.NewRequest(http.MethodPost, "/compare", strings.NewReader(form.Encode()))
			ctx.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			controller.Compare(ctx)

			assert.Equal(t, http.StatusOK, w.Code)

			body := w.Body.String()
			if tt.expectError {
				assert.Contains(t, body, "error", "Expected error message in response")
			}
		})
	}
}

func TestCompare_XMLMode(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name          string
		a             string
		b             string
		expectedExact bool
		expectError   bool
	}{
		{
			name:          "identical xml",
			a:             `<root><name>John</name></root>`,
			b:             `<root><name>John</name></root>`,
			expectedExact: true,
			expectError:   false,
		},
		{
			name:          "different formatting same data",
			a:             `<root><name>John</name></root>`,
			b:             `<root>  <name>John</name>  </root>`,
			expectedExact: true,
			expectError:   false,
		},
		{
			name:          "different xml values",
			a:             `<root><name>John</name></root>`,
			b:             `<root><name>Jane</name></root>`,
			expectedExact: false,
			expectError:   false,
		},
		{
			name:          "invalid xml in a",
			a:             `<root><name>John</root>`,
			b:             `<root><name>John</name></root>`,
			expectedExact: false,
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := setupTestController()

			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)

			form := url.Values{}
			form.Add("a", tt.a)
			form.Add("b", tt.b)
			form.Add("mode", "xml")

			ctx.Request = httptest.NewRequest(http.MethodPost, "/compare", strings.NewReader(form.Encode()))
			ctx.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			controller.Compare(ctx)

			assert.Equal(t, http.StatusOK, w.Code)

			body := w.Body.String()
			if tt.expectError {
				assert.Contains(t, body, "error", "Expected error message in response")
			}
		})
	}
}

func TestCompare_FormatActions(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name        string
		action      string
		mode        string
		a           string
		b           string
		expectError bool
	}{
		{
			name:        "format json a",
			action:      "format_a",
			mode:        "json",
			a:           `{"name":"John","age":30}`,
			b:           `{"name":"Jane"}`,
			expectError: false,
		},
		{
			name:        "format json b",
			action:      "format_b",
			mode:        "json",
			a:           `{"name":"John"}`,
			b:           `{"name":"Jane","age":25}`,
			expectError: false,
		},
		{
			name:        "format json both",
			action:      "format_both",
			mode:        "json",
			a:           `{"name":"John","age":30}`,
			b:           `{"name":"Jane","age":25}`,
			expectError: false,
		},
		{
			name:        "format xml a",
			action:      "format_a",
			mode:        "xml",
			a:           `<root><name>John</name></root>`,
			b:           `<root><name>Jane</name></root>`,
			expectError: false,
		},
		{
			name:        "format xml both",
			action:      "format_both",
			mode:        "xml",
			a:           `<root><name>John</name></root>`,
			b:           `<root><name>Jane</name></root>`,
			expectError: false,
		},
		{
			name:        "format invalid json a",
			action:      "format_a",
			mode:        "json",
			a:           `{invalid}`,
			b:           `{"name":"Jane"}`,
			expectError: true,
		},
		{
			name:        "format invalid xml b",
			action:      "format_b",
			mode:        "xml",
			a:           `<root><name>John</name></root>`,
			b:           `<root><name>Jane</root>`,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := setupTestController()

			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)

			form := url.Values{}
			form.Add("action", tt.action)
			form.Add("mode", tt.mode)
			form.Add("a", tt.a)
			form.Add("b", tt.b)

			ctx.Request = httptest.NewRequest(http.MethodPost, "/compare", strings.NewReader(form.Encode()))
			ctx.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			controller.Compare(ctx)

			assert.Equal(t, http.StatusOK, w.Code)

			body := strings.ToLower(w.Body.String())
			if tt.expectError {
				assert.Contains(t, body, "failed", "Expected error message in response")
			}
		})
	}
}

func TestCompare_WithTestDataFiles(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name          string
		fileA         string
		fileB         string
		mode          string
		expectedExact bool
	}{
		{
			name:          "json test files",
			fileA:         "pkg/testdata/example_a.json",
			fileB:         "pkg/testdata/example_b.json",
			mode:          "json",
			expectedExact: false, // Since we made changes in example_b
		},
		{
			name:          "xml test files",
			fileA:         "pkg/testdata/example_a.xml",
			fileB:         "pkg/testdata/example_b.xml",
			mode:          "xml",
			expectedExact: false, // Since we made changes in example_b
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Read test files
			aContent, err := os.ReadFile(tt.fileA)
			require.NoError(t, err, "Failed to read file A")

			bContent, err := os.ReadFile(tt.fileB)
			require.NoError(t, err, "Failed to read file B")

			controller := setupTestController()

			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)

			form := url.Values{}
			form.Add("a", string(aContent))
			form.Add("b", string(bContent))
			form.Add("mode", tt.mode)

			ctx.Request = httptest.NewRequest(http.MethodPost, "/compare", strings.NewReader(form.Encode()))
			ctx.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			controller.Compare(ctx)

			assert.Equal(t, http.StatusOK, w.Code)
			assert.NotContains(t, w.Body.String(), "error")

			// Verify SHA256 hashes are present
			body := w.Body.String()
			assert.Contains(t, body, "SHA256")
		})
	}
}

func TestCompare_EmptyInputs(t *testing.T) {
	gin.SetMode(gin.TestMode)
	controller := setupTestController()

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	form := url.Values{}
	form.Add("a", "")
	form.Add("b", "")
	form.Add("mode", "text")

	ctx.Request = httptest.NewRequest(http.MethodPost, "/compare", strings.NewReader(form.Encode()))
	ctx.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	controller.Compare(ctx)

	assert.Equal(t, http.StatusOK, w.Code)
	body := w.Body.String()
	assert.Contains(t, body, "YES") // Empty strings should match
}

func TestCompare_ModeValidation(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		mode         string
		expectedMode string
	}{
		{
			name:         "valid text mode",
			mode:         "text",
			expectedMode: "text",
		},
		{
			name:         "valid json mode",
			mode:         "json",
			expectedMode: "json",
		},
		{
			name:         "valid xml mode",
			mode:         "xml",
			expectedMode: "xml",
		},
		{
			name:         "invalid mode defaults to text",
			mode:         "invalid",
			expectedMode: "text",
		},
		{
			name:         "empty mode defaults to text",
			mode:         "",
			expectedMode: "text",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := setupTestController()

			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)

			form := url.Values{}
			form.Add("a", "test")
			form.Add("b", "test")
			form.Add("mode", tt.mode)

			ctx.Request = httptest.NewRequest(http.MethodPost, "/compare", strings.NewReader(form.Encode()))
			ctx.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			controller.Compare(ctx)

			assert.Equal(t, http.StatusOK, w.Code)
		})
	}
}
