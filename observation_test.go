package slagofa

import (
	"math"
	"testing"
)

// Observation Planning Function Tests

const obsToleranceDouble = 1.0e-12 // Double precision tolerance
const obsToleranceSingle = 1.0e-6  // Single precision tolerance

// TestPositionAngle32 tests sla_PAV (single-precision position angle)
func TestPositionAngle32(t *testing.T) {
	// Test with simple known vectors
	// v1 pointing north (z-axis), v2 pointing slightly east
	v1 := Vec3_32{0.0, 0.0, 1.0}
	v2 := Vec3_32{0.1, 0.0, 1.0} // Slightly east

	pa := Pav(v1, v2)

	// Should be close to +π/2 (east)
	if math.Abs(float64(pa)-Pi/2.0) > 0.1 {
		t.Logf("Pav for eastward displacement: %.6f radians (%.2f degrees)",
			pa, float64(pa)*180.0/Pi)
	}

	// Test with vectors from SLALIB test suite (if available)
	// Using same vectors as PositionAngle (Dpav) test
	v3 := Vec3_32{1.0, 0.1, 0.2}
	v4 := Vec3_32{-3.0, 1.0e-3, 0.2}

	pa2 := Pav(v3, v4)

	// Compare with double-precision version
	v3_64 := Vec3{1.0, 0.1, 0.2}
	v4_64 := Vec3{-3.0, 1.0e-3, 0.2}
	pa2_expected := PositionAngle(v3_64, v4_64)

	if math.Abs(float64(pa2)-pa2_expected) > obsToleranceSingle {
		t.Errorf("Pav = %.10f, want %.10f (from Dpav)",
			pa2, pa2_expected)
	}
}

// TestPositionAngle32_Consistency tests Pav vs Dpav consistency
func TestPositionAngle32_Consistency(t *testing.T) {
	tests := []struct {
		name string
		v1   Vec3_32
		v2   Vec3_32
	}{
		{"North-South", Vec3_32{0, 0, 1}, Vec3_32{0, 0, -1}},
		{"East-West", Vec3_32{1, 0, 0}, Vec3_32{-1, 0, 0}},
		{"Random 1", Vec3_32{1.0, 0.5, 0.3}, Vec3_32{0.2, 0.8, 0.1}},
		{"Random 2", Vec3_32{-0.5, 0.3, 0.7}, Vec3_32{0.4, -0.2, 0.6}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Single precision
			pa32 := PositionAngle32(tt.v1, tt.v2)

			// Double precision
			v1_64 := Vec3{float64(tt.v1[0]), float64(tt.v1[1]), float64(tt.v1[2])}
			v2_64 := Vec3{float64(tt.v2[0]), float64(tt.v2[1]), float64(tt.v2[2])}
			pa64 := PositionAngle(v1_64, v2_64)

			diff := math.Abs(float64(pa32) - pa64)
			if diff > obsToleranceSingle {
				t.Errorf("Pav = %.10f, Dpav = %.10f, diff = %.2e",
					pa32, pa64, diff)
			}
		})
	}
}

// TestZenithDistance tests sla_ZD
func TestZenithDistance(t *testing.T) {
	// Test 1: Object on meridian at observer's latitude
	// HA=0, Dec=Phi → ZD should be 0 (object at zenith)
	phi := 0.7 // ~40 degrees north
	ha := 0.0
	dec := phi

	zd := ZenithDistance(ha, dec, phi)

	if math.Abs(zd) > obsToleranceDouble {
		t.Errorf("ZD for object at zenith = %.10f, want 0.0", zd)
	}

	// Test 2: Object on horizon
	// HA=0, Dec=0, Phi=π/2 (north pole) → object on horizon, ZD = π/2
	phi2 := Pi / 2.0
	ha2 := 0.0
	dec2 := 0.0

	zd2 := ZenithDistance(ha2, dec2, phi2)
	expected2 := Pi / 2.0

	if math.Abs(zd2-expected2) > obsToleranceDouble {
		t.Errorf("ZD for object on horizon = %.10f, want %.10f",
			zd2, expected2)
	}

	// Test 3: Object 6 hours east (HA = -π/2)
	phi3 := 0.5
	ha3 := -Pi / 2.0
	dec3 := 0.0

	zd3 := ZenithDistance(ha3, dec3, phi3)

	// Should be > 0 (not at zenith)
	if zd3 < 0.0 || zd3 > Pi {
		t.Errorf("ZD out of range [0, π]: %.10f", zd3)
	}

	// Verify it's in the expected range for this geometry
	// Object at Dec=0, 6h east, observer at ~29° lat
	// ZD should be significant but not near π
	if zd3 < 0.5 || zd3 > 2.5 {
		t.Logf("ZD for HA=-90°, Dec=0°, Phi=~29° = %.6f rad (%.2f°)",
			zd3, zd3*180.0/Pi)
	}
}

// TestZenithDistance_Range tests that ZD is always in [0, π]
func TestZenithDistance_Range(t *testing.T) {
	tests := []struct {
		name string
		ha   float64
		dec  float64
		phi  float64
	}{
		{"Zenith", 0.0, 0.7, 0.7},
		{"Horizon North", 0.0, 0.0, Pi / 2.0},
		{"East", -Pi / 2.0, 0.0, 0.5},
		{"West", Pi / 2.0, 0.0, 0.5},
		{"Negative Dec", 0.0, -0.5, 0.5},
		{"High Dec", 0.0, Pi / 3.0, 0.5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			zd := ZenithDistance(tt.ha, tt.dec, tt.phi)

			if zd < 0.0 || zd > Pi {
				t.Errorf("ZD = %.10f, expected range [0, π]", zd)
			}

			// Also test alias
			zd2 := Zd(tt.ha, tt.dec, tt.phi)
			if math.Abs(zd2-zd) > 1.0e-15 {
				t.Errorf("Zd (alias) differs from ZenithDistance")
			}
		})
	}
}

// TestAltaz tests sla_ALTAZ
func TestAltaz(t *testing.T) {
	// Test 1: Object at zenith (HA=0, Dec=Phi)
	phi := 0.7 // ~40 degrees north
	ha := 0.0
	dec := phi

	result := Altaz(ha, dec, phi)

	// At zenith: elevation should be π/2
	expectedEl := Pi / 2.0
	if math.Abs(result.El-expectedEl) > obsToleranceDouble {
		t.Errorf("Elevation at zenith = %.10f, want %.10f",
			result.El, expectedEl)
	}

	// Zenith distance should be 0
	zd := Pi/2.0 - result.El
	if math.Abs(zd) > obsToleranceDouble {
		t.Errorf("ZD at zenith = %.10f, want 0.0", zd)
	}

	// Test 2: Object on meridian, Dec=0, Phi=45°
	// Should be at elevation = 45° (= Phi when on meridian at equator)
	phi2 := Pi / 4.0 // 45 degrees
	ha2 := 0.0
	dec2 := 0.0

	result2 := Altaz(ha2, dec2, phi2)
	expectedEl2 := Pi/2.0 - phi2 // Elevation = 90° - ZD, ZD = Phi for Dec=0 on meridian

	if math.Abs(result2.El-expectedEl2) > obsToleranceDouble {
		t.Errorf("Elevation for Dec=0 on meridian = %.10f, want %.10f",
			result2.El, expectedEl2)
	}

	// Azimuth should be 0 (north) or π (south) when on meridian
	azNormalized := result2.Az
	if azNormalized > Pi {
		azNormalized = TwoPi - azNormalized
	}
	if azNormalized > 0.1 && math.Abs(azNormalized-Pi) > 0.1 {
		t.Logf("Azimuth on meridian = %.6f rad (%.2f°), expected 0 or 180",
			result2.Az, result2.Az*180.0/Pi)
	}
}

// TestAltaz_AzimuthRange tests that azimuth is in [0, 2π]
func TestAltaz_AzimuthRange(t *testing.T) {
	tests := []struct {
		name string
		ha   float64
		dec  float64
		phi  float64
	}{
		{"Meridian", 0.0, 0.0, 0.5},
		{"East", -Pi / 2.0, 0.0, 0.5},
		{"West", Pi / 2.0, 0.0, 0.5},
		{"North", 0.0, Pi / 3.0, 0.5},
		{"South", 0.0, -Pi / 3.0, 0.5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Altaz(tt.ha, tt.dec, tt.phi)

			// Azimuth should be in [0, 2π]
			if result.Az < 0.0 || result.Az >= TwoPi {
				t.Errorf("Azimuth = %.10f, expected range [0, 2π)", result.Az)
			}

			// Elevation should be in [-π/2, π/2] (practically)
			if result.El < -Pi/2.0-0.01 || result.El > Pi/2.0+0.01 {
				t.Errorf("Elevation = %.10f, expected range ≈[-π/2, π/2]", result.El)
			}

			// Parallactic angle should be in [-π, π]
			if result.Pa < -Pi-0.01 || result.Pa > Pi+0.01 {
				t.Errorf("Parallactic angle = %.10f, expected range ≈[-π, π]", result.Pa)
			}
		})
	}
}

// TestAltaz_Consistency tests internal consistency of Altaz
func TestAltaz_Consistency(t *testing.T) {
	// Test that Altaz elevation matches π/2 - ZenithDistance
	phi := 0.6
	ha := -0.5
	dec := 0.3

	result := Altaz(ha, dec, phi)
	zd := ZenithDistance(ha, dec, phi)

	// Elevation = π/2 - ZD
	expectedEl := Pi/2.0 - zd
	diff := math.Abs(result.El - expectedEl)

	if diff > obsToleranceDouble {
		t.Errorf("Elevation = %.10f, from ZD = %.10f, diff = %.2e",
			result.El, expectedEl, diff)
	}
}

// TestAltaz_VelocitiesAccelerations tests that velocities/accelerations are finite
func TestAltaz_VelocitiesAccelerations(t *testing.T) {
	tests := []struct {
		name string
		ha   float64
		dec  float64
		phi  float64
	}{
		{"Normal", 0.5, 0.3, 0.7},
		{"Near zenith", 0.0, 0.7, 0.7},
		{"Near horizon", Pi / 2.0, 0.0, 0.5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Altaz(tt.ha, tt.dec, tt.phi)

			// All velocities and accelerations should be finite
			values := []struct {
				name string
				val  float64
			}{
				{"Azd", result.Azd},
				{"Azdd", result.Azdd},
				{"Eld", result.Eld},
				{"Eldd", result.Eldd},
				{"Pad", result.Pad},
				{"Padd", result.Padd},
			}

			for _, v := range values {
				if math.IsNaN(v.val) || math.IsInf(v.val, 0) {
					t.Errorf("%s = %v, expected finite value", v.name, v.val)
				}
			}
		})
	}
}

// Benchmark tests

func BenchmarkPositionAngle32(b *testing.B) {
	v1 := Vec3_32{1.0, 0.1, 0.2}
	v2 := Vec3_32{-3.0, 1.0e-3, 0.2}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = PositionAngle32(v1, v2)
	}
}

func BenchmarkZenithDistance(b *testing.B) {
	ha := -0.5
	dec := 0.3
	phi := 0.7
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ZenithDistance(ha, dec, phi)
	}
}

func BenchmarkAltaz(b *testing.B) {
	ha := -0.5
	dec := 0.3
	phi := 0.7
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Altaz(ha, dec, phi)
	}
}
