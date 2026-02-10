package internal

import "testing"

func TestChannelValue_AsFraction(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		isPercent bool
		expected  float64
	}{
		{"absolute 0.5", "0.5", false, 0.5},
		{"percent 50%", "50", true, 0.5},
		{"percent 100%", "100", true, 1.0},
		{"percent 0%", "0", true, 0.0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cv, err := NewChannelValue(tt.value, tt.isPercent)
			if err != nil {
				t.Fatal(err)
			}
			if cv.AsFraction() != tt.expected {
				t.Errorf("expected %f, got %f", tt.expected, cv.AsFraction())
			}
		})
	}
}

func TestRGBChannel_ToRGB(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		isPercent bool
		expected  float64
	}{
		{"255 absolute", "255", false, 255.0},
		{"100% as 255", "100", true, 255.0},
		{"50% as 127.5", "50", true, 127.5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rc, err := NewRGBChannel(tt.value, tt.isPercent)
			if err != nil {
				t.Fatal(err)
			}
			if rc.ToRGB() != tt.expected {
				t.Errorf("expected %f, got %f", tt.expected, rc.ToRGB())
			}
		})
	}
}

func TestLightnessChannel_ToFraction(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		isPercent bool
		expected  float64
	}{
		{"absolute 0.5", "0.5", false, 0.5},
		{"percent 50%", "50", true, 0.5},
		{"percent 100%", "100", true, 1.0},
		{"absolute 1.0", "1.0", false, 1.0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lc, err := NewLightnessChannel(tt.value, tt.isPercent)
			if err != nil {
				t.Fatal(err)
			}
			if lc.ToFraction() != tt.expected {
				t.Errorf("expected %f, got %f", tt.expected, lc.ToFraction())
			}
		})
	}
}

func TestChromaChannel_Value(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected float64
	}{
		{"zero", "0", 0.0},
		{"max chroma", "0.4", 0.4},
		{"clamped above max", "1.0", 0.4},
		{"small chroma", "0.1", 0.1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cc, err := NewChromaChannel(tt.value)
			if err != nil {
				t.Fatal(err)
			}
			if cc.Value() != tt.expected {
				t.Errorf("expected %f, got %f", tt.expected, cc.Value())
			}
		})
	}
}

func TestHueChannel_Value(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected float64
	}{
		{"zero", "0", 0.0},
		{"90 degrees", "90", 90.0},
		{"180 degrees", "180", 180.0},
		{"270 degrees", "270", 270.0},
		{"360 degrees", "360", 360.0},
		{"clamped above 360", "400", 360.0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hc, err := NewHueChannel(tt.value)
			if err != nil {
				t.Fatal(err)
			}
			if hc.Value() != tt.expected {
				t.Errorf("expected %f, got %f", tt.expected, hc.Value())
			}
		})
	}
}

func TestChannelValue_NegativeError(t *testing.T) {
	_, err := NewChannelValue("-10", false)
	if err == nil {
		t.Error("expected error for negative value, got nil")
	}
}
