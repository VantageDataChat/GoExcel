> [English](README_en.md) | [中文](README.md)

# GoExcel

A pure Go spreadsheet processing library, inspired by [PHPOffice/PhpSpreadsheet](https://github.com/PHPOffice/PhpSpreadsheet).

Supports reading and writing XLSX and CSV formats, provides an in-memory spreadsheet object model with a formula calculation engine, styling system, conditional formatting, data validation, and more.

## Features

- XLSX read/write (Open XML format, compatible with Excel 2007+)
- CSV read/write (custom delimiters)
- Formula calculation engine (24 built-in functions)
- Complete styling system (font, border, fill, alignment, number format)
- Merged cells, freeze panes
- Row/column insert, delete, copy
- Hyperlinks, comments, rich text
- Conditional formatting, data validation, auto filter
- Page setup and print configuration
- Worksheet/workbook protection
- Document properties, named ranges
- Unicode and special character support
- 97.7% test coverage, 256 test cases

## Installation

```bash
go get github.com/VantageDataChat/GoExcel
```

## Quick Start

### Create and Write XLSX

```go
package main

import "github.com/VantageDataChat/GoExcel"

func main() {
    wb := gospreadsheet.New()
    ws := wb.GetActiveSheet()

    // Set cell values
    ws.SetCellValue("A1", "Name")
    ws.SetCellValue("B1", "Score")
    ws.SetCellValue("A2", "Alice")
    ws.SetCellValue("B2", 95.5)
    ws.SetCellValue("A3", "Bob")
    ws.SetCellValue("B3", 87)

    // Add formula
    ws.SetCellFormula("B4", "AVERAGE(B2:B3)")

    // Set style
    ws.SetCellStyle("A1", gospreadsheet.NewStyle().
        SetFont(&gospreadsheet.Font{Bold: true, Size: 14}))

    // Save
    gospreadsheet.SaveFile(wb, "output.xlsx")
}
```

### Read XLSX

```go
wb, err := gospreadsheet.OpenFile("input.xlsx")
if err != nil {
    log.Fatal(err)
}

ws := wb.GetActiveSheet()
val, _ := ws.GetCellValue("A1")
fmt.Println(val)

// Iterate all rows
rows, _ := ws.RowIterator()
for _, row := range rows {
    for _, cell := range row {
        if cell != nil {
            fmt.Print(cell.GetStringValue(), "\t")
        }
    }
    fmt.Println()
}
```

### CSV Read/Write

```go
// Read CSV
wb, _ := gospreadsheet.OpenFile("data.csv")

// Custom delimiter
reader := gospreadsheet.NewCSVReader()
reader.Delimiter = ';'
wb, _ = reader.Open("data.csv")

// Write CSV
gospreadsheet.SaveFile(wb, "output.csv")
```

### Formula Calculation

```go
wb := gospreadsheet.New()
ws := wb.GetActiveSheet()
ws.SetCellValue("A1", 100)
ws.SetCellValue("A2", 200)
ws.SetCellFormula("A3", "SUM(A1:A2)")

ce := gospreadsheet.NewCalculationEngine(wb)
result, _ := ce.CalculateCell(ws, "A3")
fmt.Println(result) // 300

// Calculate all formulas
ce.CalculateAll()
```

### Styling

```go
style := gospreadsheet.NewStyle().
    SetFont(&gospreadsheet.Font{
        Name: "Arial", Size: 12, Bold: true, Color: "FF0000",
    }).
    SetFill(&gospreadsheet.Fill{
        Type: "solid", Color: "FFFF00",
    }).
    SetBorders(&gospreadsheet.Borders{
        Bottom: gospreadsheet.Border{
            Style: gospreadsheet.BorderThin, Color: "000000",
        },
    }).
    SetAlignment(&gospreadsheet.Alignment{
        Horizontal: gospreadsheet.AlignCenter,
        WrapText:   true,
    }).
    SetNumberFormat(&gospreadsheet.FormatPercent2Dec)

ws.SetCellStyle("A1", style)
```

### Multiple Worksheets

```go
wb := gospreadsheet.New()
ws1 := wb.GetActiveSheet() // Sheet1

ws2, _ := wb.AddSheet("Data")
ws2.SetCellValue("A1", "Second worksheet")

ws3, _ := wb.AddSheet("Summary")
wb.SetActiveSheet(2) // Switch to the third worksheet
```

### Merged Cells and Freeze Panes

```go
ws.SetCellValue("A1", "Merged Title")
ws.MergeCells("A1:D1")
ws.FreezePane("A2") // Freeze the first row
```

### Row and Column Operations

```go
ws.InsertRow(2)       // Insert an empty row before row 3
ws.DeleteRow(5)       // Delete row 6
ws.InsertColumn(1)    // Insert an empty column before column B
ws.DeleteColumn(3)    // Delete column D
ws.CopyRow(0, 10)     // Copy row 1 to row 11
ws.SetColumnWidth(0, 20.0)
ws.SetRowHeight(0, 30.0)
```

### Hyperlinks and Comments

```go
ws.SetCellValue("A1", "Click to visit")
ws.SetCellHyperlink("A1", "https://example.com")
ws.SetCellComment("B1", "Reviewer", "Please check this data")
```

### Conditional Formatting

```go
cf := gospreadsheet.NewConditionalFormatting("B2:B100")
cf.AddRule(gospreadsheet.CellIsRule(
    gospreadsheet.OperatorGreaterThan, "90",
    gospreadsheet.NewStyle().SetFont(&gospreadsheet.Font{Color: "00AA00"}),
))
cf.AddRule(gospreadsheet.BetweenRule("60", "90",
    gospreadsheet.NewStyle().SetFont(&gospreadsheet.Font{Color: "FF8800"}),
))
ws.AddConditionalFormatting(cf)
```

### Data Validation

```go
// Dropdown list
dv := gospreadsheet.NewDataValidation("C2:C100")
dv.SetListValues([]string{"Pass", "Fail", "Pending"})
dv.SetErrorMessage("Error", "Please select from the list")
ws.AddDataValidation(dv)

// Numeric range
dv2 := gospreadsheet.NewDataValidation("D2:D100").
    SetType(gospreadsheet.ValidationWhole).
    SetOperator(gospreadsheet.ValOperatorBetween).
    SetFormula1("0").SetFormula2("100")
ws.AddDataValidation(dv2)
```

### Page Setup

```go
ps := ws.GetPageSetup()
ps.SetOrientation(gospreadsheet.OrientationLandscape)
ps.SetPaperSize(gospreadsheet.PaperA4)
ps.SetScale(85)
ps.SetPrintArea("A1:H50")
ps.SetRepeatRows("1:1") // Repeat header row on each page
```

### Worksheet Protection

```go
sp := gospreadsheet.NewSheetProtection().SetPassword("secret")
sp.AllowSort()
sp.AllowAutoFilter()
ws.SetSheetProtection(sp)
```

## Supported Formula Functions

| Category | Functions |
|----------|-----------|
| Math | SUM, ABS, ROUND, SQRT, POWER, MOD, INT |
| Statistics | AVERAGE, COUNT, COUNTA, MAX, MIN, MEDIAN |
| Logic | IF |
| Text | LEN, UPPER, LOWER, TRIM, LEFT, RIGHT, MID, CONCATENATE |
| Conditional | SUMIF, COUNTIF |

Supports cell references (A1), range references (A1:A10), nested formulas, and arithmetic operations.

## Project Structure

```
goexcel/
├── spreadsheet.go     # Cell type and operations
├── workbook.go        # Workbook management
├── worksheet.go       # Worksheet
├── coordinate.go      # Cell coordinate parsing
├── style.go           # Styling system
├── calculation.go     # Formula calculation engine
├── functions.go       # Built-in function implementations
├── xlsx_writer.go     # XLSX writer
├── xlsx_reader.go     # XLSX reader
├── csv.go             # CSV read/write
├── iofactory.go       # Format auto-detection
├── hyperlink.go       # Hyperlinks
├── comment.go         # Comments
├── richtext.go        # Rich text
├── conditional.go     # Conditional formatting
├── validation.go      # Data validation
├── autofilter.go      # Auto filter
├── pagesetup.go       # Page setup
└── protection.go      # Worksheet/workbook protection
```

## Testing

```bash
go test ./... -v
go test ./... -cover  # Coverage 97.7%
```

## License

MIT
