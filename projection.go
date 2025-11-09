package slagofa

import (
	"github.com/hebl/gofa"
)

// Phase 6: Tangent Plane Projections
//
// These functions handle projections between spherical coordinates and
// tangent plane coordinates, commonly used in astrometry and image processing.

// Ds2tp projects spherical coordinates to tangent plane (gnomonic projection).
//
// This function projects a point on the celestial sphere onto a tangent plane.
// The tangent point is specified by (ra0, dec0), and the point to project is (ra, dec).
//
// Original FORTRAN: sla_DS2TP by P.T. Wallace
// GoFA equivalent: gofa.Tpxes (Tangent plane: spherical to rectangular)
// SLALIB reference: SUN/67 section 41
//
// Parameters:
//   - ra: RA of point to project (radians)
//   - dec: Dec of point to project (radians)
//   - ra0: RA of tangent point (radians)
//   - dec0: Dec of tangent point (radians)
//
// Returns:
//   - xi: Tangent plane X coordinate (radians)
//   - eta: Tangent plane Y coordinate (radians)
//   - status: 0 = OK, 1 = star too far from axis, 2 = antipoint, 3 = bad dec
//
// Notes:
//   - Implements gnomonic projection (tangent plane)
//   - Status 0 = point projects OK
//   - Status 1 = point > 90° from tangent point
//   - Status 2 = point is antipoint (exactly opposite)
//   - Status 3 = invalid declination
//   - Uses modern SOFA algorithm (may differ slightly from SLALIB)
func Ds2tp(ra, dec, ra0, dec0 float64) (xi, eta float64, status int) {
	status = gofa.Tpxes(ra, dec, ra0, dec0, &xi, &eta)
	return xi, eta, status
}

// Dtp2s projects tangent plane coordinates to spherical (gnomonic deprojection).
//
// This is the inverse of Ds2tp, converting tangent plane coordinates back
// to spherical coordinates.
//
// Original FORTRAN: sla_DTP2S by P.T. Wallace
// GoFA equivalent: gofa.Tpsts (Tangent plane: rectangular to spherical)
// SLALIB reference: SUN/67 section 59
//
// Parameters:
//   - xi: Tangent plane X coordinate (radians)
//   - eta: Tangent plane Y coordinate (radians)
//   - ra0: RA of tangent point (radians)
//   - dec0: Dec of tangent point (radians)
//
// Returns:
//   - ra: RA of deprojected point (radians)
//   - dec: Dec of deprojected point (radians)
//
// Notes:
//   - Inverse of Ds2tp
//   - Always succeeds (no status return)
//   - Uses modern SOFA algorithm
func Dtp2s(xi, eta, ra0, dec0 float64) (ra, dec float64) {
	gofa.Tpsts(xi, eta, ra0, dec0, &ra, &dec)
	return ra, dec
}

// Dtps2c solves for spherical coordinates from tangent plane position.
//
// Given a tangent plane position (xi, eta) and a spherical position (ra, dec)
// that corresponds to it, this function finds the two possible tangent points.
//
// Original FORTRAN: sla_DTPS2C by P.T. Wallace
// GoFA equivalent: gofa.Tpors (Tangent plane: solve for origin, spherical)
// SLALIB reference: SUN/67 section 60
//
// Parameters:
//   - xi: Tangent plane X coordinate (radians)
//   - eta: Tangent plane Y coordinate (radians)
//   - ra: RA of point (radians)
//   - dec: Dec of point (radians)
//
// Returns:
//   - ra01: RA of solution 1 (radians)
//   - dec01: Dec of solution 1 (radians)
//   - ra02: RA of solution 2 (radians)
//   - dec02: Dec of solution 2 (radians)
//   - n: Number of solutions: 0, 1, or 2
//
// Notes:
//   - Solves the inverse problem: given (xi, eta) and (ra, dec), find tangent point
//   - Two solutions usually exist (ambiguity in tangent point location)
//   - n=0: No solution (invalid input)
//   - n=1: One solution (point at pole or xi=eta=0)
//   - n=2: Two solutions (general case)
func Dtps2c(xi, eta, ra, dec float64) (ra01, dec01, ra02, dec02 float64, n int) {
	n = gofa.Tpors(xi, eta, ra, dec, &ra01, &dec01, &ra02, &dec02)
	return ra01, dec01, ra02, dec02, n
}

// Ds2tpv projects Cartesian unit vector to tangent plane.
//
// This is the vector form of Ds2tp, taking a Cartesian unit vector
// instead of spherical coordinates.
//
// Original FORTRAN: sla_DV2TP by P.T. Wallace
// GoFA equivalent: gofa.Tpxev (Tangent plane: vector to rectangular)
// SLALIB reference: SUN/67 section 65
//
// Parameters:
//   - v: Direction vector (need not be unit)
//   - v0: Tangent point direction vector
//
// Returns:
//   - xi: Tangent plane X coordinate (radians)
//   - eta: Tangent plane Y coordinate (radians)
//   - status: 0 = OK, 1 = too far from axis, 2 = antipoint, 3 = null vector
//
// Notes:
//   - Vector version of Ds2tp
//   - Vectors need not be normalized
//   - Status codes same as Ds2tp
func Ds2tpv(v, v0 Vec3) (xi, eta float64, status int) {
	status = gofa.Tpxev(v, v0, &xi, &eta)
	return xi, eta, status
}

// Dtp2sv projects tangent plane coordinates to Cartesian unit vector.
//
// This is the vector form of Dtp2s, returning a Cartesian unit vector
// instead of spherical coordinates.
//
// Original FORTRAN: sla_DTP2V by P.T. Wallace
// GoFA equivalent: gofa.Tpstv (Tangent plane: rectangular to vector)
// SLALIB reference: SUN/67 section 61
//
// Parameters:
//   - xi: Tangent plane X coordinate (radians)
//   - eta: Tangent plane Y coordinate (radians)
//   - v0: Tangent point direction vector
//
// Returns:
//   - v: Direction vector
//
// Notes:
//   - Vector version of Dtp2s
//   - Inverse of Ds2tpv
//   - Always succeeds (no status return)
func Dtp2sv(xi, eta float64, v0 Vec3) Vec3 {
	var result [3]float64
	gofa.Tpstv(xi, eta, v0, &result)
	return Vec3(result)
}

// Dtpv2c solves for spherical coordinates from tangent plane position (vector form).
//
// This is the vector form of Dtps2c, taking and returning Cartesian vectors
// instead of spherical coordinates.
//
// Original FORTRAN: sla_DTPV2C by P.T. Wallace
// GoFA equivalent: gofa.Tporv (Tangent plane: solve for origin, vector)
// SLALIB reference: SUN/67 section 62
//
// Parameters:
//   - xi: Tangent plane X coordinate (radians)
//   - eta: Tangent plane Y coordinate (radians)
//   - v: Direction vector of point
//
// Returns:
//   - v01: Direction vector of solution 1
//   - v02: Direction vector of solution 2
//   - n: Number of solutions: 0, 1, or 2
//
// Notes:
//   - Vector version of Dtps2c
//   - Solves inverse problem: given (xi, eta) and v, find tangent point
//   - Returns direction vectors instead of RA/Dec
func Dtpv2c(xi, eta float64, v Vec3) (v01, v02 Vec3, n int) {
	var result1, result2 [3]float64
	n = gofa.Tporv(xi, eta, v, &result1, &result2)
	return Vec3(result1), Vec3(result2), n
}

// SLALIB-compatible lowercase aliases

// ds2tp is a SLALIB-compatible alias for Ds2tp (sla_DS2TP)
var ds2tp = Ds2tp

// dtp2s is a SLALIB-compatible alias for Dtp2s (sla_DTP2S)
var dtp2s = Dtp2s

// dtps2c is a SLALIB-compatible alias for Dtps2c (sla_DTPS2C)
var dtps2c = Dtps2c

// dtpv2c is a SLALIB-compatible alias for Dtpv2c (sla_DTPV2C)
var dtpv2c = Dtpv2c
