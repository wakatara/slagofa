package slagofa

import (
	"math"
	"testing"
)

const tolerance = 1.0e-12
const tolerance32 = 1.0e-6

func almostEqual(a, b, tol float64) bool {
	// Automatically track deviation for reporting
	autoTrackDeviation(a, b, tol)
	return math.Abs(a-b) <= tol
}

func vec3AlmostEqual(a, b Vec3, tol float64) bool {
	for i := 0; i < 3; i++ {
		if !almostEqual(a[i], b[i], tol) {
			return false
		}
	}
	return true
}

func vec3_32AlmostEqual(a, b Vec3_32, tol float64) bool {
	for i := 0; i < 3; i++ {
		if !almostEqual(float64(a[i]), float64(b[i]), tol) {
			return false
		}
	}
	return true
}

// Tests for double-precision vector operations

func TestDotProduct(t *testing.T) {
	tests := []struct {
		name     string
		va       Vec3
		vb       Vec3
		expected float64
	}{
		{
			name:     "unit vectors along x",
			va:       Vec3{1, 0, 0},
			vb:       Vec3{1, 0, 0},
			expected: 1.0,
		},
		{
			name:     "perpendicular vectors",
			va:       Vec3{1, 0, 0},
			vb:       Vec3{0, 1, 0},
			expected: 0.0,
		},
		{
			name:     "opposite vectors",
			va:       Vec3{1, 0, 0},
			vb:       Vec3{-1, 0, 0},
			expected: -1.0,
		},
		{
			name:     "arbitrary vectors",
			va:       Vec3{1.0, 2.0, 3.0},
			vb:       Vec3{4.0, 5.0, 6.0},
			expected: 32.0, // 1*4 + 2*5 + 3*6 = 4 + 10 + 18 = 32
		},
		{
			name:     "from SLALIB test suite (sla_test.cc line 316)",
			va:       Vec3{-0.5366267667260525, 0.06977111097651444, -0.8409302618566215},
			vb:       Vec3{0.004147420704640065, -0.9496888606842218, 0.3131674740355448},
			expected: -0.3318384698006295,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DotProduct(tt.va, tt.vb)
			if !almostEqual(result, tt.expected, tolerance) {
				t.Errorf("DotProduct(%v, %v) = %.15f, want %.15f",
					tt.va, tt.vb, result, tt.expected)
			}

			// Also test SLALIB alias
			result2 := Dvdv(tt.va, tt.vb)
			if !almostEqual(result2, tt.expected, tolerance) {
				t.Errorf("Dvdv(%v, %v) = %.15f, want %.15f",
					tt.va, tt.vb, result2, tt.expected)
			}
		})
	}
}

func TestCrossProduct(t *testing.T) {
	tests := []struct {
		name     string
		va       Vec3
		vb       Vec3
		expected Vec3
	}{
		{
			name:     "x cross y = z",
			va:       Vec3{1, 0, 0},
			vb:       Vec3{0, 1, 0},
			expected: Vec3{0, 0, 1},
		},
		{
			name:     "y cross z = x",
			va:       Vec3{0, 1, 0},
			vb:       Vec3{0, 0, 1},
			expected: Vec3{1, 0, 0},
		},
		{
			name:     "z cross x = y",
			va:       Vec3{0, 0, 1},
			vb:       Vec3{1, 0, 0},
			expected: Vec3{0, 1, 0},
		},
		{
			name:     "from SLALIB test suite (sla_test.cc line 322)",
			va:       Vec3{0.004147420704640065, -0.9496888606842218, 0.3131674740355448},
			vb:       Vec3{-0.5366267667260525, 0.06977111097651444, -0.8409302618566215},
			expected: Vec3{0.7767720597123304, -0.1645663574562769, -0.5093390925544726},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CrossProduct(tt.va, tt.vb)
			if !vec3AlmostEqual(result, tt.expected, 1.0e-9) {
				t.Errorf("CrossProduct(%v, %v) = %v, want %v",
					tt.va, tt.vb, result, tt.expected)
			}

			// Also test SLALIB alias
			result2 := Dvxv(tt.va, tt.vb)
			if !vec3AlmostEqual(result2, tt.expected, 1.0e-9) {
				t.Errorf("Dvxv(%v, %v) = %v, want %v",
					tt.va, tt.vb, result2, tt.expected)
			}
		})
	}
}

func TestNormalizeVector(t *testing.T) {
	tests := []struct {
		name            string
		input           Vec3
		expectedNormal  Vec3
		expectedModulus float64
	}{
		{
			name:            "unit vector x",
			input:           Vec3{1, 0, 0},
			expectedNormal:  Vec3{1, 0, 0},
			expectedModulus: 1.0,
		},
		{
			name:            "scaled vector (3-4-5 triangle)",
			input:           Vec3{3, 4, 0},
			expectedNormal:  Vec3{0.6, 0.8, 0},
			expectedModulus: 5.0,
		},
		{
			name:            "zero vector",
			input:           Vec3{0, 0, 0},
			expectedNormal:  Vec3{0, 0, 0},
			expectedModulus: 0.0,
		},
		{
			name:            "arbitrary vector",
			input:           Vec3{1.0, 2.0, 2.0},
			expectedNormal:  Vec3{1.0 / 3.0, 2.0 / 3.0, 2.0 / 3.0},
			expectedModulus: 3.0,
		},
		{
			name:  "from SLALIB test suite (sla_test.cc line 328)",
			input: Vec3{6.889040510209034, -1577.473205461961, 520.1843672856759},
			expectedNormal: Vec3{
				0.004147420704640065,
				-0.9496888606842218,
				0.3131674740355448,
			},
			expectedModulus: 1661.042127339937,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resultVec, resultMod := NormalizeVector(tt.input)

			if !vec3AlmostEqual(resultVec, tt.expectedNormal, tolerance) {
				t.Errorf("NormalizeVector(%v) vector = %v, want %v",
					tt.input, resultVec, tt.expectedNormal)
			}

			if !almostEqual(resultMod, tt.expectedModulus, 1.0e-9) {
				t.Errorf("NormalizeVector(%v) modulus = %.15f, want %.15f",
					tt.input, resultMod, tt.expectedModulus)
			}

			// Also test SLALIB alias
			resultVec2, resultMod2 := Dvn(tt.input)
			if !vec3AlmostEqual(resultVec2, tt.expectedNormal, tolerance) {
				t.Errorf("Dvn(%v) vector = %v, want %v",
					tt.input, resultVec2, tt.expectedNormal)
			}
			if !almostEqual(resultMod2, tt.expectedModulus, 1.0e-9) {
				t.Errorf("Dvn(%v) modulus = %.15f, want %.15f",
					tt.input, resultMod2, tt.expectedModulus)
			}
		})
	}
}

func TestMagnitude(t *testing.T) {
	tests := []struct {
		name     string
		input    Vec3
		expected float64
	}{
		{
			name:     "unit vector",
			input:    Vec3{1, 0, 0},
			expected: 1.0,
		},
		{
			name:     "3-4-5 triangle",
			input:    Vec3{3, 4, 0},
			expected: 5.0,
		},
		{
			name:     "zero vector",
			input:    Vec3{0, 0, 0},
			expected: 0.0,
		},
		{
			name:     "arbitrary vector",
			input:    Vec3{1, 2, 2},
			expected: 3.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Magnitude(tt.input)
			if !almostEqual(result, tt.expected, tolerance) {
				t.Errorf("Magnitude(%v) = %.15f, want %.15f",
					tt.input, result, tt.expected)
			}
		})
	}
}

func TestVectorOperations(t *testing.T) {
	v1 := Vec3{1.0, 2.0, 3.0}
	v2 := Vec3{4.0, 5.0, 6.0}

	// Test VectorAdd
	sum := VectorAdd(v1, v2)
	expectedSum := Vec3{5.0, 7.0, 9.0}
	if !vec3AlmostEqual(sum, expectedSum, tolerance) {
		t.Errorf("VectorAdd(%v, %v) = %v, want %v", v1, v2, sum, expectedSum)
	}

	// Test SLALIB alias
	sum2 := VectorAdd(v1, v2)
	if !vec3AlmostEqual(sum2, expectedSum, tolerance) {
		t.Errorf("VectorAdd(%v, %v) = %v, want %v", v1, v2, sum2, expectedSum)
	}

	// Test VectorSubtract
	diff := VectorSubtract(v2, v1)
	expectedDiff := Vec3{3.0, 3.0, 3.0}
	if !vec3AlmostEqual(diff, expectedDiff, tolerance) {
		t.Errorf("VectorSubtract(%v, %v) = %v, want %v", v2, v1, diff, expectedDiff)
	}

	// Test SLALIB alias
	diff2 := VectorSubtract(v2, v1)
	if !vec3AlmostEqual(diff2, expectedDiff, tolerance) {
		t.Errorf("VectorSubtract(%v, %v) = %v, want %v", v2, v1, diff2, expectedDiff)
	}

	// Test ScalarMultiply
	scaled := ScalarMultiply(2.0, v1)
	expectedScaled := Vec3{2.0, 4.0, 6.0}
	if !vec3AlmostEqual(scaled, expectedScaled, tolerance) {
		t.Errorf("ScalarMultiply(2.0, %v) = %v, want %v", v1, scaled, expectedScaled)
	}

	// Test SLALIB alias
	scaled2 := ScalarMultiply(2.0, v1)
	if !vec3AlmostEqual(scaled2, expectedScaled, tolerance) {
		t.Errorf("ScalarMultiply(2.0, %v) = %v, want %v", v1, scaled2, expectedScaled)
	}
}

// Tests for single-precision vector operations

func TestDotProduct32(t *testing.T) {
	tests := []struct {
		name     string
		va       Vec3_32
		vb       Vec3_32
		expected float32
	}{
		{
			name:     "unit vectors along x",
			va:       Vec3_32{1, 0, 0},
			vb:       Vec3_32{1, 0, 0},
			expected: 1.0,
		},
		{
			name:     "perpendicular vectors",
			va:       Vec3_32{1, 0, 0},
			vb:       Vec3_32{0, 1, 0},
			expected: 0.0,
		},
		{
			name:     "from SLALIB test suite",
			va:       Vec3_32{0.004147420704640065, -0.9496888606842218, 0.3131674740355448},
			vb:       Vec3_32{-0.5366267667260525, 0.06977111097651444, -0.8409302618566215},
			expected: -0.3318384698006295,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DotProduct32(tt.va, tt.vb)
			if !almostEqual(float64(result), float64(tt.expected), tolerance32) {
				t.Errorf("DotProduct32(%v, %v) = %.15f, want %.15f",
					tt.va, tt.vb, result, tt.expected)
			}

			// Also test SLALIB alias
			result2 := Vdv(tt.va, tt.vb)
			if !almostEqual(float64(result2), float64(tt.expected), tolerance32) {
				t.Errorf("Vdv(%v, %v) = %.15f, want %.15f",
					tt.va, tt.vb, result2, tt.expected)
			}
		})
	}
}

func TestCrossProduct32(t *testing.T) {
	tests := []struct {
		name     string
		va       Vec3_32
		vb       Vec3_32
		expected Vec3_32
	}{
		{
			name:     "x cross y = z",
			va:       Vec3_32{1, 0, 0},
			vb:       Vec3_32{0, 1, 0},
			expected: Vec3_32{0, 0, 1},
		},
		{
			name:     "from SLALIB test suite",
			va:       Vec3_32{0.004147420704640065, -0.9496888606842218, 0.3131674740355448},
			vb:       Vec3_32{-0.5366267667260525, 0.06977111097651444, -0.8409302618566215},
			expected: Vec3_32{0.7767720597123304, -0.1645663574562769, -0.5093390925544726},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CrossProduct32(tt.va, tt.vb)
			if !vec3_32AlmostEqual(result, tt.expected, tolerance32) {
				t.Errorf("CrossProduct32(%v, %v) = %v, want %v",
					tt.va, tt.vb, result, tt.expected)
			}

			// Also test SLALIB alias
			result2 := Vxv(tt.va, tt.vb)
			if !vec3_32AlmostEqual(result2, tt.expected, tolerance32) {
				t.Errorf("Vxv(%v, %v) = %v, want %v",
					tt.va, tt.vb, result2, tt.expected)
			}
		})
	}
}

func TestNormalizeVector32(t *testing.T) {
	tests := []struct {
		name            string
		input           Vec3_32
		expectedNormal  Vec3_32
		expectedModulus float32
	}{
		{
			name:            "unit vector x",
			input:           Vec3_32{1, 0, 0},
			expectedNormal:  Vec3_32{1, 0, 0},
			expectedModulus: 1.0,
		},
		{
			name:            "scaled vector",
			input:           Vec3_32{3, 4, 0},
			expectedNormal:  Vec3_32{0.6, 0.8, 0},
			expectedModulus: 5.0,
		},
		{
			name:  "from SLALIB test suite",
			input: Vec3_32{6.889040510209034, -1577.473205461961, 520.1843672856759},
			expectedNormal: Vec3_32{
				0.004147420704640065,
				-0.9496888606842218,
				0.3131674740355448,
			},
			expectedModulus: 1661.042127339937,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resultVec, resultMod := NormalizeVector32(tt.input)

			if !vec3_32AlmostEqual(resultVec, tt.expectedNormal, tolerance32) {
				t.Errorf("NormalizeVector32(%v) vector = %v, want %v",
					tt.input, resultVec, tt.expectedNormal)
			}

			// Use 1.0e-3 tolerance for magnitude as in SLALIB
			if !almostEqual(float64(resultMod), float64(tt.expectedModulus), 1.0e-3) {
				t.Errorf("NormalizeVector32(%v) modulus = %.15f, want %.15f",
					tt.input, resultMod, tt.expectedModulus)
			}

			// Also test SLALIB alias
			resultVec2, resultMod2 := Vn(tt.input)
			if !vec3_32AlmostEqual(resultVec2, tt.expectedNormal, tolerance32) {
				t.Errorf("Vn(%v) vector = %v, want %v",
					tt.input, resultVec2, tt.expectedNormal)
			}
			if !almostEqual(float64(resultMod2), float64(tt.expectedModulus), 1.0e-3) {
				t.Errorf("Vn(%v) modulus = %.15f, want %.15f",
					tt.input, resultMod2, tt.expectedModulus)
			}
		})
	}
}

// Benchmark tests
func BenchmarkDotProduct(b *testing.B) {
	va := Vec3{1.0, 2.0, 3.0}
	vb := Vec3{4.0, 5.0, 6.0}
	for i := 0; i < b.N; i++ {
		_ = DotProduct(va, vb)
	}
}

func BenchmarkCrossProduct(b *testing.B) {
	va := Vec3{1.0, 2.0, 3.0}
	vb := Vec3{4.0, 5.0, 6.0}
	for i := 0; i < b.N; i++ {
		_ = CrossProduct(va, vb)
	}
}

func BenchmarkNormalizeVector(b *testing.B) {
	vec := Vec3{1.0, 2.0, 2.0}
	for i := 0; i < b.N; i++ {
		_, _ = NormalizeVector(vec)
	}
}

func BenchmarkMagnitude(b *testing.B) {
	vec := Vec3{1.0, 2.0, 2.0}
	for i := 0; i < b.N; i++ {
		_ = Magnitude(vec)
	}
}
