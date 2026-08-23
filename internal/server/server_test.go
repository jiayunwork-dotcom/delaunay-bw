package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// newTestServer builds a Server backed by the real web/ and example/
// directories of the repository.
func newTestServer(t *testing.T) http.Handler {
	t.Helper()
	srv, err := New("../../web", "../../example")
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return srv.Handler()
}

func postJSON(h http.Handler, path, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodPost, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

// TestAPITriangulateValid drives the happy path: three points produce exactly
// one triangle that is declared Delaunay.
func TestAPITriangulateValid(t *testing.T) {
	h := newTestServer(t)
	rec := postJSON(h, "/api/triangulate", `{"points":[{"x":0,"y":0},{"x":2,"y":0},{"x":0,"y":2}]}`)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200 (body %s)", rec.Code, rec.Body.String())
	}
	var resp struct {
		Triangles    [][3]int `json:"triangles"`
		IsDelaunay   bool     `json:"delaunay"`
		NumTriangles int      `json:"num_triangles"`
		HullSize     int      `json:"hull_size"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("bad JSON: %v", err)
	}
	if resp.NumTriangles != 1 || len(resp.Triangles) != 1 {
		t.Errorf("triangles = %d, want 1", resp.NumTriangles)
	}
	if !resp.IsDelaunay {
		t.Error("response claims mesh is not Delaunay")
	}
	if resp.HullSize != 3 {
		t.Errorf("hull_size = %d, want 3", resp.HullSize)
	}
	if resp.Triangles[0] != [3]int{0, 1, 2} {
		t.Errorf("triangles[0] = %v, want [0 1 2]", resp.Triangles[0])
	}
}

// TestAPITriangulateInvalid pins the failure path: duplicate points are
// rejected with a 400 whose body carries the backend's error message.
func TestAPITriangulateInvalid(t *testing.T) {
	h := newTestServer(t)
	rec := postJSON(h, "/api/triangulate", `{"points":[{"x":0,"y":0},{"x":1,"y":0},{"x":0.5,"y":0.9},{"x":1e-12,"y":0}]}`)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400 (body %s)", rec.Code, rec.Body.String())
	}
	var body struct {
		Error string `json:"error"`
		Code  int    `json:"code"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("bad error JSON: %v", err)
	}
	if body.Code != http.StatusBadRequest {
		t.Errorf("error code = %d, want 400", body.Code)
	}
	if !strings.Contains(body.Error, "duplicate") {
		t.Errorf("error text %q does not mention duplicate points", body.Error)
	}
}

// TestAPITriangulateTooFew pins the minimum-count rule through the API.
func TestAPITriangulateTooFew(t *testing.T) {
	h := newTestServer(t)
	rec := postJSON(h, "/api/triangulate", `{"points":[{"x":0,"y":0},{"x":1,"y":1}]}`)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
	var body struct {
		Error string `json:"error"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("bad JSON: %v", err)
	}
	if !strings.Contains(body.Error, "three") {
		t.Errorf("error %q should mention the three-point minimum", body.Error)
	}
}

// TestAPIVoronoiValid checks that the voronoi endpoint returns one vertex per
// triangle and the expected dual edge for a square.
func TestAPIVoronoiValid(t *testing.T) {
	h := newTestServer(t)
	rec := postJSON(h, "/api/voronoi", `{"points":[{"x":0,"y":0},{"x":2,"y":0},{"x":2,"y":2},{"x":0,"y":2}]}`)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200 (body %s)", rec.Code, rec.Body.String())
	}
	var resp struct {
		Vertices []struct {
			X, Y float64
		} `json:"vertices"`
		Edges        [][2]int `json:"edges"`
		NumVertices  int      `json:"num_vertices"`
		NumEdges     int      `json:"num_edges"`
		NumTriangles int      `json:"num_triangles"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("bad JSON: %v", err)
	}
	if resp.NumTriangles != 2 {
		t.Errorf("num_triangles = %d, want 2", resp.NumTriangles)
	}
	if resp.NumVertices != 2 {
		t.Errorf("num_vertices = %d, want 2", resp.NumVertices)
	}
	if resp.NumEdges != 1 {
		t.Errorf("num_edges = %d, want 1", resp.NumEdges)
	}
}

// TestAPIRejectsWrongMethod pins that non-POST requests get a structured
// error instead of a silent success.
func TestAPIRejectsWrongMethod(t *testing.T) {
	h := newTestServer(t)
	req := httptest.NewRequest(http.MethodGet, "/api/triangulate", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want 405", rec.Code)
	}
	var body struct {
		Error string `json:"error"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("bad JSON: %v", err)
	}
	if body.Error == "" {
		t.Error("method rejection has empty error body")
	}
}

// TestAPIExampleServed checks the preloaded example is reachable over HTTP so
// the web console can load it before computing.
func TestAPIExampleServed(t *testing.T) {
	h := newTestServer(t)
	req := httptest.NewRequest(http.MethodGet, "/example/grid-jitter.json", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	var doc struct {
		Points []json.RawMessage `json:"points"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &doc); err != nil {
		t.Fatalf("example is not JSON: %v", err)
	}
	if len(doc.Points) != 20 {
		t.Errorf("example has %d points, want 20", len(doc.Points))
	}
}

// TestWebIndexServed checks the console page itself is served at /.
func TestWebIndexServed(t *testing.T) {
	h := newTestServer(t)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "Bowyer") {
		t.Error("index page does not reference the app")
	}
}
