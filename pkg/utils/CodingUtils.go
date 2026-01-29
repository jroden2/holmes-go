package utils

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"io"

	"github.com/jroden2/holmes-go/pkg/domain"
)

func CreateMagicLink(p domain.DiffPayload) (string, error) {
	b, err := json.Marshal(p)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	_, err = gz.Write(b)
	if err != nil {
		return "", err
	}
	err = gz.Close()
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(buf.Bytes()), nil
}

func DecodeMagicLink(raw string) (*domain.DiffPayload, error) {
	compressed, err := base64.URLEncoding.DecodeString(raw)
	if err != nil {
		return nil, err
	}

	gz, _ := gzip.NewReader(bytes.NewReader(compressed))
	defer gz.Close()
	data, _ := io.ReadAll(gz)

	var payload domain.DiffPayload
	if err = json.Unmarshal(data, &payload); err != nil {
		return nil, err
	}

	return &payload, nil
}
