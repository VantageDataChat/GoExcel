package gospreadsheet

import (
	"testing"
)

func TestColumnIndexToName(t *testing.T) {
	tests := []struct {
		index int
		want  string
		err   bool
	}{
		{0, "A", false},
		{1, "B", false},
		{25, "Z", false},
		{26, "AA", false},
		{27, "AB", false},
		{51, "AZ", false},
		{52, "BA", false},
		{701, "ZZ", false},
		{702, "AAA", false},
		{-1, "", true},
	}

	for _, tt := range tests {
		got, err := ColumnIndexToName(tt.index)
		if tt.err {
			if err == nil {
				t.Errorf("ColumnIndexToName(%d) expected error, got %q", tt.index, got)
			}
			continue
		}
		if err != nil {
			t.Errorf("ColumnIndexToName(%d) unexpected error: %v", tt.index, err)
			continue
		}
		if got != tt.want {
			t.Errorf("ColumnIndexToName(%d) = %q, want %q", tt.index, got, tt.want)
		}
	}
}

func TestColumnNameToIndex(t *testing.T) {
	tests := []struct {
		name string
		want int
		err  bool
	}{
		{"A", 0, false},
		{"B", 1, false},
		{"Z", 25, false},
		{"AA", 26, false},
		{"AB", 27, false},
		{"AZ", 51, false},
		{"BA", 52, false},
		{"ZZ", 701, false},
		{"AAA", 702, false},
		{"a", 0, false},   // case insensitive
		{"aa", 26, false},  // case insensitive
		{"", -1, true},
		{"1A", -1, true},
	}

	for _, tt := range tests {
		got, err := ColumnNameToIndex(tt.name)
		if tt.err {
			if err == nil {
				t.Errorf("ColumnNameToIndex(%q) expected error, got %d", tt.name, got)
			}
			continue
		}
		if err != nil {
			t.Errorf("ColumnNameToIndex(%q) unexpected error: %v", tt.name, err)
			continue
		}
		if got != tt.want {
			t.Errorf("ColumnNameToIndex(%q) = %d, want %d", tt.name, got, tt.want)
		}
	}
}

func TestColumnRoundTrip(t *testing.T) {
	// Test that converting index -> name -> index gives the same result
	for i := 0; i < 1000; i++ {
		name, err := ColumnIndexToName(i)
		if err != nil {
			t.Fatalf("ColumnIndexToName(%d) error: %v", i, err)
		}
		idx, err := ColumnNameToIndex(name)
		if err != nil {
			t.Fatalf("ColumnNameToIndex(%q) error: %v", name, err)
		}
		if idx != i {
			t.Errorf("Round trip failed: %d -> %q -> %d", i, name, idx)
		}
	}
}

func TestParseCellReference(t *testing.T) {
	tests := []struct {
		ref    string
		col    string
		colIdx int
		row    int
		err    bool
	}{
		{"A1", "A", 0, 1, false},
		{"B2", "B", 1, 2, false},
		{"Z100", "Z", 25, 100, false},
		{"AA1", "AA", 26, 1, false},
		{"AB10", "AB", 27, 10, false},
		{"a1", "A", 0, 1, false}, // case insensitive
		{"", "", 0, 0, true},
		{"A", "", 0, 0, true},
		{"1", "", 0, 0, true},
		{"A0", "", 0, 0, true},
		{"1A", "", 0, 0, true},
	}

	for _, tt := range tests {
		cr, err := ParseCellReference(tt.ref)
		if tt.err {
			if err == nil {
				t.Errorf("ParseCellReference(%q) expected error", tt.ref)
			}
			continue
		}
		if err != nil {
			t.Errorf("ParseCellReference(%q) unexpected error: %v", tt.ref, err)
			continue
		}
		if cr.Column != tt.col {
			t.Errorf("ParseCellReference(%q).Column = %q, want %q", tt.ref, cr.Column, tt.col)
		}
		if cr.ColumnIdx != tt.colIdx {
			t.Errorf("ParseCellReference(%q).ColumnIdx = %d, want %d", tt.ref, cr.ColumnIdx, tt.colIdx)
		}
		if cr.Row != tt.row {
			t.Errorf("ParseCellReference(%q).Row = %d, want %d", tt.ref, cr.Row, tt.row)
		}
	}
}

func TestCellName(t *testing.T) {
	tests := []struct {
		row  int
		col  int
		want string
	}{
		{0, 0, "A1"},
		{0, 1, "B1"},
		{1, 0, "A2"},
		{9, 25, "Z10"},
		{0, 26, "AA1"},
	}

	for _, tt := range tests {
		got, err := CellName(tt.row, tt.col)
		if err != nil {
			t.Errorf("CellName(%d, %d) error: %v", tt.row, tt.col, err)
			continue
		}
		if got != tt.want {
			t.Errorf("CellName(%d, %d) = %q, want %q", tt.row, tt.col, got, tt.want)
		}
	}
}

func TestParseRange(t *testing.T) {
	start, end, err := ParseRange("A1:C3")
	if err != nil {
		t.Fatalf("ParseRange error: %v", err)
	}
	if start.Column != "A" || start.Row != 1 {
		t.Errorf("start = %v, want A1", start)
	}
	if end.Column != "C" || end.Row != 3 {
		t.Errorf("end = %v, want C3", end)
	}

	_, _, err = ParseRange("invalid")
	if err == nil {
		t.Error("ParseRange(\"invalid\") expected error")
	}
}
