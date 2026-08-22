package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"delaunay-bw/internal/geom"
)

// TriangulateRequest is the JSON body accepted by both API endpoints. The
// points may be given either directly as an array or wrapped in a {"points":
// [...]} object. tolerance is optional and, when omitted, the kernel default
// applies.
type TriangulateRequest struct {
	Points    []geom.Point `json:"points"`
	Tolerance float64      `json:"tolerance,omitempty"`
}

// decodeRequest parses and validates the JSON body of an API call.
//
// It accepts both wire shapes mentioned in the README, caps the point count
// to keep the O(n²) duplicate scan and the incircle loop bounded, and checks
// finiteness early so the kernel only ever sees well-formed input. All errors
// returned are meant to be embedded verbatim in the error JSON body.
func decodeRequest(r *http.Request) (TriangulateRequest, error) {
	var req TriangulateRequest

	body, err := io.ReadAll(io.LimitReader(r.Body, 1<<20))
	if err != nil {
		return req, fmt.Errorf("request body unreadable: %v", err)
	}
	if len(body) == 0 {
		return req, fmt.Errorf("empty request body: expected a JSON point set")
	}

	// Accept both {"points":[...]} and the bare array.
	var wrapped struct {
		Points    []geom.Point `json:"points"`
		Tolerance float64      `json:"tolerance,omitempty"`
	}
	if err := json.Unmarshal(body, &wrapped); err != nil {
		// Fall back to a bare array.
		var arr []geom.Point
		if arrErr := json.Unmarshal(body, &arr); arrErr != nil {
			return req, fmt.Errorf("body is not a JSON point array or {\"points\":[...]}: %v", err)
		}
		req.Points = arr
	} else {
		req.Points = wrapped.Points
		req.Tolerance = wrapped.Tolerance
	}

	if len(req.Points) > maxPoints {
		return req, fmt.Errorf("too many points: at most %d supported, got %d", maxPoints, len(req.Points))
	}
	if len(req.Points) == 0 {
		return req, geom.ErrTooFewPoints
	}
	return req, nil
}

// maxPoints caps the accepted point count. The kernel is O(n²) in the
// duplicate scan and O(n·t) in the empty-circle validation; the cap keeps a
// public endpoint from being a trivially cheap CPU amplification vector.
const maxPoints = 2000
