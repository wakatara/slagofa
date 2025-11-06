package slagofa

import (
	"math"
	"testing"
)

func TestNormalizeAngle(t *testing.T) {
	tests := []struct {
		name     string
		input    float64
		expected float64
	}{
		{
			name:     "zero angle",
			input:    0.0,
			expected: 0.0,
		},
		{
			name:     "positive angle in range",
			input:    1.5,
			expected: 1.5,
		},
		{
			name:     "negative angle in range",
			input:    -1.5,
			expected: -1.5,
		},
		{
			name:     "from SLALIB test suite (sla_test.cc line 514)",
			input:    -4.0,
			expected: 2.283185307179586,
		},
		{
			name:     "just over pi",
			input:    Pi + 0.1,
			expected: -Pi + 0.1,
		},
		{
			name:     "large positive angle",
			input:    10 * TwoPi,
			expected: 0.0,
		},
		{
			name:     "large negative angle",
			input:    -10 * TwoPi,
			expected: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NormalizeAngle(tt.input)
			if !almostEqual(result, tt.expected, tolerance) {
				t.Errorf("NormalizeAngle(%.15f) = %.15f, want %.15f",
					tt.input, result, tt.expected)
			}

			// Also test SLALIB alias
			result2 := Drange(tt.input)
			if !almostEqual(result2, tt.expected, tolerance) {
				t.Errorf("Drange(%.15f) = %.15f, want %.15f",
					tt.input, result2, tt.expected)
			}
		})
	}
}

func TestNormalizeAnglePositive(t *testing.T) {
	tests := []struct {
		name     string
		input    float64
		expected float64
	}{
		{
			name:     "zero angle",
			input:    0.0,
			expected: 0.0,
		},
		{
			name:     "positive angle in range",
			input:    1.5,
			expected: 1.5,
		},
		{
			name:     "from SLALIB test suite (sla_test.cc line 520)",
			input:    -0.1,
			expected: 6.183185307179587,
		},
		{
			name:     "negative angle",
			input:    -1.5,
			expected: TwoPi - 1.5,
		},
		{
			name:     "large positive angle",
			input:    10*TwoPi + 1.5,
			expected: 1.5,
		},
		{
			name:     "large negative angle",
			input:    -10*TwoPi - 1.5,
			expected: TwoPi - 1.5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NormalizeAnglePositive(tt.input)
			if !almostEqual(result, tt.expected, tolerance) {
				t.Errorf("NormalizeAnglePositive(%.15f) = %.15f, want %.15f",
					tt.input, result, tt.expected)
			}

			// Also test SLALIB alias
			result2 := Dranrm(tt.input)
			if !almostEqual(result2, tt.expected, tolerance) {
				t.Errorf("Dranrm(%.15f) = %.15f, want %.15f",
					tt.input, result2, tt.expected)
			}
		})
	}
}

func TestNormalizeAngle32(t *testing.T) {
	tests := []struct {
		name     string
		input    float32
		expected float32
	}{
		{
			name:     "from SLALIB test suite",
			input:    -4.0,
			expected: 2.283185307179586,
		},
		{
			name:     "positive in range",
			input:    1.5,
			expected: 1.5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NormalizeAngle32(tt.input)
			if !almostEqual(float64(result), float64(tt.expected), tolerance32) {
				t.Errorf("NormalizeAngle32(%.15f) = %.15f, want %.15f",
					tt.input, result, tt.expected)
			}

			// Also test SLALIB alias
			result2 := Range(tt.input)
			if !almostEqual(float64(result2), float64(tt.expected), tolerance32) {
				t.Errorf("Range(%.15f) = %.15f, want %.15f",
					tt.input, result2, tt.expected)
			}
		})
	}
}

func TestNormalizeAnglePositive32(t *testing.T) {
	tests := []struct {
		name     string
		input    float32
		expected float32
	}{
		{
			name:     "from SLALIB test suite",
			input:    -0.1,
			expected: 6.183185307179587,
		},
		{
			name:     "positive in range",
			input:    1.5,
			expected: 1.5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NormalizeAnglePositive32(tt.input)
			if !almostEqual(float64(result), float64(tt.expected), 1.0e-5) {
				t.Errorf("NormalizeAnglePositive32(%.15f) = %.15f, want %.15f",
					tt.input, result, tt.expected)
			}

			// Also test SLALIB alias
			result2 := Ranorm(tt.input)
			if !almostEqual(float64(result2), float64(tt.expected), 1.0e-5) {
				t.Errorf("Ranorm(%.15f) = %.15f, want %.15f",
					tt.input, result2, tt.expected)
			}
		})
	}
}

func TestRadiansToAngle(t *testing.T) {
	// Test conversion of 45 degrees
	angle := 45.0 * DegreesToRadians
	sign, deg, arcmin, arcsec, _ := RadiansToAngle(4, angle)

	if sign != '+' {
		t.Errorf("RadiansToAngle sign = %c, want +", sign)
	}
	if deg != 45 {
		t.Errorf("RadiansToAngle degrees = %d, want 45", deg)
	}
	if arcmin != 0 {
		t.Errorf("RadiansToAngle arcminutes = %d, want 0", arcmin)
	}
	if arcsec != 0 {
		t.Errorf("RadiansToAngle arcseconds = %d, want 0", arcsec)
	}

	// Test SLALIB alias
	sign2, idmsf := Dr2af(4, angle)
	if sign2 != '+' || idmsf[0] != 45 {
		t.Errorf("Dr2af failed: sign=%c, deg=%d", sign2, idmsf[0])
	}
}

func TestRadiansToTime(t *testing.T) {
	// Test conversion of 6 hours (π/2 radians)
	angle := Pi / 2.0
	sign, hours, min, sec, _ := RadiansToTime(4, angle)

	if sign != '+' {
		t.Errorf("RadiansToTime sign = %c, want +", sign)
	}
	if hours != 6 {
		t.Errorf("RadiansToTime hours = %d, want 6", hours)
	}
	if min != 0 {
		t.Errorf("RadiansToTime minutes = %d, want 0", min)
	}
	if sec != 0 {
		t.Errorf("RadiansToTime seconds = %d, want 0", sec)
	}

	// Test SLALIB alias
	sign2, ihmsf := Dr2tf(4, angle)
	if sign2 != '+' || ihmsf[0] != 6 {
		t.Errorf("Dr2tf failed: sign=%c, hours=%d", sign2, ihmsf[0])
	}
}

func TestAngleToRadians(t *testing.T) {
	// Test conversion of 45 degrees, 30 arcminutes, 15 arcseconds
	rad, status := AngleToRadians('+', 45, 30, 15.0)

	if status != 0 {
		t.Errorf("AngleToRadians status = %d, want 0", status)
	}

	expected := (45.0 + 30.0/60.0 + 15.0/3600.0) * DegreesToRadians
	if !almostEqual(rad, expected, tolerance) {
		t.Errorf("AngleToRadians = %.15f, want %.15f", rad, expected)
	}

	// Test SLALIB alias
	rad2, status2 := Daf2r('+', 45, 30, 15.0)
	if status2 != 0 || !almostEqual(rad2, expected, tolerance) {
		t.Errorf("Daf2r failed: rad=%.15f, status=%d", rad2, status2)
	}
}

func TestTimeToRadians(t *testing.T) {
	// Test conversion of 6 hours
	rad, status := TimeToRadians('+', 6, 0, 0.0)

	if status != 0 {
		t.Errorf("TimeToRadians status = %d, want 0", status)
	}

	expected := Pi / 2.0
	if !almostEqual(rad, expected, tolerance) {
		t.Errorf("TimeToRadians = %.15f, want %.15f", rad, expected)
	}

	// Test SLALIB alias
	rad2, status2 := Dtf2r('+', 6, 0, 0.0)
	if status2 != 0 || !almostEqual(rad2, expected, tolerance) {
		t.Errorf("Dtf2r failed: rad=%.15f, status=%d", rad2, status2)
	}
}

func TestAngularSeparation(t *testing.T) {
	// Test separation of two points
	// Point 1: (0, 0), Point 2: (π/2, 0) - should be π/2 apart
	sep := AngularSeparation(0, 0, Pi/2, 0)
	expected := Pi / 2

	if !almostEqual(sep, expected, tolerance) {
		t.Errorf("AngularSeparation = %.15f, want %.15f", sep, expected)
	}

	// Test SLALIB alias
	sep2 := Dsep(0, 0, Pi/2, 0)
	if !almostEqual(sep2, expected, tolerance) {
		t.Errorf("Dsep = %.15f, want %.15f", sep2, expected)
	}

	// Test identical points
	sep3 := AngularSeparation(1.5, 0.3, 1.5, 0.3)
	if !almostEqual(sep3, 0.0, tolerance) {
		t.Errorf("AngularSeparation (identical points) = %.15f, want 0.0", sep3)
	}
}

func TestAngularSeparationVec(t *testing.T) {
	// Test separation of perpendicular unit vectors
	v1 := Vec3{1, 0, 0}
	v2 := Vec3{0, 1, 0}
	sep := AngularSeparationVec(v1, v2)
	expected := Pi / 2

	if !almostEqual(sep, expected, tolerance) {
		t.Errorf("AngularSeparationVec = %.15f, want %.15f", sep, expected)
	}

	// Test SLALIB alias
	sep2 := Dsepv(v1, v2)
	if !almostEqual(sep2, expected, tolerance) {
		t.Errorf("Dsepv = %.15f, want %.15f", sep2, expected)
	}

	// Test parallel vectors
	v3 := Vec3{2, 0, 0}
	sep3 := AngularSeparationVec(v1, v3)
	if !almostEqual(sep3, 0.0, tolerance) {
		t.Errorf("AngularSeparationVec (parallel) = %.15f, want 0.0", sep3)
	}

	// Test anti-parallel vectors
	v4 := Vec3{-1, 0, 0}
	sep4 := AngularSeparationVec(v1, v4)
	if !almostEqual(sep4, Pi, tolerance) {
		t.Errorf("AngularSeparationVec (anti-parallel) = %.15f, want %.15f", sep4, Pi)
	}
}

func TestPositionAngle(t *testing.T) {
	// Test position angle of unit vectors
	v1 := Vec3{1, 0, 0}
	v2 := Vec3{0, 1, 0}
	pa := PositionAngle(v1, v2)

	// Just verify we get a result in the valid range
	if pa < -Pi || pa > Pi {
		t.Errorf("PositionAngle = %.15f, outside range [-π, π]", pa)
	}

	// Test SLALIB alias
	pa2 := Dpav(v1, v2)
	if !almostEqual(pa2, pa, tolerance) {
		t.Errorf("Dpav = %.15f, want %.15f", pa2, pa)
	}
}

func TestBearing(t *testing.T) {
	// Test bearing between two points on the celestial sphere
	// Point 1: (0, 0), Point 2: (0, π/4) - north of point 1
	bearing := Bearing(0, 0, 0, Pi/4)

	// Bearing should be approximately 0 (north)
	if !almostEqual(bearing, 0.0, 1.0e-10) {
		t.Errorf("Bearing (north) = %.15f, want 0.0", bearing)
	}

	// Test SLALIB alias
	bearing2 := Dbear(0, 0, 0, Pi/4)
	if !almostEqual(bearing2, 0.0, 1.0e-10) {
		t.Errorf("Dbear (north) = %.15f, want 0.0", bearing2)
	}

	// Test point to the east
	bearing3 := Bearing(0, 0, Pi/2, 0)
	// Should be approximately π/2 (east)
	if math.Abs(bearing3-Pi/2) > 0.1 {
		t.Errorf("Bearing (east) = %.15f, want ~%.15f", bearing3, Pi/2)
	}
}

// Benchmark tests
func BenchmarkNormalizeAngle(b *testing.B) {
	angle := -4.0
	for i := 0; i < b.N; i++ {
		_ = NormalizeAngle(angle)
	}
}

func BenchmarkNormalizeAnglePositive(b *testing.B) {
	angle := -0.1
	for i := 0; i < b.N; i++ {
		_ = NormalizeAnglePositive(angle)
	}
}

func BenchmarkAngularSeparation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = AngularSeparation(0, 0, 1.5, 0.3)
	}
}
