package state

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"minicloudstack/internal/model"
	"sync"
)

type Memory struct {
	mu   sync.RWMutex
	data map[string]model.Resource
	ctx  context.Context
}

func NewMemory(data ...map[string]model.Resource) *Memory {
	var resources map[string]model.Resource
	if len(data) > 0 {
		resources = data[0]
	}
	if resources == nil {
		resources = make(map[string]model.Resource)
	}

	return &Memory{
		mu:   sync.RWMutex{},
		data: resources,
	}
}

func (m *Memory) Save(_ context.Context, key string, resource model.Resource) (model.Resource, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if key == "" {
		key = resource.Key()
	}
	if _, exists := m.data[key]; exists {
		return nil, fmt.Errorf(
			"resource %q: %w",
			key,
			ErrResourceAlreadyExists,
		)
	}

	m.data[key] = resource
	return m.data[key], nil
}

func (m *Memory) List(_ context.Context) (map[string]model.Resource, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	resources := make(map[string]model.Resource, len(m.data))

	maps.Copy(resources, m.data)

	return resources, nil
}

func (m *Memory) Delete(_ context.Context, key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.data[key]; !ok {
		return fmt.Errorf("resource %s  %q ", key, ErrResourceNotFound)
	}
	delete(m.data, key)

	return nil
}

func (m *Memory) Get(_ context.Context, key string) (model.Resource, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	resource, exists := m.data[key]

	if !exists {
		return resource, ErrResourceNotFound
	}

	return resource, nil
}

func (m *Memory) Ready(_ context.Context) (bool, error) {
	if m == nil {
		return false, errors.New("memory state is nil")
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.data == nil {
		return false, errors.New("memory store data is not initialized")
	}
	return true, nil
}
