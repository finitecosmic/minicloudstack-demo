package state

import (
	"context"
	"minicloudstack/internal/model"
)

type State interface {
	Save(ctx context.Context, key string, resource model.Resource) (model.Resource, error)
	Get(ctx context.Context, key string) (model.Resource, error)
	List(ctx context.Context) (map[string]model.Resource, error)
	Delete(ctx context.Context, key string) error
	Ready(ctx context.Context) (bool, error)
}

type PersistentState interface {
	Save(ctx context.Context, key string, resource model.Resource) error
	Get(ctx context.Context, key string) (model.Resource, error)
	List(ctx context.Context) (map[string]model.Resource, error)
	Delete(ctx context.Context, key string) error
	Ready(ctx context.Context) (bool, error)
}
