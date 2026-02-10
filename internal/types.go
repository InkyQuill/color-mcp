package internal

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// ColorFormat represents the type of color format
type ColorFormat string

const (
	FormatHEX   ColorFormat = "hex"
	FormatRGB   ColorFormat = "rgb"
	FormatRGBA  ColorFormat = "rgba"
	FormatHSL   ColorFormat = "hsl"
	FormatHSLA  ColorFormat = "hsla"
	FormatHSB   ColorFormat = "hsb"
	FormatHSV   ColorFormat = "hsv"
	FormatOKLCH ColorFormat = "oklch"
	FormatLAB   ColorFormat = "lab"
	FormatXYZ   ColorFormat = "xyz"
	FormatHWB   ColorFormat = "hwb"
	FormatCMYK  ColorFormat = "cmyk"
)

// Color represents a color in RGB format with optional alpha
type Color struct {
	R, G, B float64 // 0-255
	A       float64 // 0-1, 1 if no alpha
}

// ColorData represents parsed color with format information
type ColorData struct {
	Color    Color
	Format   ColorFormat
	Original string
}

// Regex patterns for color format detection
var (
	hexPattern   = regexp.MustCompile(`^#([0-9a-fA-F]{3}|[0-9a-fA-F]{4}|[0-9a-fA-F]{6}|[0-9a-fA-F]{8})$`)
	rgbPattern   = regexp.MustCompile(`^rgba?\s*\(\s*([0-9]+\.?[0-9]*)(%?)\s*,\s*([0-9]+\.?[0-9]*)(%?)\s*,\s*([0-9]+\.?[0-9]*)(%?)\s*(?:,\s*([0-9]*\.?[0-9]+)\s*)?\)$`)
	hslPattern   = regexp.MustCompile(`^hsla?\s*\(\s*([0-9]+\.?[0-9]*)\s*,\s*([0-9]+\.?[0-9]*)%\s*,\s*([0-9]+\.?[0-9]*)%\s*(?:,\s*([0-9]*\.?[0-9]+)\s*)?\)$`)
	hsbPattern   = regexp.MustCompile(`(?i)^hs[bcv]\s*\(\s*([0-9]+\.?[0-9]*)\s*,\s*([0-9]+\.?[0-9]*)%\s*,\s*([0-9]+\.?[0-9]*)%\s*(?:,\s*([0-9]*\.?[0-9]+)\s*)?\)$`)
	oklchPattern = regexp.MustCompile(`(?i)^oklch\s*\(\s*([0-9]*\.?[0-9]+)(%?)\s+([0-9]*\.?[0-9]+)(?:\s+([0-9]*\.?[0-9]+))?\s*(?:/\s*([0-9]*\.?[0-9]+)\s*)?\)$`)
	labPattern   = regexp.MustCompile(`(?i)^lab\s*\(\s*([0-9]*\.?[0-9]+)\s+(-?[0-9]*\.?[0-9]+)\s+(-?[0-9]*\.?[0-9]+)\s*(?:/\s*([0-9]*\.?[0-9]+)\s*)?\)$`)
	xyzPattern   = regexp.MustCompile(`(?i)^xyz\s*\(\s*(-?[0-9]*\.?[0-9]+)\s+(-?[0-9]*\.?[0-9]+)\s+(-?[0-9]*\.?[0-9]+)\s*(?:/\s*([0-9]*\.?[0-9]+)\s*)?\)$`)
	hwbPattern   = regexp.MustCompile(`(?i)^hwb\s*\(\s*([0-9]+\.?[0-9]*)\s+([0-9]+\.?[0-9]*)%\s+([0-9]+\.?[0-9]*)%\s*(?:/\s*([0-9]*\.?[0-9]+)\s*)?\)$`)
	cmykPattern  = regexp.MustCompile(`(?i)^cmyk\s*\(\s*([0-9]+\.?[0-9]*)%\s+([0-9]+\.?[0-9]*)%\s+([0-9]+\.?[0-9]*)%\s+([0-9]+\.?[0-9]*)%\s*(?:/\s*([0-9]*\.?[0-9]+)\s*)?\)$`)
)

// DetectFormat detects the color format from the input string
func DetectFormat(input string) (ColorData, error) {
	input = strings.TrimSpace(input)

	// Try HEX
	if hexPattern.MatchString(input) {
		color, err := parseHEX(input)
		if err != nil {
			return ColorData{}, err
		}
		return ColorData{
			Color:    color,
			Format:   FormatHEX,
			Original: input,
		}, nil
	}

	// Try RGB/RGBA (both numeric and percentage)
	if rgbPattern.MatchString(input) {
		color, hasAlpha, err := parseRGB(input)
		if err != nil {
			return ColorData{}, err
		}
		format := FormatRGB
		if hasAlpha {
			format = FormatRGBA
		}
		return ColorData{
			Color:    color,
			Format:   format,
			Original: input,
		}, nil
	}

	// Try HSL/HSLA
	if hslPattern.MatchString(input) {
		color, hasAlpha, err := parseHSL(input)
		if err != nil {
			return ColorData{}, err
		}
		format := FormatHSL
		if hasAlpha {
			format = FormatHSLA
		}
		return ColorData{
			Color:    color,
			Format:   format,
			Original: input,
		}, nil
	}

	// Try HSB/HSV
	if hsbPattern.MatchString(input) {
		color, _, err := parseHSB(input)
		if err != nil {
			return ColorData{}, err
		}
		// Detect actual format name from input
		format := FormatHSB
		inputLower := strings.ToLower(input)
		if len(inputLower) >= 4 && inputLower[0:3] == "hsv" {
			format = FormatHSV
		}
		return ColorData{
			Color:    color,
			Format:   format,
			Original: input,
		}, nil
	}

	// Try OKLCH
	if oklchPattern.MatchString(input) {
		color, err := parseOKLCH(input)
		if err != nil {
			return ColorData{}, err
		}
		return ColorData{
			Color:    color,
			Format:   FormatOKLCH,
			Original: input,
		}, nil
	}

	// Try LAB
	if labPattern.MatchString(input) {
		color, err := parseLAB(input)
		if err != nil {
			return ColorData{}, err
		}
		return ColorData{
			Color:    color,
			Format:   FormatLAB,
			Original: input,
		}, nil
	}

	// Try XYZ
	if xyzPattern.MatchString(input) {
		color, err := parseXYZ(input)
		if err != nil {
			return ColorData{}, err
		}
		return ColorData{
			Color:    color,
			Format:   FormatXYZ,
			Original: input,
		}, nil
	}

	// Try HWB
	if hwbPattern.MatchString(input) {
		color, _, err := parseHWB(input)
		if err != nil {
			return ColorData{}, err
		}
		return ColorData{
			Color:    color,
			Format:   FormatHWB,
			Original: input,
		}, nil
	}

	// Try CMYK
	if cmykPattern.MatchString(input) {
		color, err := parseCMYK(input)
		if err != nil {
			return ColorData{}, err
		}
		return ColorData{
			Color:    color,
			Format:   FormatCMYK,
			Original: input,
		}, nil
	}

	return ColorData{}, fmt.Errorf("unrecognized color format: %s", input)
}

// parseHEX parses a HEX color string
func parseHEX(input string) (Color, error) {
	hex := strings.TrimPrefix(input, "#")

	var r, g, b int
	a := 255

	switch len(hex) {
	case 3: // #RGB
		r = parseHexDigit(hex[0]) * 17
		g = parseHexDigit(hex[1]) * 17
		b = parseHexDigit(hex[2]) * 17
	case 4: // #RGBA
		r = parseHexDigit(hex[0]) * 17
		g = parseHexDigit(hex[1]) * 17
		b = parseHexDigit(hex[2]) * 17
		a = parseHexDigit(hex[3]) * 17
	case 6: // #RRGGBB
		r = parseHexByte(hex[0:2])
		g = parseHexByte(hex[2:4])
		b = parseHexByte(hex[4:6])
	case 8: // #RRGGBBAA
		r = parseHexByte(hex[0:2])
		g = parseHexByte(hex[2:4])
		b = parseHexByte(hex[4:6])
		a = parseHexByte(hex[6:8])
	default:
		return Color{}, fmt.Errorf("invalid hex color: %s", input)
	}

	return Color{
		R: float64(r),
		G: float64(g),
		B: float64(b),
		A: float64(a) / RGBMax,
	}, nil
}

func parseHexDigit(c byte) int {
	switch {
	case '0' <= c && c <= '9':
		return int(c - '0')
	case 'a' <= c && c <= 'f':
		return int(c-'a') + 10
	case 'A' <= c && c <= 'F':
		return int(c-'A') + 10
	default:
		return 0
	}
}

func parseHexByte(s string) int {
	b, _ := strconv.ParseInt(s, 16, 0)
	return int(b)
}

// parseRGB parses an RGB/RGBA color string
func parseRGB(input string) (Color, bool, error) {
	matches := rgbPattern.FindStringSubmatch(input)
	if matches == nil {
		return Color{}, false, fmt.Errorf("invalid RGB format: %s", input)
	}

	rChannel, err := NewRGBChannel(matches[rgbRValueIdx], matches[rgbRPercentIdx] == "%")
	if err != nil {
		return Color{}, false, fmt.Errorf("invalid red value: %w", err)
	}

	gChannel, err := NewRGBChannel(matches[rgbGValueIdx], matches[rgbGPercentIdx] == "%")
	if err != nil {
		return Color{}, false, fmt.Errorf("invalid green value: %w", err)
	}

	bChannel, err := NewRGBChannel(matches[rgbBValueIdx], matches[rgbBPercentIdx] == "%")
	if err != nil {
		return Color{}, false, fmt.Errorf("invalid blue value: %w", err)
	}

	r := rChannel.ToRGB()
	g := gChannel.ToRGB()
	b := bChannel.ToRGB()

	a := AlphaMax
	hasAlpha := false
	if matches[rgbAValueIdx] != "" {
		hasAlpha = true
		a, _ = strconv.ParseFloat(matches[rgbAValueIdx], 64)
		a = clamp(a, AlphaMin, AlphaMax)
	}

	return Color{R: r, G: g, B: b, A: a}, hasAlpha, nil
}

// parseHSL parses an HSL/HSLA color string and converts to RGB
func parseHSL(input string) (Color, bool, error) {
	matches := hslPattern.FindStringSubmatch(input)
	if matches == nil {
		return Color{}, false, fmt.Errorf("invalid HSL format: %s", input)
	}

	h, _ := strconv.ParseFloat(matches[1], 64)
	s, _ := strconv.ParseFloat(matches[2], 64)
	l, _ := strconv.ParseFloat(matches[3], 64)

	a := AlphaMax
	hasAlpha := false

	if matches[4] != "" {
		hasAlpha = true
		a, _ = strconv.ParseFloat(matches[4], 64)
	}

	// Clamp values
	h = clamp(h, 0, HueMax)
	s = clamp(s, 0, SaturationMax)
	l = clamp(l, 0, LightnessMax)
	a = clamp(a, AlphaMin, AlphaMax)

	// Convert HSL to RGB
	r, g, b := hslToRGB(h, s, l)

	return Color{R: r, G: g, B: b, A: a}, hasAlpha, nil
}

// parseHSB parses an HSB/HSV color string and converts to RGB
func parseHSB(input string) (Color, bool, error) {
	matches := hsbPattern.FindStringSubmatch(input)
	if matches == nil {
		return Color{}, false, fmt.Errorf("invalid HSB format: %s", input)
	}

	h, _ := strconv.ParseFloat(matches[1], 64)
	s, _ := strconv.ParseFloat(matches[2], 64)
	v, _ := strconv.ParseFloat(matches[3], 64)

	a := AlphaMax
	hasAlpha := false

	if matches[4] != "" {
		hasAlpha = true
		a, _ = strconv.ParseFloat(matches[4], 64)
	}

	// Clamp values
	h = clamp(h, 0, HueMax)
	s = clamp(s, 0, SaturationMax)
	v = clamp(v, 0, SaturationMax)
	a = clamp(a, AlphaMin, AlphaMax)

	// Convert HSB to RGB
	r, g, b := hsbToRGB(h, s, v)

	return Color{R: r, G: g, B: b, A: a}, hasAlpha, nil
}

// parseOKLCH parses an OKLCH color string and converts to RGB
func parseOKLCH(input string) (Color, error) {
	matches := oklchPattern.FindStringSubmatch(input)
	if matches == nil {
		return Color{}, fmt.Errorf("invalid OKLCH format: %s", input)
	}

	// Parse lightness (can be 0-1 or 0-100%)
	lChannel, err := NewLightnessChannel(matches[oklchLValueIdx], matches[oklchLPercentIdx] == "%")
	if err != nil {
		return Color{}, fmt.Errorf("invalid lightness: %w", err)
	}
	l := lChannel.ToFraction()

	// Parse chroma (0-0.4)
	cChannel, err := NewChromaChannel(matches[oklchCValueIdx])
	if err != nil {
		return Color{}, fmt.Errorf("invalid chroma: %w", err)
	}
	c := cChannel.Value()

	// Parse hue (0-360, optional)
	h := 0.0
	if matches[oklchHValueIdx] != "" {
		hChannel, err := NewHueChannel(matches[oklchHValueIdx])
		if err != nil {
			return Color{}, fmt.Errorf("invalid hue: %w", err)
		}
		h = hChannel.Value()
	}

	// Parse alpha
	a := AlphaMax
	if matches[oklchAValueIdx] != "" {
		a, _ = strconv.ParseFloat(matches[oklchAValueIdx], 64)
		a = clamp(a, AlphaMin, AlphaMax)
	}

	r, g, b := oklchToRGB(l, c, h)
	return Color{R: r, G: g, B: b, A: a}, nil
}

// parseLAB parses a LAB color string and converts to RGB
func parseLAB(input string) (Color, error) {
	matches := labPattern.FindStringSubmatch(input)
	if matches == nil {
		return Color{}, fmt.Errorf("invalid LAB format: %s", input)
	}

	l, _ := strconv.ParseFloat(matches[1], 64)
	a, _ := strconv.ParseFloat(matches[2], 64)
	bVal, _ := strconv.ParseFloat(matches[3], 64)

	alpha := AlphaMax
	if matches[4] != "" {
		alpha, _ = strconv.ParseFloat(matches[4], 64)
	}

	// Convert LAB to RGB via XYZ
	r, g, bVal := labToRGB(l, a, bVal)

	return Color{R: r, G: g, B: bVal, A: alpha}, nil
}

// parseXYZ parses an XYZ color string and converts to RGB
func parseXYZ(input string) (Color, error) {
	matches := xyzPattern.FindStringSubmatch(input)
	if matches == nil {
		return Color{}, fmt.Errorf("invalid XYZ format: %s", input)
	}

	x, _ := strconv.ParseFloat(matches[1], 64)
	y, _ := strconv.ParseFloat(matches[2], 64)
	z, _ := strconv.ParseFloat(matches[3], 64)

	alpha := AlphaMax
	if matches[4] != "" {
		alpha, _ = strconv.ParseFloat(matches[4], 64)
	}

	// Convert XYZ to RGB
	r, g, b := xyzToRGB(x, y, z)

	return Color{R: r, G: g, B: b, A: alpha}, nil
}

// parseHWB parses an HWB color string and converts to RGB
func parseHWB(input string) (Color, bool, error) {
	matches := hwbPattern.FindStringSubmatch(input)
	if matches == nil {
		return Color{}, false, fmt.Errorf("invalid HWB format: %s", input)
	}

	h, _ := strconv.ParseFloat(matches[1], 64)
	w, _ := strconv.ParseFloat(matches[2], 64)
	bVal, _ := strconv.ParseFloat(matches[3], 64)

	a := AlphaMax
	hasAlpha := false

	if matches[4] != "" {
		hasAlpha = true
		a, _ = strconv.ParseFloat(matches[4], 64)
	}

	// Clamp values
	h = clamp(h, 0, HueMax)
	w = clamp(w, 0, LightnessMax)
	bVal = clamp(bVal, 0, LightnessMax)
	a = clamp(a, AlphaMin, AlphaMax)

	// Convert HWB to RGB
	r, g, b := hwbToRGB(h, w, bVal)

	return Color{R: r, G: g, B: b, A: a}, hasAlpha, nil
}

// parseCMYK parses a CMYK color string and converts to RGB
func parseCMYK(input string) (Color, error) {
	matches := cmykPattern.FindStringSubmatch(input)
	if matches == nil {
		return Color{}, fmt.Errorf("invalid CMYK format: %s", input)
	}

	c, _ := strconv.ParseFloat(matches[1], 64)
	m, _ := strconv.ParseFloat(matches[2], 64)
	y, _ := strconv.ParseFloat(matches[3], 64)
	k, _ := strconv.ParseFloat(matches[4], 64)

	a := AlphaMax
	if matches[5] != "" {
		a, _ = strconv.ParseFloat(matches[5], 64)
	}

	// Clamp values
	c = clamp(c, 0, SaturationMax)
	m = clamp(m, 0, SaturationMax)
	y = clamp(y, 0, SaturationMax)
	k = clamp(k, 0, SaturationMax)
	a = clamp(a, AlphaMin, AlphaMax)

	// Convert CMYK to RGB
	r, g, b := cmykToRGB(c, m, y, k)

	return Color{R: r, G: g, B: b, A: a}, nil
}

// clamp clamps a value between min and max
func clamp(v, min, max float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}
