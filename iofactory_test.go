package gospreadsheet

import (
	"path/filepath"
	"testing"
)

func TestOpenFileXLSX(t *testing.T) {
	// Create a test file first
	wb := New()
	ws := wb.GetActiveSheet()
	ws.SetCellValue("A1", "IOFactory Test")

	tmpDir := t.TempDir()
	filename := filepath.Join(tmpDir, "test.xlsx")
	if err := SaveFile(wb, filename); err != nil {
		t.Fatalf("SaveFile error: %v", err)
	}

	wb2, err := OpenFile(filename)
	if err != nil {
		t.Fatalf("OpenFile error: %v", err)
	}

	v, _ := wb2.GetActiveSheet().GetCellValue("A1")
	if v != "IOFactory Test" {
		t.Errorf("A1 = %v, want 'IOFactory Test'", v)
	}
}

func TestOpenFileCSV(t *testing.T) {
	wb := New()
	ws := wb.GetActiveSheet()
	ws.SetCellValue("A1", "CSV Test")

	tmpDir := t.TempDir()
	filename := filepath.Join(tmpDir, "test.csv")
	if err := SaveFile(wb, filename); err != nil {
		t.Fatalf("SaveFile error: %v", err)
	}

	wb2, err := OpenFile(filename)
	if err != nil {
		t.Fatalf("OpenFile error: %v", err)
	}

	v, _ := wb2.GetActiveSheet().GetCellValue("A1")
	if v != "CSV Test" {
		t.Errorf("A1 = %v, want 'CSV Test'", v)
	}
}

func TestOpenFileUnsupported(t *testing.T) {
	_, err := OpenFile("test.xyz")
	if err == nil {
		t.Error("expected error for unsupported format")
	}
}

func TestSaveFileUnsupported(t *testing.T) {
	wb := New()
	err := SaveFile(wb, "test.xyz")
	if err == nil {
		t.Error("expected error for unsupported format")
	}
}
