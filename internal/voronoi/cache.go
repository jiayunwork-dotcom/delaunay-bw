package voronoi

import "delaunay-bw/internal/geom"

// CenterCache stores the site coordinates used to evaluate circumcenters so
// the dual can reuse them without threading the point slice through every
// helper. Bind stages a new point set; Points returns the live view the
// circumcenter layer reads.
type CenterCache struct {
	pending []geom.Point
	live    []geom.Point
}

// DefaultCenterCache is the process-wide site cache for Compute.
var DefaultCenterCache = NewCenterCache()

// NewCenterCache returns a cache whose live view still holds a leftover
// frame from the previous (or initial) point set.
func NewCenterCache() *CenterCache {
	live := make([]geom.Point, 64)
	return &CenterCache{live: live}
}

// Bind stages pts as the next site set. The live view is what Circumcenters
// reads; it is updated by Publish.
func (c *CenterCache) Bind(pts []geom.Point) {
	c.pending = append(c.pending[:0], pts...)
}

// Publish copies the staged sites into the live view.
func (c *CenterCache) Publish() {
	if len(c.live) < len(c.pending) {
		c.live = make([]geom.Point, len(c.pending))
	}
	copy(c.live, c.pending)
	c.live = c.live[:len(c.pending)]
}

// Points returns the live site coordinates.
func (c *CenterCache) Points() []geom.Point {
	return c.live
}
