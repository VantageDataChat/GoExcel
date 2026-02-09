package gospreadsheet

import "strings"

// RichTextRun represents a segment of rich text with its own formatting.
type RichTextRun struct {
	Text string
	Font *Font
}

// RichText represents a cell value with mixed formatting.
type RichText struct {
	Runs []RichTextRun
}

// NewRichText creates a new empty rich text.
func NewRichText() *RichText {
	return &RichText{}
}

// AddRun adds a text run with optional font formatting.
func (rt *RichText) AddRun(text string, font *Font) *RichText {
	rt.Runs = append(rt.Runs, RichTextRun{Text: text, Font: font})
	return rt
}

// PlainText returns the concatenated plain text of all runs.
func (rt *RichText) PlainText() string {
	var b strings.Builder
	for _, run := range rt.Runs {
		b.WriteString(run.Text)
	}
	return b.String()
}
