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
	version       string
	buildTime     string
	commitSHA     string
	domain        string
	canonicalName string
}

func newMetadata() *metadata {
	return &metadata{
		version:       Version,
		buildTime:     BuildTime,
		commitSHA:     CommitSHA,
		domain:        Domain,
		canonicalName: CanonicalName,
	}
}

// Handler is a wrapper for http.Handler
// that allows you to register liveness and readiness checkers.
type Handler interface {
	// Handler is a http.Handler that provides
	// /metadata
	http.Handler

	// LiveEndpoint is a HTTP handler for /live endpoint.
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

	// if not ?full=1 - return empty body. Kubernetes checks just HTTP code.
	// if r.URL.Query().Get("full") != "1" {
	// 	_, _ = w.Write([]byte("{}\n"))
	// 	return
	// }

	// otherwise write JSON body.
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "    ")
	_ = encoder.Encode(meta)
}
