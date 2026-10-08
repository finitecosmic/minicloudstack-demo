package testutil

import (
	"context"
	"fmt"
	"minicloudstack/internal/model"
	"minicloudstack/internal/state"
)

type MockState struct {
	ctx       context.Context
	resources *map[string]model.Resource
	Err       error
}

func NewMockState(ctx context.Context, resource *map[string]model.Resource, err error) *MockState {
	if resource == nil {
		resource = &map[string]model.Resource{}
	}

	return &MockState{
		ctx:       context.Background(),
		resources: resource,
		Err:       err,
	}
}

func (f *MockState) List(_ context.Context) (map[string]model.Resource, error) {
	if f.Err != nil {
		return nil, f.Err
	}
	return *f.resources, nil
}

func (f *MockState) Delete(_ context.Context, key string) error {
	if f.Err != nil {
		return f.Err
	}
	delete(*f.resources, key)
	return nil
}

func (f *MockState) Get(_ context.Context, key string) (model.Resource, error) {
	if f.Err != nil {
		return nil, f.Err
	}
	resource, exists := (*f.resources)[key]
	if !exists {
		return nil, fmt.Errorf("container %s %s", key, state.ErrResourceNotFound)
	}

	return resource, nil
}

func (f *MockState) Save(_ context.Context, key string, resource model.Resource) (model.Resource, error) {
	if f.Err != nil {
		return nil, f.Err
	}
	resource, exists := (*f.resources)[key]
	if !exists {
		return nil, fmt.Errorf("container %s %s", key, state.ErrResourceAlreadyExists)
	}
	(*f.resources)[key] = resource

	return resource, nil
}

func (f *MockState) Health(_ context.Context) (bool, error) {
	if f.Err != nil {
		return false, f.Err
	}
	return true, nil
}
