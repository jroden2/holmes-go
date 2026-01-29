package public

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jroden2/holmes-go/pkg/domain"
	"github.com/jroden2/holmes-go/pkg/services"
	"github.com/jroden2/holmes-go/pkg/utils"
	"github.com/rs/zerolog"
)

type baseController struct {
	logger *zerolog.Logger
	sonic  services.CacheService
}

func NewBaseController(logger *zerolog.Logger) BaseController {
	return &baseController{
		logger: logger,
		sonic:  services.NewCacheService(),
	}
}

type BaseController interface {
	Home(ctx *gin.Context)
	Compare(ctx *gin.Context)

	// Used for cachable content
	CreateMagicKey(ctx *gin.Context)
	CompareUsingMagicLink(ctx *gin.Context)
	PeekMagicKeys(ctx *gin.Context)
}

func (c *baseController) Home(ctx *gin.Context) {
	tpl, err := loadTemplates()
	if err != nil {
		c.logger.Fatal().Err(err).Msg("Failed to load templates")
	}
	utils.Render(ctx, tpl, domain.PageData{Mode: "auto"})
}

func (c *baseController) CreateMagicKey(ctx *gin.Context) {
	a := strings.TrimRight(ctx.PostForm("a"), "\r\n")
	b := strings.TrimRight(ctx.PostForm("b"), "\r\n")

	id, _ := utils.Generate32CharString()
	var MagicPayload domain.DiffPayload
	MagicPayload.ID = id
	MagicPayload.ShortID = id[:8]
	MagicPayload.Original = a
	MagicPayload.New = b

	blob, err := json.Marshal(MagicPayload)
	if err != nil {
		c.logger.Error().Err(err).Msg("Failed to marshal magic link")
		ctx.Redirect(http.StatusFound, "/error")
		return
	}

	c.sonic.Add(MagicPayload.GetID(), blob)
	ctx.JSON(http.StatusOK, gin.H{
		"id": MagicPayload.GetID(),
	})
}

func (c *baseController) PeekMagicKeys(ctx *gin.Context) {
	kvp := c.sonic.PeekAll()
	keys := make([]string, 0, len(kvp))
	for k := range kvp {
		if keyStr, ok := k.(string); ok {
			keys = append(keys, keyStr)
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"keys": keys,
	})
}

func (c *baseController) CompareUsingMagicLink(ctx *gin.Context) {
	magicLink := ctx.Query("magicLink")
	if magicLink == "" {
		c.logger.Error().Msg("No magic link provided")
		ctx.Redirect(http.StatusFound, "/error") // 302
		return
	}

	var sessionPayload domain.DiffPayload
	if blob, exists := c.sonic.Get(magicLink); !exists {
		c.logger.Error().Msg("Magic link does not exist")
		ctx.Redirect(http.StatusFound, "/error")
		return
	} else {
		err := json.Unmarshal(blob, &sessionPayload)
		if err != nil {
			c.logger.Error().Err(err).Msg("Failed to unmarshal magic link")
			ctx.Redirect(http.StatusFound, "/error")
			return
		}
	}

	session := sessions.Default(ctx)
	session.Set("a", sessionPayload.Original)
	session.Set("a", sessionPayload.New)
	err := session.Save()
	if err != nil {
		c.logger.Error().Err(err).Msg("Failed to save session")
		ctx.Redirect(http.StatusFound, "/error")
		return
	}

	ctx.Redirect(http.StatusFound, "/compare")
}

func (c *baseController) Compare(ctx *gin.Context) {
	tpl, err := loadTemplates()
	if err != nil {
		c.logger.Fatal().Err(err).Msg("Failed to load templates")
	}
	action := ctx.PostForm("action") // compare | format_a | format_b | format_both

	mode := ctx.PostForm("mode")
	if mode != "json" && mode != "xml" {
		mode = "text"
	}

	ignoreWS := ctx.PostForm("ignore_ws") == "on"
	ignoreCase := ctx.PostForm("ignore_case") == "on"

	// Textareas
	a := strings.TrimRight(ctx.PostForm("a"), "\r\n")
	b := strings.TrimRight(ctx.PostForm("b"), "\r\n")

	// Uploaded files override textarea if present
	if fa, _ := utils.ReadGinFile(ctx, "file_a"); fa != "" {
		a = fa
	}
	if fb, _ := utils.ReadGinFile(ctx, "file_b"); fb != "" {
		b = fb
	}

	// Pretty-print actions
	if action == "format_a" || action == "format_b" || action == "format_both" {
		var err error

		switch mode {
		case "json":
			if action == "format_a" || action == "format_both" {
				a, err = utils.PrettyJSON(a)
				if err != nil {
					utils.Render(ctx, tpl, domain.PageData{
						A:          a,
						B:          b,
						Mode:       mode,
						IgnoreWS:   ignoreWS,
						IgnoreCase: ignoreCase,
						Error:      "Pretty JSON A failed: " + err.Error(),
					})
					return
				}
			}
			if action == "format_b" || action == "format_both" {
				b, err = utils.PrettyJSON(b)
				if err != nil {
					utils.Render(ctx, tpl, domain.PageData{
						A:          a,
						B:          b,
						Mode:       mode,
						IgnoreWS:   ignoreWS,
						IgnoreCase: ignoreCase,
						Error:      "Pretty JSON B failed: " + err.Error(),
					})
					return
				}
			}

		case "xml":
			if action == "format_a" || action == "format_both" {
				a, err = utils.PrettyXML(a)
				if err != nil {
					utils.Render(ctx, tpl, domain.PageData{
						A:          a,
						B:          b,
						Mode:       mode,
						IgnoreWS:   ignoreWS,
						IgnoreCase: ignoreCase,
						Error:      "Pretty XML A failed: " + err.Error(),
					})
					return
				}
			}
			if action == "format_b" || action == "format_both" {
				b, err = utils.PrettyXML(b)
				if err != nil {
					utils.Render(ctx, tpl, domain.PageData{
						A:          a,
						B:          b,
						Mode:       mode,
						IgnoreWS:   ignoreWS,
						IgnoreCase: ignoreCase,
						Error:      "Pretty XML B failed: " + err.Error(),
					})
					return
				}
			}

		default:
			// text mode: do nothing
		}

		utils.Render(ctx, tpl, domain.PageData{
			A:          a,
			B:          b,
			Mode:       mode,
			IgnoreWS:   ignoreWS,
			IgnoreCase: ignoreCase,
		})
		return
	}

	// Compare action (default)
	if action == "" {
		action = "compare"
	}

	compareA := a
	compareB := b

	// If mode is json/xml, compare normalized/pretty versions for stable diffs
	if mode == "json" {
		var err error
		compareA, err = utils.PrettyJSON(a)
		if err != nil {
			utils.Render(ctx, tpl, domain.PageData{
				A:          a,
				B:          b,
				Mode:       mode,
				IgnoreWS:   ignoreWS,
				IgnoreCase: ignoreCase,
				Error:      "JSON parse error for A: " + err.Error(),
			})
			return
		}
		compareB, err = utils.PrettyJSON(b)
		if err != nil {
			utils.Render(ctx, tpl, domain.PageData{
				A:          a,
				B:          b,
				Mode:       mode,
				IgnoreWS:   ignoreWS,
				IgnoreCase: ignoreCase,
				Error:      "JSON parse error for B: " + err.Error(),
			})
			return
		}
	} else if mode == "xml" {
		var err error
		compareA, err = utils.PrettyXML(a)
		if err != nil {
			utils.Render(ctx, tpl, domain.PageData{
				A:          a,
				B:          b,
				Mode:       mode,
				IgnoreWS:   ignoreWS,
				IgnoreCase: ignoreCase,
				Error:      "XML parse error for A: " + err.Error(),
			})
			return
		}
		compareB, err = utils.PrettyXML(b)
		if err != nil {
			utils.Render(ctx, tpl, domain.PageData{
				A:          a,
				B:          b,
				Mode:       mode,
				IgnoreWS:   ignoreWS,
				IgnoreCase: ignoreCase,
				Error:      "XML parse error for B: " + err.Error(),
			})
			return
		}
	}

	exact := compareA == compareB

	na := compareA
	nb := compareB
	if ignoreWS {
		na = utils.NormalizeWhitespace(na)
		nb = utils.NormalizeWhitespace(nb)
	}
	if ignoreCase {
		na = strings.ToLower(na)
		nb = strings.ToLower(nb)
	}

	normalized := na == nb

	data := domain.PageData{
		A:          a,
		B:          b,
		Mode:       mode,
		IgnoreWS:   ignoreWS,
		IgnoreCase: ignoreCase,

		ExactMatch:      exact,
		NormalizedMatch: normalized,

		ALen: len(compareA),
		BLen: len(compareB),

		AHash: utils.Sha256Hex(compareA),
		BHash: utils.Sha256Hex(compareB),

		LineDiff: utils.BasicLineDiffWithHighlight(compareA, compareB),
	}

	utils.Render(ctx, tpl, data)
}

func loadTemplates() (*template.Template, error) {
	return template.ParseFiles("./templates/index.html")
}
