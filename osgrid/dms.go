package osgrid

import "math"

/* - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -  */
/* Geodesy representation conversion functions                        (c) Chris Veness 2002-2020  */
/*                                                                                   MIT Licence  */
/* www.movable-type.co.uk/scripts/latlong.html                                                    */
/* www.movable-type.co.uk/scripts/js/geodesy/geodesy-library.html#dms                             */
/* - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -  */

/**
 * Latitude/longitude points may be represented as decimal degrees, or subdivided into sexagesimal
 * minutes and seconds. This module provides methods for parsing and representing degrees / minutes
 * / seconds.
 *
 * @module dms
 */

/**
 * Constrain degrees to range -90..+90 (for latitude); e.g. -91 => -89, 91 => 89.
 *
 * @private
 * @param {number} degrees
 * @returns degrees within range -90..+90.
 */
func Wrap90(degrees float64) float64 {
	// avoid rounding due to arithmetic ops if within range
	if -90 <= degrees && degrees <= 90 {
		return degrees
	}

	// latitude wrapping requires a triangle wave function; a general triangle wave is
	//     f(x) = 4a/p ⋅ | (x-p/4)%p - p/2 | - a
	// where a = amplitude, p = period, % = modulo; however, JavaScript '%' is a remainder operator
	// not a modulo operator - for modulo, replace 'x%n' with '((x%n)+n)%n'
	var (
		x = degrees
		a = 90.0
		p = 360.0
	)

	return 4*a/p*math.Abs(
		math.Mod(math.Mod(x-p/4, p)+p, p)-p/2) - a
}

/**
 * Constrain degrees to range -180..+180 (for longitude); e.g. -181 => 179, 181 => -179.
 *
 * @private
 * @param {number} degrees
 * @returns degrees within range -180..+180.
 */
func Wrap180(degrees float64) float64 {
	// avoid rounding due to arithmetic ops if within range
	if -180 <= degrees && degrees <= 180 {
		return degrees
	}

	// longitude wrapping requires a sawtooth wave function; a general sawtooth wave is
	//     f(x) = (2ax/p - p/2) % p - a
	// where a = amplitude, p = period, % = modulo; however, JavaScript '%' is a remainder operator
	// not a modulo operator - for modulo, replace 'x%n' with '((x%n)+n)%n'
	var (
		x = degrees
		a = 180.0
		p = 360.0
	)
	return math.Mod((math.Mod(2*a*x/p-p/2, p))+p, p) - a
}

/**
 * Constrain degrees to range 0..360 (for bearings); e.g. -1 => 359, 361 => 1.
 *
 * @private
 * @param {number} degrees
 * @returns degrees within range 0..360.
 */
func Wrap360(degrees float64) float64 {
	// avoid rounding due to arithmetic ops if within range
	if 0 <= degrees && degrees <= 360 {
		return degrees
	}

	// bearing wrapping requires a sawtooth wave function with a vertical offset equal to the
	// amplitude and a corresponding phase shift; this changes the general sawtooth wave function from
	//     f(x) = (2ax/p - p/2) % p - a
	// to
	//     f(x) = (2ax/p) % p
	// where a = amplitude, p = period, % = modulo; however, JavaScript '%' is a remainder operator
	// not a modulo operator - for modulo, replace 'x%n' with '((x%n)+n)%n'
	var (
		x = degrees
		a = 180.0
		p = 360.0
	)
	return math.Mod((math.Mod(2 * a * x / p, p)) + p, p)
}
