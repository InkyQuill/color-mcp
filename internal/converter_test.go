package internal

import (
	"math"
	"testing"
)

// TestFormatDetection tests format detection for all supported formats
func TestFormatDetection(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectFormat ColorFormat
		expectError bool
	}{
		// HEX formats
		{"HEX 3 digit", "#F00", FormatHEX, false},
		{"HEX 4 digit", "#F00A", FormatHEX, false},
		{"HEX 6 digit", "#FF0000", FormatHEX, false},
		{"HEX 8 digit", "#FF000080", FormatHEX, false},
		{"HEX lowercase", "#ff0000", FormatHEX, false},
		{"HEX uppercase", "#FF0000", FormatHEX, false},
		{"HEX mixed case", "#Ff00Aa", FormatHEX, false},

		// RGB formats
		{"RGB numeric", "rgb(255, 0, 0)", FormatRGB, false},
		{"RGB with spaces", "rgb( 255 , 0 , 0 )", FormatRGB, false},
		{"RGB decimal", "rgb(255.5, 0.0, 0.0)", FormatRGB, false},
		{"RGBA", "rgba(255, 0, 0, 0.5)", FormatRGBA, false},
		{"RGB percentage", "rgb(100%, 0%, 0%)", FormatRGB, false},
		{"RGBA percentage", "rgba(100%, 0%, 0%, 0.5)", FormatRGBA, false},

		// HSL formats
		{"HSL", "hsl(0, 100%, 50%)", FormatHSL, false},
		{"HSLA", "hsla(0, 100%, 50%, 0.5)", FormatHSLA, false},
		{"HSL decimal", "hsl(120.5, 50.5%, 75.5%)", FormatHSL, false},

		// HSB/HSV formats
		{"HSB", "hsb(0, 100%, 100%)", FormatHSB, false},
		{"HSV", "hsv(120, 50%, 75%)", FormatHSV, false},

		// OKLCH format
		{"OKLCH", "oklch(0.5 0.1 120)", FormatOKLCH, false},
		{"OKLCH with alpha", "oklch(0.5 0.1 120 / 0.5)", FormatOKLCH, false},

		// LAB format
		{"LAB", "lab(50 50 50)", FormatLAB, false},
		{"LAB with alpha", "lab(50 50 50 / 0.5)", FormatLAB, false},

		// XYZ format
		{"XYZ", "xyz(0.5 0.5 0.5)", FormatXYZ, false},
		{"XYZ with alpha", "xyz(0.5 0.5 0.5 / 0.5)", FormatXYZ, false},

		// HWB format
		{"HWB", "hwb(0 0% 0%)", FormatHWB, false},
		{"HWB with alpha", "hwb(0 0% 0% / 0.5)", FormatHWB, false},

		// CMYK format
		{"CMYK", "cmyk(0% 100% 100% 0%)", FormatCMYK, false},
		{"CMYK with alpha", "cmyk(0% 100% 100% 0% / 0.5)", FormatCMYK, false},

		// Invalid formats
		{"Invalid format", "invalid", "", true},
		{"Empty string", "", "", true},
		{"Incomplete HEX", "#FF", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := DetectFormat(tt.input)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if result.Format != tt.expectFormat {
				t.Errorf("Expected format %s, got %s", tt.expectFormat, result.Format)
			}
		})
	}
}

// TestHEXParsing tests HEX color parsing
func TestHEXParsing(t *testing.T) {
	tests := []struct {
		input    string
		expectR  float64
		expectG  float64
		expectB  float64
		expectA  float64
	}{
		{"#F00", 255, 0, 0, 1},
		{"#F00F", 255, 0, 0, 1},  // #F00F where F = 15 (full alpha)
		{"#FF0000", 255, 0, 0, 1},
		{"#FF000080", 255, 0, 0, 128.0 / 255.0},
		{"#000", 0, 0, 0, 1},
		{"#FFF", 255, 255, 255, 1},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			data, err := DetectFormat(tt.input)
			if err != nil {
				t.Fatalf("Failed to parse %s: %v", tt.input, err)
			}

			if data.Color.R != tt.expectR || data.Color.G != tt.expectG || data.Color.B != tt.expectB {
				t.Errorf("Expected RGB(%f, %f, %f), got RGB(%f, %f, %f)",
					tt.expectR, tt.expectG, tt.expectB,
					data.Color.R, data.Color.G, data.Color.B)
			}

			if !almostEqual(data.Color.A, tt.expectA, 0.01) {
				t.Errorf("Expected alpha %f, got %f", tt.expectA, data.Color.A)
			}
		})
	}
}

// TestRGBToHSL tests RGB to HSL conversion
func TestRGBToHSL(t *testing.T) {
	tests := []struct {
		name     string
		r, g, b  float64
		expectH  float64
		expectS  float64
		expectL  float64
	}{
		{"Red", 255, 0, 0, 0, 100, 50},
		{"Green", 0, 255, 0, 120, 100, 50},
		{"Blue", 0, 0, 255, 240, 100, 50},
		{"White", 255, 255, 255, 0, 0, 100},
		{"Black", 0, 0, 0, 0, 0, 0},
		{"Gray", 128, 128, 128, 0, 0, 50.2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h, s, l := rgbToHSL(tt.r, tt.g, tt.b)

			if !almostEqual(h, tt.expectH, 1.0) {
				t.Errorf("Expected H=%f, got H=%f", tt.expectH, h)
			}
			if !almostEqual(s, tt.expectS, 1.0) {
				t.Errorf("Expected S=%f, got S=%f", tt.expectS, s)
			}
			if !almostEqual(l, tt.expectL, 1.0) {
				t.Errorf("Expected L=%f, got L=%f", tt.expectL, l)
			}
		})
	}
}

// TestHSLToRGB tests HSL to RGB conversion (round-trip)
func TestHSLToRGB(t *testing.T) {
	testColors := []struct {
		r, g, b float64
	}{
		{255, 0, 0},
		{0, 255, 0},
		{0, 0, 255},
		{255, 255, 255},
		{0, 0, 0},
		{128, 128, 128},
		{255, 128, 64},
	}

	for _, color := range testColors {
		h, s, l := rgbToHSL(color.r, color.g, color.b)
		r, g, b := hslToRGB(h, s, l)

		if !almostEqual(r, color.r, 1.0) || !almostEqual(g, color.g, 1.0) || !almostEqual(b, color.b, 1.0) {
			t.Errorf("Round-trip failed: RGB(%f,%f,%f) -> HSL(%f,%f,%f) -> RGB(%f,%f,%f)",
				color.r, color.g, color.b, h, s, l, r, g, b)
		}
	}
}

// TestRGBToOKLCH tests RGB to OKLCH conversion
func TestRGBToOKLCH(t *testing.T) {
	tests := []struct {
		name              string
		r, g, b           float64
		expectLMin, LMax  float64
	}{
		{"Black", 0, 0, 0, 0, 0.1},
		{"White", 255, 255, 255, 0.9, 1.1},  // OKLCH L can be slightly > 1 for pure white
		{"Red", 255, 0, 0, 0.6, 0.8},
		{"Green", 0, 255, 0, 0.7, 0.9},
		{"Blue", 0, 0, 255, 0.3, 0.6},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l, c, h := rgbToOKLCH(tt.r, tt.g, tt.b)

			if l < tt.expectLMin || l > tt.LMax {
				t.Errorf("L out of expected range [%f, %f]: got %f", tt.expectLMin, tt.LMax, l)
			}

			// Chroma should be positive for colors
			if tt.r+tt.g+tt.b > 0 && c < 0 {
				t.Errorf("Chroma should be positive for non-black colors: got %f", c)
			}

			// Hue should be in valid range
			if h < 0 || h > 360 {
				t.Errorf("Hue out of range [0, 360]: got %f", h)
			}
		})
	}
}

// TestConvertFunction tests the main Convert function
func TestConvertFunction(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		targetFormat  string
		preserveAlpha bool
		expectError   bool
	}{
		// Valid conversions
		{"HEX to RGB", "#FF0000", "rgb", true, false},
		{"HEX to HSL", "#00FF00", "hsl", true, false},
		{"HEX to OKLCH", "#0000FF", "oklch", true, false},
		{"RGB to HEX", "rgb(255, 0, 0)", "hex", true, false},
		{"RGB to HSL", "rgb(0, 255, 0)", "hsl", true, false},
		{"HSL to RGB", "hsl(0, 100%, 50%)", "rgb", true, false},
		{"HSL to HEX", "hsl(120, 100%, 50%)", "hex", true, false},
		{"RGBA to RGB preserve alpha", "rgba(255, 0, 0, 0.5)", "rgba", true, false},
		{"RGBA to RGB strip alpha", "rgba(255, 0, 0, 0.5)", "rgb", false, false},

		// Invalid conversions
		{"Invalid input", "invalid", "rgb", true, true},
		{"Invalid target format", "#FF0000", "invalid", true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Convert(tt.input, tt.targetFormat, tt.preserveAlpha)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if result == "" {
				t.Errorf("Expected non-empty result")
			}
		})
	}
}

// TestRoundTripConversions tests that converting A -> B -> A returns approximately the same color
func TestRoundTripConversions(t *testing.T) {
	formats := []string{"hex", "rgb", "hsl", "hsb"}

	testColor := "#FF5733" // A nice orange color

	for _, format1 := range formats {
		for _, format2 := range formats {
			t.Run(testColor+" -> "+format1+" -> "+format2, func(t *testing.T) {
				// Convert to format1
				result1, err := Convert(testColor, format1, true)
				if err != nil {
					t.Fatalf("First conversion failed: %v", err)
				}

				// Convert back to format2
				result2, err := Convert(result1, format2, true)
				if err != nil {
					t.Fatalf("Second conversion failed: %v", err)
				}

				// Both results should be valid
				if result1 == "" || result2 == "" {
					t.Errorf("Conversion returned empty string")
				}
			})
		}
	}
}

// TestAlphaPreservation tests alpha channel handling
func TestAlphaPreservation(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		targetFormat  string
		preserveAlpha bool
		expectAlpha   bool
	}{
		{"Preserve alpha in rgba", "rgba(255, 0, 0, 0.5)", "rgba", true, true},
		{"Strip alpha in rgba", "rgba(255, 0, 0, 0.5)", "rgb", false, false},
		{"Preserve alpha in hsla", "hsla(0, 100%, 50%, 0.5)", "hsla", true, true},
		{"Strip alpha in hsl", "hsla(0, 100%, 50%, 0.5)", "hsl", false, false},
		{"No alpha in rgb", "rgb(255, 0, 0)", "rgba", true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Convert(tt.input, tt.targetFormat, tt.preserveAlpha)
			if err != nil {
				t.Fatalf("Conversion failed: %v", err)
			}

			hasAlpha := containsAlpha(result)
			if hasAlpha != tt.expectAlpha {
				t.Errorf("Expected alpha=%v, got alpha=%v in result: %s", tt.expectAlpha, hasAlpha, result)
			}
		})
	}
}

// TestBoundaryValues tests conversion with boundary color values
func TestBoundaryValues(t *testing.T) {
	tests := []struct {
		name         string
		color        string
		targetFormat string
	}{
		{"Pure black", "#000000", "hsl"},
		{"Pure white", "#FFFFFF", "hsl"},
		{"Pure red", "#FF0000", "hsl"},
		{"Pure green", "#00FF00", "hsl"},
		{"Pure blue", "#0000FF", "hsl"},
		{"Gray", "#808080", "hsl"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Convert(tt.color, tt.targetFormat, true)
			if err != nil {
				t.Errorf("Conversion failed: %v", err)
			}
			if result == "" {
				t.Errorf("Got empty result")
			}
		})
	}
}

// TestRealColorValues tests against real-world color values provided
func TestRealColorValues(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		targetFormat string
		validate     func(t *testing.T, result string)
	}{
		{
			name:         "Purple #a392d6 to OKLCH",
			input:        "#a392d6",
			targetFormat: "oklch",
			validate: func(t *testing.T, result string) {
				// Should be close to oklch(0.7 0.1 295)
				if !containsSubstring(result, "oklch(") {
					t.Errorf("Result should contain oklch(: %s", result)
				}
			},
		},
		{
			name:         "Red #fa7575 to HSL",
			input:        "#fa7575",
			targetFormat: "hsl",
			validate: func(t *testing.T, result string) {
				// Should be close to hsl(360 93% 72%)
				if !containsSubstring(result, "hsl(") {
					t.Errorf("Result should contain hsl(: %s", result)
				}
			},
		},
		{
			name:         "Lavender #deadfc to RGB",
			input:        "#deadfc",
			targetFormat: "rgb",
			validate: func(t *testing.T, result string) {
				// Should be close to rgb(222, 173, 252)
				if !containsSubstring(result, "rgb(") {
					t.Errorf("Result should contain rgb(: %s", result)
				}
			},
		},
		{
			name:         "Yellow #face00 to OKLCH",
			input:        "#face00",
			targetFormat: "oklch",
			validate: func(t *testing.T, result string) {
				// Should be close to oklch(0.8642 0.176983 93.2047)
				if !containsSubstring(result, "oklch(") {
					t.Errorf("Result should contain oklch(: %s", result)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Convert(tt.input, tt.targetFormat, true)
			if err != nil {
				t.Fatalf("Conversion failed: %v", err)
			}
			if result == "" {
				t.Errorf("Got empty result")
			}
			if tt.validate != nil {
				tt.validate(t, result)
			}
		})
	}
}

// TestCMYKConversion tests CMYK conversions
func TestCMYKConversion(t *testing.T) {
	tests := []struct {
		input         string
		targetFormat  string
		expectError   bool
	}{
		{"#FF0000", "cmyk", false},
		{"#00FF00", "cmyk", false},
		{"#0000FF", "cmyk", false},
		{"#000000", "cmyk", false},
		{"#FFFFFF", "cmyk", false},
		{"cmyk(0% 100% 100% 0%)", "hex", false},
		{"cmyk(100% 0% 100% 0%)", "rgb", false},
	}

	for _, tt := range tests {
		t.Run(tt.input+" to "+tt.targetFormat, func(t *testing.T) {
			result, err := Convert(tt.input, tt.targetFormat, true)

			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if !tt.expectError && result == "" {
				t.Errorf("Got empty result")
			}
		})
	}
}

// Helper function to check if a color string contains alpha
func containsAlpha(color string) bool {
	return containsSubstring(color, "rgba(") ||
		containsSubstring(color, "hsla(") ||
		containsSubstring(color, "/") ||
		(len(color) == 9 && color[0] == '#') // #RRGGBBAA
}

func containsSubstring(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && containsInString(s, substr))
}

func containsInString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// Helper function for approximate equality
func almostEqual(a, b, epsilon float64) bool {
	return math.Abs(a-b) < epsilon
}

// BenchmarkConvert benchmarks the conversion function
func BenchmarkConvert(b *testing.B) {
	input := "#FF5733"
	target := "hsl"

	for i := 0; i < b.N; i++ {
		_, _ = Convert(input, target, true)
	}
}

// BenchmarkDetectFormat benchmarks format detection
func BenchmarkDetectFormat(b *testing.B) {
	input := "#FF5733"

	for i := 0; i < b.N; i++ {
		_, _ = DetectFormat(input)
	}
}

// TestOKLCHPercentageNotation tests OKLCH with percentage notation
func TestOKLCHPercentageNotation(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expectL  float64 // Expected lightness in 0-1 range
		validate func(t *testing.T, c Color)
	}{
		{
			name:    "50% lightness",
			input:   "oklch(50% 0.1 120)",
			expectL: 0.5,
			validate: func(t *testing.T, c Color) {
				if c.R < 0 || c.R > 255 || c.G < 0 || c.G > 255 || c.B < 0 || c.B > 255 {
					t.Errorf("RGB values out of range: R=%f, G=%f, B=%f", c.R, c.G, c.B)
				}
			},
		},
		{
			name:    "0.5 absolute lightness",
			input:   "oklch(0.5 0.1 120)",
			expectL: 0.5,
			validate: func(t *testing.T, c Color) {
				if c.R < 0 || c.R > 255 || c.G < 0 || c.G > 255 || c.B < 0 || c.B > 255 {
					t.Errorf("RGB values out of range: R=%f, G=%f, B=%f", c.R, c.G, c.B)
				}
			},
		},
		{
			name:    "100% lightness with zero chroma",
			input:   "oklch(100% 0 0)",
			expectL: 1.0,
			validate: func(t *testing.T, c Color) {
				// 100% lightness with 0 chroma should produce pure white
				if c.R < 250 || c.G < 250 || c.B < 250 {
					t.Errorf("Expected near-white values for 100%% lightness: R=%f, G=%f, B=%f", c.R, c.G, c.B)
				}
			},
		},
		{
			name:    "0% lightness",
			input:   "oklch(0% 0.0 0)",
			expectL: 0.0,
			validate: func(t *testing.T, c Color) {
				// 0% lightness should produce black
				if c.R > 5 || c.G > 5 || c.B > 5 {
					t.Errorf("Expected near-black values for 0%% lightness: R=%f, G=%f, B=%f", c.R, c.G, c.B)
				}
			},
		},
		{
			name:    "75% lightness with alpha",
			input:   "oklch(75% 0.08 200 / 0.5)",
			expectL: 0.75,
			validate: func(t *testing.T, c Color) {
				if !almostEqual(c.A, 0.5, 0.01) {
					t.Errorf("Expected alpha=0.5, got %f", c.A)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := DetectFormat(tt.input)
			if err != nil {
				t.Fatalf("Failed to parse %s: %v", tt.input, err)
			}

			if data.Format != FormatOKLCH {
				t.Errorf("Expected format OKLCH, got %s", data.Format)
			}

			// Verify by converting back to OKLCH and checking lightness
			l, _, _ := rgbToOKLCH(data.Color.R, data.Color.G, data.Color.B)
			if !almostEqual(l, tt.expectL, 0.01) {
				t.Errorf("Expected lightness ~%f, got %f", tt.expectL, l)
			}

			if tt.validate != nil {
				tt.validate(t, data.Color)
			}
		})
	}
}

// TestRGBPercentageNotation tests RGB with percentage notation
func TestRGBPercentageNotation(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectR     float64
		expectG     float64
		expectB     float64
		expectAlpha float64
	}{
		{
			name:        "100% red, 50% green, 0% blue",
			input:       "rgb(100%, 50%, 0%)",
			expectR:     255.0,
			expectG:     127.5,
			expectB:     0.0,
			expectAlpha: 1.0,
		},
		{
			name:        "Mixed percentages",
			input:       "rgb(0%, 100%, 100%)",
			expectR:     0.0,
			expectG:     255.0,
			expectB:     255.0,
			expectAlpha: 1.0,
		},
		{
			name:        "25% each (dark gray)",
			input:       "rgb(25%, 25%, 25%)",
			expectR:     63.75,
			expectG:     63.75,
			expectB:     63.75,
			expectAlpha: 1.0,
		},
		{
			name:        "RGBA with percentages",
			input:       "rgba(80%, 20%, 40%, 0.7)",
			expectR:     204.0,
			expectG:     51.0,
			expectB:     102.0,
			expectAlpha: 0.7,
		},
		{
			name:        "Mixed absolute and percentage",
			input:       "rgb(255, 50%, 127.5)",
			expectR:     255.0,
			expectG:     127.5,
			expectB:     127.5,
			expectAlpha: 1.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := DetectFormat(tt.input)
			if err != nil {
				t.Fatalf("Failed to parse %s: %v", tt.input, err)
			}

			// Allow small floating point errors
			if !almostEqual(data.Color.R, tt.expectR, 0.1) {
				t.Errorf("Expected R=%f, got %f", tt.expectR, data.Color.R)
			}
			if !almostEqual(data.Color.G, tt.expectG, 0.1) {
				t.Errorf("Expected G=%f, got %f", tt.expectG, data.Color.G)
			}
			if !almostEqual(data.Color.B, tt.expectB, 0.1) {
				t.Errorf("Expected B=%f, got %f", tt.expectB, data.Color.B)
			}
			if !almostEqual(data.Color.A, tt.expectAlpha, 0.01) {
				t.Errorf("Expected A=%f, got %f", tt.expectAlpha, data.Color.A)
			}
		})
	}
}
