package server

import (
	"encoding/json"
	"log"
	"net/http"
)

// errorBody is the uniform JSON shape of every error response:
//
//	{"error":"...", "code":400}
//
// The console renders the error field verbatim, so the message is written to
// be readable by an end user (it names the violated rule, not a stack trace).
type errorBody struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}

// writeError serializes err into an errorBody and sends it with the given
// status code. It never leaks internal stack traces; the message is the
// error's own text.
func writeError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	body := errorBody{Error: msg, Code: status}
	if err := json.NewEncoder(w).Encode(body); err != nil {
		// Encoding a two-field struct cannot fail in practice; log defensively
		// but do not interrupt the already-committed response.
		log.Printf("writeError: encode: %v", err)
	}
}

// NotAcceptable writes the error body used when the request wants a content
// type the server does not serve.
func notFound(w http.ResponseWriter) {
	writeError(w, http.StatusNotFound, "no such resource on this server")
}
