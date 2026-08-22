package server

import (
	"encoding/json"
	"net/http"

	"delaunay-bw/internal/geom"
)

// TriangulateResponse is the JSON answer of POST /api/triangulate. Every
// field is computed by the kernel; no value is synthesized on the server
// front end.
type TriangulateResponse struct {
	Triangles         [][3]int `json:"triangles"`
	IsDelaunay        bool     `json:"delaunay"`
	NumPoints         int      `json:"num_points"`
	NumTriangles      int      `json:"num_triangles"`
	HullSize          int      `json:"hull_size"`
	HullArea          float64  `json:"hull_area"`
	TriangleArea      float64  `json:"triangle_area"`
	AreaDelta         float64  `json:"area_delta"`
	ExpectedTriangles int      `json:"expected_triangles"`
	MatchesFormula    bool     `json:"matches_formula"`
}

// VoronoiResponse is the JSON answer of POST /api/voronoi. Vertices holds one
// circumcenter per triangle and Edges pairs vertex indices; together they
// form the finite part of the Voronoi diagram.
type VoronoiResponse struct {
	Vertices     []geom.Point `json:"vertices"`
	Edges        [][2]int     `json:"edges"`
	NumVertices  int          `json:"num_vertices"`
	NumEdges     int          `json:"num_edges"`
	NumTriangles int          `json:"num_triangles"`
}

// writeJSON writes v as a JSON document with the given status code.
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
