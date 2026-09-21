package testutil

import (
	"context"
	"fmt"
	"minicloudstack/internal/model"
	"minicloudstack/internal/state"
)

type FakeState struct {
	resource map[string]model.Resource
	err      error
	ctx      context.Context
}

func NewFakeState(ctx context.Context, err error) *FakeState {
	return &FakeState{
		resource: make(map[string]model.Resource),
		err:      err,
		ctx:      context.Background(),
	}
}

func (f *FakeState) List(_ context.Context) (map[string]model.Resource, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.resource, nil
}

func (f *FakeState) Delete(_ context.Context, key string) error {
	if f.err != nil {
		return f.err
	}
	delete(f.resource, key)
	return nil
}

func (f *FakeState) Get(_ context.Context, key string) (model.Resource, error) {
	if f.err != nil {
		return nil, f.err
	}
	resource, exists := f.resource[key]
	if !exists {
		return nil, fmt.Errorf("container %s %s", key, state.ErrResourceNotFound)
	}

	return resource, nil
}

func (f *FakeState) Save(_ context.Context, key string, resource model.Resource) (model.Resource, error) {
	if f.err != nil {
		return nil, f.err
	}
	resource, exists := f.resource[key]
	if !exists {
		return nil, fmt.Errorf("container %s %s", key, state.ErrResourceAlreadyExists)
	}
	f.resource[key] = resource

	return resource, nil
}

func (f *FakeState) Ready(_ context.Context) (bool, error) {
	if f.err != nil {
		return false, f.err
	}
	return true, nil
}
