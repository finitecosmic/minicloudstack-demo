package state

import "errors"

var (
	ErrResourceNotFound      = errors.New("resource not found")
	ErrResourceAlreadyExists = errors.New("resource already exists")
	ErrResourceDoesNotExist  = errors.New("resource	 not found")
	ErrInvalidKey            = errors.New("invalid resource key")
)
