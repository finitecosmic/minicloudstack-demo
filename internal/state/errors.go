package state

import "errors"

var (
	ErrResourceNotFound      = errors.New("resource not found")
	ErrResourceAlreadyExists = errors.New("resource already exists")
	ErrMissingResourceType   = errors.New("missing resource type")
	ErrUnsupportedVersion    = errors.New("version not supported")
	ErrResourceTypeMismatch  = errors.New("resource type mismatch")
	ErrMissingDeleteKey      = errors.New("missing delete key")
	ErrMissingName           = errors.New("missing name")
	ErrMissingResource       = errors.New("missing resource")
	ErrUnmarshalResource     = errors.New("unmarshal resource")
	ErrBeginTransaction      = errors.New("begin transaction")
	ErrQueryRowContext       = errors.New("query row context")
	ErrMigrateV1toV2         = errors.New("migrate schema v1 to v2")
	ErrSetSchemaVersion      = errors.New("set schema version")
	ErrTxCommit              = errors.New("tx commit")
	ErrColumnDoesNotExist    = errors.New("column does not exist")
)
