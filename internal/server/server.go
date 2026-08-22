// Package server exposes the triangulation kernel over HTTP. It serves the
// static web console, the preloaded example files, and two JSON endpoints:
//
//	POST /api/triangulate   point set → triangles, Delaunay flag, statistics
//	POST /api/voronoi       point set → circumcenters and dual edges
//
// Every handler returns structured JSON, including error responses, so the
// console can render backend failures verbatim.
package server

import (
	"fmt"
	"net/http"
)

// Server owns the HTTP routing for the application.
type Server struct {
	handler http.Handler
}

// New builds a Server serving the web assets from webDir, the example files
// from exampleDir, and the JSON API. It validates the directories eagerly so
// a misconfiguration fails at startup instead of at first request.
func New(webDir, exampleDir string) (*Server, error) {
	mux := http.NewServeMux()

	if err := registerStatic(mux, webDir, exampleDir); err != nil {
		return nil, err
	}

	mux.HandleFunc("/api/triangulate", handleTriangulate)
	mux.HandleFunc("/api/voronoi", handleVoronoi)

	return &Server{handler: mux}, nil
}

// ListenAndServe starts the HTTP server on addr.
func (s *Server) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, s.handler)
}

// Handler returns the underlying http.Handler, useful for tests that drive
// the server with httptest.
func (s *Server) Handler() http.Handler {
	return s.handler
}

// AddressOf prints a human-readable serving address for startup messages.
func AddressOf(addr string) string {
	return fmt.Sprintf("http://%s", addr)
}
