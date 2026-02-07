package gospreadsheet

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"
)

func TestCSVWriteAndRead(t *testing.T) {
	wb := New()
	ws := wb.GetActiveSheet()
	ws.SetCellValue("A1", "Name")
	ws.SetCellValue("B1", "Score")
	ws.SetCellValue("A2", "Alice")
	ws.SetCellValue("B2", 95.5)
	ws.SetCellValue("A3", "Bob")
	ws.SetCellValue("B3", 87.0)

	tmpDir := t.TempDir()
	filename := filepath.Join(tmpDir, "test.csv")

	writer := NewCSVWriter()
	if err := writer.Save(wb, filename); err != nil {
		t.Fatalf("Save error: %v", err)
	}

	reader := NewCSVReader()
	wb2, err := reader.Open(filename)
	if err != nil {
		t.Fatalf("Open error: %v", err)
	}

	ws2 := wb2.GetActiveSheet()
	v, _ := ws2.GetCellValue("A1")
	if v != "Name" {
		t.Errorf("A1 = %v, want 'Name'", v)
	}

	v, _ = ws2.GetCellValue("B2")
	if v != 95.5 {
		t.Errorf("B2 = %v, want 95.5", v)
	}
}

func TestCSVReadFromString(t *testing.T) {
	csvData := "Name,Age,City\nAlice,30,NYC\nBob,25,LA\n"
	reader := NewCSVReader()
	wb, err := reader.Read(strings.NewReader(csvData))
	if err != nil {
		t.Fatalf("Read error: %v", err)
	}

	ws := wb.GetActiveSheet()
	v, _ := ws.GetCellValue("A1")
	if v != "Name" {
		t.Errorf("A1 = %v, want 'Name'", v)
	}

	v, _ = ws.GetCellValue("B2")
	if v != float64(30) {
		t.Errorf("B2 = %v, want 30", v)
	}

	v, _ = ws.GetCellValue("C3")
	if v != "LA" {
		t.Errorf("C3 = %v, want 'LA'", v)
	}
}

func TestCSVWriteToBuffer(t *testing.T) {
	wb := New()
	ws := wb.GetActiveSheet()
	ws.SetCellValue("A1", "hello")
	ws.SetCellValue("B1", "world")
	ws.SetCellValue("A2", 1)
	ws.SetCellValue("B2", 2)

	var buf bytes.Buffer
	writer := NewCSVWriter()
	if err := writer.Write(wb, &buf); err != nil {
		t.Fatalf("Write error: %v", err)
	}

	expected := "hello,world\n1,2\n"
	if buf.String() != expected {
		t.Errorf("CSV output = %q, want %q", buf.String(), expected)
	}
}

func TestCSVCustomDelimiter(t *testing.T) {
	csvData := "Name;Age;City\nAlice;30;NYC\n"

	reader := NewCSVReader()
	reader.Delimiter = ';'
	wb, err := reader.Read(strings.NewReader(csvData))
	if err != nil {
		t.Fatalf("Read error: %v", err)
	}

	ws := wb.GetActiveSheet()
	v, _ := ws.GetCellValue("C1")
	if v != "City" {
		t.Errorf("C1 = %v, want 'City'", v)
	}

	// Write with custom delimiter
	var buf bytes.Buffer
	writer := NewCSVWriter()
	writer.Delimiter = '\t'
	if err := writer.Write(wb, &buf); err != nil {
		t.Fatalf("Write error: %v", err)
	}

	if !strings.Contains(buf.String(), "\t") {
		t.Error("output should contain tab delimiter")
	}
}

func TestCSVEmptySheet(t *testing.T) {
	wb := New()
	var buf bytes.Buffer
	writer := NewCSVWriter()
	if err := writer.Write(wb, &buf); err != nil {
		t.Fatalf("Write error: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("empty sheet should produce empty CSV, got %q", buf.String())
	}
}

func TestCSVBooleanValues(t *testing.T) {
	csvData := "Active\ntrue\nfalse\n"
	reader := NewCSVReader()
	wb, err := reader.Read(strings.NewReader(csvData))
	if err != nil {
		t.Fatalf("Read error: %v", err)
	}

	ws := wb.GetActiveSheet()
	v, _ := ws.GetCellValue("A2")
	if v != true {
		t.Errorf("A2 = %v (%T), want true", v, v)
	}

	v, _ = ws.GetCellValue("A3")
	if v != false {
		t.Errorf("A3 = %v (%T), want false", v, v)
	}
}

func TestCSVSpecialCharacters(t *testing.T) {
	wb := New()
	ws := wb.GetActiveSheet()
	ws.SetCellValue("A1", "hello, world")  // contains comma
	ws.SetCellValue("A2", `say "hi"`)      // contains quotes
	ws.SetCellValue("A3", "line1\nline2")  // contains newline

	var buf bytes.Buffer
	writer := NewCSVWriter()
	if err := writer.Write(wb, &buf); err != nil {
		t.Fatalf("Write error: %v", err)
	}

	// Read back
	reader := NewCSVReader()
	wb2, err := reader.Read(strings.NewReader(buf.String()))
	if err != nil {
		t.Fatalf("Read error: %v", err)
	}

	ws2 := wb2.GetActiveSheet()
	v, _ := ws2.GetCellValue("A1")
	if v != "hello, world" {
		t.Errorf("A1 = %v, want 'hello, world'", v)
	}

	v, _ = ws2.GetCellValue("A2")
	if v != `say "hi"` {
		t.Errorf("A2 = %v, want 'say \"hi\"'", v)
	}

	v, _ = ws2.GetCellValue("A3")
	if v != "line1\nline2" {
		t.Errorf("A3 = %v, want 'line1\\nline2'", v)
	}
}

func TestCSVMultipleSheetSelection(t *testing.T) {
	wb := New()
	ws1 := wb.GetActiveSheet()
	ws1.SetCellValue("A1", "Sheet1")

	ws2, _ := wb.AddSheet("Sheet2")
	ws2.SetCellValue("A1", "Sheet2")

	var buf bytes.Buffer
	writer := NewCSVWriter()
	writer.SheetIndex = 1 // Write Sheet2
	if err := writer.Write(wb, &buf); err != nil {
		t.Fatalf("Write error: %v", err)
	}

	if !strings.Contains(buf.String(), "Sheet2") {
		t.Errorf("output should contain 'Sheet2', got %q", buf.String())
	}
	if strings.Contains(buf.String(), "Sheet1") {
		t.Errorf("output should not contain 'Sheet1', got %q", buf.String())
	}
}
