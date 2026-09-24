package mock

import (
	"context"
	"minicloudstack/internal/model"
)

type FakeMemory struct {
	Data  map[string]model.Resource
	ctx   context.Context
	ready bool
	err   error
}

func NewFakeMemory() *FakeMemory {
	data := make(map[string]model.Resource)
	ctx := context.Background()

	return &FakeMemory{
		ctx:  ctx,
		Data: data,
	}
}

func (m *FakeMemory) Save(_ context.Context, key string, resource model.Resource) (model.Resource, error) {
	if m.err != nil {
		return nil, m.err
	}

	m.Data[key] = resource
	return m.Data[key], nil
}

func (m *FakeMemory) Get(_ context.Context, key string) (model.Resource, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.Data[key], nil
}

func (m *FakeMemory) List(_ context.Context) (map[string]model.Resource, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.Data, nil
}

func (m *FakeMemory) Delete(_ context.Context, key string) error {
	if m.err != nil {
		return m.err
	}
	return nil
}

func (m *FakeMemory) Ready(_ context.Context) (bool, error) {
	if m.err != nil {
		return false, m.err
	}
	return m.ready, nil
}
