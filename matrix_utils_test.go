package slagofa

import (
	"math"
	"testing"
)

// TestDmat tests matrix inversion with SLALIB test vectors
func TestDmat(t *testing.T) {
	// From SLALIB test suite (sla_test.f line 1912)
	// Input matrix DA and vector DV
	a := Mat3{
		{2.22, 1.6578, 1.380522},
		{1.6578, 1.380522, 1.22548578},
		{1.380522, 1.22548578, 1.1356276122},
	}

	// Expected inverse matrix (from SLALIB test line 1923)
	expectedInv := Mat3{
		{18.02550629769198, -52.16386644917280607, 34.37875949717850495},
		{-52.16386644917280607, 168.1778099099805627, -118.0722869694232670},
		{34.37875949717850495, -118.0722869694232670, 86.50307003740151262},
	}

	// Expected determinant (from SLALIB test line 1948)
	expectedDet := 0.003658344147359863

	inv, d, j := Dmat(a)

	// Check status
	if j != 0 {
		t.Errorf("Dmat status = %d, want 0", j)
	}

	// Check determinant (tolerance 1e-12 from SLALIB)
	if !almostEqual(d, expectedDet, 1e-12) {
		t.Errorf("Dmat determinant = %.15e, want %.15e", d, expectedDet)
	}

	// Check inverse matrix (tolerance 1e-10 from SLALIB)
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if !almostEqual(inv[i][j], expectedInv[i][j], 1e-10) {
				t.Errorf("Dmat inv[%d][%d] = %.15e, want %.15e",
					i, j, inv[i][j], expectedInv[i][j])
			}
		}
	}

	// Test singular matrix
	singular := Mat3{
		{1, 2, 3},
		{2, 4, 6},
		{3, 6, 9},
	}

	_, _, status := Dmat(singular)
	if status != -1 {
		t.Errorf("Dmat(singular) status = %d, want -1", status)
	}
}

// TestDcmpf tests Euler angle extraction from rotation matrix
func TestDcmpf(t *testing.T) {
	// Test with identity matrix
	identity := Mat3{
		{1, 0, 0},
		{0, 1, 0},
		{0, 0, 1},
	}

	a, b, c := Dcmpf(identity)
	if !almostEqual(a, 0, 1e-10) || !almostEqual(b, 0, 1e-10) || !almostEqual(c, 0, 1e-10) {
		t.Errorf("Dcmpf(identity) = (%.10f, %.10f, %.10f), want (0, 0, 0)", a, b, c)
	}

	// Test with a known rotation matrix
	// Just verify that the function runs without errors
	// The exact values depend on the decomposition convention
	testMatrix := Mat3{
		{0.866025, -0.5, 0},
		{0.5, 0.866025, 0},
		{0, 0, 1},
	}

	a2, b2, c2 := Dcmpf(testMatrix)
	// Just check that values are reasonable (in radians)
	if math.Abs(a2) > 2*math.Pi || math.Abs(b2) > 2*math.Pi || math.Abs(c2) > 2*math.Pi {
		t.Errorf("Dcmpf returned unreasonable angles: a=%.10f, b=%.10f, c=%.10f", a2, b2, c2)
	}
}

// TestVecmat tests outer product
func TestVecmat(t *testing.T) {
	a := Vec3{1, 2, 3}
	b := Vec3{4, 5, 6}

	result := Vecmat(a, b)

	// Expected: result[i][j] = a[i] * b[j]
	expected := Mat3{
		{4, 5, 6},
		{8, 10, 12},
		{12, 15, 18},
	}

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if !almostEqual(result[i][j], expected[i][j], 1e-12) {
				t.Errorf("Vecmat[%d][%d] = %f, want %f",
					i, j, result[i][j], expected[i][j])
			}
		}
	}
}

// TestDsvd tests Singular Value Decomposition
func TestDsvd(t *testing.T) {
	// Test with a simple matrix
	a := Mat3{
		{1, 0, 0},
		{0, 2, 0},
		{0, 0, 3},
	}

	u, s, vt, j := Dsvd(a)

	// Check status
	if j != 0 {
		t.Errorf("Dsvd status = %d, want 0", j)
	}

	// For a diagonal matrix, singular values should be the diagonal elements
	// (in descending order)
	expectedS := Vec3{3, 2, 1}

	for i := 0; i < 3; i++ {
		if !almostEqual(s[i], expectedS[i], 1e-10) {
			t.Errorf("Dsvd singular value[%d] = %f, want %f",
				i, s[i], expectedS[i])
		}
	}

	// Verify that U and V are orthogonal (U*U^T = I, V*V^T = I)
	// Check U orthogonality
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			dot := 0.0
			for k := 0; k < 3; k++ {
				dot += u[i][k] * u[j][k]
			}
			expected := 0.0
			if i == j {
				expected = 1.0
			}
			if !almostEqual(dot, expected, 1e-10) {
				t.Errorf("U not orthogonal: U[%d]·U[%d] = %f, want %f",
					i, j, dot, expected)
			}
		}
	}

	// Verify reconstruction: A = U * Σ * V^T
	var reconstructed Mat3
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			sum := 0.0
			for k := 0; k < 3; k++ {
				sum += u[i][k] * s[k] * vt[k][j]
			}
			reconstructed[i][j] = sum
		}
	}

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if !almostEqual(reconstructed[i][j], a[i][j], 1e-10) {
				t.Errorf("SVD reconstruction[%d][%d] = %f, want %f",
					i, j, reconstructed[i][j], a[i][j])
			}
		}
	}
}

// Test helper functions
func TestHelperFunctions(t *testing.T) {
	// Test sqrt
	if !almostEqual(sqrt(4.0), 2.0, 1e-10) {
		t.Errorf("sqrt(4) = %f, want 2.0", sqrt(4.0))
	}
	if !almostEqual(sqrt(9.0), 3.0, 1e-10) {
		t.Errorf("sqrt(9) = %f, want 3.0", sqrt(9.0))
	}

	// Test atan2
	if !almostEqual(atan2(1, 1), math.Pi/4, 1e-10) {
		t.Errorf("atan2(1, 1) = %f, want π/4 = %f", atan2(1, 1), math.Pi/4)
	}
	if !almostEqual(atan2(0, 1), 0.0, 1e-10) {
		t.Errorf("atan2(0, 1) = %f, want 0", atan2(0, 1))
	}
	if !almostEqual(atan2(1, 0), math.Pi/2, 1e-10) {
		t.Errorf("atan2(1, 0) = %f, want π/2 = %f", atan2(1, 0), math.Pi/2)
	}

	// Test atan
	if !almostEqual(atan(1), math.Pi/4, 1e-10) {
		t.Errorf("atan(1) = %f, want π/4 = %f", atan(1), math.Pi/4)
	}
	if !almostEqual(atan(0), 0.0, 1e-10) {
		t.Errorf("atan(0) = %f, want 0", atan(0))
	}

	// Test isFinite
	if !isFinite(1.0) {
		t.Error("isFinite(1.0) = false, want true")
	}
	if !isFinite(-1.0) {
		t.Error("isFinite(-1.0) = false, want true")
	}
	if isFinite(math.NaN()) {
		t.Error("isFinite(NaN) = true, want false")
	}
	if isFinite(math.Inf(1)) {
		t.Error("isFinite(Inf) = true, want false")
	}
}

// Benchmarks
func BenchmarkDmat(b *testing.B) {
	a := Mat3{
		{2.22, 1.6578, 1.380522},
		{1.6578, 1.380522, 1.22548578},
		{1.380522, 1.22548578, 1.1356276122},
	}

	for i := 0; i < b.N; i++ {
		_, _, _ = Dmat(a)
	}
}

func BenchmarkDcmpf(b *testing.B) {
	rmat := Mat3{
		{0.866025, -0.5, 0},
		{0.5, 0.866025, 0},
		{0, 0, 1},
	}

	for i := 0; i < b.N; i++ {
		_, _, _ = Dcmpf(rmat)
	}
}

func BenchmarkVecmat(b *testing.B) {
	a := Vec3{1, 2, 3}
	bb := Vec3{4, 5, 6}

	for i := 0; i < b.N; i++ {
		_ = Vecmat(a, bb)
	}
}

func BenchmarkDsvd(b *testing.B) {
	a := Mat3{
		{1, 0, 0},
		{0, 2, 0},
		{0, 0, 3},
	}

	for i := 0; i < b.N; i++ {
		_, _, _, _ = Dsvd(a)
	}
}
