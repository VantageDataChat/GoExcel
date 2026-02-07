package gospreadsheet

// FilterType represents the type of auto filter.
type FilterType string

const (
	FilterCustom  FilterType = "custom"
	FilterDynamic FilterType = "dynamic"
	FilterTop10   FilterType = "top10"
	FilterValues  FilterType = "values"
)

// FilterOperator represents a filter comparison operator.
type FilterOperator string

const (
	FilterOpEqual          FilterOperator = "equal"
	FilterOpNotEqual       FilterOperator = "notEqual"
	FilterOpGreaterThan    FilterOperator = "greaterThan"
	FilterOpGreaterOrEqual FilterOperator = "greaterThanOrEqual"
	FilterOpLessThan       FilterOperator = "lessThan"
	FilterOpLessOrEqual    FilterOperator = "lessThanOrEqual"
)

// FilterCondition represents a single filter condition.
type FilterCondition struct {
	Operator FilterOperator
	Value    string
}

// AutoFilterColumn represents filter settings for a single column.
type AutoFilterColumn struct {
	ColumnIndex int
	FilterType  FilterType
	Conditions  []FilterCondition
	Values      []string // for value-based filtering
	ShowButton  bool
}

// AutoFilter represents auto-filter settings for a worksheet.
type AutoFilter struct {
	Range   string // e.g., "A1:D100"
	Columns []AutoFilterColumn
}

// NewAutoFilter creates a new auto filter for the given range.
func NewAutoFilter(rangeStr string) *AutoFilter {
	return &AutoFilter{
		Range:   rangeStr,
		Columns: make([]AutoFilterColumn, 0),
	}
}

// AddColumn adds a filter column configuration.
func (af *AutoFilter) AddColumn(col AutoFilterColumn) *AutoFilter {
	col.ShowButton = true
	af.Columns = append(af.Columns, col)
	return af
}

// AddValueFilter adds a value-based filter for a column.
func (af *AutoFilter) AddValueFilter(colIndex int, values []string) *AutoFilter {
	af.Columns = append(af.Columns, AutoFilterColumn{
		ColumnIndex: colIndex,
		FilterType:  FilterValues,
		Values:      values,
		ShowButton:  true,
	})
	return af
}

// AddCustomFilter adds a custom filter with one or two conditions.
func (af *AutoFilter) AddCustomFilter(colIndex int, conditions ...FilterCondition) *AutoFilter {
	af.Columns = append(af.Columns, AutoFilterColumn{
		ColumnIndex: colIndex,
		FilterType:  FilterCustom,
		Conditions:  conditions,
		ShowButton:  true,
	})
	return af
}
