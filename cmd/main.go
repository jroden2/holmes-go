package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"html/template"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jroden2/holmes-go/pkg/domain"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// Zerolog setup (console output)
	zerolog.TimeFieldFormat = time.RFC3339Nano
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: gin.DefaultWriter})

	// Parse template
	tpl := template.Must(template.ParseFiles("templates/index.html"))

	// Gin setup
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(zerologMiddleware())

	// Local-only bind
	addr := "127.0.0.1:8080"

	r.GET("/", func(c *gin.Context) {
		render(c, tpl, domain.PageData{Mode: "text"})
	})

	// Single endpoint: format/pretty + compare
	r.POST("/action", func(c *gin.Context) {
		action := c.PostForm("action") // compare | format_a | format_b | format_both

		mode := c.PostForm("mode")
		if mode != "json" && mode != "xml" {
			mode = "text"
		}

		ignoreWS := c.PostForm("ignore_ws") == "on"
		ignoreCase := c.PostForm("ignore_case") == "on"

		// Textareas
		a := strings.TrimRight(c.PostForm("a"), "\r\n")
		b := strings.TrimRight(c.PostForm("b"), "\r\n")

		// Uploaded files override textarea if present
		if fa, _ := readGinFile(c, "file_a"); fa != "" {
			a = fa
		}
		if fb, _ := readGinFile(c, "file_b"); fb != "" {
			b = fb
		}

		// Pretty-print actions
		if action == "format_a" || action == "format_b" || action == "format_both" {
			var err error

			switch mode {
			case "json":
				if action == "format_a" || action == "format_both" {
					a, err = prettyJSON(a)
					if err != nil {
						render(c, tpl, PageData{
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
					b, err = prettyJSON(b)
					if err != nil {
						render(c, tpl, PageData{
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
					a, err = prettyXML(a)
					if err != nil {
						render(c, tpl, PageData{
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
					b, err = prettyXML(b)
					if err != nil {
						render(c, tpl, PageData{
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

			render(c, tpl, PageData{
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
			compareA, err = prettyJSON(a)
			if err != nil {
				render(c, tpl, PageData{
					A:          a,
					B:          b,
					Mode:       mode,
					IgnoreWS:   ignoreWS,
					IgnoreCase: ignoreCase,
					Error:      "JSON parse error for A: " + err.Error(),
				})
				return
			}
			compareB, err = prettyJSON(b)
			if err != nil {
				render(c, tpl, PageData{
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
			compareA, err = prettyXML(a)
			if err != nil {
				render(c, tpl, PageData{
					A:          a,
					B:          b,
					Mode:       mode,
					IgnoreWS:   ignoreWS,
					IgnoreCase: ignoreCase,
					Error:      "XML parse error for A: " + err.Error(),
				})
				return
			}
			compareB, err = prettyXML(b)
			if err != nil {
				render(c, tpl, PageData{
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
			na = normalizeWhitespace(na)
			nb = normalizeWhitespace(nb)
		}
		if ignoreCase {
			na = strings.ToLower(na)
			nb = strings.ToLower(nb)
		}

		normalized := na == nb

		data := PageData{
			A:          a,
			B:          b,
			Mode:       mode,
			IgnoreWS:   ignoreWS,
			IgnoreCase: ignoreCase,

			ExactMatch:      exact,
			NormalizedMatch: normalized,

			ALen: len(compareA),
			BLen: len(compareB),

			AHash: sha256Hex(compareA),
			BHash: sha256Hex(compareB),

			LineDiff: basicLineDiffWithHighlight(compareA, compareB),
		}

		render(c, tpl, data)
	})

	log.Info().Str("addr", addr).Msg("starting local server (gin/nethttp)")
	if err := r.Run(addr); err != nil {
		log.Fatal().Err(err).Msg("server stopped")
	}
}

func zerologMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		log.Info().
			Str("method", c.Request.Method).
			Str("path", c.Request.URL.Path).
			Int("status", c.Writer.Status()).
			Dur("duration", time.Since(start)).
			Msg("request")
	}
}

func render(c *gin.Context, tpl *template.Template, data PageData) {
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.Status(http.StatusOK)

	if err := tpl.Execute(c.Writer, data); err != nil {
		log.Error().Err(err).Msg("template render failed")
		c.String(http.StatusInternalServerError, "template error")
	}
}

func readGinFile(c *gin.Context, field string) (string, string) {
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

// ===== formatters =====

func prettyJSON(s string) (string, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return "", nil
	}

	var v any
	dec := json.NewDecoder(strings.NewReader(s))
	dec.UseNumber()

	if err := dec.Decode(&v); err != nil {
		return "", err
	}

	out, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", err
	}

	return string(out), nil
}

func prettyXML(s string) (string, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return "", nil
	}

	dec := xml.NewDecoder(strings.NewReader(s))

	var buf bytes.Buffer
	enc := xml.NewEncoder(&buf)
	enc.Indent("", "  ")

	for {
		tok, err := dec.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}
		if err := enc.EncodeToken(tok); err != nil {
			return "", err
		}
	}

	if err := enc.Flush(); err != nil {
		return "", err
	}

	out := strings.TrimSpace(buf.String()) + "\n"
	return out, nil
}

// ===== diff helpers =====

func sha256Hex(s string) string {
	sum := sha256.Sum256([]byte(s))
	return hex.EncodeToString(sum[:])
}

func normalizeWhitespace(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

func splitLines(s string) []string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	if s == "" {
		return []string{}
	}
	return strings.Split(s, "\n")
}

func basicLineDiffWithHighlight(a, b string) []LineDiffRow {
	aLines := splitLines(a)
	bLines := splitLines(b)

	max := len(aLines)
	if len(bLines) > max {
		max = len(bLines)
	}

	out := make([]LineDiffRow, 0, max)

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

		row := LineDiffRow{
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
