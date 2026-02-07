package gospreadsheet

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestXLSXWriteAndRead(t *testing.T) {
	// Create a workbook with data
	wb := New()
	wb.Properties.Creator = "GoSpreadsheet Test"
	wb.Properties.Title = "Test Document"

	ws := wb.GetActiveSheet()
	ws.SetCellValue("A1", "Name")
	ws.SetCellValue("B1", "Age")
	ws.SetCellValue("C1", "Active")
	ws.SetCellValue("A2", "Alice")
	ws.SetCellValue("B2", 30)
	ws.SetCellValue("C2", true)
	ws.SetCellValue("A3", "Bob")
	ws.SetCellValue("B3", 25)
	ws.SetCellValue("C3", false)

	// Write to temp file
	tmpDir := t.TempDir()
	filename := filepath.Join(tmpDir, "test.xlsx")

	writer := NewXLSXWriter()
	if err := writer.Save(wb, filename); err != nil {
		t.Fatalf("Save error: %v", err)
	}

	// Verify file exists
	info, err := os.Stat(filename)
	if err != nil {
		t.Fatalf("file not created: %v", err)
	}
	if info.Size() == 0 {
		t.Fatal("file is empty")
	}

	// Read back
	reader := NewXLSXReader()
	wb2, err := reader.Open(filename)
	if err != nil {
		t.Fatalf("Open error: %v", err)
	}

	if wb2.SheetCount() != 1 {
		t.Errorf("sheet count = %d, want 1", wb2.SheetCount())
	}

	ws2 := wb2.GetActiveSheet()
	if ws2 == nil {
		t.Fatal("active sheet is nil")
	}

	// Verify data
	v, _ := ws2.GetCellValue("A1")
	if v != "Name" {
		t.Errorf("A1 = %v, want 'Name'", v)
	}

	v, _ = ws2.GetCellValue("A2")
	if v != "Alice" {
		t.Errorf("A2 = %v, want 'Alice'", v)
	}

	v, _ = ws2.GetCellValue("B2")
	if v != float64(30) {
		t.Errorf("B2 = %v, want 30", v)
	}

	v, _ = ws2.GetCellValue("C2")
	if v != true {
		t.Errorf("C2 = %v, want true", v)
	}

	v, _ = ws2.GetCellValue("C3")
	if v != false {
		t.Errorf("C3 = %v, want false", v)
	}
}

func TestXLSXMultipleSheets(t *testing.T) {
	wb := New()
	ws1 := wb.GetActiveSheet()
	ws1.SetCellValue("A1", "Sheet1 Data")

	ws2, _ := wb.AddSheet("Sheet2")
	ws2.SetCellValue("A1", "Sheet2 Data")

	ws3, _ := wb.AddSheet("Sheet3")
	ws3.SetCellValue("A1", "Sheet3 Data")

	tmpDir := t.TempDir()
	filename := filepath.Join(tmpDir, "multi.xlsx")

	if err := NewXLSXWriter().Save(wb, filename); err != nil {
		t.Fatalf("Save error: %v", err)
	}

	wb2, err := NewXLSXReader().Open(filename)
	if err != nil {
		t.Fatalf("Open error: %v", err)
	}

	if wb2.SheetCount() != 3 {
		t.Errorf("sheet count = %d, want 3", wb2.SheetCount())
	}

	names := wb2.GetSheetNames()
	expected := []string{"Sheet1", "Sheet2", "Sheet3"}
	for i, name := range names {
		if name != expected[i] {
			t.Errorf("sheet[%d] = %q, want %q", i, name, expected[i])
		}
	}

	ws, _ := wb2.GetSheetByName("Sheet2")
	v, _ := ws.GetCellValue("A1")
	if v != "Sheet2 Data" {
		t.Errorf("Sheet2 A1 = %v, want 'Sheet2 Data'", v)
	}
}

func TestXLSXFormula(t *testing.T) {
	wb := New()
	ws := wb.GetActiveSheet()
	ws.SetCellValue("A1", 10)
	ws.SetCellValue("A2", 20)
	ws.SetCellFormula("A3", "SUM(A1:A2)")

	tmpDir := t.TempDir()
	filename := filepath.Join(tmpDir, "formula.xlsx")

	if err := NewXLSXWriter().Save(wb, filename); err != nil {
		t.Fatalf("Save error: %v", err)
	}

	wb2, err := NewXLSXReader().Open(filename)
	if err != nil {
		t.Fatalf("Open error: %v", err)
	}

	ws2 := wb2.GetActiveSheet()
	c, _ := ws2.GetCellByName("A3")
	if c.Type != CellTypeFormula {
		t.Errorf("A3 type = %d, want CellTypeFormula", c.Type)
	}
	if c.Formula != "SUM(A1:A2)" {
		t.Errorf("A3 formula = %q, want 'SUM(A1:A2)'", c.Formula)
	}
}

func TestXLSXMergeCells(t *testing.T) {
	wb := New()
	ws := wb.GetActiveSheet()
	ws.SetCellValue("A1", "Merged")
	ws.MergeCells("A1:C1")

	tmpDir := t.TempDir()
	filename := filepath.Join(tmpDir, "merge.xlsx")

	if err := NewXLSXWriter().Save(wb, filename); err != nil {
		t.Fatalf("Save error: %v", err)
	}

	wb2, err := NewXLSXReader().Open(filename)
	if err != nil {
		t.Fatalf("Open error: %v", err)
	}

	ws2 := wb2.GetActiveSheet()
	merges := ws2.GetMergeCells()
	if len(merges) != 1 {
		t.Fatalf("merge count = %d, want 1", len(merges))
	}
	if merges[0].StartCol != 0 || merges[0].EndCol != 2 {
		t.Errorf("merge cols = (%d,%d), want (0,2)", merges[0].StartCol, merges[0].EndCol)
	}
}

func TestXLSXWriteToBuffer(t *testing.T) {
	wb := New()
	ws := wb.GetActiveSheet()
	ws.SetCellValue("A1", "Buffer Test")

	var buf bytes.Buffer
	if err := NewXLSXWriter().Write(wb, &buf); err != nil {
		t.Fatalf("Write error: %v", err)
	}

	if buf.Len() == 0 {
		t.Fatal("buffer is empty")
	}

	// Read back from buffer
	reader := bytes.NewReader(buf.Bytes())
	wb2, err := NewXLSXReader().Read(reader, int64(buf.Len()))
	if err != nil {
		t.Fatalf("Read error: %v", err)
	}

	v, _ := wb2.GetActiveSheet().GetCellValue("A1")
	if v != "Buffer Test" {
		t.Errorf("A1 = %v, want 'Buffer Test'", v)
	}
}

func TestXLSXSpecialCharacters(t *testing.T) {
	wb := New()
	ws := wb.GetActiveSheet()
	ws.SetCellValue("A1", `Hello <World> & "Quotes"`)
	ws.SetCellValue("A2", "日本語テスト")
	ws.SetCellValue("A3", "Ñoño café")

	tmpDir := t.TempDir()
	filename := filepath.Join(tmpDir, "special.xlsx")

	if err := NewXLSXWriter().Save(wb, filename); err != nil {
		t.Fatalf("Save error: %v", err)
	}

	wb2, err := NewXLSXReader().Open(filename)
	if err != nil {
		t.Fatalf("Open error: %v", err)
	}

	ws2 := wb2.GetActiveSheet()

	v, _ := ws2.GetCellValue("A1")
	if v != `Hello <World> & "Quotes"` {
		t.Errorf("A1 = %v, want special chars", v)
	}

	v, _ = ws2.GetCellValue("A2")
	if v != "日本語テスト" {
		t.Errorf("A2 = %v, want Japanese text", v)
	}

	v, _ = ws2.GetCellValue("A3")
	if v != "Ñoño café" {
		t.Errorf("A3 = %v, want Spanish text", v)
	}
}

func TestXLSXProperties(t *testing.T) {
	wb := New()
	wb.Properties.Creator = "Test Author"
	wb.Properties.Title = "Test Title"
	wb.Properties.Subject = "Test Subject"
	wb.Properties.Description = "Test Description"
	wb.Properties.Keywords = "test, xlsx"
	wb.Properties.Category = "Testing"
	wb.Properties.LastModifiedBy = "Test Modifier"

	tmpDir := t.TempDir()
	filename := filepath.Join(tmpDir, "props.xlsx")

	if err := NewXLSXWriter().Save(wb, filename); err != nil {
		t.Fatalf("Save error: %v", err)
	}

	wb2, err := NewXLSXReader().Open(filename)
	if err != nil {
		t.Fatalf("Open error: %v", err)
	}

	if wb2.Properties.Creator != "Test Author" {
		t.Errorf("Creator = %q, want 'Test Author'", wb2.Properties.Creator)
	}
	if wb2.Properties.Title != "Test Title" {
		t.Errorf("Title = %q, want 'Test Title'", wb2.Properties.Title)
	}
	if wb2.Properties.Subject != "Test Subject" {
		t.Errorf("Subject = %q, want 'Test Subject'", wb2.Properties.Subject)
	}
}

func TestDateSerialization(t *testing.T) {
	// Test known date conversions
	tests := []struct {
		date   time.Time
		serial float64
	}{
		{time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC), 2.0},   // Excel serial 2 (due to 1900 bug)
		{time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), 45292.0},
	}

	for _, tt := range tests {
		serial := dateToSerial(tt.date)
		if serial != tt.serial {
			t.Errorf("dateToSerial(%v) = %f, want %f", tt.date, serial, tt.serial)
		}

		date := serialToDate(tt.serial)
		if !date.Equal(tt.date) {
			t.Errorf("serialToDate(%f) = %v, want %v", tt.serial, date, tt.date)
		}
	}
}

func TestXLSXEmptyWorkbook(t *testing.T) {
	wb := New()
	// Don't add any data

	tmpDir := t.TempDir()
	filename := filepath.Join(tmpDir, "empty.xlsx")

	if err := NewXLSXWriter().Save(wb, filename); err != nil {
		t.Fatalf("Save error: %v", err)
	}

	wb2, err := NewXLSXReader().Open(filename)
	if err != nil {
		t.Fatalf("Open error: %v", err)
	}

	if wb2.SheetCount() != 1 {
		t.Errorf("sheet count = %d, want 1", wb2.SheetCount())
	}
}

func TestXLSXLargeData(t *testing.T) {
	wb := New()
	ws := wb.GetActiveSheet()

	// Write 1000 rows x 10 columns
	for row := 0; row < 1000; row++ {
		for col := 0; col < 10; col++ {
			ws.GetCell(row, col).SetValue(float64(row*10 + col))
		}
	}

	tmpDir := t.TempDir()
	filename := filepath.Join(tmpDir, "large.xlsx")

	if err := NewXLSXWriter().Save(wb, filename); err != nil {
		t.Fatalf("Save error: %v", err)
	}

	wb2, err := NewXLSXReader().Open(filename)
	if err != nil {
		t.Fatalf("Open error: %v", err)
	}

	ws2 := wb2.GetActiveSheet()
	if ws2.CellCount() != 10000 {
		t.Errorf("cell count = %d, want 10000", ws2.CellCount())
	}

	// Spot check
	v, _ := ws2.GetCellValue("A1")
	if v != float64(0) {
		t.Errorf("A1 = %v, want 0", v)
	}

	c := ws2.GetCell(999, 9)
	numVal, err := c.GetNumericValue()
	if err != nil {
		t.Fatalf("GetNumericValue error: %v", err)
	}
	if numVal != float64(9999) {
		t.Errorf("last cell = %f, want 9999", numVal)
	}
}
