package slagofa

import (
	"math"
	"testing"
)

// Test tolerance constants (from SLALIB test suite)
const (
	matrixTolerance    = 1.0e-12 // Double precision matrix elements
	matrixTolerance32  = 1.0e-6  // Single precision matrix elements
	magnitudeTolerance = 1.0e-9  // For vector magnitudes (less strict)
)

// Helper function to compare matrices with tolerance
func matricesAlmostEqual(a, b Mat3, tolerance float64) bool {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if math.Abs(a[i][j]-b[i][j]) > tolerance {
				return false
			}
		}
	}
	return true
}

// Helper function to compare matrices (float32) with tolerance
func matricesAlmostEqual32(a, b Mat3_32, tolerance float32) bool {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if math.Abs(float64(a[i][j]-b[i][j])) > float64(tolerance) {
				return false
			}
		}
	}
	return true
}

// TestRotationVectorToMatrix tests sla_DAV2M (axis vector to rotation matrix)
//
// Test vector from SLALIB test suite (sla_test.f line 6205)
// Input: DAV = [-0.123, 0.0987, 0.0654]
func TestRotationVectorToMatrix(t *testing.T) {
	// Input: rotation vector (axis-angle representation)
	axisVector := Vec3{-0.123, 0.0987, 0.0654}

	// Expected output from SLALIB test suite
	expected := Mat3{
		{0.9930075842721269, 0.05902743090199868, -0.1022335560329612},
		{-0.07113807138648245, 0.9903204657727545, -0.1191836812279541},
		{0.09420887631983825, 0.1256229973879967, 0.9875948309655174},
	}

	// Call function
	result := RotationVectorToMatrix(axisVector)

	// Verify each matrix element
	if !matricesAlmostEqual(result, expected, matrixTolerance) {
		t.Errorf("RotationVectorToMatrix failed")
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				diff := math.Abs(result[i][j] - expected[i][j])
				t.Errorf("  [%d][%d]: got %.15f, want %.15f (diff: %.2e)",
					i, j, result[i][j], expected[i][j], diff)
			}
		}
	}

	// Test SLALIB alias
	result2 := Dav2m(axisVector)
	if !matricesAlmostEqual(result2, expected, matrixTolerance) {
		t.Errorf("Dav2m (alias) produced different result")
	}
}

// TestRotationVectorToMatrix32 tests sla_AV2M (single precision)
func TestRotationVectorToMatrix32(t *testing.T) {
	axisVector := Vec3_32{-0.123, 0.0987, 0.0654}

	expected := Mat3_32{
		{0.9930076, 0.05902743, -0.1022336},
		{-0.07113807, 0.9903205, -0.1191837},
		{0.09420888, 0.1256230, 0.9875948},
	}

	result := RotationVectorToMatrix32(axisVector)

	if !matricesAlmostEqual32(result, expected, matrixTolerance32) {
		t.Errorf("RotationVectorToMatrix32 failed")
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				diff := math.Abs(float64(result[i][j] - expected[i][j]))
				t.Errorf("  [%d][%d]: got %.7f, want %.7f (diff: %.2e)",
					i, j, result[i][j], expected[i][j], diff)
			}
		}
	}

	// Test SLALIB alias
	result2 := Av2m(axisVector)
	if !matricesAlmostEqual32(result2, expected, matrixTolerance32) {
		t.Errorf("Av2m (alias) produced different result")
	}
}

// TestEulerMatrix tests sla_DEULER (Euler angles to rotation matrix)
//
// Test vector from SLALIB test suite (sla_test.f line 6218)
// Input: 'YZY', phi=2.345, theta=-0.333, psi=2.222
func TestEulerMatrix(t *testing.T) {
	// Input: Euler angles (order "YZY")
	order := "YZY"
	phi := 2.345
	theta := -0.333
	psi := 2.222

	// Expected output from SLALIB test suite
	expected := Mat3{
		{-0.1681574770810878, 0.1981362273264315, 0.9656423242187410},
		{-0.2285369373983370, 0.9450659587140423, -0.2337117924378156},
		{-0.9589024617479674, -0.2599853247796050, -0.1136384607117296},
	}

	// Call function
	result := EulerMatrix(order, phi, theta, psi)

	// Verify each matrix element
	if !matricesAlmostEqual(result, expected, matrixTolerance) {
		t.Errorf("EulerMatrix failed")
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				diff := math.Abs(result[i][j] - expected[i][j])
				t.Errorf("  [%d][%d]: got %.15f, want %.15f (diff: %.2e)",
					i, j, result[i][j], expected[i][j], diff)
			}
		}
	}

	// Test SLALIB alias
	result2 := Deuler(order, phi, theta, psi)
	if !matricesAlmostEqual(result2, expected, matrixTolerance) {
		t.Errorf("Deuler (alias) produced different result")
	}
}

// TestEulerMatrix32 tests sla_EULER (single precision)
func TestEulerMatrix32(t *testing.T) {
	order := "YZY"
	phi := float32(2.345)
	theta := float32(-0.333)
	psi := float32(2.222)

	expected := Mat3_32{
		{-0.1681575, 0.1981362, 0.9656423},
		{-0.2285369, 0.9450660, -0.2337118},
		{-0.9589025, -0.2599853, -0.1136385},
	}

	result := EulerMatrix32(order, phi, theta, psi)

	if !matricesAlmostEqual32(result, expected, matrixTolerance32) {
		t.Errorf("EulerMatrix32 failed")
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				diff := math.Abs(float64(result[i][j] - expected[i][j]))
				t.Errorf("  [%d][%d]: got %.7f, want %.7f (diff: %.2e)",
					i, j, result[i][j], expected[i][j], diff)
			}
		}
	}

	// Test SLALIB alias
	result2 := Euler(order, phi, theta, psi)
	if !matricesAlmostEqual32(result2, expected, matrixTolerance32) {
		t.Errorf("Euler (alias) produced different result")
	}
}

// TestMatrixMultiply tests sla_DMXM (matrix multiplication)
//
// Test vector from SLALIB test suite (sla_test.f line 6231)
// Multiplies the two matrices from previous tests: DRM2 × DRM1
func TestMatrixMultiply(t *testing.T) {
	// First matrix: from TestEulerMatrix
	drm2 := Mat3{
		{-0.1681574770810878, 0.1981362273264315, 0.9656423242187410},
		{-0.2285369373983370, 0.9450659587140423, -0.2337117924378156},
		{-0.9589024617479674, -0.2599853247796050, -0.1136384607117296},
	}

	// Second matrix: from TestRotationVectorToMatrix
	drm1 := Mat3{
		{0.9930075842721269, 0.05902743090199868, -0.1022335560329612},
		{-0.07113807138648245, 0.9903204657727545, -0.1191836812279541},
		{0.09420887631983825, 0.1256229973879967, 0.9875948309655174},
	}

	// Expected output: DRM2 × DRM1
	expected := Mat3{
		{-0.09010460088585805, 0.3075993402463796, 0.9472400998581048},
		{-0.3161868071070688, 0.8930686362478707, -0.3200848543149236},
		{-0.9444083141897035, -0.3283459407855694, 0.01678926022795169},
	}

	// Call function
	result := MatrixMultiply(drm2, drm1)

	// Verify each matrix element
	if !matricesAlmostEqual(result, expected, matrixTolerance) {
		t.Errorf("MatrixMultiply failed")
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				diff := math.Abs(result[i][j] - expected[i][j])
				t.Errorf("  [%d][%d]: got %.15f, want %.15f (diff: %.2e)",
					i, j, result[i][j], expected[i][j], diff)
			}
		}
	}

	// Test SLALIB alias
	result2 := Dmxm(drm2, drm1)
	if !matricesAlmostEqual(result2, expected, matrixTolerance) {
		t.Errorf("Dmxm (alias) produced different result")
	}
}

// TestMatrixMultiply32 tests sla_MXM (single precision)
func TestMatrixMultiply32(t *testing.T) {
	drm2 := Mat3_32{
		{-0.1681575, 0.1981362, 0.9656423},
		{-0.2285369, 0.9450660, -0.2337118},
		{-0.9589025, -0.2599853, -0.1136385},
	}

	drm1 := Mat3_32{
		{0.9930076, 0.05902743, -0.1022336},
		{-0.07113807, 0.9903205, -0.1191837},
		{0.09420888, 0.1256230, 0.9875948},
	}

	expected := Mat3_32{
		{-0.09010460, 0.3075993, 0.9472401},
		{-0.3161868, 0.8930686, -0.3200849},
		{-0.9444083, -0.3283459, 0.01678926},
	}

	result := MatrixMultiply32(drm2, drm1)

	if !matricesAlmostEqual32(result, expected, matrixTolerance32) {
		t.Errorf("MatrixMultiply32 failed")
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				diff := math.Abs(float64(result[i][j] - expected[i][j]))
				t.Errorf("  [%d][%d]: got %.7f, want %.7f (diff: %.2e)",
					i, j, result[i][j], expected[i][j], diff)
			}
		}
	}

	// Test SLALIB alias
	result2 := Mxm(drm2, drm1)
	if !matricesAlmostEqual32(result2, expected, matrixTolerance32) {
		t.Errorf("Mxm (alias) produced different result")
	}
}

// TestMatrixVectorMultiply tests sla_DMXV (matrix × vector)
//
// Test vector from SLALIB test suite (sla_test.f line 6250)
// Uses the combined rotation matrix (DRM) to rotate a vector
func TestMatrixVectorMultiply(t *testing.T) {
	// Matrix from TestMatrixMultiply result
	drm := Mat3{
		{-0.09010460088585805, 0.3075993402463796, 0.9472400998581048},
		{-0.3161868071070688, 0.8930686362478707, -0.3200848543149236},
		{-0.9444083141897035, -0.3283459407855694, 0.01678926022795169},
	}

	// Input vector (from SLALIB test: spherical coordinates converted to Cartesian)
	dv1 := Vec3{-0.5366267667260525, 0.06977111097651444, -0.8409302618566215}

	// Expected output: rotated vector
	expected := Vec3{-0.7267487768696160, 0.5011537352639822, 0.4697671220397141}

	// Call function
	result := MatrixVectorMultiply(drm, dv1)

	// Verify each element
	if !vec3AlmostEqual(result, expected, matrixTolerance) {
		t.Errorf("MatrixVectorMultiply failed")
		for i := 0; i < 3; i++ {
			diff := math.Abs(result[i] - expected[i])
			t.Errorf("  [%d]: got %.15f, want %.15f (diff: %.2e)",
				i, result[i], expected[i], diff)
		}
	}

	// Test SLALIB alias
	result2 := Dmxv(drm, dv1)
	if !vec3AlmostEqual(result2, expected, matrixTolerance) {
		t.Errorf("Dmxv (alias) produced different result")
	}
}

// TestMatrixVectorMultiply32 tests sla_MXV (single precision)
func TestMatrixVectorMultiply32(t *testing.T) {
	drm := Mat3_32{
		{-0.09010460, 0.3075993, 0.9472401},
		{-0.3161868, 0.8930686, -0.3200849},
		{-0.9444083, -0.3283459, 0.01678926},
	}

	dv1 := Vec3_32{-0.5366268, 0.06977111, -0.8409303}

	expected := Vec3_32{-0.7267488, 0.5011537, 0.4697671}

	result := MatrixVectorMultiply32(drm, dv1)

	if !vec3_32AlmostEqual(result, expected, matrixTolerance32) {
		t.Errorf("MatrixVectorMultiply32 failed")
		for i := 0; i < 3; i++ {
			diff := math.Abs(float64(result[i] - expected[i]))
			t.Errorf("  [%d]: got %.7f, want %.7f (diff: %.2e)",
				i, result[i], expected[i], diff)
		}
	}

	// Test SLALIB alias
	result2 := Mxv(drm, dv1)
	if !vec3_32AlmostEqual(result2, expected, matrixTolerance32) {
		t.Errorf("Mxv (alias) produced different result")
	}
}

// TestInverseMatrixVectorMultiply tests sla_DIMXV (inverse/transpose matrix × vector)
//
// Test vector from SLALIB test suite (sla_test.f line 6267)
// This should recover the original vector by applying inverse transformation
func TestInverseMatrixVectorMultiply(t *testing.T) {
	// Matrix from TestMatrixMultiply result
	drm := Mat3{
		{-0.09010460088585805, 0.3075993402463796, 0.9472400998581048},
		{-0.3161868071070688, 0.8930686362478707, -0.3200848543149236},
		{-0.9444083141897035, -0.3283459407855694, 0.01678926022795169},
	}

	// Input: rotated vector from TestMatrixVectorMultiply
	dv3 := Vec3{-0.7267487768696160, 0.5011537352639822, 0.4697671220397141}

	// Expected: should recover original vector (within tolerance)
	expected := Vec3{-0.5366267667260526, 0.06977111097651445, -0.8409302618566215}

	// Call function
	result := InverseMatrixVectorMultiply(drm, dv3)

	// Verify each element
	if !vec3AlmostEqual(result, expected, matrixTolerance) {
		t.Errorf("InverseMatrixVectorMultiply failed")
		for i := 0; i < 3; i++ {
			diff := math.Abs(result[i] - expected[i])
			t.Errorf("  [%d]: got %.15f, want %.15f (diff: %.2e)",
				i, result[i], expected[i], diff)
		}
	}

	// Test SLALIB alias
	result2 := Dimxv(drm, dv3)
	if !vec3AlmostEqual(result2, expected, matrixTolerance) {
		t.Errorf("Dimxv (alias) produced different result")
	}
}

// TestInverseMatrixVectorMultiply32 tests sla_IMXV (single precision)
func TestInverseMatrixVectorMultiply32(t *testing.T) {
	drm := Mat3_32{
		{-0.09010460, 0.3075993, 0.9472401},
		{-0.3161868, 0.8930686, -0.3200849},
		{-0.9444083, -0.3283459, 0.01678926},
	}

	dv3 := Vec3_32{-0.7267488, 0.5011537, 0.4697671}

	expected := Vec3_32{-0.5366268, 0.06977111, -0.8409303}

	result := InverseMatrixVectorMultiply32(drm, dv3)

	if !vec3_32AlmostEqual(result, expected, matrixTolerance32) {
		t.Errorf("InverseMatrixVectorMultiply32 failed")
		for i := 0; i < 3; i++ {
			diff := math.Abs(float64(result[i] - expected[i]))
			t.Errorf("  [%d]: got %.7f, want %.7f (diff: %.2e)",
				i, result[i], expected[i], diff)
		}
	}

	// Test SLALIB alias
	result2 := Imxv(drm, dv3)
	if !vec3_32AlmostEqual(result2, expected, matrixTolerance32) {
		t.Errorf("Imxv (alias) produced different result")
	}
}

// TestMatrixToRotationVector tests sla_DM2AV (rotation matrix to axis vector)
//
// Test vector from SLALIB test suite (sla_test.f line 6280)
// Converts the combined rotation matrix back to rotation vector form
func TestMatrixToRotationVector(t *testing.T) {
	// Matrix from TestMatrixMultiply result
	drm := Mat3{
		{-0.09010460088585805, 0.3075993402463796, 0.9472400998581048},
		{-0.3161868071070688, 0.8930686362478707, -0.3200848543149236},
		{-0.9444083141897035, -0.3283459407855694, 0.01678926022795169},
	}

	// Expected output: rotation vector (axis-angle representation)
	expected := Vec3{0.006889040510209034, -1.577473205461961, 0.5201843672856759}

	// Call function
	result := MatrixToRotationVector(drm)

	// Verify each element
	if !vec3AlmostEqual(result, expected, matrixTolerance) {
		t.Errorf("MatrixToRotationVector failed")
		for i := 0; i < 3; i++ {
			diff := math.Abs(result[i] - expected[i])
			t.Errorf("  [%d]: got %.15f, want %.15f (diff: %.2e)",
				i, result[i], expected[i], diff)
		}
	}

	// Test SLALIB alias
	result2 := Dm2av(drm)
	if !vec3AlmostEqual(result2, expected, matrixTolerance) {
		t.Errorf("Dm2av (alias) produced different result")
	}
}

// TestMatrixToRotationVector32 tests sla_M2AV (single precision)
func TestMatrixToRotationVector32(t *testing.T) {
	drm := Mat3_32{
		{-0.09010460, 0.3075993, 0.9472401},
		{-0.3161868, 0.8930686, -0.3200849},
		{-0.9444083, -0.3283459, 0.01678926},
	}

	expected := Vec3_32{0.006889041, -1.577473, 0.5201844}

	result := MatrixToRotationVector32(drm)

	if !vec3_32AlmostEqual(result, expected, matrixTolerance32) {
		t.Errorf("MatrixToRotationVector32 failed")
		for i := 0; i < 3; i++ {
			diff := math.Abs(float64(result[i] - expected[i]))
			t.Errorf("  [%d]: got %.7f, want %.7f (diff: %.2e)",
				i, result[i], expected[i], diff)
		}
	}

	// Test SLALIB alias
	result2 := M2av(drm)
	if !vec3_32AlmostEqual(result2, expected, matrixTolerance32) {
		t.Errorf("M2av (alias) produced different result")
	}
}

// TestRoundTripRotationVector tests rotation vector → matrix → rotation vector
//
// Verifies that converting a rotation vector to matrix and back produces
// the same rotation vector (within numerical precision)
func TestRoundTripRotationVector(t *testing.T) {
	original := Vec3{0.1, 0.2, 0.3}

	// Convert to matrix
	matrix := RotationVectorToMatrix(original)

	// Convert back to rotation vector
	recovered := MatrixToRotationVector(matrix)

	// Should recover original
	if !vec3AlmostEqual(original, recovered, matrixTolerance) {
		t.Errorf("Round-trip rotation vector conversion failed")
		for i := 0; i < 3; i++ {
			diff := math.Abs(original[i] - recovered[i])
			t.Errorf("  [%d]: original %.15f, recovered %.15f (diff: %.2e)",
				i, original[i], recovered[i], diff)
		}
	}
}

// TestMatrixMultiplicationAssociativity tests (A × B) × C = A × (B × C)
func TestMatrixMultiplicationAssociativity(t *testing.T) {
	// Create three test matrices
	a := EulerMatrix("XYZ", 0.1, 0.2, 0.3)
	b := EulerMatrix("ZYX", 0.4, 0.5, 0.6)
	c := RotationVectorToMatrix(Vec3{0.7, 0.8, 0.9})

	// Left-associative: (A × B) × C
	ab := MatrixMultiply(a, b)
	left := MatrixMultiply(ab, c)

	// Right-associative: A × (B × C)
	bc := MatrixMultiply(b, c)
	right := MatrixMultiply(a, bc)

	// Should be equal
	if !matricesAlmostEqual(left, right, matrixTolerance) {
		t.Errorf("Matrix multiplication associativity failed")
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				diff := math.Abs(left[i][j] - right[i][j])
				if diff > matrixTolerance {
					t.Errorf("  [%d][%d]: left %.15f, right %.15f (diff: %.2e)",
						i, j, left[i][j], right[i][j], diff)
				}
			}
		}
	}
}

// TestInverseMatrixProperty tests M × M^T = I (rotation matrices are orthogonal)
func TestInverseMatrixProperty(t *testing.T) {
	// Create a rotation matrix
	m := EulerMatrix("YZY", 1.0, 2.0, 3.0)

	// Apply matrix and then its transpose to a vector
	v := Vec3{1.0, 0.0, 0.0}

	// Forward: v' = M × v
	vPrime := MatrixVectorMultiply(m, v)

	// Reverse: v'' = M^T × v' (should recover v)
	vRecovered := InverseMatrixVectorMultiply(m, vPrime)

	// Should recover original vector
	if !vec3AlmostEqual(v, vRecovered, matrixTolerance) {
		t.Errorf("Inverse matrix property failed")
		for i := 0; i < 3; i++ {
			diff := math.Abs(v[i] - vRecovered[i])
			t.Errorf("  [%d]: original %.15f, recovered %.15f (diff: %.2e)",
				i, v[i], vRecovered[i], diff)
		}
	}
}

// Benchmark tests to verify zero allocations

func BenchmarkMatrixMultiply(b *testing.B) {
	m1 := EulerMatrix("XYZ", 0.1, 0.2, 0.3)
	m2 := EulerMatrix("ZYX", 0.4, 0.5, 0.6)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = MatrixMultiply(m1, m2)
	}
}

func BenchmarkMatrixVectorMultiply(b *testing.B) {
	m := EulerMatrix("XYZ", 0.1, 0.2, 0.3)
	v := Vec3{1.0, 2.0, 3.0}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = MatrixVectorMultiply(m, v)
	}
}

func BenchmarkInverseMatrixVectorMultiply(b *testing.B) {
	m := EulerMatrix("XYZ", 0.1, 0.2, 0.3)
	v := Vec3{1.0, 2.0, 3.0}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = InverseMatrixVectorMultiply(m, v)
	}
}

func BenchmarkRotationVectorToMatrix(b *testing.B) {
	w := Vec3{0.1, 0.2, 0.3}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = RotationVectorToMatrix(w)
	}
}

func BenchmarkMatrixToRotationVector(b *testing.B) {
	m := EulerMatrix("XYZ", 0.1, 0.2, 0.3)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = MatrixToRotationVector(m)
	}
}

func BenchmarkEulerMatrix(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = EulerMatrix("YZY", 2.345, -0.333, 2.222)
	}
}
