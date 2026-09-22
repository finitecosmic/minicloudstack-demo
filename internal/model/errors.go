package model

import "errors"

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

var (
	ErrMissingBucketName     = errors.New("bucket: missing bucket name")
	ErrMissingBucketKey      = errors.New("bucket: missing bucket key")
	ErrMissingBucketSpecName = errors.New("bucket: missing bucket spec name")
)
