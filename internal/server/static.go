package server

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

// staticRoutes wires up the file serving for the web console and the example
// data:
//
//	/                 → web/index.html
//	/web/             → other static assets (app.js, style.css)
//	/example/         → example JSON files (grid-jitter.json)
//
// The example endpoint lets the console fetch a point set and then hand it
// straight to the kernel, keeping the "load example then compute" workflow
// entirely backend-driven.
func registerStatic(mux *http.ServeMux, webDir, exampleDir string) error {
	webAbs, err := filepath.Abs(webDir)
	if err != nil {
		return fmt.Errorf("web dir: %w", err)
	}
	if st, err := os.Stat(webAbs); err != nil || !st.IsDir() {
		return fmt.Errorf("web dir %q is not a directory: %v", webDir, err)
	}

	exAbs, err := filepath.Abs(exampleDir)
	if err != nil {
		return fmt.Errorf("example dir: %w", err)
	}
	if st, err := os.Stat(exAbs); err != nil || !st.IsDir() {
		return fmt.Errorf("example dir %q is not a directory: %v", exampleDir, err)
	}

	webFS := http.FileServer(http.Dir(webAbs))
	exampleFS := http.FileServer(http.Dir(exAbs))

	mux.Handle("/", webFS)
	mux.Handle("/web/", http.StripPrefix("/web/", webFS))
	mux.Handle("/example/", http.StripPrefix("/example/", exampleFS))
	return nil
}
