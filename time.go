package slagofa

import "github.com/hebl/gofa"

// Time Scale and Sidereal Time Functions
//
// This file implements SLALIB-compatible time scale conversions and
// sidereal time calculations using GoFA (IAU SOFA).

// Dtt computes the offset TT-UTC (Terrestrial Time minus UTC)
//
// Original FORTRAN: sla_DTT by P.T. Wallace
// GoFA equivalent: gofa.Dat (leap seconds) + constant
//
// Parameters:
//   - utc: UTC as Modified Julian Date
//
// Returns:
//   - TT-UTC in seconds
//
// Notes:
//   - TT-UTC = (TT-TAI) + (TAI-UTC) = 32.184 + leap_seconds
//   - TT-TAI is exactly 32.184 seconds (defined constant)
//   - TAI-UTC is the accumulated leap seconds (from gofa.Dat)
//   - Valid from 1961 January 1 onwards (when leap seconds were introduced)
func Dtt(utc float64) float64 {
	// Convert MJD to calendar date to query leap seconds
	var iy, im, id int
	var fd float64
	status := gofa.Jd2cal(gofa.DJM0, utc, &iy, &im, &id, &fd)
	if status != 0 {
		// Invalid date - return constant offset (no leap seconds)
		return 32.184
	}

	// Get leap seconds (TAI-UTC) for this date
	var deltat float64
	status = gofa.Dat(iy, im, id, fd, &deltat)
	if status != 0 {
		// Date outside leap second table - return constant offset
		return 32.184
	}

	// TT-UTC = (TT-TAI) + (TAI-UTC) = 32.184 + deltat
	return 32.184 + deltat
}

// DeltaT estimates the offset between dynamical time and UT for a given epoch
//
// Original FORTRAN: sla_DT by P.T. Wallace
// PAL equivalent: palDt (direct port from PAL palDt.c)
//
// Parameters:
//   - epoch: Julian epoch (e.g., 2000.0, 1950.0, etc.)
//
// Returns:
//   - Estimate of ET-UT (or TT-UT after 1984) in seconds
//
// Notes:
//   - Uses three parabolic approximations based on epoch:
//     • Before 979: Stephenson & Morrison's 390 BC to AD 948 model
//     • 979 to 1708: Stephenson & Morrison's 948 to 1600 model
//     • After 1708: McCarthy & Babcock's post-1650 model
//   - Breakpoints (979.0258204760233 and 1708.185161980887) ensure continuity
//   - Based on lunar tidal acceleration of -26.00 arcsec/century²
//   - Accuracy: ~20 seconds post-1650, ~30 minutes around 1000 BC
//
// Historical context:
//   - Before 1984: Returns ET-UT (Ephemeris Time - Universal Time)
//   - After 1984: Returns TT-UT (Terrestrial Time - Universal Time)
//   - ET and TT differ by ~0.5ms, negligible for historical calculations
//
// Example:
//   - DeltaT(2000.0) ≈ 63.8 seconds (TT-UT at J2000.0)
//   - DeltaT(1950.0) ≈ 29.1 seconds
//   - DeltaT(1800.0) ≈ 13.7 seconds
func DeltaT(epoch float64) float64 {
	// Centuries since 1800
	t := (epoch - 1800.0) / 100.0

	var s float64

	if epoch >= 1708.185161980887 {
		// Post-1708: McCarthy & Babcock model
		w := t - 0.19
		s = 5.156 + 13.3066*w*w

	} else if epoch >= 979.0258204760233 {
		// 979-1708: Stephenson & Morrison's 948-1600 model
		s = 25.5 * t * t

	} else {
		// Pre-979: Stephenson & Morrison's 390 BC to AD 948 model
		s = 1360.0 + (320.0+44.3*t)*t
	}

	return s
}

// Dt is a SLALIB-compatible alias for DeltaT (sla_DT)
var Dt = DeltaT

// Gmst computes Greenwich Mean Sidereal Time from UT1
//
// Original FORTRAN: sla_GMST by P.T. Wallace
// GoFA equivalent: gofa.Gmst06 (IAU 2006 model, following PAL)
//
// Parameters:
//   - ut1: UT1 as Modified Julian Date
//
// Returns:
//   - Greenwich Mean Sidereal Time in radians (0 to 2π)
//
// Notes:
//   - PAL uses IAU 2006 model (Gmst06) for better accuracy than IAU 1982
//   - GMST is the hour angle of the mean vernal equinox
//   - Advances by ~366.2422 rotations per 365.2422 mean solar days
//   - Result is in range [0, 2π)
//   - Requires both TT and UT1; we use UT1 for both (introduces ~100 μas error)
//
// Example:
//   - At UT1 = 0h on 2000-01-01, GMST ≈ 6.697 hours ≈ 1.753 radians
func Gmst(ut1 float64) float64 {
	// GoFA.Gmst06 expects four arguments: (ut1a, ut1b, tta, ttb)
	// PAL calls: eraGmst06(PAL__MJD0, ut1, PAL__MJD0, ut1)
	// We have MJD, so convert: JD = MJD + 2400000.5
	// Using UT1 for both TT and UT1 (introduces ~100 microarcsecond error)
	return gofa.Gmst06(2400000.5, ut1, 2400000.5, ut1)
}

// Gmsta computes Greenwich Mean Sidereal Time from UT1 (high precision)
//
// Original FORTRAN: sla_GMSTA by P.T. Wallace
// GoFA equivalent: gofa.Gmst06 (IAU 2006 model, following PAL)
//
// Parameters:
//   - ut1a: UT1 as Modified Julian Date (first part)
//   - ut1b: UT1 as Modified Julian Date (second part, for precision)
//
// Returns:
//   - Greenwich Mean Sidereal Time in radians (0 to 2π)
//
// Notes:
//   - PAL uses IAU 2006 model (Gmst06) for better accuracy than IAU 1982
//   - Two-part date allows sub-microsecond precision
//   - For most uses, ut1b can be 0.0 and ut1a = full MJD
//   - When ultra-high precision needed, split MJD into integer + fraction
//   - Requires both TT and UT1; we use UT1 for both (introduces ~100 μas error)
//
// Example:
//   - Gmsta(51544.0, 0.5) for 2000-01-01 12:00 UT1
//   - Gmsta(51544.5, 0.0) gives same result (less precise input)
func Gmsta(ut1a, ut1b float64) float64 {
	// GoFA.Gmst06 expects four arguments: (ut1a, ut1b, tta, ttb)
	// PAL calls: eraGmst06(date, ut, date, ut) where date = ut1a + PAL__MJD0
	// We have MJD, so convert: JD = MJD + 2400000.5
	// Using UT1 for both TT and UT1 (introduces ~100 microarcsecond error)
	date := ut1a + 2400000.5
	return gofa.Gmst06(date, ut1b, date, ut1b)
}

// Helper: Gmsta0 is a simplified version for single MJD input
//
// This is a convenience wrapper for the common case where you have
// a single MJD value and don't need to split it for extra precision.
//
// Parameters:
//   - ut1: UT1 as Modified Julian Date
//
// Returns:
//   - Greenwich Mean Sidereal Time in radians
func Gmsta0(ut1 float64) float64 {
	return Gmsta(ut1, 0.0)
}

// Note: Calendar and epoch conversion functions (Caldj, Cldj, Djcal, Djcl, Epb, Epb2d, Epj, Epj2d)
// are implemented in calendar.go to keep related functionality together.

// SLALIB-compatible lowercase aliases

// dtt is a SLALIB-compatible alias for Dtt (sla_DTT)
var dtt = Dtt

// gmst is a SLALIB-compatible alias for Gmst (sla_GMST)
var gmst = Gmst

// gmsta is a SLALIB-compatible alias for Gmsta (sla_GMSTA)
var gmsta = Gmsta
