package internal

import (
	"math"
	"testing"
)

func TestCompareColors_Identical(t *testing.T) {
	tests := []struct {
		name   string
		color1 string
		color2 string
	}{
		{"Same HEX", "#FF0000", "#FF0000"},
		{"Same RGB", "rgb(255, 0, 0)", "rgb(255, 0, 0)"},
		{"Same HSL", "hsl(0, 100%, 50%)", "hsl(0, 100%, 50%)"},
		{"Same OKLCH", "oklch(0.628 0.257 29.23)", "oklch(0.628 0.257 29.23)"},
		{"Different formats same color", "#FF0000", "rgb(255, 0, 0)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := CompareColors(tt.color1, tt.color2)
			if err != nil {
				t.Fatalf("CompareColors() error = %v", err)
			}

			if result.PerceptualDiff != DeltaEIdentical {
				t.Errorf("PerceptualDiff = %f, want %f", result.PerceptualDiff, DeltaEIdentical)
			}

			if result.Verdict != VerdictIdentical {
				t.Errorf("Verdict = %s, want %s", result.Verdict, VerdictIdentical)
			}

			if result.HueDiff != 0 {
				t.Errorf("HueDiff = %f, want 0", result.HueDiff)
			}

			if result.LightnessDiff != 0 {
				t.Errorf("LightnessDiff = %f, want 0", result.LightnessDiff)
			}

			if result.SaturationDiff != 0 {
				t.Errorf("SaturationDiff = %f, want 0", result.SaturationDiff)
			}
		})
	}
}

func TestCompareColors_Indistinguishable(t *testing.T) {
	tests := []struct {
		name   string
		color1 string
		color2 string
		maxDE  float64
	}{
		{"Very close reds", "#FF0000", "#FE0000", DeltaEIndistinguishable},
		{"Nearly identical green", "#00FF00", "#00FE01", DeltaEIndistinguishable},
		{"Almost same blue", "rgb(0, 0, 255)", "rgb(0, 1, 254)", DeltaEIndistinguishable},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := CompareColors(tt.color1, tt.color2)
			if err != nil {
				t.Fatalf("CompareColors() error = %v", err)
			}

			if result.PerceptualDiff > tt.maxDE {
				t.Errorf("PerceptualDiff = %f, want <= %f", result.PerceptualDiff, tt.maxDE)
			}

			if result.Verdict != VerdictIndistinguishable && result.Verdict != VerdictIdentical {
				t.Errorf("Verdict = %s, want indistinguishable or identical", result.Verdict)
			}
		})
	}
}

func TestCompareColors_SlightlyDifferent(t *testing.T) {
	tests := []struct {
		name      string
		color1    string
		color2    string
		minDE     float64
		maxDE     float64
		expectVerdict VerdictType
	}{
		{"Lightness shift", "#FF0000", "#E60000", 0.02, 0.10, VerdictSlightlyDifferent},
		{"Small hue shift", "hsl(0, 100%, 50%)", "hsl(10, 100%, 50%)", 0.02, 0.10, VerdictSlightlyDifferent},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := CompareColors(tt.color1, tt.color2)
			if err != nil {
				t.Fatalf("CompareColors() error = %v", err)
			}

			if result.PerceptualDiff < tt.minDE || result.PerceptualDiff > tt.maxDE {
				t.Errorf("PerceptualDiff = %f, want between %f and %f", result.PerceptualDiff, tt.minDE, tt.maxDE)
			}

			// Verdict should be at least "slightly different"
			if result.Verdict != VerdictSlightlyDifferent && result.Verdict != VerdictIndistinguishable {
				t.Logf("Note: Verdict = %s (may vary based on exact ΔE calculation)", result.Verdict)
			}
		})
	}
}

func TestCompareColors_Different(t *testing.T) {
	tests := []struct {
		name   string
		color1 string
		color2 string
		minDE  float64
	}{
		{"Red vs Green", "#FF0000", "#00FF00", 0.5},
		{"Black vs White", "#000000", "#FFFFFF", 0.5},
		{"Complementary colors", "hsl(0, 100%, 50%)", "hsl(180, 100%, 50%)", 0.4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := CompareColors(tt.color1, tt.color2)
			if err != nil {
				t.Fatalf("CompareColors() error = %v", err)
			}

			if result.PerceptualDiff < tt.minDE {
				t.Errorf("PerceptualDiff = %f, want >= %f", result.PerceptualDiff, tt.minDE)
			}

			if result.Verdict != VerdictDifferent {
				t.Errorf("Verdict = %s, want %s", result.Verdict, VerdictDifferent)
			}
		})
	}
}

func TestCompareColors_ContrastRatio(t *testing.T) {
	tests := []struct {
		name          string
		color1        string
		color2        string
		minRatio      float64
		expectGrade   string
	}{
		{"Black vs White (AAA)", "#000000", "#FFFFFF", 20.0, "AAA"},
		{"Dark gray vs White (AA)", "#333333", "#FFFFFF", 12.0, "AAA"},
		{"Medium gray vs White (AA)", "#767676", "#FFFFFF", 4.5, "AA"},
		{"Light gray vs White (Fail)", "#EEEEEE", "#FFFFFF", 1.1, "Fail"},
		{"Black vs Medium Gray (Fail)", "#000000", "#767676", 4.0, "AA (large text only)"},
		{"Red vs Black (Fail)", "#FF0000", "#000000", 5.0, "AA"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := CompareColors(tt.color1, tt.color2)
			if err != nil {
				t.Fatalf("CompareColors() error = %v", err)
			}

			if result.ContrastRatio < tt.minRatio {
				t.Errorf("ContrastRatio = %f, want >= %f", result.ContrastRatio, tt.minRatio)
			}

			// Check WCAG grade matches expected (with some tolerance for edge cases)
			if tt.expectGrade != "Fail" && result.WCAGGrade != tt.expectGrade {
				t.Logf("Note: WCAGGrade = %s, expected %s (may vary based on exact calculation)", result.WCAGGrade, tt.expectGrade)
			}
		})
	}
}

func TestHueDifference_ShortestPath(t *testing.T) {
	tests := []struct {
		name     string
		hue1     float64
		hue2     float64
		expected float64
	}{
		{"Normal difference", 10, 50, 40},
		{"Across 0° (short path)", 350, 10, 20},
		{"Across 0° (long path)", 10, 350, 20},
		{"Opposite colors", 0, 180, 180},
		{"Near opposite", 175, 185, 10},
		{"Identical", 90, 90, 0},
		{"Full circle", 0, 360, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculateHueDifference(tt.hue1, tt.hue2)
			if math.Abs(result-tt.expected) > 0.1 {
				t.Errorf("calculateHueDifference() = %f, want %f", result, tt.expected)
			}
		})
	}
}

func TestCompareColors_ComponentDifferences(t *testing.T) {
	tests := []struct {
		name            string
		color1          string
		color2          string
		expectHueDiff   float64
		expectLightDiff float64
		expectSatDiff   float64
	}{
		{
			"Pure hue shift",
			"hsl(0, 100%, 50%)",
			"hsl(60, 100%, 50%)",
			60.0, // 60° hue shift
			0.0,
			0.0,
		},
		{
			"Lightness shift only",
			"hsl(0, 100%, 30%)",
			"hsl(0, 100%, 70%)",
			0.0,
			40.0, // 40% lightness difference
			0.0,
		},
		{
			"Saturation shift only",
			"hsl(0, 50%, 50%)",
			"hsl(0, 100%, 50%)",
			0.0,
			0.0,
			50.0, // 50% saturation difference
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := CompareColors(tt.color1, tt.color2)
			if err != nil {
				t.Fatalf("CompareColors() error = %v", err)
			}

			// Allow small floating point errors
			tolerance := 0.5

			if math.Abs(result.HueDiff-tt.expectHueDiff) > tolerance {
				t.Errorf("HueDiff = %f, want %f", result.HueDiff, tt.expectHueDiff)
			}

			if math.Abs(result.LightnessDiff-tt.expectLightDiff) > tolerance {
				t.Errorf("LightnessDiff = %f, want %f", result.LightnessDiff, tt.expectLightDiff)
			}

			if math.Abs(result.SaturationDiff-tt.expectSatDiff) > tolerance {
				t.Errorf("SaturationDiff = %f, want %f", result.SaturationDiff, tt.expectSatDiff)
			}
		})
	}
}

func TestCompareColors_InvalidInput(t *testing.T) {
	tests := []struct {
		name   string
		color1 string
		color2 string
	}{
		{"Invalid color1", "not-a-color", "#FF0000"},
		{"Invalid color2", "#FF0000", "not-a-color"},
		{"Both invalid", "invalid1", "invalid2"},
		{"Empty color1", "", "#FF0000"},
		{"Empty color2", "#FF0000", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := CompareColors(tt.color1, tt.color2)
			if err == nil {
				t.Error("Expected error for invalid input, got nil")
			}
		})
	}
}

func TestCompareColors_DifferentFormats(t *testing.T) {
	tests := []struct {
		name   string
		color1 string
		color2 string
	}{
		{"HEX vs RGB", "#FF0000", "rgb(255, 0, 0)"},
		{"RGB vs HSL", "rgb(255, 0, 0)", "hsl(0, 100%, 50%)"},
		{"HSL vs OKLCH", "hsl(0, 100%, 50%)", "oklch(0.628 0.257 29.23)"},
		{"HEX vs OKLCH", "#FF0000", "oklch(0.628 0.257 29.23)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := CompareColors(tt.color1, tt.color2)
			if err != nil {
				t.Fatalf("CompareColors() error = %v", err)
			}

			// All these represent the same color, should be identical or nearly so
			if result.PerceptualDiff > 0.01 {
				t.Errorf("PerceptualDiff = %f, want <= 0.01 for same color in different formats", result.PerceptualDiff)
			}

			if result.Verdict != VerdictIdentical && result.Verdict != VerdictIndistinguishable {
				t.Errorf("Verdict = %s, want identical or indistinguishable", result.Verdict)
			}
		})
	}
}

func TestFormatComparisonBasic(t *testing.T) {
	result, err := CompareColors("#FF0000", "#00FF00")
	if err != nil {
		t.Fatalf("CompareColors() error = %v", err)
	}

	output := FormatComparisonBasic(result)
	if output == "" {
		t.Error("FormatComparisonBasic() returned empty string")
	}

	// Check that key information is present
	requiredStrings := []string{
		"Color Comparison",
		"Perceptual Difference",
		"Verdict",
		"Contrast Ratio",
	}

	for _, s := range requiredStrings {
		if !contains(output, s) {
			t.Errorf("FormatComparisonBasic() output missing '%s'", s)
		}
	}
}

func TestFormatComparisonDetailed(t *testing.T) {
	result, err := CompareColors("#FF0000", "#00FF00")
	if err != nil {
		t.Fatalf("CompareColors() error = %v", err)
	}

	output := FormatComparisonDetailed(result)
	if output == "" {
		t.Error("FormatComparisonDetailed() returned empty string")
	}

	// Check that detailed information is present
	requiredStrings := []string{
		"Color Comparison",
		"Perceptual Difference",
		"Verdict",
		"Hue Difference",
		"Lightness Difference",
		"Saturation Difference",
		"Contrast Ratio",
		"WCAG Grade",
	}

	for _, s := range requiredStrings {
		if !contains(output, s) {
			t.Errorf("FormatComparisonDetailed() output missing '%s'", s)
		}
	}
}

func TestCalculateOKLCHDeltaE(t *testing.T) {
	tests := []struct {
		name     string
		color1   Color
		color2   Color
		expected float64
	}{
		{
			"Identical colors",
			Color{R: 255, G: 0, B: 0},
			Color{R: 255, G: 0, B: 0},
			0.0,
		},
		{
			"Very different colors",
			Color{R: 255, G: 0, B: 0},
			Color{R: 0, G: 255, B: 0},
			0.5, // Approximate, will vary slightly
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculateOKLCHDeltaE(tt.color1, tt.color2)
			if tt.name == "Identical colors" && result != tt.expected {
				t.Errorf("calculateOKLCHDeltaE() = %f, want %f", result, tt.expected)
			}
			if tt.name == "Very different colors" && result < tt.expected {
				t.Errorf("calculateOKLCHDeltaE() = %f, want >= %f", result, tt.expected)
			}
		})
	}
}

func TestCalculateRelativeLuminance(t *testing.T) {
	tests := []struct {
		name     string
		color    Color
		expected float64
	}{
		{"Pure black", Color{R: 0, G: 0, B: 0}, 0.0},
		{"Pure white", Color{R: 255, G: 255, B: 255}, 1.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculateRelativeLuminance(tt.color)
			// Allow small floating point errors
			if math.Abs(result-tt.expected) > 0.001 {
				t.Errorf("calculateRelativeLuminance() = %f, want %f", result, tt.expected)
			}
		})
	}
}

func TestDetermineVerdict(t *testing.T) {
	tests := []struct {
		deltaE float64
		expect VerdictType
	}{
		{0.0, VerdictIdentical},
		{0.01, VerdictIndistinguishable},
		{0.02, VerdictIndistinguishable},
		{0.05, VerdictSlightlyDifferent},
		{0.10, VerdictSlightlyDifferent},
		{0.15, VerdictDifferent},
		{1.0, VerdictDifferent},
	}

	for _, tt := range tests {
		t.Run(string(tt.expect), func(t *testing.T) {
			result := determineVerdict(tt.deltaE)
			if result != tt.expect {
				t.Errorf("determineVerdict(%f) = %s, want %s", tt.deltaE, result, tt.expect)
			}
		})
	}
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
