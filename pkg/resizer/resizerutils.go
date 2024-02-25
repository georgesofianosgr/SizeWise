package resizer

import (
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
)

var ErrInvalidModifier = errors.New("invalid modifier")

type Modifiers struct {
	Width            int
	Height           int
	Multiplier       int
	CalculatedWidth  int // This is the Width * Multiplier
	CalculatedHeight int // This is the Height * Multiplier
}

func ParseModifiers(modifier string) (Modifiers, error) {
	var width, height int
	multiplier := 1
	for _, m := range strings.Split(modifier, "-") {
		switch m[0] {
		case 'w':
			width, _ = strconv.Atoi(m[1:])
		case 'h':
			height, _ = strconv.Atoi(m[1:])
		case 'x':
			multiplier, _ = strconv.Atoi(m[1:])
		default:
			return Modifiers{}, fmt.Errorf("cannot parse modifier %v, %w", m, ErrInvalidModifier)
		}
	}

	return Modifiers{width, height, multiplier, width * multiplier, height * multiplier}, nil
}

func StringifyModifiers(modifiers Modifiers) string {
	return fmt.Sprintf("w%d-h%d-x%d", modifiers.Width, modifiers.Height, modifiers.Multiplier)
}

func ModifiedPath(path string, modifiers Modifiers) string {
	modifier := StringifyModifiers(modifiers)
	pathWithoutExt := strings.TrimSuffix(path, filepath.Ext(path))
	extension := filepath.Ext(path)
	pathWithModifier := fmt.Sprintf("%s-%s%s", pathWithoutExt, modifier, extension)

	return pathWithModifier
}
