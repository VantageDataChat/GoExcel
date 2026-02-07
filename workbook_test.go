package gospreadsheet

import (
	"testing"
)

func TestNewWorkbook(t *testing.T) {
	wb := New()
	if wb.SheetCount() != 1 {
		t.Errorf("new workbook should have 1 sheet, got %d", wb.SheetCount())
	}
	names := wb.GetSheetNames()
	if len(names) != 1 || names[0] != "Sheet1" {
		t.Errorf("sheet names = %v, want [Sheet1]", names)
	}
}

func TestNewEmptyWorkbook(t *testing.T) {
	wb := NewEmpty()
	if wb.SheetCount() != 0 {
		t.Errorf("empty workbook should have 0 sheets, got %d", wb.SheetCount())
	}
}

func TestWorkbookAddSheet(t *testing.T) {
	wb := New()
	ws, err := wb.AddSheet("Sheet2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ws.Title() != "Sheet2" {
		t.Errorf("sheet title = %q, want 'Sheet2'", ws.Title())
	}
	if wb.SheetCount() != 2 {
		t.Errorf("sheet count = %d, want 2", wb.SheetCount())
	}

	// Duplicate name
	_, err = wb.AddSheet("Sheet2")
	if err == nil {
		t.Error("expected error for duplicate sheet name")
	}
}

func TestWorkbookGetSheet(t *testing.T) {
	wb := New()
	ws, err := wb.GetSheet(0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ws.Title() != "Sheet1" {
		t.Errorf("sheet title = %q, want 'Sheet1'", ws.Title())
	}

	_, err = wb.GetSheet(1)
	if err == nil {
		t.Error("expected error for out of range index")
	}

	_, err = wb.GetSheet(-1)
	if err == nil {
		t.Error("expected error for negative index")
	}
}

func TestWorkbookGetSheetByName(t *testing.T) {
	wb := New()
	wb.AddSheet("MySheet")

	ws, err := wb.GetSheetByName("MySheet")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ws.Title() != "MySheet" {
		t.Errorf("sheet title = %q, want 'MySheet'", ws.Title())
	}

	_, err = wb.GetSheetByName("NonExistent")
	if err == nil {
		t.Error("expected error for non-existent sheet")
	}
}

func TestWorkbookRemoveSheet(t *testing.T) {
	wb := New()
	wb.AddSheet("Sheet2")
	wb.AddSheet("Sheet3")

	if err := wb.RemoveSheet(1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if wb.SheetCount() != 2 {
		t.Errorf("sheet count = %d, want 2", wb.SheetCount())
	}
	names := wb.GetSheetNames()
	if names[0] != "Sheet1" || names[1] != "Sheet3" {
		t.Errorf("sheet names = %v, want [Sheet1, Sheet3]", names)
	}

	// Cannot remove last sheet
	wb.RemoveSheet(0)
	err := wb.RemoveSheet(0)
	if err == nil {
		t.Error("expected error when removing last sheet")
	}

	// Out of range
	err = wb.RemoveSheet(99)
	if err == nil {
		t.Error("expected error for out of range index")
	}
}

func TestWorkbookActiveSheet(t *testing.T) {
	wb := New()
	wb.AddSheet("Sheet2")

	ws := wb.GetActiveSheet()
	if ws.Title() != "Sheet1" {
		t.Errorf("active sheet = %q, want 'Sheet1'", ws.Title())
	}

	if err := wb.SetActiveSheet(1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	ws = wb.GetActiveSheet()
	if ws.Title() != "Sheet2" {
		t.Errorf("active sheet = %q, want 'Sheet2'", ws.Title())
	}

	err := wb.SetActiveSheet(99)
	if err == nil {
		t.Error("expected error for out of range index")
	}
}

func TestWorkbookProperties(t *testing.T) {
	wb := New()
	wb.Properties.Creator = "Test User"
	wb.Properties.Title = "Test Document"
	wb.Properties.Subject = "Testing"
	wb.Properties.Description = "A test document"
	wb.Properties.Keywords = "test, go"
	wb.Properties.Category = "Test"

	if wb.Properties.Creator != "Test User" {
		t.Errorf("Creator = %q, want 'Test User'", wb.Properties.Creator)
	}
	if wb.Properties.Title != "Test Document" {
		t.Errorf("Title = %q, want 'Test Document'", wb.Properties.Title)
	}
}

func TestWorkbookNamedRanges(t *testing.T) {
	wb := New()
	wb.AddNamedRange("MyRange", "Sheet1!A1:B10")

	ref, err := wb.GetNamedRange("MyRange")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ref != "Sheet1!A1:B10" {
		t.Errorf("named range = %q, want 'Sheet1!A1:B10'", ref)
	}

	_, err = wb.GetNamedRange("NonExistent")
	if err == nil {
		t.Error("expected error for non-existent named range")
	}

	ranges := wb.GetNamedRanges()
	if len(ranges) != 1 {
		t.Errorf("named ranges count = %d, want 1", len(ranges))
	}
}

func TestWorkbookRemoveSheetAdjustsActiveSheet(t *testing.T) {
	wb := New()
	wb.AddSheet("Sheet2")
	wb.SetActiveSheet(1)

	wb.RemoveSheet(1)
	ws := wb.GetActiveSheet()
	if ws.Title() != "Sheet1" {
		t.Errorf("active sheet after removal = %q, want 'Sheet1'", ws.Title())
	}
}
