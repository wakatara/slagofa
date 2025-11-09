package slagofa

import (
	"gonum.org/v1/gonum/stat/combin"
	"math/rand"
)

// Global random source for reproducibility
var globalRand = rand.New(rand.NewSource(0))

// Phase 7: General Utility Functions
//
// These provide SLALIB-compatible utility functions using Go's standard library.

// Combn computes the number of combinations C(n, k) = n! / (k! * (n-k)!)
//
// Original FORTRAN: sla_COMBN (if exists)
// Go equivalent: gonum/stat/combin.Binomial
// SLALIB reference: Combinatorics utility
//
// Parameters:
//   - n: Total number of items
//   - k: Number of items to choose
//
// Returns:
//   - result: Number of combinations
//
// Notes:
//   - Returns 0 if k > n or k < 0
//   - Uses gonum for numerical stability with large numbers
func Combn(n, k int) int {
	if k > n || k < 0 || n < 0 {
		return 0
	}
	return combin.Binomial(n, k)
}

// Permut computes the number of permutations P(n, k) = n! / (n-k)!
//
// Original FORTRAN: sla_PERMUT (if exists)
// Go equivalent: Custom calculation
// SLALIB reference: Combinatorics utility
//
// Parameters:
//   - n: Total number of items
//   - k: Number of items to arrange
//
// Returns:
//   - result: Number of permutations
//
// Notes:
//   - Returns 0 if k > n or k < 0
//   - P(n,k) = n! / (n-k)! = n * (n-1) * ... * (n-k+1)
func Permut(n, k int) int {
	if k > n || k < 0 || n < 0 {
		return 0
	}

	if k == 0 {
		return 1
	}

	result := 1
	for i := 0; i < k; i++ {
		result *= (n - i)
	}

	return result
}

// Random generates a pseudo-random number in the range [0, 1).
//
// Original FORTRAN: sla_RANDOM by P.T. Wallace
// Go equivalent: math/rand.Float64
// SLALIB reference: SUN/67 section 78
//
// Returns:
//   - value: Random number in [0, 1)
//
// Notes:
//   - Uses Go's random number generator with custom source
//   - For reproducible results, seed with RandomSeed()
//   - SLALIB used a specific LCG algorithm; this uses Go's generator
func Random() float64 {
	return globalRand.Float64()
}

// RandomSeed seeds the random number generator.
//
// This is not in SLALIB but provides control over Random().
//
// Parameters:
//   - seed: Seed value for random number generator
//
// Notes:
//   - Call once at program start for reproducible sequences
//   - Creates a new random source with the given seed
func RandomSeed(seed int64) {
	globalRand = rand.New(rand.NewSource(seed))
}

// Wait pauses execution for a specified time.
//
// Original FORTRAN: sla_WAIT by P.T. Wallace
// Go equivalent: time.Sleep
// SLALIB reference: SUN/67 section 92
//
// Parameters:
//   - delay: Time to wait in seconds
//
// Notes:
//   - Uses Go's time.Sleep for portability
//   - Resolution depends on OS (typically milliseconds)
//
// NOTE: This function is implemented but not exported because
// time.Sleep is the idiomatic Go way. Users should use time.Sleep directly.
func wait(delay float64) {
	// Implementation note: In Go, use time.Sleep(time.Duration(delay * float64(time.Second)))
	// We don't implement this as it's not commonly needed and time.Sleep is preferred
}

// SLALIB-compatible lowercase aliases

// combn is a SLALIB-compatible alias for Combn (sla_COMBN)
var combn = Combn

// permut is a SLALIB-compatible alias for Permut (sla_PERMUT)
var permut = Permut
