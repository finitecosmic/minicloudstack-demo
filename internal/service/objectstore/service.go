package objectstore

import (
	"minicloudstack/internal/state"
)

type Service struct {
	state state.State
}

func New(stateStore state.State) *Service {
	return &Service{
		state: stateStore,
	}
}
