package geom

import "math"

// AreCollinear reports whether the three points a, b, c are collinear within
// the relative tolerance tol. The tolerance is applied to the area term after
// scaling by the coordinate magnitude, so a set of points around (1e6, 1e6)
// is judged with the same relative precision as one around the origin.
func AreCollinear(a, b, c Point, tol float64) bool {
	cross := CrossFrom(a, b, c)
	scale := abs(a.X) + abs(a.Y) + abs(b.X) + abs(b.Y) + abs(c.X) + abs(c.Y)
	return abs(cross) <= tol*(scale+1)*(scale+1)
}

// AllCollinear reports whether every point in pts lies on a single straight
// line within the relative tolerance tol. It is used as a pre-condition check
// before triangulation: a collinear point set spans no area.
func AllCollinear(pts []Point, tol float64) bool {
	n := len(pts)
	if n < 3 {
		// A pair of points is trivially collinear, but the caller decides
		// whether that is an error via its own count check.
		return n < 3
	}
	// Pick the two points that are farthest apart as the reference line so the
	// test is not fooled by a lucky first pair.
	p0, p1 := farthestPair(pts)
	for i := range pts {
		if i == p0 || i == p1 {
			continue
		}
		if !AreCollinear(pts[p0], pts[p1], pts[i], tol) {
			return false
		}
	}
	return true
}

// farthestPair returns the indices of the two points with maximal squared
// distance. Choosing the longest reference segment makes the collinearity
// check maximally stable.
func farthestPair(pts []Point) (int, int) {
	bestI, bestJ := 0, 1
	best := pts[0].DistanceSq(pts[1])
	for i := 0; i < len(pts); i++ {
		for j := i + 1; j < len(pts); j++ {
			d := pts[i].DistanceSq(pts[j])
			if d > best {
				best = d
				bestI, bestJ = i, j
			}
		}
	}
	return bestI, bestJ
}

// Tolerance0 is the default relative tolerance used by predicates that take
// an explicit tol argument when the caller wants the standard guard.
const Tolerance0 = 1e-9

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

func sqrt(x float64) float64 {
	return math.Sqrt(x)
}
