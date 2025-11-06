package slagofa

import (
	"gonum.org/v1/gonum/mat"
	"math"
)

// Phase 7: Coordinate Fitting and Transformation Functions
//
// These functions fit linear transformations between coordinate systems,
// commonly used for astrometric plate solutions.

// FitxyCoeffs represents the 6 coefficients of a linear transformation.
//
// The transformation is:
//   x' = a + b*x + c*y
//   y' = d + e*x + f*y
type FitxyCoeffs struct {
	A, B, C float64 // Coefficients for x transformation
	D, E, F float64 // Coefficients for y transformation
}

// Fitxy fits a linear transformation between two sets of (x,y) coordinates.
//
// Given n points in two coordinate systems, this finds a linear
// transformation that best maps one set to the other.
//
// Original FORTRAN: sla_FITXY by P.T. Wallace
// Go equivalent: Least squares fit using gonum
// SLALIB reference: SUN/67 section 47
//
// Parameters:
//   - itype: Type of model (4 or 6)
//     - 4: Solid body rotation (scale, rotation, translation; no shear/squash)
//     - 6: Full 6-parameter fit (includes shear and squash)
//   - npts: Number of points
//   - xye: Expected/model (x,y) coordinates (output system)
//   - xym: Measured (x,y) coordinates (input system)
//
// Returns:
//   - coeffs: Transformation coefficients (always 6 coefficients)
//   - j: Status (0=OK, -1=illegal ITYPE, -2=insufficient data, -3=no solution)
//
// Notes:
//   - For ITYPE=4: requires at least 2 points
//   - For ITYPE=6: requires at least 3 non-collinear points
//   - Uses least squares for best fit
//   - Transformation: XE = A + B*XM + C*YM, YE = D + E*XM + F*YM
//   - For ITYPE=4: coefficients constrained so |B|=|F| and |C|=|E|
func Fitxy(itype, npts int, xye, xym [][2]float64) (coeffs FitxyCoeffs, j int) {
	// Check ITYPE validity
	if itype != 4 && itype != 6 {
		return FitxyCoeffs{}, -1 // Illegal ITYPE
	}

	// Check minimum points
	minPts := 3
	if itype == 4 {
		minPts = 2
	}
	if npts < minPts {
		return FitxyCoeffs{}, -2 // Insufficient data
	}

	if itype == 6 {
		// Six-coefficient linear model (full fit with shear/squash)
		return fitxy6(npts, xye, xym)
	}

	// Four-coefficient solid body rotation model
	return fitxy4(npts, xye, xym)
}

// fitxy6 performs a 6-parameter linear fit
func fitxy6(npts int, xye, xym [][2]float64) (coeffs FitxyCoeffs, j int) {
	// Build design matrix and observation vectors
	// We want to find XE = A + B*XM + C*YM (and similarly for YE)
	// where XE/YE are expected (xye) and XM/YM are measured (xym)
	A := mat.NewDense(npts, 3, nil)
	bx := mat.NewVecDense(npts, nil)
	by := mat.NewVecDense(npts, nil)

	for i := 0; i < npts; i++ {
		xm := xym[i][0] // Measured X
		ym := xym[i][1] // Measured Y

		// Design matrix: [1, XM, YM]
		A.Set(i, 0, 1.0)
		A.Set(i, 1, xm)
		A.Set(i, 2, ym)

		// Observation vectors (expected values)
		bx.SetVec(i, xye[i][0])
		by.SetVec(i, xye[i][1])
	}

	// Solve normal equations: A^T * A * x = A^T * b
	var ata mat.Dense
	ata.Mul(A.T(), A)

	// Check if singular
	if math.Abs(mat.Det(&ata)) < 1e-10 {
		return FitxyCoeffs{}, -3 // No solution
	}

	// Solve for x coefficients
	var atbx mat.VecDense
	atbx.MulVec(A.T(), bx)

	var xSol mat.VecDense
	err := xSol.SolveVec(&ata, &atbx)
	if err != nil {
		return FitxyCoeffs{}, -3 // No solution
	}

	// Solve for y coefficients
	var atby mat.VecDense
	atby.MulVec(A.T(), by)

	var ySol mat.VecDense
	err = ySol.SolveVec(&ata, &atby)
	if err != nil {
		return FitxyCoeffs{}, -3 // No solution
	}

	// Extract coefficients
	coeffs = FitxyCoeffs{
		A: xSol.AtVec(0),
		B: xSol.AtVec(1),
		C: xSol.AtVec(2),
		D: ySol.AtVec(0),
		E: ySol.AtVec(1),
		F: ySol.AtVec(2),
	}

	return coeffs, 0
}

// fitxy4 performs a 4-parameter solid body rotation fit
// Constrains |B|=|F| and |C|=|E| (no shear or squash)
func fitxy4(npts int, xye, xym [][2]float64) (coeffs FitxyCoeffs, j int) {
	// Try two solutions: with and without sign flip in X
	// Pick the one with smaller residuals

	var bestCoeffs FitxyCoeffs
	var bestRMS float64 = 1e99
	found := false

	for _, sgn := range []float64{1.0, -1.0} {
		// Form summations for this solution
		// Solve for A,B,C,D in: sgn*XE = A + B*XM - C*YM
		//                           YE = D + C*XM + B*YM

		p := float64(npts)
		var sxe, sxxyy, sxyyx, sye, sxm, sym, sx2y2 float64

		for i := 0; i < npts; i++ {
			xe := xye[i][0] * sgn
			ye := xye[i][1]
			xm := xym[i][0]
			ym := xym[i][1]

			sxe += xe
			sxxyy += xe*xm + ye*ym
			sxyyx += xe*ym - ye*xm
			sye += ye
			sxm += xm
			sym += ym
			sx2y2 += xm*xm + ym*ym
		}

		// Set up 4x4 system
		M := mat.NewDense(4, 4, []float64{
			p, sxm, -sym, 0,
			sxm, sx2y2, 0, sym,
			sym, 0, -sx2y2, -sxm,
			0, sym, sxm, p,
		})

		b := mat.NewVecDense(4, []float64{sxe, sxxyy, sxyyx, sye})

		// Solve system
		var sol mat.VecDense
		err := sol.SolveVec(M, b)
		if err != nil {
			continue // Try other sign
		}

		a := sol.AtVec(0)
		bCoeff := sol.AtVec(1)
		cCoeff := sol.AtVec(2)
		d := sol.AtVec(3)

		// Calculate RMS residual
		var sumSqRes float64
		for i := 0; i < npts; i++ {
			xe := xye[i][0]
			ye := xye[i][1]
			xm := xym[i][0]
			ym := xym[i][1]

			xr := sgn*(a + bCoeff*xm - cCoeff*ym) - xe
			yr := d + cCoeff*xm + bCoeff*ym - ye
			sumSqRes += xr*xr + yr*yr
		}

		rms := math.Sqrt(sumSqRes / float64(npts))

		if rms < bestRMS {
			bestRMS = rms
			bestCoeffs = FitxyCoeffs{
				A: sgn * a,
				B: sgn * bCoeff,
				C: sgn * (-cCoeff),
				D: d,
				E: cCoeff,
				F: bCoeff,
			}
			found = true
		}
	}

	if !found {
		return FitxyCoeffs{}, -3 // No solution
	}

	return bestCoeffs, 0
}

// Pxy applies a linear transformation to (x,y) coordinates.
//
// Applies the transformation fitted by Fitxy.
//
// Original FORTRAN: sla_PXY by P.T. Wallace
// Go equivalent: Matrix multiplication
// SLALIB reference: SUN/67 section 76
//
// Parameters:
//   - npts: Number of points to transform
//   - xye: Input (x,y) coordinates
//   - coeffs: Transformation coefficients from Fitxy
//
// Returns:
//   - xym: Transformed (x,y) coordinates
//
// Notes:
//   - Applies: x' = a + b*x + c*y
//              y' = d + e*x + f*y
func Pxy(npts int, xye [][2]float64, coeffs FitxyCoeffs) [][2]float64 {
	xym := make([][2]float64, npts)

	for i := 0; i < npts; i++ {
		x := xye[i][0]
		y := xye[i][1]

		xym[i][0] = coeffs.A + coeffs.B*x + coeffs.C*y
		xym[i][1] = coeffs.D + coeffs.E*x + coeffs.F*y
	}

	return xym
}

// Invf inverts a linear transformation.
//
// Computes the inverse of the transformation represented by coeffs.
//
// Original FORTRAN: sla_INVF by P.T. Wallace
// Go equivalent: Matrix inversion
// SLALIB reference: SUN/67 section 55
//
// Parameters:
//   - coeffs: Forward transformation coefficients
//
// Returns:
//   - invCoeffs: Inverse transformation coefficients
//   - j: Status (0=OK, -1=singular, non-invertible)
//
// Notes:
//   - If T is the forward transform, computes T^{-1}
//   - Fails if transformation is singular (e.g., maps plane to line)
func Invf(coeffs FitxyCoeffs) (invCoeffs FitxyCoeffs, j int) {
	// The transformation matrix is:
	// [x']   [b c] [x]   [a]
	// [y'] = [e f] [y] + [d]
	//
	// To invert:
	// [x]   [b c]^-1 ([x']   [a])
	// [y] = [e f]    ([y'] - [d])

	// Compute determinant
	det := coeffs.B*coeffs.F - coeffs.C*coeffs.E

	if math.Abs(det) < 1e-10 {
		return FitxyCoeffs{}, -1 // Singular
	}

	// Invert the 2×2 linear part
	invB := coeffs.F / det
	invC := -coeffs.C / det
	invE := -coeffs.E / det
	invF := coeffs.B / det

	// Transform the translation part
	invA := -(invB*coeffs.A + invC*coeffs.D)
	invD := -(invE*coeffs.A + invF*coeffs.D)

	invCoeffs = FitxyCoeffs{
		A: invA, B: invB, C: invC,
		D: invD, E: invE, F: invF,
	}

	return invCoeffs, 0
}

// Xy2xy transforms coordinates from one system to another using fitted coefficients.
//
// This is a convenience wrapper that combines Pxy with the transformation.
//
// Original FORTRAN: sla_XY2XY by P.T. Wallace
// Go equivalent: Wrapper around Pxy
// SLALIB reference: SUN/67 section 94
//
// Parameters:
//   - x, y: Input coordinates
//   - coeffs: Transformation coefficients
//
// Returns:
//   - xp, yp: Transformed coordinates
func Xy2xy(x, y float64, coeffs FitxyCoeffs) (xp, yp float64) {
	xp = coeffs.A + coeffs.B*x + coeffs.C*y
	yp = coeffs.D + coeffs.E*x + coeffs.F*y
	return xp, yp
}
