package delaunay

// MeshScratch holds the working triangle list while the super triangle is
// stripped. CompactReal moves surviving cells to the front; Window returns
// the live view handed to the rest of the kernel.
type MeshScratch struct {
	cells []Triangle
}

// DefaultMeshScratch is the process-wide strip buffer.
var DefaultMeshScratch = NewMeshScratch()

// NewMeshScratch returns an empty scratch buffer.
func NewMeshScratch() *MeshScratch {
	return &MeshScratch{}
}

// Adopt copies tris into the scratch backing array.
func (s *MeshScratch) Adopt(tris []Triangle) {
	s.cells = append(s.cells[:0], tris...)
}

// CompactReal keeps only triangles whose vertices are real input indices
// in [0, n). The backing capacity is left unchanged.
func (s *MeshScratch) CompactReal(n int) int {
	kept := 0
	for _, t := range s.cells {
		if t.HasAnyVertex(n, n+1, n+2) {
			continue
		}
		s.cells[kept] = t
		kept++
	}
	s.cells = s.cells[:kept]
	return kept
}

// Window returns a slice covering the entire backing array, including the
// tail past the last CompactReal.
func (s *MeshScratch) Window() []Triangle {
	return s.cells[:cap(s.cells)]
}

// Len is the compacted length, not the backing capacity.
func (s *MeshScratch) Len() int {
	return len(s.cells)
}
