package gospreadsheet

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// ColumnIndexToName converts a 0-based column index to a column name (e.g., 0 -> "A", 25 -> "Z", 26 -> "AA").
func ColumnIndexToName(index int) (string, error) {
	if index < 0 {
		return "", errors.New("column index must be non-negative")
	}
	result := ""
	for {
		result = string(rune('A'+index%26)) + result
		index = index/26 - 1
		if index < 0 {
			break
		}
	}
	return result, nil
}

// ColumnNameToIndex converts a column name to a 0-based column index (e.g., "A" -> 0, "Z" -> 25, "AA" -> 26).
func ColumnNameToIndex(name string) (int, error) {
	name = strings.ToUpper(strings.TrimSpace(name))
	if name == "" {
		return -1, errors.New("column name cannot be empty")
	}
	result := 0
	for _, ch := range name {
		if ch < 'A' || ch > 'Z' {
			return -1, fmt.Errorf("invalid column name character: %c", ch)
		}
		result = result*26 + int(ch-'A') + 1
	}
	return result - 1, nil
}

// CellReference represents a parsed cell reference like "A1", "B2", etc.
type CellReference struct {
	Column    string
	ColumnIdx int
	Row       int // 1-based row number
}

// ParseCellReference parses a cell reference string (e.g., "A1") into its components.
func ParseCellReference(ref string) (*CellReference, error) {
	ref = strings.TrimSpace(ref)
	if ref == "" {
		return nil, errors.New("cell reference cannot be empty")
	}

	// Split into column letters and row number
	colPart := ""
	rowPart := ""
	for _, ch := range ref {
		if unicode.IsLetter(ch) {
			if rowPart != "" {
				return nil, fmt.Errorf("invalid cell reference: %s", ref)
			}
			colPart += string(ch)
		} else if unicode.IsDigit(ch) {
			rowPart += string(ch)
		} else {
			return nil, fmt.Errorf("invalid character in cell reference: %c", ch)
		}
	}

	if colPart == "" || rowPart == "" {
		return nil, fmt.Errorf("invalid cell reference: %s", ref)
	}

	colIdx, err := ColumnNameToIndex(colPart)
	if err != nil {
		return nil, err
	}

	row, err := strconv.Atoi(rowPart)
	if err != nil || row < 1 {
		return nil, fmt.Errorf("invalid row number in cell reference: %s", ref)
	}

	return &CellReference{
		Column:    strings.ToUpper(colPart),
		ColumnIdx: colIdx,
		Row:       row,
	}, nil
}

// CellName returns the cell reference string (e.g., "A1") for a 0-based row and column.
func CellName(row, col int) (string, error) {
	colName, err := ColumnIndexToName(col)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s%d", colName, row+1), nil
}

// ParseRange parses a range string like "A1:C3" into two CellReferences.
func ParseRange(rangeStr string) (*CellReference, *CellReference, error) {
	parts := strings.Split(rangeStr, ":")
	if len(parts) != 2 {
		return nil, nil, fmt.Errorf("invalid range: %s", rangeStr)
	}
	start, err := ParseCellReference(parts[0])
	if err != nil {
		return nil, nil, fmt.Errorf("invalid range start: %w", err)
	}
	end, err := ParseCellReference(parts[1])
	if err != nil {
		return nil, nil, fmt.Errorf("invalid range end: %w", err)
	}
	return start, end, nil
}
