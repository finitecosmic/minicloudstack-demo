package state

import "errors"

var (
	ErrResourceNotFound      = errors.New("resource not found")
	ErrResourceAlreadyExists = errors.New("resource already exists")
	ErrMissingResourceType   = errors.New("missing resource type")
	ErrResourceTypeMismatch  = errors.New("resource type mismatch")
	ErrMissingDeleteKey      = errors.New("missing delete key")
	ErrMissingName           = errors.New("missing name")
	ErrMissingResource       = errors.New("missing resource")
	ErrUnmarshalResource     = errors.New("unmarshal resource")
)
