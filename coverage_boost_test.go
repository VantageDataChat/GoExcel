package gospreadsheet

import (
	"bytes"
	"path/filepath"
	"testing"
	"time"
)

// --- matchesCriteria coverage (was 40%) ---

func TestMatchesCriteriaExactNumeric(t *testing.T) {
	if !matchesCriteria(float64(10), "10") {
		t.Error("10 should match '10'")
	}
}

func TestMatchesCriteriaExactString(t *testing.T) {
	if !matchesCriteria("hello", "hello") {
		t.Error("'hello' should match 'hello'")
	}
	if !matchesCriteria("Hello", "hello") {
		t.Error("case insensitive match should work")
	}
}

func TestMatchesCriteriaGreaterThan(t *testing.T) {
	if !matchesCriteria(float64(20), ">10") {
		t.Error("20 > 10")
	}
	if matchesCriteria(float64(5), ">10") {
		t.Error("5 is not > 10")
	}
}

func TestMatchesCriteriaLessThan(t *testing.T) {
	if !matchesCriteria(float64(5), "<10") {
		t.Error("5 < 10")
	}
	if matchesCriteria(float64(20), "<10") {
		t.Error("20 is not < 10")
	}
}

func TestMatchesCriteriaGreaterOrEqual(t *testing.T) {
	if !matchesCriteria(float64(10), ">=10") {
		t.Error("10 >= 10")
	}
	if !matchesCriteria(float64(15), ">=10") {
		t.Error("15 >= 10")
	}
	if matchesCriteria(float64(5), ">=10") {
		t.Error("5 is not >= 10")
	}
}

func TestMatchesCriteriaLessOrEqual(t *testing.T) {
	if !matchesCriteria(float64(10), "<=10") {
		t.Error("10 <= 10")
	}
	if !matchesCriteria(float64(5), "<=10") {
		t.Error("5 <= 10")
	}
	if matchesCriteria(float64(15), "<=10") {
		t.Error("15 is not <= 10")
	}
}

func TestMatchesCriteriaNotEqual(t *testing.T) {
	if !matchesCriteria("hello", "<>world") {
		t.Error("hello <> world")
	}
	if matchesCriteria("world", "<>world") {
		t.Error("world is not <> world")
	}
}

func TestMatchesCriteriaInvalidOperator(t *testing.T) {
	// ">abc" should fail to parse and return false
	if matchesCriteria(float64(10), ">abc") {
		t.Error("should not match invalid numeric criteria")
	}
	if matchesCriteria(float64(10), "<abc") {
		t.Error("should not match invalid numeric criteria")
	}
	if matchesCriteria(float64(10), ">=abc") {
		t.Error("should not match invalid numeric criteria")
	}
	if matchesCriteria(float64(10), "<=abc") {
		t.Error("should not match invalid numeric criteria")
	}
}

// --- toFloat coverage (was 22%) ---

func TestToFloatInt(t *testing.T) {
	if toFloat(42) != 42.0 {
		t.Error("int conversion failed")
	}
}

func TestToFloatBool(t *testing.T) {
	if toFloat(true) != 1.0 {
		t.Error("true should be 1")
	}
	if toFloat(false) != 0.0 {
		t.Error("false should be 0")
	}
}

func TestToFloatString(t *testing.T) {
	if toFloat("3.14") != 3.14 {
		t.Error("string '3.14' should convert")
	}
	if toFloat("not a number") != 0 {
		t.Error("non-numeric string should be 0")
	}
}

func TestToFloatNil(t *testing.T) {
	if toFloat(nil) != 0 {
		t.Error("nil should be 0")
	}
}

// --- parseArgs coverage (was 76%) ---

func TestParseArgsNested(t *testing.T) {
	args := parseArgs("SUM(A1:A3),B1,IF(C1,1,2)")
	if len(args) != 3 {
		t.Fatalf("expected 3 args, got %d: %v", len(args), args)
	}
	if args[0] != "SUM(A1:A3)" {
		t.Errorf("arg[0] = %q", args[0])
	}
	if args[2] != "IF(C1,1,2)" {
		t.Errorf("arg[2] = %q", args[2])
	}
}

func TestParseArgsEmpty(t *testing.T) {
	args := parseArgs("")
	if len(args) != 0 {
		t.Errorf("expected 0 args, got %d", len(args))
	}
}

func TestParseArgsSingle(t *testing.T) {
	args := parseArgs("A1:A10")
	if len(args) != 1 || args[0] != "A1:A10" {
		t.Errorf("args = %v", args)
	}
}

// --- SetCellStyle coverage (was 0%) ---

func TestWorksheetSetCellStyle(t *testing.T) {
	ws := NewWorksheet("Test")
	ws.SetCellValue("A1", "styled")
	style := NewStyle().SetFont(&Font{Bold: true})
	err := ws.SetCellStyle("A1", style)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	c, _ := ws.GetCellByName("A1")
	if c.Style == nil || !c.Style.Font.Bold {
		t.Error("style not applied")
	}

	err = ws.SetCellStyle("invalid!", style)
	if err == nil {
		t.Error("expected error for invalid ref")
	}
}

// --- SetPageSetup coverage (was 0%) ---

func TestWorksheetSetPageSetupExplicit(t *testing.T) {
	ws := NewWorksheet("Test")
	ps := NewPageSetup().SetOrientation(OrientationLandscape)
	ws.SetPageSetup(ps)
	if ws.GetPageSetup().Orientation != OrientationLandscape {
		t.Error("page setup not set")
	}
}

// --- SetMargins coverage (was 0%) ---

func TestPageSetupSetMargins(t *testing.T) {
	ps := NewPageSetup()
	m := PageMargins{Top: 1.0, Bottom: 1.0, Left: 0.5, Right: 0.5, Header: 0.4, Footer: 0.4}
	ps.SetMargins(m)
	if ps.Margins.Top != 1.0 {
		t.Errorf("Top = %f, want 1.0", ps.Margins.Top)
	}
	if ps.Margins.Left != 0.5 {
		t.Errorf("Left = %f, want 0.5", ps.Margins.Left)
	}
}

// --- AddColumn coverage (was 0%) ---

func TestAutoFilterAddColumn(t *testing.T) {
	af := NewAutoFilter("A1:D10")
	af.AddColumn(AutoFilterColumn{
		ColumnIndex: 0,
		FilterType:  FilterValues,
		Values:      []string{"A"},
	})
	if len(af.Columns) != 1 {
		t.Fatalf("columns = %d", len(af.Columns))
	}
	if !af.Columns[0].ShowButton {
		t.Error("ShowButton should be true")
	}
}

// --- GetActiveSheet nil path (was 66%) ---

func TestGetActiveSheetEmpty(t *testing.T) {
	wb := NewEmpty()
	if wb.GetActiveSheet() != nil {
		t.Error("should return nil for empty workbook")
	}
}

// --- DeleteRow full coverage (was 53%) ---

func TestDeleteRowWithMergeCells(t *testing.T) {
	ws := NewWorksheet("Test")
	ws.SetCellValue("A1", "keep")
	ws.SetCellValue("A2", "delete")
	ws.SetCellValue("A3", "shift up")
	ws.MergeCells("A2:C2") // merge in deleted row
	ws.MergeCells("A3:C3") // merge below deleted row
	ws.SetRowHeight(3, 25.0)

	ws.DeleteRow(1)

	merges := ws.GetMergeCells()
	// The A2:C2 merge should be removed, A3:C3 should shift to A2:C2
	if len(merges) != 1 {
		t.Fatalf("merges = %d, want 1", len(merges))
	}
	if merges[0].StartRow != 1 {
		t.Errorf("merge start row = %d, want 1", merges[0].StartRow)
	}

	// Row height should shift
	if ws.GetRowHeight(2) != 25.0 {
		t.Errorf("row height = %f, want 25", ws.GetRowHeight(2))
	}
}

func TestDeleteRowHeightShift(t *testing.T) {
	ws := NewWorksheet("Test")
	ws.SetRowHeight(0, 20.0)
	ws.SetRowHeight(1, 30.0)
	ws.SetRowHeight(2, 40.0)
	ws.SetCellValue("A1", "x") // need data for DeleteRow to work on

	ws.DeleteRow(1)

	if ws.GetRowHeight(0) != 20.0 {
		t.Errorf("row 0 height = %f, want 20", ws.GetRowHeight(0))
	}
	if ws.GetRowHeight(1) != 40.0 {
		t.Errorf("row 1 height = %f, want 40", ws.GetRowHeight(1))
	}
}

// --- InsertColumn full coverage (was 61%) ---

func TestInsertColumnWithMergesAndWidths(t *testing.T) {
	ws := NewWorksheet("Test")
	ws.SetCellValue("A1", "a")
	ws.SetCellValue("B1", "b")
	ws.MergeCells("B1:C1")
	ws.SetColumnWidth(1, 20.0)
	ws.SetColumnWidth(2, 30.0)

	ws.InsertColumn(1)

	merges := ws.GetMergeCells()
	if len(merges) != 1 {
		t.Fatalf("merges = %d", len(merges))
	}
	if merges[0].StartCol != 2 {
		t.Errorf("merge start col = %d, want 2", merges[0].StartCol)
	}
	if merges[0].EndCol != 3 {
		t.Errorf("merge end col = %d, want 3", merges[0].EndCol)
	}

	if ws.GetColumnWidth(2) != 20.0 {
		t.Errorf("col 2 width = %f, want 20", ws.GetColumnWidth(2))
	}
	if ws.GetColumnWidth(3) != 30.0 {
		t.Errorf("col 3 width = %f, want 30", ws.GetColumnWidth(3))
	}
}

// --- DeleteColumn full coverage (was 65%) ---

func TestDeleteColumnWithMergesAndWidths(t *testing.T) {
	ws := NewWorksheet("Test")
	ws.SetCellValue("A1", "a")
	ws.SetCellValue("B1", "b")
	ws.SetCellValue("C1", "c")
	ws.MergeCells("B1:B1") // single col merge, will be removed
	ws.MergeCells("C1:D1") // merge after deleted col
	ws.SetColumnWidth(0, 10.0)
	ws.SetColumnWidth(2, 30.0)

	ws.DeleteColumn(1)

	merges := ws.GetMergeCells()
	if len(merges) != 1 {
		t.Fatalf("merges = %d, want 1", len(merges))
	}
	if merges[0].StartCol != 1 {
		t.Errorf("merge start col = %d, want 1", merges[0].StartCol)
	}

	if ws.GetColumnWidth(0) != 10.0 {
		t.Errorf("col 0 width = %f, want 10", ws.GetColumnWidth(0))
	}
	if ws.GetColumnWidth(1) != 30.0 {
		t.Errorf("col 1 width = %f, want 30", ws.GetColumnWidth(1))
	}
}

// --- isCellRef edge cases (was 87%) ---

func TestIsCellRef(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"A1", true},
		{"Z99", true},
		{"AA100", true},
		{"", false},
		{"1A", false},
		{"A", false},
		{"1", false},
		{"A1B", false}, // letter after digit
		{"A 1", false}, // space
	}
	for _, tt := range tests {
		got := isCellRef(tt.input)
		if got != tt.want {
			t.Errorf("isCellRef(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

// --- evaluate edge cases ---

func TestEvaluateEmptyFormula(t *testing.T) {
	_, ws, ce := setupCalcSheet()
	_, err := ce.evaluate(ws, "")
	if err == nil {
		t.Error("expected error for empty formula")
	}
}

func TestEvaluateStringLiteral(t *testing.T) {
	_, ws, ce := setupCalcSheet()
	val, err := ce.evaluate(ws, `"hello"`)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if val != "hello" {
		t.Errorf("val = %v, want 'hello'", val)
	}
}

func TestEvaluateUnknownFormula(t *testing.T) {
	_, ws, ce := setupCalcSheet()
	_, err := ce.evaluate(ws, "UNKNOWN_THING")
	if err == nil {
		t.Error("expected error for unknown formula")
	}
}

// --- CalculateCell error path ---

func TestCalculateCellInvalidRef(t *testing.T) {
	_, ws, ce := setupCalcSheet()
	_, err := ce.CalculateCell(ws, "!!!")
	if err == nil {
		t.Error("expected error for invalid ref")
	}
}

// --- CalculateAll with error formula ---

func TestCalculateAllWithError(t *testing.T) {
	wb, ws, _ := setupCalcSheet()
	ws.SetCellFormula("A1", "UNKNOWN_FUNC()")
	ce := NewCalculationEngine(wb)
	ce.CalculateAll()

	c, _ := ws.GetCellByName("A1")
	if c.Value != "#ERROR!" {
		t.Errorf("value = %v, want '#ERROR!'", c.Value)
	}
}

// --- findOperator edge cases ---

func TestFindOperatorWithParens(t *testing.T) {
	// Should not find + inside parentheses
	idx := findOperator("SUM(A1+B1)", "+")
	if idx != -1 {
		t.Errorf("should not find + inside parens, got idx=%d", idx)
	}
}

// --- applyOp unknown operator ---

func TestApplyOpUnknown(t *testing.T) {
	_, err := applyOp(1, 2, "^")
	if err == nil {
		t.Error("expected error for unknown operator")
	}
}

// --- XLSX writer: freeze pane + column widths + row heights ---

func TestXLSXFreezePaneRoundTrip(t *testing.T) {
	wb := New()
	ws := wb.GetActiveSheet()
	ws.SetCellValue("A1", "Header")
	ws.FreezePane("B2")
	ws.SetColumnWidth(0, 20.0)
	ws.SetRowHeight(0, 30.0)

	var buf bytes.Buffer
	NewXLSXWriter().Write(wb, &buf)

	reader := bytes.NewReader(buf.Bytes())
	wb2, err := NewXLSXReader().Read(reader, int64(buf.Len()))
	if err != nil {
		t.Fatalf("Read error: %v", err)
	}
	// Just verify it reads without error and data is intact
	v, _ := wb2.GetActiveSheet().GetCellValue("A1")
	if v != "Header" {
		t.Errorf("A1 = %v", v)
	}
}

// --- XLSX writer: date cell ---

func TestXLSXDateRoundTrip(t *testing.T) {
	wb := New()
	ws := wb.GetActiveSheet()
	date := time.Date(2025, 6, 15, 0, 0, 0, 0, time.UTC)
	ws.GetCell(0, 0).SetValue(date)

	var buf bytes.Buffer
	NewXLSXWriter().Write(wb, &buf)

	reader := bytes.NewReader(buf.Bytes())
	wb2, err := NewXLSXReader().Read(reader, int64(buf.Len()))
	if err != nil {
		t.Fatalf("Read error: %v", err)
	}
	// Date is stored as serial number
	c := wb2.GetActiveSheet().GetCell(0, 0)
	v, err := c.GetNumericValue()
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if v < 45000 { // 2025 dates are around 45000+
		t.Errorf("serial = %f, expected > 45000", v)
	}
}

// --- XLSX writer: bool cell false ---

func TestXLSXBoolFalseRoundTrip(t *testing.T) {
	wb := New()
	ws := wb.GetActiveSheet()
	ws.SetCellValue("A1", false)

	var buf bytes.Buffer
	NewXLSXWriter().Write(wb, &buf)

	reader := bytes.NewReader(buf.Bytes())
	wb2, err := NewXLSXReader().Read(reader, int64(buf.Len()))
	if err != nil {
		t.Fatalf("Read error: %v", err)
	}
	v, _ := wb2.GetActiveSheet().GetCellValue("A1")
	if v != false {
		t.Errorf("A1 = %v, want false", v)
	}
}

// --- CSV error paths ---

func TestCSVWriterInvalidSheet(t *testing.T) {
	wb := New()
	var buf bytes.Buffer
	writer := NewCSVWriter()
	writer.SheetIndex = 99
	err := writer.Write(wb, &buf)
	if err == nil {
		t.Error("expected error for invalid sheet index")
	}
}

func TestCSVOpenNonExistent(t *testing.T) {
	_, err := NewCSVReader().Open("/nonexistent/path/file.csv")
	if err == nil {
		t.Error("expected error for non-existent file")
	}
}

func TestCSVSaveInvalidPath(t *testing.T) {
	wb := New()
	err := NewCSVWriter().Save(wb, "/nonexistent/path/file.csv")
	if err == nil {
		t.Error("expected error for invalid path")
	}
}

// --- XLSX error paths ---

func TestXLSXOpenNonExistent(t *testing.T) {
	_, err := NewXLSXReader().Open("/nonexistent/path/file.xlsx")
	if err == nil {
		t.Error("expected error for non-existent file")
	}
}

func TestXLSXSaveInvalidPath(t *testing.T) {
	wb := New()
	err := NewXLSXWriter().Save(wb, "/nonexistent/path/file.xlsx")
	if err == nil {
		t.Error("expected error for invalid path")
	}
}

// --- CellName error path ---

func TestCellNameNegativeCol(t *testing.T) {
	_, err := CellName(0, -1)
	if err == nil {
		t.Error("expected error for negative column")
	}
}

// --- SetValue default case ---

func TestCellSetValueUnknownType(t *testing.T) {
	c := NewCell(0, 0)
	type custom struct{ X int }
	c.SetValue(custom{X: 1})
	if c.Type != CellTypeString {
		t.Errorf("type = %d, want CellTypeString", c.Type)
	}
}

// --- GetStringValue default case ---

func TestCellGetStringValueDefault(t *testing.T) {
	c := NewCell(0, 0)
	type custom struct{ X int }
	c.Value = custom{X: 42}
	c.Type = CellTypeString
	s := c.GetStringValue()
	if s == "" {
		t.Error("should return non-empty string for custom type")
	}
}

// --- GetStringValue date ---

func TestCellGetStringValueDate(t *testing.T) {
	c := NewCell(0, 0)
	c.SetValue(time.Date(2025, 1, 15, 10, 30, 0, 0, time.UTC))
	s := c.GetStringValue()
	if s != "2025-01-15 10:30:00" {
		t.Errorf("date string = %q", s)
	}
}

// --- XLSX with inline string type ---

func TestXLSXReadInlineString(t *testing.T) {
	// This tests the "str" type path in readSheet
	wb := New()
	ws := wb.GetActiveSheet()
	ws.SetCellValue("A1", "test")

	tmpDir := t.TempDir()
	filename := filepath.Join(tmpDir, "test.xlsx")
	NewXLSXWriter().Save(wb, filename)

	wb2, err := NewXLSXReader().Open(filename)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	v, _ := wb2.GetActiveSheet().GetCellValue("A1")
	if v != "test" {
		t.Errorf("A1 = %v", v)
	}
}

// --- Worksheet SetCellFormula/SetCellValue error paths ---

func TestSetCellValueInvalidRef(t *testing.T) {
	ws := NewWorksheet("Test")
	err := ws.SetCellValue("!!!", "hello")
	if err == nil {
		t.Error("expected error")
	}
}

func TestSetCellFormulaInvalidRef(t *testing.T) {
	ws := NewWorksheet("Test")
	err := ws.SetCellFormula("!!!", "SUM(A1)")
	if err == nil {
		t.Error("expected error")
	}
}

func TestGetCellValueInvalidRef(t *testing.T) {
	ws := NewWorksheet("Test")
	_, err := ws.GetCellValue("!!!")
	if err == nil {
		t.Error("expected error")
	}
}

// --- SetCellHyperlink/SetCellComment error paths ---

func TestSetCellHyperlinkInvalidRef(t *testing.T) {
	ws := NewWorksheet("Test")
	err := ws.SetCellHyperlink("!!!", "https://example.com")
	if err == nil {
		t.Error("expected error")
	}
}

func TestSetCellCommentInvalidRef(t *testing.T) {
	ws := NewWorksheet("Test")
	err := ws.SetCellComment("!!!", "author", "text")
	if err == nil {
		t.Error("expected error")
	}
}

// --- FreezePane error path ---

func TestFreezePaneInvalidRef(t *testing.T) {
	ws := NewWorksheet("Test")
	err := ws.FreezePane("!!!")
	if err == nil {
		t.Error("expected error")
	}
}

// --- CopyRow on empty sheet ---

func TestCopyRowEmptySheet(t *testing.T) {
	ws := NewWorksheet("Test")
	// Should not panic
	ws.CopyRow(0, 1)
}

// --- ParseRange error paths ---

func TestParseRangeInvalidStart(t *testing.T) {
	_, _, err := ParseRange("!!!:A1")
	if err == nil {
		t.Error("expected error")
	}
}

func TestParseRangeInvalidEnd(t *testing.T) {
	_, _, err := ParseRange("A1:!!!")
	if err == nil {
		t.Error("expected error")
	}
}

// --- resolveRange with formula cells ---

func TestResolveRangeWithFormulas(t *testing.T) {
	_, ws, ce := setupCalcSheet()
	ws.SetCellValue("A1", 10)
	ws.SetCellFormula("A2", "A1+5")
	ws.SetCellFormula("A3", "SUM(A1:A2)")

	val, err := ce.CalculateCell(ws, "A3")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	// A1=10, A2=10+5=15, SUM=10+15=25
	if val != 25.0 {
		t.Errorf("SUM = %v, want 25", val)
	}
}

// --- SharedStrings GetIndex miss ---

func TestSharedStringsGetIndexMiss(t *testing.T) {
	ss := newSharedStrings()
	ss.Add("hello")
	idx := ss.GetIndex("nonexistent")
	if idx != -1 {
		t.Errorf("expected -1, got %d", idx)
	}
}

// --- SharedStrings Add duplicate ---

func TestSharedStringsAddDuplicate(t *testing.T) {
	ss := newSharedStrings()
	idx1 := ss.Add("hello")
	idx2 := ss.Add("hello")
	if idx1 != idx2 {
		t.Errorf("duplicate should return same index: %d vs %d", idx1, idx2)
	}
	if len(ss.strings) != 1 {
		t.Errorf("should have 1 string, got %d", len(ss.strings))
	}
}
