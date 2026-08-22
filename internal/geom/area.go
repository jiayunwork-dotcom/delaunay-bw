package geom

import "math"

// SignedAreaTwice returns twice the signed area of the triangle (a, b, c).
// It is positive for counter-clockwise triples and negative for clockwise
// triples. The factor of two keeps the computation free of divisions.
func SignedAreaTwice(a, b, c Point) float64 {
	return CrossFrom(a, b, c)
}

// TriangleArea returns the unsigned area of the triangle (a, b, c).
func TriangleArea(a, b, c Point) float64 {
	return 0.5 * abs(CrossFrom(a, b, c))
}

// PolygonArea returns the signed area of a simple polygon whose vertices are
// given in order. Positive for counter-clockwise order. Uses the shoelace
// formula.
func PolygonArea(pts []Point) float64 {
	return 0.5 * PolygonAreaTwice(pts)
}

// PolygonAreaTwice returns twice the signed area of a simple polygon given by
// ordered vertices.
func PolygonAreaTwice(pts []Point) float64 {
	if len(pts) < 3 {
		return 0
	}
	var sum float64
	for i := range pts {
		j := (i + 1) % len(pts)
		sum += pts[i].X*pts[j].Y - pts[j].X*pts[i].Y
	}
	return sum
}

// PolygonAreaFromIndices returns the unsigned area of the polygon obtained by
// walking the index slice idx in order through the point slice pts.
func PolygonAreaFromIndices(pts []Point, idx []int) float64 {
	ordered := make([]Point, 0, len(idx))
	for _, i := range idx {
		ordered = append(ordered, pts[i])
	}
	return abs(PolygonArea(ordered))
}

// TriangleAreaSum returns the sum of the unsigned areas of the triangles whose
// vertex indices are listed in tris.
func TriangleAreaSum(pts []Point, tris [][3]int) float64 {
	var sum float64
	for _, t := range tris {
		sum += TriangleArea(pts[t[0]], pts[t[1]], pts[t[2]])
	}
	return sum
}

// PointInConvexPolygon reports whether p lies inside (or on the boundary of)
// the convex polygon given by ordered vertices. The order may be clockwise or
// counter-clockwise; the test works in either case because it only compares
// signs of successive cross products.
func PointInConvexPolygon(p Point, poly []Point, tol float64) bool {
	if len(poly) < 3 {
		return false
	}
	var ref float64
	first := true
	for i := range poly {
		j := (i + 1) % len(poly)
		c := CrossFrom(poly[i], poly[j], p)
		// A zero (or near-zero) cross product means p lies on the supporting
		// line of an edge; treat it as inside.
		if first {
			ref = c
			first = false
			if abs(c) <= tol {
				continue
			}
		}
		if c < -tol && ref > tol || c > tol && ref < -tol {
			return false
		}
	}
	return true
}

// ScaleForTolerance returns a characteristic length for a point set: the
// diagonal of its bounding box, or 1 if every point coincides. Relative
// tolerances multiply this scale.
func ScaleForTolerance(pts []Point) float64 {
	if len(pts) == 0 {
		return 1
	}
	b := BoundsOf(pts)
	d := b.Max.Sub(b.Min)
	s := math.Sqrt(d.X*d.X + d.Y*d.Y)
	if s == 0 {
		return 1
	}
	return s
}
