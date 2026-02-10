package internal

// Color space ranges
const (
	RGBMax          float64 = 255.0
	RGBPercent      float64 = 100.0
	HueMax          float64 = 360.0
	SaturationMax   float64 = 100.0
	LightnessMax    float64 = 100.0
	OKLCH_L_Max     float64 = 1.0
	OKLCH_L_Percent float64 = 100.0
	OKLCH_C_Max     float64 = 0.4
	OKLCH_H_Max     float64 = 360.0
	AlphaMin        float64 = 0.0
	AlphaMax        float64 = 1.0
)

// Gamma correction constants
const (
	sRGBGammaThreshold   = 0.0031308
	sRGBInverseThreshold = 0.04045
	sRGBGammaFactor      = 12.92
	sRGBGammaPower       = 1.0 / 2.4
	sRGBGammaOffset      = 1.055
	sRGBGammaSubtract    = 0.055
)

// LAB conversion constants
const (
	labK = 29.0 * 29.0 * 29.0 / (3.0 * 3.0 * 3.0) // ≈ 903.2962962962963
	labE = 6.0 * 6.0 * 6.0 / (29.0 * 29.0 * 29.0) // ≈ 0.008856451679035631
)

// Circle helpers
const (
	FullCircle = 360.0
	OneThird   = 1.0 / 3.0
	TwoThirds  = 2.0 / 3.0
	OneSixth   = 1.0 / 6.0
)

// Regex group indices
const (
	// RGB groups
	rgbRValueIdx   = 1
	rgbRPercentIdx = 2
	rgbGValueIdx   = 3
	rgbGPercentIdx = 4
	rgbBValueIdx   = 5
	rgbBPercentIdx = 6
	rgbAValueIdx   = 7
	// OKLCH groups
	oklchLValueIdx   = 1
	oklchLPercentIdx = 2
	oklchCValueIdx   = 3
	oklchHValueIdx   = 4
	oklchAValueIdx   = 5
)

// Comparison thresholds (based on OKLCH ΔE research)
const (
	DeltaEIdentical         float64 = 0.0  // Exact match
	DeltaEIndistinguishable float64 = 0.02 // Just Noticeable Difference (JND)
	DeltaESlightlyDifferent float64 = 0.10 // Noticeable but similar
)

// WCAG contrast thresholds
const (
	WCAGAAANormal float64 = 7.0
	WCAGAANormal  float64 = 4.5
	WCAGAALarge   float64 = 3.0
)
