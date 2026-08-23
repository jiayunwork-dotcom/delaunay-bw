package voronoi

import (
	"math"
	"math/rand"
	"testing"

	"delaunay-bw/internal/delaunay"
	"delaunay-bw/internal/geom"
)

// jitteredMesh builds a triangulation of a jittered grid for tests.
func jitteredMesh(t *testing.T, cols, rows int, seed int64) ([]geom.Point, []delaunay.Triangle, *delaunay.Result) {
	t.Helper()
	r := rand.New(rand.NewSource(seed))
	pts := make([]geom.Point, 0, cols*rows)
	for i := 0; i < cols; i++ {
		for j := 0; j < rows; j++ {
			pts = append(pts, geom.Point{
				X: float64(i) + r.Float64()*0.3,
				Y: float64(j) + r.Float64()*0.3,
			})
		}
	}
	res, err := delaunay.Triangulate(pts, delaunay.Options{})
	if err != nil {
		t.Fatalf("Triangulate: %v", err)
	}
	return pts, res.Triangles, res
}

// TestVoronoiVerticesAreCircumcenters pins the defining property of the dual
// construction: every Voronoi vertex is exactly the circumcenter of its
// triangle, within the relative tolerance.
func TestVoronoiVerticesAreCircumcenters(t *testing.T) {
	pts, tris, _ := jitteredMesh(t, 5, 4, 3)
	diagram, err := Compute(pts, tris)
	if err != nil {
		t.Fatalf("Compute: %v", err)
	}
	if err := ValidateVertexCenters(pts, tris, diagram.Vertices, 1e-9); err != nil {
		t.Fatalf("vertex centers: %v", err)
	}
}

// TestVoronoiDualEdgeCount pins the dual edge count invariant: one dual edge
// per interior mesh edge, (3T - h)/2.
func TestVoronoiDualEdgeCount(t *testing.T) {
	pts, tris, res := jitteredMesh(t, 6, 5, 11)
	diagram, err := Compute(pts, tris)
	if err != nil {
		t.Fatalf("Compute: %v", err)
	}
	want := InteriorEdgeCount(len(tris), len(res.Hull))
	if diagram.EdgeCount() != want {
		t.Errorf("dual edges = %d, want %d (interior edges)", diagram.EdgeCount(), want)
	}
	if err := ValidateDualEdgeCount(len(tris), len(res.Hull), diagram.EdgeCount()); err != nil {
		t.Errorf("dual edge count check: %v", err)
	}
}

// TestVoronoiSquareDiagram pins a minimal exact case: a square yields two
// triangles whose shared diagonal is the only interior edge, so the Voronoi
// diagram has two vertices joined by exactly one dual edge.
func TestVoronoiSquareDiagram(t *testing.T) {
	square := []geom.Point{
		{X: 0, Y: 0},
		{X: 2, Y: 0},
		{X: 2, Y: 2},
		{X: 0, Y: 2},
	}
	res, err := delaunay.Triangulate(square, delaunay.Options{})
	if err != nil {
		t.Fatalf("Triangulate: %v", err)
	}
	diagram, err := Compute(square, res.Triangles)
	if err != nil {
		t.Fatalf("Compute: %v", err)
	}
	if diagram.VertexCount() != 2 {
		t.Errorf("Voronoi vertices = %d, want 2", diagram.VertexCount())
	}
	if diagram.EdgeCount() != 1 {
		t.Errorf("Voronoi edges = %d, want 1", diagram.EdgeCount())
	}
	// The single dual edge joins the two circumcenters.
	a, b := diagram.Edges[0][0], diagram.Edges[0][1]
	if math.Abs(diagram.Vertices[a].X-1) > 1e-12 || math.Abs(diagram.Vertices[a].Y-1) > 1e-12 {
		t.Errorf("vertex %d = %s, want (1,1)", a, diagram.Vertices[a].String())
	}
	if b != a && (math.Abs(diagram.Vertices[b].X-1) > 1e-12 || math.Abs(diagram.Vertices[b].Y-1) > 1e-12) {
		t.Errorf("vertex %d = %s, want (1,1)", b, diagram.Vertices[b].String())
	}
}

// TestVoronoiCentersDeviateOnlyByNoise checks the circumradius spread stays
// within floating-point noise for a general point set, confirming that each
// Voronoi vertex really is equidistant from the sites of its triangle.
func TestVoronoiCentersDeviateOnlyByNoise(t *testing.T) {
	pts, tris, _ := jitteredMesh(t, 7, 5, 42)
	dev, err := MaxCenterDeviation(pts, tris)
	if err != nil {
		t.Fatalf("MaxCenterDeviation: %v", err)
	}
	scale := geom.ScaleForTolerance(pts)
	if dev > 1e-6*scale*scale {
		t.Errorf("circumradius spread %g exceeds noise budget (scale²=%.3g)", dev, scale*scale)
	}
}

// TestVoronoiNeighborsMatchEdges checks that the neighbor map pairs each
// shared edge exactly once and never pairs a triangle with itself.
func TestVoronoiNeighborsMatchEdges(t *testing.T) {
	_, tris, _ := jitteredMesh(t, 5, 5, 7)
	neighbors := neighborMap(tris)
	if len(neighbors) != len(tris) {
		t.Fatalf("neighbor map length %d, want %d", len(neighbors), len(tris))
	}
	for i, ns := range neighbors {
		seen := map[int]bool{}
		for _, j := range ns {
			if j == i {
				t.Errorf("triangle %d listed as its own neighbor", i)
			}
			if seen[j] {
				t.Errorf("triangle %d has duplicate neighbor %d", i, j)
			}
			seen[j] = true
		}
	}
}
