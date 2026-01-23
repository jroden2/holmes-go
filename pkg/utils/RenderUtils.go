package utils

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"html/template"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jroden2/holmes-go/pkg/domain"
	"github.com/rs/zerolog/log"
)

func Render(c *gin.Context, tpl *template.Template, data domain.PageData) {
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.Status(http.StatusOK)

	if err := tpl.Execute(c.Writer, data); err != nil {
		log.Error().Err(err).Msg("template render failed")
		c.String(http.StatusInternalServerError, "template error")
	}
}

func Sha256Hex(s string) string {
	sum := sha256.Sum256([]byte(s))
	return hex.EncodeToString(sum[:])
}

func NormalizeWhitespace(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

func SplitLines(s string) []string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	if s == "" {
		return []string{}
	}
	return strings.Split(s, "\n")
}

func BasicLineDiffWithHighlight(a, b string) []domain.LineDiffRow {
	aLines := SplitLines(a)
	bLines := SplitLines(b)

	max := len(aLines)
	if len(bLines) > max {
		max = len(bLines)
	}

	out := make([]domain.LineDiffRow, 0, max)

	for i := 0; i < max; i++ {
		var av, bv string
		hasA := i < len(aLines)
		hasB := i < len(bLines)

		if hasA {
			av = aLines[i]
		}
		if hasB {
			bv = bLines[i]
		}

		status := "same"
		switch {
		case hasA && hasB && av == bv:
			status = "same"
		case hasA && hasB && av != bv:
			status = "changed"
		case hasA && !hasB:
			status = "removed"
		case !hasA && hasB:
			status = "added"
		}

		row := domain.LineDiffRow{
			LineNum: i + 1,
			A:       av,
			B:       bv,
			Status:  status,
		}

		if status == "changed" {
			aHTML, bHTML := highlightCharDiff(av, bv)
			row.AHTML = aHTML
			row.BHTML = bHTML
		} else {
			row.AHTML = template.HTML(template.HTMLEscapeString(av))
			row.BHTML = template.HTML(template.HTMLEscapeString(bv))
		}

		out = append(out, row)
	}

	return out
}

func highlightCharDiff(a, b string) (template.HTML, template.HTML) {
	ar := []rune(a)
	br := []rune(b)

	// common prefix
	p := 0
	for p < len(ar) && p < len(br) && ar[p] == br[p] {
		p++
	}

	// common suffix
	as := len(ar)
	bs := len(br)
	for as > p && bs > p && ar[as-1] == br[bs-1] {
		as--
		bs--
	}

	aPrefix := string(ar[:p])
	aMid := string(ar[p:as])
	aSuffix := string(ar[as:])

	bPrefix := string(br[:p])
	bMid := string(br[p:bs])
	bSuffix := string(br[bs:])

	var bufA bytes.Buffer
	bufA.WriteString(template.HTMLEscapeString(aPrefix))
	if aMid != "" {
		bufA.WriteString("<mark>")
		bufA.WriteString(template.HTMLEscapeString(aMid))
		bufA.WriteString("</mark>")
	}
	bufA.WriteString(template.HTMLEscapeString(aSuffix))

	var bufB bytes.Buffer
	bufB.WriteString(template.HTMLEscapeString(bPrefix))
	if bMid != "" {
		bufB.WriteString("<mark>")
		bufB.WriteString(template.HTMLEscapeString(bMid))
		bufB.WriteString("</mark>")
	}
	bufB.WriteString(template.HTMLEscapeString(bSuffix))

	return template.HTML(bufA.String()), template.HTML(bufB.String())
}
