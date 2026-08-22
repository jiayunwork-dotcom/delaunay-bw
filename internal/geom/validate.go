package geom

import (
	"errors"
	"fmt"
	"math"
)

// Sentinel errors returned by ValidatePoints. They carry stable, descriptive
// messages so the HTTP layer and the CLI can surface the exact failing rule to
// the user.
var (
	// ErrTooFewPoints is reported when fewer than three points are supplied;
	// a triangle needs at least three vertices.
	ErrTooFewPoints = errors.New("geom: at least three points are required")

	// ErrNonFinite is reported when any coordinate is NaN or ±Inf. Such
	// points make every geometric predicate meaningless.
	ErrNonFinite = errors.New("geom: coordinates must be finite (no NaN or Inf)")

	// ErrDuplicatePoint is reported when two points coincide within tolerance.
	// Duplicated points collapse the cavity loop and must be rejected before
	// triangulation.
	ErrDuplicatePoint = errors.New("geom: duplicate point detected")

	// ErrCollinearSet is reported when all points lie on a single line, so no
	// triangle with non-zero area can be formed.
	ErrCollinearSet = errors.New("geom: all points are collinear; they span no area")
)

// ValidatePoints checks the preconditions required by the triangulator:
// at least three points, finite coordinates, no duplicates within the
// relative tolerance, and not all points collinear.
//
// The returned error is one of the sentinel errors above, so callers can
// compare with errors.Is. A nil error means the point set is ready for
// Bowyer–Watson insertion.
func ValidatePoints(pts []Point, tol float64) error {
	if len(pts) < 3 {
		return ErrTooFewPoints
	}
	scale := ScaleForTolerance(pts)
	for i, p := range pts {
		if math.IsNaN(p.X) || math.IsNaN(p.Y) || math.IsInf(p.X, 0) || math.IsInf(p.Y, 0) {
			return fmt.Errorf("%w (point %d: %s)", ErrNonFinite, i, p.String())
		}
	}
	dupTol := tol * scale
	for i := 0; i < len(pts); i++ {
		for j := i + 1; j < len(pts); j++ {
			if pts[i].Equal(pts[j], dupTol) {
				return fmt.Errorf("%w (points %d and %d both at %s)",
					ErrDuplicatePoint, i, j, pts[i].String())
			}
		}
	}
	if AllCollinear(pts, tol) {
		return ErrCollinearSet
	}
	return nil
}

// ValidateFiniteOnly runs only the finiteness and minimum-count checks,
// skipping the O(n^2) duplicate scan. It is useful for callers that already
// know the point set is distinct.
func ValidateFiniteOnly(pts []Point, tol float64) error {
	if len(pts) < 3 {
		return ErrTooFewPoints
	}
	for i, p := range pts {
		if math.IsNaN(p.X) || math.IsNaN(p.Y) || math.IsInf(p.X, 0) || math.IsInf(p.Y, 0) {
			return fmt.Errorf("%w (point %d: %s)", ErrNonFinite, i, p.String())
		}
	}
	return nil
}

// IndexError wraps a failing input index together with the root cause, making
// API diagnostics precise about which point violated which rule.
type IndexError struct {
	Index int
	Err   error
}

// Error implements the error interface.
func (e *IndexError) Error() string {
	return fmt.Sprintf("geom: point %d: %v", e.Index, e.Err)
}

// Unwrap exposes the wrapped sentinel for errors.Is.
func (e *IndexError) Unwrap() error {
	return e.Err
}
