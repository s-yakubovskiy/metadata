package metadata

import (
	"encoding/json"
	"net/http"
)

const (
	MetadataHandlerPath = "/metadata"
)

// Variables defined by the Makefile and passed in with ldflags
var (
	Version       = "latest"
	BuildTime     = "build-time"
	CommitSHA     = "commit-sha"
	Domain        = "domain"
	CanonicalName = "canonical-name"
)

type metadata struct {
	Version       string
	BuildTime     string
	CommitSHA     string
	Domain        string
	CanonicalName string
}

func newMetadata() *metadata {
	return &metadata{
		Version:       Version,
		BuildTime:     BuildTime,
		CommitSHA:     CommitSHA,
		Domain:        Domain,
		CanonicalName: CanonicalName,
	}
}

// Handler is a wrapper for http.Handler
// that allows you to register liveness and readiness checkers.
type Handler interface {
	// Handler is a http.Handler that provides
	// /metadata
	http.Handler

	// MetadataEndpoint is a HTTP handler for /metadata endpoint.
	MetadataEndpoint(http.ResponseWriter, *http.Request)
}

// NewHandler creates new base Handler
func NewHandler() Handler {
	h := NewBasicHandler()
	h.Handle(MetadataHandlerPath, http.HandlerFunc(h.MetadataEndpoint))
	return h
}

// basicHandler is a basic Handler implementation.
type basicHandler struct {
	http.ServeMux
}

func NewBasicHandler() *basicHandler {
	return &basicHandler{}
}

func (s *basicHandler) MetadataEndpoint(w http.ResponseWriter, r *http.Request) {
	s.handle(w, r)
}

func (s *basicHandler) handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	status := http.StatusOK

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	w.WriteHeader(status)

	meta := newMetadata()

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "    ")
	_ = encoder.Encode(meta)
}
