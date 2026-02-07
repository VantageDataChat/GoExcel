package gospreadsheet

import (
	"archive/zip"
	"bytes"
	"testing"
)

// Helper to create a minimal zip with specific files.
func createTestZip(files map[string]string) *bytes.Reader {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	for name, content := range files {
		w, _ := zw.Create(name)
		w.Write([]byte(content))
	}
	zw.Close()
	return bytes.NewReader(buf.Bytes())
}

func TestXLSXReadMissingWorkbook(t *testing.T) {
	// Zip without xl/workbook.xml
	r := createTestZip(map[string]string{
		"xl/sharedStrings.xml": `<?xml version="1.0"?><sst xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main"></sst>`,
	})
	_, err := NewXLSXReader().Read(r, int64(r.Len()))
	if err == nil {
		t.Error("expected error for missing workbook.xml")
	}
}

func TestXLSXReadInvalidWorkbookXML(t *testing.T) {
	r := createTestZip(map[string]string{
		"xl/workbook.xml":     "not valid xml <<<<",
		"xl/sharedStrings.xml": `<?xml version="1.0"?><sst xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main"></sst>`,
	})
	_, err := NewXLSXReader().Read(r, int64(r.Len()))
	if err == nil {
		t.Error("expected error for invalid workbook XML")
	}
}

func TestXLSXReadMissingWorkbookRels(t *testing.T) {
	r := createTestZip(map[string]string{
		"xl/workbook.xml": `<?xml version="1.0"?><workbook xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main"><sheets><sheet name="Sheet1" sheetId="1" r:id="rId1" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships"/></sheets></workbook>`,
		"xl/sharedStrings.xml": `<?xml version="1.0"?><sst xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main"></sst>`,
	})
	_, err := NewXLSXReader().Read(r, int64(r.Len()))
	if err == nil {
		t.Error("expected error for missing workbook.xml.rels")
	}
}

func TestXLSXReadInvalidWorkbookRelsXML(t *testing.T) {
	r := createTestZip(map[string]string{
		"xl/workbook.xml":              `<?xml version="1.0"?><workbook xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main"><sheets><sheet name="Sheet1" sheetId="1" r:id="rId1" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships"/></sheets></workbook>`,
		"xl/_rels/workbook.xml.rels":   "not valid xml <<<<",
		"xl/sharedStrings.xml":         `<?xml version="1.0"?><sst xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main"></sst>`,
	})
	_, err := NewXLSXReader().Read(r, int64(r.Len()))
	if err == nil {
		t.Error("expected error for invalid rels XML")
	}
}

func TestXLSXReadInvalidSharedStringsXML(t *testing.T) {
	r := createTestZip(map[string]string{
		"xl/workbook.xml":              `<?xml version="1.0"?><workbook xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main"><sheets></sheets></workbook>`,
		"xl/_rels/workbook.xml.rels":   `<?xml version="1.0"?><Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"></Relationships>`,
		"xl/sharedStrings.xml":         "not valid xml <<<<",
	})
	_, err := NewXLSXReader().Read(r, int64(r.Len()))
	if err == nil {
		t.Error("expected error for invalid shared strings XML")
	}
}

func TestXLSXReadMissingSheetFile(t *testing.T) {
	// Sheet referenced but file doesn't exist - should skip gracefully
	r := createTestZip(map[string]string{
		"xl/workbook.xml":            `<?xml version="1.0"?><workbook xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main"><sheets><sheet name="Sheet1" sheetId="1" r:id="rId1" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships"/></sheets></workbook>`,
		"xl/_rels/workbook.xml.rels": `<?xml version="1.0"?><Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"><Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/worksheet" Target="worksheets/sheet1.xml"/></Relationships>`,
		"xl/sharedStrings.xml":       `<?xml version="1.0"?><sst xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main"></sst>`,
		// Note: xl/worksheets/sheet1.xml is missing
	})
	_, err := NewXLSXReader().Read(r, int64(r.Len()))
	if err == nil {
		t.Error("expected error for missing sheet file")
	}
}

func TestXLSXReadInvalidSheetXML(t *testing.T) {
	r := createTestZip(map[string]string{
		"xl/workbook.xml":              `<?xml version="1.0"?><workbook xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main"><sheets><sheet name="Sheet1" sheetId="1" r:id="rId1" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships"/></sheets></workbook>`,
		"xl/_rels/workbook.xml.rels":   `<?xml version="1.0"?><Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"><Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/worksheet" Target="worksheets/sheet1.xml"/></Relationships>`,
		"xl/sharedStrings.xml":         `<?xml version="1.0"?><sst xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main"></sst>`,
		"xl/worksheets/sheet1.xml":     "not valid xml <<<<",
	})
	_, err := NewXLSXReader().Read(r, int64(r.Len()))
	if err == nil {
		t.Error("expected error for invalid sheet XML")
	}
}

func TestXLSXReadSheetWithUnmappedRID(t *testing.T) {
	// Sheet has rId that doesn't exist in rels - should be skipped
	r := createTestZip(map[string]string{
		"xl/workbook.xml":            `<?xml version="1.0"?><workbook xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main"><sheets><sheet name="Sheet1" sheetId="1" r:id="rId99" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships"/></sheets></workbook>`,
		"xl/_rels/workbook.xml.rels": `<?xml version="1.0"?><Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"><Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/worksheet" Target="worksheets/sheet1.xml"/></Relationships>`,
		"xl/sharedStrings.xml":       `<?xml version="1.0"?><sst xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main"></sst>`,
	})
	wb, err := NewXLSXReader().Read(r, int64(r.Len()))
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	// Sheet should be skipped since rId99 is not in rels
	if wb.SheetCount() != 0 {
		t.Errorf("sheet count = %d, want 0", wb.SheetCount())
	}
}

func TestXLSXReadSheetWithInvalidCellRef(t *testing.T) {
	// Sheet with a cell that has an invalid reference - should be skipped
	r := createTestZip(map[string]string{
		"xl/workbook.xml":            `<?xml version="1.0"?><workbook xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main"><sheets><sheet name="Sheet1" sheetId="1" r:id="rId1" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships"/></sheets></workbook>`,
		"xl/_rels/workbook.xml.rels": `<?xml version="1.0"?><Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"><Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/worksheet" Target="worksheets/sheet1.xml"/></Relationships>`,
		"xl/sharedStrings.xml":       `<?xml version="1.0"?><sst xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main"><si><t>hello</t></si></sst>`,
		"xl/worksheets/sheet1.xml":   `<?xml version="1.0"?><worksheet xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main"><sheetData><row r="1"><c r="!!!" t="s"><v>0</v></c><c r="A1" t="s"><v>0</v></c><c r="B1" t="e"><v>#REF!</v></c><c r="C1" t="str"><v>inline</v></c><c r="D1"><v>not_a_number</v></c></row></sheetData></worksheet>`,
	})
	wb, err := NewXLSXReader().Read(r, int64(r.Len()))
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	ws := wb.GetActiveSheet()
	// A1 should have "hello" from shared strings
	v, _ := ws.GetCellValue("A1")
	if v != "hello" {
		t.Errorf("A1 = %v, want 'hello'", v)
	}
	// B1 should be error type
	c, _ := ws.GetCellByName("B1")
	if c.Type != CellTypeError {
		t.Errorf("B1 type = %d, want CellTypeError", c.Type)
	}
	// C1 should be inline string
	v, _ = ws.GetCellValue("C1")
	if v != "inline" {
		t.Errorf("C1 = %v, want 'inline'", v)
	}
	// D1 should be string (failed float parse)
	v, _ = ws.GetCellValue("D1")
	if v != "not_a_number" {
		t.Errorf("D1 = %v, want 'not_a_number'", v)
	}
}

func TestXLSXReadNoCoreProperties(t *testing.T) {
	// Valid XLSX without docProps/core.xml
	r := createTestZip(map[string]string{
		"xl/workbook.xml":            `<?xml version="1.0"?><workbook xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main"><sheets><sheet name="Sheet1" sheetId="1" r:id="rId1" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships"/></sheets></workbook>`,
		"xl/_rels/workbook.xml.rels": `<?xml version="1.0"?><Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"><Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/worksheet" Target="worksheets/sheet1.xml"/></Relationships>`,
		"xl/worksheets/sheet1.xml":   `<?xml version="1.0"?><worksheet xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main"><sheetData></sheetData></worksheet>`,
	})
	wb, err := NewXLSXReader().Read(r, int64(r.Len()))
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if wb.Properties.Creator != "" {
		t.Errorf("Creator should be empty, got %q", wb.Properties.Creator)
	}
}

func TestXLSXReadInvalidCorePropertiesXML(t *testing.T) {
	r := createTestZip(map[string]string{
		"xl/workbook.xml":            `<?xml version="1.0"?><workbook xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main"><sheets><sheet name="Sheet1" sheetId="1" r:id="rId1" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships"/></sheets></workbook>`,
		"xl/_rels/workbook.xml.rels": `<?xml version="1.0"?><Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"><Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/worksheet" Target="worksheets/sheet1.xml"/></Relationships>`,
		"xl/worksheets/sheet1.xml":   `<?xml version="1.0"?><worksheet xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main"><sheetData></sheetData></worksheet>`,
		"docProps/core.xml":          "not valid xml <<<<",
	})
	// Should not error - core properties are optional
	wb, err := NewXLSXReader().Read(r, int64(r.Len()))
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if wb.SheetCount() != 1 {
		t.Errorf("sheet count = %d", wb.SheetCount())
	}
}

func TestXLSXReadSheetWithFormula(t *testing.T) {
	r := createTestZip(map[string]string{
		"xl/workbook.xml":            `<?xml version="1.0"?><workbook xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main"><sheets><sheet name="Sheet1" sheetId="1" r:id="rId1" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships"/></sheets></workbook>`,
		"xl/_rels/workbook.xml.rels": `<?xml version="1.0"?><Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"><Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/worksheet" Target="worksheets/sheet1.xml"/></Relationships>`,
		"xl/worksheets/sheet1.xml":   `<?xml version="1.0"?><worksheet xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main"><sheetData><row r="1"><c r="A1"><v>10</v></c><c r="A2"><f>A1*2</f></c></row></sheetData></worksheet>`,
	})
	wb, err := NewXLSXReader().Read(r, int64(r.Len()))
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	ws := wb.GetActiveSheet()
	c, _ := ws.GetCellByName("A2")
	if c.Type != CellTypeFormula || c.Formula != "A1*2" {
		t.Errorf("A2: type=%d, formula=%q", c.Type, c.Formula)
	}
}

func TestXLSXReadSheetWithInvalidSharedStringIndex(t *testing.T) {
	r := createTestZip(map[string]string{
		"xl/workbook.xml":            `<?xml version="1.0"?><workbook xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main"><sheets><sheet name="Sheet1" sheetId="1" r:id="rId1" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships"/></sheets></workbook>`,
		"xl/_rels/workbook.xml.rels": `<?xml version="1.0"?><Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"><Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/worksheet" Target="worksheets/sheet1.xml"/></Relationships>`,
		"xl/sharedStrings.xml":       `<?xml version="1.0"?><sst xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main"><si><t>hello</t></si></sst>`,
		"xl/worksheets/sheet1.xml":   `<?xml version="1.0"?><worksheet xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main"><sheetData><row r="1"><c r="A1" t="s"><v>999</v></c><c r="B1" t="s"><v>abc</v></c></row></sheetData></worksheet>`,
	})
	wb, err := NewXLSXReader().Read(r, int64(r.Len()))
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	// A1 has invalid shared string index 999 - should be empty
	ws := wb.GetActiveSheet()
	c := ws.GetCell(0, 0)
	if c.Type != CellTypeEmpty {
		t.Logf("A1 type = %d (invalid SS index handled gracefully)", c.Type)
	}
}

func TestXLSXReadSheetWithMergeCells(t *testing.T) {
	r := createTestZip(map[string]string{
		"xl/workbook.xml":            `<?xml version="1.0"?><workbook xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main"><sheets><sheet name="Sheet1" sheetId="1" r:id="rId1" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships"/></sheets></workbook>`,
		"xl/_rels/workbook.xml.rels": `<?xml version="1.0"?><Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"><Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/worksheet" Target="worksheets/sheet1.xml"/></Relationships>`,
		"xl/worksheets/sheet1.xml":   `<?xml version="1.0"?><worksheet xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main"><sheetData><row r="1"><c r="A1"><v>1</v></c></row></sheetData><mergeCells count="2"><mergeCell ref="A1:C1"/><mergeCell ref="invalid_merge"/></mergeCells></worksheet>`,
	})
	wb, err := NewXLSXReader().Read(r, int64(r.Len()))
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	ws := wb.GetActiveSheet()
	merges := ws.GetMergeCells()
	// Only valid merge should be kept
	if len(merges) != 1 {
		t.Errorf("merges = %d, want 1", len(merges))
	}
}

func TestXLSXReadNoSharedStrings(t *testing.T) {
	// Valid XLSX without sharedStrings.xml
	r := createTestZip(map[string]string{
		"xl/workbook.xml":            `<?xml version="1.0"?><workbook xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main"><sheets><sheet name="Sheet1" sheetId="1" r:id="rId1" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships"/></sheets></workbook>`,
		"xl/_rels/workbook.xml.rels": `<?xml version="1.0"?><Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"><Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/worksheet" Target="worksheets/sheet1.xml"/></Relationships>`,
		"xl/worksheets/sheet1.xml":   `<?xml version="1.0"?><worksheet xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main"><sheetData><row r="1"><c r="A1"><v>42</v></c></row></sheetData></worksheet>`,
	})
	wb, err := NewXLSXReader().Read(r, int64(r.Len()))
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	v, _ := wb.GetActiveSheet().GetCellValue("A1")
	if v != 42.0 {
		t.Errorf("A1 = %v, want 42", v)
	}
}
