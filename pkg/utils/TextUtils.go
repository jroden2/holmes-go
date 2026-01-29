package utils

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"io"
	"strings"
)

func PrettyJSON(s string) (string, error) {
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

func PrettyXML(s string) (string, error) {
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

func Generate32CharString() (string, error) {
	b := make([]byte, 16) // 16 bytes = 32 hex chars
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
