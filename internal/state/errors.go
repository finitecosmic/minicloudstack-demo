package state

import "errors"

var (
	ErrResourceNotFound        = errors.New("resource not found")
	ErrResourceAlreadyExists   = errors.New("resource already exists")
	ErrResourceDoesNotExist    = errors.New("resource	 not found")
	ErrInvalidKey              = errors.New("invalid resource key")
	ErrInvalidResourceType     = errors.New("invalid resource type")
	ErrMissingResourceType     = errors.New("missing resource type")
	ErrResourceTypeMismatch    = errors.New("resource type mismatch")
	ErrUnexpectedResourceType  = errors.New("unexpected resource type")
	ErrUnsupportedResourceType = errors.New("unsupported resource type")
	ErrMissingDeleteKey        = errors.New("missing delete key")
	ErrMissingName             = errors.New("missing name")
	ErrMissingResource         = errors.New("missing resource")
	ErrSQLiteTableDoesNotExist = errors.New("sqlite table does not exist")
)
