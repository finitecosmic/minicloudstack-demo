package handlers

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"mime"
	"net/http"
	"strings"
)

func parseContentType(r *http.Request) (string, error) {
	contentType := r.Header.Get("Content-Type")
	if contentType == "" {
		return "", fmt.Errorf("content type is empty")
	}
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return "", fmt.Errorf("invalid content type: %w", err)
	}
	return mediaType, nil
}

func writeResponse(w http.ResponseWriter, r *http.Request, response interface{}) {
	var err error

	accept := r.Header.Get("Accept")
	if strings.Contains(accept, "application/xml") {
		w.Header().Set("Content-Type", "application/xml")
		if err = xml.NewEncoder(w).Encode(response); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(response); err != nil {
		return
	}
	return
}
