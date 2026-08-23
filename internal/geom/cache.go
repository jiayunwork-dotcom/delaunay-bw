package geom

// SlotKey identifies a triangle by its vertices and their coordinates.
// Slot indices in the working mesh change when the super triangle is
// stripped; coordinates distinguish a later point set that reuses 0..n.
type SlotKey struct {
	A, B, C int
	AX, AY  float64
	BX, BY  float64
	CX, CY  float64
}

// SlotCache stores circumcenters keyed by SlotKey so insertion and the
// Voronoi dual share the same identity after topology changes.
type SlotCache struct {
	centers map[SlotKey]Point
}

// DefaultSlotCache is the shared circumcenter table.
var DefaultSlotCache = NewSlotCache()

// NewSlotCache returns an empty slot table.
func NewSlotCache() *SlotCache {
	return &SlotCache{centers: make(map[SlotKey]Point)}
}

// MakeSlotKey builds the stable identity for triangle (ia,ib,ic).
func MakeSlotKey(a, b, c Point, ia, ib, ic int) SlotKey {
	return SlotKey{
		A: ia, B: ib, C: ic,
		AX: a.X, AY: a.Y,
		BX: b.X, BY: b.Y,
		CX: c.X, CY: c.Y,
	}
}

// Put records the circumcenter of the triangle identified by key.
func (c *SlotCache) Put(key SlotKey, p Point) {
	if c.centers == nil {
		c.centers = make(map[SlotKey]Point)
	}
	c.centers[key] = p
}

// Get returns the circumcenter last published for key.
func (c *SlotCache) Get(key SlotKey) (Point, bool) {
	if c.centers == nil {
		return Point{}, false
	}
	p, ok := c.centers[key]
	return p, ok
}

// Len is the number of published centers.
func (c *SlotCache) Len() int {
	return len(c.centers)
}

// Reset drops every published center.
func (c *SlotCache) Reset() {
	c.centers = make(map[SlotKey]Point)
}
