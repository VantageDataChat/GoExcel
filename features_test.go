package gospreadsheet

import (
	"testing"
)

// --- Hyperlink Tests ---

func TestHyperlink(t *testing.T) {
	h := NewHyperlink("https://example.com").SetTooltip("Example")
	if h.URL != "https://example.com" {
		t.Errorf("URL = %q, want 'https://example.com'", h.URL)
	}
	if h.Tooltip != "Example" {
		t.Errorf("Tooltip = %q, want 'Example'", h.Tooltip)
	}
}

func TestCellHyperlink(t *testing.T) {
	c := NewCell(0, 0)
	c.SetValue("Click me")
	c.SetHyperlink(NewHyperlink("https://example.com"))

	if c.Hyperlink == nil {
		t.Fatal("hyperlink should not be nil")
	}
	if c.Hyperlink.URL != "https://example.com" {
		t.Errorf("URL = %q", c.Hyperlink.URL)
	}
}

func TestWorksheetSetCellHyperlink(t *testing.T) {
	ws := NewWorksheet("Test")
	ws.SetCellValue("A1", "Link")
	ws.SetCellHyperlink("A1", "https://example.com")

	c, _ := ws.GetCellByName("A1")
	if c.Hyperlink == nil || c.Hyperlink.URL != "https://example.com" {
		t.Error("hyperlink not set correctly")
	}
}

// --- Comment Tests ---

func TestComment(t *testing.T) {
	c := NewComment("Author", "This is a comment")
	if c.Author != "Author" {
		t.Errorf("Author = %q", c.Author)
	}
	if c.Text != "This is a comment" {
		t.Errorf("Text = %q", c.Text)
	}
}

func TestCellComment(t *testing.T) {
	cell := NewCell(0, 0)
	cell.SetComment(NewComment("User", "Note"))

	if cell.Comment == nil {
		t.Fatal("comment should not be nil")
	}
	if cell.Comment.Author != "User" {
		t.Errorf("Author = %q", cell.Comment.Author)
	}
}

func TestWorksheetSetCellComment(t *testing.T) {
	ws := NewWorksheet("Test")
	ws.SetCellValue("A1", "Data")
	ws.SetCellComment("A1", "Admin", "Important note")

	c, _ := ws.GetCellByName("A1")
	if c.Comment == nil || c.Comment.Text != "Important note" {
		t.Error("comment not set correctly")
	}
}

// --- RichText Tests ---

func TestRichText(t *testing.T) {
	rt := NewRichText().
		AddRun("Hello ", &Font{Bold: true}).
		AddRun("World", &Font{Italic: true})

	if len(rt.Runs) != 2 {
		t.Fatalf("runs = %d, want 2", len(rt.Runs))
	}
	if rt.PlainText() != "Hello World" {
		t.Errorf("PlainText = %q, want 'Hello World'", rt.PlainText())
	}
	if !rt.Runs[0].Font.Bold {
		t.Error("first run should be bold")
	}
	if !rt.Runs[1].Font.Italic {
		t.Error("second run should be italic")
	}
}

func TestCellRichText(t *testing.T) {
	cell := NewCell(0, 0)
	rt := NewRichText().AddRun("Test", nil)
	cell.SetRichText(rt)

	if cell.RichText == nil {
		t.Fatal("richtext should not be nil")
	}
	if cell.Type != CellTypeString {
		t.Errorf("type = %d, want CellTypeString", cell.Type)
	}
	if cell.Value != "Test" {
		t.Errorf("value = %v, want 'Test'", cell.Value)
	}
}

// --- Conditional Formatting Tests ---

func TestConditionalFormatting(t *testing.T) {
	cf := NewConditionalFormatting("A1:A10")
	style := NewStyle().SetFont(&Font{Bold: true, Color: "FF0000"})

	cf.AddRule(CellIsRule(OperatorGreaterThan, "50", style))
	cf.AddRule(BetweenRule("10", "50", style))
	cf.AddRule(ExpressionRule("A1>0", style))

	if cf.Range != "A1:A10" {
		t.Errorf("Range = %q", cf.Range)
	}
	if len(cf.Rules) != 3 {
		t.Fatalf("rules = %d, want 3", len(cf.Rules))
	}
	if cf.Rules[0].Type != ConditionalCellIs {
		t.Errorf("rule[0] type = %q", cf.Rules[0].Type)
	}
	if cf.Rules[0].Operator != OperatorGreaterThan {
		t.Errorf("rule[0] operator = %q", cf.Rules[0].Operator)
	}
	if cf.Rules[1].Operator != OperatorBetween {
		t.Errorf("rule[1] operator = %q", cf.Rules[1].Operator)
	}
	if cf.Rules[2].Type != ConditionalExpression {
		t.Errorf("rule[2] type = %q", cf.Rules[2].Type)
	}
}

func TestConditionalFormattingPriority(t *testing.T) {
	cf := NewConditionalFormatting("A1:A10")
	cf.AddRule(ConditionalRule{Type: ConditionalCellIs})
	cf.AddRule(ConditionalRule{Type: ConditionalCellIs})

	if cf.Rules[0].Priority != 1 {
		t.Errorf("rule[0] priority = %d, want 1", cf.Rules[0].Priority)
	}
	if cf.Rules[1].Priority != 2 {
		t.Errorf("rule[1] priority = %d, want 2", cf.Rules[1].Priority)
	}
}

func TestWorksheetConditionalFormatting(t *testing.T) {
	ws := NewWorksheet("Test")
	cf := NewConditionalFormatting("A1:A10")
	ws.AddConditionalFormatting(cf)

	cfs := ws.GetConditionalFormattings()
	if len(cfs) != 1 {
		t.Errorf("conditional formattings = %d, want 1", len(cfs))
	}
}

// --- Data Validation Tests ---

func TestDataValidation(t *testing.T) {
	dv := NewDataValidation("A1:A100").
		SetType(ValidationWhole).
		SetOperator(ValOperatorBetween).
		SetFormula1("1").
		SetFormula2("100").
		SetErrorMessage("Error", "Value must be between 1 and 100").
		SetPromptMessage("Input", "Enter a number 1-100")

	if dv.Range != "A1:A100" {
		t.Errorf("Range = %q", dv.Range)
	}
	if dv.Type != ValidationWhole {
		t.Errorf("Type = %q", dv.Type)
	}
	if dv.Operator != ValOperatorBetween {
		t.Errorf("Operator = %q", dv.Operator)
	}
	if dv.Formula1 != "1" || dv.Formula2 != "100" {
		t.Errorf("Formulas = %q, %q", dv.Formula1, dv.Formula2)
	}
	if dv.ErrorTitle != "Error" {
		t.Errorf("ErrorTitle = %q", dv.ErrorTitle)
	}
	if dv.PromptTitle != "Input" {
		t.Errorf("PromptTitle = %q", dv.PromptTitle)
	}
}

func TestDataValidationList(t *testing.T) {
	dv := NewDataValidation("B1:B10")
	dv.SetListValues([]string{"Red", "Green", "Blue"})

	if dv.Type != ValidationList {
		t.Errorf("Type = %q, want list", dv.Type)
	}
	if dv.Formula1 != `"Red","Green","Blue"` {
		t.Errorf("Formula1 = %q", dv.Formula1)
	}
}

func TestWorksheetDataValidation(t *testing.T) {
	ws := NewWorksheet("Test")
	dv := NewDataValidation("A1:A10").SetType(ValidationDecimal)
	ws.AddDataValidation(dv)

	dvs := ws.GetDataValidations()
	if len(dvs) != 1 {
		t.Errorf("data validations = %d, want 1", len(dvs))
	}
}

// --- AutoFilter Tests ---

func TestAutoFilter(t *testing.T) {
	af := NewAutoFilter("A1:D100")
	af.AddValueFilter(0, []string{"Apple", "Banana"})
	af.AddCustomFilter(1,
		FilterCondition{Operator: FilterOpGreaterThan, Value: "10"},
	)

	if af.Range != "A1:D100" {
		t.Errorf("Range = %q", af.Range)
	}
	if len(af.Columns) != 2 {
		t.Fatalf("columns = %d, want 2", len(af.Columns))
	}
	if af.Columns[0].FilterType != FilterValues {
		t.Errorf("col[0] type = %q", af.Columns[0].FilterType)
	}
	if len(af.Columns[0].Values) != 2 {
		t.Errorf("col[0] values = %d", len(af.Columns[0].Values))
	}
	if af.Columns[1].FilterType != FilterCustom {
		t.Errorf("col[1] type = %q", af.Columns[1].FilterType)
	}
}

func TestWorksheetAutoFilter(t *testing.T) {
	ws := NewWorksheet("Test")
	if ws.GetAutoFilter() != nil {
		t.Error("auto filter should be nil initially")
	}

	af := NewAutoFilter("A1:C10")
	ws.SetAutoFilter(af)

	if ws.GetAutoFilter() == nil {
		t.Error("auto filter should not be nil")
	}
	if ws.GetAutoFilter().Range != "A1:C10" {
		t.Errorf("Range = %q", ws.GetAutoFilter().Range)
	}
}

// --- Page Setup Tests ---

func TestPageSetup(t *testing.T) {
	ps := NewPageSetup().
		SetPaperSize(PaperLetter).
		SetOrientation(OrientationLandscape).
		SetScale(75).
		SetFitToPage(1, 0).
		SetPrintArea("A1:H50").
		SetRepeatRows("1:2").
		SetRepeatColumns("A:B")

	if ps.PaperSize != PaperLetter {
		t.Errorf("PaperSize = %d", ps.PaperSize)
	}
	if ps.Orientation != OrientationLandscape {
		t.Errorf("Orientation = %q", ps.Orientation)
	}
	if ps.Scale != 75 {
		t.Errorf("Scale = %d", ps.Scale)
	}
	if ps.FitToWidth != 1 || ps.FitToHeight != 0 {
		t.Errorf("FitTo = %d x %d", ps.FitToWidth, ps.FitToHeight)
	}
	if ps.PrintArea == nil || ps.PrintArea.Range != "A1:H50" {
		t.Error("PrintArea not set correctly")
	}
	if ps.RepeatRows != "1:2" {
		t.Errorf("RepeatRows = %q", ps.RepeatRows)
	}
	if ps.RepeatColumns != "A:B" {
		t.Errorf("RepeatColumns = %q", ps.RepeatColumns)
	}
}

func TestPageSetupScaleBounds(t *testing.T) {
	ps := NewPageSetup()
	ps.SetScale(5)
	if ps.Scale != 10 {
		t.Errorf("Scale = %d, want 10 (min)", ps.Scale)
	}
	ps.SetScale(500)
	if ps.Scale != 400 {
		t.Errorf("Scale = %d, want 400 (max)", ps.Scale)
	}
}

func TestPageMargins(t *testing.T) {
	m := DefaultPageMargins()
	if m.Top != 0.75 || m.Bottom != 0.75 {
		t.Errorf("margins top/bottom = %f/%f", m.Top, m.Bottom)
	}
	if m.Left != 0.7 || m.Right != 0.7 {
		t.Errorf("margins left/right = %f/%f", m.Left, m.Right)
	}
}

func TestWorksheetPageSetup(t *testing.T) {
	ws := NewWorksheet("Test")
	ps := ws.GetPageSetup()
	if ps == nil {
		t.Fatal("page setup should not be nil")
	}
	if ps.PaperSize != PaperA4 {
		t.Errorf("default paper = %d, want A4", ps.PaperSize)
	}
}

// --- Sheet Protection Tests ---

func TestSheetProtection(t *testing.T) {
	sp := NewSheetProtection()
	if !sp.Sheet {
		t.Error("Sheet should be true by default")
	}
	if !sp.Objects {
		t.Error("Objects should be true by default")
	}

	sp.SetPassword("secret123")
	if sp.Password == "" {
		t.Error("password should be hashed")
	}
	if sp.Password == "secret123" {
		t.Error("password should not be stored in plain text")
	}
}

func TestSheetProtectionAllow(t *testing.T) {
	sp := NewSheetProtection().
		AllowFormatCells().
		AllowInsertRows().
		AllowDeleteRows().
		AllowSort().
		AllowAutoFilter()

	if sp.FormatCells {
		t.Error("FormatCells should be false (allowed)")
	}
	if sp.InsertRows {
		t.Error("InsertRows should be false (allowed)")
	}
}

func TestWorksheetProtection(t *testing.T) {
	ws := NewWorksheet("Test")
	if ws.GetSheetProtection() != nil {
		t.Error("protection should be nil initially")
	}

	sp := NewSheetProtection()
	ws.SetSheetProtection(sp)
	if ws.GetSheetProtection() == nil {
		t.Error("protection should not be nil")
	}
}

// --- Tab Color Tests ---

func TestWorksheetTabColor(t *testing.T) {
	ws := NewWorksheet("Test")
	if ws.GetTabColor() != "" {
		t.Error("tab color should be empty initially")
	}

	ws.SetTabColor("FF0000")
	if ws.GetTabColor() != "FF0000" {
		t.Errorf("tab color = %q, want 'FF0000'", ws.GetTabColor())
	}
}

// --- Insert/Delete Row/Column Tests ---

func TestInsertRow(t *testing.T) {
	ws := NewWorksheet("Test")
	ws.SetCellValue("A1", "row1")
	ws.SetCellValue("A2", "row2")
	ws.SetCellValue("A3", "row3")

	ws.InsertRow(1) // Insert before row 2

	v, _ := ws.GetCellValue("A1")
	if v != "row1" {
		t.Errorf("A1 = %v, want 'row1'", v)
	}

	// row2 should now be at A3
	c := ws.GetCell(2, 0)
	if c.GetStringValue() != "row2" {
		t.Errorf("A3 = %q, want 'row2'", c.GetStringValue())
	}

	// row3 should now be at A4
	c = ws.GetCell(3, 0)
	if c.GetStringValue() != "row3" {
		t.Errorf("A4 = %q, want 'row3'", c.GetStringValue())
	}
}

func TestDeleteRow(t *testing.T) {
	ws := NewWorksheet("Test")
	ws.SetCellValue("A1", "row1")
	ws.SetCellValue("A2", "row2")
	ws.SetCellValue("A3", "row3")

	ws.DeleteRow(1) // Delete row 2

	v, _ := ws.GetCellValue("A1")
	if v != "row1" {
		t.Errorf("A1 = %v, want 'row1'", v)
	}

	// row3 should now be at A2
	c := ws.GetCell(1, 0)
	if c.GetStringValue() != "row3" {
		t.Errorf("A2 = %q, want 'row3'", c.GetStringValue())
	}

	if ws.CellCount() != 2 {
		t.Errorf("cell count = %d, want 2", ws.CellCount())
	}
}

func TestInsertColumn(t *testing.T) {
	ws := NewWorksheet("Test")
	ws.SetCellValue("A1", "col1")
	ws.SetCellValue("B1", "col2")
	ws.SetCellValue("C1", "col3")

	ws.InsertColumn(1) // Insert before column B

	c := ws.GetCell(0, 0)
	if c.GetStringValue() != "col1" {
		t.Errorf("A1 = %q, want 'col1'", c.GetStringValue())
	}

	// col2 should now be at C1
	c = ws.GetCell(0, 2)
	if c.GetStringValue() != "col2" {
		t.Errorf("C1 = %q, want 'col2'", c.GetStringValue())
	}

	// col3 should now be at D1
	c = ws.GetCell(0, 3)
	if c.GetStringValue() != "col3" {
		t.Errorf("D1 = %q, want 'col3'", c.GetStringValue())
	}
}

func TestDeleteColumn(t *testing.T) {
	ws := NewWorksheet("Test")
	ws.SetCellValue("A1", "col1")
	ws.SetCellValue("B1", "col2")
	ws.SetCellValue("C1", "col3")

	ws.DeleteColumn(1) // Delete column B

	c := ws.GetCell(0, 0)
	if c.GetStringValue() != "col1" {
		t.Errorf("A1 = %q, want 'col1'", c.GetStringValue())
	}

	// col3 should now be at B1
	c = ws.GetCell(0, 1)
	if c.GetStringValue() != "col3" {
		t.Errorf("B1 = %q, want 'col3'", c.GetStringValue())
	}

	if ws.CellCount() != 2 {
		t.Errorf("cell count = %d, want 2", ws.CellCount())
	}
}

func TestCopyRow(t *testing.T) {
	ws := NewWorksheet("Test")
	ws.SetCellValue("A1", "hello")
	ws.SetCellValue("B1", 42)

	ws.CopyRow(0, 2) // Copy row 1 to row 3

	c := ws.GetCell(2, 0)
	if c.GetStringValue() != "hello" {
		t.Errorf("A3 = %q, want 'hello'", c.GetStringValue())
	}
	c = ws.GetCell(2, 1)
	v, _ := c.GetNumericValue()
	if v != 42 {
		t.Errorf("B3 = %f, want 42", v)
	}
}

func TestInsertRowShiftsMergeCells(t *testing.T) {
	ws := NewWorksheet("Test")
	ws.MergeCells("A2:C2")
	ws.InsertRow(1)

	merges := ws.GetMergeCells()
	if len(merges) != 1 {
		t.Fatalf("merges = %d, want 1", len(merges))
	}
	if merges[0].StartRow != 2 {
		t.Errorf("merge start row = %d, want 2", merges[0].StartRow)
	}
}

func TestInsertRowShiftsRowHeights(t *testing.T) {
	ws := NewWorksheet("Test")
	ws.SetRowHeight(2, 30.0)
	ws.InsertRow(1)

	if ws.GetRowHeight(3) != 30.0 {
		t.Errorf("row height = %f, want 30", ws.GetRowHeight(3))
	}
}

func TestDeleteColumnShiftsWidths(t *testing.T) {
	ws := NewWorksheet("Test")
	ws.SetColumnWidth(2, 25.0)
	ws.DeleteColumn(0)

	if ws.GetColumnWidth(1) != 25.0 {
		t.Errorf("col width = %f, want 25", ws.GetColumnWidth(1))
	}
}

// --- Workbook Protection Tests ---

func TestWorkbookProtection(t *testing.T) {
	wp := NewWorkbookProtection()
	if !wp.LockStructure {
		t.Error("LockStructure should be true by default")
	}

	wp.SetPassword("test")
	if wp.Password == "" || wp.Password == "test" {
		t.Error("password should be hashed")
	}
}
