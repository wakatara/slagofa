package slagofa

import (
	"github.com/hebl/gofa"
)

// Helper function to convert MJD to two-part JD for GoFA
// Returns jd1, jd2 where JD = jd1 + jd2
func mjdToJD(mjd float64) (jd1, jd2 float64) {
	// MJD = JD - 2400000.5
	// So JD = MJD + 2400000.5
	// Split into two parts for precision: 2400000.5 + MJD
	return ModifiedJulianDateOffset, mjd
}

// Astrometry Functions
//
// This file implements SLALIB-compatible astrometry functions using GoFA.
// These functions handle proper motion, precession, nutation, and aberration.

// ProperMotion applies proper motion to star coordinates.
//
// Updates star position from epoch EP0 to EP1, accounting for proper motion,
// parallax, and radial velocity.
//
// Original FORTRAN: sla_PM by P.T. Wallace
// GoFA equivalent: gofa.Pmsafe (safe proper motion with validity checking)
// SLALIB reference: SUN/67 section 73
//
// Parameters:
//   - ra0: Right ascension at epoch EP0 (radians)
//   - dec0: Declination at epoch EP0 (radians)
//   - pmRA: Proper motion in RA (radians/year, including cos(Dec) factor)
//   - pmDec: Proper motion in Dec (radians/year)
//   - parallax: Parallax (arcseconds)
//   - rv: Radial velocity (km/s, positive if receding)
//   - ep0: Start epoch (Julian years, e.g., 2000.0)
//   - ep1: End epoch (Julian years)
//
// Returns:
//   - ra1: Right ascension at epoch EP1 (radians)
//   - dec1: Declination at epoch EP1 (radians)
//   - status: 0=OK, -1=system error, 1=distance overridden, 2=excessive velocity
//
// Notes:
//   - Proper motions are dRA/dt (not cos(Dec)*dRA/dt) in same system as input
//   - For pre-FK5 data, use Besselian epochs and scale RV by 365.2422/365.25
//   - Uses Pmsafe which prevents arithmetic problems with extreme values
//   - Accounts for perspective effects (parallax and radial velocity)
func ProperMotion(ra0, dec0, pmRA, pmDec, parallax, rv, ep0, ep1 float64) (ra1, dec1 float64, status int) {
	// Convert epochs to Julian Date
	// Use J2000.0 as reference: JD = 2451545.0 + (epoch - 2000.0) * 365.25
	jd0a := 2451545.0 + (ep0-2000.0)*365.25
	jd0b := 0.0
	jd1a := 2451545.0 + (ep1-2000.0)*365.25
	jd1b := 0.0

	// Dummy variables for proper motion, parallax, RV at end epoch
	// (Pmsafe updates these but we ignore them for SLALIB compatibility)
	var pmRA2, pmDec2, px2, rv2 float64

	status = gofa.Pmsafe(ra0, dec0, pmRA, pmDec, parallax, rv,
		jd0a, jd0b, jd1a, jd1b,
		&ra1, &dec1, &pmRA2, &pmDec2, &px2, &rv2)

	// Normalize RA to [0, 2π)
	ra1 = gofa.Anp(ra1)

	return ra1, dec1, status
}

// Pm is a SLALIB-compatible alias for ProperMotion (sla_PM)
var Pm = ProperMotion

// PrecessionMatrix computes the precession matrix from J2000.0 to a given date.
//
// Returns the precession matrix for transforming J2000.0 mean equator and
// equinox coordinates to the mean equator and equinox of a specified date.
//
// Original FORTRAN: sla_PREC by P.T. Wallace
// GoFA equivalent: gofa.Pmat76 (IAU 1976) or gofa.Pmat06 (IAU 2006)
// SLALIB reference: SUN/67 section 74
//
// Parameters:
//   - ep0: Beginning epoch (Julian years, e.g., 2000.0)
//   - ep1: Ending epoch (Julian years)
//
// Returns:
//   - pmat: Precession matrix [3][3] to rotate vectors from ep0 to ep1
//
// Notes:
//   - Uses IAU 2006 precession model (more accurate than SLALIB's IAU 1976)
//   - To apply: v1 = pmat × v0
//   - For J2000.0 → date, use ep0=2000.0
//   - Matrix is orthogonal (transpose = inverse)
func PrecessionMatrix(ep0, ep1 float64) Mat3 {
	// Convert epochs to Julian Date
	jd0 := 2451545.0 + (ep0-2000.0)*365.25
	jd1 := 2451545.0 + (ep1-2000.0)*365.25

	var pmat [3][3]float64

	// Use IAU 2006 precession (more accurate than SLALIB's IAU 1976)
	// For J2000.0 to date
	if ep0 == 2000.0 {
		gofa.Pmat06(jd1, 0.0, &pmat)
	} else {
		// For arbitrary epochs, need to compute from J2000 to each epoch
		// then combine: P(ep0→ep1) = P(J2000→ep1) × P(J2000→ep0)^T
		var p0, p1 [3][3]float64
		gofa.Pmat06(jd0, 0.0, &p0)
		gofa.Pmat06(jd1, 0.0, &p1)

		// pmat = p1 × p0^T
		var p0t [3][3]float64
		gofa.Tr(p0, &p0t) // Transpose p0
		gofa.Rxr(p1, p0t, &pmat)
	}

	return Mat3(pmat)
}

// Prec is a SLALIB-compatible alias for PrecessionMatrix (sla_PREC)
var Prec = PrecessionMatrix

// PrecessionMatrix76 computes the IAU 1976 precession matrix (for compatibility).
//
// Like PrecessionMatrix but uses the older IAU 1976 model for exact
// SLALIB compatibility when needed.
//
// Original FORTRAN: sla_PREC by P.T. Wallace (using IAU 1976)
// GoFA equivalent: gofa.Pmat76
//
// Parameters:
//   - ep0: Beginning epoch (Julian years)
//   - ep1: Ending epoch (Julian years)
//
// Returns:
//   - pmat: IAU 1976 precession matrix
func PrecessionMatrix76(ep0, ep1 float64) Mat3 {
	jd0 := 2451545.0 + (ep0-2000.0)*365.25
	jd1 := 2451545.0 + (ep1-2000.0)*365.25

	var pmat [3][3]float64

	if ep0 == 2000.0 {
		gofa.Pmat76(jd1, 0.0, &pmat)
	} else {
		var p0, p1 [3][3]float64
		gofa.Pmat76(jd0, 0.0, &p0)
		gofa.Pmat76(jd1, 0.0, &p1)

		var p0t [3][3]float64
		gofa.Tr(p0, &p0t)
		gofa.Rxr(p1, p0t, &pmat)
	}

	return Mat3(pmat)
}

// NutationComponents computes nutation in longitude and obliquity.
//
// Returns the nutation components (nutation in longitude and obliquity)
// for a given date using IAU 2006/2000A nutation model.
//
// Original FORTRAN: sla_NUTC by P.T. Wallace
// GoFA equivalent: gofa.Nut06a (IAU 2006/2000A) or gofa.Nut80 (IAU 1980)
// SLALIB reference: SUN/67 section 68
//
// Parameters:
//   - mjd: Modified Julian Date (TT)
//
// Returns:
//   - dpsi: Nutation in longitude (radians)
//   - deps: Nutation in obliquity (radians)
//
// Notes:
//   - Uses IAU 2006/2000A nutation model (more accurate than SLALIB's IAU 1980)
//   - Date should be in TT (Terrestrial Time)
//   - For J2000.0: dpsi and deps are small (~0.0001 rad = 20 arcsec)
func NutationComponents(mjd float64) (dpsi, deps float64) {
	// Convert MJD to JD
	jd1, jd2 := mjdToJD(mjd)

	// Use IAU 2006/2000A nutation model
	gofa.Nut06a(jd1, jd2, &dpsi, &deps)

	return dpsi, deps
}

// Nutc is a SLALIB-compatible alias for NutationComponents (sla_NUTC)
var Nutc = NutationComponents

// NutationComponents80 computes IAU 1980 nutation (for compatibility).
//
// Like NutationComponents but uses the older IAU 1980 model for exact
// SLALIB compatibility when needed.
//
// Original FORTRAN: sla_NUTC by P.T. Wallace (using IAU 1980)
// GoFA equivalent: gofa.Nut80
//
// Parameters:
//   - mjd: Modified Julian Date (TT)
//
// Returns:
//   - dpsi: Nutation in longitude (radians, IAU 1980)
//   - deps: Nutation in obliquity (radians, IAU 1980)
func NutationComponents80(mjd float64) (dpsi, deps float64) {
	jd1, jd2 := mjdToJD(mjd)
	gofa.Nut80(jd1, jd2, &dpsi, &deps)
	return dpsi, deps
}

// NutationMatrix computes the nutation matrix.
//
// Returns the nutation matrix for transforming from mean to true equator
// and equinox of date.
//
// Original FORTRAN: sla_NUT by P.T. Wallace
// GoFA equivalent: gofa.Nutm80 (IAU 1980)
// SLALIB reference: SUN/67 section 67
//
// Parameters:
//   - mjd: Modified Julian Date (TT)
//
// Returns:
//   - nmat: Nutation matrix [3][3] (mean → true equator & equinox)
//
// Notes:
//   - Transforms from mean equator/equinox to true equator/equinox
//   - To apply: v_true = nmat × v_mean
//   - Uses IAU 1980 nutation for SLALIB compatibility
func NutationMatrix(mjd float64) Mat3 {
	jd1, jd2 := mjdToJD(mjd)

	var nmat [3][3]float64
	gofa.Nutm80(jd1, jd2, &nmat)

	return Mat3(nmat)
}

// Nut is a SLALIB-compatible alias for NutationMatrix (sla_NUT)
var Nut = NutationMatrix

// Precess applies precession to coordinates.
//
// Convenience function that applies precession to spherical coordinates.
//
// Original FORTRAN: sla_PRECES by P.T. Wallace
// GoFA equivalent: Combination of Pmat06 and vector operations
// SLALIB reference: SUN/67 section 75
//
// Parameters:
//   - ra0: Right ascension at epoch EP0 (radians)
//   - dec0: Declination at epoch EP0 (radians)
//   - ep0: Beginning epoch (Julian years)
//   - ep1: Ending epoch (Julian years)
//
// Returns:
//   - ra1: Right ascension at epoch EP1 (radians)
//   - dec1: Declination at epoch EP1 (radians)
//
// Notes:
//   - Applies precession only (not nutation)
//   - Uses IAU 2006 precession model
//   - For mean equator and equinox
func Precess(ra0, dec0, ep0, ep1 float64) (ra1, dec1 float64) {
	// Convert to Cartesian
	var v0 [3]float64
	gofa.S2c(ra0, dec0, &v0)

	// Get precession matrix
	pmat := PrecessionMatrix(ep0, ep1)

	// Apply precession
	var v1 [3]float64
	gofa.Rxp([3][3]float64(pmat), v0, &v1)

	// Convert back to spherical
	gofa.C2s(v1, &ra1, &dec1)

	// Normalize RA
	ra1 = gofa.Anp(ra1)

	return ra1, dec1
}

// Preces is a SLALIB-compatible alias for Precess (sla_PRECES)
var Preces = Precess

// PrecessionNutationMatrix computes combined precession-nutation matrix.
//
// Returns the combined matrix for transforming from J2000.0 mean equator
// and equinox to true equator and equinox of date.
//
// Original FORTRAN: sla_PRENUT by P.T. Wallace
// GoFA equivalent: Combination of Pmat06 and Nutm80
// SLALIB reference: SUN/67 section 76
//
// Parameters:
//   - ep0: Beginning epoch (Julian years, typically 2000.0)
//   - mjd: Modified Julian Date for end epoch (TT)
//
// Returns:
//   - pnmat: Combined precession-nutation matrix [3][3]
//
// Notes:
//   - Combines precession from ep0 to date, then nutation
//   - Result: PN = N × P (nutation after precession)
//   - To apply: v_true = pnmat × v_mean_J2000
func PrecessionNutationMatrix(ep0, mjd float64) Mat3 {
	jd1, jd2 := mjdToJD(mjd)

	// Get precession matrix from ep0 to date
	var pmat [3][3]float64
	if ep0 == 2000.0 {
		gofa.Pmat06(jd1, jd2, &pmat)
	} else {
		jd0 := 2451545.0 + (ep0-2000.0)*365.25
		var p0, p1 [3][3]float64
		gofa.Pmat06(jd0, 0.0, &p0)
		gofa.Pmat06(jd1, jd2, &p1)
		var p0t [3][3]float64
		gofa.Tr(p0, &p0t)
		gofa.Rxr(p1, p0t, &pmat)
	}

	// Get nutation matrix for date
	var nmat [3][3]float64
	gofa.Nutm80(jd1, jd2, &nmat)

	// Combine: PN = N × P
	var pnmat [3][3]float64
	gofa.Rxr(nmat, pmat, &pnmat)

	return Mat3(pnmat)
}

// Prenut is a SLALIB-compatible alias for PrecessionNutationMatrix (sla_PRENUT)
var Prenut = PrecessionNutationMatrix
