package gospreadsheet

import (
	"testing"
	"time"
)

func TestCellSetValue(t *testing.T) {
	c := NewCell(0, 0)

	// String
	c.SetValue("hello")
	if c.Type != CellTypeString {
		t.Errorf("expected CellTypeString, got %d", c.Type)
	}
	if c.Value != "hello" {
		t.Errorf("expected 'hello', got %v", c.Value)
	}

	// Int
	c.SetValue(42)
	if c.Type != CellTypeNumeric {
		t.Errorf("expected CellTypeNumeric, got %d", c.Type)
	}
	if c.Value != float64(42) {
		t.Errorf("expected 42.0, got %v", c.Value)
	}

	// Int32
	c.SetValue(int32(100))
	if c.Type != CellTypeNumeric {
		t.Errorf("expected CellTypeNumeric for int32, got %d", c.Type)
	}

	// Int64
	c.SetValue(int64(200))
	if c.Type != CellTypeNumeric {
		t.Errorf("expected CellTypeNumeric for int64, got %d", c.Type)
	}

	// Float32
	c.SetValue(float32(3.14))
	if c.Type != CellTypeNumeric {
		t.Errorf("expected CellTypeNumeric for float32, got %d", c.Type)
	}

	// Float64
	c.SetValue(3.14)
	if c.Type != CellTypeNumeric {
		t.Errorf("expected CellTypeNumeric, got %d", c.Type)
	}

	// Bool
	c.SetValue(true)
	if c.Type != CellTypeBool {
		t.Errorf("expected CellTypeBool, got %d", c.Type)
	}
	if c.Value != true {
		t.Errorf("expected true, got %v", c.Value)
	}

	// Time
	now := time.Now()
	c.SetValue(now)
	if c.Type != CellTypeDate {
		t.Errorf("expected CellTypeDate, got %d", c.Type)
	}

	// Nil
	c.SetValue(nil)
	if c.Type != CellTypeEmpty {
		t.Errorf("expected CellTypeEmpty, got %d", c.Type)
	}
}

func TestCellFormula(t *testing.T) {
	c := NewCell(0, 0)
	c.SetFormula("SUM(A1:A10)")
	if c.Type != CellTypeFormula {
		t.Errorf("expected CellTypeFormula, got %d", c.Type)
	}
	if c.Formula != "SUM(A1:A10)" {
		t.Errorf("expected 'SUM(A1:A10)', got %q", c.Formula)
	}
}

func TestCellGetStringValue(t *testing.T) {
	tests := []struct {
		name  string
		value interface{}
		want  string
	}{
		{"nil", nil, ""},
		{"string", "hello", "hello"},
		{"int", 42, "42"},
		{"float", 3.14, "3.14"},
		{"bool true", true, "TRUE"},
		{"bool false", false, "FALSE"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCell(0, 0)
			c.SetValue(tt.value)
			got := c.GetStringValue()
			if got != tt.want {
				t.Errorf("GetStringValue() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestCellGetNumericValue(t *testing.T) {
	c := NewCell(0, 0)
	c.SetValue(42.5)
	v, err := c.GetNumericValue()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v != 42.5 {
		t.Errorf("expected 42.5, got %f", v)
	}

	c.SetValue("not a number")
	_, err = c.GetNumericValue()
	if err == nil {
		t.Error("expected error for non-numeric cell")
	}
}

func TestCellGetBoolValue(t *testing.T) {
	c := NewCell(0, 0)
	c.SetValue(true)
	v, err := c.GetBoolValue()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !v {
		t.Error("expected true")
	}

	c.SetValue("not a bool")
	_, err = c.GetBoolValue()
	if err == nil {
		t.Error("expected error for non-bool cell")
	}
}

func TestCellGetDateValue(t *testing.T) {
	c := NewCell(0, 0)
	now := time.Date(2025, 1, 15, 10, 30, 0, 0, time.UTC)
	c.SetValue(now)
	v, err := c.GetDateValue()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !v.Equal(now) {
		t.Errorf("expected %v, got %v", now, v)
	}

	c.SetValue("not a date")
	_, err = c.GetDateValue()
	if err == nil {
		t.Error("expected error for non-date cell")
	}
}

func TestCellRowCol(t *testing.T) {
	c := NewCell(5, 3)
	if c.Row() != 5 {
		t.Errorf("Row() = %d, want 5", c.Row())
	}
	if c.Col() != 3 {
		t.Errorf("Col() = %d, want 3", c.Col())
	}
}

func TestCellStyle(t *testing.T) {
	c := NewCell(0, 0)
	s := NewStyle().SetFont(&Font{Bold: true, Size: 14})
	c.SetStyle(s)
	if c.Style == nil {
		t.Fatal("style should not be nil")
	}
	if !c.Style.Font.Bold {
		t.Error("font should be bold")
	}
	if c.Style.Font.Size != 14 {
		t.Errorf("font size = %f, want 14", c.Style.Font.Size)
	}
}
