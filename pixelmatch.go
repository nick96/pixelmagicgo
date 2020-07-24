package pixelmatch

import (
	"errors"
	"fmt"
)

var (
	ErrNotImplemented = errors.New("Not implemented")
)

// Key to map an option to the provided value. By specifying a type we get a bit
// more type safety and can use `exhaustive`
type optionKey int

const (
	thresholdKey optionKey = iota
	antiAliasDetectionKey
	alphaKey
	antiAliasColorKey
	diffColourKey
	diffColourAltKey
)

// RGB is a representation of an RGB colour with R, G, B mapping to the red
// green and blue respectively.
type RGB struct {
	R, G, B uint8
}

// optionalRGB reprsents an optional RGB value. It's just a wrapper of a pointer
// to an RGB struct but provides the nicety of an `Optional` type.
type optionalRGB struct {
	value *RGB
}

// Generic option. The effect of providing this option is defined by its name and value.
type option struct {
	// Key uniquely identifying the option. This defines what the option is
	// about (e.g. threshold).
	key optionKey
	// Value of the option. This defines the effect that the option has.
	value interface{}
}

// Config mapping options to their typed values.
//
// The option pattern is nice but it doesn't allow us to have a type safe value
// beacuse the value is an empty interface. By parsing the options into a struct
// with typed fields we keep the case of use of the option pattern but get type
// safety in code that uses the config.
type config struct {
	threshold             float32
	antiAliasingDetection bool
	alpha                 float32
	antiAliasingColour    RGB
	diffColour            RGB
	diffColourAlt         optionalRGB
}

// Threshold specifies the matching threshold where a smaller value means more sensitive.
func Threshold(threshold float32) option {
	return option{key: thresholdKey, value: threshold}
}

// AntiAliasDetection specifies whether to include anti-aliasing detection.
func AntiAliasDetection(enable bool) option {
	return option{key: antiAliasDetectionKey, value: enable}
}

// Alpha specifies the opacity of the original image in the diff output.
func Alpha(alpha float32) option {
	return option{key: alphaKey, value: alpha}
}

// AntiAliasColour specifies the colour of anti-aliasing pixels in the output.
func AntiAliasColour(rgb RGB) option {
	return option{key: antiAliasColorKey, value: rgb}
}

// DiffColour specifies the colour of different pixels in the output.
func DiffColour(rgb RGB) option {
	return option{key: diffColourKey, value: rgb}
}

// DiffColourAlt specifies the colour the differentiate between dark on light differences.
func DiffColourAlt(rgb RGB) option {
	return option{key: diffColourAltKey, value: rgb}
}

// PixelMatch compares two images, given as bytes and returns an image (in
// bytes) showing the differences between the two. It also returns the number of
// different pixels.
//
// The behaviour of PixelMatch can be modified by providing options. For more
// info about each option, read its documentation but here I'll provide the
// default values for each:
//
// - Threshold = 0.1
// - AntiAliasDetection = false
// - Alpha = 0.1
// - AntiAliasColour = 255, 255, 0 (yellow)
// - DiffColour  = 255, 0, 0 (red)
// - DiffColourAlg = none
func PixelMatch(actual, expected []byte, options ...option) ([]byte, int, error) {
	config, err := parseOptions(options)
	if err != nil {
		return []byte{}, 0, err
	}

	_ = config

	return []byte{}, 0, ErrNotImplemented
}

func noneRGB() optionalRGB {
	return optionalRGB{nil}
}

func someRGB(value RGB) optionalRGB {
	return optionalRGB{&value}
}

// Check if the optional is the some variant.
func (r optionalRGB) isSome() bool {
	return r.value != nil
}

// Check if the optional is the none variant.
func (r optionalRGB) isNone() bool {
	return r.value == nil
}

// Retrieve the some variant of the optional. If the optional is None and error
// will be returned.
func (r optionalRGB) some() (RGB, error) {
	if r.value == nil {
		return RGB{}, errors.New("Tried to get the Some variant of a None optional")
	}
	return *r.value, nil
}

// Parse the list of type-unsafe options into a type-safe config. If there are
// any any issue with types or such, returns a descriptive error.
func parseOptions(options []option) (config, error) {
	config := config{
		threshold:             0.1,
		antiAliasingDetection: false,
		alpha:                 0.1,
		antiAliasingColour:    RGB{255, 255, 0},
		diffColour:            RGB{255, 0, 0},
		diffColourAlt:         noneRGB(),
	}
	for _, opt := range options {
		switch opt.key {
		case thresholdKey:
			threshold, ok := opt.value.(float32)
			if !ok {
				return config, fmt.Errorf(
					"invalid threshold value %v, expected type float32",
					opt.value,
				)
			}
			config.threshold = threshold
		case antiAliasDetectionKey:
			antiAliasingDetection, ok := opt.value.(bool)
			if !ok {
				return config, fmt.Errorf(
					"invalid anti alias detection value %v, expected type bool",
					opt.value,
				)
			}
			config.antiAliasingDetection = antiAliasingDetection
		case alphaKey:
			alpha, ok := opt.value.(float32)
			if !ok {
				return config, fmt.Errorf("invalid alpha value %v, expected type float32", opt.value)
			}
			config.alpha = alpha
		case antiAliasColorKey:
			antiAliasingColour, ok := opt.value.(RGB)
			if !ok {
				return config, fmt.Errorf(
					"invalid anti aliasing colour value %v, expected type RGB",
					opt.value,
				)
			}
			config.antiAliasingColour = antiAliasingColour
		case diffColourKey:
			diffColour, ok := opt.value.(RGB)
			if !ok {
				return config, fmt.Errorf("invalid diff colour valie %v, expected type RGB", opt.value)
			}
			config.diffColour = diffColour
		case diffColourAltKey:
			diffColourAlt, ok := opt.value.(RGB)
			if !ok {
				return config, fmt.Errorf("invalid diff colour alt value %v, expected type RGB", opt.value)
			}
			config.diffColourAlt = someRGB(diffColourAlt)
		}
	}
	return config, ErrNotImplemented
}
