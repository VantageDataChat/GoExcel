package gospreadsheet

// HorizontalAlignment represents horizontal text alignment.
type HorizontalAlignment string

const (
	AlignLeft    HorizontalAlignment = "left"
	AlignCenter  HorizontalAlignment = "center"
	AlignRight   HorizontalAlignment = "right"
	AlignJustify HorizontalAlignment = "justify"
	AlignGeneral HorizontalAlignment = "general"
)

// VerticalAlignment represents vertical text alignment.
type VerticalAlignment string

const (
	AlignTop    VerticalAlignment = "top"
	AlignMiddle VerticalAlignment = "center"
	AlignBottom VerticalAlignment = "bottom"
)

// BorderStyle represents a border line style.
type BorderStyle string

const (
	BorderNone   BorderStyle = "none"
	BorderThin   BorderStyle = "thin"
	BorderMedium BorderStyle = "medium"
	BorderThick  BorderStyle = "thick"
	BorderDashed BorderStyle = "dashed"
	BorderDotted BorderStyle = "dotted"
	BorderDouble BorderStyle = "double"
)

// Border represents a single border edge.
type Border struct {
	Style BorderStyle
	Color string // hex color, e.g., "FF0000"
}

// Borders represents all four borders of a cell.
type Borders struct {
	Left   Border
	Right  Border
	Top    Border
	Bottom Border
}

// Font represents cell font properties.
type Font struct {
	Name      string
	Size      float64
	Bold      bool
	Italic    bool
	Underline bool
	Strikethrough bool
	Color     string // hex color
}

// Fill represents cell fill/background properties.
type Fill struct {
	Type    string // "solid", "pattern", "none"
	Color   string // hex color
	Pattern string // pattern type for pattern fills
}

// Alignment represents cell text alignment.
type Alignment struct {
	Horizontal HorizontalAlignment
	Vertical   VerticalAlignment
	WrapText   bool
	TextRotation int // degrees, -90 to 90
	Indent     int
}

// NumberFormat represents a cell number format.
type NumberFormat struct {
	FormatCode string
}

// Common number format codes.
var (
	FormatGeneral     = NumberFormat{FormatCode: "General"}
	FormatNumber      = NumberFormat{FormatCode: "0"}
	FormatNumber2Dec  = NumberFormat{FormatCode: "0.00"}
	FormatPercent     = NumberFormat{FormatCode: "0%"}
	FormatPercent2Dec = NumberFormat{FormatCode: "0.00%"}
	FormatDate        = NumberFormat{FormatCode: "yyyy-mm-dd"}
	FormatDateTime    = NumberFormat{FormatCode: "yyyy-mm-dd hh:mm:ss"}
	FormatTime        = NumberFormat{FormatCode: "hh:mm:ss"}
	FormatCurrency    = NumberFormat{FormatCode: `#,##0.00"$"`}
	FormatAccounting  = NumberFormat{FormatCode: `_("$"* #,##0.00_)`}
	FormatText        = NumberFormat{FormatCode: "@"}
)

// Style represents the complete style of a cell.
type Style struct {
	Font         *Font
	Fill         *Fill
	Borders      *Borders
	Alignment    *Alignment
	NumberFormat *NumberFormat
}

// NewStyle creates a new empty style.
func NewStyle() *Style {
	return &Style{}
}

// SetFont sets the font and returns the style for chaining.
func (s *Style) SetFont(f *Font) *Style {
	s.Font = f
	return s
}

// SetFill sets the fill and returns the style for chaining.
func (s *Style) SetFill(f *Fill) *Style {
	s.Fill = f
	return s
}

// SetBorders sets the borders and returns the style for chaining.
func (s *Style) SetBorders(b *Borders) *Style {
	s.Borders = b
	return s
}

// SetAlignment sets the alignment and returns the style for chaining.
func (s *Style) SetAlignment(a *Alignment) *Style {
	s.Alignment = a
	return s
}

// SetNumberFormat sets the number format and returns the style for chaining.
func (s *Style) SetNumberFormat(nf *NumberFormat) *Style {
	s.NumberFormat = nf
	return s
}

// DefaultFont returns the default font.
func DefaultFont() *Font {
	return &Font{
		Name: "Calibri",
		Size: 11,
	}
}
