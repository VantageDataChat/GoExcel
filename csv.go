package gospreadsheet

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
)

// CSVReader reads CSV files into a Workbook.
type CSVReader struct {
	Delimiter rune
	LazyQuotes bool
}

// NewCSVReader creates a new CSV reader with default settings.
func NewCSVReader() *CSVReader {
	return &CSVReader{
		Delimiter: ',',
	}
}

// Open reads a CSV file and returns a Workbook with a single worksheet.
func (r *CSVReader) Open(filename string) (*Workbook, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("opening file: %w", err)
	}
	defer f.Close()
	return r.Read(f)
}

// Read reads CSV data from an io.Reader.
func (r *CSVReader) Read(reader io.Reader) (*Workbook, error) {
	csvReader := csv.NewReader(reader)
	csvReader.Comma = r.Delimiter
	csvReader.LazyQuotes = r.LazyQuotes
	csvReader.FieldsPerRecord = -1 // allow variable field count

	wb := New()
	ws := wb.GetActiveSheet()

	row := 0
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("reading CSV row %d: %w", row+1, err)
		}
		for col, val := range record {
			cell := ws.GetCell(row, col)
			// Try to detect numeric values
			if v, err := strconv.ParseFloat(val, 64); err == nil {
				cell.SetValue(v)
			} else if v, err := strconv.ParseBool(val); err == nil {
				cell.SetValue(v)
			} else {
				cell.SetValue(val)
			}
		}
		row++
	}

	return wb, nil
}

// CSVWriter writes a Workbook to CSV format.
type CSVWriter struct {
	Delimiter rune
	SheetIndex int // which sheet to write (0-based)
}

// NewCSVWriter creates a new CSV writer with default settings.
func NewCSVWriter() *CSVWriter {
	return &CSVWriter{
		Delimiter: ',',
	}
}

// Save writes the workbook to a CSV file.
func (w *CSVWriter) Save(wb *Workbook, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("creating file: %w", err)
	}
	defer f.Close()
	return w.Write(wb, f)
}

// Write writes the workbook to an io.Writer.
func (w *CSVWriter) Write(wb *Workbook, writer io.Writer) error {
	ws, err := wb.GetSheet(w.SheetIndex)
	if err != nil {
		return err
	}

	csvWriter := csv.NewWriter(writer)
	csvWriter.Comma = w.Delimiter
	defer csvWriter.Flush()

	rows, err := ws.RowIterator()
	if err != nil {
		// Empty sheet, write nothing
		return nil
	}

	for _, row := range rows {
		record := make([]string, len(row))
		for i, cell := range row {
			if cell != nil {
				record[i] = cell.GetStringValue()
			}
		}
		if err := csvWriter.Write(record); err != nil {
			return fmt.Errorf("writing CSV row: %w", err)
		}
	}

	return csvWriter.Error()
}
