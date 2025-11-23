package server

import (
	"io"
	"log"
	"net/http"
	"net/url"
)

func (s *Server) media(w http.ResponseWriter, r *http.Request) {
	// Пробуем получить путь из PathValue (wildcard)
	file := r.PathValue("file")
	
	// Если PathValue не работает, извлекаем путь вручную из URL
	if file == "" {
		// Убираем префикс /api/media/ из пути
		path := r.URL.Path
		prefix := "/api/media/"
		if len(path) > len(prefix) && path[:len(prefix)] == prefix {
			file = path[len(prefix):]
		} else {
			log.Printf("Error: invalid path format: %s", path)
			Error(w, "invalid_path", "file path required", http.StatusBadRequest)
			return
		}
	}

	// Декодируем URL-encoded путь (на случай, если есть %2F и т.д.)
	decodedFile, err := url.PathUnescape(file)
	if err != nil {
		log.Printf("Warning: failed to decode file path %s: %v, using original", file, err)
		decodedFile = file
	}

	// Всегда используем дефолтный bucket "media"
	bucket := "media"

	// Логирование для отладки
	log.Printf("Media request: raw_path=%s, decoded_path=%s, bucket=%s", file, decodedFile, bucket)
	log.Printf("Full URL: %s, Path: %s", r.URL.String(), r.URL.Path)

	stream, err := s.svc.GetMedia(bucket, decodedFile)
	if err != nil {
		log.Printf("Media error: %v", err)
		Error(w, "not_found", err.Error(), http.StatusNotFound)
		return
	}

	if stream.ContentType != "" {
		w.Header().Set("Content-Type", stream.ContentType)
	}

	// Закрываем reader после копирования, если он реализует Close
	_, err = io.Copy(w, stream.Reader)
	if closer, ok := stream.Reader.(io.Closer); ok {
		closer.Close()
	}
	if err != nil {
		log.Printf("Error copying media stream: %v", err)
	}
}
