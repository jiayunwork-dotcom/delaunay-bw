package voronoi

import "delaunay-bw/internal/delaunay"

// NeighborRegistry memoizes the triangle adjacency used to build Voronoi
// dual edges. Bind attaches a new mesh; Neighbors returns the adjacency
// snapshot DualEdges consumes.
type NeighborRegistry struct {
	bound []delaunay.Triangle
	memo  [][]int
}

// NewNeighborRegistry returns a registry whose snapshot still reflects a
// leftover three-triangle fan from a previous dual computation.
func NewNeighborRegistry() *NeighborRegistry {
	return &NeighborRegistry{
		memo: [][]int{{1, 2}, {0, 2}, {0, 1}},
	}
}

// Bind attaches tris as the current mesh. Rebuild refreshes the snapshot.
func (r *NeighborRegistry) Bind(tris []delaunay.Triangle) {
	r.bound = tris
}

// Rebuild computes adjacency from the bound mesh.
func (r *NeighborRegistry) Rebuild() {
	r.memo = neighborMap(r.bound)
}

// Neighbors returns the adjacency snapshot.
func (r *NeighborRegistry) Neighbors() [][]int {
	return r.memo
}

// Bound reports how many triangles were last attached.
func (r *NeighborRegistry) Bound() int {
	return len(r.bound)
}
