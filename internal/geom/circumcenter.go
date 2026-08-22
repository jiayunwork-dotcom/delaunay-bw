package geom

import "errors"

// ErrCollinear is returned by Circumcenter when the three input points do not
// span a non-zero area, so no circle passes through them uniquely.
var ErrCollinear = errors.New("geom: circumcenter undefined for collinear points")

// Circumcenter returns the center of the circle passing through a, b and c
// using the closed-form barycentric solution. The three points must be
// non-collinear; the guard threshold scales with the coordinate magnitude so
// that large point sets do not overflow the denominator prematurely.
//
// The returned center is exact in the sense that its distance to a, b and c
// agrees within floating point error. This is the Voronoi vertex of the
// triangle in the dual diagram.
func Circumcenter(a, b, c Point) (Point, error) {
	d := 2 * (a.X*(b.Y-c.Y) + b.X*(c.Y-a.Y) + c.X*(a.Y-b.Y))
	scale := abs(a.X) + abs(a.Y) + abs(b.X) + abs(b.Y) + abs(c.X) + abs(c.Y)
	// The determinant d has magnitude ~ coordinate^2; scale the threshold
	// accordingly so the check is relative, not absolute.
	if abs(d) <= 1e-12*(scale+1)*(scale+1) {
		return Point{}, ErrCollinear
	}
	aa := a.X*a.X + a.Y*a.Y
	bb := b.X*b.X + b.Y*b.Y
	cc := c.X*c.X + c.Y*c.Y
	ux := (aa*(b.Y-c.Y) + bb*(c.Y-a.Y) + cc*(a.Y-b.Y)) / d
	uy := (aa*(c.X-b.X) + bb*(a.X-c.X) + cc*(b.X-a.X)) / d
	return Point{X: ux, Y: uy}, nil
}

// CircumradiusSq returns the squared radius of the circle through a, b and c.
// It is a convenience wrapper over Circumcenter and propagates its error for
// degenerate triangles.
func CircumradiusSq(a, b, c Point) (float64, error) {
	o, err := Circumcenter(a, b, c)
	if err != nil {
		return 0, err
	}
	return a.DistanceSq(o), nil
}

// PointInTriangle uses barycentric coordinates to test whether p lies inside
// or on the boundary of the triangle (a, b, c). The test is orientation-aware:
// it computes signed areas against a consistent winding, so it never depends
// on the input order of a, b, c.
func PointInTriangle(p, a, b, c Point, tol float64) bool {
	s1 := CrossFrom(a, b, p)
	s2 := CrossFrom(b, c, p)
	s3 := CrossFrom(c, a, p)
	// All three signed areas must share the same sign (allowing near-zero for
	// points on the boundary).
	hasNeg := s1 < -tol || s2 < -tol || s3 < -tol
	hasPos := s1 > tol || s2 > tol || s3 > tol
	return !(hasNeg && hasPos)
}
