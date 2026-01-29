package domain

import (
	"html/template"
)

type PageData struct {
	A, B                 string
	IgnoreWS, IgnoreCase bool
	Mode                 string // "text" | "json" | "xml"

	ExactMatch, NormalizedMatch bool

	ALen, BLen int

	AHash, BHash string

	LineDiff []LineDiffRow
	Error    string
}

type LineDiffRow struct {
	LineNum      int
	A, B         string
	AHTML, BHTML template.HTML
	Status       string
}

type DiffPayload struct {
	ID       string `json:"id"`
	ShortID  string `json:"short_id"`
	Original string `json:"a"`
	New      string `json:"b"`
	Format   string `json:"f"` // Json, XML, Text
}

func NewDiffPayload(original, new string) DiffPayload {
	return DiffPayload{"", "", original, new, "text"}
}

func (dp *DiffPayload) GetID() string {
	return dp.ShortID
}
