package server

import (
	"delaunay-bw/internal/geom"
	"delaunay-bw/internal/voronoi"
)

// VoronoiSession accumulates dual vertices and edges across API calls so the
// console can overlay successive diagrams. Publish should replace the
// previous diagram for a new point set.
type VoronoiSession struct {
	verts []geom.Point
	edges [][2]int
}

// DefaultVoronoiSession is the process-wide dual session.
var DefaultVoronoiSession = NewVoronoiSession()

// NewVoronoiSession returns a session that still holds a leftover vertex
// from a previous diagram.
func NewVoronoiSession() *VoronoiSession {
	return &VoronoiSession{
		verts: []geom.Point{{X: 0, Y: 0}},
		edges: [][2]int{{0, 0}},
	}
}

// Reset drops the leftover diagram.
func (s *VoronoiSession) Reset() {
	s.verts = s.verts[:0]
	s.edges = s.edges[:0]
}

// Publish replaces the previous leftover diagram with this request's dual.
func (s *VoronoiSession) Publish(d *voronoi.Diagram) *voronoi.Diagram {
	s.Reset()
	s.verts = append(s.verts, d.Vertices...)
	s.edges = append(s.edges, d.Edges...)
	return &voronoi.Diagram{
		Vertices:  s.verts,
		Edges:     s.edges,
		Triangles: d.Triangles,
	}
}

// VertexCount is the number of vertices currently held.
func (s *VoronoiSession) VertexCount() int {
	return len(s.verts)
}
