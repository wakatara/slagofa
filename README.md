# slagofa - SLALIB-Compatible API for GoFA

A Go library providing a sla_DVNLALIB-compatible API layer on top of the
International Astronomic Union's standard [GoFA](https://github.com/hebl/gofa)
(Golang Standards Of Fundamental Astronomy) library.

`slagofa` seeks to provide a modern, accessible, drop-in, and maintainable
replacement for astronomers preferring or having code dependent on the SLA
linraries while reconciling this with recognized global standards from [IAC
SOFA](https://www.iausofa.org/) but in Go lang while maintaining blistering
performance, the friendly API astronomers love, and the advantages of the GO
ecosystem.

SLA aliases are provided to try to reach 100% API compatability. Where SOFA does
not provide a function, Go community scientific and mathematics libraries have
been used to supplement 100% coverage.

While primary meant to be used in a Go ecosystem, we are also investigating
performant C bindings so this can be used as a full, modern, drop-in replacement
for older code bases that are using `slalib.so`

THIS IS CURRENTLY A WORK IN PROGRESS.

SLALIB contains 190 functions overall.

Of these 190, as of Nov 8th, 2025:

- **92 implemented (48.4%)** with verified SLALIB-compatible aliases
- 98 remaining (51.6%) to be implemented

**Verification:** Run `./check_aliases.sh` to verify all functions have SLALIB aliases

## Overview

**slagofa** follows the same approach as [PAL (Positional Astronomy
Library)](https://github.com/Starlink/pal) - but better - providing a familiar
SLALIB API while using modern IAU standards underneath via GoFA/SOFA.

- **100% SLALIB API compatibility** where possible
- **Idiomatic Go** design with SLALIB aliases for drop-in replacement
- **Modern IAU standards** via GoFA (based on SOFA)
- **High performance** - zero allocations on vector/matrix operations
- **Well tested** - using exact SLALIB test vectors

Tests follow this pattern:
// From SLALIB test suite (sla_test.cc line 1152-1157)
input := 5.43 // Exact SLALIB input
expected := 5.429855087793875 // Exact SLALIB output
tolerance := 1.0e-12 // Same as SLALIB

## Architecture

```
SLALIB (Fortran)
    ↓
PAL (C/SOFA)        →  slagofa (Go/GoFA)
    ↓                       ↓
SOFA (C)            →  GoFA (Go)
```

## Installation

```bash
go get github.com/uhawaii/slagofa
```

## Usage

### Go-Idiomatic API

```go
import "github.com/uhawaii/slagofa"

// Normalize angle to ±π range
angle := slagofa.NormalizeAngle(-4.0)

// Vector operations
v1 := slagofa.Vec3{1.0, 2.0, 3.0}
v2 := slagofa.Vec3{4.0, 5.0, 6.0}
dot := slagofa.DotProduct(v1, v2)
cross := slagofa.CrossProduct(v1, v2)
unit, mag := slagofa.NormalizeVector(v1)

// Angular separation
sep := slagofa.AngularSeparation(ra1, dec1, ra2, dec2)

// Atmospheric refraction (highest precision)
hm := 2000.0      // Observer height (m)
tdk := 280.0      // Temperature (K)
pmb := 800.0      // Pressure (mbar)
rh := 0.5         // Relative humidity
wl := 0.55        // Wavelength (μm, optical)
phi := 0.3        // Latitude (rad)
tlr := 0.0065     // Temperature lapse rate (K/m)
eps := 1e-8       // Precision (rad)

// Observed → true zenith distance
zobs := 1.2  // radians
ref := slagofa.Refro(zobs, hm, tdk, pmb, rh, wl, phi, tlr, eps)
ztrue := zobs + ref

// Or fit refraction constants for multiple calculations
refa, refb := slagofa.Refco(hm, tdk, pmb, rh, wl, phi, tlr, eps)
zr := slagofa.Refz(ztrue, refa, refb)  // true → observed (fast)
```

### SLALIB-Compatible API

For direct SLALIB compatibility, use the lowercase aliases:

```go
import sla "github.com/uhawaii/slagofa"

// Exact SLALIB API
angle := sla.Drange(-4.0)        // sla_DRANGE
dot := sla.Dvdv(v1, v2)          // sla_DVDV
cross := sla.Dvxv(v1, v2)        // sla_DVXV
unit, mag := sla.Dvn(v1)         // sla_DVN
sep := sla.Dsep(ra1, dec1, ra2, dec2)  // sla_DSEP
```

## Function Categories

### Implemented Functions

- ✅ **Vector Operations** - Dot product, cross product, normalize, magnitude
- ✅ **Matrix Operations** - Matrix multiply, matrix-vector multiply, inversion, SVD
- ✅ **Angle Operations** - Normalization, formatting (radians ↔ DMS/HMS)
- ✅ **Angular Separation** - Between vectors and spherical coordinates
- ✅ **Coordinate Conversions** - Spherical ↔ Cartesian, coordinate fitting
- ✅ **Calendar/Time** - Julian date conversions, epochs, GMST, Delta T
- ✅ **Atmospheric Refraction** - Full integration (Refco, Refcoq, Refro, Refv, Airmas, Pa)
- ✅ **Utility Functions** - Random numbers, combinations, permutations
- ✅ **String I/O** - DMS/HMS string parsing and formatting

### Planned Functions

- **Coordinate Transformations** - Equatorial ↔ Horizon, Galactic, Ecliptic
- **Proper Motion** - Star catalog updates, precession
- **Observation Planning** - Rise/transit/set times
- **Map Projections** - Gnomonic, zenithal projections

See [MAPPING.md](MAPPING.md) for complete SLALIB → GoFA function mapping.

## Design Principles

### 1. Dual API

Every function has two names:

- **PascalCase** (Go-idiomatic): `NormalizeAngle`, `DotProduct`
- **Lowercase SLALIB**: `Drange`, `Dvdv` (aliases to Go functions)

### 2. Type Safety

Separate types for different precisions:

```go
type Vec3 [3]float64      // Double precision
type Vec3_32 [3]float32   // Single precision
```

### 3. Zero Allocations

All vector/matrix operations use value types and avoid heap allocations:

```go
func DotProduct(a, b Vec3) float64  // No allocations
```

### 4. Exact Test Vectors

All tests use exact values from SLALIB test suite with matching tolerances:

- Double precision: 1.0e-12
- Single precision: 1.0e-6

## Performance

Benchmark results (Apple M-series, Go 1.21):

```
Operation              ns/op    allocs/op
DotProduct             0.24     0
CrossProduct           0.24     0
NormalizeVector        2.71     0
NormalizeAngle         2.99     0
AngularSeparation      4.50     0
```

## Implementation Notes

### Atmospheric Refraction: Direct SLALIB Port

**Important:** The atmospheric refraction functions (`Refco`, `Refcoq`, `Refro`, `Refv`)
are **direct ports from SLALIB/PAL** rather than wrappers around GoFA.

**Why?** GoFA provides only a simplified refraction model (`gofa.Refco`) based on
empirical formulas. For **highest precision refraction** (0.1 arcsec accuracy), we
need SLALIB's full Hohenkerk & Sinclair atmospheric integration method.

**What we ported:**

- Full numerical integration through atmospheric layers (troposphere/stratosphere)
- Simpson's Rule adaptive integration (up to 16384 strips)
- Complete atmospheric modeling with temperature/pressure/humidity profiles
- ~500 lines of direct SLALIB/PAL translation to Go

**Accuracy achieved:**

- `Refco`: 0.5 arcsec for ZD < 80°, 0.01 arcsec for ZD < 60°
- `Refro`: 0.1 arcsec (vs. 35% error with simplified model)
- Exact match to SLALIB test vectors within 1e-8 to 1e-12

**Performance:** Zero heap allocations maintained despite complex integration.

This divergence from the "wrap GoFA" approach was necessary to maintain **full
SLALIB fidelity** for precision astronomy applications requiring sub-arcsecond
refraction corrections.

### Other Differences from SLALIB

While maintaining API compatibility, slagofa uses modern IAU standards where appropriate:

1. **Precession/Nutation** - Uses IAU 2006/2000A models (vs. IAU 1976)
2. **Time scales** - Modern TT/UT1 handling
3. **Coordinate systems** - ICRS instead of FK5 where appropriate

These changes provide more accurate results while maintaining the familiar SLALIB interface.

## Contributing

Contributions welcome! Priority areas:

1. Additional coordinate transformation functions
2. Refraction and atmospheric models
3. Proper motion functions
4. Additional test coverage

## References

- **SLALIB**: [Starlink SLALIB](http://star-www.rl.ac.uk/docs/sun67.htx/sun67.html)
- **PAL**: [Positional Astronomy Library](https://github.com/Starlink/pal)
- **GoFA**: [Golang SOFA](https://github.com/hebl/gofa)
- **SOFA**: [IAU SOFA](http://www.iausofa.org/)

## License

MIT License (see [LICENSE](LICENSE))

## Citation

If you use `slagofa` in your research, please consider citing it and:

- The GoFA library
- The IAU SOFA library
- The original SLALIB: Wallace, P.T., "SLALIB - A Library of Subroutines"

## Testing & Validation

We track **deviation from SLALIB test vectors** for every test to ensure accuracy.

### Test Accuracy Tracking

Every test records its deviation from the original SLALIB test suite:

```go
AssertAlmostEqual(t, "NormalizeAngle", "test description",
    FormatSLALIBSource("sla_test.cc", 514),
    result, expected, tolerance)
```

This provides:

- **Perfect match tracking** - Identifies zero-deviation tests
- **Deviation analysis** - Shows largest deviations even in passing tests
- **Source attribution** - Every test cites sla_test.f or sla_test.cc line numbers
- **Per-function accuracy** - Summary statistics by function

### Running Tests with Deviation Reports

```bash
# Run all tests with deviation tracking
go test -v

# Generate deviation summary
go test -v -run TestDeviationReport

# Full test validation report
./test_validation_report.sh
```

**See [TESTING.md](TESTING.md) for complete testing guide.**
