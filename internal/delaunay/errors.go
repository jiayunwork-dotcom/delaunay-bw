package delaunay

import "errors"

// Errors reported by the triangulation kernel. Each message names the rule
// that was violated so the HTTP layer can surface it verbatim to the client.
var (
	// ErrTooFewPoints mirrors geom.ErrTooFewPoints for callers that prefer to
	// stay within this package.
	ErrTooFewPoints = errors.New("delaunay: at least three points are required")

	// ErrDegenerateTriangle is returned when the kernel encounters a triangle
	// whose circumcircle cannot be computed, meaning its vertices are
	// collinear or coincident. This should not happen after validation; it is
	// a defensive guard.
	ErrDegenerateTriangle = errors.New("delaunay: degenerate triangle has no circumcircle")

	// ErrInvalidTolerance is returned when the configured relative tolerance
	// is not a positive finite number.
	ErrInvalidTolerance = errors.New("delaunay: tolerance must be a positive finite number")

	// ErrEmptyCircleViolated is returned by ValidateDelaunay when some point
	// lies strictly inside another triangle's circumcircle.
	ErrEmptyCircleViolated = errors.New("delaunay: empty-circle property violated")
)
