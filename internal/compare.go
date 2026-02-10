package internal

import (
	"fmt"
	"math"
)

// VerdictType represents the perceptual similarity verdict
type VerdictType string

const (
	VerdictIdentical         VerdictType = "identical"
	VerdictIndistinguishable VerdictType = "indistinguishable"
	VerdictSlightlyDifferent VerdictType = "slightly different"
	VerdictDifferent         VerdictType = "different"
)

// ComparisonResult contains detailed comparison metrics between two colors
type ComparisonResult struct {
	Color1, Color2 ColorData
	PerceptualDiff float64 // OKLCH ΔE (0-1+)
	Verdict        VerdictType
	HueDiff        float64 // 0-360° (HSL-based)
	LightnessDiff  float64 // 0-100% (HSL-based)
	SaturationDiff float64 // 0-100% (HSL-based)
	ContrastRatio  float64 // WCAG ratio (1-21)
	WCAGGrade      string
}

// CompareColors compares two colors for perceptual similarity, component differences, and contrast ratio
func CompareColors(color1, color2 string) (*ComparisonResult, error) {
	// Parse both colors using existing DetectFormat
	data1, err := DetectFormat(color1)
	if err != nil {
		return nil, fmt.Errorf("invalid color1: %w", err)
	}

	data2, err := DetectFormat(color2)
	if err != nil {
		return nil, fmt.Errorf("invalid color2: %w", err)
	}

	// Calculate all metrics
	deltaE := calculateOKLCHDeltaE(data1.Color, data2.Color)
	verdict := determineVerdict(deltaE)

	h1, s1, l1 := rgbToHSL(data1.Color.R, data1.Color.G, data1.Color.B)
	h2, s2, l2 := rgbToHSL(data2.Color.R, data2.Color.G, data2.Color.B)

	hueDiff := calculateHueDifference(h1, h2)
	lightnessDiff := math.Abs(l2 - l1)
	saturationDiff := math.Abs(s2 - s1)

	contrast := calculateContrastRatio(data1.Color, data2.Color)
	wcagGrade := getWCAGGrade(contrast)

	return &ComparisonResult{
		Color1:         data1,
		Color2:         data2,
		PerceptualDiff: deltaE,
		Verdict:        verdict,
		HueDiff:        hueDiff,
		LightnessDiff:  lightnessDiff,
		SaturationDiff: saturationDiff,
		ContrastRatio:  contrast,
		WCAGGrade:      wcagGrade,
	}, nil
}

// calculateOKLCHDeltaE calculates perceptual difference using OKLCH color space
// OKLCH is perceptually uniform - equal distances correspond to equal perceived differences
func calculateOKLCHDeltaE(c1, c2 Color) float64 {
	l1, c1c, h1 := rgbToOKLCH(c1.R, c1.G, c1.B)
	l2, c2c, h2 := rgbToOKLCH(c2.R, c2.G, c2.B)

	// Convert polar to Cartesian for Euclidean distance
	a1 := c1c * math.Cos(h1*math.Pi/180)
	b1 := c1c * math.Sin(h1*math.Pi/180)
	a2 := c2c * math.Cos(h2*math.Pi/180)
	b2 := c2c * math.Sin(h2*math.Pi/180)

	deltaL := l2 - l1
	deltaA := a2 - a1
	deltaB := b2 - b1

	return math.Sqrt(deltaL*deltaL + deltaA*deltaA + deltaB*deltaB)
}

// calculateHueDifference calculates the shortest path difference between two hues
func calculateHueDifference(h1, h2 float64) float64 {
	diff := math.Abs(h1 - h2)
	if diff > 180 {
		diff = 360 - diff // Shortest path around color wheel
	}
	return diff
}

// calculateContrastRatio calculates WCAG contrast ratio between two colors
func calculateContrastRatio(c1, c2 Color) float64 {
	l1 := calculateRelativeLuminance(c1)
	l2 := calculateRelativeLuminance(c2)

	lighter := math.Max(l1, l2)
	darker := math.Min(l1, l2)

	return (lighter + 0.05) / (darker + 0.05)
}

// calculateRelativeLuminance calculates WCAG relative luminance
func calculateRelativeLuminance(c Color) float64 {
	// Convert to linear RGB
	rLin := srgbInverseGamma(c.R / 255)
	gLin := srgbInverseGamma(c.G / 255)
	bLin := srgbInverseGamma(c.B / 255)

	// WCAG 2.0 formula
	return 0.2126*rLin + 0.7152*gLin + 0.0722*bLin
}

// determineVerdict maps ΔE to human-readable verdict
func determineVerdict(deltaE float64) VerdictType {
	if deltaE == DeltaEIdentical {
		return VerdictIdentical
	}
	if deltaE <= DeltaEIndistinguishable {
		return VerdictIndistinguishable
	}
	if deltaE <= DeltaESlightlyDifferent {
		return VerdictSlightlyDifferent
	}
	return VerdictDifferent
}

// getWCAGGrade returns WCAG grade based on contrast ratio
func getWCAGGrade(contrast float64) string {
	if contrast >= WCAGAAANormal {
		return "AAA"
	}
	if contrast >= WCAGAANormal {
		return "AA"
	}
	if contrast >= WCAGAALarge {
		return "AA (large text only)"
	}
	return "Fail"
}

// FormatComparisonBasic formats comparison result with basic information
func FormatComparisonBasic(result *ComparisonResult) string {
	return fmt.Sprintf(
		"Color Comparison: %s vs %s\n"+
			"Perceptual Difference: %.3f ΔE\n"+
			"Verdict: %s\n"+
			"Contrast Ratio: %.2f:1 (%s)",
		result.Color1.Original, result.Color2.Original,
		result.PerceptualDiff,
		result.Verdict,
		result.ContrastRatio, result.WCAGGrade,
	)
}

// FormatComparisonDetailed formats comparison result with detailed breakdown
func FormatComparisonDetailed(result *ComparisonResult) string {
	return fmt.Sprintf(
		"Color Comparison: %s (%s) vs %s (%s)\n\n"+
			"Perceptual Difference: %.3f ΔE\n"+
			"Verdict: %s\n\n"+
			"Component Breakdown:\n"+
			"  Hue Difference: %.1f°\n"+
			"  Lightness Difference: %.1f%%\n"+
			"  Saturation Difference: %.1f%%\n\n"+
			"Contrast Ratio: %.2f:1\n"+
			"WCAG Grade: %s",
		result.Color1.Original, result.Color1.Format,
		result.Color2.Original, result.Color2.Format,
		result.PerceptualDiff,
		result.Verdict,
		result.HueDiff,
		result.LightnessDiff,
		result.SaturationDiff,
		result.ContrastRatio,
		result.WCAGGrade,
	)
}
