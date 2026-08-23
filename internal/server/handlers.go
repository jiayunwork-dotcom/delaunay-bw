package server

import (
	"net/http"

	"delaunay-bw/internal/delaunay"
	"delaunay-bw/internal/voronoi"
)

// handleTriangulate serves POST /api/triangulate. It decodes the point set,
// runs the Bowyer–Watson kernel and answers with the triangle indices, the
// Delaunay flag and the statistics table. Any validation failure is mapped to
// a 400 with an error body that names the violated rule.
func handleTriangulate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method must be POST")
		return
	}
	req, err := decodeRequest(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	opts := delaunay.Options{Tolerance: req.Tolerance}
	result, err := delaunay.Triangulate(req.Points, opts)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	// The response carries a per-triangle Delaunay verification so the caller
	// can audit the empty-circle property without recomputing it.
	resp := TriangulateResponse{
		Triangles:   trianglesAsJSON(result.Triangles),
		IsDelaunay:  delaunay.IsDelaunay(req.Points, result.Triangles, opts.EffectiveTolerance()),
		NumPoints:   result.Stats.NumPoints,
		NumTriangles: result.Stats.NumTriangles,
		HullSize:    result.Stats.HullSize,
		HullArea:    result.Stats.HullArea,
		TriangleArea: result.Stats.TriangleArea,
		AreaDelta:   result.Stats.AreaDelta,
		ExpectedTriangles: result.Stats.ExpectFormula,
		MatchesFormula:    result.Stats.MatchesFormula,
	}
	writeJSON(w, http.StatusOK, resp)
}

// handleVoronoi serves POST /api/voronoi. It triangulates the point set and
// returns the Voronoi vertex list (the circumcenters) together with the dual
// edge pairs. The dual edges reference indices into the vertex list.
func handleVoronoi(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method must be POST")
		return
	}
	req, err := decodeRequest(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	opts := delaunay.Options{Tolerance: req.Tolerance}
	result, err := delaunay.Triangulate(req.Points, opts)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	diagram, err := voronoi.Compute(req.Points, result.Triangles)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	published := DefaultVoronoiSession.Publish(diagram)
	resp := VoronoiResponse{
		Vertices:     published.Vertices,
		Edges:        published.Edges,
		NumVertices:  published.VertexCount(),
		NumEdges:     published.EdgeCount(),
		NumTriangles: len(result.Triangles),
	}
	writeJSON(w, http.StatusOK, resp)
}

// trianglesAsJSON converts kernel triangles to [][3]int for the wire format.
func trianglesAsJSON(tris []delaunay.Triangle) [][3]int {
	out := make([][3]int, 0, len(tris))
	for _, t := range tris {
		out = append(out, [3]int{t[0], t[1], t[2]})
	}
	return out
}
