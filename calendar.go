package slagofa

import "github.com/hebl/gofa"

// Calendar and Epoch Conversion Functions
//
// This file implements SLALIB-compatible calendar and epoch conversion functions
// using GoFA (IAU SOFA) as the foundation.

// Cldj converts Gregorian calendar date to Modified Julian Date (MJD)
//
// Original FORTRAN: sla_CLDJ by P.T. Wallace
// GoFA equivalent: gofa.Cal2jd
//
// Parameters:
//   - iy: Year in Gregorian calendar
//   - im: Month (1-12)
//   - id: Day (1-31)
//
// Returns:
//   - djm: Modified Julian Date (JD - 2400000.5)
//   - status: 0 = OK, -1 = bad year, -2 = bad month, -3 = bad day
//
// Notes:
//   - Gregorian calendar is used for all dates
//   - Years may be negative (B.C. dates)
//   - MJD = JD - 2400000.5
func Cldj(iy, im, id int) (djm float64, status int) {
	var djm0, djm1 float64
	status = gofa.Cal2jd(iy, im, id, &djm0, &djm1)
	// GoFA returns two-part JD: djm0 = MJD epoch (2400000.5), djm1 = offset
	// Return just the MJD part for SLALIB compatibility
	return djm1, status
}

// Djcl converts Modified Julian Date to Gregorian calendar date
//
// Original FORTRAN: sla_DJCL by P.T. Wallace
// GoFA equivalent: gofa.Jd2cal
//
// Parameters:
//   - djm: Modified Julian Date (JD - 2400000.5)
//
// Returns:
//   - iy: Year in Gregorian calendar
//   - im: Month (1-12)
//   - id: Day (1-31)
//   - fd: Fraction of day (0.0 to 1.0)
//   - status: 0 = OK, -1 = unacceptable date
//
// Notes:
//   - Valid for any date with |MJD| < 1e9
//   - The earliest date is 4800BC March 1
func Djcl(djm float64) (iy, im, id int, fd float64, status int) {
	// Convert MJD to two-part JD for GoFA
	// djm0 = 2400000.5 (MJD epoch), djm = offset
	status = gofa.Jd2cal(gofa.DJM0, djm, &iy, &im, &id, &fd)
	return iy, im, id, fd, status
}

// Djcal converts Modified Julian Date to Gregorian calendar with formatted output
//
// Original FORTRAN: sla_DJCAL by P.T. Wallace
// GoFA equivalent: gofa.Jd2cal (with formatting)
//
// Parameters:
//   - ndp: Number of decimal places for fraction (recommend ≤ 4 to avoid overflow)
//   - djm: Modified Julian Date (JD - 2400000.5)
//
// Returns:
//   - iy: Year in Gregorian calendar
//   - im: Month (1-12)
//   - id: Day (1-31)
//   - fraction: Fraction of day scaled to ndp decimal places (as integer)
//   - status: 0 = OK, non-zero = out of range
//
// Notes:
//   - Any date after 4701 BC March 1 is accepted
//   - ndp should be 4 or less to avoid integer overflow on 32-bit systems
//   - The fraction is scaled: for ndp=4, 0.9999 days → 9999
//   - This format is convenient for formatted output
//
// Example:
//   - MJD 50123.9999 with ndp=4 → (1996, 2, 10, 9999, 0)
func Djcal(ndp int, djm float64) (iy, im, id int, fraction int, status int) {
	// Get calendar date and fractional day
	var fd float64
	status = gofa.Jd2cal(gofa.DJM0, djm, &iy, &im, &id, &fd)
	if status != 0 {
		return iy, im, id, 0, status
	}

	// Validate ndp
	if ndp < 0 || ndp > 9 {
		return iy, im, id, 0, -1
	}

	// Scale fraction to ndp decimal places
	// For ndp=4: 0.9999 * 10000 = 9999
	scale := 1.0
	for i := 0; i < ndp; i++ {
		scale *= 10.0
	}

	// Round to nearest integer
	fraction = int(fd*scale + 0.5)

	return iy, im, id, fraction, 0
}

// Caldj converts Gregorian calendar date to Modified Julian Date with 2-digit year handling
//
// Original FORTRAN: sla_CALDJ by P.T. Wallace
// GoFA equivalent: gofa.Cal2jd (with preprocessing)
//
// Parameters:
//   - iy: Year (4-digit or 2-digit)
//   - im: Month (1-12)
//   - id: Day (1-31)
//
// Returns:
//   - djm: Modified Julian Date
//   - status: 0 = OK, -1 = bad year, -2 = bad month, -3 = bad day
//
// Notes:
//   - 2-digit years are interpreted as:
//     00-49 → 2000-2049
//     50-99 → 1950-1999
//   - 4-digit years are used as-is
//   - This matches the PAL implementation (palCaldj.c)
func Caldj(iy, im, id int) (djm float64, status int) {
	// Handle 2-digit years
	if iy >= 0 && iy <= 49 {
		iy += 2000
	} else if iy >= 50 && iy <= 99 {
		iy += 1900
	}

	// Use standard calendar conversion
	return Cldj(iy, im, id)
}

// Calyd converts Gregorian calendar date to year and day number
//
// Original FORTRAN: sla_CALYD (not in standard SLALIB docs, but in some implementations)
//
// Parameters:
//   - iy: Year in Gregorian calendar
//   - im: Month (1-12)
//   - id: Day (1-31)
//
// Returns:
//   - ny: Year (same as input)
//   - nd: Day number in year (1-366)
//   - status: 0 = OK, non-zero = error from Cal2jd
//
// Example:
//   - (2024, 1, 1) → (2024, 1)   - January 1st is day 1
//   - (2024, 12, 31) → (2024, 366) - Leap year
//   - (2023, 12, 31) → (2023, 365) - Non-leap year
func Calyd(iy, im, id int) (ny, nd int, status int) {
	var djm0, djm float64

	// Convert given date to MJD
	status = gofa.Cal2jd(iy, im, id, &djm0, &djm)
	if status != 0 {
		return 0, 0, status
	}

	// Convert January 1 of same year to MJD
	var djm0_jan1, djm_jan1 float64
	status = gofa.Cal2jd(iy, 1, 1, &djm0_jan1, &djm_jan1)
	if status != 0 {
		return 0, 0, status
	}

	// Day number = difference + 1 (January 1 is day 1)
	nd = int(djm-djm_jan1) + 1
	ny = iy

	return ny, nd, 0
}

// Clyd converts year and day number to Gregorian calendar date
//
// Original FORTRAN: sla_CLYD (not in standard SLALIB docs)
//
// Parameters:
//   - iy: Year in Gregorian calendar
//   - id: Day number in year (1-366)
//
// Returns:
//   - ny: Year (same as input)
//   - nm: Month (1-12)
//   - nd: Day (1-31)
//   - status: 0 = OK, non-zero = error
//
// Example:
//   - (2024, 1) → (2024, 1, 1)   - Day 1 is January 1st
//   - (2024, 366) → (2024, 12, 31) - Day 366 in leap year
//   - (2023, 365) → (2023, 12, 31) - Day 365 in non-leap year
func Clyd(iy, id int) (ny, nm, nd int, status int) {
	var djm0, djm_jan1 float64

	// Get MJD of January 1
	status = gofa.Cal2jd(iy, 1, 1, &djm0, &djm_jan1)
	if status != 0 {
		return 0, 0, 0, status
	}

	// Add day offset (id-1 because January 1 is day 1, not day 0)
	djm := djm_jan1 + float64(id-1)

	// Convert back to calendar date
	var fd float64
	status = gofa.Jd2cal(djm0, djm, &ny, &nm, &nd, &fd)
	if status != 0 {
		return 0, 0, 0, status
	}

	return ny, nm, nd, 0
}

// Epb converts Julian Date to Besselian epoch
//
// Original FORTRAN: sla_EPB by P.T. Wallace
// GoFA equivalent: gofa.Epb
//
// Parameters:
//   - dj: Julian Date (TDB)
//
// Returns:
//   - Besselian epoch
//
// Notes:
//   - Besselian epoch is the epoch of the FK4 system
//   - B1900.0 = JD 2415020.31352
//   - B1950.0 = JD 2433282.42345905
//   - Tropical year = 365.242198781 days
func Epb(dj float64) float64 {
	// GoFA expects two-part JD for precision
	// We split MJD into epoch (2400000.5) + offset
	return gofa.Epb(gofa.DJM0, dj)
}

// Epb2d converts Besselian epoch to Julian Date
//
// Original FORTRAN: sla_EPB2D by P.T. Wallace
// GoFA equivalent: gofa.Epb2jd
//
// Parameters:
//   - epb: Besselian epoch
//
// Returns:
//   - Julian Date (TDB)
//
// Notes:
//   - Inverse of Epb
//   - Returns Modified Julian Date for SLALIB compatibility
func Epb2d(epb float64) float64 {
	var djm0, djm float64
	gofa.Epb2jd(epb, &djm0, &djm)
	// Return MJD part only
	return djm
}

// Epj converts Julian Date to Julian epoch
//
// Original FORTRAN: sla_EPJ by P.T. Wallace
// GoFA equivalent: gofa.Epj
//
// Parameters:
//   - dj: Julian Date (TT)
//
// Returns:
//   - Julian epoch
//
// Notes:
//   - Julian epoch is the standard epoch for modern astronomy
//   - J2000.0 = JD 2451545.0 (2000 January 1.5 TT)
//   - Julian year = exactly 365.25 days
func Epj(dj float64) float64 {
	// GoFA expects two-part JD for precision
	return gofa.Epj(gofa.DJM0, dj)
}

// Epj2d converts Julian epoch to Julian Date
//
// Original FORTRAN: sla_EPJ2D by P.T. Wallace
// GoFA equivalent: gofa.Epj2jd
//
// Parameters:
//   - epj: Julian epoch
//
// Returns:
//   - Julian Date (TT)
//
// Notes:
//   - Inverse of Epj
//   - Returns Modified Julian Date for SLALIB compatibility
func Epj2d(epj float64) float64 {
	var djm0, djm float64
	gofa.Epj2jd(epj, &djm0, &djm)
	// Return MJD part only
	return djm
}

// SLALIB-compatible lowercase aliases

// caldj is a SLALIB-compatible alias for Caldj (sla_CALDJ)
var caldj = Caldj

// cldj is a SLALIB-compatible alias for Cldj (sla_CLDJ)
var cldj = Cldj

// djcl is a SLALIB-compatible alias for Djcl (sla_DJCL)
var djcl = Djcl

// djcal is a SLALIB-compatible alias for Djcal (sla_DJCAL)
var djcal = Djcal

// calyd is a SLALIB-compatible alias for Calyd (sla_CALYD)
var calyd = Calyd

// clyd is a SLALIB-compatible alias for Clyd (sla_CLYD)
var clyd = Clyd

// epb is a SLALIB-compatible alias for Epb (sla_EPB)
var epb = Epb

// epb2d is a SLALIB-compatible alias for Epb2d (sla_EPB2D)
var epb2d = Epb2d

// epj is a SLALIB-compatible alias for Epj (sla_EPJ)
var epj = Epj

// epj2d is a SLALIB-compatible alias for Epj2d (sla_EPJ2D)
var epj2d = Epj2d
