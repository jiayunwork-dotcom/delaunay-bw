package delaunay

import (
	"encoding/json"
	"os"
	"testing"

	"delaunay-bw/internal/geom"
)

// loadExamplePoints reads a {"points":[...]} JSON example file and returns the
// point slice. Tests use it to drive the kernel with the same data the web
// console loads.
func loadExamplePoints(path string) ([]geom.Point, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var doc struct {
		Points []geom.Point `json:"points"`
	}
	if err := json.NewDecoder(f).Decode(&doc); err != nil {
		return nil, err
	}
	return doc.Points, nil
}

// triangleSet normalizes triangles into sorted vertex keys so two meshes can
// be compared as sets regardless of vertex order.
func triangleSet(tris []Triangle) map[[3]int]bool {
	out := make(map[[3]int]bool, len(tris))
	for _, t := range tris {
		out[t.key()] = true
	}
	return out
}

// TestExampleFilePresent guards the shipped example file against accidental
// deletion or reformatting to a different wire shape.
func TestExampleFilePresent(t *testing.T) {
	pts, err := loadExamplePoints("../../example/grid-jitter.json")
	if err != nil {
		t.Fatalf("example file missing or unreadable: %v", err)
	}
	if len(pts) < 3 {
		t.Fatalf("example has only %d points", len(pts))
	}
}
