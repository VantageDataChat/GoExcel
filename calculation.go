package gospreadsheet

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// CalculationEngine evaluates cell formulas.
type CalculationEngine struct {
	workbook *Workbook
}

// NewCalculationEngine creates a new calculation engine for a workbook.
func NewCalculationEngine(wb *Workbook) *CalculationEngine {
	return &CalculationEngine{workbook: wb}
}

// CalculateCell evaluates the formula in a cell and returns the result.
func (ce *CalculationEngine) CalculateCell(ws *Worksheet, ref string) (interface{}, error) {
	cell, err := ws.GetCellByName(ref)
	if err != nil {
		return nil, err
	}
	if cell.Type != CellTypeFormula {
		return cell.Value, nil
	}
	return ce.evaluate(ws, cell.Formula)
}

// CalculateAll evaluates all formula cells in the workbook.
func (ce *CalculationEngine) CalculateAll() error {
	for i := 0; i < ce.workbook.SheetCount(); i++ {
		ws, _ := ce.workbook.GetSheet(i)
		for _, cell := range ws.AllCells() {
			if cell.Type == CellTypeFormula {
				val, err := ce.evaluate(ws, cell.Formula)
				if err != nil {
					cell.Value = "#ERROR!"
				} else {
					cell.Value = val
				}
			}
		}
	}
	return nil
}

// evaluate parses and evaluates a formula string.
func (ce *CalculationEngine) evaluate(ws *Worksheet, formula string) (interface{}, error) {
	formula = strings.TrimSpace(formula)
	if formula == "" {
		return nil, fmt.Errorf("empty formula")
	}

	// Check for function calls
	upperFormula := strings.ToUpper(formula)
	for fname, fn := range getBuiltinFunctions() {
		if strings.HasPrefix(upperFormula, fname+"(") && strings.HasSuffix(formula, ")") {
			args := formula[len(fname)+1 : len(formula)-1]
			return fn(ce, ws, args)
		}
	}

	// Check for cell reference
	if isCellRef(formula) {
		return ce.resolveCellValue(ws, formula)
	}

	// Try numeric literal
	if v, err := strconv.ParseFloat(formula, 64); err == nil {
		return v, nil
	}

	// Try string literal
	if strings.HasPrefix(formula, `"`) && strings.HasSuffix(formula, `"`) {
		return formula[1 : len(formula)-1], nil
	}

	// Simple binary operations: A1+B1, A1-B1, A1*B1, A1/B1
	for _, op := range []string{"+", "-", "*", "/"} {
		if idx := findOperator(formula, op); idx > 0 {
			left, err := ce.evaluate(ws, formula[:idx])
			if err != nil {
				return nil, err
			}
			right, err := ce.evaluate(ws, formula[idx+1:])
			if err != nil {
				return nil, err
			}
			return applyOp(toFloat(left), toFloat(right), op)
		}
	}

	return nil, fmt.Errorf("cannot evaluate: %s", formula)
}

func findOperator(formula string, op string) int {
	depth := 0
	// Search from right to left for + and - (lower precedence)
	// Search from left to right for * and /
	if op == "+" || op == "-" {
		for i := len(formula) - 1; i > 0; i-- {
			ch := formula[i]
			if ch == ')' {
				depth++
			} else if ch == '(' {
				depth--
			} else if depth == 0 && string(ch) == op {
				return i
			}
		}
	} else {
		for i := len(formula) - 1; i > 0; i-- {
			ch := formula[i]
			if ch == ')' {
				depth++
			} else if ch == '(' {
				depth--
			} else if depth == 0 && string(ch) == op {
				return i
			}
		}
	}
	return -1
}

func applyOp(left, right float64, op string) (float64, error) {
	switch op {
	case "+":
		return left + right, nil
	case "-":
		return left - right, nil
	case "*":
		return left * right, nil
	case "/":
		if right == 0 {
			return 0, fmt.Errorf("#DIV/0!")
		}
		return left / right, nil
	}
	return 0, fmt.Errorf("unknown operator: %s", op)
}

func toFloat(v interface{}) float64 {
	switch val := v.(type) {
	case float64:
		return val
	case int:
		return float64(val)
	case bool:
		if val {
			return 1
		}
		return 0
	case string:
		if f, err := strconv.ParseFloat(val, 64); err == nil {
			return f
		}
	}
	return 0
}

func isCellRef(s string) bool {
	s = strings.TrimSpace(s)
	if s == "" {
		return false
	}
	hasLetter := false
	hasDigit := false
	for i, ch := range s {
		if unicode.IsLetter(ch) {
			if hasDigit {
				return false
			}
			hasLetter = true
		} else if unicode.IsDigit(ch) {
			if i == 0 {
				return false
			}
			hasDigit = true
		} else {
			return false
		}
	}
	return hasLetter && hasDigit
}

func (ce *CalculationEngine) resolveCellValue(ws *Worksheet, ref string) (interface{}, error) {
	cell, err := ws.GetCellByName(ref)
	if err != nil {
		return nil, err
	}
	if cell.Type == CellTypeFormula {
		return ce.evaluate(ws, cell.Formula)
	}
	return cell.Value, nil
}

// resolveRange resolves a range like "A1:A10" to a slice of values.
func (ce *CalculationEngine) resolveRange(ws *Worksheet, rangeStr string) ([]interface{}, error) {
	start, end, err := ParseRange(rangeStr)
	if err != nil {
		return nil, err
	}
	var values []interface{}
	for row := start.Row - 1; row <= end.Row-1; row++ {
		for col := start.ColumnIdx; col <= end.ColumnIdx; col++ {
			cell := ws.GetCell(row, col)
			var val interface{}
			if cell.Type == CellTypeFormula {
				val, _ = ce.evaluate(ws, cell.Formula)
			} else {
				val = cell.Value
			}
			if val != nil {
				values = append(values, val)
			}
		}
	}
	return values, nil
}

// parseArgs splits function arguments, respecting nested parentheses.
func parseArgs(args string) []string {
	var result []string
	depth := 0
	current := ""
	for _, ch := range args {
		if ch == '(' {
			depth++
			current += string(ch)
		} else if ch == ')' {
			depth--
			current += string(ch)
		} else if ch == ',' && depth == 0 {
			result = append(result, strings.TrimSpace(current))
			current = ""
		} else {
			current += string(ch)
		}
	}
	if current != "" {
		result = append(result, strings.TrimSpace(current))
	}
	return result
}
