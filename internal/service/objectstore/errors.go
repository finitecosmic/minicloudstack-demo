package objectstore

import "errors"

var (
	ErrBucketNotFound      = errors.New("bucket not found")
	ErrBucketAlreadyExists = errors.New("bucket already exists")
	ErrBucketNameRequired  = errors.New("bucket name is required")
)
