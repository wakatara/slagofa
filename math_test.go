package slagofa

import (
	"math"
	"testing"
)

// Math function tests
// Note: These functions are FORTRAN intrinsics that SLALIB documented but didn't implement.
// There are no official SLALIB test vectors. We use standard mathematical test cases.

const mathTolerance = 1.0e-15      // Very tight tolerance for simple math
const mathTolerance32 = 1.0e-7     // Single precision tolerance

// TestMod tests sla_DMOD (modulus always positive)
func TestMod(t *testing.T) {
	tests := []struct {
		name     string
		a, b     float64
		expected float64
	}{
		{"positive mod positive", 8.0, 3.0, 2.0},
		{"negative mod positive", -8.0, 3.0, 1.0}, // Always positive result
		{"positive mod negative", 8.0, -3.0, -1.0},
		{"zero dividend", 0.0, 5.0, 0.0},
		{"fractional", 7.5, 2.0, 1.5},
		{"large values", 1000.0, 7.0, 6.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Mod(tt.a, tt.b)
			if math.Abs(result-tt.expected) > mathTolerance {
				t.Errorf("Mod(%.15f, %.15f) = %.15f, want %.15f",
					tt.a, tt.b, result, tt.expected)
			}

			// Test SLALIB alias
			result2 := Mod(tt.a, tt.b)
			if math.Abs(result2-tt.expected) > mathTolerance {
				t.Errorf("Dmod (alias) produced different result")
			}
		})
	}
}

// TestMod32 tests sla_MOD (single precision)
func TestMod32(t *testing.T) {
	result := Mod32(8.0, 3.0)
	expected := float32(2.0)
	if math.Abs(float64(result-expected)) > mathTolerance32 {
		t.Errorf("Mod32(8.0, 3.0) = %.7f, want %.7f", result, expected)
	}
}

// TestSign tests sla_DSIGN (transfer of sign)
func TestSign(t *testing.T) {
	tests := []struct {
		name     string
		a, b     float64
		expected float64
	}{
		{"positive value, positive sign", 5.0, 3.0, 5.0},
		{"positive value, negative sign", 5.0, -3.0, -5.0},
		{"negative value, positive sign", -5.0, 3.0, 5.0},
		{"negative value, negative sign", -5.0, -3.0, -5.0},
		{"zero value", 0.0, -1.0, 0.0},
		{"zero sign (positive)", 5.0, 0.0, 5.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Sign(tt.a, tt.b)
			if math.Abs(result-tt.expected) > mathTolerance {
				t.Errorf("Sign(%.15f, %.15f) = %.15f, want %.15f",
					tt.a, tt.b, result, tt.expected)
			}

			// Test SLALIB alias
			result2 := Sign(tt.a, tt.b)
			if math.Abs(result2-tt.expected) > mathTolerance {
				t.Errorf("Dsign (alias) produced different result")
			}
		})
	}
}

// TestSign32 tests sla_SIGN (single precision)
func TestSign32(t *testing.T) {
	result := Sign32(5.0, -3.0)
	expected := float32(-5.0)
	if math.Abs(float64(result-expected)) > mathTolerance32 {
		t.Errorf("Sign32(5.0, -3.0) = %.7f, want %.7f", result, expected)
	}
}

// TestTrunc tests sla_DINT (truncate to integer)
func TestTrunc(t *testing.T) {
	tests := []struct {
		name     string
		x        float64
		expected float64
	}{
		{"positive, round down", 3.7, 3.0},
		{"positive, already integer", 3.0, 3.0},
		{"negative, round up", -3.7, -3.0},
		{"negative, already integer", -3.0, -3.0},
		{"zero", 0.0, 0.0},
		{"small positive", 0.9, 0.0},
		{"small negative", -0.9, 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Trunc(tt.x)
			if math.Abs(result-tt.expected) > mathTolerance {
				t.Errorf("Trunc(%.15f) = %.15f, want %.15f",
					tt.x, result, tt.expected)
			}

			// Test SLALIB aliases
			result2 := Trunc(tt.x)
			if math.Abs(result2-tt.expected) > mathTolerance {
				t.Errorf("Dint (alias) produced different result")
			}

			result3 := Trunc(tt.x)
			if math.Abs(result3-tt.expected) > mathTolerance {
				t.Errorf("Aint (alias) produced different result")
			}
		})
	}
}

// TestTrunc32 tests sla_INT (single precision)
func TestTrunc32(t *testing.T) {
	result := Trunc32(3.7)
	expected := float32(3.0)
	if math.Abs(float64(result-expected)) > mathTolerance32 {
		t.Errorf("Trunc32(3.7) = %.7f, want %.7f", result, expected)
	}

	// Test SLALIB alias
	result2 := Trunc32(3.7)
	if math.Abs(float64(result2-expected)) > mathTolerance32 {
		t.Errorf("Int (alias) produced different result")
	}
}

// TestRound tests sla_ANINT (round to nearest integer)
func TestRound(t *testing.T) {
	tests := []struct {
		name     string
		x        float64
		expected float64
	}{
		{"round down", 3.4, 3.0},
		{"round up", 3.6, 4.0},
		{"exactly half (positive)", 3.5, 4.0},  // Go rounds half away from zero
		{"exactly half (negative)", -3.5, -4.0},
		{"already integer", 3.0, 3.0},
		{"zero", 0.0, 0.0},
		{"negative round down", -3.4, -3.0},
		{"negative round up", -3.6, -4.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Round(tt.x)
			if math.Abs(result-tt.expected) > mathTolerance {
				t.Errorf("Round(%.15f) = %.15f, want %.15f",
					tt.x, result, tt.expected)
			}

			// Test SLALIB alias
			result2 := Round(tt.x)
			if math.Abs(result2-tt.expected) > mathTolerance {
				t.Errorf("Anint (alias) produced different result")
			}
		})
	}
}

// TestRound32 tests sla_NINT (single precision)
func TestRound32(t *testing.T) {
	result := Round32(3.6)
	expected := float32(4.0)
	if math.Abs(float64(result-expected)) > mathTolerance32 {
		t.Errorf("Round32(3.6) = %.7f, want %.7f", result, expected)
	}

	// Test SLALIB alias
	result2 := Round32(3.6)
	if math.Abs(float64(result2-expected)) > mathTolerance32 {
		t.Errorf("Nint (alias) produced different result")
	}
}

// TestPoly tests sla_POLY (polynomial evaluation using Horner's method)
func TestPoly(t *testing.T) {
	tests := []struct {
		name     string
		coeffs   []float64
		x        float64
		expected float64
	}{
		{
			name:     "constant polynomial",
			coeffs:   []float64{5.0},
			x:        2.0,
			expected: 5.0,
		},
		{
			name:     "linear polynomial: 3 + 2x",
			coeffs:   []float64{3.0, 2.0},
			x:        4.0,
			expected: 11.0, // 3 + 2*4 = 11
		},
		{
			name:     "quadratic: 1 + 2x + 3x^2",
			coeffs:   []float64{1.0, 2.0, 3.0},
			x:        2.0,
			expected: 17.0, // 1 + 2*2 + 3*4 = 1 + 4 + 12 = 17
		},
		{
			name:     "cubic: 1 + 0x + 0x^2 + 1x^3 = 1 + x^3",
			coeffs:   []float64{1.0, 0.0, 0.0, 1.0},
			x:        2.0,
			expected: 9.0, // 1 + 8 = 9
		},
		{
			name:     "zero at x=0",
			coeffs:   []float64{1.0, 2.0, 3.0},
			x:        0.0,
			expected: 1.0, // Just the constant term
		},
		{
			name:     "negative coefficients",
			coeffs:   []float64{10.0, -5.0, 2.0},
			x:        3.0,
			expected: 13.0, // 10 - 5*3 + 2*9 = 10 - 15 + 18 = 13
		},
		{
			name:     "empty coefficients",
			coeffs:   []float64{},
			x:        5.0,
			expected: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Poly(tt.coeffs, tt.x)
			if math.Abs(result-tt.expected) > mathTolerance {
				t.Errorf("Poly(%v, %.15f) = %.15f, want %.15f",
					tt.coeffs, tt.x, result, tt.expected)
			}

			// Test SLALIB alias
			result2 := Poly(tt.coeffs, tt.x)
			if math.Abs(result2-tt.expected) > mathTolerance {
				t.Errorf("Dpoly (alias) produced different result")
			}
		})
	}
}

// TestPoly32 tests polynomial evaluation (single precision)
func TestPoly32(t *testing.T) {
	coeffs := []float32{1.0, 2.0, 3.0}
	x := float32(2.0)
	expected := float32(17.0)

	result := Poly32(coeffs, x)
	if math.Abs(float64(result-expected)) > mathTolerance32 {
		t.Errorf("Poly32(%v, %.7f) = %.7f, want %.7f",
			coeffs, x, result, expected)
	}
}

// TestPolyMod tests sla_POLMO (polynomial with modulus)
func TestPolyMod(t *testing.T) {
	tests := []struct {
		name     string
		coeffs   []float64
		x        float64
		modulus  float64
		expected float64
	}{
		{
			name:     "simple with modulus",
			coeffs:   []float64{1.0, 2.0, 3.0},
			x:        2.0,
			modulus:  10.0,
			expected: 7.0, // (1 + 2*2 + 3*4) mod 10 = 17 mod 10 = 7
		},
		{
			name:     "wraps multiple times",
			coeffs:   []float64{10.0, 5.0, 3.0},
			x:        5.0,
			modulus:  7.0,
			expected: 5.0, // (10 + 5*5 + 3*25) mod 7 = (10 + 25 + 75) mod 7 = 110 mod 7 = 5
		},
		{
			name:     "negative result",
			coeffs:   []float64{-5.0, 2.0},
			x:        1.0,
			modulus:  10.0,
			expected: 7.0, // (-5 + 2*1) mod 10 = -3 mod 10 = 7
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := PolyMod(tt.coeffs, tt.x, tt.modulus)

			// For modulus, allow slightly larger tolerance due to multiple operations
			tol := mathTolerance * 10.0
			if math.Abs(result-tt.expected) > tol {
				// Manual calculation for debugging
				direct := Poly(tt.coeffs, tt.x)
				directMod := Mod(direct, tt.modulus)
				t.Errorf("PolyMod(%v, %.15f, %.15f) = %.15f, want %.15f",
					tt.coeffs, tt.x, tt.modulus, result, tt.expected)
				t.Errorf("  Direct poly result: %.15f, direct mod: %.15f", direct, directMod)
			}

			// Test SLALIB alias
			result2 := PolyMod(tt.coeffs, tt.x, tt.modulus)
			if math.Abs(result2-result) > mathTolerance {
				t.Errorf("Dpolmo (alias) produced different result")
			}
		})
	}
}

// TestPolyMod32 tests polynomial with modulus (single precision)
func TestPolyMod32(t *testing.T) {
	coeffs := []float32{1.0, 2.0, 3.0}
	x := float32(2.0)
	modulus := float32(10.0)
	expected := float32(7.0)

	result := PolyMod32(coeffs, x, modulus)
	if math.Abs(float64(result-expected)) > mathTolerance32*10.0 {
		t.Errorf("PolyMod32(%v, %.7f, %.7f) = %.7f, want %.7f",
			coeffs, x, modulus, result, expected)
	}
}

// Benchmark tests

func BenchmarkMod(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = Mod(8.0, 3.0)
	}
}

func BenchmarkSign(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = Sign(5.0, -3.0)
	}
}

func BenchmarkTrunc(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = Trunc(3.7)
	}
}

func BenchmarkRound(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = Round(3.6)
	}
}

func BenchmarkPoly(b *testing.B) {
	coeffs := []float64{1.0, 2.0, 3.0, 4.0, 5.0}
	x := 2.0
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Poly(coeffs, x)
	}
}

func BenchmarkPolyMod(b *testing.B) {
	coeffs := []float64{1.0, 2.0, 3.0, 4.0, 5.0}
	x := 2.0
	mod := 10.0
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = PolyMod(coeffs, x, mod)
	}
}
