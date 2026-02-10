package internal

import (
	"math"
)

// hslToRGB converts HSL values to RGB
// h: 0-360, s: 0-100, l: 0-100
// Returns RGB values in 0-255 range
func hslToRGB(h, s, l float64) (r, g, b float64) {
	// Convert to 0-1 range
	s /= SaturationMax
	l /= LightnessMax

	if s == 0 {
		// Achromatic (gray)
		return l * RGBMax, l * RGBMax, l * RGBMax
	}

	var q float64
	if l < 0.5 {
		q = l * (1 + s)
	} else {
		q = l + s - l*s
	}

	p := 2*l - q

	hk := h / FullCircle

	// Calculate RGB
	r = hueToRGB(p, q, hk+OneThird)
	g = hueToRGB(p, q, hk)
	b = hueToRGB(p, q, hk-OneThird)

	return r * RGBMax, g * RGBMax, b * RGBMax
}

func hueToRGB(p, q, t float64) float64 {
	if t < 0 {
		t += 1
	}
	if t > 1 {
		t -= 1
	}
	if t < OneSixth {
		return p + (q-p)*6*t
	}
	if t < 0.5 {
		return q
	}
	if t < TwoThirds {
		return p + (q-p)*(TwoThirds-t)*6
	}
	return p
}

// rgbToHSL converts RGB to HSL
// r, g, b: 0-255
// Returns h: 0-360, s: 0-100, l: 0-100
func rgbToHSL(r, g, b float64) (h, s, l float64) {
	// Convert to 0-1 range
	r /= RGBMax
	g /= RGBMax
	b /= RGBMax

	max := math.Max(r, math.Max(g, b))
	min := math.Min(r, math.Min(g, b))
	delta := max - min

	// Lightness
	l = (max + min) / 2

	// Saturation
	if delta == 0 {
		s = 0
		h = 0
	} else {
		if l < 0.5 {
			s = delta / (max + min)
		} else {
			s = delta / (2 - max - min)
		}

		// Hue
		switch {
		case r == max:
			h = (g - b) / delta
		case g == max:
			h = 2 + (b-r)/delta
		default:
			h = 4 + (r-g)/delta
		}

		h *= 60
		if h < 0 {
			h += HueMax
		}
	}

	return h, s * SaturationMax, l * LightnessMax
}

// hsbToRGB converts HSB/HSV values to RGB
// h: 0-360, s: 0-100, v: 0-100
// Returns RGB values in 0-255 range
func hsbToRGB(h, s, v float64) (r, g, b float64) {
	// Convert to 0-1 range
	s /= SaturationMax
	v /= SaturationMax

	c := v * s
	hk := h / 60
	x := c * (1 - math.Abs(math.Mod(hk, 2)-1))

	var r1, g1, b1 float64

	switch {
	case hk < 1:
		r1, g1, b1 = c, x, 0
	case hk < 2:
		r1, g1, b1 = x, c, 0
	case hk < 3:
		r1, g1, b1 = 0, c, x
	case hk < 4:
		r1, g1, b1 = 0, x, c
	case hk < 5:
		r1, g1, b1 = x, 0, c
	default:
		r1, g1, b1 = c, 0, x
	}

	m := v - c

	return (r1 + m) * RGBMax, (g1 + m) * RGBMax, (b1 + m) * RGBMax
}

// rgbToHSB converts RGB to HSB/HSV
// r, g, b: 0-255
// Returns h: 0-360, s: 0-100, v: 0-100
func rgbToHSB(r, g, b float64) (h, s, v float64) {
	// Convert to 0-1 range
	r /= RGBMax
	g /= RGBMax
	b /= RGBMax

	max := math.Max(r, math.Max(g, b))
	min := math.Min(r, math.Min(g, b))
	delta := max - min

	// Value
	v = max

	// Saturation
	if max == 0 {
		s = 0
	} else {
		s = delta / max
	}

	// Hue
	if delta == 0 {
		h = 0
	} else {
		switch {
		case r == max:
			h = (g - b) / delta
		case g == max:
			h = 2 + (b-r)/delta
		default:
			h = 4 + (r-g)/delta
		}

		h *= 60
		if h < 0 {
			h += HueMax
		}
	}

	return h, s * SaturationMax, v * SaturationMax
}

// oklchToRGB converts OKLCH to RGB
// l: 0-1, c: 0-0.4, h: 0-360
// Returns RGB values in 0-255 range
// Based on formulas from culori library
func oklchToRGB(l, c, h float64) (r, g, b float64) {
	// Convert OKLCH to OKLab
	hRad := h * math.Pi / 180
	a := c * math.Cos(hRad)
	bVal := c * math.Sin(hRad)

	// Convert OKLab to LMS (using culori formulas)
	L := math.Pow(l+0.3963377773761749*a+0.2158037573099136*bVal, 3)
	M := math.Pow(l-0.1055613458156586*a-0.0638541728258133*bVal, 3)
	S := math.Pow(l-0.0894841775298119*a-1.2914855480194092*bVal, 3)

	// Convert LMS to linear RGB (using culori formulas)
	rLin := 4.0767416360759574*L - 3.3077115392580616*M + 0.2309699031821044*S
	gLin := -1.2684379732850317*L + 2.6097573492876887*M - 0.3413193760026573*S
	bLin := -0.0041960761386756*L - 0.7034186179359362*M + 1.7076146940746117*S

	// Gamma correction (sRGB)
	r = srgbGamma(rLin) * RGBMax
	g = srgbGamma(gLin) * RGBMax
	b = srgbGamma(bLin) * RGBMax

	return clamp(r, 0, RGBMax), clamp(g, 0, RGBMax), clamp(b, 0, RGBMax)
}

// rgbToOKLCH converts RGB to OKLCH
// r, g, b: 0-255
// Returns l: 0-1, c: 0-0.4, h: 0-360
// Based on formulas from culori library
func rgbToOKLCH(r, g, b float64) (l, c, h float64) {
	// Convert sRGB to linear RGB
	rLin := srgbInverseGamma(r / RGBMax)
	gLin := srgbInverseGamma(g / RGBMax)
	bLin := srgbInverseGamma(b / RGBMax)

	// Convert linear RGB to LMS (using culori formulas)
	cbrtL := cbrt(0.412221469470763*rLin + 0.5363325372617348*gLin + 0.0514459932675022*bLin)
	cbrtM := cbrt(0.2119034958178252*rLin + 0.6806995506452344*gLin + 0.1073969535369406*bLin)
	cbrtS := cbrt(0.0883024591900564*rLin + 0.2817188391361215*gLin + 0.6299787016738222*bLin)

	// Convert LMS to OKLab (using culori formulas)
	l = 0.210454268309314*cbrtL + 0.7936177747023054*cbrtM - 0.0040720430116193*cbrtS
	a := 1.9779985324311684*cbrtL - 2.4285922420485799*cbrtM + 0.450593709617411*cbrtS
	bVal := 0.0259040424655478*cbrtL + 0.7827717124575296*cbrtM - 0.8086757549230774*cbrtS

	// For achromatic colors (gray), set a and b to 0
	if r == g && g == b {
		a = 0
		bVal = 0
	}

	// Convert OKLab to OKLCH
	c = math.Sqrt(a*a + bVal*bVal)
	h = math.Atan2(bVal, a) * 180 / math.Pi
	if h < 0 {
		h += HueMax
	}

	return l, c, h
}

// labToRGB converts LAB to RGB via XYZ
// Using updated XYZ -> RGB matrix from CSS Color Module / culori
func labToRGB(lVal, a, bVal float64) (r, g, b float64) {
	// LAB to XYZ
	y := (lVal + 16) / 116
	x := y + a/500
	z := y - bVal/200

	// Inverse labF function
	fInv := func(f float64) float64 {
		f3 := f * f * f
		if f3 > labE {
			return f3
		}
		return (116*f - 16) / labK
	}

	x = xyzD65[0] * fInv(x)
	y = xyzD65[1] * fInv(y)
	z = xyzD65[2] * fInv(z)

	// XYZ to RGB (using updated CSS Color Module matrix)
	rLin := 3.240969941904521*x - 1.537383177570093*y - 0.498610760293*z
	gLin := -0.96924363628087*x + 1.8759675015077202*y + 0.041555057407175*z
	bLin := 0.055630079696993*x - 0.20397695888897*y + 1.0569715142428786*z

	r = srgbGamma(rLin) * RGBMax
	g = srgbGamma(gLin) * RGBMax
	b = srgbGamma(bLin) * RGBMax

	return clamp(r, 0, RGBMax), clamp(g, 0, RGBMax), clamp(b, 0, RGBMax)
}

// rgbToLAB converts RGB to LAB via XYZ
// Based on culori implementation with achromatic color fix
func rgbToLAB(r, g, b float64) (l, a, bVal float64) {
	// sRGB to linear RGB
	rLin := srgbInverseGamma(r / RGBMax)
	gLin := srgbInverseGamma(g / RGBMax)
	bLin := srgbInverseGamma(b / RGBMax)

	// Linear RGB to XYZ (using updated CSS Color Module matrix)
	x := 0.41239079926595934*rLin + 0.357584339383878*gLin + 0.1804807884018343*bLin
	y := 0.21263900587151027*rLin + 0.715168678767756*gLin + 0.07219231536073371*bLin
	z := 0.019330818715591841*rLin + 0.11919477979462587*gLin + 0.9505321522496607*bLin

	// XYZ to LAB
	xNorm := x / xyzD65[0]
	yNorm := y / xyzD65[1]
	zNorm := z / xyzD65[2]

	fx := labF(xNorm)
	fy := labF(yNorm)
	fz := labF(zNorm)

	l = 116*fy - 16
	a = 500 * (fx - fy)
	bVal = 200 * (fy - fz)

	// Fixes achromatic RGB colors having a slight chroma due to floating-point errors
	// See: https://github.com/d3/d3-color/pull/46
	if r == g && g == b {
		a = 0
		bVal = 0
	}

	return l, a, bVal
}

// xyzToRGB converts XYZ to RGB
// Using inverse sRGB transformation matrix from CSS Color Module / culori
func xyzToRGB(x, y, z float64) (r, g, b float64) {
	rLin := 3.240969941904521*x - 1.537383177570093*y - 0.498610760293*z
	gLin := -0.96924363628087*x + 1.8759675015077202*y + 0.041555057407175*z
	bLin := 0.055630079696993*x - 0.20397695888897*y + 1.0569715142428786*z

	r = srgbGamma(rLin) * RGBMax
	g = srgbGamma(gLin) * RGBMax
	b = srgbGamma(bLin) * RGBMax

	return clamp(r, 0, RGBMax), clamp(g, 0, RGBMax), clamp(b, 0, RGBMax)
}

// rgbToXYZ converts RGB to XYZ
// Using sRGB transformation matrix from CSS Color Module / culori
func rgbToXYZ(r, g, b float64) (x, y, z float64) {
	rLin := srgbInverseGamma(r / RGBMax)
	gLin := srgbInverseGamma(g / RGBMax)
	bLin := srgbInverseGamma(b / RGBMax)

	x = 0.41239079926595934*rLin + 0.357584339383878*gLin + 0.1804807884018343*bLin
	y = 0.21263900587151027*rLin + 0.715168678767756*gLin + 0.07219231536073371*bLin
	z = 0.019330818715591841*rLin + 0.11919477979462587*gLin + 0.9505321522496607*bLin

	return x, y, z
}

// hwbToRGB converts HWB to RGB
// h: 0-360, w: 0-100, b: 0-100
func hwbToRGB(h, w, bVal float64) (r, g, b float64) {
	// Convert to 0-1 range
	w /= LightnessMax
	bVal /= LightnessMax

	// First get RGB from H
	r, g, b = hslToRGB(h, SaturationMax, LightnessMax/2) // Use full saturation/lightness

	// Then mix with white and black
	r = r/RGBMax*(1-w-bVal) + w
	g = g/RGBMax*(1-w-bVal) + w
	b = b/RGBMax*(1-w-bVal) + w

	return r * RGBMax, g * RGBMax, b * RGBMax
}

// rgbToHWB converts RGB to HWB
func rgbToHWB(r, g, b float64) (h, w, bVal float64) {
	// RGB to HSL first (we only need hue)
	h, _, _ = rgbToHSL(r, g, b)

	// Whiteness and blackness
	min := math.Min(r/RGBMax, math.Min(g/RGBMax, b/RGBMax))
	max := math.Max(r/RGBMax, math.Max(g/RGBMax, b/RGBMax))

	w = min * LightnessMax
	bVal = (1 - max) * LightnessMax

	return h, w, bVal
}

// cmykToRGB converts CMYK to RGB
// c, m, y, k: 0-100
func cmykToRGB(c, m, y, k float64) (r, g, b float64) {
	// Convert to 0-1 range
	c /= SaturationMax
	m /= SaturationMax
	y /= SaturationMax
	k /= SaturationMax

	r = (1 - c) * (1 - k) * RGBMax
	g = (1 - m) * (1 - k) * RGBMax
	b = (1 - y) * (1 - k) * RGBMax

	return clamp(r, 0, RGBMax), clamp(g, 0, RGBMax), clamp(b, 0, RGBMax)
}

// rgbToCMYK converts RGB to CMYK
func rgbToCMYK(r, g, b float64) (c, m, y, k float64) {
	// Convert to 0-1 range
	r = 1 - r/RGBMax
	g = 1 - g/RGBMax
	b = 1 - b/RGBMax

	k = math.Min(r, math.Min(g, b))

	c = (r - k) / (1 - k)
	m = (g - k) / (1 - k)
	y = (b - k) / (1 - k)

	if k == 1 {
		c, m, y = 0, 0, 0
	}

	return c * SaturationMax, m * SaturationMax, y * SaturationMax, k * SaturationMax
}

// Helper functions

func srgbGamma(v float64) float64 {
	if v <= sRGBGammaThreshold {
		return sRGBGammaFactor * v
	}
	return sRGBGammaOffset*math.Pow(v, sRGBGammaPower) - sRGBGammaSubtract
}

func srgbInverseGamma(v float64) float64 {
	if v <= sRGBInverseThreshold {
		return v / sRGBGammaFactor
	}
	return math.Pow((v+sRGBGammaSubtract)/sRGBGammaOffset, 2.4)
}

func labF(t float64) float64 {
	if t > labE {
		return math.Cbrt(t)
	}
	return (labK*t + 16) / 116
}

func cbrt(x float64) float64 {
	return math.Pow(x, 1.0/3.0)
}

// D65 illuminant XYZ values (from CSS Color Module / culori)
var xyzD65 = [3]float64{
	0.3127 / 0.329, // ≈ 0.9504559270516716
	1.0,
	(1 - 0.3127 - 0.329) / 0.329, // ≈ 1.08905775075988
}
