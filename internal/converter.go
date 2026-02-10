package internal

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// Convert converts a color from one format to another
// color: input color string
// targetFormat: target format (hex, rgb, hsl, hsla, hsb, oklch, lab, xyz, hwb, cmyk)
// preserveAlpha: whether to preserve alpha channel
func Convert(color string, targetFormat string, preserveAlpha bool) (string, error) {
	// Detect input format
	data, err := DetectFormat(color)
	if err != nil {
		return "", fmt.Errorf("failed to detect color format: %w", err)
	}

	// Parse target format
	format := ColorFormat(strings.ToLower(targetFormat))
	if !isValidFormat(format) {
		return "", fmt.Errorf("invalid target format: %s (supported: hex, rgb, hsl, hsla, hsb, oklch, lab, xyz, hwb, cmyk)", targetFormat)
	}

	// Get RGB values
	r, g, b := data.Color.R, data.Color.G, data.Color.B
	a := data.Color.A

	// Handle alpha preservation
	if !preserveAlpha {
		a = 1.0
	}

	// Convert to target format
	switch format {
	case FormatHEX:
		return formatHEX(r, g, b, a), nil
	case FormatRGB, FormatRGBA:
		return formatRGB(r, g, b, a, format == FormatRGBA), nil
	case FormatHSL, FormatHSLA:
		return formatHSL(r, g, b, a, format == FormatHSLA), nil
	case FormatHSB, FormatHSV:
		return formatHSB(r, g, b, a), nil
	case FormatOKLCH:
		return formatOKLCH(r, g, b, a), nil
	case FormatLAB:
		return formatLAB(r, g, b, a), nil
	case FormatXYZ:
		return formatXYZ(r, g, b, a), nil
	case FormatHWB:
		return formatHWB(r, g, b, a), nil
	case FormatCMYK:
		return formatCMYK(r, g, b, a), nil
	default:
		return "", fmt.Errorf("unsupported target format: %s", format)
	}
}

// isValidFormat checks if a format is valid
func isValidFormat(format ColorFormat) bool {
	switch format {
	case FormatHEX, FormatRGB, FormatRGBA, FormatHSL, FormatHSLA,
		FormatHSB, FormatHSV, FormatOKLCH, FormatLAB, FormatXYZ,
		FormatHWB, FormatCMYK:
		return true
	default:
		return false
	}
}

// formatHEX formats RGB values as HEX
func formatHEX(r, g, b, a float64) string {
	rInt := int(math.Round(r))
	gInt := int(math.Round(g))
	bInt := int(math.Round(b))

	if a < 1.0 {
		aInt := int(math.Round(a * 255))
		return fmt.Sprintf("#%02X%02X%02X%02X", rInt, gInt, bInt, aInt)
	}
	return fmt.Sprintf("#%02X%02X%02X", rInt, gInt, bInt)
}

// formatRGB formats RGB values as rgb() or rgba()
func formatRGB(r, g, b, a float64, includeAlpha bool) string {
	rStr := strconv.FormatFloat(r, 'f', -1, 64)
	gStr := strconv.FormatFloat(g, 'f', -1, 64)
	bStr := strconv.FormatFloat(b, 'f', -1, 64)

	if includeAlpha {
		return fmt.Sprintf("rgba(%s, %s, %s, %.2f)", rStr, gStr, bStr, a)
	}
	return fmt.Sprintf("rgb(%s, %s, %s)", rStr, gStr, bStr)
}

// formatHSL formats RGB values as HSL
func formatHSL(r, g, b, a float64, includeAlpha bool) string {
	h, s, l := rgbToHSL(r, g, b)

	hStr := strconv.FormatFloat(h, 'f', -1, 64)
	sStr := strconv.FormatFloat(s, 'f', -1, 64)
	lStr := strconv.FormatFloat(l, 'f', -1, 64)

	if includeAlpha {
		return fmt.Sprintf("hsla(%s, %s%%, %s%%, %.2f)", hStr, sStr, lStr, a)
	}
	return fmt.Sprintf("hsl(%s, %s%%, %s%%)", hStr, sStr, lStr)
}

// formatHSB formats RGB values as HSB
func formatHSB(r, g, b, a float64) string {
	h, s, v := rgbToHSB(r, g, b)

	hStr := strconv.FormatFloat(h, 'f', -1, 64)
	sStr := strconv.FormatFloat(s, 'f', -1, 64)
	vStr := strconv.FormatFloat(v, 'f', -1, 64)

	return fmt.Sprintf("hsb(%s, %s%%, %s%%)", hStr, sStr, vStr)
}

// formatOKLCH formats RGB values as OKLCH
func formatOKLCH(r, g, b, a float64) string {
	l, c, h := rgbToOKLCH(r, g, b)

	lStr := strconv.FormatFloat(l, 'f', 4, 64)
	cStr := strconv.FormatFloat(c, 'f', 4, 64)
	hStr := strconv.FormatFloat(h, 'f', 2, 64)

	if a < 1.0 {
		return fmt.Sprintf("oklch(%s %s %s / %.2f)", lStr, cStr, hStr, a)
	}
	return fmt.Sprintf("oklch(%s %s %s)", lStr, cStr, hStr)
}

// formatLAB formats RGB values as LAB
func formatLAB(r, g, b, a float64) string {
	l, aVal, bVal := rgbToLAB(r, g, b)

	lStr := strconv.FormatFloat(l, 'f', 2, 64)
	aStr := strconv.FormatFloat(aVal, 'f', 2, 64)
	bStr := strconv.FormatFloat(bVal, 'f', 2, 64)

	if a < 1.0 {
		return fmt.Sprintf("lab(%s %s %s / %.2f)", lStr, aStr, bStr, a)
	}
	return fmt.Sprintf("lab(%s %s %s)", lStr, aStr, bStr)
}

// formatXYZ formats RGB values as XYZ
func formatXYZ(r, g, b, a float64) string {
	x, y, z := rgbToXYZ(r, g, b)

	xStr := strconv.FormatFloat(x, 'f', 4, 64)
	yStr := strconv.FormatFloat(y, 'f', 4, 64)
	zStr := strconv.FormatFloat(z, 'f', 4, 64)

	if a < 1.0 {
		return fmt.Sprintf("xyz(%s %s %s / %.2f)", xStr, yStr, zStr, a)
	}
	return fmt.Sprintf("xyz(%s %s %s)", xStr, yStr, zStr)
}

// formatHWB formats RGB values as HWB
func formatHWB(r, g, b, a float64) string {
	h, w, bVal := rgbToHWB(r, g, b)

	hStr := strconv.FormatFloat(h, 'f', -1, 64)
	wStr := strconv.FormatFloat(w, 'f', -1, 64)
	bStr := strconv.FormatFloat(bVal, 'f', -1, 64)

	if a < 1.0 {
		return fmt.Sprintf("hwb(%s %s%% %s%% / %.2f)", hStr, wStr, bStr, a)
	}
	return fmt.Sprintf("hwb(%s %s%% %s%%)", hStr, wStr, bStr)
}

// formatCMYK formats RGB values as CMYK
func formatCMYK(r, g, b, a float64) string {
	c, m, y, k := rgbToCMYK(r, g, b)

	cStr := strconv.FormatFloat(c, 'f', 2, 64)
	mStr := strconv.FormatFloat(m, 'f', 2, 64)
	yStr := strconv.FormatFloat(y, 'f', 2, 64)
	kStr := strconv.FormatFloat(k, 'f', 2, 64)

	if a < 1.0 {
		return fmt.Sprintf("cmyk(%s%% %s%% %s%% %s%% / %.2f)", cStr, mStr, yStr, kStr, a)
	}
	return fmt.Sprintf("cmyk(%s%% %s%% %s%% %s%%)", cStr, mStr, yStr, kStr)
}

// GetSupportedFormats returns a list of supported color formats
func GetSupportedFormats() []string {
	return []string{
		"hex", "rgb", "rgba", "hsl", "hsla",
		"hsb", "hsv", "oklch", "lab", "xyz",
		"hwb", "cmyk",
	}
}

// DetectInputFormat returns the format of an input color string
func DetectInputFormat(input string) (string, error) {
	data, err := DetectFormat(input)
	if err != nil {
		return "", err
	}
	return string(data.Format), nil
}
