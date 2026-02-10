package gospreadsheet

import (
	"testing"
	"time"
)

// ============================================================
// Cell: Clear, HasFormula, IsEmpty, IsNumber, IsBool, etc.
// ============================================================

func TestCellClear(t *testing.T) {
	c := NewCell(0, 0)
	c.SetValue("hello")
	c.SetStyle(NewStyle().SetFont(&Font{Bold: true}))
	c.SetHyperlink(NewHyperlink("https://example.com"))
	c.SetComment(NewComment("author", "text"))

	c.Clear()

	if c.Value != nil {
		t.Errorf("Value should be nil after Clear, got %v", c.Value)
	}
	if c.Type != CellTypeEmpty {
		t.Errorf("Type should be CellTypeEmpty after Clear, got %d", c.Type)
	}
	if c.Formula != "" {
		t.Errorf("Formula should be empty after Clear, got %q", c.Formula)
	}
	if c.Style != nil {
		t.Error("Style should be nil after Clear")
	}
	if c.Hyperlink != nil {
		t.Error("Hyperlink should be nil after Clear")
	}
	if c.Comment != nil {
		t.Error("Comment should be nil after Clear")
	}
	if c.RichText != nil {
		t.Error("RichText should be nil after Clear")
	}
}

func TestCellClearFormula(t *testing.T) {
	c := NewCell(0, 0)
	c.SetFormula("SUM(A1:A10)")
	c.Clear()
	if c.HasFormula() {
		t.Error("HasFormula should be false after Clear")
	}
}

func TestCellHasFormula(t *testing.T) {
	c := NewCell(0, 0)
	if c.HasFormula() {
		t.Error("empty cell should not have formula")
	}
	c.SetFormula("SUM(A1:A10)")
	if !c.HasFormula() {
		t.Error("cell with formula should return true")
	}
	c.SetValue("hello")
	if c.HasFormula() {
		t.Error("cell with string value should not have formula")
	}
}

func TestCellIsEmpty(t *testing.T) {
	c := NewCell(0, 0)
	if !c.IsEmpty() {
		t.Error("new cell should be empty")
	}
	c.SetValue("hello")
	if c.IsEmpty() {
		t.Error("cell with value should not be empty")
	}
	c.Clear()
	if !c.IsEmpty() {
		t.Error("cleared cell should be empty")
	}
}

func TestCellIsNumber(t *testing.T) {
	c := NewCell(0, 0)
	if c.IsNumber() {
		t.Error("empty cell should not be number")
	}
	c.SetValue(42)
	if !c.IsNumber() {
		t.Error("cell with int should be number")
	}
	c.SetValue(3.14)
	if !c.IsNumber() {
		t.Error("cell with float should be number")
	}
	c.SetValue("hello")
	if c.IsNumber() {
		t.Error("cell with string should not be number")
	}
}

func TestCellIsBool(t *testing.T) {
	c := NewCell(0, 0)
	if c.IsBool() {
		t.Error("empty cell should not be bool")
	}
	c.SetValue(true)
	if !c.IsBool() {
		t.Error("cell with true should be bool")
	}
	c.SetValue(false)
	if !c.IsBool() {
		t.Error("cell with false should be bool")
	}
}

func TestCellIsDate(t *testing.T) {
	c := NewCell(0, 0)
	if c.IsDate() {
		t.Error("empty cell should not be date")
	}
	c.SetValue(time.Now())
	if !c.IsDate() {
		t.Error("cell with time should be date")
	}
}

func TestCellIsString(t *testing.T) {
	c := NewCell(0, 0)
	c.SetValue("hello")
	if !c.IsString() {
		t.Error("cell with string should be string")
	}
	c.SetValue(42)
	if c.IsString() {
		t.Error("cell with number should not be string")
	}
}

func TestCellIsError(t *testing.T) {
	c := NewCell(0, 0)
	c.Value = "#DIV/0!"
	c.Type = CellTypeError
	if !c.IsError() {
		t.Error("cell with error should be error")
	}
	c.SetValue("hello")
	if c.IsError() {
		t.Error("cell with string should not be error")
	}
}

// ============================================================
// Cell: SetDateWithStyle, SetNumberWithFormat
// ============================================================

func TestCellSetDateWithStyle(t *testing.T) {
	c := NewCell(0, 0)
	now := time.Date(2025, 6, 15, 10, 30, 0, 0, time.UTC)
	c.SetDateWithStyle(now)

	if c.Type != CellTypeDate {
		t.Errorf("expected CellTypeDate, got %d", c.Type)
	}
	if c.Style == nil || c.Style.NumberFormat == nil {
		t.Fatal("style and number format should be set")
	}
	if c.Style.NumberFormat.FormatCode != "yyyy-mm-dd" {
		t.Errorf("format code = %q, want 'yyyy-mm-dd'", c.Style.NumberFormat.FormatCode)
	}
	v, err := c.GetDateValue()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !v.Equal(now) {
		t.Errorf("date = %v, want %v", v, now)
	}
}

func TestCellSetDateWithStylePreservesExistingStyle(t *testing.T) {
	c := NewCell(0, 0)
	c.SetStyle(NewStyle().SetFont(&Font{Bold: true}))
	c.SetDateWithStyle(time.Now())

	if c.Style.Font == nil || !c.Style.Font.Bold {
		t.Error("existing font style should be preserved")
	}
	if c.Style.NumberFormat == nil {
		t.Error("number format should be set")
	}
}

func TestCellSetNumberWithFormat(t *testing.T) {
	c := NewCell(0, 0)
	c.SetNumberWithFormat(0.85, FormatPercent)

	if c.Type != CellTypeNumeric {
		t.Errorf("expected CellTypeNumeric, got %d", c.Type)
	}
	v, err := c.GetNumericValue()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v != 0.85 {
		t.Errorf("value = %f, want 0.85", v)
	}
	if c.Style == nil || c.Style.NumberFormat == nil {
		t.Fatal("style and number format should be set")
	}
	if c.Style.NumberFormat.FormatCode != "0%" {
		t.Errorf("format code = %q, want '0%%'", c.Style.NumberFormat.FormatCode)
	}
}

// ============================================================
// Cell: SetFormulaArray, SetInlineString
// ============================================================

func TestCellSetFormulaArray(t *testing.T) {
	c := NewCell(0, 0)
	c.SetFormulaArray("SUM(A1:A10*B1:B10)")

	if c.Type != CellTypeFormula {
		t.Errorf("expected CellTypeFormula, got %d", c.Type)
	}
	if c.Formula != "{SUM(A1:A10*B1:B10)}" {
		t.Errorf("formula = %q, want '{SUM(A1:A10*B1:B10)}'", c.Formula)
	}
	if !c.HasFormula() {
		t.Error("HasFormula should be true")
	}
}

func TestCellSetInlineString(t *testing.T) {
	c := NewCell(0, 0)
	c.SetInlineString("inline text")

	if c.Type != CellTypeString {
		t.Errorf("expected CellTypeString, got %d", c.Type)
	}
	if c.GetStringValue() != "inline text" {
		t.Errorf("value = %q, want 'inline text'", c.GetStringValue())
	}
}

// ============================================================
// Cell: GetFormattedValue
// ============================================================

func TestCellGetFormattedValue(t *testing.T) {
	tests := []struct {
		name   string
		value  interface{}
		format *NumberFormat
		want   string
	}{
		{"no format string", "hello", nil, "hello"},
		{"no format number", 42, nil, "42"},
		{"general int", 42.0, &NumberFormat{FormatCode: "General"}, "42"},
		{"general float", 3.14, &NumberFormat{FormatCode: "General"}, "3.14"},
		{"integer format", 42.7, &NumberFormat{FormatCode: "0"}, "43"},
		{"two decimal", 42.7, &NumberFormat{FormatCode: "0.00"}, "42.70"},
		{"percent", 0.85, &NumberFormat{FormatCode: "0%"}, "85%"},
		{"percent 2dec", 0.8567, &NumberFormat{FormatCode: "0.00%"}, "85.67%"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCell(0, 0)
			c.SetValue(tt.value)
			if tt.format != nil {
				c.SetStyle(NewStyle().SetNumberFormat(tt.format))
			}
			got := c.GetFormattedValue()
			if got != tt.want {
				t.Errorf("GetFormattedValue() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestCellGetFormattedValueDate(t *testing.T) {
	c := NewCell(0, 0)
	dt := time.Date(2025, 3, 15, 14, 30, 45, 0, time.UTC)
	c.SetValue(dt)
	c.SetStyle(NewStyle().SetNumberFormat(&NumberFormat{FormatCode: "yyyy-mm-dd"}))

	got := c.GetFormattedValue()
	if got != "2025-03-15" {
		t.Errorf("GetFormattedValue() = %q, want '2025-03-15'", got)
	}

	c.Style.NumberFormat.FormatCode = "yyyy-mm-dd hh:mm:ss"
	got = c.GetFormattedValue()
	if got != "2025-03-15 14:30:45" {
		t.Errorf("GetFormattedValue() = %q, want '2025-03-15 14:30:45'", got)
	}

	c.Style.NumberFormat.FormatCode = "hh:mm:ss"
	got = c.GetFormattedValue()
	if got != "14:30:45" {
		t.Errorf("GetFormattedValue() = %q, want '14:30:45'", got)
	}
}

// ============================================================
// Style: ShrinkToFit, Diagonal borders
// ============================================================

func TestAlignmentShrinkToFit(t *testing.T) {
	s := NewStyle().SetAlignment(&Alignment{
		Horizontal:  AlignCenter,
		ShrinkToFit: true,
	})
	if !s.Alignment.ShrinkToFit {
		t.Error("ShrinkToFit should be true")
	}
}

func TestDiagonalBorder(t *testing.T) {
	s := NewStyle().SetBorders(&Borders{
		Diagonal:     Border{Style: BorderThin, Color: "FF0000"},
		DiagonalUp:   true,
		DiagonalDown: false,
	})
	if s.Borders.Diagonal.Style != BorderThin {
		t.Errorf("Diagonal.Style = %q, want 'thin'", s.Borders.Diagonal.Style)
	}
	if !s.Borders.DiagonalUp {
		t.Error("DiagonalUp should be true")
	}
	if s.Borders.DiagonalDown {
		t.Error("DiagonalDown should be false")
	}
}

// ============================================================
// Worksheet: RemoveMergedCell
// ============================================================

func TestRemoveMergedCell(t *testing.T) {
	ws := NewWorksheet("Test")
	ws.MergeCells("A1:C3")
	ws.MergeCells("D1:F3")

	if len(ws.GetMergeCells()) != 2 {
		t.Fatalf("expected 2 merges, got %d", len(ws.GetMergeCells()))
	}

	err := ws.RemoveMergedCell("A1:C3")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ws.GetMergeCells()) != 1 {
		t.Fatalf("expected 1 merge after removal, got %d", len(ws.GetMergeCells()))
	}

	// Remaining merge should be D1:F3
	mc := ws.GetMergeCells()[0]
	if mc.StartRow != 0 || mc.StartCol != 3 {
		t.Errorf("remaining merge start = (%d,%d), want (0,3)", mc.StartRow, mc.StartCol)
	}
}

func TestRemoveMergedCellNotFound(t *testing.T) {
	ws := NewWorksheet("Test")
	ws.MergeCells("A1:C3")

	err := ws.RemoveMergedCell("D1:F3")
	if err == nil {
		t.Error("expected error for non-existent merge")
	}
}

func TestRemoveMergedCellInvalidRange(t *testing.T) {
	ws := NewWorksheet("Test")
	err := ws.RemoveMergedCell("invalid")
	if err == nil {
		t.Error("expected error for invalid range")
	}
}

// ============================================================
// Worksheet: Row/Column visibility
// ============================================================

func TestRowHidden(t *testing.T) {
	ws := NewWorksheet("Test")

	if ws.IsRowHidden(0) {
		t.Error("row should not be hidden by default")
	}

	ws.SetRowHidden(0, true)
	if !ws.IsRowHidden(0) {
		t.Error("row should be hidden after SetRowHidden(true)")
	}

	ws.SetRowHidden(0, false)
	if ws.IsRowHidden(0) {
		t.Error("row should not be hidden after SetRowHidden(false)")
	}
}

func TestColumnHidden(t *testing.T) {
	ws := NewWorksheet("Test")

	if ws.IsColumnHidden(0) {
		t.Error("column should not be hidden by default")
	}

	ws.SetColumnHidden(0, true)
	if !ws.IsColumnHidden(0) {
		t.Error("column should be hidden after SetColumnHidden(true)")
	}

	ws.SetColumnHidden(0, false)
	if ws.IsColumnHidden(0) {
		t.Error("column should not be hidden after SetColumnHidden(false)")
	}
}

func TestMultipleRowsHidden(t *testing.T) {
	ws := NewWorksheet("Test")
	ws.SetRowHidden(1, true)
	ws.SetRowHidden(3, true)
	ws.SetRowHidden(5, true)

	if !ws.IsRowHidden(1) || !ws.IsRowHidden(3) || !ws.IsRowHidden(5) {
		t.Error("rows 1, 3, 5 should be hidden")
	}
	if ws.IsRowHidden(0) || ws.IsRowHidden(2) || ws.IsRowHidden(4) {
		t.Error("rows 0, 2, 4 should not be hidden")
	}
}

// ============================================================
// Worksheet: SheetView
// ============================================================

func TestSheetViewDefaults(t *testing.T) {
	ws := NewWorksheet("Test")
	sv := ws.GetSheetView()

	if sv.ZoomScale != 100 {
		t.Errorf("default zoom = %d, want 100", sv.ZoomScale)
	}
	if !sv.ShowGridlines {
		t.Error("gridlines should be shown by default")
	}
	if !sv.ShowRowColHeaders {
		t.Error("row/col headers should be shown by default")
	}
}

func TestSheetViewCustom(t *testing.T) {
	ws := NewWorksheet("Test")
	ws.SetSheetView(&SheetView{
		ZoomScale:         150,
		ShowGridlines:     false,
		ShowRowColHeaders: false,
		ShowRuler:         true,
		TopLeftCell:       "C5",
	})

	sv := ws.GetSheetView()
	if sv.ZoomScale != 150 {
		t.Errorf("zoom = %d, want 150", sv.ZoomScale)
	}
	if sv.ShowGridlines {
		t.Error("gridlines should be hidden")
	}
	if sv.TopLeftCell != "C5" {
		t.Errorf("TopLeftCell = %q, want 'C5'", sv.TopLeftCell)
	}
}

// ============================================================
// Worksheet: Sort
// ============================================================

func TestSortAscending(t *testing.T) {
	ws := NewWorksheet("Test")
	// Header row
	ws.SetCellValue("A1", "Name")
	ws.SetCellValue("B1", "Score")
	// Data rows
	ws.SetCellValue("A2", "Charlie")
	ws.SetCellValue("B2", 70)
	ws.SetCellValue("A3", "Alice")
	ws.SetCellValue("B3", 90)
	ws.SetCellValue("A4", "Bob")
	ws.SetCellValue("B4", 80)

	err := ws.Sort("A", 1, SortOrderAscending)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// After sort by column A ascending: Alice, Bob, Charlie
	v, _ := ws.GetCellValue("A2")
	if v != "Alice" {
		t.Errorf("A2 = %v, want 'Alice'", v)
	}
	v, _ = ws.GetCellValue("A3")
	if v != "Bob" {
		t.Errorf("A3 = %v, want 'Bob'", v)
	}
	v, _ = ws.GetCellValue("A4")
	if v != "Charlie" {
		t.Errorf("A4 = %v, want 'Charlie'", v)
	}
}

func TestSortDescending(t *testing.T) {
	ws := NewWorksheet("Test")
	ws.SetCellValue("A1", "Name")
	ws.SetCellValue("B1", "Score")
	ws.SetCellValue("A2", "Alice")
	ws.SetCellValue("B2", 90)
	ws.SetCellValue("A3", "Bob")
	ws.SetCellValue("B3", 80)
	ws.SetCellValue("A4", "Charlie")
	ws.SetCellValue("B4", 70)

	err := ws.Sort("B", 1, SortOrderDescending)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// After sort by column B descending: 90, 80, 70
	v, _ := ws.GetCellValue("B2")
	if v != float64(90) {
		t.Errorf("B2 = %v, want 90", v)
	}
	v, _ = ws.GetCellValue("B3")
	if v != float64(80) {
		t.Errorf("B3 = %v, want 80", v)
	}
	v, _ = ws.GetCellValue("B4")
	if v != float64(70) {
		t.Errorf("B4 = %v, want 70", v)
	}
}

func TestSortInvalidColumn(t *testing.T) {
	ws := NewWorksheet("Test")
	ws.SetCellValue("A1", "data")

	err := ws.Sort("123", 0, SortOrderAscending)
	if err == nil {
		t.Error("expected error for invalid column")
	}
}

func TestSortEmptySheet(t *testing.T) {
	ws := NewWorksheet("Test")
	err := ws.Sort("A", 0, SortOrderAscending)
	if err == nil {
		t.Error("expected error for empty sheet")
	}
}

func TestSortPreservesOtherColumns(t *testing.T) {
	ws := NewWorksheet("Test")
	ws.SetCellValue("A1", "Charlie")
	ws.SetCellValue("B1", 30)
	ws.SetCellValue("A2", "Alice")
	ws.SetCellValue("B2", 10)
	ws.SetCellValue("A3", "Bob")
	ws.SetCellValue("B3", 20)

	ws.Sort("A", 0, SortOrderAscending)

	// Alice should be in row 0 with score 10
	v, _ := ws.GetCellValue("A1")
	if v != "Alice" {
		t.Errorf("A1 = %v, want 'Alice'", v)
	}
	v, _ = ws.GetCellValue("B1")
	if v != float64(10) {
		t.Errorf("B1 = %v, want 10", v)
	}
}

// ============================================================
// Worksheet: SetBorderRange
// ============================================================

func TestSetBorderRange(t *testing.T) {
	ws := NewWorksheet("Test")
	ws.SetCellValue("A1", "a")
	ws.SetCellValue("B1", "b")
	ws.SetCellValue("A2", "c")
	ws.SetCellValue("B2", "d")

	border := Border{Style: BorderThin, Color: "000000"}
	err := ws.SetBorderRange("A1:B2", border)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Top-left corner (A1): should have top and left borders
	c, _ := ws.GetCellByName("A1")
	if c.Style == nil || c.Style.Borders == nil {
		t.Fatal("A1 should have borders")
	}
	if c.Style.Borders.Top.Style != BorderThin {
		t.Error("A1 should have top border")
	}
	if c.Style.Borders.Left.Style != BorderThin {
		t.Error("A1 should have left border")
	}

	// Top-right corner (B1): should have top and right borders
	c, _ = ws.GetCellByName("B1")
	if c.Style.Borders.Top.Style != BorderThin {
		t.Error("B1 should have top border")
	}
	if c.Style.Borders.Right.Style != BorderThin {
		t.Error("B1 should have right border")
	}

	// Bottom-left corner (A2): should have bottom and left borders
	c, _ = ws.GetCellByName("A2")
	if c.Style.Borders.Bottom.Style != BorderThin {
		t.Error("A2 should have bottom border")
	}
	if c.Style.Borders.Left.Style != BorderThin {
		t.Error("A2 should have left border")
	}

	// Bottom-right corner (B2): should have bottom and right borders
	c, _ = ws.GetCellByName("B2")
	if c.Style.Borders.Bottom.Style != BorderThin {
		t.Error("B2 should have bottom border")
	}
	if c.Style.Borders.Right.Style != BorderThin {
		t.Error("B2 should have right border")
	}
}

func TestSetBorderRangeInvalid(t *testing.T) {
	ws := NewWorksheet("Test")
	err := ws.SetBorderRange("invalid", Border{})
	if err == nil {
		t.Error("expected error for invalid range")
	}
}

// ============================================================
// Worksheet: Validate
// ============================================================

func TestWorksheetValidate(t *testing.T) {
	ws := NewWorksheet("Test")
	if err := ws.Validate(); err != nil {
		t.Errorf("valid worksheet should pass validation: %v", err)
	}

	ws.SetTitle("")
	if err := ws.Validate(); err == nil {
		t.Error("worksheet with empty title should fail validation")
	}
}

func TestWorksheetValidateOverlappingMerges(t *testing.T) {
	ws := NewWorksheet("Test")
	ws.MergeCells("A1:C3")
	// Manually add overlapping merge (bypassing normal checks)
	ws.mergeCells = append(ws.mergeCells, MergeCell{
		StartRow: 1, StartCol: 1, EndRow: 4, EndCol: 4,
	})

	if err := ws.Validate(); err == nil {
		t.Error("overlapping merges should fail validation")
	}
}

// ============================================================
// Workbook: RemoveNamedRange, Sheets, Validate, Close
// ============================================================

func TestWorkbookRemoveNamedRange(t *testing.T) {
	wb := New()
	wb.AddNamedRange("MyRange", "Sheet1!A1:B10")

	err := wb.RemoveNamedRange("MyRange")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(wb.GetNamedRanges()) != 0 {
		t.Error("named ranges should be empty after removal")
	}

	err = wb.RemoveNamedRange("NonExistent")
	if err == nil {
		t.Error("expected error for non-existent named range")
	}
}

func TestWorkbookSheets(t *testing.T) {
	wb := New()
	wb.AddSheet("Sheet2")
	wb.AddSheet("Sheet3")

	sheets := wb.Sheets()
	if len(sheets) != 3 {
		t.Fatalf("expected 3 sheets, got %d", len(sheets))
	}
	if sheets[0].Title() != "Sheet1" {
		t.Errorf("first sheet = %q, want 'Sheet1'", sheets[0].Title())
	}
	if sheets[2].Title() != "Sheet3" {
		t.Errorf("third sheet = %q, want 'Sheet3'", sheets[2].Title())
	}
}

func TestWorkbookValidate(t *testing.T) {
	wb := New()
	if err := wb.Validate(); err != nil {
		t.Errorf("valid workbook should pass: %v", err)
	}

	// Empty workbook
	wb2 := NewEmpty()
	if err := wb2.Validate(); err == nil {
		t.Error("empty workbook should fail validation")
	}
}

func TestWorkbookValidateDuplicateNames(t *testing.T) {
	wb := New()
	ws, _ := wb.AddSheet("Sheet2")
	// Force duplicate name
	ws.SetTitle("Sheet1")

	if err := wb.Validate(); err == nil {
		t.Error("workbook with duplicate sheet names should fail validation")
	}
}

func TestWorkbookClose(t *testing.T) {
	wb := New()
	if err := wb.Close(); err != nil {
		t.Errorf("Close should not return error: %v", err)
	}
}

// ============================================================
// Workbook: Protection
// ============================================================

func TestWorkbookProtectionSetAndClear(t *testing.T) {
	wb := New()
	if wb.GetWorkbookProtection() != nil {
		t.Error("protection should be nil initially")
	}

	wp := NewWorkbookProtection()
	wp.SetPassword("secret")
	wb.SetWorkbookProtection(wp)

	if wb.GetWorkbookProtection() == nil {
		t.Error("protection should not be nil after setting")
	}
	if !wb.GetWorkbookProtection().LockStructure {
		t.Error("LockStructure should be true by default")
	}

	wb.ClearProtection()
	if wb.GetWorkbookProtection() != nil {
		t.Error("protection should be nil after clearing")
	}
}

// ============================================================
// Workbook: CopySheet
// ============================================================

func TestWorkbookCopySheet(t *testing.T) {
	wb := New()
	ws := wb.GetActiveSheet()
	ws.SetCellValue("A1", "Hello")
	ws.SetCellValue("B1", 42)
	ws.SetCellStyle("A1", NewStyle().SetFont(&Font{Bold: true}))
	ws.MergeCells("A1:B1")
	ws.SetColumnWidth(0, 20.0)
	ws.SetRowHeight(0, 30.0)
	ws.SetTabColor("FF0000")
	ws.SetRowHidden(1, true)
	ws.SetColumnHidden(2, true)

	copy, err := wb.CopySheet(0, "Copy")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify copied values
	v, _ := copy.GetCellValue("A1")
	if v != "Hello" {
		t.Errorf("A1 = %v, want 'Hello'", v)
	}
	v, _ = copy.GetCellValue("B1")
	if v != float64(42) {
		t.Errorf("B1 = %v, want 42", v)
	}

	// Verify copied style
	c, _ := copy.GetCellByName("A1")
	if c.Style == nil || c.Style.Font == nil || !c.Style.Font.Bold {
		t.Error("copied cell should have bold font")
	}

	// Verify merge cells
	merges := copy.GetMergeCells()
	if len(merges) != 1 {
		t.Fatalf("expected 1 merge, got %d", len(merges))
	}

	// Verify dimensions
	if copy.GetColumnWidth(0) != 20.0 {
		t.Errorf("column width = %f, want 20.0", copy.GetColumnWidth(0))
	}
	if copy.GetRowHeight(0) != 30.0 {
		t.Errorf("row height = %f, want 30.0", copy.GetRowHeight(0))
	}

	// Verify tab color
	if copy.GetTabColor() != "FF0000" {
		t.Errorf("tab color = %q, want 'FF0000'", copy.GetTabColor())
	}

	// Verify hidden rows/cols
	if !copy.IsRowHidden(1) {
		t.Error("row 1 should be hidden in copy")
	}
	if !copy.IsColumnHidden(2) {
		t.Error("column 2 should be hidden in copy")
	}

	// Verify independence (modifying copy doesn't affect original)
	copy.SetCellValue("A1", "Modified")
	origV, _ := ws.GetCellValue("A1")
	if origV != "Hello" {
		t.Error("modifying copy should not affect original")
	}
}

func TestWorkbookCopySheetDuplicateName(t *testing.T) {
	wb := New()
	_, err := wb.CopySheet(0, "Sheet1")
	if err == nil {
		t.Error("expected error for duplicate sheet name")
	}
}

func TestWorkbookCopySheetInvalidIndex(t *testing.T) {
	wb := New()
	_, err := wb.CopySheet(99, "Copy")
	if err == nil {
		t.Error("expected error for invalid index")
	}
}

// ============================================================
// Workbook: MoveSheet
// ============================================================

func TestWorkbookMoveSheet(t *testing.T) {
	wb := New()
	wb.AddSheet("Sheet2")
	wb.AddSheet("Sheet3")

	err := wb.MoveSheet(2, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	names := wb.GetSheetNames()
	if names[0] != "Sheet3" || names[1] != "Sheet1" || names[2] != "Sheet2" {
		t.Errorf("sheet order = %v, want [Sheet3, Sheet1, Sheet2]", names)
	}
}

func TestWorkbookMoveSheetSameIndex(t *testing.T) {
	wb := New()
	wb.AddSheet("Sheet2")

	err := wb.MoveSheet(0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if wb.GetSheetNames()[0] != "Sheet1" {
		t.Error("sheet should not move when from == to")
	}
}

func TestWorkbookMoveSheetInvalidIndex(t *testing.T) {
	wb := New()
	if err := wb.MoveSheet(-1, 0); err == nil {
		t.Error("expected error for negative source index")
	}
	if err := wb.MoveSheet(0, 99); err == nil {
		t.Error("expected error for out of range destination index")
	}
}

// ============================================================
// Integration: Sort with styles preserved
// ============================================================

func TestSortPreservesStyles(t *testing.T) {
	ws := NewWorksheet("Test")
	ws.SetCellValue("A1", "Zebra")
	ws.SetCellStyle("A1", NewStyle().SetFont(&Font{Bold: true}))
	ws.SetCellValue("A2", "Apple")
	ws.SetCellStyle("A2", NewStyle().SetFont(&Font{Italic: true}))

	ws.Sort("A", 0, SortOrderAscending)

	// Apple should now be in row 0 with italic style
	c, _ := ws.GetCellByName("A1")
	if c.GetStringValue() != "Apple" {
		t.Errorf("A1 = %q, want 'Apple'", c.GetStringValue())
	}
	if c.Style == nil || c.Style.Font == nil || !c.Style.Font.Italic {
		t.Error("Apple cell should have italic font after sort")
	}

	// Zebra should now be in row 1 with bold style
	c, _ = ws.GetCellByName("A2")
	if c.GetStringValue() != "Zebra" {
		t.Errorf("A2 = %q, want 'Zebra'", c.GetStringValue())
	}
	if c.Style == nil || c.Style.Font == nil || !c.Style.Font.Bold {
		t.Error("Zebra cell should have bold font after sort")
	}
}

// ============================================================
// Integration: Cell type transitions
// ============================================================

func TestCellTypeTransitions(t *testing.T) {
	c := NewCell(0, 0)

	// Empty -> String
	c.SetValue("hello")
	if !c.IsString() {
		t.Error("should be string")
	}

	// String -> Number
	c.SetValue(42)
	if !c.IsNumber() {
		t.Error("should be number")
	}

	// Number -> Bool
	c.SetValue(true)
	if !c.IsBool() {
		t.Error("should be bool")
	}

	// Bool -> Date
	c.SetValue(time.Now())
	if !c.IsDate() {
		t.Error("should be date")
	}

	// Date -> Formula
	c.SetFormula("SUM(A1:A10)")
	if !c.HasFormula() {
		t.Error("should have formula")
	}

	// Formula -> Empty
	c.Clear()
	if !c.IsEmpty() {
		t.Error("should be empty")
	}
}

// ============================================================
// Edge cases
// ============================================================

func TestCellClearAndReuse(t *testing.T) {
	c := NewCell(0, 0)
	c.SetValue("first")
	c.Clear()
	c.SetValue("second")

	if c.GetStringValue() != "second" {
		t.Errorf("value = %q, want 'second'", c.GetStringValue())
	}
	if !c.IsString() {
		t.Error("should be string after reuse")
	}
}

func TestSetBorderRangeSingleCell(t *testing.T) {
	ws := NewWorksheet("Test")
	border := Border{Style: BorderThick, Color: "FF0000"}
	err := ws.SetBorderRange("A1:A1", border)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	c := ws.GetCell(0, 0)
	if c.Style == nil || c.Style.Borders == nil {
		t.Fatal("cell should have borders")
	}
	b := c.Style.Borders
	if b.Top.Style != BorderThick || b.Bottom.Style != BorderThick ||
		b.Left.Style != BorderThick || b.Right.Style != BorderThick {
		t.Error("single cell should have all four borders")
	}
}

func TestSortWithEmptyCells(t *testing.T) {
	ws := NewWorksheet("Test")
	ws.SetCellValue("A1", "Zebra")
	ws.SetCellValue("A3", "Apple")
	// A2 is empty

	err := ws.Sort("A", 0, SortOrderAscending)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetFormattedValueNoStyle(t *testing.T) {
	c := NewCell(0, 0)
	c.SetValue(42)
	// No style set - should return default string representation
	got := c.GetFormattedValue()
	if got != "42" {
		t.Errorf("GetFormattedValue() = %q, want '42'", got)
	}
}

func TestGetFormattedValueTextFormat(t *testing.T) {
	c := NewCell(0, 0)
	c.SetValue(42.0)
	c.SetStyle(NewStyle().SetNumberFormat(&NumberFormat{FormatCode: "@"}))
	got := c.GetFormattedValue()
	if got != "42" {
		t.Errorf("GetFormattedValue() = %q, want '42'", got)
	}
}
