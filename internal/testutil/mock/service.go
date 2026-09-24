package mock

import "minicloudstack/internal/state"

type FakeService struct {
	state        state.State
	resourceType string
}

func NewFakeService(stateStore state.State, resourceType string) *FakeService {
	return &FakeService{
		state:        stateStore,
		resourceType: resourceType,
	}
}
