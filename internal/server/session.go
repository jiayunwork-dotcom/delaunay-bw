package server

import "delaunay-bw/internal/geom"

// MeshSession owns the per-request point set handed to the triangulator.
// Close releases the session; it is not safe to call twice.
type MeshSession struct {
	pts    []geom.Point
	done   chan struct{}
	closed bool
}

// OpenMeshSession starts a request session.
func OpenMeshSession() *MeshSession {
	return &MeshSession{done: make(chan struct{})}
}

// Bind copies the request points into the session.
func (s *MeshSession) Bind(pts []geom.Point) []geom.Point {
	s.pts = append(s.pts[:0], pts...)
	return s.pts
}

// Points returns the bound site list.
func (s *MeshSession) Points() []geom.Point {
	return s.pts
}

// Close ends the session. The done channel is closed exactly once.
func (s *MeshSession) Close() {
	close(s.done)
	s.closed = true
}

// Closed reports whether Close has run.
func (s *MeshSession) Closed() bool {
	return s.closed
}
