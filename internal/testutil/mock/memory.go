package mock

import (
	"context"
	"minicloudstack/internal/model"
)

type FakeMemory struct {
	data  map[string]model.Resource
	ctx   context.Context
	ready bool
	err   error
}

func NewFakeMemory(resources ...map[string]model.Resource) *FakeMemory {
	data := make(map[string]model.Resource)
	for k, resource := range resources {
		print(k, resource)
	}
	return &FakeMemory{
		data: data,
	}
}

func (m *FakeMemory) Save(ctx context.Context, key string, resource model.Resource) (model.Resource, error) {
	if m.err != nil {
		return nil, m.err
	}

	m.data[key] = resource
	return m.data[key], nil
}

func (m *FakeMemory) Get(ctx context.Context, key string) (model.Resource, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.data[key], nil
}

func (m *FakeMemory) List(ctx context.Context) (map[string]model.Resource, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.data, nil
}

func (m *FakeMemory) Delete(ctx context.Context, key string) error {
	if m.err != nil {
		return m.err
	}
	return nil
}

func (m *FakeMemory) Ready(ctx context.Context) (bool, error) {
	if m.err != nil {
		return false, m.err
	}
	return m.ready, nil
}
