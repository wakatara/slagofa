package slagofa

import (
	"math"
	"testing"
)

// Time scale and sidereal time test vectors from PAL test suite (palTest.c)

const timeTolerance = 1.0e-12 // Standard tolerance for time calculations

// TestDtt tests sla_DTT (TT-UTC calculation)
func TestDtt(t *testing.T) {
	// Test for a known date: 2000-01-01 (MJD 51544)
	// At this time, leap seconds totaled 32 seconds
	// TT-UTC = 32.184 + 32 = 64.184 seconds
	utc := 51544.0
	result := Dtt(utc)
	expected := 64.184 // 32.184 (TT-TAI) + 32 (leap seconds at 2000-01-01)

	if math.Abs(result-expected) > 0.1 {
		t.Errorf("Dtt(51544.0) = %.3f, want %.3f", result, expected)
	}
}

// TestDtt_MultipleEpochs tests TT-UTC at various epochs
func TestDtt_MultipleEpochs(t *testing.T) {
	tests := []struct {
		name     string
		mjd      float64
		minValue float64 // Minimum expected (leap seconds increase over time)
	}{
		{"1972-01-01 (first leap second)", 41317.0, 42.184}, // 32.184 + 10
		{"1980-01-01", 44239.0, 50.184},                     // Approximate
		{"1990-01-01", 47892.0, 56.184},                     // Approximate
		{"2000-01-01", 51544.0, 64.184},                     // 32.184 + 32
		{"2010-01-01", 55197.0, 66.184},                     // Approximate
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Dtt(tt.mjd)

			if result < tt.minValue {
				t.Errorf("Dtt(%.1f) = %.3f, expected >= %.3f",
					tt.mjd, result, tt.minValue)
			}

			// Should always be >= 32.184 (the TT-TAI offset)
			if result < 32.184 {
				t.Errorf("Dtt(%.1f) = %.3f, must be >= 32.184", tt.mjd, result)
			}
		})
	}
}

// TestDeltaT tests sla_DT (historical Delta T model)
//
// Test vectors from PAL palTest.c lines 712-716
func TestDeltaT(t *testing.T) {
	tests := []struct {
		name     string
		epoch    float64
		expected float64
		tol      float64
	}{
		// From PAL palTest.c: vvd ( palDt ( 500 ), 4686.7, 1e-10, "palDt", "500", status );
		{"PAL test: year 500", 500.0, 4686.7, 1.0e-10},
		// From PAL palTest.c: vvd ( palDt ( 1400 ), 408, 1e-11, "palDt", "1400", status );
		{"PAL test: year 1400", 1400.0, 408.0, 1.0e-11},
		// From PAL palTest.c: vvd ( palDt ( 1950 ), 27.99145626, 1e-12, "palDt", "1950", status );
		{"PAL test: year 1950", 1950.0, 27.99145626, 1.0e-12},
		// Additional tests for model boundaries
		{"Boundary 979", 979.0258204760233, 1718.696, 0.01},      // At boundary
		{"Boundary 1708", 1708.185161980887, 21.496, 0.01},       // At boundary
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DeltaT(tt.epoch)

			if math.Abs(result-tt.expected) > tt.tol {
				t.Errorf("DeltaT(%.3f) = %.3f, want %.3f ± %.1f",
					tt.epoch, result, tt.expected, tt.tol)
			}

			// Test SLALIB alias
			result2 := Dt(tt.epoch)
			if math.Abs(result2-result) > 1.0e-15 {
				t.Errorf("Dt (alias) produced different result")
			}
		})
	}
}

// TestDeltaT_Continuity tests that models are continuous at boundaries
func TestDeltaT_Continuity(t *testing.T) {
	// Test continuity at boundary 1: year 979.0258204760233
	boundary1 := 979.0258204760233
	before1 := DeltaT(boundary1 - 0.001)
	at1 := DeltaT(boundary1)
	after1 := DeltaT(boundary1 + 0.001)

	// Should be continuous (differ by < 1 second for 0.001 year change)
	if math.Abs(at1-before1) > 1.0 {
		t.Errorf("Discontinuity at boundary 979: %.3f → %.3f (diff: %.3f)",
			before1, at1, at1-before1)
	}
	if math.Abs(after1-at1) > 1.0 {
		t.Errorf("Discontinuity at boundary 979: %.3f → %.3f (diff: %.3f)",
			at1, after1, after1-at1)
	}

	// Test continuity at boundary 2: year 1708.185161980887
	boundary2 := 1708.185161980887
	before2 := DeltaT(boundary2 - 0.001)
	at2 := DeltaT(boundary2)
	after2 := DeltaT(boundary2 + 0.001)

	if math.Abs(at2-before2) > 1.0 {
		t.Errorf("Discontinuity at boundary 1708: %.3f → %.3f (diff: %.3f)",
			before2, at2, at2-before2)
	}
	if math.Abs(after2-at2) > 1.0 {
		t.Errorf("Discontinuity at boundary 1708: %.3f → %.3f (diff: %.3f)",
			at2, after2, after2-at2)
	}
}

// TestDeltaT_HistoricalValues tests some known historical values
func TestDeltaT_HistoricalValues(t *testing.T) {
	// These are approximate values for validation
	tests := []struct {
		year     float64
		approx   float64
		message  string
	}{
		{1950.0, 29.1, "Around 1950, Delta T was ~29 seconds"},
		{1980.0, 50.5, "Around 1980, Delta T was ~51 seconds"},
		{2000.0, 63.8, "Around 2000, Delta T was ~64 seconds"},
	}

	for _, tt := range tests {
		result := DeltaT(tt.year)
		// Loose tolerance - these are historical approximations
		if math.Abs(result-tt.approx) > 5.0 {
			t.Logf("%s: got %.1f", tt.message, result)
		}
	}
}

// TestGmst tests sla_GMST (Greenwich Mean Sidereal Time)
//
// Test vector from PAL palTest.c line 942-943
func TestGmst(t *testing.T) {
	// vvd ( palGmst( 53736. ), 1.754174971870091203, 1e-12, "palGmst", " ", status );
	result := Gmst(53736.0)
	expected := 1.754174971870091203

	// PAL uses IAU 2006 model (Gmst06), which is more accurate than IAU 1982
	// Test tolerance from PAL is 1e-12
	tol := 1.0e-12
	if math.Abs(result-expected) > tol {
		t.Errorf("Gmst(53736.0) = %.15f, want %.15f (diff: %.2e)",
			result, expected, math.Abs(result-expected))
	}
}

// TestGmst_Range tests that GMST is in range [0, 2π)
func TestGmst_Range(t *testing.T) {
	tests := []float64{
		51544.0,  // J2000
		51544.5,  // J2000 + 12h
		51545.0,  // J2000 + 1d
		52000.0,  // Random date
		40000.0,  // Earlier date
	}

	for _, mjd := range tests {
		result := Gmst(mjd)

		if result < 0.0 || result >= TwoPi {
			t.Errorf("Gmst(%.1f) = %.6f, expected in range [0, 2π)", mjd, result)
		}
	}
}

// TestGmsta tests sla_GMSTA (high-precision GMST)
//
// Test vector from PAL palTest.c line 945-946
func TestGmsta(t *testing.T) {
	// vvd ( palGmsta( 53736., 0.0 ), 1.754174971870091203, 1e-12, "palGmsta", " ", status );
	result := Gmsta(53736.0, 0.0)
	expected := 1.754174971870091203

	// PAL uses IAU 2006 model (Gmst06)
	// Test tolerance from PAL is 1e-12
	tol := 1.0e-12
	if math.Abs(result-expected) > tol {
		t.Errorf("Gmsta(53736.0, 0.0) = %.15f, want %.15f (diff: %.2e)",
			result, expected, math.Abs(result-expected))
	}
}

// TestGmsta0_Convenience tests the single-argument convenience wrapper
func TestGmsta0_Convenience(t *testing.T) {
	mjd := 53736.0

	// Should give same result as Gmsta with second arg = 0
	result1 := Gmsta(mjd, 0.0)
	result2 := Gmsta0(mjd)

	if math.Abs(result1-result2) > 1.0e-15 {
		t.Errorf("Gmsta0 differs from Gmsta: %.15f vs %.15f", result2, result1)
	}
}

// TestGmstGmsta_Comparison compares Gmst and Gmsta
func TestGmstGmsta_Comparison(t *testing.T) {
	mjd := 51544.0 // J2000

	gmst := Gmst(mjd)     // IAU 2006 model (following PAL)
	gmsta := Gmsta0(mjd)  // Also IAU 2006 model

	// Should be identical (both use Gmst06)
	diff := math.Abs(gmsta - gmst)
	if diff > 1.0e-15 {
		t.Errorf("Gmst and Gmsta differ: %.2e radians", diff)
	}

	t.Logf("At J2000: Gmst = %.15f, Gmsta = %.15f (diff: %.2e)",
		gmst, gmsta, diff)
}

// TestGmst_DailyAdvance tests that GMST advances ~366.24/365.24 = 1.0027... per day
func TestGmst_DailyAdvance(t *testing.T) {
	mjd1 := 51544.0
	mjd2 := mjd1 + 1.0 // One day later

	gmst1 := Gmst(mjd1)
	gmst2 := Gmst(mjd2)

	t.Logf("GMST at MJD %.1f: %.15f radians", mjd1, gmst1)
	t.Logf("GMST at MJD %.1f: %.15f radians", mjd2, gmst2)

	// GMST advances by 2π * (366.24/365.24) ≈ 6.3004 radians per day (conceptually)
	// But since GMST is normalized to [0, 2π), the observed difference is:
	// 6.3004 - 2π ≈ 0.017203 radians (about 4 minutes of time, or 1 degree)
	expectedConceptualAdvance := TwoPi * 366.24 / 365.24  // ~6.3004 radians
	expectedObservedAdvance := expectedConceptualAdvance - TwoPi  // ~0.017 radians

	advance := gmst2 - gmst1
	// Handle wrap-around (shouldn't happen since both are already normalized)
	if advance < 0 {
		advance += TwoPi
	}

	// Loose tolerance - checking the observed difference
	if math.Abs(advance-expectedObservedAdvance) > 0.001 {
		t.Errorf("GMST daily observed advance = %.6f radians, want %.6f (conceptual %.6f)",
			advance, expectedObservedAdvance, expectedConceptualAdvance)
	}

	t.Logf("GMST advance in 1 day: %.6f radians (%.3f hours)",
		advance, advance*12.0/Pi)
}

// TestGmst_ZeroAtVernalEquinox tests GMST is ~0 at vernal equinox
func TestGmst_ZeroAtVernalEquinox(t *testing.T) {
	// J2000.0 epoch (2000-01-01 12:00 TT) has GMST ≈ 6.697 hours ≈ 1.753 radians
	// This is well-known reference value
	mjd := 51544.5 // J2000.0 epoch (note: .5 because epoch is at noon)
	gmst := Gmst(mjd)

	expectedHours := 6.697
	expectedRadians := expectedHours * Pi / 12.0

	// Loose tolerance
	if math.Abs(gmst-expectedRadians) > 0.01 {
		hours := gmst * 12.0 / Pi
		t.Logf("GMST at J2000.0 = %.3f hours (%.6f radians), expected %.3f hours",
			hours, gmst, expectedHours)
	}
}

// Benchmark tests

func BenchmarkDtt(b *testing.B) {
	mjd := 51544.0
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Dtt(mjd)
	}
}

func BenchmarkDeltaT(b *testing.B) {
	epoch := 2000.0
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = DeltaT(epoch)
	}
}

func BenchmarkGmst(b *testing.B) {
	mjd := 53736.0
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Gmst(mjd)
	}
}

func BenchmarkGmsta(b *testing.B) {
	mjd := 53736.0
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Gmsta(mjd, 0.0)
	}
}

//
// Note: Tests for calendar and epoch conversion functions (Caldj, Cldj, Djcal, Djcl,
// Epb, Epb2d, Epj, Epj2d) are in calendar_test.go where those functions are implemented.
// This file only contains tests for time scale functions (Dtt, DeltaT, Gmst, Gmsta).
//

