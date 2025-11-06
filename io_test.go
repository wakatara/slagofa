package slagofa

import (
	"math"
	"testing"
)

// TestDafin tests the angle parsing function with SLALIB test vectors
func TestDafin(t *testing.T) {
	// From SLALIB test suite (sla_test.f line 543)
	// Test string: "12 34 56.7 |"
	// Expected: 0.2196045986911432 radians at position 12, status 0
	s := "12 34 56.7 |"
	iptr := 0

	angle, j := Dafin(s, iptr)

	// Check status
	if j != 0 {
		t.Errorf("Dafin status = %d, want 0", j)
	}

	// Check result (tolerance 1e-12 from SLALIB test)
	expected := 0.2196045986911432
	if !almostEqual(angle, expected, 1e-12) {
		t.Errorf("Dafin = %.15e, want %.15e", angle, expected)
	}

	// Test various formats
	tests := []struct {
		input    string
		expected float64
	}{
		{"45.5", 45.5 * DegreesToRadians},
		{"-30.25", -30.25 * DegreesToRadians},
		{"12 30", (12.0 + 30.0/60.0) * DegreesToRadians},
		{"12 30 45", (12.0 + 30.0/60.0 + 45.0/3600.0) * DegreesToRadians},
		{"-12 30 45", -(12.0 + 30.0/60.0 + 45.0/3600.0) * DegreesToRadians},
		{"12:30:45", (12.0 + 30.0/60.0 + 45.0/3600.0) * DegreesToRadians},
	}

	for _, tt := range tests {
		result, status := Dafin(tt.input, 0)
		if status != 0 {
			t.Errorf("Dafin(%q) status = %d, want 0", tt.input, status)
		}
		if !almostEqual(result, tt.expected, 1e-10) {
			t.Errorf("Dafin(%q) = %.15e, want %.15e", tt.input, result, tt.expected)
		}
	}
}

// TestDfltin tests floating point parsing
func TestDfltin(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"123.456", 123.456},
		{"-987.654", -987.654},
		{"  42.0  ", 42.0},
		{"3.14159265", 3.14159265},
		{"-0.001", -0.001},
	}

	for _, tt := range tests {
		result, status := Dfltin(tt.input, 0)
		if status != 0 {
			t.Errorf("Dfltin(%q) status = %d, want 0", tt.input, status)
		}
		if !almostEqual(result, tt.expected, 1e-12) {
			t.Errorf("Dfltin(%q) = %.15e, want %.15e", tt.input, result, tt.expected)
		}
	}

	// Test error cases
	errorTests := []string{
		"",
		"   ",
		"not a number",
	}

	for _, input := range errorTests {
		_, status := Dfltin(input, 0)
		if status != 1 {
			t.Errorf("Dfltin(%q) status = %d, want 1 (error)", input, status)
		}
	}
}

// TestFlotin tests single-precision float parsing
func TestFlotin(t *testing.T) {
	tests := []struct {
		input    string
		expected float32
	}{
		{"123.456", 123.456},
		{"-987.654", -987.654},
		{"  42.0  ", 42.0},
		{"3.14159", 3.14159},
	}

	for _, tt := range tests {
		result, status := Flotin(tt.input, 0)
		if status != 0 {
			t.Errorf("Flotin(%q) status = %d, want 0", tt.input, status)
		}
		if math.Abs(float64(result-tt.expected)) > 1e-5 {
			t.Errorf("Flotin(%q) = %f, want %f", tt.input, result, tt.expected)
		}
	}
}

// TestIntin tests integer parsing
func TestIntin(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"123", 123},
		{"-456", -456},
		{"  789  ", 789},
		{"+42", 42},
		{"0", 0},
	}

	for _, tt := range tests {
		result, status := Intin(tt.input, 0)
		if status != 0 {
			t.Errorf("Intin(%q) status = %d, want 0", tt.input, status)
		}
		if result != tt.expected {
			t.Errorf("Intin(%q) = %d, want %d", tt.input, result, tt.expected)
		}
	}

	// Test error cases
	errorTests := []string{
		"",
		"   ",
		"not a number",
		"12.34", // Not an integer
	}

	for _, input := range errorTests {
		_, status := Intin(input, 0)
		if status != 1 {
			t.Errorf("Intin(%q) status = %d, want 1 (error)", input, status)
		}
	}
}

// TestObs tests observatory lookup
func TestObs(t *testing.T) {
	// Test lookup by number
	tests := []struct {
		n    int
		name string
	}{
		{82, "Mauna Kea"},
		{83, "Gemini-N"},
		{84, "CFHT"},
		{85, "Keck"},
		{86, "Subaru"},
		{24, "JCMT"},
		{10, "Palomar"},
	}

	for _, tt := range tests {
		w, p, h, err := Obs(tt.n, "")
		if err != nil {
			t.Errorf("Obs(%d) error: %v", tt.n, err)
			continue
		}

		// Check that values are reasonable
		if w < 0 || w > 2*math.Pi {
			t.Errorf("Obs(%d) longitude %f out of range [0, 2π]", tt.n, w)
		}
		if math.Abs(p) > math.Pi/2 {
			t.Errorf("Obs(%d) latitude %f out of range [-π/2, π/2]", tt.n, p)
		}
		if h < 0 || h > 10000 {
			t.Errorf("Obs(%d) height %f out of reasonable range", tt.n, h)
		}
	}

	// Test lookup by name
	nameTests := []string{
		"Mauna Kea",
		"KECK",
		"subaru",
		"Palomar",
	}

	for _, name := range nameTests {
		w, p, h, err := Obs(0, name)
		if err != nil {
			t.Errorf("Obs(0, %q) error: %v", name, err)
			continue
		}

		// Check that values are reasonable
		if w < 0 || w > 2*math.Pi {
			t.Errorf("Obs(0, %q) longitude %f out of range", name, w)
		}
		if math.Abs(p) > math.Pi/2 {
			t.Errorf("Obs(0, %q) latitude %f out of range", name, p)
		}
		if h < 0 || h > 10000 {
			t.Errorf("Obs(0, %q) height %f out of reasonable range", name, h)
		}
	}

	// Test error case - non-existent observatory
	_, _, _, err := Obs(999, "")
	if err == nil {
		t.Error("Obs(999) should return error for non-existent observatory")
	}

	_, _, _, err = Obs(0, "NonExistent")
	if err == nil {
		t.Error("Obs(0, 'NonExistent') should return error")
	}
}

// TestGresid tests Gregorian calendar residual
func TestGresid(t *testing.T) {
	// Test that residual is between 0 and 1
	tests := []struct {
		year, month, day int
	}{
		{2000, 1, 1},
		{2023, 6, 15},
		{1985, 12, 31},
	}

	for _, tt := range tests {
		residual := Gresid(tt.year, tt.month, tt.day)
		if residual < 0 || residual >= 1 {
			t.Errorf("Gresid(%d, %d, %d) = %f, want [0, 1)",
				tt.year, tt.month, tt.day, residual)
		}
	}
}

// Benchmarks
func BenchmarkDafin(b *testing.B) {
	s := "12 34 56.7"
	for i := 0; i < b.N; i++ {
		_, _ = Dafin(s, 0)
	}
}

func BenchmarkDfltin(b *testing.B) {
	s := "123.456789"
	for i := 0; i < b.N; i++ {
		_, _ = Dfltin(s, 0)
	}
}

func BenchmarkIntin(b *testing.B) {
	s := "12345"
	for i := 0; i < b.N; i++ {
		_, _ = Intin(s, 0)
	}
}

func BenchmarkObs(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _, _, _ = Obs(82, "")
	}
}
