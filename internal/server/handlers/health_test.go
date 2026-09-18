package handlers

import (
	"context"
	"errors"
	"minicloudstack/internal/testutil"
	"testing"
)

func TestHealthHandler_StoreReadiness(t *testing.T) {
	var tests = []struct {
		name      string
		stateErr  error
		wantReady bool
	}{
		{
			name:      "ok state",
			stateErr:  nil,
			wantReady: true,
		},
		{
			name:      "not ok state",
			stateErr:  errors.New("memory state is nil"),
			wantReady: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			stateStore := testutil.NewFakeState(ctx, tt.stateErr)

			handler := NewHealthHandler(stateStore)
			got, _ := handler.state.Ready(ctx)
			if got != tt.wantReady {
				t.Errorf("HealthHandler.Ready() got = %v, want %v", got, tt.wantReady)
			}
		})
	}
}
