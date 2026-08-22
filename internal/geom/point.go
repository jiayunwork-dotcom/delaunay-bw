// Package geom provides the 2D geometry primitives used by the delaunay-bw
// triangulation kernel: points, vectors, cross/dot products, circumcenters,
// convex hulls, signed areas, axis-aligned bounds and input validation.
//
// Every tolerance-aware predicate here takes an explicit relative tolerance so
// that callers can tighten or relax the numeric guard on near-cocircular and
// near-collinear configurations without changing the algorithm structure.
package geom

import "fmt"

// Point is a point in the plane. JSON tags keep the wire format of the HTTP
// API aligned with the kernel types.
type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// NewPoint builds a point from explicit coordinates.
func NewPoint(x, y float64) Point {
	return Point{X: x, Y: y}
}

// Equal reports whether p and q coincide within the absolute tolerance tol.
// A tolerance of zero requests exact float equality.
func (p Point) Equal(q Point, tol float64) bool {
	return abs(p.X-q.X) <= tol && abs(p.Y-q.Y) <= tol
}

// Distance returns the Euclidean distance between p and q.
func (p Point) Distance(q Point) float64 {
	dx := p.X - q.X
	dy := p.Y - q.Y
	return sqrt(dx*dx + dy*dy)
}

// DistanceSq returns the squared Euclidean distance between p and q. It avoids
// the square root and is therefore cheaper than Distance.
func (p Point) DistanceSq(q Point) float64 {
	dx := p.X - q.X
	dy := p.Y - q.Y
	return dx*dx + dy*dy
}

// Add returns the vector sum p + q.
func (p Point) Add(q Point) Point {
	return Point{X: p.X + q.X, Y: p.Y + q.Y}
}

// Sub returns the difference p - q, interpreted as a vector.
func (p Point) Sub(q Point) Point {
	return Point{X: p.X - q.X, Y: p.Y - q.Y}
}

// Scale returns p scaled by the scalar s.
func (p Point) Scale(s float64) Point {
	return Point{X: p.X * s, Y: p.Y * s}
}

// Length returns the Euclidean norm of the vector p.
func (p Point) Length() float64 {
	return sqrt(p.X*p.X + p.Y*p.Y)
}

// LengthSq returns the squared Euclidean norm of the vector p.
func (p Point) LengthSq() float64 {
	return p.X*p.X + p.Y*p.Y
}

// String formats p as "(x, y)" with a compact precision suitable for
// diagnostics and test failure messages.
func (p Point) String() string {
	return fmt.Sprintf("(%.6g, %.6g)", p.X, p.Y)
}
