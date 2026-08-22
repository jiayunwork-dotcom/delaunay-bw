package geom

// Bounds is an axis-aligned bounding box around a point set.
type Bounds struct {
	Min Point
	Max Point
}

// BoundsOf returns the axis-aligned bounding box of pts. The empty box is
// returned when pts is empty; the returned box always satisfies Min <= Max.
func BoundsOf(pts []Point) Bounds {
	if len(pts) == 0 {
		return Bounds{}
	}
	minX, minY := pts[0].X, pts[0].Y
	maxX, maxY := pts[0].X, pts[0].Y
	for _, p := range pts[1:] {
		if p.X < minX {
			minX = p.X
		}
		if p.X > maxX {
			maxX = p.X
		}
		if p.Y < minY {
			minY = p.Y
		}
		if p.Y > maxY {
			maxY = p.Y
		}
	}
	return Bounds{Min: Point{X: minX, Y: minY}, Max: Point{X: maxX, Y: maxY}}
}

// Width returns the horizontal extent of the box.
func (b Bounds) Width() float64 {
	return b.Max.X - b.Min.X
}

// Height returns the vertical extent of the box.
func (b Bounds) Height() float64 {
	return b.Max.Y - b.Min.Y
}

// Diagonal returns the distance between the two corners of the box.
func (b Bounds) Diagonal() float64 {
	return b.Min.Distance(b.Max)
}

// Center returns the midpoint of the box.
func (b Bounds) Center() Point {
	return Point{
		X: 0.5 * (b.Min.X + b.Max.X),
		Y: 0.5 * (b.Min.Y + b.Max.Y),
	}
}

// Contains reports whether the point p lies inside the box, treating the
// boundary as inclusive.
func (b Bounds) Contains(p Point) bool {
	return p.X >= b.Min.X && p.X <= b.Max.X &&
		p.Y >= b.Min.Y && p.Y <= b.Max.Y
}

// Pad returns a copy of the box expanded by m on every side.
func (b Bounds) Pad(m float64) Bounds {
	return Bounds{
		Min: Point{X: b.Min.X - m, Y: b.Min.Y - m},
		Max: Point{X: b.Max.X + m, Y: b.Max.Y + m},
	}
}

// LargestSide returns the larger of the width and the height of the box.
func (b Bounds) LargestSide() float64 {
	w := b.Width()
	h := b.Height()
	if w > h {
		return w
	}
	return h
}

// AspectRatio returns height/width (width/height when width is zero),
// a guard that never divides by zero.
func (b Bounds) AspectRatio() float64 {
	w := b.Width()
	if w == 0 {
		return 1
	}
	return b.Height() / w
}
