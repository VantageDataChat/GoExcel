package gospreadsheet

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

// XLSXWriter writes a Workbook to XLSX format.
type XLSXWriter struct{}

// NewXLSXWriter creates a new XLSX writer.
func NewXLSXWriter() *XLSXWriter {
	return &XLSXWriter{}
}

// Save writes the workbook to a file.
func (w *XLSXWriter) Save(wb *Workbook, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("creating file: %w", err)
	}

	writeErr := w.Write(wb, f)
	closeErr := f.Close()

	if writeErr != nil {
		return writeErr
	}
	if closeErr != nil {
		return fmt.Errorf("closing file: %w", closeErr)
	}
	return nil
}

// Write writes the workbook to an io.Writer.
func (w *XLSXWriter) Write(wb *Workbook, writer io.Writer) error {
	zw := zip.NewWriter(writer)
	defer zw.Close()

	// Collect shared strings
	sharedStrings := newSharedStrings()
	for _, ws := range wb.worksheets {
		for _, c := range ws.cells {
			if c.Type == CellTypeString {
				if s, ok := c.Value.(string); ok {
					sharedStrings.Add(s)
				}
			}
		}
	}

	// [Content_Types].xml
	if err := w.writeContentTypes(zw, wb); err != nil {
		return err
	}
	// _rels/.rels
	if err := w.writeRels(zw); err != nil {
		return err
	}
	// xl/_rels/workbook.xml.rels
	if err := w.writeWorkbookRels(zw, wb); err != nil {
		return err
	}
	// xl/workbook.xml
	if err := w.writeWorkbook(zw, wb); err != nil {
		return err
	}
	// xl/styles.xml
	if err := w.writeStyles(zw, wb); err != nil {
		return err
	}
	// xl/sharedStrings.xml
	if err := w.writeSharedStrings(zw, sharedStrings); err != nil {
		return err
	}
	// xl/worksheets/sheet{n}.xml
	for i, ws := range wb.worksheets {
		if err := w.writeSheet(zw, ws, i+1, sharedStrings); err != nil {
			return err
		}
	}
	// docProps/core.xml
	if err := w.writeCoreProperties(zw, wb); err != nil {
		return err
	}

	return nil
}

// sharedStrings manages the shared string table for XLSX.
type sharedStrings struct {
	strings []string
	index   map[string]int
}

func newSharedStrings() *sharedStrings {
	return &sharedStrings{
		strings: make([]string, 0),
		index:   make(map[string]int),
	}
}

func (ss *sharedStrings) Add(s string) int {
	if idx, ok := ss.index[s]; ok {
		return idx
	}
	idx := len(ss.strings)
	ss.strings = append(ss.strings, s)
	ss.index[s] = idx
	return idx
}

func (ss *sharedStrings) GetIndex(s string) int {
	if idx, ok := ss.index[s]; ok {
		return idx
	}
	return -1
}

func (w *XLSXWriter) writeContentTypes(zw *zip.Writer, wb *Workbook) error {
	f, err := zw.Create("[Content_Types].xml")
	if err != nil {
		return err
	}
	ct := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` + "\n"
	ct += `<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">`
	ct += `<Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>`
	ct += `<Default Extension="xml" ContentType="application/xml"/>`
	ct += `<Override PartName="/xl/workbook.xml" ContentType="application/vnd.openxmlformats-officedocument.spreadsheetml.sheet.main+xml"/>`
	ct += `<Override PartName="/xl/styles.xml" ContentType="application/vnd.openxmlformats-officedocument.spreadsheetml.styles+xml"/>`
	ct += `<Override PartName="/xl/sharedStrings.xml" ContentType="application/vnd.openxmlformats-officedocument.spreadsheetml.sharedStrings+xml"/>`
	ct += `<Override PartName="/docProps/core.xml" ContentType="application/vnd.openxmlformats-package.core-properties+xml"/>`
	for i := range wb.worksheets {
		ct += fmt.Sprintf(`<Override PartName="/xl/worksheets/sheet%d.xml" ContentType="application/vnd.openxmlformats-officedocument.spreadsheetml.worksheet+xml"/>`, i+1)
	}
	ct += `</Types>`
	_, err = f.Write([]byte(ct))
	return err
}

func (w *XLSXWriter) writeRels(zw *zip.Writer) error {
	f, err := zw.Create("_rels/.rels")
	if err != nil {
		return err
	}
	rels := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` + "\n"
	rels += `<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">`
	rels += `<Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="xl/workbook.xml"/>`
	rels += `<Relationship Id="rId2" Type="http://schemas.openxmlformats.org/package/2006/relationships/metadata/core-properties" Target="docProps/core.xml"/>`
	rels += `</Relationships>`
	_, err = f.Write([]byte(rels))
	return err
}

func (w *XLSXWriter) writeWorkbookRels(zw *zip.Writer, wb *Workbook) error {
	f, err := zw.Create("xl/_rels/workbook.xml.rels")
	if err != nil {
		return err
	}
	rels := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` + "\n"
	rels += `<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">`
	for i := range wb.worksheets {
		rels += fmt.Sprintf(`<Relationship Id="rId%d" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/worksheet" Target="worksheets/sheet%d.xml"/>`, i+1, i+1)
	}
	rels += fmt.Sprintf(`<Relationship Id="rId%d" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/styles" Target="styles.xml"/>`, len(wb.worksheets)+1)
	rels += fmt.Sprintf(`<Relationship Id="rId%d" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/sharedStrings" Target="sharedStrings.xml"/>`, len(wb.worksheets)+2)
	rels += `</Relationships>`
	_, err = f.Write([]byte(rels))
	return err
}

func (w *XLSXWriter) writeWorkbook(zw *zip.Writer, wb *Workbook) error {
	f, err := zw.Create("xl/workbook.xml")
	if err != nil {
		return err
	}
	wbXML := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` + "\n"
	wbXML += `<workbook xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships">`
	wbXML += `<sheets>`
	for i, ws := range wb.worksheets {
		wbXML += fmt.Sprintf(`<sheet name="%s" sheetId="%d" r:id="rId%d"/>`,
			xmlEscape(ws.title), i+1, i+1)
	}
	wbXML += `</sheets>`
	wbXML += `</workbook>`
	_, err = f.Write([]byte(wbXML))
	return err
}

func (w *XLSXWriter) writeStyles(zw *zip.Writer, wb *Workbook) error {
	f, err := zw.Create("xl/styles.xml")
	if err != nil {
		return err
	}
	// TODO: Serialize user-defined styles (Font, Fill, Borders, Alignment, NumberFormat)
	// from cell Style objects. Currently only writes minimal defaults, so custom
	// styling set via SetStyle() is lost on save.
	styles := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` + "\n"
	styles += `<styleSheet xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main">`
	styles += `<fonts count="1"><font><sz val="11"/><name val="Calibri"/></font></fonts>`
	styles += `<fills count="2"><fill><patternFill patternType="none"/></fill><fill><patternFill patternType="gray125"/></fill></fills>`
	styles += `<borders count="1"><border><left/><right/><top/><bottom/><diagonal/></border></borders>`
	styles += `<cellStyleXfs count="1"><xf numFmtId="0" fontId="0" fillId="0" borderId="0"/></cellStyleXfs>`
	styles += `<cellXfs count="1"><xf numFmtId="0" fontId="0" fillId="0" borderId="0" xfId="0"/></cellXfs>`
	styles += `</styleSheet>`
	_, err = f.Write([]byte(styles))
	return err
}

func (w *XLSXWriter) writeSharedStrings(zw *zip.Writer, ss *sharedStrings) error {
	f, err := zw.Create("xl/sharedStrings.xml")
	if err != nil {
		return err
	}
	ssXML := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` + "\n"
	ssXML += fmt.Sprintf(`<sst xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main" count="%d" uniqueCount="%d">`,
		len(ss.strings), len(ss.strings))
	for _, s := range ss.strings {
		ssXML += fmt.Sprintf(`<si><t>%s</t></si>`, xmlEscape(s))
	}
	ssXML += `</sst>`
	_, err = f.Write([]byte(ssXML))
	return err
}

func (w *XLSXWriter) writeSheet(zw *zip.Writer, ws *Worksheet, sheetNum int, ss *sharedStrings) error {
	f, err := zw.Create(fmt.Sprintf("xl/worksheets/sheet%d.xml", sheetNum))
	if err != nil {
		return err
	}

	var b strings.Builder
	b.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` + "\n")
	b.WriteString(`<worksheet xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships">`)

	// Freeze pane
	if ws.frozen != nil {
		b.WriteString(`<sheetViews><sheetView tabSelected="1" workbookViewId="0">`)
		fmt.Fprintf(&b, `<pane xSplit="%d" ySplit="%d" topLeftCell="%s%d" activePane="bottomRight" state="frozen"/>`,
			ws.frozen.ColumnIdx, ws.frozen.Row-1, ws.frozen.Column, ws.frozen.Row)
		b.WriteString(`</sheetView></sheetViews>`)
	}

	// Column widths
	if len(ws.colWidths) > 0 {
		b.WriteString(`<cols>`)
		for col, width := range ws.colWidths {
			fmt.Fprintf(&b, `<col min="%d" max="%d" width="%.2f" customWidth="1"/>`, col+1, col+1, width)
		}
		b.WriteString(`</cols>`)
	}

	b.WriteString(`<sheetData>`)

	cells := ws.AllCells()
	if len(cells) > 0 {
		currentRow := -1
		for _, c := range cells {
			if c.row != currentRow {
				if currentRow >= 0 {
					b.WriteString(`</row>`)
				}
				currentRow = c.row
				if h, ok := ws.rowHeights[c.row]; ok {
					fmt.Fprintf(&b, `<row r="%d" ht="%.2f" customHeight="1">`, c.row+1, h)
				} else {
					fmt.Fprintf(&b, `<row r="%d">`, c.row+1)
				}
			}

			cellName, _ := CellName(c.row, c.col)
			switch c.Type {
			case CellTypeString:
				if s, ok := c.Value.(string); ok {
					idx := ss.GetIndex(s)
					fmt.Fprintf(&b, `<c r="%s" t="s"><v>%d</v></c>`, cellName, idx)
				}
			case CellTypeNumeric:
				if v, ok := c.Value.(float64); ok {
					fmt.Fprintf(&b, `<c r="%s"><v>%g</v></c>`, cellName, v)
				}
			case CellTypeBool:
				if v, ok := c.Value.(bool); ok {
					boolVal := 0
					if v {
						boolVal = 1
					}
					fmt.Fprintf(&b, `<c r="%s" t="b"><v>%d</v></c>`, cellName, boolVal)
				}
			case CellTypeFormula:
				fmt.Fprintf(&b, `<c r="%s"><f>%s</f></c>`, cellName, xmlEscape(c.Formula))
			case CellTypeDate:
				if t, ok := c.Value.(time.Time); ok {
					serial := dateToSerial(t)
					fmt.Fprintf(&b, `<c r="%s"><v>%g</v></c>`, cellName, serial)
				}
			}
		}
		if currentRow >= 0 {
			b.WriteString(`</row>`)
		}
	}

	b.WriteString(`</sheetData>`)

	// Merge cells
	if len(ws.mergeCells) > 0 {
		fmt.Fprintf(&b, `<mergeCells count="%d">`, len(ws.mergeCells))
		for _, mc := range ws.mergeCells {
			startName, _ := CellName(mc.StartRow, mc.StartCol)
			endName, _ := CellName(mc.EndRow, mc.EndCol)
			fmt.Fprintf(&b, `<mergeCell ref="%s:%s"/>`, startName, endName)
		}
		b.WriteString(`</mergeCells>`)
	}

	b.WriteString(`</worksheet>`)
	_, err = f.Write([]byte(b.String()))
	return err
}

func (w *XLSXWriter) writeCoreProperties(zw *zip.Writer, wb *Workbook) error {
	f, err := zw.Create("docProps/core.xml")
	if err != nil {
		return err
	}
	props := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` + "\n"
	props += `<cp:coreProperties xmlns:cp="http://schemas.openxmlformats.org/package/2006/metadata/core-properties" xmlns:dc="http://purl.org/dc/elements/1.1/" xmlns:dcterms="http://purl.org/dc/terms/" xmlns:dcmitype="http://purl.org/dc/dcmitype/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">`
	if wb.Properties.Creator != "" {
		props += fmt.Sprintf(`<dc:creator>%s</dc:creator>`, xmlEscape(wb.Properties.Creator))
	}
	if wb.Properties.Title != "" {
		props += fmt.Sprintf(`<dc:title>%s</dc:title>`, xmlEscape(wb.Properties.Title))
	}
	if wb.Properties.Subject != "" {
		props += fmt.Sprintf(`<dc:subject>%s</dc:subject>`, xmlEscape(wb.Properties.Subject))
	}
	if wb.Properties.Description != "" {
		props += fmt.Sprintf(`<dc:description>%s</dc:description>`, xmlEscape(wb.Properties.Description))
	}
	if wb.Properties.Keywords != "" {
		props += fmt.Sprintf(`<cp:keywords>%s</cp:keywords>`, xmlEscape(wb.Properties.Keywords))
	}
	if wb.Properties.Category != "" {
		props += fmt.Sprintf(`<cp:category>%s</cp:category>`, xmlEscape(wb.Properties.Category))
	}
	if wb.Properties.LastModifiedBy != "" {
		props += fmt.Sprintf(`<cp:lastModifiedBy>%s</cp:lastModifiedBy>`, xmlEscape(wb.Properties.LastModifiedBy))
	}
	props += `</cp:coreProperties>`
	_, err = f.Write([]byte(props))
	return err
}

// xmlEscape escapes special XML characters.
func xmlEscape(s string) string {
	var buf strings.Builder
	xml.EscapeText(&buf, []byte(s))
	return buf.String()
}

// dateToSerial converts a time.Time to an Excel serial date number.
// Excel's epoch is January 1, 1900 (serial number 1).
// Note: Excel incorrectly treats 1900 as a leap year (Lotus 1-2-3 bug).
func dateToSerial(t time.Time) float64 {
	epoch := time.Date(1899, 12, 30, 0, 0, 0, 0, time.UTC)
	duration := t.Sub(epoch)
	return duration.Hours() / 24.0
}

// serialToDate converts an Excel serial date number to time.Time.
func serialToDate(serial float64) time.Time {
	epoch := time.Date(1899, 12, 30, 0, 0, 0, 0, time.UTC)
	duration := time.Duration(serial * 24 * float64(time.Hour))
	return epoch.Add(duration)
}
