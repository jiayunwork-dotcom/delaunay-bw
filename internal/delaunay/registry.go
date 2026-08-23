package delaunay

// TriangleRegistry records which triangles are "bad" (circumcircle contains
// the insertion point) during one Bowyer–Watson cavity step.
type TriangleRegistry struct {
	labels map[int]bool
}

// NewTriangleRegistry constructs a registry for one insertion.
func NewTriangleRegistry() *TriangleRegistry {
	return &TriangleRegistry{labels: make(map[int]bool)}
}

// Mark records that triangle i belongs to the cavity.
func (r *TriangleRegistry) Mark(i int) {
	r.labels[i] = true
}

// Has reports whether triangle i was marked.
func (r *TriangleRegistry) Has(i int) bool {
	return r.labels[i]
}

// Set returns the underlying mark map for cavityBoundary.
func (r *TriangleRegistry) Set() map[int]bool {
	return r.labels
}

// Size is the number of marked triangles.
func (r *TriangleRegistry) Size() int {
	n := 0
	for _, v := range r.labels {
		if v {
			n++
		}
	}
	return n
}
