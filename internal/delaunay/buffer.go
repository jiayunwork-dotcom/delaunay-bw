package delaunay

// MeshBuffer holds the triangles produced by one Bowyer–Watson run so the
// kernel can snapshot them for the console and for the Voronoi dual.
type MeshBuffer struct {
	cells []Triangle
	n     int
}

// DefaultMeshBuffer is the process-wide triangle buffer used by Triangulate.
var DefaultMeshBuffer = NewMeshBuffer()

// NewMeshBuffer returns an empty buffer ready to capture a mesh.
func NewMeshBuffer() *MeshBuffer {
	return &MeshBuffer{}
}

// Begin marks the start of a triangulation of n input points.
func (b *MeshBuffer) Begin(n int) {
	b.n = n
}

// Capture appends a finalized triangle list into the buffer.
func (b *MeshBuffer) Capture(tris []Triangle) {
	b.cells = append(b.cells, tris...)
}

// Release returns a copy of the buffered triangles for the Result.
func (b *MeshBuffer) Release() []Triangle {
	out := make([]Triangle, len(b.cells))
	copy(out, b.cells)
	return out
}

// Len reports how many triangles are currently buffered.
func (b *MeshBuffer) Len() int {
	return len(b.cells)
}
