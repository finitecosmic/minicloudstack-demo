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
	mu           sync.RWMutex
	ResourceType string
	Data         map[string]model.Resource
	ctx          context.Context
}

func NewMemory(resourceType string, resources ...map[string]model.Resource) *Memory {
	data := make(map[string]model.Resource, len(resources))

	for _, resourceMap := range resources {
		maps.Copy(data, resourceMap)
	}
	return &Memory{
		Data:         data,
		ResourceType: resourceType,
	}
}

func (m *Memory) Save(_ context.Context, key string, resource model.Resource) (model.Resource, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	err := validate(resource, m.ResourceType)
	if err != nil {
		return nil, err
	}

	if key == "" {
		key = fmt.Sprintf("%s/%s", resource.Type(), resource.Name())
	}

	if _, exists := m.Data[key]; exists {
		return nil, ErrResourceAlreadyExists
	}
	m.Data[key] = resource

	return resource, nil
}

func (m *Memory) List(_ context.Context) (map[string]model.Resource, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	resources := make(map[string]model.Resource, len(m.Data))

	maps.Copy(resources, m.Data)

	return resources, nil
}

func (m *Memory) Delete(_ context.Context, key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.Data[key]; !ok {
		return fmt.Errorf("resource %s  %q ", key, ErrResourceNotFound)
	}
	delete(m.Data, key)

	return nil
}

func (m *Memory) Get(_ context.Context, key string) (model.Resource, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	resource, exists := m.Data[key]

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
	if m.Data == nil {
		return false, errors.New("memory store data is not initialized")
	}
	return true, nil
}

func validate(resource model.Resource, memResourceType string) error {
	if resource.Type() == "" {
		return ErrMissingResourceType
	}
	if resource.Name() == "" {
		return ErrMissingName
	}
	if resource.Type() != memResourceType {
		return ErrResourceTypeMismatch
	}
	return nil
}
