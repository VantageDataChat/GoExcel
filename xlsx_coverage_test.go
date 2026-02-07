package gospreadsheet

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

// errWriter is an io.Writer that always returns an error.
type errWriter struct{}

func (e *errWriter) Write(p []byte) (int, error) {
	return 0, errors.New("write error")
}

func TestXLSXWriteToErrWriter(t *testing.T) {
	wb := New()
	wb.GetActiveSheet().SetCellValue("A1", "test")
	// The zip writer buffers internally, so errors may only surface on Close.
	// We just verify it doesn't panic.
	_ = NewXLSXWriter().Write(wb, &errWriter{})
}

func TestXLSXReadInvalidZip(t *testing.T) {
	data := []byte("this is not a zip file")
	reader := bytes.NewReader(data)
	_, err := NewXLSXReader().Read(reader, int64(len(data)))
	if err == nil {
		t.Error("expected error for invalid zip")
	}
}

func TestXLSXReadErrorCellType(t *testing.T) {
	// Create an XLSX with an error cell type
	wb := New()
	ws := wb.GetActiveSheet()
	// Manually set an error cell
	c := ws.GetCell(0, 0)
	c.Value = "#REF!"
	c.Type = CellTypeError

	// Write and read back - error cells won't be written as error type
	// but we can verify the writer doesn't crash
	var buf bytes.Buffer
	err := NewXLSXWriter().Write(wb, &buf)
	if err != nil {
		t.Fatalf("Write error: %v", err)
	}
}

func TestXLSXWriterAllCellTypes(t *testing.T) {
	// Ensure all cell type branches in writeSheet are covered
	wb := New()
	ws := wb.GetActiveSheet()
	ws.SetCellValue("A1", "string")
	ws.SetCellValue("A2", 42.5)
	ws.SetCellValue("A3", true)
	ws.SetCellValue("A4", false)
	ws.SetCellFormula("A5", "SUM(A2)")
	ws.MergeCells("B1:C1")
	ws.SetColumnWidth(0, 15.0)
	ws.SetRowHeight(0, 20.0)
	ws.FreezePane("A2")

	var buf bytes.Buffer
	err := NewXLSXWriter().Write(wb, &buf)
	if err != nil {
		t.Fatalf("Write error: %v", err)
	}

	// Read back and verify
	reader := bytes.NewReader(buf.Bytes())
	wb2, err := NewXLSXReader().Read(reader, int64(buf.Len()))
	if err != nil {
		t.Fatalf("Read error: %v", err)
	}

	ws2 := wb2.GetActiveSheet()
	v, _ := ws2.GetCellValue("A1")
	if v != "string" {
		t.Errorf("A1 = %v", v)
	}
	v, _ = ws2.GetCellValue("A2")
	if v != 42.5 {
		t.Errorf("A2 = %v", v)
	}
	v, _ = ws2.GetCellValue("A3")
	if v != true {
		t.Errorf("A3 = %v", v)
	}
	v, _ = ws2.GetCellValue("A4")
	if v != false {
		t.Errorf("A4 = %v", v)
	}
	c, _ := ws2.GetCellByName("A5")
	if c.Type != CellTypeFormula || c.Formula != "SUM(A2)" {
		t.Errorf("A5 formula = %q, type = %d", c.Formula, c.Type)
	}
}

func TestXLSXWriterNoProperties(t *testing.T) {
	// Test with empty properties (all branches in writeCoreProperties)
	wb := New()
	ws := wb.GetActiveSheet()
	ws.SetCellValue("A1", "test")
	// Properties are all empty by default

	var buf bytes.Buffer
	err := NewXLSXWriter().Write(wb, &buf)
	if err != nil {
		t.Fatalf("Write error: %v", err)
	}
}

func TestXLSXWriterAllProperties(t *testing.T) {
	wb := New()
	ws := wb.GetActiveSheet()
	ws.SetCellValue("A1", "test")
	wb.Properties = DocumentProperties{
		Creator:        "Creator",
		LastModifiedBy: "Modifier",
		Title:          "Title",
		Subject:        "Subject",
		Description:    "Description",
		Keywords:       "Keywords",
		Category:       "Category",
	}

	var buf bytes.Buffer
	err := NewXLSXWriter().Write(wb, &buf)
	if err != nil {
		t.Fatalf("Write error: %v", err)
	}

	// Read back
	reader := bytes.NewReader(buf.Bytes())
	wb2, err := NewXLSXReader().Read(reader, int64(buf.Len()))
	if err != nil {
		t.Fatalf("Read error: %v", err)
	}
	if wb2.Properties.Creator != "Creator" {
		t.Errorf("Creator = %q", wb2.Properties.Creator)
	}
	if wb2.Properties.LastModifiedBy != "Modifier" {
		t.Errorf("LastModifiedBy = %q", wb2.Properties.LastModifiedBy)
	}
	if wb2.Properties.Description != "Description" {
		t.Errorf("Description = %q", wb2.Properties.Description)
	}
	if wb2.Properties.Keywords != "Keywords" {
		t.Errorf("Keywords = %q", wb2.Properties.Keywords)
	}
	if wb2.Properties.Category != "Category" {
		t.Errorf("Category = %q", wb2.Properties.Category)
	}
}

func TestXLSXWriterMultiRow(t *testing.T) {
	// Test multiple rows to cover the row-closing logic
	wb := New()
	ws := wb.GetActiveSheet()
	ws.SetCellValue("A1", "r1")
	ws.SetCellValue("A2", "r2")
	ws.SetCellValue("A3", "r3")
	ws.SetCellValue("B1", "c2")
	ws.SetRowHeight(1, 25.0)

	var buf bytes.Buffer
	err := NewXLSXWriter().Write(wb, &buf)
	if err != nil {
		t.Fatalf("Write error: %v", err)
	}

	reader := bytes.NewReader(buf.Bytes())
	wb2, err := NewXLSXReader().Read(reader, int64(buf.Len()))
	if err != nil {
		t.Fatalf("Read error: %v", err)
	}
	if wb2.GetActiveSheet().CellCount() != 4 {
		t.Errorf("cell count = %d, want 4", wb2.GetActiveSheet().CellCount())
	}
}

// --- fnMin empty values ---

func TestFnMinEmpty(t *testing.T) {
	_, ws, ce := setupCalcSheet()
	// MIN with a range that has values
	ws.SetCellValue("A1", 5)
	ws.SetCellValue("A2", 3)
	ws.SetCellValue("A3", 8)
	ws.SetCellFormula("A4", "MIN(A1:A3)")
	val, err := ce.CalculateCell(ws, "A4")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if val != 3.0 {
		t.Errorf("MIN = %v, want 3", val)
	}
}

// --- fnAverage with single value ---

func TestFnAverageSingle(t *testing.T) {
	_, ws, ce := setupCalcSheet()
	ws.SetCellValue("A1", 42)
	ws.SetCellFormula("B1", "AVERAGE(A1)")
	val, err := ce.CalculateCell(ws, "B1")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if val != 42.0 {
		t.Errorf("AVERAGE = %v, want 42", val)
	}
}

// --- fnMedian single value ---

func TestFnMedianSingle(t *testing.T) {
	_, ws, ce := setupCalcSheet()
	ws.SetCellValue("A1", 7)
	ws.SetCellFormula("B1", "MEDIAN(A1)")
	val, err := ce.CalculateCell(ws, "B1")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if val != 7.0 {
		t.Errorf("MEDIAN = %v, want 7", val)
	}
}

// --- fnMax empty ---

func TestFnMaxEmpty(t *testing.T) {
	_, ws, ce := setupCalcSheet()
	ws.SetCellFormula("A1", "MAX(Z1:Z1)")
	val, err := ce.CalculateCell(ws, "A1")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	// Empty range returns 0
	if val != 0.0 {
		t.Errorf("MAX = %v, want 0", val)
	}
}

// --- fnMin empty ---

func TestFnMinEmptyRange(t *testing.T) {
	_, ws, ce := setupCalcSheet()
	ws.SetCellFormula("A1", "MIN(Z1:Z1)")
	val, err := ce.CalculateCell(ws, "A1")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if val != 0.0 {
		t.Errorf("MIN = %v, want 0", val)
	}
}

// --- XLSX reader: readSheet with "str" type ---

func TestXLSXReaderInlineString(t *testing.T) {
	// We can't easily create an XLSX with inline strings via our writer,
	// but we can verify the reader handles all cell types by writing
	// and reading back various types
	wb := New()
	ws := wb.GetActiveSheet()
	ws.SetCellValue("A1", "inline test")

	var buf bytes.Buffer
	NewXLSXWriter().Write(wb, &buf)

	reader := bytes.NewReader(buf.Bytes())
	wb2, _ := NewXLSXReader().Read(reader, int64(buf.Len()))
	v, _ := wb2.GetActiveSheet().GetCellValue("A1")
	if v != "inline test" {
		t.Errorf("A1 = %v", v)
	}
}

// --- XLSX writer: shared string not found (GetIndex returns -1) ---

func TestSharedStringsGetIndexNotFound(t *testing.T) {
	ss := newSharedStrings()
	idx := ss.GetIndex("missing")
	if idx != -1 {
		t.Errorf("expected -1, got %d", idx)
	}
}

// --- CSV reader with LazyQuotes ---

func TestCSVReaderLazyQuotes(t *testing.T) {
	csvData := `Name,Value
"Alice",10
Bob,"it's "quoted""
`
	reader := NewCSVReader()
	reader.LazyQuotes = true
	wb, err := reader.Read(strings.NewReader(csvData))
	if err != nil {
		t.Fatalf("Read error: %v", err)
	}
	ws := wb.GetActiveSheet()
	if ws.CellCount() == 0 {
		t.Error("should have cells")
	}
}
