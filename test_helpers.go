package slagofa

import (
	"fmt"
	"math"
	"os"
	"runtime"
	"strings"
	"sync"
	"testing"
)

// TestResult tracks accuracy of a single test against SLALIB
type TestResult struct {
	FunctionName string
	TestName     string
	Source       string // "sla_test.f" or "sla_test.cc" with line number
	Expected     float64
	Actual       float64
	Tolerance    float64
	Deviation    float64
	Passed       bool
}

// Global test results for aggregation
var (
	testResults     []TestResult
	testResultMutex sync.Mutex
	autoTrackingEnabled = true // Can be disabled via env var
)

func init() {
	// Allow disabling auto-tracking via environment variable
	if os.Getenv("SLAGOFA_DISABLE_AUTO_TRACKING") != "" {
		autoTrackingEnabled = false
	}
}

// autoTrackDeviation automatically logs deviations when called from test context
// This is called by almostEqual helpers to transparently track all tests
func autoTrackDeviation(actual, expected, tolerance float64) {
	if !autoTrackingEnabled {
		return
	}

	deviation := math.Abs(actual - expected)
	passed := deviation <= tolerance

	// Get caller information to determine function and test name
	funcName := "unknown"
	testName := "auto-tracked"
	source := "auto"

	// Walk up the call stack to find the test function
	for i := 2; i < 10; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}

		fn := runtime.FuncForPC(pc)
		if fn == nil {
			continue
		}

		name := fn.Name()

		// Look for Test* function
		if strings.Contains(name, ".Test") {
			parts := strings.Split(name, ".")
			if len(parts) > 0 {
				testName = parts[len(parts)-1]
			}
			source = fmt.Sprintf("%s:%d", file, line)

			// Try to extract function being tested from test name
			// e.g., TestNormalizeAngle -> NormalizeAngle
			if strings.HasPrefix(testName, "Test") {
				funcName = strings.TrimPrefix(testName, "Test")
			}
			break
		}
	}

	result := TestResult{
		FunctionName: funcName,
		TestName:     testName,
		Source:       source,
		Expected:     expected,
		Actual:       actual,
		Tolerance:    tolerance,
		Deviation:    deviation,
		Passed:       passed,
	}

	testResultMutex.Lock()
	testResults = append(testResults, result)
	testResultMutex.Unlock()
}

// AssertAlmostEqual checks if two float64 values are equal within tolerance
// and tracks the result for deviation analysis
func AssertAlmostEqual(t *testing.T, funcName, testName, source string, actual, expected, tolerance float64) {
	t.Helper()

	deviation := math.Abs(actual - expected)
	passed := deviation <= tolerance

	result := TestResult{
		FunctionName: funcName,
		TestName:     testName,
		Source:       source,
		Expected:     expected,
		Actual:       actual,
		Tolerance:    tolerance,
		Deviation:    deviation,
		Passed:       passed,
	}

	testResults = append(testResults, result)

	if !passed {
		t.Errorf("%s: %s\n  Expected: %.15e\n  Actual:   %.15e\n  Tolerance: %.2e\n  Deviation: %.2e (%.2e × tolerance)\n  Source: %s",
			funcName, testName, expected, actual, tolerance, deviation, deviation/tolerance, source)
	} else if deviation > 0 {
		// Log successful but non-zero deviation
		t.Logf("%s: %s - PASS with deviation %.2e (%.1f%% of tolerance) - Source: %s",
			funcName, testName, deviation, (deviation/tolerance)*100, source)
	}
}

// AssertAlmostEqual32 is the float32 version
func AssertAlmostEqual32(t *testing.T, funcName, testName, source string, actual, expected, tolerance float32) {
	t.Helper()
	AssertAlmostEqual(t, funcName, testName, source, float64(actual), float64(expected), float64(tolerance))
}

// GenerateDeviationReport creates a summary report of all test deviations
func GenerateDeviationReport(t *testing.T) {
	t.Helper()

	if len(testResults) == 0 {
		return
	}

	totalTests := len(testResults)
	passed := 0
	failed := 0
	perfectMatches := 0

	var maxDeviation TestResult
	var worstRatio TestResult
	worstRatio.Deviation = 0
	maxDeviation.Deviation = 0

	for _, result := range testResults {
		if result.Passed {
			passed++
			if result.Deviation == 0 {
				perfectMatches++
			}

			// Track maximum deviation (even in passing tests)
			if result.Deviation > maxDeviation.Deviation {
				maxDeviation = result
			}

			// Track worst ratio (deviation / tolerance)
			ratio := result.Deviation / result.Tolerance
			worstRatioValue := worstRatio.Deviation / worstRatio.Tolerance
			if ratio > worstRatioValue {
				worstRatio = result
			}
		} else {
			failed++
		}
	}

	t.Logf("\n=== SLALIB Test Deviation Report ===")
	t.Logf("Total tests: %d", totalTests)
	t.Logf("Passed: %d (%.1f%%)", passed, float64(passed)/float64(totalTests)*100)
	t.Logf("Failed: %d", failed)
	t.Logf("Perfect matches (zero deviation): %d (%.1f%%)", perfectMatches, float64(perfectMatches)/float64(totalTests)*100)
	t.Logf("")

	if maxDeviation.Deviation > 0 {
		t.Logf("Largest deviation (in passing tests):")
		t.Logf("  Function: %s - %s", maxDeviation.FunctionName, maxDeviation.TestName)
		t.Logf("  Deviation: %.2e (tolerance: %.2e, ratio: %.2f%%)",
			maxDeviation.Deviation, maxDeviation.Tolerance,
			(maxDeviation.Deviation/maxDeviation.Tolerance)*100)
		t.Logf("  Source: %s", maxDeviation.Source)
		t.Logf("")
	}

	if worstRatio.Deviation > 0 {
		ratio := worstRatio.Deviation / worstRatio.Tolerance
		t.Logf("Worst tolerance ratio (in passing tests):")
		t.Logf("  Function: %s - %s", worstRatio.FunctionName, worstRatio.TestName)
		t.Logf("  Deviation: %.2e (tolerance: %.2e, ratio: %.2f%%)",
			worstRatio.Deviation, worstRatio.Tolerance, ratio*100)
		t.Logf("  Source: %s", worstRatio.Source)
		t.Logf("")
	}

	// Report by function
	funcStats := make(map[string]struct{ total, passed, perfect int })
	for _, result := range testResults {
		stats := funcStats[result.FunctionName]
		stats.total++
		if result.Passed {
			stats.passed++
			if result.Deviation == 0 {
				stats.perfect++
			}
		}
		funcStats[result.FunctionName] = stats
	}

	t.Logf("Per-function accuracy:")
	for funcName, stats := range funcStats {
		accuracy := float64(stats.passed) / float64(stats.total) * 100
		perfectRate := float64(stats.perfect) / float64(stats.total) * 100
		t.Logf("  %s: %d/%d passed (%.1f%%), %d perfect (%.1f%%)",
			funcName, stats.passed, stats.total, accuracy, stats.perfect, perfectRate)
	}
}

// ResetTestResults clears the accumulated test results
func ResetTestResults() {
	testResults = nil
}

// Helper for common test patterns
func slalibTest(t *testing.T, funcName, testName, source string, actual, expected, tolerance float64) {
	t.Helper()
	AssertAlmostEqual(t, funcName, testName, source, actual, expected, tolerance)
}

// FormatSLALIBSource creates a standard source reference string
func FormatSLALIBSource(filename string, lineNumber int) string {
	return fmt.Sprintf("%s:%d", filename, lineNumber)
}
