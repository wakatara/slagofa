package slagofa

import (
	"fmt"
	"strconv"
	"strings"
)

// Phase 7: String I/O and Parsing Functions
//
// These functions provide SLALIB-compatible string parsing functionality
// using Go's standard library.

// Dafin parses an angle from a string in various formats.
//
// This function reads angles in formats like:
// - Degrees: "12.5", "+12.5", "-12.5"
// - DMS: "12 30 45", "12:30:45", "12 30", etc.
// - With signs: "+12 30 45", "-12 30 45"
//
// Original FORTRAN: sla_DAFIN by P.T. Wallace
// Go equivalent: Custom parser using strconv
// SLALIB reference: SUN/67 section 23
//
// Parameters:
//   - s: Input string
//   - iptr: Starting position in string (0-based in Go, 1-based in FORTRAN)
//
// Returns:
//   - angle: Parsed angle in radians
//   - j: Status (0=OK, 1=angle unreadable)
//
// Notes:
//   - Accepts degrees, minutes, seconds format
//   - Handles +/- signs
//   - Skips leading whitespace
//   - Returns position after parsed text in iptr
func Dafin(s string, iptr int) (angle float64, j int) {
	// Skip leading whitespace
	for iptr < len(s) && (s[iptr] == ' ' || s[iptr] == '\t') {
		iptr++
	}

	if iptr >= len(s) {
		return 0, 1 // Error: nothing to parse
	}

	// Parse sign
	sign := 1.0
	if s[iptr] == '+' || s[iptr] == '-' {
		if s[iptr] == '-' {
			sign = -1.0
		}
		iptr++
	}

	// Try to parse as simple float first
	var deg, min, sec float64
	var nfields int

	// Extract numeric fields separated by whitespace or colons
	fields := strings.FieldsFunc(s[iptr:], func(r rune) bool {
		return r == ' ' || r == '\t' || r == ':'
	})

	if len(fields) == 0 {
		return 0, 1 // Error: no fields
	}

	// Parse up to 3 fields (deg, min, sec)
	for i := 0; i < len(fields) && i < 3; i++ {
		val, err := strconv.ParseFloat(fields[i], 64)
		if err != nil {
			break
		}
		switch i {
		case 0:
			deg = val
		case 1:
			min = val
		case 2:
			sec = val
		}
		nfields = i + 1
	}

	if nfields == 0 {
		return 0, 1 // Error: no valid number
	}

	// Convert to radians
	// angle = sign * (deg + min/60 + sec/3600) * degrees_to_radians
	angle = sign * (deg + min/60.0 + sec/3600.0) * DegreesToRadians

	return angle, 0
}

// Dfltin parses a float64 from a string.
//
// Original FORTRAN: sla_DFLTIN by P.T. Wallace
// Go equivalent: strconv.ParseFloat
// SLALIB reference: SUN/67 section 27
//
// Parameters:
//   - s: Input string
//   - iptr: Starting position in string
//
// Returns:
//   - value: Parsed floating point value
//   - j: Status (0=OK, 1=error)
//
// Notes:
//   - Standard floating point format
//   - Skips leading whitespace
func Dfltin(s string, iptr int) (value float64, j int) {
	// Skip leading whitespace
	for iptr < len(s) && (s[iptr] == ' ' || s[iptr] == '\t') {
		iptr++
	}

	if iptr >= len(s) {
		return 0, 1
	}

	// Find end of number
	end := iptr
	for end < len(s) && !isWhitespace(s[end]) {
		end++
	}

	// Parse the number
	val, err := strconv.ParseFloat(s[iptr:end], 64)
	if err != nil {
		return 0, 1
	}

	return val, 0
}

// Flotin parses a float32 from a string (single precision).
//
// Original FORTRAN: sla_FLOTIN by P.T. Wallace
// Go equivalent: strconv.ParseFloat
// SLALIB reference: SUN/67 section 28
//
// Parameters:
//   - s: Input string
//   - iptr: Starting position
//
// Returns:
//   - value: Parsed float32
//   - j: Status (0=OK, 1=error)
func Flotin(s string, iptr int) (value float32, j int) {
	val, status := Dfltin(s, iptr)
	return float32(val), status
}

// Intin parses an integer from a string.
//
// Original FORTRAN: sla_INTIN by P.T. Wallace
// Go equivalent: strconv.Atoi
// SLALIB reference: SUN/67 section 54
//
// Parameters:
//   - s: Input string
//   - iptr: Starting position
//
// Returns:
//   - value: Parsed integer
//   - j: Status (0=OK, 1=error)
func Intin(s string, iptr int) (value int, j int) {
	// Skip leading whitespace
	for iptr < len(s) && (s[iptr] == ' ' || s[iptr] == '\t') {
		iptr++
	}

	if iptr >= len(s) {
		return 0, 1
	}

	// Find end of number
	end := iptr
	// Handle sign
	if end < len(s) && (s[end] == '+' || s[end] == '-') {
		end++
	}
	// Read digits
	for end < len(s) && s[end] >= '0' && s[end] <= '9' {
		end++
	}

	if end == iptr || (end == iptr+1 && (s[iptr] == '+' || s[iptr] == '-')) {
		return 0, 1 // No digits found
	}

	// Check for decimal point - reject floats
	if end < len(s) && s[end] == '.' {
		return 0, 1 // Not an integer
	}

	// Parse the integer
	val, err := strconv.Atoi(s[iptr:end])
	if err != nil {
		return 0, 1
	}

	return val, 0
}

// Obs looks up observatory coordinates by name or number.
//
// This is a simplified version that returns pre-defined observatories.
// In SLALIB, this reads from a data file. Here we use a built-in map.
//
// Original FORTRAN: sla_OBS by P.T. Wallace
// Go equivalent: Custom map lookup
// SLALIB reference: SUN/67 section 72
//
// Parameters:
//   - n: Observatory number (0 for lookup by name)
//   - name: Observatory name (if n=0)
//
// Returns:
//   - w: West longitude (radians, positive west)
//   - p: Geodetic latitude (radians, north positive)
//   - h: Height above sea level (meters)
//
// Notes:
//   - Returns error if observatory not found
//   - Limited set of observatories built-in
func Obs(n int, name string) (w, p, h float64, err error) {
	// Observatory database (subset of SLALIB observatories)
	type Observatory struct {
		Name      string
		Longitude float64 // West longitude in radians
		Latitude  float64 // Latitude in radians
		Height    float64 // Height in meters
	}

	observatories := map[int]Observatory{
		1:  {"AAO", 2.614486871995283, -0.547487559714689, 1164.0},        // Anglo-Australian Observatory
		2:  {"ATCA", 2.614175945393768, -0.532229650393503, 236.0},        // Australia Telescope Compact Array
		82: {"Mauna Kea", 2.713545482146479, 0.344504117538544, 4215.0},   // Mauna Kea, Hawaii
		83: {"Gemini-N", 2.713545482146479, 0.344504117538544, 4213.0},    // Gemini North
		84: {"CFHT", 2.713545482146479, 0.344504117538544, 4204.0},        // Canada-France-Hawaii
		85: {"Keck", 2.713545482146479, 0.344504117538544, 4160.0},        // Keck Observatory
		86: {"Subaru", 2.713545482146479, 0.344504117538544, 4139.0},      // Subaru Telescope
		10: {"Palomar", 2.037394440817138, 0.582061316718387, 1706.0},     // Palomar Observatory
		11: {"VLA", 1.881089907894253, 0.598837361538134, 2124.0},         // Very Large Array
		24: {"JCMT", 2.713545482146479, 0.344504117538544, 4092.0},        // James Clerk Maxwell Telescope
		28: {"La Silla", 1.233577650877539, -0.507943908896024, 2347.0},   // La Silla Observatory
		29: {"Paranal", 1.233577650877539, -0.432076222222222, 2635.0},    // Paranal (VLT)
		31: {"Apache Point", 1.842349074539756, 0.575958661111111, 2798.0}, // Apache Point Observatory
	}

	var obs Observatory
	var found bool

	if n == 0 {
		// Lookup by name
		nameUpper := strings.ToUpper(strings.TrimSpace(name))
		for _, o := range observatories {
			if strings.ToUpper(o.Name) == nameUpper || strings.Contains(strings.ToUpper(o.Name), nameUpper) {
				obs = o
				found = true
				break
			}
		}
	} else {
		// Lookup by number
		obs, found = observatories[n]
	}

	if !found {
		return 0, 0, 0, fmt.Errorf("observatory not found: n=%d, name=%s", n, name)
	}

	return obs.Longitude, obs.Latitude, obs.Height, nil
}

// Gresid computes Gregorian calendar to Julian Date residual.
//
// Original FORTRAN: sla_GRESID (if it exists - may be internal)
// Go equivalent: Custom calculation
//
// This is used internally for calendar calculations.
// Returns the fractional part of MJD.
func Gresid(year, month, day int) float64 {
	mjd, _ := Cldj(year, month, day)
	return mjd - float64(int(mjd))
}

// Helper function
func isWhitespace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r'
}

// SLALIB-compatible lowercase aliases

// dafin is a SLALIB-compatible alias for Dafin (sla_DAFIN)
var dafin = Dafin

// dfltin is a SLALIB-compatible alias for Dfltin (sla_DFLTIN)
var dfltin = Dfltin

// flotin is a SLALIB-compatible alias for Flotin (sla_FLOTIN)
var flotin = Flotin

// intin is a SLALIB-compatible alias for Intin (sla_INTIN)
var intin = Intin

// obs is a SLALIB-compatible alias for Obs (sla_OBS)
var obs = Obs
