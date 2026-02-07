package gospreadsheet

import (
	"testing"
)

func TestWorksheetTitle(t *testing.T) {
	ws := NewWorksheet("Test")
	if ws.Title() != "Test" {
		t.Errorf("Title() = %q, want %q", ws.Title(), "Test")
	}
	ws.SetTitle("New Title")
	if ws.Title() != "New Title" {
		t.Errorf("Title() = %q, want %q", ws.Title(), "New Title")
	}
}

func TestWorksheetGetCell(t *testing.T) {
	ws := NewWorksheet("Test")

	// Getting a cell should create it
	c := ws.GetCell(0, 0)
	if c == nil {
		t.Fatal("GetCell returned nil")
	}
	if c.Type != CellTypeEmpty {
		t.Errorf("new cell type = %d, want CellTypeEmpty", c.Type)
	}

	// Getting the same cell again should return the same instance
	c2 := ws.GetCell(0, 0)
	if c != c2 {
		t.Error("GetCell should return the same cell instance")
	}
}

func TestWorksheetGetCellByName(t *testing.T) {
	ws := NewWorksheet("Test")

	c, err := ws.GetCellByName("A1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Row() != 0 || c.Col() != 0 {
		t.Errorf("A1 should be (0,0), got (%d,%d)", c.Row(), c.Col())
	}

	c, err = ws.GetCellByName("C5")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Row() != 4 || c.Col() != 2 {
		t.Errorf("C5 should be (4,2), got (%d,%d)", c.Row(), c.Col())
	}

	_, err = ws.GetCellByName("invalid")
	if err == nil {
		t.Error("expected error for invalid reference")
	}
}

func TestWorksheetSetCellValue(t *testing.T) {
	ws := NewWorksheet("Test")

	if err := ws.SetCellValue("A1", "Hello"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := ws.SetCellValue("B1", 42); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := ws.SetCellValue("C1", true); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	v, _ := ws.GetCellValue("A1")
	if v != "Hello" {
		t.Errorf("A1 = %v, want 'Hello'", v)
	}

	v, _ = ws.GetCellValue("B1")
	if v != float64(42) {
		t.Errorf("B1 = %v, want 42", v)
	}

	v, _ = ws.GetCellValue("C1")
	if v != true {
		t.Errorf("C1 = %v, want true", v)
	}
}

func TestWorksheetFormula(t *testing.T) {
	ws := NewWorksheet("Test")
	ws.SetCellValue("A1", 10)
	ws.SetCellValue("A2", 20)
	ws.SetCellFormula("A3", "SUM(A1:A2)")

	c, _ := ws.GetCellByName("A3")
	if c.Type != CellTypeFormula {
		t.Errorf("expected formula type, got %d", c.Type)
	}
	if c.Formula != "SUM(A1:A2)" {
		t.Errorf("formula = %q, want 'SUM(A1:A2)'", c.Formula)
	}
}

func TestWorksheetMergeCells(t *testing.T) {
	ws := NewWorksheet("Test")
	if err := ws.MergeCells("A1:C3"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	merges := ws.GetMergeCells()
	if len(merges) != 1 {
		t.Fatalf("expected 1 merge, got %d", len(merges))
	}
	if merges[0].StartRow != 0 || merges[0].StartCol != 0 {
		t.Errorf("merge start = (%d,%d), want (0,0)", merges[0].StartRow, merges[0].StartCol)
	}
	if merges[0].EndRow != 2 || merges[0].EndCol != 2 {
		t.Errorf("merge end = (%d,%d), want (2,2)", merges[0].EndRow, merges[0].EndCol)
	}

	err := ws.MergeCells("invalid")
	if err == nil {
		t.Error("expected error for invalid range")
	}
}

func TestWorksheetColumnWidth(t *testing.T) {
	ws := NewWorksheet("Test")

	// Default width
	if w := ws.GetColumnWidth(0); w != 8.43 {
		t.Errorf("default width = %f, want 8.43", w)
	}

	ws.SetColumnWidth(0, 20.0)
	if w := ws.GetColumnWidth(0); w != 20.0 {
		t.Errorf("width = %f, want 20.0", w)
	}
}

func TestWorksheetRowHeight(t *testing.T) {
	ws := NewWorksheet("Test")

	// Default height
	if h := ws.GetRowHeight(0); h != 15.0 {
		t.Errorf("default height = %f, want 15.0", h)
	}

	ws.SetRowHeight(0, 30.0)
	if h := ws.GetRowHeight(0); h != 30.0 {
		t.Errorf("height = %f, want 30.0", h)
	}
}

func TestWorksheetFreezePane(t *testing.T) {
	ws := NewWorksheet("Test")

	if ws.GetFreezePane() != nil {
		t.Error("freeze pane should be nil initially")
	}

	if err := ws.FreezePane("B2"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	fp := ws.GetFreezePane()
	if fp == nil {
		t.Fatal("freeze pane should not be nil")
	}
	if fp.Column != "B" || fp.Row != 2 {
		t.Errorf("freeze pane = %s%d, want B2", fp.Column, fp.Row)
	}
}

func TestWorksheetDimensions(t *testing.T) {
	ws := NewWorksheet("Test")

	// Empty sheet
	_, _, _, _, err := ws.Dimensions()
	if err == nil {
		t.Error("expected error for empty sheet")
	}

	ws.SetCellValue("B2", "hello")
	ws.SetCellValue("D5", "world")

	minRow, minCol, maxRow, maxCol, err := ws.Dimensions()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if minRow != 1 || minCol != 1 {
		t.Errorf("min = (%d,%d), want (1,1)", minRow, minCol)
	}
	if maxRow != 4 || maxCol != 3 {
		t.Errorf("max = (%d,%d), want (4,3)", maxRow, maxCol)
	}
}

func TestWorksheetRowIterator(t *testing.T) {
	ws := NewWorksheet("Test")
	ws.SetCellValue("A1", "a")
	ws.SetCellValue("B1", "b")
	ws.SetCellValue("A2", "c")
	ws.SetCellValue("B2", "d")

	rows, err := ws.RowIterator()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rows) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(rows))
	}
	if len(rows[0]) != 2 {
		t.Fatalf("expected 2 cols, got %d", len(rows[0]))
	}
	if rows[0][0].GetStringValue() != "a" {
		t.Errorf("A1 = %q, want 'a'", rows[0][0].GetStringValue())
	}
	if rows[1][1].GetStringValue() != "d" {
		t.Errorf("B2 = %q, want 'd'", rows[1][1].GetStringValue())
	}
}

func TestWorksheetCellCount(t *testing.T) {
	ws := NewWorksheet("Test")
	if ws.CellCount() != 0 {
		t.Errorf("empty sheet cell count = %d, want 0", ws.CellCount())
	}

	ws.SetCellValue("A1", "hello")
	ws.SetCellValue("B1", 42)
	if ws.CellCount() != 2 {
		t.Errorf("cell count = %d, want 2", ws.CellCount())
	}
}

func TestWorksheetAllCells(t *testing.T) {
	ws := NewWorksheet("Test")
	ws.SetCellValue("C1", "third")
	ws.SetCellValue("A1", "first")
	ws.SetCellValue("B1", "second")

	cells := ws.AllCells()
	if len(cells) != 3 {
		t.Fatalf("expected 3 cells, got %d", len(cells))
	}
	// Should be sorted by column
	if cells[0].GetStringValue() != "first" {
		t.Errorf("first cell = %q, want 'first'", cells[0].GetStringValue())
	}
	if cells[1].GetStringValue() != "second" {
		t.Errorf("second cell = %q, want 'second'", cells[1].GetStringValue())
	}
	if cells[2].GetStringValue() != "third" {
		t.Errorf("third cell = %q, want 'third'", cells[2].GetStringValue())
	}
}
