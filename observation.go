package slagofa

import (
	"math"

	"github.com/hebl/gofa"
)

// Observation Planning Functions
//
// This file implements SLALIB-compatible functions for observation planning
// including position angles, zenith distance, and telescope mount calculations.

// PositionAngle32 computes the position angle of one point with respect
// to another, given as direction vectors (single precision).
//
// Original FORTRAN: sla_PAV by P.T. Wallace
// Implementation: Converts float32 to float64, calls PositionAngle
//
// Parameters:
//   - v1: Direction of reference point
//   - v2: Direction of point whose PA is required
//
// Returns:
//   - Position angle in radians (-π to +π)
//
// Notes:
//   - If v2 is north of v1, PA ≈ 0
//   - If v2 is east of v1, PA ≈ +π/2
//   - Zero is returned if the two points are coincident
//   - v1 and v2 do not have to be unit vectors
func PositionAngle32(v1, v2 Vec3_32) float32 {
	v1_64 := Vec3{float64(v1[0]), float64(v1[1]), float64(v1[2])}
	v2_64 := Vec3{float64(v2[0]), float64(v2[1]), float64(v2[2])}
	return float32(gofa.Pap(v1_64, v2_64))
}

// Pav is a SLALIB-compatible alias for PositionAngle32 (sla_PAV)
var Pav = PositionAngle32

// ZenithDistance computes the zenith distance from hour angle, declination,
// and observatory latitude.
//
// Original FORTRAN: sla_ZD by P.T. Wallace
// Implementation: Direct port of SLALIB trigonometric calculation
//
// Parameters:
//   - ha: Hour angle in radians
//   - dec: Declination in radians
//   - phi: Observatory latitude in radians (geodetic)
//
// Returns:
//   - Zenith distance in radians (range 0 to π)
//
// Notes:
//   - The latitude must be geodetic. In critical applications,
//     corrections for polar motion should be applied.
//   - It may be important to distinguish between zenith distance
//     as affected by refraction (use "observed" HA,Dec) and zenith
//     distance in vacuo (use "topocentric" HA,Dec).
//   - If diurnal aberration can be neglected, "apparent" HA,Dec
//     may be used instead of topocentric HA,Dec.
//   - No range checking of arguments is performed.
func ZenithDistance(ha, dec, phi float64) float64 {
	// Compute sines and cosines
	sh := math.Sin(ha)
	ch := math.Cos(ha)
	sd := math.Sin(dec)
	cd := math.Cos(dec)
	sp := math.Sin(phi)
	cp := math.Cos(phi)

	// Compute zenith distance using vector components
	x := ch*cd*sp - sd*cp
	y := sh * cd
	z := ch*cd*cp + sd*sp

	return math.Atan2(math.Sqrt(x*x+y*y), z)
}

// Zd is a SLALIB-compatible alias for ZenithDistance (sla_ZD)
var Zd = ZenithDistance

// AltazResult holds the output from the Altaz function
type AltazResult struct {
	Az   float64 // Azimuth (radians, range 0-2π, north=0, east=π/2)
	Azd  float64 // Azimuth velocity (radians per radian of HA)
	Azdd float64 // Azimuth acceleration (radians per radian of HA squared)
	El   float64 // Elevation (radians, range ±π)
	Eld  float64 // Elevation velocity (radians per radian of HA)
	Eldd float64 // Elevation acceleration (radians per radian of HA squared)
	Pa   float64 // Parallactic angle (radians, range ±π, +ve west of meridian)
	Pad  float64 // Parallactic angle velocity (radians per radian of HA)
	Padd float64 // Parallactic angle acceleration (radians per radian of HA squared)
}

// Altaz computes positions, velocities and accelerations for an
// altazimuth telescope mount.
//
// Original FORTRAN: sla_ALTAZ by P.T. Wallace
// Implementation: Direct port of SLALIB calculation
//
// Parameters:
//   - ha: Hour angle (radians)
//   - dec: Declination (radians)
//   - phi: Observatory latitude (radians, geodetic)
//
// Returns:
//   - AltazResult containing azimuth, elevation, parallactic angle,
//     and their velocities and accelerations
//
// Notes:
//   - Natural units are used throughout. HA, DEC, PHI, AZ, EL are in radians.
//   - The velocities and accelerations assume constant declination and
//     constant rate of change of hour angle (as for tracking a star).
//   - Units of velocities (Azd, Eld, Pad): radians per radian of HA
//   - Units of accelerations (Azdd, Eldd, Padd): radians per radian of HA squared
//   - To convert to practical units:
//       angles × 360/2π → degrees
//       velocities × (2π/86400)×(360/2π) → degree/sec
//       accelerations × ((2π/86400)²)×(360/2π) → degree/sec/sec
//   - The velocity and acceleration factors assume sidereal tracking.
//     Their numerical values are (exactly) 1/240 and (approx) 1/3300236.9
//   - Azimuth is in range 0-2π; north is zero, east is +π/2
//   - Elevation and parallactic angle are in range ±π
//   - Parallactic angle is +ve for a star west of the meridian
//     and is the angle NP-star-zenith
//   - The latitude is geodetic (not geocentric)
//   - The hour angle and declination are topocentric
//   - Refraction and telescope mounting deficiencies are ignored
//   - No range checking of arguments is performed
func Altaz(ha, dec, phi float64) AltazResult {
	const tiny = 1.0e-30 // Clamp value to avoid division by zero at zenith/nadir

	// Compute sines and cosines
	sh := math.Sin(ha)
	ch := math.Cos(ha)
	sd := math.Sin(dec)
	cd := math.Cos(dec)
	sp := math.Sin(phi)
	cp := math.Cos(phi)
	chcd := ch * cd
	sdcp := sd * cp

	// Compute position vector components
	x := -chcd*sp + sdcp
	y := -sh * cd
	z := chcd*cp + sd*sp
	rsq := x*x + y*y
	r := math.Sqrt(rsq)

	// Azimuth and elevation
	var a float64
	if rsq == 0.0 {
		a = 0.0
	} else {
		a = math.Atan2(y, x)
	}
	if a < 0.0 {
		a += TwoPi
	}
	e := math.Atan2(z, r)

	// Parallactic angle
	c := cd*sp - ch*sdcp
	s := sh * cp
	var q float64
	if c*c+s*s > 0 {
		q = math.Atan2(s, c)
	} else {
		q = Pi - ha
	}

	// Velocities and accelerations (clamped at zenith/nadir)
	rsqClamped := rsq
	rClamped := r
	if rsq < tiny {
		rsqClamped = tiny
		rClamped = math.Sqrt(rsqClamped)
	}

	qd := -x * cp / rsqClamped
	ad := sp + z*qd
	ed := cp * y / rClamped
	edr := ed / rClamped
	add := edr * (z*sp + (2.0-rsqClamped)*qd)
	edd := -rClamped * qd * ad
	qdd := edr * (sp + 2.0*z*qd)

	return AltazResult{
		Az:   a,
		Azd:  ad,
		Azdd: add,
		El:   e,
		Eld:  ed,
		Eldd: edd,
		Pa:   q,
		Pad:  qd,
		Padd: qdd,
	}
}

// SLALIB-compatible lowercase aliases

// altaz is a SLALIB-compatible alias for Altaz (sla_ALTAZ)
var altaz = Altaz
