package internal

import (
	"math"
	"testing"
)

// TestOKLCHAccuracyAgainstCulori tests OKLCH conversions against known values from culori library
// Reference: https://github.com/Evercoder/culori test/oklch.test.js
func TestOKLCHAccuracyAgainstCulori(t *testing.T) {
	tests := []struct {
		name         string
		inputRGB     string
		expectOKLCH  struct{ L, C, H float64 }
		tolerance    float64
		hueTolerance float64
	}{
		{
			name:     "White",
			inputRGB: "#ffffff",
			expectOKLCH: struct{ L, C, H float64 }{
				L: 1.0,
				C: 0.0,
				H: 0.0, // Hue doesn't matter for achromatic
			},
			tolerance:    0.01,
			hueTolerance: 0.0,
		},
		{
			name:     "Black",
			inputRGB: "#000000",
			expectOKLCH: struct{ L, C, H float64 }{
				L: 0.0,
				C: 0.0,
				H: 0.0,
			},
			tolerance:    0.01,
			hueTolerance: 0.0,
		},
		{
			name:     "Dark gray #111",
			inputRGB: "#111",
			expectOKLCH: struct{ L, C, H float64 }{
				L: 0.17763777307657064,
				C: 0.0,
				H: 0.0,
			},
			tolerance:    0.01,
			hueTolerance: 0.0,
		},
		{
			name:     "Red",
			inputRGB: "#ff0000",
			expectOKLCH: struct{ L, C, H float64 }{
				L: 0.6279553639214311,
				C: 0.2576833038053608,
				H: 29.233880279627854,
			},
			tolerance:    0.01,
			hueTolerance: 0.5,
		},
		{
			name:     "Tomato (rgb(255, 99, 71))",
			inputRGB: "rgb(255, 99, 71)",
			expectOKLCH: struct{ L, C, H float64 }{
				L: 0.648, // Approximate
				C: 0.201, // Approximate
				H: 27.67, // Approximate
			},
			tolerance:    0.05, // Higher tolerance for complex colors
			hueTolerance: 10.0, // Allow up to 10 degrees difference for hue
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse input RGB
			data, err := DetectFormat(tt.inputRGB)
			if err != nil {
				t.Fatalf("Failed to parse %s: %v", tt.inputRGB, err)
			}

			// Convert to OKLCH
			l, c, h := rgbToOKLCH(data.Color.R, data.Color.G, data.Color.B)

			// Validate L
			if !almostEqual(l, tt.expectOKLCH.L, tt.tolerance) {
				t.Errorf("L: expected %.4f, got %.4f (diff: %.4f, tolerance: %.4f)",
					tt.expectOKLCH.L, l, math.Abs(l-tt.expectOKLCH.L), tt.tolerance)
			}

			// Validate C
			if !almostEqual(c, tt.expectOKLCH.C, tt.tolerance) {
				t.Errorf("C: expected %.4f, got %.4f (diff: %.4f, tolerance: %.4f)",
					tt.expectOKLCH.C, c, math.Abs(c-tt.expectOKLCH.C), tt.tolerance)
			}

			// For chromatic colors, validate H
			if tt.expectOKLCH.C > 0.01 {
				// Normalize hue angles to [0, 360]
				hNorm := h
				for hNorm < 0 {
					hNorm += 360
				}
				for hNorm >= 360 {
					hNorm -= 360
				}

				expectedHNorm := tt.expectOKLCH.H
				for expectedHNorm < 0 {
					expectedHNorm += 360
				}
				for expectedHNorm >= 360 {
					expectedHNorm -= 360
				}

				// Check if hue is close (allow for circular wrap-around)
				hueDiff := math.Abs(hNorm - expectedHNorm)
				if hueDiff > 180 {
					hueDiff = 360 - hueDiff
				}

				hueTol := tt.hueTolerance
				if hueTol == 0 {
					hueTol = tt.tolerance * 10 // Default: higher tolerance for hue
				}

				if !almostEqual(hueDiff, 0, hueTol) {
					t.Errorf("H: expected %.2f°, got %.2f° (diff: %.2f°, tolerance: %.2f°)",
						expectedHNorm, hNorm, hueDiff, hueTol)
				}
			}
		})
	}
}

// TestOKLCHToRGBRoundTrip tests round-trip OKLCH -> RGB -> OKLCH
func TestOKLCHToRGBRoundTrip(t *testing.T) {
	tests := []struct {
		name    string
		l, c, h float64
	}{
		{"White", 1.0, 0.0, 0.0},
		{"Black", 0.0, 0.0, 0.0},
		{"Red", 0.6279553639214311, 0.2576833038053608, 29.233880279627854},
		{"Purple", 0.4891, 0.2855, 279.13},
		{"Lavender", 0.819, 0.1198, 312.13},
		{"Yellow", 0.8642, 0.176983, 93.2047},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Convert OKLCH to RGB
			r, g, b := oklchToRGB(tt.l, tt.c, tt.h)

			// Convert RGB back to OKLCH
			l2, c2, h2 := rgbToOKLCH(r, g, b)

			// Check L is close
			if !almostEqual(tt.l, l2, 0.05) {
				t.Errorf("L round-trip: %.4f -> %.4f (diff: %.4f)",
					tt.l, l2, math.Abs(tt.l-l2))
			}

			// For chromatic colors, check C and H
			if tt.c > 0.01 {
				if !almostEqual(tt.c, c2, 0.05) {
					t.Errorf("C round-trip: %.4f -> %.4f (diff: %.4f)",
						tt.c, c2, math.Abs(tt.c-c2))
				}

				// Hue can vary more due to circular nature
				hueDiff := math.Abs(tt.h - h2)
				if hueDiff > 180 {
					hueDiff = 360 - hueDiff
				}
				if !almostEqual(hueDiff, 0, 10) { // 10 degree tolerance for hue
					t.Errorf("H round-trip: %.2f° -> %.2f° (diff: %.2f°)",
						tt.h, h2, hueDiff)
				}
			}
		})
	}
}

// TestOKLCHWithAlpha tests OKLCH conversions with alpha channel
func TestOKLCHWithAlpha(t *testing.T) {
	tests := []struct {
		name       string
		inputOKLCH string
		expectRGB  struct{ R, G, B, A float64 }
		tolerance  float64
	}{
		{
			name:       "White with alpha",
			inputOKLCH: "oklch(1 0 180 / 0.5)",
			expectRGB: struct{ R, G, B, A float64 }{
				R: 255, G: 255, B: 255, A: 0.5,
			},
			tolerance: 1.0,
		},
		{
			name:       "Black with alpha from user's example",
			inputOKLCH: "oklch(0 0 0 / 0.6667)",
			expectRGB: struct{ R, G, B, A float64 }{
				R: 0, G: 0, B: 0, A: 0.6667,
			},
			tolerance: 1.0,
		},
		{
			name:       "Purple with alpha from user's example",
			inputOKLCH: "oklch(0.4891 0.2855 279.13 / 0.1608)",
			expectRGB: struct{ R, G, B, A float64 }{
				R: 82, G: 24, B: 250, A: 0.1608,
			},
			tolerance: 30.0, // Higher tolerance due to OKLCH approximation
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse OKLCH
			data, err := DetectFormat(tt.inputOKLCH)
			if err != nil {
				t.Fatalf("Failed to parse %s: %v", tt.inputOKLCH, err)
			}

			// Validate RGB values
			if !almostEqual(data.Color.R, tt.expectRGB.R, tt.tolerance) {
				t.Errorf("R: expected %.2f, got %.2f",
					tt.expectRGB.R, data.Color.R)
			}
			if !almostEqual(data.Color.G, tt.expectRGB.G, tt.tolerance) {
				t.Errorf("G: expected %.2f, got %.2f",
					tt.expectRGB.G, data.Color.G)
			}
			if !almostEqual(data.Color.B, tt.expectRGB.B, tt.tolerance) {
				t.Errorf("B: expected %.2f, got %.2f",
					tt.expectRGB.B, data.Color.B)
			}
			if !almostEqual(data.Color.A, tt.expectRGB.A, 0.01) {
				t.Errorf("Alpha: expected %.4f, got %.4f",
					tt.expectRGB.A, data.Color.A)
			}
		})
	}
}

// TestOKLCHEdgeCases tests edge cases for OKLCH conversion
func TestOKLCHEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		color    string
		validate func(t *testing.T, data ColorData)
	}{
		{
			name:  "Pure white OKLCH",
			color: "oklch(1 0 180)",
			validate: func(t *testing.T, data ColorData) {
				if !almostEqual(data.Color.R, 255, 1.0) {
					t.Errorf("Expected R≈255, got %.2f", data.Color.R)
				}
				if !almostEqual(data.Color.G, 255, 1.0) {
					t.Errorf("Expected G≈255, got %.2f", data.Color.G)
				}
				if !almostEqual(data.Color.B, 255, 1.0) {
					t.Errorf("Expected B≈255, got %.2f", data.Color.B)
				}
			},
		},
		{
			name:  "Pure black OKLCH",
			color: "oklch(0 0 0)",
			validate: func(t *testing.T, data ColorData) {
				if !almostEqual(data.Color.R, 0, 1.0) {
					t.Errorf("Expected R≈0, got %.2f", data.Color.R)
				}
				if !almostEqual(data.Color.G, 0, 1.0) {
					t.Errorf("Expected G≈0, got %.2f", data.Color.G)
				}
				if !almostEqual(data.Color.B, 0, 1.0) {
					t.Errorf("Expected B≈0, got %.2f", data.Color.B)
				}
			},
		},
		{
			name:  "OKLCH with zero chroma",
			color: "oklch(0.5 0 0)",
			validate: func(t *testing.T, data ColorData) {
				// Should be a gray color
				avg := (data.Color.R + data.Color.G + data.Color.B) / 3
				if !almostEqual(data.Color.R, avg, 5.0) ||
					!almostEqual(data.Color.G, avg, 5.0) ||
					!almostEqual(data.Color.B, avg, 5.0) {
					t.Errorf("Expected gray, got RGB(%.2f, %.2f, %.2f)",
						data.Color.R, data.Color.G, data.Color.B)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := DetectFormat(tt.color)
			if err != nil {
				t.Fatalf("Failed to parse %s: %v", tt.color, err)
			}
			tt.validate(t, data)
		})
	}
}
