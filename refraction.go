package slagofa

import (
	"math"
)

// Phase 8: Refraction and Atmospheric Functions
//
// These functions handle atmospheric refraction corrections for astronomical
// observations, critical for converting between observed and true positions.
//
// IMPORTANT IMPLEMENTATION NOTE:
// Unlike most slagofa functions (which wrap GoFA), the refraction functions are
// DIRECT PORTS from SLALIB/PAL. This architectural divergence is necessary because:
//
// 1. GoFA provides only simplified refraction (gofa.Refco with empirical formulas)
// 2. High-precision astronomy requires full atmospheric integration (0.1 arcsec)
// 3. SLALIB's Hohenkerk & Sinclair method achieves this via numerical integration
//
// We ported ~500 lines of SLALIB/PAL code to maintain full SLALIB fidelity:
// - atmt/atms: Troposphere/stratosphere refractive index models
// - Refro: Full Simpson's Rule integration through atmospheric layers
// - Refco: Fits tan(Z) model by calling Refro at two zenith distances
//
// Performance: Zero heap allocations maintained despite complex integration.
// Accuracy: Matches SLALIB test vectors within 1e-8 to 1e-12.
//
// For applications requiring highest precision refraction corrections (sub-arcsecond),
// this direct port approach ensures compatibility with existing SLALIB-based systems.

// Internal helper functions for atmospheric refraction calculations

// atmt computes refractive index and derivative for the troposphere.
//
// Calculates temperature, refractive index, and rate of change at a given
// distance from Earth's center, modeling the troposphere (0-11km).
//
// Based on PAL's pal1Atmt, from SLALIB's internal ATMT subroutine.
func atmt(r0, t0, alpha, gamm2, delm2, c1, c2, c3, c4, c5, c6, r float64) (t, dn, rdndr float64) {
	// Temperature at height r (clamped to safe range)
	t = math.Max(math.Min(t0-alpha*(r-r0), 320.0), 100.0)

	// Temperature ratios
	tt0 := t / t0
	tt0gm2 := math.Pow(tt0, gamm2)
	tt0dm2 := math.Pow(tt0, delm2)

	// Refractive index
	dn = 1.0 + (c1*tt0gm2-(c2-c5/t)*tt0dm2)*tt0

	// r * rate of change of refractive index
	rdndr = r * (-c3*tt0gm2 + (c4-c6/tt0)*tt0dm2)

	return t, dn, rdndr
}

// atms computes refractive index and derivative for the stratosphere.
//
// Calculates refractive index and rate of change at a given distance from
// Earth's center, modeling the stratosphere (11-80km).
//
// Based on PAL's pal1Atms, from SLALIB's internal ATMS subroutine.
func atms(rt, tt, dnt, gamal, r float64) (dn, rdndr float64) {
	b := gamal / tt
	w := (dnt - 1.0) * math.Exp(-b*(r-rt))
	dn = 1.0 + w
	rdndr = -r * b * w
	return dn, rdndr
}

// Refco determines the constants A and B in the atmospheric refraction model.
//
// The refraction model is: dZ = A*tan(Z) + B*tan³(Z)
// where Z is the observed zenith distance and dZ is the refraction correction.
//
// Original FORTRAN: sla_REFCO by P.T. Wallace
// Go equivalent: Direct port from PAL/SLALIB
// SLALIB reference: SUN/67 section 79
//
// Parameters:
//   - hm: Height of observer above sea level (meters)
//   - tdk: Ambient temperature at observer (Kelvin)
//   - pmb: Pressure at observer (millibars)
//   - rh: Relative humidity at observer (0-1)
//   - wl: Effective wavelength (micrometers)
//   - phi: Latitude of observer (radians, astronomical)
//   - tlr: Temperature lapse rate in troposphere (K/meter)
//   - eps: Precision required to terminate iteration (radians)
//
// Returns:
//   - refa: tan(Z) coefficient (radians)
//   - refb: tan³(Z) coefficient (radians)
//
// Notes:
//   - Achieves 0.5 arcsec accuracy for ZD < 80°
//   - Achieves 0.01 arcsec accuracy for ZD < 60°
//   - Achieves 0.001 arcsec accuracy for ZD < 45°
//   - Fits constants using Refro at two sample zenith distances
//   - Typical values: tlr=0.0065, eps=1e-10
func Refco(hm, tdk, pmb, rh, wl, phi, tlr, eps float64) (refa, refb float64) {
	// Sample zenith distances: arctan(1) = 45° and arctan(4) ≈ 76°
	const atn1 = 0.7853981633974483 // arctan(1)
	const atn4 = 1.325817663668033  // arctan(4)

	// Determine refraction for the two sample zenith distances
	r1 := Refro(atn1, hm, tdk, pmb, rh, wl, phi, tlr, eps)
	r2 := Refro(atn4, hm, tdk, pmb, rh, wl, phi, tlr, eps)

	// Solve for refraction constants
	refa = (64.0*r1 - r2) / 60.0
	refb = (r2 - 4.0*r1) / 60.0

	return refa, refb
}

// Refcoq is a quick refraction coefficient calculation.
//
// Simpler version of Refco that doesn't account for height, latitude,
// or temperature lapse rate. Suitable for most applications.
//
// Original FORTRAN: sla_REFCOQ by P.T. Wallace
// Go equivalent: Simplified atmospheric model (Stone/Green formula)
// SLALIB reference: SUN/67 section 80
//
// Parameters:
//   - tdk: Ambient temperature (Kelvin)
//   - pmb: Pressure (millibars)
//   - rh: Relative humidity (0-1)
//   - wl: Wavelength (micrometers)
//
// Returns:
//   - refa: tan(Z) coefficient (radians)
//   - refb: tan³(Z) coefficient (radians)
//
// Notes:
//   - This is the most commonly used refraction function
//   - Standard conditions: 1010 mb, 283K (10°C), 0.5 RH
//   - Uses empirical formulas from Stone (1996) and Green
func Refcoq(tdk, pmb, rh, wl float64) (refa, refb float64) {
	// Decide whether optical/IR or radio case - switch at 100 microns
	optic := wl <= 100.0

	// Restrict parameters to safe values
	t := math.Max(tdk-273.15, -150.0) // Convert to Celsius
	t = math.Min(t, 200.0)
	p := math.Max(pmb, 0.0)
	p = math.Min(p, 10000.0)
	r := math.Max(rh, 0.0)
	r = math.Min(r, 1.0)
	w := math.Max(wl, 0.1)
	w = math.Min(w, 1e6)

	// Water vapour pressure at the observer
	var pw float64
	if p > 0.0 {
		ps := math.Pow(10.0, (0.7859+0.03477*t)/(1.0+0.00412*t)) *
			(1.0 + p*(4.5e-6+6e-10*t*t))
		pw = r * ps / (1.0 - (1.0-r)*ps/p)
	} else {
		pw = 0.0
	}

	// Refractive index minus 1 at the observer
	tk := t + 273.15
	var gamma float64
	if optic {
		wlsq := w * w
		gamma = ((77.53484e-6+(4.39108e-7+3.666e-9/wlsq)/wlsq)*p -
			11.2684e-6*pw) / tk
	} else {
		gamma = (77.6890e-6*p - (6.3938e-6-0.375463/tk)*pw) / tk
	}

	// Formula for beta from Stone, with empirical adjustments
	beta := 4.4474e-6 * tk
	if !optic {
		beta -= 0.0074 * pw * beta
	}

	// Refraction constants from Green
	refa = gamma * (1.0 - beta)
	refb = -gamma * (beta - gamma/2.0)

	return refa, refb
}

// Refz computes refraction for a given zenith distance (fast version).
//
// Applies the refraction model: dZ = A*tan(Z) + B*tan³(Z)
//
// Original FORTRAN: sla_REFZ by P.T. Wallace
// Go equivalent: Custom using refraction model
// SLALIB reference: SUN/67 section 82
//
// Parameters:
//   - zu: Unrefracted zenith distance (radians)
//   - refa: tan(Z) coefficient from Refco (radians)
//   - refb: tan³(Z) coefficient from Refco (radians)
//
// Returns:
//   - zr: Refracted (observed) zenith distance (radians)
//
// Notes:
//   - Fast, non-iterative
//   - Goes from true → observed (adds refraction)
//   - For observed → true, use Refro
func Refz(zu, refa, refb float64) float64 {
	// Compute tan(zu)
	tanzu := math.Tan(zu)

	// Apply refraction model: dZ = A*tan(Z) + B*tan³(Z)
	dz := refa*tanzu + refb*tanzu*tanzu*tanzu

	// Refracted zenith distance
	return zu + dz
}

// Refro computes refraction for observed zenith distance using full atmospheric integration.
//
// Removes refraction to get true zenith distance from observed using the
// Hohenkerk & Sinclair method with numerical integration through atmospheric layers.
//
// Original FORTRAN: sla_REFRO by P.T. Wallace
// Go equivalent: Direct port from PAL/SLALIB
// SLALIB reference: SUN/67 section 81
//
// Parameters:
//   - zobs: Observed (refracted) zenith distance (radians)
//   - hm: Height above sea level (meters)
//   - tdk: Ambient temperature (Kelvin)
//   - pmb: Pressure (millibars)
//   - rh: Relative humidity (0-1)
//   - wl: Wavelength (micrometers)
//   - phi: Latitude (radians)
//   - tlr: Temperature lapse rate (K/meter)
//   - eps: Precision (radians)
//
// Returns:
//   - ref: Refraction correction (radians): in vacuo ZD minus observed ZD
//
// Notes:
//   - Achieves 0.1 arcsec accuracy for moderate ZD
//   - Uses Simpson's Rule integration with adaptive refinement
//   - Separate integration through troposphere and stratosphere
//   - Radio refraction chosen by wl > 100 micrometers
//   - Typical values: tlr=0.0065, eps=1e-8
func Refro(zobs, hm, tdk, pmb, rh, wl, phi, tlr, eps float64) float64 {
	// Fixed parameters
	const d93 = 1.623156204      // 93 degrees in radians
	const gcr = 8314.32          // Universal gas constant
	const dmd = 28.9644          // Molecular weight of dry air
	const dmw = 18.0152          // Molecular weight of water vapour
	const s = 6378120.0          // Mean Earth radius (meters)
	const delta = 18.36          // Exponent of temp dependence of water vapour
	const ht = 11000.0           // Height of tropopause (meters)
	const hs = 80000.0           // Upper limit for refractive effects (meters)
	const ismax = 16384          // Max number of strips for integration

	// The refraction integrand
	refi := func(dn, rdndr float64) float64 {
		return rdndr / (dn + rdndr)
	}

	// Transform ZOBS into the normal range
	zobs1 := NormalizeAngle(zobs)
	zobs2 := math.Min(math.Abs(zobs1), d93)

	// Keep other arguments within safe bounds
	hmok := math.Min(math.Max(hm, -1e3), hs)
	tdkok := math.Min(math.Max(tdk, 100.0), 500.0)
	pmbok := math.Min(math.Max(pmb, 0.0), 10000.0)
	rhok := math.Min(math.Max(rh, 0.0), 1.0)
	wlok := math.Max(wl, 0.1)
	alpha := math.Min(math.Max(math.Abs(tlr), 0.001), 0.01)

	// Tolerance for iteration
	tol := math.Min(math.Max(math.Abs(eps), 1e-12), 0.1) / 2.0

	// Decide whether optical/IR or radio case - switch at 100 microns
	optic := wlok <= 100.0

	// Set up model atmosphere parameters defined at the observer
	wlsq := wlok * wlok
	gb := 9.784 * (1.0 - 0.0026*math.Cos(phi+phi) - 0.00000028*hmok)
	var a float64
	if optic {
		a = (287.6155 + (1.62887+0.01360/wlsq)/wlsq) * 273.15e-6 / 1013.25
	} else {
		a = 77.6890e-6
	}
	gamal := (gb * dmd) / gcr
	gamma := gamal / alpha
	gamm2 := gamma - 2.0
	delm2 := delta - 2.0
	tdc := tdkok - 273.15
	psat := math.Pow(10.0, (0.7859+0.03477*tdc)/(1.0+0.00412*tdc)) *
		(1.0 + pmbok*(4.5e-6+6.0e-10*tdc*tdc))
	var pwo float64
	if pmbok > 0.0 {
		pwo = rhok * psat / (1.0 - (1.0-rhok)*psat/pmbok)
	} else {
		pwo = 0.0
	}
	w := pwo * (1.0 - dmw/dmd) * gamma / (delta - gamma)
	c1 := a * (pmbok + w) / tdkok
	var c2 float64
	if optic {
		c2 = (a*w + 11.2684e-6*pwo) / tdkok
	} else {
		c2 = (a*w + 6.3938e-6*pwo) / tdkok
	}
	c3 := (gamma - 1.0) * alpha * c1 / tdkok
	c4 := (delta - 1.0) * alpha * c2 / tdkok
	var c5, c6 float64
	if optic {
		c5 = 0.0
		c6 = 0.0
	} else {
		c5 = 375463e-6 * pwo / tdkok
		c6 = c5 * delm2 * alpha / (tdkok * tdkok)
	}

	// Conditions at the observer
	r0 := s + hmok
	_, dn0, rdndr0 := atmt(r0, tdkok, alpha, gamm2, delm2, c1, c2, c3, c4, c5, c6, r0)
	sk0 := dn0 * r0 * math.Sin(zobs2)
	f0 := refi(dn0, rdndr0)

	// Conditions in the troposphere at the tropopause
	rt := s + math.Max(ht, hmok)
	tt, dnt, rdndrt := atmt(r0, tdkok, alpha, gamm2, delm2, c1, c2, c3, c4, c5, c6, rt)
	sine := sk0 / (rt * dnt)
	zt := math.Atan2(sine, math.Sqrt(math.Max(1.0-sine*sine, 0.0)))
	ft := refi(dnt, rdndrt)

	// Conditions in the stratosphere at the tropopause
	dnts, rdndrp := atms(rt, tt, dnt, gamal, rt)
	sine = sk0 / (rt * dnts)
	zts := math.Atan2(sine, math.Sqrt(math.Max(1.0-sine*sine, 0.0)))
	fts := refi(dnts, rdndrp)

	// Conditions at the stratosphere limit
	rs := s + hs
	dns, rdndrs := atms(rt, tt, dnt, gamal, rs)
	sine = sk0 / (rs * dns)
	zs := math.Atan2(sine, math.Sqrt(math.Max(1.0-sine*sine, 0.0)))
	fs := refi(dns, rdndrs)

	// Variable initialization
	var reft, refs float64 // Troposphere and stratosphere components

	// Integrate the refraction integral in two parts:
	// first in troposphere (k=1), then in stratosphere (k=2)
	for k := 1; k <= 2; k++ {
		// Initialize previous refraction to ensure at least two iterations
		refold := 1.0

		// Start off with 8 strips
		is := 8

		// Start z, z range, and start and end values
		var z0, zrange, fb, ff float64
		if k == 1 {
			z0 = zobs2
			zrange = zt - z0
			fb = f0
			ff = ft
		} else {
			z0 = zts
			zrange = zs - z0
			fb = fts
			ff = fs
		}

		// Sums of odd and even values
		fo := 0.0
		fe := 0.0

		// First time through the loop we have to do every point
		n := 1

		// Start of iteration loop (terminates at specified precision)
		loop := true
		var refp float64
		for loop {
			// Strip width
			h := zrange / float64(is)

			// Initialize distance from Earth centre for quadrature pass
			var r float64
			if k == 1 {
				r = r0
			} else {
				r = rt
			}

			// One pass (no need to compute evens after first time)
			for i := 1; i < is; i += n {
				// Sine of observed zenith distance
				sz := math.Sin(z0 + h*float64(i))

				// Find r (to the nearest metre, maximum four iterations)
				if sz > 1e-20 {
					ww := sk0 / sz
					rg := r
					dr := 1.0e6
					for j := 0; j < 4 && math.Abs(dr) > 1.0; j++ {
						var dn, rdndr float64
						if k == 1 {
							_, dn, rdndr = atmt(r0, tdkok, alpha, gamm2, delm2,
								c1, c2, c3, c4, c5, c6, rg)
						} else {
							dn, rdndr = atms(rt, tt, dnt, gamal, rg)
						}
						dr = (rg*dn - ww) / (dn + rdndr)
						rg = rg - dr
					}
					r = rg
				}

				// Find the refractive index and integrand at r
				var dn, rdndr float64
				if k == 1 {
					_, dn, rdndr = atmt(r0, tdkok, alpha, gamm2, delm2,
						c1, c2, c3, c4, c5, c6, r)
				} else {
					dn, rdndr = atms(rt, tt, dnt, gamal, r)
				}
				f := refi(dn, rdndr)

				// Accumulate odd and (first time only) even values
				if n == 1 && i%2 == 0 {
					fe += f
				} else {
					fo += f
				}
			}

			// Evaluate the integrand using Simpson's Rule
			refp = h * (fb + 4.0*fo + 2.0*fe + ff) / 3.0

			// Has the required precision been achieved (or can't be)?
			if math.Abs(refp-refold) > tol && is < ismax {
				// No: prepare for next iteration
				refold = refp
				is += is      // Double the number of strips
				fe += fo      // Sum of all current values = next pass's even values
				fo = 0.0      // Prepare for new odd values
				n = 2         // Skip even values next time
			} else {
				// Yes: save component and terminate the loop
				if k == 1 {
					reft = refp
				} else {
					refs = refp
				}
				loop = false
			}
		}
	}

	// Result: sum of troposphere and stratosphere components
	ref := reft + refs

	if zobs1 < 0.0 {
		ref = -ref
	}

	return ref
}

// Refv applies refraction to a direction vector.
//
// Adjusts a direction vector for atmospheric refraction.
//
// Original FORTRAN: sla_REFV by P.T. Wallace
// Go equivalent: Custom using Refco + vector operations
// SLALIB reference: SUN/67 section 83
//
// Parameters:
//   - vu: Unrefracted (true) direction vector [x,y,z]
//   - refa: tan(Z) coefficient from Refco (radians)
//   - refb: tan³(Z) coefficient from Refco (radians)
//
// Returns:
//   - vr: Refracted (observed) direction vector [x,y,z]
//
// Notes:
//   - Vector must be a unit vector
//   - Refraction applied in direction of zenith
//   - Returns normalized vector
func Refv(vu Vec3, refa, refb float64) Vec3 {
	// Extract altitude from vector
	// Altitude = arcsin(z-component)
	alt := math.Asin(vu[2])

	// Zenith distance
	zu := math.Pi/2 - alt

	// Apply refraction
	zr := Refz(zu, refa, refb)

	// New altitude
	altr := math.Pi/2 - zr

	// Azimuth is unchanged
	// Calculate azimuth from x,y components
	az := math.Atan2(vu[1], vu[0])

	// Reconstruct vector with refracted altitude
	cosalt := math.Cos(altr)
	vr := Vec3{
		cosalt * math.Cos(az),
		cosalt * math.Sin(az),
		math.Sin(altr),
	}

	return vr
}

// Airmas computes air mass for a given zenith distance.
//
// Air mass is the relative path length through the atmosphere compared
// to the zenith direction.
//
// Original FORTRAN: sla_AIRMAS by P.T. Wallace
// Go equivalent: Hardie's (1962) polynomial fit
// SLALIB reference: SUN/67 section 6
//
// Parameters:
//   - zd: Observed zenith distance (radians)
//
// Returns:
//   - am: Air mass (relative to zenith)
//
// Notes:
//   - Uses Hardie's (1962) polynomial fit to Bemporad's data
//   - Accurate to better than 0.1% up to X = 6.8
//   - ZD is clamped at 87 degrees to avoid arithmetic overflows
//   - Sign of ZD is ignored
func Airmas(zd float64) float64 {
	// Clamp zenith distance at 87 degrees (1.52 radians)
	// and take absolute value
	zdAbs := math.Abs(zd)
	if zdAbs > 1.52 {
		zdAbs = 1.52
	}

	// sec(z) - 1
	seczm1 := 1.0/math.Cos(zdAbs) - 1.0

	// Hardie's polynomial fit
	am := 1.0 + seczm1*(0.9981833-seczm1*(0.002875+0.0008083*seczm1))

	return am
}

// Pa computes parallactic angle for an object.
//
// The parallactic angle is the angle between the direction to the
// celestial north pole and the direction to the zenith, measured
// at the object.
//
// Original FORTRAN: sla_PA by P.T. Wallace
// Go equivalent: Custom calculation
// SLALIB reference: SUN/67 section 74
//
// Parameters:
//   - ha: Hour angle (radians)
//   - dec: Declination (radians)
//   - phi: Observer's latitude (radians)
//
// Returns:
//   - Parallactic angle (radians, -π to +π)
//
// Notes:
//   - Positive when field rotates anticlockwise as seen from below
//   - Zero at transit
//   - ±π/2 at horizon
//   - Returns zero at the pole (when both components are zero)
func Pa(ha, dec, phi float64) float64 {
	// Compute parallactic angle using:
	// tan(PA) = sin(HA) / (cos(dec)*tan(phi) - sin(dec)*cos(HA))

	cosPhi := math.Cos(phi)
	sinPhi := math.Sin(phi)
	sinHA := math.Sin(ha)
	cosHA := math.Cos(ha)
	sinDec := math.Sin(dec)
	cosDec := math.Cos(dec)

	// Numerator: cos(phi) * sin(HA)
	sqsz := cosPhi * sinHA

	// Denominator: sin(phi)*cos(dec) - cos(phi)*sin(dec)*cos(HA)
	cqsz := sinPhi*cosDec - cosPhi*sinDec*cosHA

	// Special case at pole: if both are zero, set denominator to 1
	if sqsz == 0.0 && cqsz == 0.0 {
		cqsz = 1.0
	}

	// Parallactic angle
	pa := math.Atan2(sqsz, cqsz)

	return pa
}
