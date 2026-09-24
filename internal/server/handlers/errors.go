package handlers

import "errors"

var (
	ErrBucketNameRequired = errors.New("bucket name is required")
	ErrBucketExists       = errors.New("bucket already exists")
	ErrBucketNotFound     = errors.New("bucket not found")
	ErrInvalidResource    = errors.New("invalid resource")
	ErrInvalidContentType = errors.New("invalid content type")
	ErrContentTypeEmpty   = errors.New("content type is empty")
	ErrUnsupportedMedia   = errors.New("unsupported media type")
	ErrInvalidJSON        = errors.New("invalid JSON")
	ErrParseMediaType     = errors.New("error parsing media type")
	ErrInvalidXML         = errors.New("invalid XML")
	ErrParseXML           = errors.New("error parsing XML")
	ErrInvalidUri         = errors.New("invalid uri")
	ErrDeleteMissingKey   = errors.New("missing key")
)
