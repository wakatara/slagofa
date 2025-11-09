package slagofa

import (
	"math"

	"github.com/hebl/gofa"
)

// Coordinate Conversion Functions
// These functions wrap GoFA's spherical/Cartesian conversion functions
// to provide SLALIB-compatible API.

// SphericalToCartesian converts spherical coordinates to Cartesian direction cosines.
//
// Given spherical coordinates (longitude, latitude), this function computes
// the corresponding unit vector in Cartesian coordinates.
//
// Original FORTRAN: sla_DCS2C by P.T. Wallace
// GoFA equivalent: gofa.S2c (spherical coordinates to unit vector)
//
// Parameters:
//   - longitude: Longitude angle (e.g., RA, azimuth) in radians
//   - latitude: Latitude angle (e.g., Dec, altitude) in radians
//
// Returns:
//   - Direction cosines (unit vector)
func SphericalToCartesian(longitude, latitude float64) Vec3 {
	var result [3]float64
	gofa.S2c(longitude, latitude, &result)
	return Vec3(result)
}

// Dcs2c is a SLALIB-compatible alias for SphericalToCartesian (sla_DCS2C)
func Dcs2c(longitude, latitude float64) Vec3 {
	return SphericalToCartesian(longitude, latitude)
}

// CartesianToSpherical converts Cartesian coordinates to spherical.
//
// Given a direction vector (which need not be normalized), this function
// computes the corresponding spherical coordinates (longitude, latitude).
//
// Original FORTRAN: sla_DCC2S by P.T. Wallace
// GoFA equivalent: gofa.C2s (p-vector to spherical coordinates)
//
// Parameters:
//   - v: Direction vector (any magnitude)
//
// Returns:
//   - longitude: Longitude in radians
//   - latitude: Latitude in radians
//
// Notes:
//   - If v is null, zero longitude and latitude are returned
//   - At either pole, zero longitude is returned
func CartesianToSpherical(v Vec3) (longitude, latitude float64) {
	gofa.C2s(v, &longitude, &latitude)
	return longitude, latitude
}

// Dcc2s is a SLALIB-compatible wrapper for CartesianToSpherical (sla_DCC2S)
func Dcc2s(v Vec3) (float64, float64) {
	return CartesianToSpherical(v)
}

// SphericalPolarToCartesian converts spherical polar coordinates to Cartesian.
//
// Given spherical polar coordinates (longitude, latitude, radius), this function
// computes the corresponding Cartesian position vector.
//
// Original FORTRAN: sla_DTP2S by P.T. Wallace
// GoFA equivalent: gofa.S2p (spherical polar to p-vector)
//
// Parameters:
//   - longitude: Longitude angle in radians
//   - latitude: Latitude angle in radians
//   - radius: Radial distance
//
// Returns:
//   - Cartesian position vector
func SphericalPolarToCartesian(longitude, latitude, radius float64) Vec3 {
	var result [3]float64
	gofa.S2p(longitude, latitude, radius, &result)
	return Vec3(result)
}

// NOTE: SphericalPolarToCartesian does not have a direct SLALIB equivalent.
// This is a GoFA utility function (gofa.S2p).

// CartesianToSphericalPolar converts Cartesian to spherical polar coordinates.
//
// Given a Cartesian position vector, this function computes the corresponding
// spherical polar coordinates (longitude, latitude, radius).
//
// Original FORTRAN: sla_DTS2C by P.T. Wallace
// GoFA equivalent: gofa.P2s (p-vector to spherical polar)
//
// Parameters:
//   - v: Cartesian position vector
//
// Returns:
//   - longitude: Longitude in radians
//   - latitude: Latitude in radians
//   - radius: Radial distance
//
// Notes:
//   - If v is null, zero values are returned
//   - At either pole, zero longitude is returned
func CartesianToSphericalPolar(v Vec3) (longitude, latitude, radius float64) {
	gofa.P2s(v, &longitude, &latitude, &radius)
	return longitude, latitude, radius
}

// NOTE: CartesianToSphericalPolar does not have a direct SLALIB equivalent.
// This is a GoFA utility function (gofa.P2s).

// Single-precision coordinate conversions

// SphericalToCartesian32 converts spherical coordinates to Cartesian (single precision).
//
// Original FORTRAN: sla_CS2C by P.T. Wallace
//
// Parameters:
//   - longitude: Longitude angle in radians
//   - latitude: Latitude angle in radians
//
// Returns:
//   - Direction cosines (unit vector)
func SphericalToCartesian32(longitude, latitude float32) Vec3_32 {
	var result [3]float64
	gofa.S2c(float64(longitude), float64(latitude), &result)
	return Vec3_32{float32(result[0]), float32(result[1]), float32(result[2])}
}

// Cs2c is a SLALIB-compatible alias for SphericalToCartesian32 (sla_CS2C)
func Cs2c(longitude, latitude float32) Vec3_32 {
	return SphericalToCartesian32(longitude, latitude)
}

// CartesianToSpherical32 converts Cartesian to spherical coordinates (single precision).
//
// Original FORTRAN: sla_CC2S by P.T. Wallace
//
// Parameters:
//   - v: Direction vector (any magnitude)
//
// Returns:
//   - longitude: Longitude in radians
//   - latitude: Latitude in radians
func CartesianToSpherical32(v Vec3_32) (longitude, latitude float32) {
	v64 := [3]float64{float64(v[0]), float64(v[1]), float64(v[2])}
	var lon, lat float64
	gofa.C2s(v64, &lon, &lat)
	return float32(lon), float32(lat)
}

// Cc2s is a SLALIB-compatible wrapper for CartesianToSpherical32 (sla_CC2S)
func Cc2s(v Vec3_32) (float32, float32) {
	return CartesianToSpherical32(v)
}

// Phase 5: Coordinate System Transformations
//
// These functions convert between different astronomical coordinate systems.

// De2h converts equatorial to horizon coordinates (HA, Dec → Az, El)
//
// This function converts hour angle and declination to azimuth and elevation
// for an observer at a given geodetic latitude.
//
// Original FORTRAN: sla_DE2H by P.T. Wallace
// GoFA equivalent: gofa.Ae2hd
// SLALIB reference: SUN/67 section 42
//
// Parameters:
//   - ha: Hour angle (radians, local apparent)
//   - dec: Declination (radians)
//   - phi: Observatory latitude (radians, geodetic)
//
// Returns:
//   - az: Azimuth (radians, N=0, E=π/2, range [0,2π))
//   - el: Elevation/altitude (radians, range [-π/2, +π/2])
//
// Notes:
//   - All angles in radians
//   - Azimuth: North=0, East=+π/2 (SLALIB convention)
//   - No range checking (SLALIB behavior)
//   - Latitude must be geodetic (apply polar motion corrections for critical applications)
//   - See SLALIB documentation for distinctions between observed/topocentric/apparent coordinates
func De2h(ha, dec, phi float64) (az, el float64) {
	gofa.Ae2hd(ha, dec, phi, &az, &el)
	return az, el
}

// Dh2e converts horizon to equatorial coordinates (Az, El → HA, Dec)
//
// This function is the inverse of De2h, converting azimuth and elevation
// to hour angle and declination.
//
// Original FORTRAN: sla_DH2E by P.T. Wallace
// GoFA equivalent: gofa.Ae2hd (note: Ae2hd converts Az/El → HA/Dec)
// SLALIB reference: SUN/67 section 52
//
// Parameters:
//   - az: Azimuth (radians, N=0, E=π/2)
//   - el: Elevation (radians)
//   - phi: Observatory latitude (radians, geodetic)
//
// Returns:
//   - ha: Hour angle (radians, range [-π, +π])
//   - dec: Declination (radians, range [-π/2, +π/2])
//
// Notes:
//   - Inverse of De2h
//   - HA returned in range [-π, +π] (SLALIB convention)
//   - No range checking (SLALIB behavior)
func Dh2e(az, el, phi float64) (ha, dec float64) {
	gofa.Ae2hd(az, el, phi, &ha, &dec)
	// Normalize HA to range [-π, +π] (SLALIB convention)
	// GoFA may return it in [0, 2π), so normalize using Anpm
	ha = gofa.Anpm(ha)
	return ha, dec
}

// Eqgal converts J2000.0 equatorial to IAU 1958 galactic coordinates
//
// This function transforms J2000.0 FK5 equatorial coordinates (RA, Dec)
// to IAU 1958 galactic coordinates (l, b).
//
// Original FORTRAN: sla_EQGAL by P.T. Wallace
// GoFA equivalent: gofa.Icrs2g
// SLALIB reference: SUN/67 section 46
//
// Parameters:
//   - dr: J2000.0 RA (radians, FK5)
//   - dd: J2000.0 Dec (radians, FK5)
//
// Returns:
//   - dl: Galactic longitude (radians, IAU 1958)
//   - db: Galactic latitude (radians, IAU 1958)
//
// Notes:
//   - Input coordinates MUST be J2000.0 FK5 equatorial
//   - Output is IAU 1958 galactic system
//   - For B1950.0 FK4 coordinates, use sla_GE50 first (not implemented yet)
//   - Reference: Blaauw et al, Mon.Not.R.Astron.Soc., 121, 123 (1960)
func Eqgal(dr, dd float64) (dl, db float64) {
	gofa.Icrs2g(dr, dd, &dl, &db)
	return dl, db
}

// Galeq converts IAU 1958 galactic to J2000.0 equatorial coordinates
//
// This function transforms IAU 1958 galactic coordinates (l, b)
// to J2000.0 FK5 equatorial coordinates (RA, Dec).
//
// Original FORTRAN: sla_GALEQ by P.T. Wallace
// GoFA equivalent: gofa.G2icrs
// SLALIB reference: SUN/67 section 49
//
// Parameters:
//   - dl: Galactic longitude (radians, IAU 1958)
//   - db: Galactic latitude (radians, IAU 1958)
//
// Returns:
//   - dr: J2000.0 RA (radians, FK5)
//   - dd: J2000.0 Dec (radians, FK5)
//
// Notes:
//   - Inverse of Eqgal
//   - Output is J2000.0 FK5 equatorial
//   - For B1950.0 FK4 output, use sla_G50E after (not implemented yet)
func Galeq(dl, db float64) (dr, dd float64) {
	gofa.G2icrs(dl, db, &dr, &dd)
	return dr, dd
}

// Eqecl converts equatorial to ecliptic coordinates (IAU 2006)
//
// This function transforms equatorial coordinates (RA, Dec) to ecliptic
// coordinates (longitude, latitude) using the IAU 2006 obliquity model.
//
// Original FORTRAN: sla_EQECL by P.T. Wallace
// GoFA equivalent: gofa.Eqec06
// SLALIB reference: SUN/67 section 45
//
// Parameters:
//   - dr: RA (radians)
//   - dd: Dec (radians)
//   - mjd: Modified Julian Date (TT)
//
// Returns:
//   - dl: Ecliptic longitude (radians)
//   - db: Ecliptic latitude (radians)
//
// Notes:
//   - Uses IAU 2006 obliquity model (more accurate than SLALIB's IAU 1980)
//   - Date is used to compute obliquity of ecliptic at that epoch
//   - Input date should be Terrestrial Time (TT)
//   - For mean ecliptic, use date = MJD of epoch
func Eqecl(dr, dd, mjd float64) (dl, db float64) {
	// GoFA.Eqec06 expects two-part Julian Date (TT)
	// Convert MJD to JD: JD = MJD + 2400000.5
	gofa.Eqec06(2400000.5, mjd, dr, dd, &dl, &db)
	return dl, db
}

// Ecleq converts ecliptic to equatorial coordinates (IAU 2006)
//
// This function transforms ecliptic coordinates (longitude, latitude) to
// equatorial coordinates (RA, Dec) using the IAU 2006 obliquity model.
//
// Original FORTRAN: sla_ECLEQ by P.T. Wallace
// GoFA equivalent: gofa.Eceq06
// SLALIB reference: SUN/67 section 39
//
// Parameters:
//   - dl: Ecliptic longitude (radians)
//   - db: Ecliptic latitude (radians)
//   - mjd: Modified Julian Date (TT)
//
// Returns:
//   - dr: RA (radians)
//   - dd: Dec (radians)
//
// Notes:
//   - Inverse of Eqecl
//   - Uses IAU 2006 obliquity model (more accurate than SLALIB's IAU 1980)
//   - Date is used to compute obliquity of ecliptic at that epoch
//   - Input date should be Terrestrial Time (TT)
func Ecleq(dl, db, mjd float64) (dr, dd float64) {
	// GoFA.Eceq06 expects two-part Julian Date (TT)
	// Convert MJD to JD: JD = MJD + 2400000.5
	gofa.Eceq06(2400000.5, mjd, dl, db, &dr, &dd)
	return dr, dd
}

// Geoc converts geodetic position to geocentric
//
// This function converts geodetic coordinates (longitude, latitude, height)
// to geocentric Cartesian coordinates (X, Y, Z).
//
// Original FORTRAN: sla_GEOC by P.T. Wallace
// GoFA equivalent: gofa.Gd2gc
// SLALIB reference: SUN/67 section 50
//
// Parameters:
//   - phi: Geodetic latitude (radians)
//   - height: Height above reference ellipsoid (meters)
//
// Returns:
//   - rho: Distance from Earth's center (Earth radii)
//   - z: Height above equatorial plane (Earth radii)
//
// Notes:
//   - Uses WGS84 ellipsoid via GoFA
//   - Returns geocentric cylindrical coordinates (rho, z)
//   - Rho = distance from polar axis
//   - SLALIB used different internal ellipsoid, small differences expected
func Geoc(phi, height float64) (rho, z float64) {
	// GoFA.Gd2gc returns Cartesian XYZ coordinates
	// We need to convert to SLALIB's (rho, z) cylindrical system
	// For Geoc, we assume longitude = 0
	var xyz [3]float64
	// n=1 for WGS84 ellipsoid
	gofa.Gd2gc(1, 0.0, phi, height, &xyz)

	// Convert Cartesian to cylindrical: rho = sqrt(x² + y²), z = z
	rho = math.Sqrt(xyz[0]*xyz[0] + xyz[1]*xyz[1])
	z = xyz[2]

	// Convert from meters to Earth radii (WGS84 equatorial radius)
	const earthRadius = 6378137.0 // meters
	rho /= earthRadius
	z /= earthRadius

	return rho, z
}

// SLALIB-compatible lowercase aliases

// cc2s is a SLALIB-compatible alias for Cc2s (sla_CC2S)
var cc2s = Cc2s

// cs2c is a SLALIB-compatible alias for Cs2c (sla_CS2C)
var cs2c = Cs2c

// dcc2s is a SLALIB-compatible alias for Dcc2s (sla_DCC2S)
var dcc2s = Dcc2s

// dcs2c is a SLALIB-compatible alias for Dcs2c (sla_DCS2C)
var dcs2c = Dcs2c

// de2h is a SLALIB-compatible alias for De2h (sla_DE2H)
var de2h = De2h

// dh2e is a SLALIB-compatible alias for Dh2e (sla_DH2E)
var dh2e = Dh2e

// ecleq is a SLALIB-compatible alias for Ecleq (sla_ECLEQ)
var ecleq = Ecleq

// eqecl is a SLALIB-compatible alias for Eqecl (sla_EQECL)
var eqecl = Eqecl

// eqgal is a SLALIB-compatible alias for Eqgal (sla_EQGAL)
var eqgal = Eqgal

// galeq is a SLALIB-compatible alias for Galeq (sla_GALEQ)
var galeq = Galeq

// geoc is a SLALIB-compatible alias for Geoc (sla_GEOC)
var geoc = Geoc
