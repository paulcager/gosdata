package gosdata

import (
	"errors"
)

var InvalidGridLetter = errors.New("Invalid grid letter")

type Point struct{ X, Y int }

func gridFirstLetter(letter byte) (Point, error) {
	switch letter {
	case 'H':
		return Point{0, 1000000}, nil
	case 'N':
		return Point{0, 500000}, nil
	case 'O':
		return Point{500000, 500000}, nil
	case 'S':
		return Point{0, 0}, nil
	case 'T':
		return Point{500000, 0}, nil
	}
	return Point{}, InvalidGridLetter
}

func gridSecondLetter(letter byte) (Point, error) {
	switch {
	case letter >= 'A' && letter <= 'E':
		return Point{int(letter-'A') * 100000, 400000}, nil
	case letter >= 'F' && letter <= 'H':
		return Point{int(letter-'F') * 100000, 300000}, nil
	case letter >= 'J' && letter <= 'K':
		return Point{int(letter-'G') * 100000, 300000}, nil
	case letter >= 'L' && letter <= 'P':
		return Point{int(letter-'L') * 100000, 200000}, nil
	case letter >= 'Q' && letter <= 'U':
		return Point{int(letter-'Q') * 100000, 100000}, nil
	case letter >= 'V' && letter <= 'Z':
		return Point{int(letter-'V') * 100000, 00000}, nil
	}
	return Point{}, InvalidGridLetter
}
