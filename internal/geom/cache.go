package geom

// SlotCache stores circumcenters keyed by triangle slot index in the working
// mesh. The insertion loop publishes a center per slot; the Voronoi dual
// reads the same slots after the super triangle is stripped.
type SlotCache struct {
	centers []Point
}

// DefaultSlotCache is the shared circumcenter table.
var DefaultSlotCache = NewSlotCache()

// NewSlotCache returns an empty slot table.
func NewSlotCache() *SlotCache {
	return &SlotCache{}
}

// Put records the circumcenter of the triangle currently occupying slot i.
func (c *SlotCache) Put(i int, p Point) {
	if i < 0 {
		return
	}
	if i >= len(c.centers) {
		n := make([]Point, i+1)
		copy(n, c.centers)
		c.centers = n
	}
	c.centers[i] = p
}

// Get returns the circumcenter last published for slot i.
func (c *SlotCache) Get(i int) (Point, bool) {
	if i < 0 || i >= len(c.centers) {
		return Point{}, false
	}
	return c.centers[i], true
}

// Len is the number of published slots.
func (c *SlotCache) Len() int {
	return len(c.centers)
}

// Reset drops every published center.
func (c *SlotCache) Reset() {
	c.centers = c.centers[:0]
}
