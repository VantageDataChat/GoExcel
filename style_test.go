package gospreadsheet

import (
	"testing"
)

func TestNewStyle(t *testing.T) {
	s := NewStyle()
	if s == nil {
		t.Fatal("NewStyle returned nil")
	}
	if s.Font != nil {
		t.Error("Font should be nil initially")
	}
	if s.Fill != nil {
		t.Error("Fill should be nil initially")
	}
	if s.Borders != nil {
		t.Error("Borders should be nil initially")
	}
	if s.Alignment != nil {
		t.Error("Alignment should be nil initially")
	}
	if s.NumberFormat != nil {
		t.Error("NumberFormat should be nil initially")
	}
}

func TestStyleChaining(t *testing.T) {
	s := NewStyle().
		SetFont(&Font{
			Name:      "Arial",
			Size:      12,
			Bold:      true,
			Italic:    true,
			Underline: true,
			Color:     "FF0000",
		}).
		SetFill(&Fill{
			Type:  "solid",
			Color: "FFFF00",
		}).
		SetBorders(&Borders{
			Left:   Border{Style: BorderThin, Color: "000000"},
			Right:  Border{Style: BorderThin, Color: "000000"},
			Top:    Border{Style: BorderMedium, Color: "000000"},
			Bottom: Border{Style: BorderMedium, Color: "000000"},
		}).
		SetAlignment(&Alignment{
			Horizontal:   AlignCenter,
			Vertical:     AlignMiddle,
			WrapText:     true,
			TextRotation: 45,
		}).
		SetNumberFormat(&FormatNumber2Dec)

	if s.Font.Name != "Arial" {
		t.Errorf("Font.Name = %q, want 'Arial'", s.Font.Name)
	}
	if s.Font.Size != 12 {
		t.Errorf("Font.Size = %f, want 12", s.Font.Size)
	}
	if !s.Font.Bold {
		t.Error("Font.Bold should be true")
	}
	if !s.Font.Italic {
		t.Error("Font.Italic should be true")
	}
	if !s.Font.Underline {
		t.Error("Font.Underline should be true")
	}
	if s.Font.Color != "FF0000" {
		t.Errorf("Font.Color = %q, want 'FF0000'", s.Font.Color)
	}

	if s.Fill.Type != "solid" {
		t.Errorf("Fill.Type = %q, want 'solid'", s.Fill.Type)
	}
	if s.Fill.Color != "FFFF00" {
		t.Errorf("Fill.Color = %q, want 'FFFF00'", s.Fill.Color)
	}

	if s.Borders.Left.Style != BorderThin {
		t.Errorf("Borders.Left.Style = %q, want 'thin'", s.Borders.Left.Style)
	}
	if s.Borders.Top.Style != BorderMedium {
		t.Errorf("Borders.Top.Style = %q, want 'medium'", s.Borders.Top.Style)
	}

	if s.Alignment.Horizontal != AlignCenter {
		t.Errorf("Alignment.Horizontal = %q, want 'center'", s.Alignment.Horizontal)
	}
	if s.Alignment.Vertical != AlignMiddle {
		t.Errorf("Alignment.Vertical = %q, want 'center'", s.Alignment.Vertical)
	}
	if !s.Alignment.WrapText {
		t.Error("Alignment.WrapText should be true")
	}
	if s.Alignment.TextRotation != 45 {
		t.Errorf("Alignment.TextRotation = %d, want 45", s.Alignment.TextRotation)
	}

	if s.NumberFormat.FormatCode != "0.00" {
		t.Errorf("NumberFormat.FormatCode = %q, want '0.00'", s.NumberFormat.FormatCode)
	}
}

func TestDefaultFont(t *testing.T) {
	f := DefaultFont()
	if f.Name != "Calibri" {
		t.Errorf("default font name = %q, want 'Calibri'", f.Name)
	}
	if f.Size != 11 {
		t.Errorf("default font size = %f, want 11", f.Size)
	}
	if f.Bold {
		t.Error("default font should not be bold")
	}
}

func TestNumberFormats(t *testing.T) {
	formats := []struct {
		name string
		nf   NumberFormat
		code string
	}{
		{"General", FormatGeneral, "General"},
		{"Number", FormatNumber, "0"},
		{"Number2Dec", FormatNumber2Dec, "0.00"},
		{"Percent", FormatPercent, "0%"},
		{"Percent2Dec", FormatPercent2Dec, "0.00%"},
		{"Date", FormatDate, "yyyy-mm-dd"},
		{"DateTime", FormatDateTime, "yyyy-mm-dd hh:mm:ss"},
		{"Time", FormatTime, "hh:mm:ss"},
		{"Text", FormatText, "@"},
	}

	for _, tt := range formats {
		t.Run(tt.name, func(t *testing.T) {
			if tt.nf.FormatCode != tt.code {
				t.Errorf("%s format code = %q, want %q", tt.name, tt.nf.FormatCode, tt.code)
			}
		})
	}
}

func TestBorderStyles(t *testing.T) {
	styles := []BorderStyle{
		BorderNone, BorderThin, BorderMedium, BorderThick,
		BorderDashed, BorderDotted, BorderDouble,
	}
	for _, s := range styles {
		if s == "" {
			t.Error("border style should not be empty")
		}
	}
}
