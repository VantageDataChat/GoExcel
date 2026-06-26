package gospreadsheet

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// XLSXReader reads XLSX files into a Workbook.
type XLSXReader struct{}

// NewXLSXReader creates a new XLSX reader.
func NewXLSXReader() *XLSXReader {
	return &XLSXReader{}
}

// Open reads an XLSX file and returns a Workbook.
func (r *XLSXReader) Open(filename string) (*Workbook, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("opening file: %w", err)
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		return nil, fmt.Errorf("stat file: %w", err)
	}

	return r.Read(f, info.Size())
}

// Read reads an XLSX from an io.ReaderAt with the given size.
func (r *XLSXReader) Read(reader io.ReaderAt, size int64) (*Workbook, error) {
	zr, err := zip.NewReader(reader, size)
	if err != nil {
		return nil, fmt.Errorf("opening zip: %w", err)
	}

	wb := NewEmpty()

	// Read shared strings
	ss, err := r.readSharedStrings(zr)
	if err != nil {
		return nil, fmt.Errorf("reading shared strings: %w", err)
	}

	// Read workbook to get sheet names and order
	sheetInfos, err := r.readWorkbookSheets(zr)
	if err != nil {
		return nil, fmt.Errorf("reading workbook: %w", err)
	}

	// Read workbook relationships to map rId -> sheet file
	relMap, err := r.readWorkbookRels(zr)
	if err != nil {
		return nil, fmt.Errorf("reading workbook rels: %w", err)
	}

	// Read styles to detect date-formatted cells
	styleNumFmts := r.readStyles(zr)

	// Read core properties
	r.readCoreProperties(zr, wb)

	// Read each sheet
	for _, si := range sheetInfos {
		target, ok := relMap[si.rID]
		if !ok {
			continue
		}
		ws, err := wb.AddSheet(si.name)
		if err != nil {
			return nil, err
		}
		if err := r.readSheet(zr, ws, target, ss, styleNumFmts); err != nil {
			return nil, fmt.Errorf("reading sheet %q: %w", si.name, err)
		}
	}

	return wb, nil
}

type sheetInfo struct {
	name string
	rID  string
}

func (r *XLSXReader) readSheet(zr *zip.Reader, ws *Worksheet, path string, ss []string, styleNumFmts []int) error {
	data, err := readZipFile(zr, path)
	if err != nil {
		return err
	}

	type xmlCellValue struct {
		Value string `xml:",chardata"`
	}
	type xmlCell struct {
		Ref     string       `xml:"r,attr"`
		Type    string       `xml:"t,attr"`
		Style   string       `xml:"s,attr"`
		Value   xmlCellValue `xml:"v"`
		Formula string       `xml:"f"`
	}
	type xmlRow struct {
		R     string    `xml:"r,attr"`
		Cells []xmlCell `xml:"c"`
	}
	type xmlMergeCell struct {
		Ref string `xml:"ref,attr"`
	}
	type xmlSheetData struct {
		Rows       []xmlRow       `xml:"sheetData>row"`
		MergeCells []xmlMergeCell `xml:"mergeCells>mergeCell"`
	}

	var sheetData xmlSheetData
	if err := xml.Unmarshal(data, &sheetData); err != nil {
		return fmt.Errorf("parsing sheet xml: %w", err)
	}

	for _, row := range sheetData.Rows {
		for _, c := range row.Cells {
			cr, err := ParseCellReference(c.Ref)
			if err != nil {
				continue
			}
			cell := ws.GetCell(cr.Row-1, cr.ColumnIdx)

			if c.Formula != "" {
				cell.SetFormula(c.Formula)
				continue
			}

			switch c.Type {
			case "s": // shared string
				idx, err := strconv.Atoi(c.Value.Value)
				if err == nil && idx >= 0 && idx < len(ss) {
					cell.SetValue(ss[idx])
				}
			case "b": // boolean
				cell.SetValue(c.Value.Value == "1")
			case "e": // error
				cell.Value = c.Value.Value
				cell.Type = CellTypeError
			case "str": // inline string
				cell.SetValue(c.Value.Value)
			default: // number or date
				if c.Value.Value != "" {
					if v, err := strconv.ParseFloat(c.Value.Value, 64); err == nil {
						// Check if the style indicates a date format
						styleIdx, _ := strconv.Atoi(c.Style)
						if styleIdx >= 0 && styleIdx < len(styleNumFmts) && isDateFormatID(styleNumFmts[styleIdx]) {
							cell.SetValue(serialToDate(v))
						} else {
							cell.SetValue(v)
						}
					} else {
						cell.SetValue(c.Value.Value)
					}
				}
			}
		}
	}

	// Read merge cells
	for _, mc := range sheetData.MergeCells {
		if err := ws.MergeCells(mc.Ref); err != nil {
			continue // skip invalid merge ranges
		}
	}

	return nil
}

func (r *XLSXReader) readWorkbookSheets(zr *zip.Reader) ([]sheetInfo, error) {
	data, err := readZipFile(zr, "xl/workbook.xml")
	if err != nil {
		return nil, err
	}

	type xmlSheet struct {
		Name    string `xml:"name,attr"`
		SheetID string `xml:"sheetId,attr"`
		RID     string `xml:"id,attr"`
	}
	type xmlSheets struct {
		Sheets []xmlSheet `xml:"sheets>sheet"`
	}

	var sheets xmlSheets
	if err := xml.Unmarshal(data, &sheets); err != nil {
		return nil, fmt.Errorf("parsing workbook.xml: %w", err)
	}

	infos := make([]sheetInfo, len(sheets.Sheets))
	for i, s := range sheets.Sheets {
		infos[i] = sheetInfo{name: s.Name, rID: s.RID}
	}
	return infos, nil
}

func (r *XLSXReader) readWorkbookRels(zr *zip.Reader) (map[string]string, error) {
	data, err := readZipFile(zr, "xl/_rels/workbook.xml.rels")
	if err != nil {
		return nil, err
	}

	type xmlRel struct {
		ID     string `xml:"Id,attr"`
		Target string `xml:"Target,attr"`
	}
	type xmlRels struct {
		Rels []xmlRel `xml:"Relationship"`
	}

	var rels xmlRels
	if err := xml.Unmarshal(data, &rels); err != nil {
		return nil, fmt.Errorf("parsing workbook.xml.rels: %w", err)
	}

	relMap := make(map[string]string)
	for _, rel := range rels.Rels {
		relMap[rel.ID] = workbookRelationshipTargetPath(rel.Target)
	}
	return relMap, nil
}

func workbookRelationshipTargetPath(target string) string {
	target = strings.ReplaceAll(strings.TrimSpace(target), "\\", "/")
	target = strings.TrimPrefix(target, "/")
	if strings.HasPrefix(target, "xl/") {
		return target
	}
	return "xl/" + target
}

func (r *XLSXReader) readSharedStrings(zr *zip.Reader) ([]string, error) {
	data, err := readZipFile(zr, "xl/sharedStrings.xml")
	if err != nil {
		// Shared strings may not exist
		return nil, nil
	}

	type xmlT struct {
		Value string `xml:",chardata"`
	}
	type xmlSI struct {
		T xmlT `xml:"t"`
	}
	type xmlSST struct {
		SI []xmlSI `xml:"si"`
	}

	var sst xmlSST
	if err := xml.Unmarshal(data, &sst); err != nil {
		return nil, fmt.Errorf("parsing sharedStrings.xml: %w", err)
	}

	strings := make([]string, len(sst.SI))
	for i, si := range sst.SI {
		strings[i] = si.T.Value
	}
	return strings, nil
}

// readStyles parses xl/styles.xml and returns a slice mapping each cellXf index
// to its numFmtId. This allows the reader to detect date-formatted cells.
func (r *XLSXReader) readStyles(zr *zip.Reader) []int {
	data, err := readZipFile(zr, "xl/styles.xml")
	if err != nil {
		return nil
	}

	type xmlNumFmt struct {
		NumFmtID   int    `xml:"numFmtId,attr"`
		FormatCode string `xml:"formatCode,attr"`
	}
	type xmlXf struct {
		NumFmtID int `xml:"numFmtId,attr"`
	}
	type xmlStyleSheet struct {
		NumFmts []xmlNumFmt `xml:"numFmts>numFmt"`
		CellXfs []xmlXf     `xml:"cellXfs>xf"`
	}

	var ss xmlStyleSheet
	if err := xml.Unmarshal(data, &ss); err != nil {
		return nil
	}

	// Build custom numFmt lookup for date detection
	customDateFmts := make(map[int]bool)
	for _, nf := range ss.NumFmts {
		if looksLikeDateFormat(nf.FormatCode) {
			customDateFmts[nf.NumFmtID] = true
		}
	}

	// Store the numFmtId for each xf, also mark custom date formats
	result := make([]int, len(ss.CellXfs))
	for i, xf := range ss.CellXfs {
		fmtID := xf.NumFmtID
		// If it's a custom format that looks like a date, remap to a known date ID
		// so isDateFormatID can catch it
		if customDateFmts[fmtID] {
			result[i] = 14 // treat as built-in date
		} else {
			result[i] = fmtID
		}
	}
	return result
}

// isDateFormatID returns true if the given numFmtId corresponds to a
// built-in Excel date/time format.
// Built-in date format IDs: 14-22, 27-36, 45-47, 50-58.
func isDateFormatID(id int) bool {
	return (id >= 14 && id <= 22) ||
		(id >= 27 && id <= 36) ||
		(id >= 45 && id <= 47) ||
		(id >= 50 && id <= 58)
}

// looksLikeDateFormat checks if a custom format code string contains
// date/time tokens like y, m, d, h, s. It strips quoted strings, escaped
// characters, and bracketed color/locale codes before checking, so that
// tokens inside literals (e.g. 0.00"mm") don't cause false positives.
func looksLikeDateFormat(code string) bool {
	if code == "" {
		return false
	}
	lower := strings.ToLower(code)
	// Skip formats that are clearly numeric/text
	if lower == "general" || lower == "@" {
		return false
	}

	// Strip content between double quotes (literal strings)
	cleaned := stripQuotedAndEscaped(lower)
	// Strip bracketed sections like [Red], [$-409], [h], etc.
	// but preserve [h] and [m] and [s] which are elapsed time tokens
	cleaned = stripBracketedExceptElapsed(cleaned)

	// Check for date/time tokens: y, m, d, h, s, am/pm, a/p
	// Single-char tokens are valid in Excel (e.g. "d-mmm", "m/d/yyyy", "h:mm")
	hasDate := strings.ContainsAny(cleaned, "yd")
	hasTime := strings.ContainsAny(cleaned, "hs")
	hasAMPM := strings.Contains(cleaned, "am/pm") || strings.Contains(cleaned, "a/p")

	if hasDate || hasAMPM {
		return true
	}
	// 'm' is ambiguous (months vs minutes). If 'h' or 's' is present,
	// it's a time format. If 'm' appears alone, we need 'y' or 'd' to
	// confirm it's a date — but those are already caught above.
	// So 'm' alone with 'h' or 's' = time format.
	if hasTime {
		return true
	}
	// 'm' alone (no y, d, h, s) could be months in a month-only format,
	// but that's extremely rare and ambiguous. We check for 'mmm' or longer
	// which is unambiguously a month name format.
	if strings.Contains(cleaned, "mmm") {
		return true
	}

	return false
}

// stripQuotedAndEscaped removes quoted strings ("...") and escaped characters
// (\x) from a format code so they don't interfere with token detection.
func stripQuotedAndEscaped(code string) string {
	var b strings.Builder
	b.Grow(len(code))
	i := 0
	for i < len(code) {
		switch code[i] {
		case '"':
			// Skip everything until closing quote
			i++
			for i < len(code) && code[i] != '"' {
				i++
			}
			if i < len(code) {
				i++ // skip closing quote
			}
		case '\\':
			// Skip escaped character
			i += 2
		default:
			b.WriteByte(code[i])
			i++
		}
	}
	return b.String()
}

// stripBracketedExceptElapsed removes bracketed sections like [Red], [$-409],
// [DBNum1] etc., but preserves elapsed time tokens [h], [m], [s], [hh], [mm], [ss].
func stripBracketedExceptElapsed(code string) string {
	var b strings.Builder
	b.Grow(len(code))
	i := 0
	for i < len(code) {
		if code[i] == '[' {
			end := strings.IndexByte(code[i:], ']')
			if end == -1 {
				b.WriteByte(code[i])
				i++
				continue
			}
			inner := code[i+1 : i+end]
			// Preserve elapsed time tokens
			switch strings.ToLower(inner) {
			case "h", "hh", "m", "mm", "s", "ss":
				b.WriteString(inner)
			}
			// All other bracketed content is dropped (colors, locale, etc.)
			i += end + 1
		} else {
			b.WriteByte(code[i])
			i++
		}
	}
	return b.String()
}

func (r *XLSXReader) readCoreProperties(zr *zip.Reader, wb *Workbook) {
	data, err := readZipFile(zr, "docProps/core.xml")
	if err != nil {
		return
	}

	type xmlCore struct {
		Creator        string `xml:"creator"`
		Title          string `xml:"title"`
		Subject        string `xml:"subject"`
		Description    string `xml:"description"`
		Keywords       string `xml:"keywords"`
		Category       string `xml:"category"`
		LastModifiedBy string `xml:"lastModifiedBy"`
	}

	var core xmlCore
	if err := xml.Unmarshal(data, &core); err != nil {
		return
	}

	wb.Properties.Creator = core.Creator
	wb.Properties.Title = core.Title
	wb.Properties.Subject = core.Subject
	wb.Properties.Description = core.Description
	wb.Properties.Keywords = core.Keywords
	wb.Properties.Category = core.Category
	wb.Properties.LastModifiedBy = core.LastModifiedBy
}

func readZipFile(zr *zip.Reader, name string) ([]byte, error) {
	for _, f := range zr.File {
		if strings.EqualFold(f.Name, name) {
			rc, err := f.Open()
			if err != nil {
				return nil, err
			}
			defer rc.Close()
			return io.ReadAll(rc)
		}
	}
	return nil, fmt.Errorf("file %q not found in archive", name)
}
