package slagofa

import "github.com/hebl/gofa"

// Matrix Operations
//
// This file implements SLALIB-compatible matrix operations using GoFA.
// All functions maintain SLALIB API compatibility while using IAU standards.

// MatrixMultiply performs 3x3 matrix multiplication: result = a × b
//
// Original FORTRAN: sla_DMXM by P.T. Wallace
// GoFA equivalent: gofa.Rxr
//
// Parameters:
//   - a: First matrix (left operand)
//   - b: Second matrix (right operand)
//
// Returns:
//   - Product matrix a × b
//
// Note: Matrix multiplication is not commutative (a×b ≠ b×a)
func MatrixMultiply(a, b Mat3) Mat3 {
	var result [3][3]float64
	gofa.Rxr(a, b, &result)
	return Mat3(result)
}

// Dmxm is a SLALIB-compatible alias for MatrixMultiply (sla_DMXM)
var Dmxm = MatrixMultiply

// MatrixMultiply32 performs 3x3 matrix multiplication (single precision)
//
// Original FORTRAN: sla_MXM by P.T. Wallace
// GoFA equivalent: gofa.Rxr (with type conversion)
func MatrixMultiply32(a, b Mat3_32) Mat3_32 {
	// Convert to float64
	a64 := Mat3{
		{float64(a[0][0]), float64(a[0][1]), float64(a[0][2])},
		{float64(a[1][0]), float64(a[1][1]), float64(a[1][2])},
		{float64(a[2][0]), float64(a[2][1]), float64(a[2][2])},
	}
	b64 := Mat3{
		{float64(b[0][0]), float64(b[0][1]), float64(b[0][2])},
		{float64(b[1][0]), float64(b[1][1]), float64(b[1][2])},
		{float64(b[2][0]), float64(b[2][1]), float64(b[2][2])},
	}

	// Compute using MatrixMultiply
	result64 := MatrixMultiply(a64, b64)

	// Convert back to float32
	return Mat3_32{
		{float32(result64[0][0]), float32(result64[0][1]), float32(result64[0][2])},
		{float32(result64[1][0]), float32(result64[1][1]), float32(result64[1][2])},
		{float32(result64[2][0]), float32(result64[2][1]), float32(result64[2][2])},
	}
}

// Mxm is a SLALIB-compatible alias for MatrixMultiply32 (sla_MXM)
var Mxm = MatrixMultiply32

// MatrixVectorMultiply performs matrix-vector multiplication: result = m × v
//
// Original FORTRAN: sla_DMXV by P.T. Wallace
// GoFA equivalent: gofa.Rxp
//
// Parameters:
//   - m: 3×3 matrix
//   - v: 3-element vector
//
// Returns:
//   - Product vector m × v
func MatrixVectorMultiply(m Mat3, v Vec3) Vec3 {
	var result [3]float64
	gofa.Rxp(m, v, &result)
	return Vec3(result)
}

// Dmxv is a SLALIB-compatible alias for MatrixVectorMultiply (sla_DMXV)
var Dmxv = MatrixVectorMultiply

// MatrixVectorMultiply32 performs matrix-vector multiplication (single precision)
//
// Original FORTRAN: sla_MXV by P.T. Wallace
// GoFA equivalent: gofa.Rxp (with type conversion)
func MatrixVectorMultiply32(m Mat3_32, v Vec3_32) Vec3_32 {
	// Convert to float64
	m64 := Mat3{
		{float64(m[0][0]), float64(m[0][1]), float64(m[0][2])},
		{float64(m[1][0]), float64(m[1][1]), float64(m[1][2])},
		{float64(m[2][0]), float64(m[2][1]), float64(m[2][2])},
	}
	v64 := Vec3{float64(v[0]), float64(v[1]), float64(v[2])}

	// Compute
	result64 := MatrixVectorMultiply(m64, v64)

	// Convert back to float32
	return Vec3_32{float32(result64[0]), float32(result64[1]), float32(result64[2])}
}

// Mxv is a SLALIB-compatible alias for MatrixVectorMultiply32 (sla_MXV)
var Mxv = MatrixVectorMultiply32

// InverseMatrixVectorMultiply performs inverse matrix-vector multiplication: result = m^T × v
//
// For rotation matrices, the inverse equals the transpose.
//
// Original FORTRAN: sla_DIMXV by P.T. Wallace
// GoFA equivalent: gofa.Trxp
//
// Parameters:
//   - m: 3×3 rotation matrix
//   - v: 3-element vector
//
// Returns:
//   - Product vector m^T × v (equivalent to m^-1 × v for rotation matrices)
func InverseMatrixVectorMultiply(m Mat3, v Vec3) Vec3 {
	var result [3]float64
	gofa.Trxp(m, v, &result)
	return Vec3(result)
}

// Dimxv is a SLALIB-compatible alias for InverseMatrixVectorMultiply (sla_DIMXV)
var Dimxv = InverseMatrixVectorMultiply

// InverseMatrixVectorMultiply32 performs inverse matrix-vector multiplication (single precision)
//
// Original FORTRAN: sla_IMXV by P.T. Wallace
// GoFA equivalent: gofa.Trxp (with type conversion)
func InverseMatrixVectorMultiply32(m Mat3_32, v Vec3_32) Vec3_32 {
	// Convert to float64
	m64 := Mat3{
		{float64(m[0][0]), float64(m[0][1]), float64(m[0][2])},
		{float64(m[1][0]), float64(m[1][1]), float64(m[1][2])},
		{float64(m[2][0]), float64(m[2][1]), float64(m[2][2])},
	}
	v64 := Vec3{float64(v[0]), float64(v[1]), float64(v[2])}

	// Compute
	result64 := InverseMatrixVectorMultiply(m64, v64)

	// Convert back to float32
	return Vec3_32{float32(result64[0]), float32(result64[1]), float32(result64[2])}
}

// Imxv is a SLALIB-compatible alias for InverseMatrixVectorMultiply32 (sla_IMXV)
var Imxv = InverseMatrixVectorMultiply32

// MatrixToRotationVector converts a rotation matrix to a rotation vector
//
// The rotation vector represents a rotation as a 3-element vector where:
//   - Direction = axis of rotation (Euler axis)
//   - Magnitude = angle of rotation in radians
//
// Original FORTRAN: sla_DM2AV by P.T. Wallace
// GoFA equivalent: gofa.Rm2v
//
// Parameters:
//   - m: 3×3 rotation matrix
//
// Returns:
//   - Rotation vector (axis-angle representation)
func MatrixToRotationVector(m Mat3) Vec3 {
	var result [3]float64
	gofa.Rm2v(m, &result)
	return Vec3(result)
}

// Dm2av is a SLALIB-compatible alias for MatrixToRotationVector (sla_DM2AV)
var Dm2av = MatrixToRotationVector

// MatrixToRotationVector32 converts rotation matrix to rotation vector (single precision)
//
// Original FORTRAN: sla_M2AV by P.T. Wallace
// GoFA equivalent: gofa.Rm2v (with type conversion)
func MatrixToRotationVector32(m Mat3_32) Vec3_32 {
	// Convert to float64
	m64 := Mat3{
		{float64(m[0][0]), float64(m[0][1]), float64(m[0][2])},
		{float64(m[1][0]), float64(m[1][1]), float64(m[1][2])},
		{float64(m[2][0]), float64(m[2][1]), float64(m[2][2])},
	}

	// Compute
	result64 := MatrixToRotationVector(m64)

	// Convert back to float32
	return Vec3_32{float32(result64[0]), float32(result64[1]), float32(result64[2])}
}

// M2av is a SLALIB-compatible alias for MatrixToRotationVector32 (sla_M2AV)
var M2av = MatrixToRotationVector32

// RotationVectorToMatrix converts a rotation vector to a rotation matrix
//
// The rotation vector represents a rotation as a 3-element vector where:
//   - Direction = axis of rotation (Euler axis)
//   - Magnitude = angle of rotation in radians
//
// Original FORTRAN: sla_DAV2M by P.T. Wallace
// GoFA equivalent: gofa.Rv2m
//
// Parameters:
//   - w: Rotation vector (axis-angle representation)
//
// Returns:
//   - 3×3 rotation matrix
func RotationVectorToMatrix(w Vec3) Mat3 {
	var result [3][3]float64
	gofa.Rv2m(w, &result)
	return Mat3(result)
}

// Dav2m is a SLALIB-compatible alias for RotationVectorToMatrix (sla_DAV2M)
var Dav2m = RotationVectorToMatrix

// RotationVectorToMatrix32 converts rotation vector to rotation matrix (single precision)
//
// Original FORTRAN: sla_AV2M by P.T. Wallace
// GoFA equivalent: gofa.Rv2m (with type conversion)
func RotationVectorToMatrix32(w Vec3_32) Mat3_32 {
	// Convert to float64
	w64 := Vec3{float64(w[0]), float64(w[1]), float64(w[2])}

	// Compute
	result64 := RotationVectorToMatrix(w64)

	// Convert back to float32
	return Mat3_32{
		{float32(result64[0][0]), float32(result64[0][1]), float32(result64[0][2])},
		{float32(result64[1][0]), float32(result64[1][1]), float32(result64[1][2])},
		{float32(result64[2][0]), float32(result64[2][1]), float32(result64[2][2])},
	}
}

// Av2m is a SLALIB-compatible alias for RotationVectorToMatrix32 (sla_AV2M)
var Av2m = RotationVectorToMatrix32

// EulerMatrix creates a rotation matrix from Euler angles
//
// This function forms a rotation matrix from a sequence of three rotations
// about specified axes. The order of rotations is given by the order string.
//
// Original FORTRAN: sla_DEULER by P.T. Wallace
// GoFA equivalent: gofa.Rx, gofa.Ry, gofa.Rz (combined)
//
// Parameters:
//   - order: String specifying rotation order (e.g., "ZYX", "XYZ", "ZXZ")
//            Each character must be 'X', 'Y', or 'Z'
//            Order is applied RIGHT TO LEFT (mathematical convention)
//   - phi: First Euler angle in radians
//   - theta: Second Euler angle in radians
//   - psi: Third Euler angle in radians
//
// Returns:
//   - 3×3 rotation matrix
//
// Example:
//   - EulerMatrix("ZYX", phi, theta, psi) means:
//     1. Rotate by psi about X
//     2. Then rotate by theta about Y
//     3. Then rotate by phi about Z
//
// Note: Rotations are applied in reverse order of the string to match
// SLALIB convention where the first character represents the outermost rotation.
func EulerMatrix(order string, phi, theta, psi float64) Mat3 {
	var result [3][3]float64

	// Start with identity matrix
	gofa.Ir(&result)

	// Angles array for mapping
	angles := []float64{phi, theta, psi}

	// Apply rotations in reverse order (right to left)
	for i := len(order) - 1; i >= 0; i-- {
		angle := angles[len(order)-1-i]

		switch order[i] {
		case 'X', 'x':
			gofa.Rx(angle, &result)
		case 'Y', 'y':
			gofa.Ry(angle, &result)
		case 'Z', 'z':
			gofa.Rz(angle, &result)
		}
	}

	return Mat3(result)
}

// Deuler is a SLALIB-compatible alias for EulerMatrix (sla_DEULER)
var Deuler = EulerMatrix

// EulerMatrix32 creates a rotation matrix from Euler angles (single precision)
//
// Original FORTRAN: sla_EULER by P.T. Wallace
// GoFA equivalent: gofa.Rx, gofa.Ry, gofa.Rz (with type conversion)
func EulerMatrix32(order string, phi, theta, psi float32) Mat3_32 {
	// Convert to float64 and compute
	result64 := EulerMatrix(order, float64(phi), float64(theta), float64(psi))

	// Convert back to float32
	return Mat3_32{
		{float32(result64[0][0]), float32(result64[0][1]), float32(result64[0][2])},
		{float32(result64[1][0]), float32(result64[1][1]), float32(result64[1][2])},
		{float32(result64[2][0]), float32(result64[2][1]), float32(result64[2][2])},
	}
}

// Euler is a SLALIB-compatible alias for EulerMatrix32 (sla_EULER)
var Euler = EulerMatrix32
