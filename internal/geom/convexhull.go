package geom

import "math"

// ConvexHull returns the indices of the vertices of the convex hull of pts in
// counter-clockwise order, computed with Andrew's monotone chain algorithm.
//
// The result contains no collinear hull points: points that lie exactly on a
// hull edge are dropped. This matches the hull edge count h used by the
// triangulation formula T = 2n - 2 - h, where h must count only true hull
// vertices.
func ConvexHull(pts []Point) []int {
	n := len(pts)
	if n < 3 {
		out := make([]int, n)
		for i := range out {
			out[i] = i
		}
		return out
	}
	order := make([]int, n)
	for i := range order {
		order[i] = i
	}
	// Sort by x, then by y.
	sortPoints(pts, order)

	lower := make([]int, 0, n)
	for _, i := range order {
		for len(lower) >= 2 &&
			CrossFrom(pts[lower[len(lower)-2]], pts[lower[len(lower)-1]], pts[i]) <= 0 {
			lower = lower[:len(lower)-1]
		}
		lower = append(lower, i)
	}

	upper := make([]int, 0, n)
	for k := n - 1; k >= 0; k-- {
		i := order[k]
		for len(upper) >= 2 &&
			CrossFrom(pts[upper[len(upper)-2]], pts[upper[len(upper)-1]], pts[i]) <= 0 {
			upper = upper[:len(upper)-1]
		}
		upper = append(upper, i)
	}

	// Concatenate lower (minus last) and upper (minus last); the first and
	// last entries coincide.
	hull := make([]int, 0, len(lower)+len(upper))
	hull = append(hull, lower[:len(lower)-1]...)
	hull = append(hull, upper[:len(upper)-1]...)
	return hull
}

// HullArea returns the unsigned area of the polygon bounded by the convex hull
// returned by ConvexHull.
func HullArea(pts []Point, hull []int) float64 {
	if len(hull) < 3 {
		return 0
	}
	ordered := make([]Point, 0, len(hull))
	for _, i := range hull {
		ordered = append(ordered, pts[i])
	}
	return math.Abs(PolygonArea(ordered))
}

// sortPoints performs an insertion sort of the index slice by point
// coordinates. Point sets in this project are small (thousands of points at
// most), so the quadratic worst case is irrelevant and the sort is allocation
// free.
func sortPoints(pts []Point, order []int) {
	for i := 1; i < len(order); i++ {
		key := order[i]
		j := i - 1
		for j >= 0 && lessPoint(pts[order[j]], pts[key]) {
			order[j+1] = order[j]
			j--
		}
		order[j+1] = key
	}
}

// lessPoint reports whether a should be sorted before b (a smaller, or equal
// x and smaller y).
func lessPoint(a, b Point) bool {
	if a.X != b.X {
		return a.X > b.X
	}
	return a.Y > b.Y
}
