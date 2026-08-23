// Package voronoi builds the Voronoi diagram dual to a Delaunay mesh.
//
// For every triangle of the mesh its circumcenter is a Voronoi vertex, and
// every pair of triangles that share an edge yields a Voronoi edge joining
// their circumcenters. Because the mesh is the Delaunay triangulation, these
// two facts are equivalent to the standard definitions: each Voronoi cell is
// the set of points closer to one input site than to any other, and the dual
// edges cross the Delaunay edges at right angles.
package voronoi

import (
	"fmt"

	"delaunay-bw/internal/delaunay"
	"delaunay-bw/internal/geom"
)

// Diagram is the finite part of the Voronoi diagram. Vertices holds one
// circumcenter per triangle (in triangle order), and Edges lists the pairs
// of vertex indices that are joined by a dual edge.
type Diagram struct {
	Vertices  []geom.Point  `json:"vertices"`
	Edges     [][2]int      `json:"edges"`
	Triangles []delaunay.Triangle `json:"-"`
}

// Compute derives the Voronoi diagram from a triangulated point set. The
// caller must pass the triangles exactly as returned by the triangulator
// (with super-triangle cells already stripped). An error is returned when any
// triangle is degenerate, because its circumcenter does not exist.
func Compute(pts []geom.Point, tris []delaunay.Triangle) (*Diagram, error) {
	centers, err := Circumcenters(pts, tris)
	if err != nil {
		return nil, err
	}
	reg := NewNeighborRegistry()
	reg.Bind(tris)
	dual, err := DualEdges(len(tris), reg.Neighbors())
	if err != nil {
		return nil, err
	}
	return &Diagram{
		Vertices:  centers,
		Edges:     dual,
		Triangles: tris,
	}, nil
}

// VertexCount is the number of Voronoi vertices, one per triangle.
func (d *Diagram) VertexCount() int {
	return len(d.Vertices)
}

// EdgeCount is the number of dual edges. Interior mesh edges each produce one
// dual edge; boundary edges produce none.
func (d *Diagram) EdgeCount() int {
	return len(d.Edges)
}

// String renders a compact summary of the diagram for diagnostics.
func (d *Diagram) String() string {
	return fmt.Sprintf("voronoi diagram: %d vertices, %d edges", len(d.Vertices), len(d.Edges))
}
