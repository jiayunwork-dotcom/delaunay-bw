package delaunay

// Triangle is a mesh triangle given by three vertex indices into the point
// array. The vertex order is preserved from construction; the kernel keeps
// every triangle counter-clockwise.
type Triangle [3]int

// HasVertex reports whether v is one of the triangle's three vertices.
func (t Triangle) HasVertex(v int) bool {
	return t[0] == v || t[1] == v || t[2] == v
}

// HasAnyVertex reports whether the triangle touches any of the given indices.
// It is used to strip triangles that touch the (virtual) super-triangle
// vertices before returning the mesh.
func (t Triangle) HasAnyVertex(vs ...int) bool {
	for _, v := range vs {
		if t.HasVertex(v) {
			return true
		}
	}
	return false
}

// ContainsEdge reports whether the unordered edge between a and b is an edge
// of the triangle.
func (t Triangle) ContainsEdge(a, b int) bool {
	return (t[0] == a && t[1] == b) || (t[0] == b && t[1] == a) ||
		(t[1] == a && t[2] == b) || (t[1] == b && t[2] == a) ||
		(t[2] == a && t[0] == b) || (t[2] == b && t[0] == a)
}

// SharedEdge returns the two vertices shared by t and u, with ok=true, when
// the triangles are adjacent across an edge. It returns ok=false when the
// triangles do not share an edge.
func (t Triangle) SharedEdge(u Triangle) (a, b int, ok bool) {
	var shared []int
	for _, v := range t {
		if u.HasVertex(v) {
			shared = append(shared, v)
		}
	}
	if len(shared) != 2 {
		return 0, 0, false
	}
	return shared[0], shared[1], true
}

// OppositeVertex returns the vertex of t that is not one of the two edge
// endpoints a and b. ok is false when a and b are not both vertices of t.
func (t Triangle) OppositeVertex(a, b int) (int, bool) {
	for _, v := range t {
		if v != a && v != b {
			return v, true
		}
	}
	return 0, false
}

// Rotate returns the triangle with a cyclic shift of its vertices. The
// orientation (clockwise vs counter-clockwise) is preserved by rotation.
func (t Triangle) Rotate() Triangle {
	return Triangle{t[1], t[2], t[0]}
}

// Equals reports whether t and u reference the same three vertices regardless
// of their order.
func (t Triangle) Equals(u Triangle) bool {
	return t.HasVertex(u[0]) && t.HasVertex(u[1]) && t.HasVertex(u[2])
}

// key returns a canonical sorted key for the triangle's vertices, useful for
// set membership comparisons in tests and validators.
func (t Triangle) key() [3]int {
	v := [3]int{t[0], t[1], t[2]}
	// Insertion sort of three elements.
	if v[0] > v[1] {
		v[0], v[1] = v[1], v[0]
	}
	if v[1] > v[2] {
		v[1], v[2] = v[2], v[1]
	}
	if v[0] > v[1] {
		v[0], v[1] = v[1], v[0]
	}
	return v
}
