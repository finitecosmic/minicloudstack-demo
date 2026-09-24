package mock

import (
	"context"
	"minicloudstack/internal/model"
	"time"
)

type FakeNetwork struct {
	name         string
	cidr         string
	data         map[string]model.Resource
	resourceType string
	ctx          context.Context
	ready        bool
	err          error
}

func NewFakeNetwork(name string, cidr string) *FakeNetwork {
	data := make(map[string]model.Resource)

	return &FakeNetwork{
		data: data,
		cidr: cidr,
	}
}

func (m *FakeNetwork) Save(ctx context.Context, key string, resource model.Resource) (model.Resource, error) {
	if m.err != nil {
		return nil, m.err
	}

	m.data[key] = resource
	return m.data[key], nil
}

func (m *FakeNetwork) Get(ctx context.Context, key string) (model.Resource, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.data[key], nil
}

func (m *FakeNetwork) List(ctx context.Context) (map[string]model.Resource, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.data, nil
}

func (m *FakeNetwork) Delete(ctx context.Context, key string) error {
	if m.err != nil {
		return m.err
	}
	return nil
}

func (m *FakeNetwork) Ready(ctx context.Context) (bool, error) {
	if m.err != nil {
		return false, m.err
	}
	return m.ready, nil
}

func (m *FakeNetwork) Type() string {
	return m.resourceType
}

func (m *FakeNetwork) SetUpdatedAt(time.Time) {}
