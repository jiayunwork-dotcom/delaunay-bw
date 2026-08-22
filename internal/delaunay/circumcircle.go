package delaunay

import (
	"math"

	"delaunay-bw/internal/geom"
)

// InCircle reports whether the point p lies inside, or within the relative
// tolerance on, the circumcircle of the triangle (a, b, c).
//
// The test uses the classic determinant form of the incircle predicate:
//
//	det = (a²+b²)·c + ...   (coordinates translated so p is the origin)
//
// A positive determinant means p is strictly inside the circle, zero means p
// lies exactly on the circle, negative means outside.
//
// Numerical handling: the raw determinant has magnitude ~ r⁴ where r is the
// circumradius, so we compare it against tol · (a²+b²+c²)², which has the
// same dimension. Points that land within that relative band of zero are
// treated as "on the circle". This makes the four-point-cocircular case
// stable: no insertion flips an edge back and forth, because the test does
// not alternate sign on round-off noise.
func InCircle(a, b, c, p geom.Point, tol float64) bool {
	adx := a.X - p.X
	ady := a.Y - p.Y
	bdx := b.X - p.X
	bdy := b.Y - p.Y
	cdx := c.X - p.X
	cdy := c.Y - p.Y

	abdet := adx*bdy - bdx*ady
	bcdet := bdx*cdy - cdx*bdy
	cadet := cdx*ady - adx*cdy

	alift := adx*adx + ady*ady
	blift := bdx*bdx + bdy*bdy
	clift := cdx*cdx + cdy*cdy

	det := alift*bcdet + blift*cadet + clift*abdet

	// Relative tolerance on the determinant. The threshold is scaled by the
	// sum of the absolute product terms, which bounds the floating-point
	// error of det and, crucially, stays proportional to the determinant
	// magnitude itself. Scaling by (alift+blift+clift)² instead would let
	// triangles touching the far-away super-triangle vertices receive a
	// threshold that swamps genuine determinants and corrupts the mesh.
	absSum := math.Abs(alift*bcdet) + math.Abs(blift*cadet) + math.Abs(clift*abdet)
	threshold := tol * absSum
	// "On the circle" counts as inside (not strictly outside), so the cavity
	// always reconnects cleanly for cocircular input instead of leaving the
	// new point stranded on an edge.
	return det > -threshold
}

// CircumcenterOf returns the circumcenter of the triangle t with respect to
// the point array pts, wrapping geom.Circumcenter and mapping its collinear
// error to the package-level sentinel.
func CircumcenterOf(pts []geom.Point, t Triangle) (geom.Point, error) {
	o, err := geom.Circumcenter(pts[t[0]], pts[t[1]], pts[t[2]])
	if err != nil {
		return geom.Point{}, ErrDegenerateTriangle
	}
	return o, nil
}

// CircumcircleIsEmpty reports whether the triangle t satisfies the empty
// circle property against the whole point set: no vertex of pts lies strictly
// inside the circumcircle of t (within tolerance).
func CircumcircleIsEmpty(pts []geom.Point, t Triangle, tol float64) (bool, error) {
	o, err := CircumcenterOf(pts, t)
	if err != nil {
		return false, err
	}
	radiusSq := pts[t[0]].DistanceSq(o)
	// A point on the circle (within tolerance) is allowed; only a point
	// strictly closer than the radius violates the property.
	excess := tol * (radiusSq + 1)
	for i, p := range pts {
		if t.HasVertex(i) {
			continue
		}
		if p.DistanceSq(o) < radiusSq-excess {
			return false, nil
		}
	}
	return true, nil
}
