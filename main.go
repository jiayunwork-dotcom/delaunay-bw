// Command delaunay-bw computes the Bowyer-Watson Delaunay triangulation of a
// planar point set and serves an interactive web console plus JSON endpoints.
//
// Usage:
//
//	go run . -http :8080
//
// The web console at / loads example/grid-jitter.json, posts the point set to
// /api/triangulate and /api/voronoi, and draws the resulting mesh as SVG.
// Every number on the page comes from the backend kernel.
package main

import (
	"flag"
	"fmt"
	"os"

	"delaunay-bw/internal/server"
)

func main() {
	httpAddr := flag.String("http", ":8080", "HTTP listen address for the web console and /api")
	webDir := flag.String("web-dir", "web", "directory of the static web assets")
	exampleDir := flag.String("example-dir", "example", "directory of example JSON files")
	flag.Parse()

	srv, err := server.New(*webDir, *exampleDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "delaunay-bw: %v\n", err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stderr, "delaunay-bw: serving %s on %s (web=%q example=%q)\n",
		"web console + /api/triangulate + /api/voronoi", *httpAddr, *webDir, *exampleDir)
	if err := srv.ListenAndServe(*httpAddr); err != nil {
		fmt.Fprintf(os.Stderr, "delaunay-bw: server error: %v\n", err)
		os.Exit(1)
	}
}
