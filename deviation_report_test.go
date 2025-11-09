package slagofa

import (
	"strings"
	"testing"
)

// TestZZZ_DeviationReport runs at the end to summarize all test deviations
// Named with ZZZ prefix to run last (tests run alphabetically)
// Run with: go test -v
func TestZZZ_DeviationReport(t *testing.T) {
	// This test runs last and reports on all accumulated test results
	// from the automatic deviation tracking in almostEqual()

	// Generate the deviation report from all automatically tracked tests
	t.Log("\n" + strings.Repeat("=", 70))
	t.Log("AUTOMATIC DEVIATION REPORT FROM ALL TESTS")
	t.Log(strings.Repeat("=", 70))

	GenerateDeviationReport(t)

	t.Log(strings.Repeat("=", 70))
	t.Log("\nNOTE: All deviations were automatically tracked via almostEqual().")
	t.Log("To see individual test deviations, run: go test -v")
	t.Log("To disable auto-tracking, set: SLAGOFA_DISABLE_AUTO_TRACKING=1")
}

// Example of how to update an existing test to use deviation tracking
func TestNormalizeAngleWithDeviation(t *testing.T) {
	tests := []struct {
		name      string
		input     float64
		expected  float64
		source    string
		tolerance float64
	}{
		{
			name:      "SLALIB test cc:514",
			input:     -4.0,
			expected:  2.283185307179586,
			source:    FormatSLALIBSource("sla_test.cc", 514),
			tolerance: 1.0e-12,
		},
		{
			name:      "SLALIB test cc:515",
			input:     -0.1,
			expected:  6.183185307179587 - TwoPi, // Should be in [-π, π]
			source:    FormatSLALIBSource("sla_test.cc", 515),
			tolerance: 1.0e-12,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NormalizeAngle(tt.input)
			AssertAlmostEqual(t, "NormalizeAngle", tt.name, tt.source,
				result, tt.expected, tt.tolerance)
		})
	}
}

/*
USAGE GUIDE FOR OTHER TEST FILES:

1. Import deviation tracking in your test file:
   (Already available in slagofa package)

2. Replace manual comparison with AssertAlmostEqual:

   OLD:
   if math.Abs(result - expected) > tolerance {
       t.Errorf("failed: got %v, want %v", result, expected)
   }

   NEW:
   AssertAlmostEqual(t, "FunctionName", "test description",
       FormatSLALIBSource("sla_test.f", 123),
       result, expected, tolerance)

3. At package level, add a TestDeviationReport that calls GenerateDeviationReport(t)

4. Run tests with: go test -v
   The deviation report will appear at the end showing:
   - Total pass/fail counts
   - Perfect matches (zero deviation)
   - Largest deviations even in passing tests
   - Worst tolerance ratios
   - Per-function accuracy statistics

EXAMPLE CONVERSION:

// OLD TEST:
func TestOldStyle(t *testing.T) {
    result := MyFunc(input)
    if math.Abs(result - expected) > 1.0e-12 {
        t.Errorf("MyFunc failed")
    }
}

// NEW TEST WITH DEVIATION TRACKING:
func TestNewStyle(t *testing.T) {
    result := MyFunc(input)
    AssertAlmostEqual(t, "MyFunc", "basic test",
        FormatSLALIBSource("sla_test.f", 456),
        result, expected, 1.0e-12)
}
*/
