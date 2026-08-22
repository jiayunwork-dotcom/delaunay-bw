package delaunay

import "fmt"

// Edge is an undirected edge between two vertex indices.
type Edge struct {
	A int
	B int
}

// NewEdge builds a normalized edge: A and B are swapped if necessary so that
// A <= B. Two edges that share endpoints in either order compare equal.
func NewEdge(a, b int) Edge {
	if a > b {
		a, b = b, a
	}
	return Edge{A: a, B: b}
}

// EdgeOf builds an edge from the two first vertices of a triangle. It is a
// small helper for code that iterates the three edges of a triangle.
func EdgeOf(t Triangle, k int) Edge {
	return NewEdge(t[k], t[(k+1)%3])
}

// Equals reports whether e and o connect the same two vertices.
func (e Edge) Equals(o Edge) bool {
	return e.A == o.A && e.B == o.B
}

// Reverse returns the edge with swapped endpoints. Because NewEdge
// normalizes, Reverse(e) equals e for undirected use.
func (e Edge) Reverse() Edge {
	return Edge{A: e.B, B: e.A}
}

// String formats the edge for diagnostics.
func (e Edge) String() string {
	return fmt.Sprintf("(%d-%d)", e.A, e.B)
}

// edgeKey is the comparable form of an Edge used as a map key.
type edgeKey struct{ a, b int }

// key converts the edge to its comparable map form.
func (e Edge) key() edgeKey {
	return edgeKey{a: e.A, b: e.B}
}

// EdgeCount counts how many distinct edges are referenced by the given
// triangles. It is used to validate that the mesh is a proper 2-manifold with
// each edge shared by at most two triangles.
func EdgeCount(tris []Triangle) int {
	seen := make(map[edgeKey]struct{}, 3*len(tris))
	for _, t := range tris {
		for k := 0; k < 3; k++ {
			seen[EdgeOf(t, k).key()] = struct{}{}
		}
	}
	return len(seen)
}
