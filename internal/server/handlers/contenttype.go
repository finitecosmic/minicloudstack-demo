package handlers

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"mime"
	"net/http"
	"strings"
)

var supportedMediaTypes = map[string]struct{}{
	"application/json": {},
	"application/xml":  {},
}

func isSupportedMediaType(contentType string) error {
	if contentType == "" {
		return ErrContentTypeEmpty
	}

	_, ok := supportedMediaTypes[contentType]
	if !ok {
		return ErrUnsupportedMedia
	}

	return nil
}

func parseContentType(r *http.Request) (string, error) {

	if err := isSupportedMediaType(r.Header.Get("Content-Type")); err != nil {
		return "", err
	}

	contentType := r.Header.Get("Content-Type")
	if contentType == "" {
		return "", fmt.Errorf("content type is empty")
	}

	mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil {
		return "", ErrParseMediaType
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
