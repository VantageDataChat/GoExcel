package gospreadsheet

// ValidationType represents the type of data validation.
type ValidationType string

const (
	ValidationNone       ValidationType = "none"
	ValidationWhole      ValidationType = "whole"
	ValidationDecimal    ValidationType = "decimal"
	ValidationList       ValidationType = "list"
	ValidationDate       ValidationType = "date"
	ValidationTime       ValidationType = "time"
	ValidationTextLength ValidationType = "textLength"
	ValidationCustom     ValidationType = "custom"
)

// ValidationErrorStyle represents the error alert style.
type ValidationErrorStyle string

const (
	ErrorStyleStop        ValidationErrorStyle = "stop"
	ErrorStyleWarning     ValidationErrorStyle = "warning"
	ErrorStyleInformation ValidationErrorStyle = "information"
)

// ValidationOperator represents the comparison operator for validation.
type ValidationOperator string

const (
	ValOperatorBetween           ValidationOperator = "between"
	ValOperatorNotBetween        ValidationOperator = "notBetween"
	ValOperatorEqual             ValidationOperator = "equal"
	ValOperatorNotEqual          ValidationOperator = "notEqual"
	ValOperatorGreaterThan       ValidationOperator = "greaterThan"
	ValOperatorLessThan          ValidationOperator = "lessThan"
	ValOperatorGreaterThanOrEqual ValidationOperator = "greaterThanOrEqual"
	ValOperatorLessThanOrEqual   ValidationOperator = "lessThanOrEqual"
)

// DataValidation represents a data validation rule for a cell range.
type DataValidation struct {
	Range         string // e.g., "A1:A100"
	Type          ValidationType
	Operator      ValidationOperator
	Formula1      string
	Formula2      string // used for between/notBetween
	AllowBlank    bool
	ShowInputMsg  bool
	ShowErrorMsg  bool
	ErrorStyle    ValidationErrorStyle
	ErrorTitle    string
	ErrorMessage  string
	PromptTitle   string
	PromptMessage string
}

// NewDataValidation creates a new data validation for a range.
func NewDataValidation(rangeStr string) *DataValidation {
	return &DataValidation{
		Range:        rangeStr,
		Type:         ValidationNone,
		AllowBlank:   true,
		ShowInputMsg: true,
		ShowErrorMsg: true,
		ErrorStyle:   ErrorStyleStop,
	}
}

// SetType sets the validation type and returns for chaining.
func (dv *DataValidation) SetType(t ValidationType) *DataValidation {
	dv.Type = t
	return dv
}

// SetOperator sets the operator and returns for chaining.
func (dv *DataValidation) SetOperator(op ValidationOperator) *DataValidation {
	dv.Operator = op
	return dv
}

// SetFormula1 sets the first formula/value.
func (dv *DataValidation) SetFormula1(f string) *DataValidation {
	dv.Formula1 = f
	return dv
}

// SetFormula2 sets the second formula/value (for between/notBetween).
func (dv *DataValidation) SetFormula2(f string) *DataValidation {
	dv.Formula2 = f
	return dv
}

// SetErrorMessage sets the error alert message.
func (dv *DataValidation) SetErrorMessage(title, message string) *DataValidation {
	dv.ErrorTitle = title
	dv.ErrorMessage = message
	dv.ShowErrorMsg = true
	return dv
}

// SetPromptMessage sets the input prompt message.
func (dv *DataValidation) SetPromptMessage(title, message string) *DataValidation {
	dv.PromptTitle = title
	dv.PromptMessage = message
	dv.ShowInputMsg = true
	return dv
}

// SetListValues is a convenience method for list validation with explicit values.
func (dv *DataValidation) SetListValues(values []string) *DataValidation {
	dv.Type = ValidationList
	joined := ""
	for i, v := range values {
		if i > 0 {
			joined += ","
		}
		joined += `"` + v + `"`
	}
	dv.Formula1 = joined
	return dv
}
