package slagofa

import (
	"github.com/hebl/gofa"
)

// Angle Operations
// These functions wrap GoFA's angle library to provide SLALIB-compatible
// API for angle normalization, formatting, and angular separation.

// NormalizeAngle normalizes an angle into the range -π to +π.
//
// Original FORTRAN: sla_DRANGE by P.T. Wallace / Rutherford Appleton Laboratory
// GoFA equivalent: gofa.Anpm (normalize angle into range -π to +π)
//
// Parameters:
//   - angle: Angle in radians
//
// Returns:
//   - Angle normalized to the range -π to +π
func NormalizeAngle(angle float64) float64 {
	return gofa.Anpm(angle)
}

// Drange is a SLALIB-compatible alias for NormalizeAngle (sla_DRANGE)
var Drange = NormalizeAngle

// NormalizeAnglePositive normalizes an angle into the range 0 to 2π.
//
// Original FORTRAN: sla_DRANRM by P.T. Wallace
// GoFA equivalent: gofa.Anp (normalize angle into range 0 to 2π)
//
// Parameters:
//   - angle: Angle in radians
//
// Returns:
//   - Angle normalized to the range 0 to 2π
func NormalizeAnglePositive(angle float64) float64 {
	return gofa.Anp(angle)
}

// Dranrm is a SLALIB-compatible alias for NormalizeAnglePositive (sla_DRANRM)
var Dranrm = NormalizeAnglePositive

// RadiansToAngle decomposes radians into degrees, arcminutes, arcseconds, fraction.
//
// Original FORTRAN: sla_DR2AF by P.T. Wallace
// GoFA equivalent: gofa.A2af (decompose radians to degrees/arcmin/arcsec)
//
// Parameters:
//   - ndp: Resolution (number of decimal places in arcseconds)
//   - angle: Angle in radians
//
// Returns:
//   - sign: '+' or '-'
//   - degrees: Degrees (0-359)
//   - arcmin: Arcminutes (0-59)
//   - arcsec: Arcseconds (0-59)
//   - fraction: Fractional arcseconds (scaled by 10^ndp)
func RadiansToAngle(ndp int, angle float64) (sign byte, degrees, arcmin, arcsec, fraction int) {
	var idmsf [4]int
	gofa.A2af(ndp, angle, &sign, &idmsf)
	return sign, idmsf[0], idmsf[1], idmsf[2], idmsf[3]
}

// Dr2af is a SLALIB-compatible wrapper for RadiansToAngle (sla_DR2AF)
func Dr2af(ndp int, angle float64) (byte, [4]int) {
	var sign byte
	var idmsf [4]int
	gofa.A2af(ndp, angle, &sign, &idmsf)
	return sign, idmsf
}

// RadiansToTime decomposes radians into hours, minutes, seconds, fraction.
//
// Original FORTRAN: sla_DR2TF by P.T. Wallace
// GoFA equivalent: gofa.A2tf (decompose radians to hours/min/sec)
//
// Parameters:
//   - ndp: Resolution (number of decimal places in seconds)
//   - angle: Angle in radians
//
// Returns:
//   - sign: '+' or '-'
//   - hours: Hours (0-23)
//   - minutes: Minutes (0-59)
//   - seconds: Seconds (0-59)
//   - fraction: Fractional seconds (scaled by 10^ndp)
func RadiansToTime(ndp int, angle float64) (sign byte, hours, minutes, seconds, fraction int) {
	var ihmsf [4]int
	gofa.A2tf(ndp, angle, &sign, &ihmsf)
	return sign, ihmsf[0], ihmsf[1], ihmsf[2], ihmsf[3]
}

// Dr2tf is a SLALIB-compatible wrapper for RadiansToTime (sla_DR2TF)
func Dr2tf(ndp int, angle float64) (byte, [4]int) {
	var sign byte
	var ihmsf [4]int
	gofa.A2tf(ndp, angle, &sign, &ihmsf)
	return sign, ihmsf
}

// AngleToRadians converts degrees, arcminutes, arcseconds to radians.
//
// Original FORTRAN: sla_DAF2R by P.T. Wallace
// GoFA equivalent: gofa.Af2a (convert deg/arcmin/arcsec to radians)
//
// Parameters:
//   - sign: '+' or '-' (negative if '-')
//   - degrees: Degrees
//   - arcmin: Arcminutes
//   - arcsec: Arcseconds
//
// Returns:
//   - radians: Angle in radians
//   - status: 0 = OK, 1 = bad degrees, 2 = bad arcmin, 3 = bad arcsec
func AngleToRadians(sign byte, degrees, arcmin int, arcsec float64) (radians float64, status int) {
	status = gofa.Af2a(sign, degrees, arcmin, arcsec, &radians)
	return radians, status
}

// Daf2r is a SLALIB-compatible wrapper for AngleToRadians (sla_DAF2R)
func Daf2r(sign byte, degrees, arcmin int, arcsec float64) (float64, int) {
	return AngleToRadians(sign, degrees, arcmin, arcsec)
}

// TimeToRadians converts hours, minutes, seconds to radians.
//
// Original FORTRAN: sla_DTF2R by P.T. Wallace
// GoFA equivalent: gofa.Tf2a (convert hours/min/sec to radians)
//
// Parameters:
//   - sign: '+' or '-' (negative if '-')
//   - hours: Hours
//   - minutes: Minutes
//   - seconds: Seconds
//
// Returns:
//   - radians: Angle in radians
//   - status: 0 = OK, 1 = bad hours, 2 = bad minutes, 3 = bad seconds
func TimeToRadians(sign byte, hours, minutes int, seconds float64) (radians float64, status int) {
	status = gofa.Tf2a(sign, hours, minutes, seconds, &radians)
	return radians, status
}

// Dtf2r is a SLALIB-compatible wrapper for TimeToRadians (sla_DTF2R)
func Dtf2r(sign byte, hours, minutes int, seconds float64) (float64, int) {
	return TimeToRadians(sign, hours, minutes, seconds)
}

// TimeToDays converts hours, minutes, seconds to days.
//
// Original FORTRAN: sla_DTF2D by P.T. Wallace
// GoFA equivalent: gofa.Tf2d (convert hours/min/sec to days)
//
// Parameters:
//   - sign: '+' or '-' (negative if '-')
//   - hours: Hours
//   - minutes: Minutes
//   - seconds: Seconds
//
// Returns:
//   - days: Interval in days
//   - status: 0 = OK, 1 = bad hours, 2 = bad minutes, 3 = bad seconds
func TimeToDays(sign byte, hours, minutes int, seconds float64) (days float64, status int) {
	status = gofa.Tf2d(sign, hours, minutes, seconds, &days)
	return days, status
}

// Dtf2d is a SLALIB-compatible wrapper for TimeToDays (sla_DTF2D)
func Dtf2d(sign byte, hours, minutes int, seconds float64) (float64, int) {
	return TimeToDays(sign, hours, minutes, seconds)
}

// DaysToTime decomposes days into hours, minutes, seconds, fraction.
//
// Original FORTRAN: sla_DD2TF by P.T. Wallace
// GoFA equivalent: gofa.D2tf (decompose days to hours/min/sec)
//
// Parameters:
//   - ndp: Resolution (number of decimal places in seconds)
//   - days: Interval in days
//
// Returns:
//   - sign: '+' or '-'
//   - hours: Hours
//   - minutes: Minutes
//   - seconds: Seconds
//   - fraction: Fractional seconds (scaled by 10^ndp)
func DaysToTime(ndp int, days float64) (sign byte, hours, minutes, seconds, fraction int) {
	var ihmsf [4]int
	gofa.D2tf(ndp, days, &sign, &ihmsf)
	return sign, ihmsf[0], ihmsf[1], ihmsf[2], ihmsf[3]
}

// Dd2tf is a SLALIB-compatible wrapper for DaysToTime (sla_DD2TF)
func Dd2tf(ndp int, days float64) (byte, [4]int) {
	var sign byte
	var ihmsf [4]int
	gofa.D2tf(ndp, days, &sign, &ihmsf)
	return sign, ihmsf
}

// Angular Separation Functions

// AngularSeparation computes the angular separation between two points
// given in spherical coordinates (e.g., RA/Dec).
//
// Original FORTRAN: sla_DSEP by P.T. Wallace
// GoFA equivalent: gofa.Seps (angular separation from spherical coordinates)
//
// Parameters:
//   - a1: Longitude of first point (e.g., RA) in radians
//   - b1: Latitude of first point (e.g., Dec) in radians
//   - a2: Longitude of second point in radians
//   - b2: Latitude of second point in radians
//
// Returns:
//   - Angular separation in radians (always positive)
func AngularSeparation(a1, b1, a2, b2 float64) float64 {
	return gofa.Seps(a1, b1, a2, b2)
}

// Dsep is a SLALIB-compatible alias for AngularSeparation (sla_DSEP)
var Dsep = AngularSeparation

// AngularSeparationVec computes the angular separation between two
// direction vectors.
//
// Original FORTRAN: sla_DSEPV by P.T. Wallace
// GoFA equivalent: gofa.Sepp (angular separation between two p-vectors)
//
// Parameters:
//   - a: First direction vector (need not be unit length)
//   - b: Second direction vector (need not be unit length)
//
// Returns:
//   - Angular separation in radians (always positive)
func AngularSeparationVec(a, b Vec3) float64 {
	return gofa.Sepp(a, b)
}

// Dsepv is a SLALIB-compatible alias for AngularSeparationVec (sla_DSEPV)
var Dsepv = AngularSeparationVec

// AngularSeparation32 computes the angular separation between two points
// on a sphere (single precision).
//
// Original FORTRAN: sla_SEP by P.T. Wallace
// Implementation: Converts float32 to float64, calls AngularSeparation
//
// Parameters:
//   - a1: Longitude of first point (e.g., RA) in radians
//   - b1: Latitude of first point (e.g., Dec) in radians
//   - a2: Longitude of second point in radians
//   - b2: Latitude of second point in radians
//
// Returns:
//   - Angular separation in radians (always positive)
func AngularSeparation32(a1, b1, a2, b2 float32) float32 {
	return float32(gofa.Seps(float64(a1), float64(b1), float64(a2), float64(b2)))
}

// Sep is a SLALIB-compatible alias for AngularSeparation32 (sla_SEP)
var Sep = AngularSeparation32

// AngularSeparationVec32 computes the angular separation between two
// direction vectors (single precision).
//
// Original FORTRAN: sla_SEPV by P.T. Wallace
// Implementation: Converts Vec3_32 to Vec3, calls AngularSeparationVec
//
// Parameters:
//   - a: First direction vector (need not be unit length)
//   - b: Second direction vector (need not be unit length)
//
// Returns:
//   - Angular separation in radians (always positive)
func AngularSeparationVec32(a, b Vec3_32) float32 {
	a64 := Vec3{float64(a[0]), float64(a[1]), float64(a[2])}
	b64 := Vec3{float64(b[0]), float64(b[1]), float64(b[2])}
	return float32(gofa.Sepp(a64, b64))
}

// Sepv is a SLALIB-compatible alias for AngularSeparationVec32 (sla_SEPV)
var Sepv = AngularSeparationVec32

// PositionAngle computes the position angle of one point with respect
// to another, given as direction vectors.
//
// Original FORTRAN: sla_DPAV by P.T. Wallace
// GoFA equivalent: gofa.Pap (position angle from two p-vectors)
//
// Parameters:
//   - a: Direction of reference point
//   - b: Direction of point whose PA is required
//
// Returns:
//   - Position angle in radians (-π to +π)
//
// Notes:
//   - If b is north of a, PA ≈ 0
//   - If b is east of a, PA ≈ +π/2
func PositionAngle(a, b Vec3) float64 {
	return gofa.Pap(a, b)
}

// Dpav is a SLALIB-compatible alias for PositionAngle (sla_DPAV)
var Dpav = PositionAngle

// Bearing computes the bearing (position angle) of one point on a sphere
// relative to another, given in spherical coordinates.
//
// Original FORTRAN: sla_DBEAR by P.T. Wallace
// GoFA equivalent: gofa.Pas (position angle from spherical coordinates)
//
// Parameters:
//   - a1: Longitude of point 1 (e.g., RA) in radians
//   - b1: Latitude of point 1 (e.g., Dec) in radians
//   - a2: Longitude of point 2 in radians
//   - b2: Latitude of point 2 in radians
//
// Returns:
//   - Bearing in radians (-π to +π)
//   - If point 2 is east of point 1, bearing ≈ +π/2
func Bearing(a1, b1, a2, b2 float64) float64 {
	return gofa.Pas(a1, b1, a2, b2)
}

// Dbear is a SLALIB-compatible alias for Bearing (sla_DBEAR)
var Dbear = Bearing

// Single-precision angle normalization

// NormalizeAngle32 normalizes an angle into the range -π to +π (single precision).
//
// Original FORTRAN: sla_RANGE by P.T. Wallace
//
// Parameters:
//   - angle: Angle in radians
//
// Returns:
//   - Angle normalized to the range -π to +π
func NormalizeAngle32(angle float32) float32 {
	return float32(gofa.Anpm(float64(angle)))
}

// Range is a SLALIB-compatible alias for NormalizeAngle32 (sla_RANGE)
var Range = NormalizeAngle32

// NormalizeAnglePositive32 normalizes an angle into the range 0 to 2π (single precision).
//
// Original FORTRAN: sla_RANORM by P.T. Wallace
//
// Parameters:
//   - angle: Angle in radians
//
// Returns:
//   - Angle normalized to the range 0 to 2π
func NormalizeAnglePositive32(angle float32) float32 {
	return float32(gofa.Anp(float64(angle)))
}

// Ranorm is a SLALIB-compatible alias for NormalizeAnglePositive32 (sla_RANORM)
var Ranorm = NormalizeAnglePositive32

// Additional SLALIB-compatible lowercase aliases

// daf2r is a SLALIB-compatible alias for Daf2r (sla_DAF2R)
var daf2r = Daf2r

// dd2tf is a SLALIB-compatible alias for Dd2tf (sla_DD2TF)
var dd2tf = Dd2tf

// dr2af is a SLALIB-compatible alias for Dr2af (sla_DR2AF)
var dr2af = Dr2af

// dr2tf is a SLALIB-compatible alias for Dr2tf (sla_DR2TF)
var dr2tf = Dr2tf

// dtf2d is a SLALIB-compatible alias for Dtf2d (sla_DTF2D)
var dtf2d = Dtf2d

// dtf2r is a SLALIB-compatible alias for Dtf2r (sla_DTF2R)
var dtf2r = Dtf2r
