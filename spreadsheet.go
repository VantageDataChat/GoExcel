// Package gospreadsheet provides a pure Go library for reading and writing
// spreadsheet files, inspired by PHPOffice/PhpSpreadsheet.
// It supports XLSX and CSV formats with an in-memory spreadsheet model.
package gospreadsheet

import (
	"errors"
	"fmt"
	"time"
)

// CellType represents the data type of a cell value.
type CellType int

const (
	CellTypeEmpty CellType = iota
	CellTypeString
	CellTypeNumeric
	CellTypeBool
	CellTypeFormula
	CellTypeDate
	CellTypeError
)

// Cell represents a single cell in a worksheet.
type Cell struct {
	Value     interface{}
	Type      CellType
	Formula   string
	Style     *Style
	Hyperlink *Hyperlink
	Comment   *Comment
	RichText  *RichText
	row       int
	col       int
}

// NewCell creates a new empty cell.
func NewCell(row, col int) *Cell {
	return &Cell{
		Type: CellTypeEmpty,
		row:  row,
		col:  col,
	}
}

// SetValue sets the cell value and auto-detects the type.
func (c *Cell) SetValue(v interface{}) *Cell {
	if v == nil {
		c.Value = nil
		c.Type = CellTypeEmpty
		return c
	}
	switch val := v.(type) {
	case string:
		c.Value = val
		c.Type = CellTypeString
	case int:
		c.Value = float64(val)
		c.Type = CellTypeNumeric
	case int32:
		c.Value = float64(val)
		c.Type = CellTypeNumeric
	case int64:
		c.Value = float64(val)
		c.Type = CellTypeNumeric
	case float32:
		c.Value = float64(val)
		c.Type = CellTypeNumeric
	case float64:
		c.Value = val
		c.Type = CellTypeNumeric
	case bool:
		c.Value = val
		c.Type = CellTypeBool
	case time.Time:
		c.Value = val
		c.Type = CellTypeDate
	default:
		c.Value = v
		c.Type = CellTypeString
	}
	return c
}

// SetFormula sets a formula for the cell.
func (c *Cell) SetFormula(formula string) *Cell {
	c.Formula = formula
	c.Type = CellTypeFormula
	return c
}

// SetStyle sets the style for the cell.
func (c *Cell) SetStyle(s *Style) *Cell {
	c.Style = s
	return c
}

// GetStringValue returns the cell value as a string.
func (c *Cell) GetStringValue() string {
	if c.Value == nil {
		return ""
	}
	switch v := c.Value.(type) {
	case string:
		return v
	case float64:
		if v == float64(int64(v)) {
			return fmt.Sprintf("%d", int64(v))
		}
		return fmt.Sprintf("%g", v)
	case bool:
		if v {
			return "TRUE"
		}
		return "FALSE"
	case time.Time:
		return v.Format("2006-01-02 15:04:05")
	default:
		return fmt.Sprintf("%v", v)
	}
}

// GetNumericValue returns the cell value as float64.
func (c *Cell) GetNumericValue() (float64, error) {
	if c.Type == CellTypeNumeric {
		if v, ok := c.Value.(float64); ok {
			return v, nil
		}
	}
	return 0, errors.New("cell is not numeric")
}

// GetBoolValue returns the cell value as bool.
func (c *Cell) GetBoolValue() (bool, error) {
	if c.Type == CellTypeBool {
		if v, ok := c.Value.(bool); ok {
			return v, nil
		}
	}
	return false, errors.New("cell is not boolean")
}

// GetDateValue returns the cell value as time.Time.
func (c *Cell) GetDateValue() (time.Time, error) {
	if c.Type == CellTypeDate {
		if v, ok := c.Value.(time.Time); ok {
			return v, nil
		}
	}
	return time.Time{}, errors.New("cell is not a date")
}

// Row returns the cell's row index (0-based).
func (c *Cell) Row() int { return c.row }

// Col returns the cell's column index (0-based).
func (c *Cell) Col() int { return c.col }

// SetHyperlink sets a hyperlink on the cell.
func (c *Cell) SetHyperlink(h *Hyperlink) *Cell {
	c.Hyperlink = h
	return c
}

// SetComment sets a comment on the cell.
func (c *Cell) SetComment(comment *Comment) *Cell {
	c.Comment = comment
	return c
}

// SetRichText sets rich text content on the cell.
func (c *Cell) SetRichText(rt *RichText) *Cell {
	c.RichText = rt
	c.Type = CellTypeString
	if rt != nil {
		c.Value = rt.PlainText()
	}
	return c
}
