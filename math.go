package slagofa

import "math"

// Mathematical Utility Functions
//
// This file implements SLALIB-compatible mathematical utility functions.
// Most are simple wrappers around Go's math package.

// Mod returns the modulus (always positive) of a / b
//
// Note: This is a utility function, not a SLALIB function.
// SLALIB uses Fortran's intrinsic DMOD/MOD, not a separate function.
//
// Go equivalent: math.Mod with sign correction
//
// Parameters:
//   - a: Dividend
//   - b: Divisor (must be non-zero)
//
// Returns:
//   - Modulus in range [0, b) for positive b, (b, 0] for negative b
//
// Note: Unlike math.Mod, this always returns a result with the same sign as b
func Mod(a, b float64) float64 {
	r := math.Mod(a, b)
	if r < 0.0 && b > 0.0 {
		r += b
	} else if r > 0.0 && b < 0.0 {
		r += b
	}
	return r
}

// Mod32 returns the modulus (single precision)
func Mod32(a, b float32) float32 {
	return float32(Mod(float64(a), float64(b)))
}

// Sign returns the value of a with the sign of b
//
// Note: This is a utility function, not a SLALIB function.
// SLALIB uses Fortran's intrinsic DSIGN/SIGN, not a separate function.
//
// Go equivalent: math.Copysign
//
// Parameters:
//   - a: Value whose absolute value is used
//   - b: Value whose sign is used
//
// Returns:
//   - |a| if b >= 0, -|a| if b < 0
func Sign(a, b float64) float64 {
	return math.Copysign(a, b)
}

// Sign32 returns the value with transferred sign (single precision)
func Sign32(a, b float32) float32 {
	return float32(math.Copysign(float64(a), float64(b)))
}

// Trunc returns the integer part of x (truncates towards zero)
//
// Note: This is a utility function, not a SLALIB function.
// SLALIB uses Fortran's intrinsic DINT/AINT/INT, not separate functions.
//
// Go equivalent: math.Trunc
//
// Parameters:
//   - x: Value to truncate
//
// Returns:
//   - Integer part of x
//
// Examples:
//   - Trunc(3.7) = 3.0
//   - Trunc(-3.7) = -3.0
func Trunc(x float64) float64 {
	return math.Trunc(x)
}

// Trunc32 returns the integer part (single precision)
func Trunc32(x float32) float32 {
	return float32(math.Trunc(float64(x)))
}

// Round returns the nearest integer value
//
// Note: This is a utility function, not a SLALIB function.
// SLALIB uses Fortran's intrinsic ANINT/NINT, not separate functions.
//
// Go equivalent: math.Round
//
// Parameters:
//   - x: Value to round
//
// Returns:
//   - Nearest integer to x, rounding half-values away from zero
//
// Examples:
//   - Round(3.5) = 4.0
//   - Round(3.4) = 3.0
//   - Round(-3.5) = -4.0
func Round(x float64) float64 {
	return math.Round(x)
}

// Round32 returns the nearest integer value (single precision)
func Round32(x float32) float32 {
	return float32(math.Round(float64(x)))
}

// Poly evaluates a polynomial using Horner's method
//
// Original FORTRAN: sla_POLY by P.T. Wallace
//
// Parameters:
//   - coeffs: Polynomial coefficients [a0, a1, a2, ..., an]
//             Polynomial is: a0 + a1*x + a2*x^2 + ... + an*x^n
//   - x: Value at which to evaluate the polynomial
//
// Returns:
//   - Value of polynomial at x
//
// Example:
//   - Poly([1, 2, 3], 2.0) evaluates 1 + 2*2 + 3*2^2 = 1 + 4 + 12 = 17
//
// Note: Uses Horner's method for numerical stability and efficiency:
//       a0 + x(a1 + x(a2 + x(...)))
func Poly(coeffs []float64, x float64) float64 {
	if len(coeffs) == 0 {
		return 0.0
	}

	// Horner's method: start with highest coefficient
	result := coeffs[len(coeffs)-1]
	for i := len(coeffs) - 2; i >= 0; i-- {
		result = result*x + coeffs[i]
	}
	return result
}


// Poly32 evaluates a polynomial (single precision)
func Poly32(coeffs []float32, x float32) float32 {
	if len(coeffs) == 0 {
		return 0.0
	}

	result := coeffs[len(coeffs)-1]
	for i := len(coeffs) - 2; i >= 0; i-- {
		result = result*x + coeffs[i]
	}
	return result
}

// PolyMod evaluates a polynomial and applies modulus
//
// Original FORTRAN: sla_POLMO by P.T. Wallace (not in standard SLALIB docs)
//
// This function evaluates a polynomial using Horner's method and applies
// modulus at each step to keep intermediate results bounded. Useful for
// computing angles that wrap around.
//
// Parameters:
//   - coeffs: Polynomial coefficients [a0, a1, a2, ..., an]
//   - x: Value at which to evaluate the polynomial
//   - modulus: Modulus to apply at each step
//
// Returns:
//   - Value of polynomial at x, with modulus applied
//
// Example:
//   - PolyMod([1, 2, 3], 2.0, 10.0) computes polynomial mod 10 at each step
func PolyMod(coeffs []float64, x, modulus float64) float64 {
	if len(coeffs) == 0 {
		return 0.0
	}

	// Horner's method with modulus at each step
	result := Mod(coeffs[len(coeffs)-1], modulus)
	for i := len(coeffs) - 2; i >= 0; i-- {
		result = Mod(result*x+coeffs[i], modulus)
	}
	return result
}


// PolyMod32 evaluates a polynomial with modulus (single precision)
func PolyMod32(coeffs []float32, x, modulus float32) float32 {
	if len(coeffs) == 0 {
		return 0.0
	}

	result := Mod32(coeffs[len(coeffs)-1], modulus)
	for i := len(coeffs) - 2; i >= 0; i-- {
		result = Mod32(result*x+coeffs[i], modulus)
	}
	return result
}
