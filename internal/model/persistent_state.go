package model

import (
	"context"
	"time"
)

type PersistentState struct {
	ID          string
	Kind        string
	Name        string
	Provider    string
	Spec        []byte
	State       []byte
	Status      string
	CreatedTime time.Time
	UpdatedTime time.Time
}

type PersistentStateStore interface {
	Get(ctx context.Context, key string) (PersistentState, error)
	Save(ctx context.Context, key string, state PersistentState) error
	Delete(ctx context.Context, key string) error
	List(ctx context.Context) ([]PersistentState, error)
}

type Config struct {
	DatabasePath string
}
