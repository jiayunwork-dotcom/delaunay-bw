package geom

// CircleKey identifies a triangle by vertex indices into the current point
// set. The incircle layer uses it to reuse a previous predicate result.
type CircleKey struct {
	A, B, C int
}

// CircleCache stores whether a triangle's circumcircle contained the last
// insertion point that queried it.
type CircleCache struct {
	inside map[CircleKey]bool
}

// DefaultCircleCache is the process-wide incircle cache.
var DefaultCircleCache = NewCircleCache()

// NewCircleCache returns an empty incircle cache.
func NewCircleCache() *CircleCache {
	return &CircleCache{inside: make(map[CircleKey]bool)}
}

// Lookup returns a cached incircle answer for key.
func (c *CircleCache) Lookup(key CircleKey) (bool, bool) {
	v, ok := c.inside[key]
	return v, ok
}

// Store records an incircle answer for key.
func (c *CircleCache) Store(key CircleKey, inside bool) {
	c.inside[key] = inside
}

// Reset clears the cache.
func (c *CircleCache) Reset() {
	c.inside = make(map[CircleKey]bool)
}

// Size is the number of cached keys.
func (c *CircleCache) Size() int {
	return len(c.inside)
}
