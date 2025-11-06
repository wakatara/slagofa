// Package slagofa provides a SLALIB-compatible API layer on top of GoFA
// (Golang Standards Of Fundamental Astronomy).
//
// This package follows the same approach as PAL (Positional Astronomy Library),
// providing familiar SLALIB function names while using modern IAU standards
// via GoFA/SOFA underneath.
//
// Each function is available in two forms:
//   - Go-idiomatic PascalCase names (e.g., NormalizeAngle, DotProduct)
//   - SLALIB-compatible lowercase aliases (e.g., Drange, Dvdv)
//
// Example usage:
//
//	// Go-idiomatic API
//	angle := slagofa.NormalizeAngle(-4.0)
//	dot := slagofa.DotProduct(v1, v2)
//
//	// SLALIB-compatible API
//	angle := slagofa.Drange(-4.0)  // sla_DRANGE
//	dot := slagofa.Dvdv(v1, v2)    // sla_DVDV
package slagofa

import (
	"math"
)

// Mathematical constants matching SLALIB/SOFA values
const (
	// Pi - Ratio of circumference to diameter
	Pi = 3.141592653589793238462643

	// TwoPi - 2 × π
	TwoPi = 6.283185307179586476925287

	// HalfPi - π/2
	HalfPi = 1.570796326794896619231322

	// DegreesToRadians - Conversion factor: degrees to radians
	DegreesToRadians = 0.017453292519943295769236907684886

	// RadiansToDegrees - Conversion factor: radians to degrees
	RadiansToDegrees = 57.295779513082320876798154814105

	// HoursToRadians - Conversion factor: hours to radians
	HoursToRadians = 0.26179938779914943653855361527329

	// RadiansToHours - Conversion factor: radians to hours
	RadiansToHours = 3.8197186342054880584532103209403

	// ArcsecondsToRadians - Conversion factor: arcseconds to radians
	ArcsecondsToRadians = 4.8481368110953599358991410235795e-6

	// RadiansToArcseconds - Conversion factor: radians to arcseconds
	RadiansToArcseconds = 206264.80624709635515647335733078

	// SecondsPerDay - Seconds in a day
	SecondsPerDay = 86400.0

	// DaysPerJulianYear - Days in a Julian year
	DaysPerJulianYear = 365.25

	// DaysPerTropicalYear - Days in a tropical year
	DaysPerTropicalYear = 365.242198781

	// AstronomicalUnit - Astronomical unit in meters (IAU 2012)
	AstronomicalUnit = 149597870700.0

	// SpeedOfLight - Speed of light in m/s (defining constant)
	SpeedOfLight = 299792458.0

	// JulianEpoch2000 - Julian epoch J2000.0 (JD)
	JulianEpoch2000 = 2451545.0

	// ModifiedJulianDateOffset - Offset between JD and MJD
	ModifiedJulianDateOffset = 2400000.5
)

// Vec3 represents a 3-dimensional vector with double precision (float64) components.
// This is used for positions, velocities, and direction cosines in astronomical calculations.
//
// In SLALIB, this corresponds to the 3-element DOUBLE PRECISION array used in
// functions like sla_DVDV, sla_DVN, sla_DVXV, etc.
type Vec3 [3]float64

// Vec3_32 represents a 3-dimensional vector with single precision (float32) components.
// This is the single-precision counterpart to Vec3.
//
// In SLALIB, this corresponds to the 3-element REAL array used in
// functions like sla_VDV, sla_VN, sla_VXV, etc.
type Vec3_32 [3]float32

// Mat3 represents a 3×3 rotation matrix with double precision components.
// Matrices are stored in row-major order: mat[row][col].
//
// In SLALIB, this corresponds to the 3×3 DOUBLE PRECISION array used in
// functions like sla_DMXM, sla_DMXV, sla_DEULER, etc.
type Mat3 [3][3]float64

// Mat3_32 represents a 3×3 rotation matrix with single precision components.
// This is the single-precision counterpart to Mat3.
//
// In SLALIB, this corresponds to the 3×3 REAL array used in
// functions like sla_MXM, sla_MXV, sla_EULER, etc.
type Mat3_32 [3][3]float32

// PosVel represents a position/velocity vector pair with double precision.
// The first element [0] is position, the second [1] is velocity.
//
// In SOFA/GoFA, this is called a "pv-vector" and is used in functions
// dealing with motion (position + velocity).
type PosVel [2][3]float64

// SphericalCoord represents spherical coordinates (longitude, latitude).
// Both components are in radians.
//
// Longitude (θ): typically right ascension, hour angle, or azimuth
// Latitude (φ): typically declination, altitude, or elevation
type SphericalCoord struct {
	Longitude float64 // θ (theta) - longitude in radians
	Latitude  float64 // φ (phi) - latitude in radians
}

// DMS represents an angle in degrees, arcminutes, arcseconds format.
// Used for output formatting of angular values.
type DMS struct {
	Sign       byte // '+' or '-'
	Degrees    int  // Degrees (0-359)
	Arcminutes int  // Arcminutes (0-59)
	Arcseconds int  // Arcseconds (0-59)
	Fraction   int  // Fractional arcseconds
}

// HMS represents an angle in hours, minutes, seconds format.
// Used for right ascension and hour angle formatting.
type HMS struct {
	Sign     byte // '+' or '-'
	Hours    int  // Hours (0-23)
	Minutes  int  // Minutes (0-59)
	Seconds  int  // Seconds (0-59)
	Fraction int  // Fractional seconds
}

// Helper functions for internal use

// abs returns the absolute value of a float64
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

// sign returns the sign of y with the magnitude of x
// This implements the FORTRAN SIGN function
func sign(x, y float64) float64 {
	if y >= 0 {
		return math.Abs(x)
	}
	return -math.Abs(x)
}

// mod returns the modulus x mod y, always positive
func mod(x, y float64) float64 {
	result := math.Mod(x, y)
	if result < 0 {
		result += y
	}
	return result
}
