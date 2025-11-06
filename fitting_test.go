package slagofa

import (
	"math"
	"testing"
)

// TestFitxy tests coordinate fitting with SLALIB test vectors
func TestFitxy(t *testing.T) {
	// From SLALIB test suite (sla_test.f line 2876)
	// 8 points to fit
	npts := 8

	// Input coordinates (measured)
	xye := [][2]float64{
		{-23.4, -12.1},
		{32.0, -15.3},
		{10.9, 23.7},
		{-3.0, 16.1},
		{45.0, 32.5},
		{8.6, -17.0},
		{15.3, 10.0},
		{121.7, -3.8},
	}

	// Model coordinates
	xym := [][2]float64{
		{-23.41, 12.12},
		{32.03, 15.34},
		{10.93, -23.72},
		{-3.01, -16.10},
		{44.90, -32.46},
		{8.55, 17.02},
		{15.31, -10.07},
		{120.92, 3.81},
	}

	// Expected coefficients for 6-parameter fit (from SLALIB line 2911)
	expectedCoeffs := FitxyCoeffs{
		A: -2.617232551841476e-2,
		B: 1.005634905041421,
		C: 2.133045023329208e-3,
		D: 3.846993364417779909e-3,
		E: 1.301671386431460e-4,
		F: -0.9994827065693964,
	}

	coeffs, j := Fitxy(6, npts, xye, xym)

	// Check status
	if j != 0 {
		t.Errorf("Fitxy status = %d, want 0", j)
	}

	// Check coefficients (tolerance 1e-12 from SLALIB)
	if !almostEqual(coeffs.A, expectedCoeffs.A, 1e-12) {
		t.Errorf("Fitxy coeff A = %.15e, want %.15e", coeffs.A, expectedCoeffs.A)
	}
	if !almostEqual(coeffs.B, expectedCoeffs.B, 1e-12) {
		t.Errorf("Fitxy coeff B = %.15e, want %.15e", coeffs.B, expectedCoeffs.B)
	}
	if !almostEqual(coeffs.C, expectedCoeffs.C, 1e-12) {
		t.Errorf("Fitxy coeff C = %.15e, want %.15e", coeffs.C, expectedCoeffs.C)
	}
	if !almostEqual(coeffs.D, expectedCoeffs.D, 1e-12) {
		t.Errorf("Fitxy coeff D = %.15e, want %.15e", coeffs.D, expectedCoeffs.D)
	}
	if !almostEqual(coeffs.E, expectedCoeffs.E, 1e-12) {
		t.Errorf("Fitxy coeff E = %.15e, want %.15e", coeffs.E, expectedCoeffs.E)
	}
	if !almostEqual(coeffs.F, expectedCoeffs.F, 1e-12) {
		t.Errorf("Fitxy coeff F = %.15e, want %.15e", coeffs.F, expectedCoeffs.F)
	}

	// Test insufficient points for 6-parameter
	_, status := Fitxy(6, 2, xye[:2], xym[:2])
	if status != -2 {
		t.Errorf("Fitxy(6) with 2 points status = %d, want -2 (insufficient data)", status)
	}

	// Test invalid ITYPE
	_, status = Fitxy(5, npts, xye, xym)
	if status != -1 {
		t.Errorf("Fitxy with invalid ITYPE status = %d, want -1", status)
	}

	// Test 4-parameter fit (solid body rotation)
	coeffs4, j4 := Fitxy(4, npts, xye, xym)
	if j4 != 0 {
		t.Errorf("Fitxy(4) status = %d, want 0", j4)
	}
	// 4-parameter fit should have |B|=|F| and |C|=|E| constraints (solid body rotation)
	// Verify the constraint is approximately satisfied
	if !almostEqual(math.Abs(coeffs4.B), math.Abs(coeffs4.F), 1e-10) {
		t.Errorf("Fitxy(4): |B|=%.10f != |F|=%.10f (should be equal for solid body rotation)",
			math.Abs(coeffs4.B), math.Abs(coeffs4.F))
	}
	if !almostEqual(math.Abs(coeffs4.C), math.Abs(coeffs4.E), 1e-10) {
		t.Errorf("Fitxy(4): |C|=%.10f != |E|=%.10f (should be equal for solid body rotation)",
			math.Abs(coeffs4.C), math.Abs(coeffs4.E))
	}
}

// TestPxy tests applying transformation with SLALIB test vectors
func TestPxy(t *testing.T) {
	// Use same test data as TestFitxy
	npts := 8

	// Measured coordinates (input to transformation)
	xym := [][2]float64{
		{-23.41, 12.12},
		{32.03, 15.34},
		{10.93, -23.72},
		{-3.01, -16.10},
		{44.90, -32.46},
		{8.55, 17.02},
		{15.31, -10.07},
		{120.92, 3.81},
	}

	// Coefficients from previous fit
	coeffs := FitxyCoeffs{
		A: -2.617232551841476e-2,
		B: 1.005634905041421,
		C: 2.133045023329208e-3,
		D: 3.846993364417779909e-3,
		E: 1.301671386431460e-4,
		F: -0.9994827065693964,
	}

	// Expected transformed coordinates (from SLALIB line 2927)
	expected := [][2]float64{
		{-23.542232946855340, -12.11293062297230597},
		{32.217034593616180, -15.324048471959370},
		{10.914821358630950, 23.712999520015880},
		{-3.087475414568693, 16.09512676604438414},
		{45.05759626938414666, 32.45290015313210889},
		{8.608310538882801, -17.006235743411300},
		{15.348618307280820, 10.07063070741086835},
		{121.5833272936291482, -3.788442308260240},
	}

	result := Pxy(npts, xym, coeffs)

	// Check transformed coordinates (tolerance 1e-12 from SLALIB)
	for i := 0; i < npts; i++ {
		if !almostEqual(result[i][0], expected[i][0], 1e-12) {
			t.Errorf("Pxy[%d][0] = %.15e, want %.15e", i, result[i][0], expected[i][0])
		}
		if !almostEqual(result[i][1], expected[i][1], 1e-12) {
			t.Errorf("Pxy[%d][1] = %.15e, want %.15e", i, result[i][1], expected[i][1])
		}
	}
}

// TestInvf tests transformation inversion with SLALIB test vectors
func TestInvf(t *testing.T) {
	// Coefficients to invert (from SLALIB line 2911)
	coeffs := FitxyCoeffs{
		A: -2.617232551841476e-2,
		B: 1.005634905041421,
		C: 2.133045023329208e-3,
		D: 3.846993364417779909e-3,
		E: 1.301671386431460e-4,
		F: -0.9994827065693964,
	}

	// Expected inverse coefficients (from SLALIB line 2977)
	expected := FitxyCoeffs{
		A: 0.02601750208015891,
		B: 0.9943963945040283,
		C: 0.002122190075497872,
		D: 0.003852372795357474353,
		E: 0.0001295047252932767,
		F: -1.000517284779212,
	}

	inv, j := Invf(coeffs)

	// Check status
	if j != 0 {
		t.Errorf("Invf status = %d, want 0", j)
	}

	// Check inverse coefficients (tolerance 1e-12 from SLALIB)
	if !almostEqual(inv.A, expected.A, 1e-12) {
		t.Errorf("Invf coeff A = %.15e, want %.15e", inv.A, expected.A)
	}
	if !almostEqual(inv.B, expected.B, 1e-12) {
		t.Errorf("Invf coeff B = %.15e, want %.15e", inv.B, expected.B)
	}
	if !almostEqual(inv.C, expected.C, 1e-12) {
		t.Errorf("Invf coeff C = %.15e, want %.15e", inv.C, expected.C)
	}
	if !almostEqual(inv.D, expected.D, 1e-12) {
		t.Errorf("Invf coeff D = %.15e, want %.15e", inv.D, expected.D)
	}
	if !almostEqual(inv.E, expected.E, 1e-12) {
		t.Errorf("Invf coeff E = %.15e, want %.15e", inv.E, expected.E)
	}
	if !almostEqual(inv.F, expected.F, 1e-12) {
		t.Errorf("Invf coeff F = %.15e, want %.15e", inv.F, expected.F)
	}

	// Test singular transformation
	singular := FitxyCoeffs{
		A: 0, B: 1, C: 2,
		D: 0, E: 2, F: 4, // E*C - B*F = 0 (singular)
	}

	_, status := Invf(singular)
	if status != -1 {
		t.Errorf("Invf(singular) status = %d, want -1", status)
	}
}

// TestXy2xy tests single point transformation with SLALIB test vectors
func TestXy2xy(t *testing.T) {
	// Coefficients from SLALIB test
	coeffs := FitxyCoeffs{
		A: -2.617232551841476e-2,
		B: 1.005634905041421,
		C: 2.133045023329208e-3,
		D: 3.846993364417779909e-3,
		E: 1.301671386431460e-4,
		F: -0.9994827065693964,
	}

	// Test point from SLALIB line 2996
	x, y := 44.5, 32.5

	// Expected transformed coordinates (from SLALIB line 2997)
	expectedX := 44.793904912083030
	expectedY := -32.473548532471330

	xp, yp := Xy2xy(x, y, coeffs)

	// Check results (tolerance 1e-11 from SLALIB)
	if !almostEqual(xp, expectedX, 1e-11) {
		t.Errorf("Xy2xy x = %.15e, want %.15e", xp, expectedX)
	}
	if !almostEqual(yp, expectedY, 1e-11) {
		t.Errorf("Xy2xy y = %.15e, want %.15e", yp, expectedY)
	}
}

// TestFittingRoundTrip tests that forward and inverse transformations cancel
func TestFittingRoundTrip(t *testing.T) {
	// Define a transformation
	coeffs := FitxyCoeffs{
		A: 1.0, B: 0.9, C: 0.1,
		D: 2.0, E: -0.05, F: 1.05,
	}

	// Invert it
	inv, status := Invf(coeffs)
	if status != 0 {
		t.Fatalf("Invf failed with status %d", status)
	}

	// Test several points
	testPoints := [][2]float64{
		{10.0, 20.0},
		{-5.0, 15.0},
		{0.0, 0.0},
		{100.0, -50.0},
	}

	for _, pt := range testPoints {
		// Transform forward
		x1, y1 := Xy2xy(pt[0], pt[1], coeffs)

		// Transform back
		x2, y2 := Xy2xy(x1, y1, inv)

		// Should get original point back
		if !almostEqual(x2, pt[0], 1e-10) || !almostEqual(y2, pt[1], 1e-10) {
			t.Errorf("Round trip failed for (%.2f, %.2f): got (%.10f, %.10f)",
				pt[0], pt[1], x2, y2)
		}
	}
}

// Benchmarks
func BenchmarkFitxy(b *testing.B) {
	npts := 8

	xye := [][2]float64{
		{-23.4, -12.1},
		{32.0, -15.3},
		{10.9, 23.7},
		{-3.0, 16.1},
		{45.0, 32.5},
		{8.6, -17.0},
		{15.3, 10.0},
		{121.7, -3.8},
	}

	xym := [][2]float64{
		{-23.41, 12.12},
		{32.03, 15.34},
		{10.93, -23.72},
		{-3.01, -16.10},
		{44.90, -32.46},
		{8.55, 17.02},
		{15.31, -10.07},
		{120.92, 3.81},
	}

	for i := 0; i < b.N; i++ {
		_, _ = Fitxy(6, npts, xye, xym)
	}
}

func BenchmarkPxy(b *testing.B) {
	npts := 8

	xye := [][2]float64{
		{-23.4, -12.1},
		{32.0, -15.3},
		{10.9, 23.7},
		{-3.0, 16.1},
		{45.0, 32.5},
		{8.6, -17.0},
		{15.3, 10.0},
		{121.7, -3.8},
	}

	coeffs := FitxyCoeffs{
		A: -2.617232551841476e-2,
		B: 1.005634905041421,
		C: 2.133045023329208e-3,
		D: 3.846993364417779909e-3,
		E: 1.301671386431460e-4,
		F: -0.9994827065693964,
	}

	for i := 0; i < b.N; i++ {
		_ = Pxy(npts, xye, coeffs)
	}
}

func BenchmarkInvf(b *testing.B) {
	coeffs := FitxyCoeffs{
		A: -2.617232551841476e-2,
		B: 1.005634905041421,
		C: 2.133045023329208e-3,
		D: 3.846993364417779909e-3,
		E: 1.301671386431460e-4,
		F: -0.9994827065693964,
	}

	for i := 0; i < b.N; i++ {
		_, _ = Invf(coeffs)
	}
}

func BenchmarkXy2xy(b *testing.B) {
	coeffs := FitxyCoeffs{
		A: -2.617232551841476e-2,
		B: 1.005634905041421,
		C: 2.133045023329208e-3,
		D: 3.846993364417779909e-3,
		E: 1.301671386431460e-4,
		F: -0.9994827065693964,
	}

	x, y := 44.5, 32.5

	for i := 0; i < b.N; i++ {
		_, _ = Xy2xy(x, y, coeffs)
	}
}
