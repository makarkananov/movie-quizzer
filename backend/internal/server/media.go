package server

import (
	"io"
	"net/http"
)

func (s *Server) media(w http.ResponseWriter, r *http.Request) {
	bucket := r.PathValue("bucket")
	file := r.PathValue("file")

	if bucket == "" || file == "" {
		Error(w, "invalid_path", "bucket and file required", http.StatusBadRequest)
		return
	}

	stream, err := s.svc.GetMedia(bucket, file)
	if err != nil {
		Error(w, "not_found", err.Error(), http.StatusNotFound)
		return
	}

	if stream.ContentType != "" {
		w.Header().Set("Content-Type", stream.ContentType)
	}

	_, _ = io.Copy(w, stream.Reader)
}
