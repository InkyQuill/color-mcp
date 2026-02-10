package internal

import (
	"fmt"
	"strconv"
)

// ChannelValue represents a numeric value that can be absolute or percentage
type ChannelValue struct {
	value     float64
	isPercent bool
}

func NewChannelValue(strValue string, hasPercent bool) (ChannelValue, error) {
	value, err := strconv.ParseFloat(strValue, 64)
	if err != nil {
		return ChannelValue{}, fmt.Errorf("invalid channel value: %w", err)
	}
	if value < 0 {
		return ChannelValue{}, fmt.Errorf("channel value cannot be negative: %f", value)
	}
	return ChannelValue{value: value, isPercent: hasPercent}, nil
}

func (cv ChannelValue) AsFraction() float64 {
	if cv.isPercent {
		return cv.value / 100.0
	}
	return cv.value
}

func (cv ChannelValue) As255() float64 {
	if cv.isPercent {
		return (cv.value / 100.0) * RGBMax
	}
	return cv.value
}

// RGBChannel represents an RGB color channel
type RGBChannel struct {
	ChannelValue
}

func NewRGBChannel(value string, hasPercent bool) (RGBChannel, error) {
	cv, err := NewChannelValue(value, hasPercent)
	if err != nil {
		return RGBChannel{}, err
	}
	return RGBChannel{ChannelValue: cv}, nil
}

func (rc RGBChannel) ToRGB() float64 {
	return clamp(rc.As255(), 0, RGBMax)
}

// LightnessChannel represents OKLCH lightness (0-1 or 0-100%)
type LightnessChannel struct {
	ChannelValue
}

func NewLightnessChannel(value string, hasPercent bool) (LightnessChannel, error) {
	cv, err := NewChannelValue(value, hasPercent)
	if err != nil {
		return LightnessChannel{}, err
	}
	return LightnessChannel{ChannelValue: cv}, nil
}

func (lc LightnessChannel) ToFraction() float64 {
	return clamp(lc.AsFraction(), 0, OKLCH_L_Max)
}

// ChromaChannel represents OKLCH chroma (0-0.4)
type ChromaChannel struct {
	value float64
}

func NewChromaChannel(value string) (ChromaChannel, error) {
	v, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return ChromaChannel{}, err
	}
	return ChromaChannel{value: clamp(v, 0, OKLCH_C_Max)}, nil
}

func (cc ChromaChannel) Value() float64 { return cc.value }

// HueChannel represents hue angle (0-360)
type HueChannel struct {
	value float64
}

func NewHueChannel(value string) (HueChannel, error) {
	v, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return HueChannel{}, err
	}
	return HueChannel{value: clamp(v, 0, HueMax)}, nil
}

func (hc HueChannel) Value() float64 { return hc.value }
