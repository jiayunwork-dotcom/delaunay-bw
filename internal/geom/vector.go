package geom

// Cross returns the 2D cross product of the vectors a and b,
// a.X*b.Y - a.Y*b.X. Its sign tells the orientation of b relative to a.
func Cross(a, b Point) float64 {
	return a.X*b.Y - a.Y*b.X
}

// Dot returns the dot product of the vectors a and b.
func Dot(a, b Point) float64 {
	return a.X*b.X + a.Y*b.Y
}

// CrossFrom returns the signed area parameter of the triangle (o, a, b):
// Cross(a-o, b-o). It is positive when b lies counter-clockwise from a as
// seen from o, negative when clockwise, and (near) zero when the three points
// are collinear.
func CrossFrom(o, a, b Point) float64 {
	return (a.X-o.X)*(b.Y-o.Y) - (a.Y-o.Y)*(b.X-o.X)
}

// Orient returns +1 when the triangle (a, b, c) is oriented counter-clockwise,
// -1 when it is clockwise, and 0 when the points are collinear within the
// relative tolerance tol. The tolerance is scaled by the triangle's coordinate
// magnitude so the result is stable across large and small point sets.
func Orient(a, b, c Point, tol float64) int {
	cross := CrossFrom(a, b, c)
	scale := abs(a.X) + abs(a.Y) + abs(b.X) + abs(b.Y) + abs(c.X) + abs(c.Y)
	eps := tol * (scale + 1)
	switch {
	case cross > eps:
		return 1
	case cross < -eps:
		return -1
	default:
		return 0
	}
}

// IsCCW reports whether the triangle (a, b, c) is oriented counter-clockwise
// within the given relative tolerance.
func IsCCW(a, b, c Point, tol float64) bool {
	return Orient(a, b, c, tol) > 0
}

// IsCW reports whether the triangle (a, b, c) is oriented clockwise within the
// given relative tolerance.
func IsCW(a, b, c Point, tol float64) bool {
	return Orient(a, b, c, tol) < 0
}
