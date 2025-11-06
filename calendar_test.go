package slagofa

import (
	"math"
	"testing"
)

// Calendar and Epoch test vectors from PAL test suite (palTest.c)

const calendarTolerance = 1.0e-9 // Tolerance for date/time calculations

// TestCldj tests sla_CLDJ (calendar to MJD)
//
// Test vector from PAL palTest.c line 455-457
func TestCldj(t *testing.T) {
	// palCaldj ( 1999, 12, 31, &djm, &j );
	// vvd ( djm, 51543, 0, "palCaldj", " ", status );
	djm, status := Cldj(1999, 12, 31)

	if status != 0 {
		t.Errorf("Cldj returned status %d, want 0", status)
	}

	expected := 51543.0
	if math.Abs(djm-expected) > calendarTolerance {
		t.Errorf("Cldj(1999, 12, 31) = %.15f, want %.15f", djm, expected)
	}
}

// TestCldj_EdgeCases tests various edge cases
func TestCldj_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		year     int
		month    int
		day      int
		expected float64
		wantErr  bool
	}{
		{"2000-01-01 (J2000.0 epoch)", 2000, 1, 1, 51544.0, false},
		{"1858-11-17 (MJD epoch)", 1858, 11, 17, 0.0, false},
		{"Leap year Feb 29", 2000, 2, 29, 51603.0, false},
		{"Non-leap year Feb 28", 1900, 2, 28, 15078.0, false},
		{"Invalid month (0)", 2000, 0, 1, 0.0, true},
		{"Invalid month (13)", 2000, 13, 1, 0.0, true},
		{"Invalid day (0)", 2000, 1, 0, 0.0, true},
		{"Invalid day (32)", 2000, 1, 32, 0.0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			djm, status := Cldj(tt.year, tt.month, tt.day)

			if tt.wantErr {
				if status == 0 {
					t.Errorf("Cldj(%d, %d, %d) succeeded, want error",
						tt.year, tt.month, tt.day)
				}
			} else {
				if status != 0 {
					t.Errorf("Cldj(%d, %d, %d) returned status %d, want 0",
						tt.year, tt.month, tt.day, status)
				}
				if math.Abs(djm-tt.expected) > calendarTolerance {
					t.Errorf("Cldj(%d, %d, %d) = %.15f, want %.15f",
						tt.year, tt.month, tt.day, djm, tt.expected)
				}
			}
		})
	}
}

// TestDjcl tests sla_DJCL (MJD to calendar)
func TestDjcl(t *testing.T) {
	// Round-trip test: convert known date to MJD and back
	origYear := 1999
	origMonth := 12
	origDay := 31

	djm, status := Cldj(origYear, origMonth, origDay)
	if status != 0 {
		t.Fatalf("Cldj failed with status %d", status)
	}

	year, month, day, fd, status := Djcl(djm)
	if status != 0 {
		t.Errorf("Djcl returned status %d, want 0", status)
	}

	if year != origYear || month != origMonth || day != origDay {
		t.Errorf("Djcl(%.1f) = (%d, %d, %d), want (%d, %d, %d)",
			djm, year, month, day, origYear, origMonth, origDay)
	}

	// Fraction should be 0.0 for midnight
	if math.Abs(fd) > calendarTolerance {
		t.Errorf("Djcl fraction = %.15f, want 0.0", fd)
	}
}

// TestDjcl_WithFraction tests MJD with fractional day
func TestDjcl_WithFraction(t *testing.T) {
	// MJD 51544.5 = 2000-01-01 12:00 UTC
	// (MJD 51544.0 = 2000-01-01 00:00 UTC)
	year, month, day, fd, status := Djcl(51544.5)
	if status != 0 {
		t.Errorf("Djcl returned status %d, want 0", status)
	}

	// Should be January 1, 2000 at 12:00 (fraction = 0.5)
	if year != 2000 || month != 1 || day != 1 {
		t.Errorf("Djcl(51544.5) date = (%d, %d, %d), want (2000, 1, 1)",
			year, month, day)
	}

	// Fraction should be 0.5 (12:00 UTC)
	if math.Abs(fd-0.5) > calendarTolerance {
		t.Errorf("Djcl(51544.5) fraction = %.15f, want 0.5", fd)
	}
}

// TestCaldj tests sla_CALDJ (calendar with 2-digit year handling)
func TestCaldj(t *testing.T) {
	tests := []struct {
		name     string
		year     int
		month    int
		day      int
		expected float64
	}{
		{"4-digit year 1999", 1999, 12, 31, 51543.0},
		{"2-digit year 00 → 2000", 0, 1, 1, 51544.0},
		{"2-digit year 49 → 2049", 49, 1, 1, 69442.0},
		{"2-digit year 50 → 1950", 50, 1, 1, 33282.0},
		{"2-digit year 99 → 1999", 99, 12, 31, 51543.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			djm, status := Caldj(tt.year, tt.month, tt.day)

			if status != 0 {
				t.Errorf("Caldj(%d, %d, %d) returned status %d, want 0",
					tt.year, tt.month, tt.day, status)
			}

			if math.Abs(djm-tt.expected) > calendarTolerance {
				t.Errorf("Caldj(%d, %d, %d) = %.15f, want %.15f",
					tt.year, tt.month, tt.day, djm, tt.expected)
			}
		})
	}
}

// TestCalyd tests sla_CALYD (calendar to year + day number)
func TestCalyd(t *testing.T) {
	tests := []struct {
		name      string
		year      int
		month     int
		day       int
		expectDay int
	}{
		{"January 1", 2024, 1, 1, 1},
		{"January 31", 2024, 1, 31, 31},
		{"March 1 (leap year)", 2024, 3, 1, 61}, // Day 61 in leap year
		{"March 1 (non-leap)", 2023, 3, 1, 60},  // Day 60 in non-leap year
		{"December 31 (leap)", 2024, 12, 31, 366},
		{"December 31 (non-leap)", 2023, 12, 31, 365},
		{"Mid-year", 2024, 7, 15, 197},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ny, nd, status := Calyd(tt.year, tt.month, tt.day)

			if status != 0 {
				t.Errorf("Calyd(%d, %d, %d) returned status %d, want 0",
					tt.year, tt.month, tt.day, status)
			}

			if ny != tt.year {
				t.Errorf("Calyd year = %d, want %d", ny, tt.year)
			}

			if nd != tt.expectDay {
				t.Errorf("Calyd(%d, %d, %d) day = %d, want %d",
					tt.year, tt.month, tt.day, nd, tt.expectDay)
			}
		})
	}
}

// TestClyd tests sla_CLYD (year + day number to calendar)
func TestClyd(t *testing.T) {
	tests := []struct {
		name        string
		year        int
		dayNum      int
		expectMonth int
		expectDay   int
	}{
		{"Day 1", 2024, 1, 1, 1},
		{"Day 31", 2024, 31, 1, 31},
		{"Day 32", 2024, 32, 2, 1},
		{"Day 60 (leap year)", 2024, 60, 2, 29},  // Feb 29
		{"Day 61 (leap year)", 2024, 61, 3, 1},   // Mar 1
		{"Day 60 (non-leap)", 2023, 60, 3, 1},    // Mar 1
		{"Day 365 (non-leap)", 2023, 365, 12, 31},
		{"Day 366 (leap year)", 2024, 366, 12, 31},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ny, nm, nd, status := Clyd(tt.year, tt.dayNum)

			if status != 0 {
				t.Errorf("Clyd(%d, %d) returned status %d, want 0",
					tt.year, tt.dayNum, status)
			}

			if ny != tt.year {
				t.Errorf("Clyd year = %d, want %d", ny, tt.year)
			}

			if nm != tt.expectMonth {
				t.Errorf("Clyd(%d, %d) month = %d, want %d",
					tt.year, tt.dayNum, nm, tt.expectMonth)
			}

			if nd != tt.expectDay {
				t.Errorf("Clyd(%d, %d) day = %d, want %d",
					tt.year, tt.dayNum, nd, tt.expectDay)
			}
		})
	}
}

// TestCalydClydRoundTrip tests round-trip conversion
func TestCalydClydRoundTrip(t *testing.T) {
	tests := []struct {
		year  int
		month int
		day   int
	}{
		{2024, 1, 1},
		{2024, 2, 29},
		{2023, 12, 31},
		{2000, 6, 15},
		{1999, 12, 31},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			// Convert to day number
			ny1, nd, status := Calyd(tt.year, tt.month, tt.day)
			if status != 0 {
				t.Fatalf("Calyd failed with status %d", status)
			}

			// Convert back to calendar
			ny2, nm, nd2, status := Clyd(ny1, nd)
			if status != 0 {
				t.Fatalf("Clyd failed with status %d", status)
			}

			// Should match original
			if ny2 != tt.year || nm != tt.month || nd2 != tt.day {
				t.Errorf("Round-trip (%d, %d, %d) → (%d, %d) → (%d, %d, %d)",
					tt.year, tt.month, tt.day, ny1, nd, ny2, nm, nd2)
			}
		})
	}
}

// TestEpb tests sla_EPB (Julian Date to Besselian epoch)
//
// Test vector from PAL palTest.c line 815-816
func TestEpb(t *testing.T) {
	// vvd ( palEpb( 45123 ), 1982.419793168669, 1e-8, "palEpb", " ", status );
	result := Epb(45123.0)
	expected := 1982.419793168669

	if math.Abs(result-expected) > 1.0e-8 {
		t.Errorf("Epb(45123) = %.15f, want %.15f", result, expected)
	}
}

// TestEpb2d tests sla_EPB2D (Besselian epoch to Julian Date)
//
// Test vector from PAL palTest.c line 820-821
func TestEpb2d(t *testing.T) {
	// vvd ( palEpb2d( 1975.5 ), 42595.5995279655, 1e-7, "palEpb2d", " ", status );
	result := Epb2d(1975.5)
	expected := 42595.5995279655

	if math.Abs(result-expected) > 1.0e-7 {
		t.Errorf("Epb2d(1975.5) = %.15f, want %.15f", result, expected)
	}
}

// TestEpbRoundTrip tests Besselian epoch round-trip
func TestEpbRoundTrip(t *testing.T) {
	tests := []float64{
		45000.0,  // ~1982
		51544.0,  // J2000
		40000.0,  // ~1968
	}

	for _, mjd := range tests {
		epoch := Epb(mjd)
		recovered := Epb2d(epoch)

		if math.Abs(recovered-mjd) > 1.0e-6 {
			t.Errorf("Epb round-trip: %.1f → %.12f → %.12f (diff: %.2e)",
				mjd, epoch, recovered, math.Abs(recovered-mjd))
		}
	}
}

// TestEpj tests sla_EPJ (Julian Date to Julian epoch)
func TestEpj(t *testing.T) {
	// J2000.0 is defined as MJD 51544.5
	result := Epj(51544.5)
	expected := 2000.0

	if math.Abs(result-expected) > 1.0e-9 {
		t.Errorf("Epj(51544.5) = %.15f, want 2000.0", result)
	}
}

// TestEpj2d tests sla_EPJ2D (Julian epoch to Julian Date)
func TestEpj2d(t *testing.T) {
	result := Epj2d(2000.0)
	expected := 51544.5 // J2000.0 = MJD 51544.5

	if math.Abs(result-expected) > 1.0e-9 {
		t.Errorf("Epj2d(2000.0) = %.15f, want 51544.5", result)
	}
}

// TestEpjRoundTrip tests Julian epoch round-trip
func TestEpjRoundTrip(t *testing.T) {
	tests := []float64{
		51544.5,  // J2000.0
		51179.5,  // J1999.0
		51910.0,  // J2001.0
	}

	for _, mjd := range tests {
		epoch := Epj(mjd)
		recovered := Epj2d(epoch)

		if math.Abs(recovered-mjd) > 1.0e-9 {
			t.Errorf("Epj round-trip: %.1f → %.12f → %.12f (diff: %.2e)",
				mjd, epoch, recovered, math.Abs(recovered-mjd))
		}
	}
}

// TestEpochComparison compares Besselian and Julian epochs
func TestEpochComparison(t *testing.T) {
	// At J2000.0
	mjd := 51544.5
	bep := Epb(mjd)
	jep := Epj(mjd)

	// Besselian epoch should be slightly greater than 2000 at J2000.0
	// The Besselian year (tropical year = 365.242198781 days) is shorter
	// than the Julian year (365.25 days), so B2000 occurs before J2000
	// At J2000.0, we're slightly past B2000.0
	if bep <= 2000.0 || bep >= 2000.01 {
		t.Errorf("Besselian epoch at J2000 = %.6f, expected slightly > 2000.0", bep)
	}

	if math.Abs(jep-2000.0) > 1.0e-9 {
		t.Errorf("Julian epoch at J2000 = %.15f, want 2000.0", jep)
	}

	t.Logf("At MJD %.1f: Besselian = %.6f, Julian = %.6f (diff = %.6f)",
		mjd, bep, jep, jep-bep)
}

// Benchmark tests

func BenchmarkCldj(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = Cldj(1999, 12, 31)
	}
}

func BenchmarkDjcl(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _, _, _, _ = Djcl(51543.0)
	}
}

func BenchmarkEpb(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = Epb(45123.0)
	}
}

func BenchmarkEpj(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = Epj(51544.5)
	}
}
