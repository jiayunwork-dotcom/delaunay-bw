package delaunay

import (
	"fmt"

	"delaunay-bw/internal/geom"
)

// IsDelaunay reports whether every triangle of the mesh satisfies the empty
// circle property: no input point lies strictly inside the circumcircle of
// any triangle. Points sitting exactly on a circle (within tolerance) are
// allowed; that is exactly the four-point-cocircular configuration the
// tolerance is designed to keep stable.
func IsDelaunay(pts []geom.Point, tris []Triangle, tol float64) bool {
	for _, t := range tris {
		ok, err := CircumcircleIsEmpty(pts, t, tol)
		if err != nil || !ok {
			return false
		}
	}
	return true
}

// ValidateDelaunay is the error-returning form of IsDelaunay. The returned
// error names the failing triangle and the offending point so the message is
// actionable in the HTTP error body.
func ValidateDelaunay(pts []geom.Point, tris []Triangle, tol float64) error {
	for _, t := range tris {
		o, err := CircumcenterOf(pts, t)
		if err != nil {
			return fmt.Errorf("%w: %v", err, t)
		}
		radiusSq := pts[t[0]].DistanceSq(o)
		excess := tol * (radiusSq + 1)
		for i, p := range pts {
			if t.HasVertex(i) {
				continue
			}
			if p.DistanceSq(o) < radiusSq-excess {
				return fmt.Errorf(
					"%w: point %d %s is inside the circumcircle of triangle %v "+
						"(radius²=%.6g, distance²=%.6g)",
					ErrEmptyCircleViolated, i, p.String(), t, radiusSq, p.DistanceSq(o))
			}
		}
	}
	return nil
}

// ValidateTopology checks the structural invariants of the mesh:
//
//   - every vertex index is a valid input index,
//   - every triangle has a positive area (non-degenerate),
//   - every interior edge is shared by exactly two triangles (a closed
//     icos-manifold patch of the convex hull).
//
// A violation means the mesh would draw cross-edges or overlap itself.
func ValidateTopology(pts []geom.Point, tris []Triangle, tol float64) error {
	for _, t := range tris {
		for _, v := range t {
			if v < 0 || v >= len(pts) {
				return fmt.Errorf("delaunay: triangle %v references out-of-range vertex %d", t, v)
			}
		}
		if geom.TriangleArea(pts[t[0]], pts[t[1]], pts[t[2]]) <= 0 {
			return fmt.Errorf("delaunay: triangle %v has zero area", t)
		}
	}

	// Count each directed edge; a valid mesh has every non-boundary edge
	// appearing in both directions exactly once.
	incidence := make(map[directedEdge]int)
	for _, t := range tris {
		for k := 0; k < 3; k++ {
			incidence[directedEdge{a: t[k], b: t[(k+1)%3]}]++
		}
	}
	for d, c := range incidence {
		rc := incidence[directedEdge{a: d.b, b: d.a}]
		if rc == 0 {
			// Boundary edge: allowed exactly once.
			if c != 1 {
				return fmt.Errorf("delaunay: boundary edge (%d-%d) seen %d times", d.a, d.b, c)
			}
			continue
		}
		if c != 1 || rc != 1 {
			return fmt.Errorf("delaunay: interior edge (%d-%d) is shared by %d/%d triangles",
				d.a, d.b, c, rc)
		}
	}
	return nil
}

// CoverageError describes a mismatch between the triangle union area and the
// convex hull area.
type CoverageError struct {
	HullArea    float64
	TriangleSum float64
	Delta       float64
}

// Error implements the error interface.
func (e *CoverageError) Error() string {
	return fmt.Sprintf(
		"delaunay: triangle union area %.6g does not match convex hull area %.6g (delta %.6g)",
		e.TriangleSum, e.HullArea, e.Delta)
}

// ValidateCoverage verifies that the mesh tiles the convex hull: the sum of
// triangle areas equals the hull area. A ghost triangle outside the hull
// (from a stripped super-triangle leak) inflates the triangle sum by a
// hull-scale amount and trips this check.
//
// The tolerance is relative to the hull area: points that lie almost exactly
// on a hull edge can make the mesh boundary differ from Andrew's collinear-
// stripped hull by a tiny sliver, and that numerical edge case is allowed.
// Real coverage errors are of order the hull area itself.
func ValidateCoverage(pts []geom.Point, tris []Triangle, hull []int, tol float64) error {
	hullArea := geom.HullArea(pts, hull)
	sum := 0.0
	for _, t := range tris {
		sum += geom.TriangleArea(pts[t[0]], pts[t[1]], pts[t[2]])
	}
	delta := hullArea - sum
	if delta < 0 {
		delta = -delta
	}
	threshold := (tol + 1e-5) * (hullArea + 1)
	if delta > threshold {
		return &CoverageError{HullArea: hullArea, TriangleSum: sum, Delta: delta}
	}
	return nil
}
