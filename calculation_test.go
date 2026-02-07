package gospreadsheet

import (
	"math"
	"testing"
)

func setupCalcSheet() (*Workbook, *Worksheet, *CalculationEngine) {
	wb := New()
	ws := wb.GetActiveSheet()
	ce := NewCalculationEngine(wb)
	return wb, ws, ce
}

func TestCalcSUM(t *testing.T) {
	_, ws, ce := setupCalcSheet()
	ws.SetCellValue("A1", 10)
	ws.SetCellValue("A2", 20)
	ws.SetCellValue("A3", 30)
	ws.SetCellFormula("A4", "SUM(A1:A3)")

	val, err := ce.CalculateCell(ws, "A4")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if val != 60.0 {
		t.Errorf("SUM = %v, want 60", val)
	}
}

func TestCalcSUMMultipleArgs(t *testing.T) {
	_, ws, ce := setupCalcSheet()
	ws.SetCellValue("A1", 10)
	ws.SetCellValue("B1", 20)
	ws.SetCellFormula("C1", "SUM(A1,B1)")

	val, err := ce.CalculateCell(ws, "C1")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if val != 30.0 {
		t.Errorf("SUM = %v, want 30", val)
	}
}

func TestCalcAVERAGE(t *testing.T) {
	_, ws, ce := setupCalcSheet()
	ws.SetCellValue("A1", 10)
	ws.SetCellValue("A2", 20)
	ws.SetCellValue("A3", 30)
	ws.SetCellFormula("A4", "AVERAGE(A1:A3)")

	val, err := ce.CalculateCell(ws, "A4")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if val != 20.0 {
		t.Errorf("AVERAGE = %v, want 20", val)
	}
}

func TestCalcCOUNT(t *testing.T) {
	_, ws, ce := setupCalcSheet()
	ws.SetCellValue("A1", 10)
	ws.SetCellValue("A2", 20)
	ws.SetCellValue("A3", 30)
	ws.SetCellFormula("A4", "COUNT(A1:A3)")

	val, err := ce.CalculateCell(ws, "A4")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if val != 3.0 {
		t.Errorf("COUNT = %v, want 3", val)
	}
}

func TestCalcMAX(t *testing.T) {
	_, ws, ce := setupCalcSheet()
	ws.SetCellValue("A1", 10)
	ws.SetCellValue("A2", 50)
	ws.SetCellValue("A3", 30)
	ws.SetCellFormula("A4", "MAX(A1:A3)")

	val, err := ce.CalculateCell(ws, "A4")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if val != 50.0 {
		t.Errorf("MAX = %v, want 50", val)
	}
}

func TestCalcMIN(t *testing.T) {
	_, ws, ce := setupCalcSheet()
	ws.SetCellValue("A1", 10)
	ws.SetCellValue("A2", 50)
	ws.SetCellValue("A3", 30)
	ws.SetCellFormula("A4", "MIN(A1:A3)")

	val, err := ce.CalculateCell(ws, "A4")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if val != 10.0 {
		t.Errorf("MIN = %v, want 10", val)
	}
}

func TestCalcIF(t *testing.T) {
	_, ws, ce := setupCalcSheet()
	ws.SetCellValue("A1", 10)
	ws.SetCellFormula("B1", "IF(A1,\"yes\",\"no\")")

	val, err := ce.CalculateCell(ws, "B1")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if val != "yes" {
		t.Errorf("IF = %v, want 'yes'", val)
	}
}

func TestCalcIFFalse(t *testing.T) {
	_, ws, ce := setupCalcSheet()
	ws.SetCellValue("A1", 0)
	ws.SetCellFormula("B1", "IF(A1,\"yes\",\"no\")")

	val, err := ce.CalculateCell(ws, "B1")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if val != "no" {
		t.Errorf("IF = %v, want 'no'", val)
	}
}

func TestCalcABS(t *testing.T) {
	_, ws, ce := setupCalcSheet()
	ws.SetCellFormula("A1", "ABS(-42)")

	val, err := ce.CalculateCell(ws, "A1")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if val != 42.0 {
		t.Errorf("ABS = %v, want 42", val)
	}
}

func TestCalcROUND(t *testing.T) {
	_, ws, ce := setupCalcSheet()
	ws.SetCellFormula("A1", "ROUND(3.14159,2)")

	val, err := ce.CalculateCell(ws, "A1")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if val != 3.14 {
		t.Errorf("ROUND = %v, want 3.14", val)
	}
}

func TestCalcSQRT(t *testing.T) {
	_, ws, ce := setupCalcSheet()
	ws.SetCellFormula("A1", "SQRT(144)")

	val, err := ce.CalculateCell(ws, "A1")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if val != 12.0 {
		t.Errorf("SQRT = %v, want 12", val)
	}
}

func TestCalcSQRTNegative(t *testing.T) {
	_, ws, ce := setupCalcSheet()
	ws.SetCellFormula("A1", "SQRT(-1)")

	_, err := ce.CalculateCell(ws, "A1")
	if err == nil {
		t.Error("expected error for SQRT(-1)")
	}
}

func TestCalcPOWER(t *testing.T) {
	_, ws, ce := setupCalcSheet()
	ws.SetCellFormula("A1", "POWER(2,10)")

	val, err := ce.CalculateCell(ws, "A1")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if val != 1024.0 {
		t.Errorf("POWER = %v, want 1024", val)
	}
}

func TestCalcMOD(t *testing.T) {
	_, ws, ce := setupCalcSheet()
	ws.SetCellFormula("A1", "MOD(10,3)")

	val, err := ce.CalculateCell(ws, "A1")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if val != 1.0 {
		t.Errorf("MOD = %v, want 1", val)
	}
}

func TestCalcMODDivZero(t *testing.T) {
	_, ws, ce := setupCalcSheet()
	ws.SetCellFormula("A1", "MOD(10,0)")

	_, err := ce.CalculateCell(ws, "A1")
	if err == nil {
		t.Error("expected error for MOD(10,0)")
	}
}

func TestCalcINT(t *testing.T) {
	_, ws, ce := setupCalcSheet()
	ws.SetCellFormula("A1", "INT(7.8)")

	val, err := ce.CalculateCell(ws, "A1")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if val != 7.0 {
		t.Errorf("INT = %v, want 7", val)
	}
}

func TestCalcLEN(t *testing.T) {
	_, ws, ce := setupCalcSheet()
	ws.SetCellValue("A1", "hello")
	ws.SetCellFormula("B1", "LEN(A1)")

	val, err := ce.CalculateCell(ws, "B1")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if val != 5.0 {
		t.Errorf("LEN = %v, want 5", val)
	}
}

func TestCalcUPPER(t *testing.T) {
	_, ws, ce := setupCalcSheet()
	ws.SetCellValue("A1", "hello")
	ws.SetCellFormula("B1", "UPPER(A1)")

	val, err := ce.CalculateCell(ws, "B1")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if val != "HELLO" {
		t.Errorf("UPPER = %v, want 'HELLO'", val)
	}
}

func TestCalcLOWER(t *testing.T) {
	_, ws, ce := setupCalcSheet()
	ws.SetCellValue("A1", "HELLO")
	ws.SetCellFormula("B1", "LOWER(A1)")

	val, err := ce.CalculateCell(ws, "B1")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if val != "hello" {
		t.Errorf("LOWER = %v, want 'hello'", val)
	}
}

func TestCalcTRIM(t *testing.T) {
	_, ws, ce := setupCalcSheet()
	ws.SetCellValue("A1", "  hello  ")
	ws.SetCellFormula("B1", "TRIM(A1)")

	val, err := ce.CalculateCell(ws, "B1")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if val != "hello" {
		t.Errorf("TRIM = %v, want 'hello'", val)
	}
}

func TestCalcLEFT(t *testing.T) {
	_, ws, ce := setupCalcSheet()
	ws.SetCellValue("A1", "Hello World")
	ws.SetCellFormula("B1", "LEFT(A1,5)")

	val, err := ce.CalculateCell(ws, "B1")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if val != "Hello" {
		t.Errorf("LEFT = %v, want 'Hello'", val)
	}
}

func TestCalcRIGHT(t *testing.T) {
	_, ws, ce := setupCalcSheet()
	ws.SetCellValue("A1", "Hello World")
	ws.SetCellFormula("B1", "RIGHT(A1,5)")

	val, err := ce.CalculateCell(ws, "B1")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if val != "World" {
		t.Errorf("RIGHT = %v, want 'World'", val)
	}
}

func TestCalcMID(t *testing.T) {
	_, ws, ce := setupCalcSheet()
	ws.SetCellValue("A1", "Hello World")
	ws.SetCellFormula("B1", "MID(A1,7,5)")

	val, err := ce.CalculateCell(ws, "B1")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if val != "World" {
		t.Errorf("MID = %v, want 'World'", val)
	}
}

func TestCalcCONCATENATE(t *testing.T) {
	_, ws, ce := setupCalcSheet()
	ws.SetCellValue("A1", "Hello")
	ws.SetCellValue("B1", " ")
	ws.SetCellValue("C1", "World")
	ws.SetCellFormula("D1", "CONCATENATE(A1,B1,C1)")

	val, err := ce.CalculateCell(ws, "D1")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if val != "Hello World" {
		t.Errorf("CONCATENATE = %v, want 'Hello World'", val)
	}
}

func TestCalcMEDIAN(t *testing.T) {
	_, ws, ce := setupCalcSheet()
	ws.SetCellValue("A1", 1)
	ws.SetCellValue("A2", 3)
	ws.SetCellValue("A3", 5)
	ws.SetCellFormula("A4", "MEDIAN(A1:A3)")

	val, err := ce.CalculateCell(ws, "A4")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if val != 3.0 {
		t.Errorf("MEDIAN = %v, want 3", val)
	}
}

func TestCalcMEDIANEven(t *testing.T) {
	_, ws, ce := setupCalcSheet()
	ws.SetCellValue("A1", 1)
	ws.SetCellValue("A2", 2)
	ws.SetCellValue("A3", 3)
	ws.SetCellValue("A4", 4)
	ws.SetCellFormula("A5", "MEDIAN(A1:A4)")

	val, err := ce.CalculateCell(ws, "A5")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if val != 2.5 {
		t.Errorf("MEDIAN = %v, want 2.5", val)
	}
}

func TestCalcSUMIF(t *testing.T) {
	_, ws, ce := setupCalcSheet()
	ws.SetCellValue("A1", "apple")
	ws.SetCellValue("A2", "banana")
	ws.SetCellValue("A3", "apple")
	ws.SetCellValue("B1", 10)
	ws.SetCellValue("B2", 20)
	ws.SetCellValue("B3", 30)
	ws.SetCellFormula("B4", `SUMIF(A1:A3,"apple",B1:B3)`)

	val, err := ce.CalculateCell(ws, "B4")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if val != 40.0 {
		t.Errorf("SUMIF = %v, want 40", val)
	}
}

func TestCalcSUMIFNumeric(t *testing.T) {
	_, ws, ce := setupCalcSheet()
	ws.SetCellValue("A1", 10)
	ws.SetCellValue("A2", 20)
	ws.SetCellValue("A3", 30)
	ws.SetCellFormula("A4", `SUMIF(A1:A3,">15")`)

	val, err := ce.CalculateCell(ws, "A4")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if val != 50.0 {
		t.Errorf("SUMIF = %v, want 50", val)
	}
}

func TestCalcCOUNTIF(t *testing.T) {
	_, ws, ce := setupCalcSheet()
	ws.SetCellValue("A1", "apple")
	ws.SetCellValue("A2", "banana")
	ws.SetCellValue("A3", "apple")
	ws.SetCellValue("A4", "cherry")
	ws.SetCellFormula("A5", `COUNTIF(A1:A4,"apple")`)

	val, err := ce.CalculateCell(ws, "A5")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if val != 2.0 {
		t.Errorf("COUNTIF = %v, want 2", val)
	}
}

func TestCalcArithmetic(t *testing.T) {
	_, ws, ce := setupCalcSheet()
	ws.SetCellValue("A1", 10)
	ws.SetCellValue("B1", 3)

	tests := []struct {
		formula string
		want    float64
	}{
		{"A1+B1", 13},
		{"A1-B1", 7},
		{"A1*B1", 30},
		{"A1/B1", 10.0 / 3.0},
	}

	for _, tt := range tests {
		ws.SetCellFormula("C1", tt.formula)
		val, err := ce.CalculateCell(ws, "C1")
		if err != nil {
			t.Errorf("%s: error: %v", tt.formula, err)
			continue
		}
		if math.Abs(toFloat(val)-tt.want) > 0.0001 {
			t.Errorf("%s = %v, want %v", tt.formula, val, tt.want)
		}
	}
}

func TestCalcDivByZero(t *testing.T) {
	_, ws, ce := setupCalcSheet()
	ws.SetCellValue("A1", 10)
	ws.SetCellValue("B1", 0)
	ws.SetCellFormula("C1", "A1/B1")

	_, err := ce.CalculateCell(ws, "C1")
	if err == nil {
		t.Error("expected error for division by zero")
	}
}

func TestCalcNestedFormula(t *testing.T) {
	_, ws, ce := setupCalcSheet()
	ws.SetCellValue("A1", 10)
	ws.SetCellValue("A2", 20)
	ws.SetCellFormula("A3", "SUM(A1:A2)")
	ws.SetCellFormula("A4", "A3+5")

	val, err := ce.CalculateCell(ws, "A4")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if val != 35.0 {
		t.Errorf("nested = %v, want 35", val)
	}
}

func TestCalcCalculateAll(t *testing.T) {
	wb, ws, ce := setupCalcSheet()
	_ = wb
	ws.SetCellValue("A1", 5)
	ws.SetCellValue("A2", 10)
	ws.SetCellFormula("A3", "SUM(A1:A2)")
	ws.SetCellFormula("A4", "A3*2")

	ce.CalculateAll()

	c, _ := ws.GetCellByName("A3")
	if c.Value != 15.0 {
		t.Errorf("A3 = %v, want 15", c.Value)
	}

	c, _ = ws.GetCellByName("A4")
	if c.Value != 30.0 {
		t.Errorf("A4 = %v, want 30", c.Value)
	}
}

func TestCalcNonFormulaCell(t *testing.T) {
	_, ws, ce := setupCalcSheet()
	ws.SetCellValue("A1", 42)

	val, err := ce.CalculateCell(ws, "A1")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if val != float64(42) {
		t.Errorf("value = %v, want 42", val)
	}
}

func TestCalcCOUNTA(t *testing.T) {
	_, ws, ce := setupCalcSheet()
	ws.SetCellValue("A1", "hello")
	ws.SetCellValue("A2", 42)
	ws.SetCellValue("A3", true)
	ws.SetCellFormula("A4", "COUNTA(A1:A3)")

	val, err := ce.CalculateCell(ws, "A4")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if val != 3.0 {
		t.Errorf("COUNTA = %v, want 3", val)
	}
}
