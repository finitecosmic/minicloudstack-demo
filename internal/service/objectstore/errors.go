package objectstore

import "errors"

var (
	ErrBucketNotFound      = errors.New("bucket not found")
	ErrBucketAlreadyExists = errors.New("bucket already exists")
	ErrBucketNameRequired  = errors.New("bucket name is required")
	ErrBucketKeyRequired   = errors.New("bucket key is required")
	ErrInvalidKey          = errors.New("invalid key")
)
