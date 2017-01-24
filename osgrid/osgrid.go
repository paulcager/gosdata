package osgrid

import (
	"errors"
	"fmt"
	"strings"

	"github.com/golang/geo/r2"
)

const (
	LargeGridSize = 500 * 1000 // 500km
	SmallGridSize = 100 * 1000 // 500km
)

var InvalidGridLetter = errors.New("Invalid grid letter")

func invalid(s string) error {
	return fmt.Errorf("Invalid GridRef: %q", s)
}


// TODO - better just to have a Parse function?
func gridFirstLetter(letter byte) (r2.Point, error) {
	switch letter {
	case 'H':
		return r2.Point{0, 2 * LargeGridSize}, nil
	case 'N':
		return r2.Point{0, 1 * LargeGridSize}, nil
	case 'O':
		return r2.Point{1 * LargeGridSize, 1 * LargeGridSize}, nil
	case 'S':
		return r2.Point{0, 0}, nil
	case 'T':
		return r2.Point{1 * LargeGridSize, 0}, nil
	}
	return r2.Point{}, InvalidGridLetter
}

func gridSecondLetter(letter byte) (r2.Point, error) {
	switch {
	case letter >= 'A' && letter <= 'E':
		return r2.Point{float64(letter-'A') * 100000, 400000}, nil
	case letter >= 'F' && letter <= 'H':
		return r2.Point{float64(letter-'F') * 100000, 300000}, nil
	case letter >= 'J' && letter <= 'K':
		return r2.Point{float64(letter-'G') * 100000, 300000}, nil
	case letter >= 'L' && letter <= 'P':
		return r2.Point{float64(letter-'L') * 100000, 200000}, nil
	case letter >= 'Q' && letter <= 'U':
		return r2.Point{float64(letter-'Q') * 100000, 100000}, nil
	case letter >= 'V' && letter <= 'Z':
		return r2.Point{float64(letter-'V') * 100000, 00000}, nil
	}
	return r2.Point{}, InvalidGridLetter
}

type OSGrid struct {
	r2.Rect
}

func Parse(s string) (OSGrid, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return OSGrid{Rect: r2.EmptyRect()}, invalid(s)
	}

	var p r2.Point
	switch s[0] {
	case 'H':
		p = r2.Point{0, 1000000}
	case 'N':
		p = r2.Point{0, 500000}
	case 'O':
		p = r2.Point{500000, 500000}
	case 'S':
		p = r2.Point{0, 0}
	case 'T':
		p = r2.Point{500000, 0}
	}
	
	r := r2.Rect{}
	
	// TODO
}
