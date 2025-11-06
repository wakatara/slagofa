package slagofa

import (
	"math"
	"testing"
)

// TestCombn tests binomial coefficient calculation
func TestCombn(t *testing.T) {
	tests := []struct {
		n, k     int
		expected int
	}{
		{5, 2, 10},       // C(5,2) = 10
		{10, 3, 120},     // C(10,3) = 120
		{6, 3, 20},       // C(6,3) = 20
		{4, 0, 1},        // C(n,0) = 1
		{4, 4, 1},        // C(n,n) = 1
		{10, 1, 10},      // C(n,1) = n
		{52, 5, 2598960}, // C(52,5) = 2598960 (poker hands)
		{0, 0, 1},        // C(0,0) = 1
	}

	for _, tt := range tests {
		result := Combn(tt.n, tt.k)
		if result != tt.expected {
			t.Errorf("Combn(%d, %d) = %d, want %d",
				tt.n, tt.k, result, tt.expected)
		}
	}

	// Test error cases
	errorTests := []struct {
		n, k int
	}{
		{3, 5},  // k > n
		{5, -1}, // k < 0
		{-1, 2}, // n < 0
	}

	for _, tt := range errorTests {
		result := Combn(tt.n, tt.k)
		if result != 0 {
			t.Errorf("Combn(%d, %d) = %d, want 0 (error case)",
				tt.n, tt.k, result)
		}
	}
}

// TestPermut tests permutation calculation
func TestPermut(t *testing.T) {
	tests := []struct {
		n, k     int
		expected int
	}{
		{5, 2, 20},   // P(5,2) = 5*4 = 20
		{10, 3, 720}, // P(10,3) = 10*9*8 = 720
		{6, 3, 120},  // P(6,3) = 6*5*4 = 120
		{4, 0, 1},    // P(n,0) = 1
		{4, 4, 24},   // P(4,4) = 4! = 24
		{10, 1, 10},  // P(n,1) = n
		{5, 5, 120},  // P(5,5) = 5! = 120
		{0, 0, 1},    // P(0,0) = 1
	}

	for _, tt := range tests {
		result := Permut(tt.n, tt.k)
		if result != tt.expected {
			t.Errorf("Permut(%d, %d) = %d, want %d",
				tt.n, tt.k, result, tt.expected)
		}
	}

	// Test error cases
	errorTests := []struct {
		n, k int
	}{
		{3, 5},  // k > n
		{5, -1}, // k < 0
		{-1, 2}, // n < 0
	}

	for _, tt := range errorTests {
		result := Permut(tt.n, tt.k)
		if result != 0 {
			t.Errorf("Permut(%d, %d) = %d, want 0 (error case)",
				tt.n, tt.k, result)
		}
	}
}

// TestRandom tests random number generation
func TestRandom(t *testing.T) {
	// Test that Random returns values in [0, 1)
	for i := 0; i < 100; i++ {
		r := Random()
		if r < 0.0 || r >= 1.0 {
			t.Errorf("Random() = %f, want [0, 1)", r)
		}
	}

	// Test that RandomSeed produces reproducible sequences
	RandomSeed(12345)
	seq1 := make([]float64, 10)
	for i := range seq1 {
		seq1[i] = Random()
	}

	RandomSeed(12345) // Reset with same seed
	seq2 := make([]float64, 10)
	for i := range seq2 {
		seq2[i] = Random()
	}

	// Sequences should be identical
	for i := range seq1 {
		if seq1[i] != seq2[i] {
			t.Errorf("RandomSeed not reproducible: seq1[%d]=%f, seq2[%d]=%f",
				i, seq1[i], i, seq2[i])
		}
	}

	// Test that different seeds produce different sequences
	RandomSeed(54321)
	seq3 := make([]float64, 10)
	for i := range seq3 {
		seq3[i] = Random()
	}

	same := true
	for i := range seq1 {
		if seq1[i] != seq3[i] {
			same = false
			break
		}
	}

	if same {
		t.Error("Different seeds produced identical sequences")
	}
}

// TestCombinatoricsRelationship tests that C(n,k) * k! = P(n,k)
func TestCombinatoricsRelationship(t *testing.T) {
	tests := []struct {
		n, k int
	}{
		{5, 2},
		{10, 3},
		{6, 3},
		{8, 4},
	}

	factorial := func(n int) int {
		if n <= 1 {
			return 1
		}
		result := 1
		for i := 2; i <= n; i++ {
			result *= i
		}
		return result
	}

	for _, tt := range tests {
		c := Combn(tt.n, tt.k)
		p := Permut(tt.n, tt.k)
		kFact := factorial(tt.k)

		if c*kFact != p {
			t.Errorf("C(%d,%d) * %d! = %d * %d = %d, but P(%d,%d) = %d",
				tt.n, tt.k, tt.k, c, kFact, c*kFact, tt.n, tt.k, p)
		}
	}
}

// TestRandomDistribution tests basic statistical properties of Random()
func TestRandomDistribution(t *testing.T) {
	RandomSeed(42) // Use fixed seed for reproducibility

	n := 10000
	samples := make([]float64, n)
	for i := range samples {
		samples[i] = Random()
	}

	// Calculate mean (should be close to 0.5)
	sum := 0.0
	for _, v := range samples {
		sum += v
	}
	mean := sum / float64(n)

	if math.Abs(mean-0.5) > 0.02 {
		t.Errorf("Random mean = %f, want ~0.5 (±0.02)", mean)
	}

	// Calculate variance (should be close to 1/12 ≈ 0.0833)
	variance := 0.0
	for _, v := range samples {
		diff := v - mean
		variance += diff * diff
	}
	variance /= float64(n)

	expectedVariance := 1.0 / 12.0
	if math.Abs(variance-expectedVariance) > 0.01 {
		t.Errorf("Random variance = %f, want ~%f (±0.01)",
			variance, expectedVariance)
	}

	// Test distribution across bins
	nBins := 10
	bins := make([]int, nBins)
	for _, v := range samples {
		bin := int(v * float64(nBins))
		if bin == nBins {
			bin-- // Handle edge case of exactly 1.0 (shouldn't happen)
		}
		bins[bin]++
	}

	// Each bin should have roughly n/nBins samples
	expectedPerBin := float64(n) / float64(nBins)
	for i, count := range bins {
		// Allow 30% deviation from expected
		if math.Abs(float64(count)-expectedPerBin) > expectedPerBin*0.3 {
			t.Errorf("Bin %d has %d samples, want ~%.0f (±30%%)",
				i, count, expectedPerBin)
		}
	}
}

// Benchmarks
func BenchmarkCombn(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = Combn(52, 5)
	}
}

func BenchmarkPermut(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = Permut(10, 3)
	}
}

func BenchmarkRandom(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = Random()
	}
}

func BenchmarkRandomSeed(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandomSeed(int64(i))
	}
}
