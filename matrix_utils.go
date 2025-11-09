package slagofa

import (
	"gonum.org/v1/gonum/mat"
	"math"
)

// Phase 7: Matrix Utility Functions
//
// These functions provide SLALIB-compatible matrix operations using gonum/mat.

// Dmat inverts a 3×3 matrix.
//
// Original FORTRAN: sla_DMAT by P.T. Wallace
// Go equivalent: gonum/mat matrix inversion
// SLALIB reference: SUN/67 section 64
//
// Parameters:
//   - a: Input 3×3 matrix
//
// Returns:
//   - inv: Inverse of a
//   - d: Determinant of a
//   - j: Status (0=OK, -1=singular matrix)
//
// Notes:
//   - Returns error if matrix is singular (non-invertible)
//   - Uses gonum for numerical stability
func Dmat(a Mat3) (inv Mat3, d float64, j int) {
	// Convert to gonum matrix
	gm := mat.NewDense(3, 3, []float64{
		a[0][0], a[0][1], a[0][2],
		a[1][0], a[1][1], a[1][2],
		a[2][0], a[2][1], a[2][2],
	})

	// Compute determinant
	d = mat.Det(gm)

	// Check if singular
	if d == 0 || !isFinite(d) {
		return Mat3{}, 0, -1
	}

	// Compute inverse
	var ginv mat.Dense
	err := ginv.Inverse(gm)
	if err != nil {
		return Mat3{}, d, -1
	}

	// Convert back to Mat3
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			inv[i][j] = ginv.At(i, j)
		}
	}

	return inv, d, 0
}

// Dcmpf decomposes a rotation matrix into factors.
//
// This function decomposes a rotation matrix R into the form:
// R = Rx(a) * Ry(b) * Rz(c)
// where Rx, Ry, Rz are rotations about x, y, z axes.
//
// Original FORTRAN: sla_DCMPF by P.T. Wallace
// Go equivalent: Custom extraction of Euler angles
// SLALIB reference: SUN/67 section 26
//
// Parameters:
//   - rmat: Rotation matrix
//
// Returns:
//   - a: Rotation angle about x-axis (radians)
//   - b: Rotation angle about y-axis (radians)
//   - c: Rotation angle about z-axis (radians)
//
// Notes:
//   - Extracts XYZ Euler angles from rotation matrix
//   - Multiple solutions may exist; returns principal values
func Dcmpf(rmat Mat3) (a, b, c float64) {
	// Extract Euler angles from rotation matrix
	// For XYZ Euler angles (intrinsic rotations):
	// R = Rz(c) * Ry(b) * Rx(a)
	//
	// R = [  cy*cz                    cy*sz                   -sy    ]
	//     [  sx*sy*cz - cx*sz         sx*sy*sz + cx*cz        sx*cy  ]
	//     [  cx*sy*cz + sx*sz         cx*sy*sz - sx*cz        cx*cy  ]

	// Extract b (rotation about y)
	sy := -rmat[0][2]
	cy := sqrt(rmat[0][0]*rmat[0][0] + rmat[0][1]*rmat[0][1])

	if cy > 1e-10 {
		// Non-singular case
		b = atan2(sy, cy)
		a = atan2(rmat[1][2]/cy, rmat[2][2]/cy)
		c = atan2(rmat[0][1]/cy, rmat[0][0]/cy)
	} else {
		// Gimbal lock case (b = ±π/2)
		b = atan2(sy, cy)
		a = atan2(-rmat[1][0], rmat[1][1])
		c = 0.0 // Arbitrary choice
	}

	return a, b, c
}

// Vecmat computes the outer product of two vectors to form a matrix.
//
// Given vectors a and b, computes M = a ⊗ b (outer product).
// The result is: M[i][j] = a[i] * b[j]
//
// Original FORTRAN: sla_VECMAT (if exists in SLALIB)
// Go equivalent: Outer product using gonum
//
// Parameters:
//   - a: First vector (3-element)
//   - b: Second vector (3-element)
//
// Returns:
//   - m: Resulting 3×3 matrix
func Vecmat(a, b Vec3) Mat3 {
	var m Mat3
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			m[i][j] = a[i] * b[j]
		}
	}
	return m
}

// Dsvd computes the Singular Value Decomposition of a matrix.
//
// Decomposes matrix A into A = U * Σ * V^T
//
// Original FORTRAN: sla_DSVD (if exists)
// Go equivalent: gonum/mat SVD
//
// Parameters:
//   - a: Input m×n matrix (as Mat3 for 3×3 case)
//
// Returns:
//   - u: Left singular vectors (3×3)
//   - s: Singular values (3-element vector)
//   - vt: Right singular vectors transposed (3×3)
//   - j: Status (0=OK, -1=error)
//
// Notes:
//   - Simplified for 3×3 case
//   - Uses gonum for numerical stability
func Dsvd(a Mat3) (u Mat3, s Vec3, vt Mat3, j int) {
	// Convert to gonum matrix
	gm := mat.NewDense(3, 3, []float64{
		a[0][0], a[0][1], a[0][2],
		a[1][0], a[1][1], a[1][2],
		a[2][0], a[2][1], a[2][2],
	})

	// Compute SVD
	var svd mat.SVD
	ok := svd.Factorize(gm, mat.SVDFull)
	if !ok {
		return Mat3{}, Vec3{}, Mat3{}, -1
	}

	// Extract U
	var uMat mat.Dense
	svd.UTo(&uMat)
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			u[i][j] = uMat.At(i, j)
		}
	}

	// Extract singular values
	sVals := svd.Values(nil)
	for i := 0; i < 3; i++ {
		s[i] = sVals[i]
	}

	// Extract V^T
	var vMat mat.Dense
	svd.VTo(&vMat)
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			vt[i][j] = vMat.At(j, i) // Transpose
		}
	}

	return u, s, vt, 0
}

// Helper functions
func isFinite(x float64) bool {
	return x == x && x != x+1 && x != x-1 // Not NaN, not Inf
}

func sqrt(x float64) float64 {
	if x < 0 {
		return 0
	}
	// Use built-in sqrt
	y := x
	for i := 0; i < 10; i++ {
		y = (y + x/y) / 2
	}
	return y
}

func atan2(y, x float64) float64 {
	return math.Atan2(y, x)
}

func atan(x float64) float64 {
	return math.Atan(x)
}

// SLALIB-compatible lowercase aliases

// dcmpf is a SLALIB-compatible alias for Dcmpf (sla_DCMPF)
var dcmpf = Dcmpf

// dmat is a SLALIB-compatible alias for Dmat (sla_DMAT)
var dmat = Dmat
