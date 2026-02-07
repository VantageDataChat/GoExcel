package gospreadsheet

import (
	"errors"
	"fmt"
)

// DocumentProperties holds metadata about the spreadsheet.
type DocumentProperties struct {
	Creator        string
	LastModifiedBy string
	Title          string
	Subject        string
	Description    string
	Keywords       string
	Category       string
}

// Workbook represents a spreadsheet workbook containing one or more worksheets.
type Workbook struct {
	worksheets  []*Worksheet
	activeSheet int
	Properties  DocumentProperties
	namedRanges map[string]string // name -> "Sheet1!A1:B2"
}

// New creates a new workbook with a single default worksheet.
func New() *Workbook {
	wb := &Workbook{
		worksheets:  make([]*Worksheet, 0),
		namedRanges: make(map[string]string),
	}
	ws := NewWorksheet("Sheet1")
	ws.parent = wb
	ws.index = 0
	wb.worksheets = append(wb.worksheets, ws)
	return wb
}

// NewEmpty creates a new workbook with no worksheets.
func NewEmpty() *Workbook {
	return &Workbook{
		worksheets:  make([]*Worksheet, 0),
		namedRanges: make(map[string]string),
	}
}

// AddSheet adds a new worksheet with the given title.
func (wb *Workbook) AddSheet(title string) (*Worksheet, error) {
	for _, ws := range wb.worksheets {
		if ws.title == title {
			return nil, fmt.Errorf("worksheet with title %q already exists", title)
		}
	}
	ws := NewWorksheet(title)
	ws.parent = wb
	ws.index = len(wb.worksheets)
	wb.worksheets = append(wb.worksheets, ws)
	return ws, nil
}

// GetSheet returns the worksheet at the given index (0-based).
func (wb *Workbook) GetSheet(index int) (*Worksheet, error) {
	if index < 0 || index >= len(wb.worksheets) {
		return nil, fmt.Errorf("sheet index %d out of range (0-%d)", index, len(wb.worksheets)-1)
	}
	return wb.worksheets[index], nil
}

// GetSheetByName returns the worksheet with the given title.
func (wb *Workbook) GetSheetByName(title string) (*Worksheet, error) {
	for _, ws := range wb.worksheets {
		if ws.title == title {
			return ws, nil
		}
	}
	return nil, fmt.Errorf("worksheet %q not found", title)
}

// SheetCount returns the number of worksheets.
func (wb *Workbook) SheetCount() int {
	return len(wb.worksheets)
}

// GetSheetNames returns the titles of all worksheets.
func (wb *Workbook) GetSheetNames() []string {
	names := make([]string, len(wb.worksheets))
	for i, ws := range wb.worksheets {
		names[i] = ws.title
	}
	return names
}

// RemoveSheet removes the worksheet at the given index.
func (wb *Workbook) RemoveSheet(index int) error {
	if index < 0 || index >= len(wb.worksheets) {
		return fmt.Errorf("sheet index %d out of range", index)
	}
	if len(wb.worksheets) <= 1 {
		return errors.New("cannot remove the last worksheet")
	}
	wb.worksheets = append(wb.worksheets[:index], wb.worksheets[index+1:]...)
	// Re-index
	for i, ws := range wb.worksheets {
		ws.index = i
	}
	if wb.activeSheet >= len(wb.worksheets) {
		wb.activeSheet = len(wb.worksheets) - 1
	}
	return nil
}

// GetActiveSheet returns the active worksheet.
func (wb *Workbook) GetActiveSheet() *Worksheet {
	if len(wb.worksheets) == 0 {
		return nil
	}
	return wb.worksheets[wb.activeSheet]
}

// SetActiveSheet sets the active worksheet by index.
func (wb *Workbook) SetActiveSheet(index int) error {
	if index < 0 || index >= len(wb.worksheets) {
		return fmt.Errorf("sheet index %d out of range", index)
	}
	wb.activeSheet = index
	return nil
}

// AddNamedRange adds a named range to the workbook.
func (wb *Workbook) AddNamedRange(name, reference string) {
	wb.namedRanges[name] = reference
}

// GetNamedRange returns the reference for a named range.
func (wb *Workbook) GetNamedRange(name string) (string, error) {
	ref, ok := wb.namedRanges[name]
	if !ok {
		return "", fmt.Errorf("named range %q not found", name)
	}
	return ref, nil
}

// GetNamedRanges returns all named ranges.
func (wb *Workbook) GetNamedRanges() map[string]string {
	result := make(map[string]string, len(wb.namedRanges))
	for k, v := range wb.namedRanges {
		result[k] = v
	}
	return result
}
