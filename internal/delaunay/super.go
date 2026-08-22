package delaunay

import "delaunay-bw/internal/geom"

// superScale is how many times the bounding-box diagonal the super-triangle
// extends beyond the box center. The triangle must be large enough that every
// boundary-adjacent insertion point falls inside the circumcircle of the
// triangle covering its visible hull arc. A factor of 1000 leaves a wide
// safety margin over the empirical threshold (~50x) for the shipped examples
// while keeping incircle determinants (~10^15) far below float64 overflow.
const superScale = 1000.0

// SuperTriangle returns three points that form a large counter-clockwise
// triangle strictly containing every point of the bounding box. The triangle
// is centered on the box center and sized from its diagonal.
//
// The returned slice always has exactly three points. The triangulator
// appends them to the input points and uses the appended indices as the
// virtual super-triangle vertices; the finalize step strips any triangle that
// references them.
func SuperTriangle(b geom.Bounds) []geom.Point {
	c := b.Center()
	d := b.Diagonal()
	if d <= 0 {
		d = 1
	}
	ext := superScale * d
	return []geom.Point{
		{X: c.X - ext, Y: c.Y - ext},
		{X: c.X + ext, Y: c.Y - ext},
		{X: c.X, Y: c.Y + ext},
	}
}

// IsSuperVertex reports whether v is one of the three appended super-triangle
// vertex indices. Input points keep indices [0, n), so any index >= n belongs
// to the super triangle.
func IsSuperVertex(v, n int) bool {
	return v >= n
}
