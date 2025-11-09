package slagofa

import (
	"github.com/hebl/gofa"
	"math"
)

// Vector Operations
// These functions wrap GoFA's vector/matrix library (vml.go) to provide
// SLALIB-compatible API for 3-dimensional vector operations.

// DotProduct computes the scalar (dot) product of two 3-element vectors.
//
// The dot product is defined as: a[0]*b[0] + a[1]*b[1] + a[2]*b[2]
//
// Original FORTRAN: sla_DVDV by P.T. Wallace / Rutherford Appleton Laboratory
// GoFA equivalent: gofa.Pdp (p-vector dot product)
//
// Parameters:
//   - a: First 3-component vector
//   - b: Second 3-component vector
//
// Returns:
//   - The scalar product a · b
func DotProduct(a, b Vec3) float64 {
	return gofa.Pdp(a, b)
}

// Dvdv is a SLALIB-compatible alias for DotProduct (sla_DVDV)
var Dvdv = DotProduct

// CrossProduct computes the vector (cross) product of two 3-element vectors.
//
// The cross product produces a vector perpendicular to both input vectors,
// with magnitude equal to the area of the parallelogram formed by the inputs.
//
// Original FORTRAN: sla_DVXV by P.T. Wallace
// GoFA equivalent: gofa.Pxp (p-vector cross product)
//
// Parameters:
//   - a: First 3-component vector
//   - b: Second 3-component vector
//
// Returns:
//   - The vector product a × b
func CrossProduct(a, b Vec3) Vec3 {
	var result [3]float64
	gofa.Pxp(a, b, &result)
	return Vec3(result)
}

// Dvxv is a SLALIB-compatible alias for CrossProduct (sla_DVXV)
var Dvxv = CrossProduct

// NormalizeVector normalizes a 3-component vector and returns its modulus.
//
// The function calculates the magnitude (modulus) of the input vector and returns
// both the normalized unit vector and the original magnitude. If the input vector
// has zero magnitude, the output vector is set to zero.
//
// Original FORTRAN: sla_DVN by P.T. Wallace
// GoFA equivalent: gofa.Pn (convert p-vector to modulus and unit vector)
//
// Parameters:
//   - v: Input 3-component vector
//
// Returns:
//   - unit: Unit vector in the same direction as v (zero if v has zero magnitude)
//   - modulus: The magnitude of v
func NormalizeVector(v Vec3) (unit Vec3, modulus float64) {
	var unitVec [3]float64
	gofa.Pn(v, &modulus, &unitVec)
	return Vec3(unitVec), modulus
}

// Dvn is a SLALIB-compatible wrapper for NormalizeVector (sla_DVN)
func Dvn(v Vec3) (Vec3, float64) {
	return NormalizeVector(v)
}

// Magnitude calculates the magnitude (length) of a 3-component vector.
//
// This is a convenience function that returns just the magnitude without
// normalizing the vector. It's equivalent to sqrt(v·v).
//
// GoFA equivalent: gofa.Pm (modulus of p-vector)
//
// Parameters:
//   - v: Input 3-component vector
//
// Returns:
//   - The magnitude of v
func Magnitude(v Vec3) float64 {
	return gofa.Pm(v)
}

// VectorAdd adds two 3-component vectors.
//
// Original FORTRAN: sla_DVP (vector plus vector)
// GoFA equivalent: gofa.Ppp (p-vector plus p-vector)
//
// Parameters:
//   - a: First vector
//   - b: Second vector
//
// Returns:
//   - a + b
func VectorAdd(a, b Vec3) Vec3 {
	var result [3]float64
	gofa.Ppp(a, b, &result)
	return Vec3(result)
}


// VectorSubtract subtracts two 3-component vectors.
//
// Original FORTRAN: sla_DVSB (vector subtract vector)
// GoFA equivalent: gofa.Pmp (p-vector minus p-vector)
//
// Parameters:
//   - a: First vector
//   - b: Second vector
//
// Returns:
//   - a - b
func VectorSubtract(a, b Vec3) Vec3 {
	var result [3]float64
	gofa.Pmp(a, b, &result)
	return Vec3(result)
}


// ScalarMultiply multiplies a 3-component vector by a scalar.
//
// Original FORTRAN: sla_DSX (scalar times vector)
// GoFA equivalent: gofa.Sxp (multiply p-vector by scalar)
//
// Parameters:
//   - s: Scalar multiplier
//   - v: Vector
//
// Returns:
//   - s * v
func ScalarMultiply(s float64, v Vec3) Vec3 {
	var result [3]float64
	gofa.Sxp(s, v, &result)
	return Vec3(result)
}


// Single-precision (float32) vector operations

// DotProduct32 computes the scalar product of two 3-element vectors (single precision).
//
// This is the single-precision version of DotProduct.
//
// Original FORTRAN: sla_VDV by P.T. Wallace
//
// Parameters:
//   - a: First 3-component vector
//   - b: Second 3-component vector
//
// Returns:
//   - The scalar product a · b
func DotProduct32(a, b Vec3_32) float32 {
	// Convert to float64 for GoFA
	a64 := Vec3{float64(a[0]), float64(a[1]), float64(a[2])}
	b64 := Vec3{float64(b[0]), float64(b[1]), float64(b[2])}
	return float32(gofa.Pdp(a64, b64))
}

// Vdv is a SLALIB-compatible alias for DotProduct32 (sla_VDV)
var Vdv = DotProduct32

// CrossProduct32 computes the vector product of two 3-element vectors (single precision).
//
// This is the single-precision version of CrossProduct.
//
// Original FORTRAN: sla_VXV by P.T. Wallace
//
// Parameters:
//   - a: First 3-component vector
//   - b: Second 3-component vector
//
// Returns:
//   - The vector product a × b
func CrossProduct32(a, b Vec3_32) Vec3_32 {
	// Convert to float64 for GoFA
	a64 := [3]float64{float64(a[0]), float64(a[1]), float64(a[2])}
	b64 := [3]float64{float64(b[0]), float64(b[1]), float64(b[2])}
	var result64 [3]float64
	gofa.Pxp(a64, b64, &result64)
	return Vec3_32{float32(result64[0]), float32(result64[1]), float32(result64[2])}
}

// Vxv is a SLALIB-compatible alias for CrossProduct32 (sla_VXV)
var Vxv = CrossProduct32

// NormalizeVector32 normalizes a 3-component vector and returns its modulus (single precision).
//
// This is the single-precision version of NormalizeVector.
//
// Original FORTRAN: sla_VN by P.T. Wallace
//
// Parameters:
//   - v: Input 3-component vector
//
// Returns:
//   - unit: Unit vector in the same direction as v (zero if v has zero magnitude)
//   - modulus: The magnitude of v
func NormalizeVector32(v Vec3_32) (unit Vec3_32, modulus float32) {
	// Convert to float64 for GoFA
	v64 := [3]float64{float64(v[0]), float64(v[1]), float64(v[2])}
	var unit64 [3]float64
	var mod64 float64
	gofa.Pn(v64, &mod64, &unit64)
	return Vec3_32{float32(unit64[0]), float32(unit64[1]), float32(unit64[2])}, float32(mod64)
}

// Vn is a SLALIB-compatible wrapper for NormalizeVector32 (sla_VN)
func Vn(v Vec3_32) (Vec3_32, float32) {
	return NormalizeVector32(v)
}

// Magnitude32 calculates the magnitude of a 3-component vector (single precision).
//
// This is the single-precision version of Magnitude.
//
// Parameters:
//   - v: Input 3-component vector
//
// Returns:
//   - The magnitude of v
func Magnitude32(v Vec3_32) float32 {
	return float32(math.Sqrt(float64(v[0]*v[0] + v[1]*v[1] + v[2]*v[2])))
}

// SLALIB-compatible lowercase aliases

// dvn is a SLALIB-compatible alias for Dvn (sla_DVN)
var dvn = Dvn

// vn is a SLALIB-compatible alias for Vn (sla_VN)
var vn = Vn
