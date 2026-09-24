// internal/state/memory_test.go
package state_test

import (
	"context"
	"errors"
	"minicloudstack/internal/model"
	"minicloudstack/internal/state"
	"minicloudstack/internal/testutil/mock"
	"strings"
	"testing"
)

func TestMemory_Save(t *testing.T) {
	tests := []struct {
		name             string
		newResource      model.Resource
		newResourceName  string
		err              error
		setup            func(m *state.Memory)
		resourceType     string
		wantName         string
		wantKey          string
		wantNumResources int
		wantErr          error
	}{
		{
			name:         "valid save bucket",
			resourceType: model.ResourceTypeBucket,
			newResource:  mock.NewFakeBucket("new-bucket", "bucket/new-bucket", nil),
			setup: func(m *state.Memory) {
				m.Data["bucket/existing-bucket"] = mock.NewFakeBucket(
					"existing-bucket",
					"bucket/new-bucket",
					nil)
			},

			wantName:         "new-bucket",
			wantKey:          "bucket/new-bucket",
			wantNumResources: 2,
			wantErr:          nil,
		},
		{
			name:         "valid save bucket",
			resourceType: model.ResourceTypeBucket,
			newResource: mock.NewFakeBucket(
				"new-bucket", "bucket/new-bucket", nil,
			),
			newResourceName:  "new-bucket",
			setup:            nil,
			wantName:         "new-bucket",
			wantKey:          "bucket/new-bucket",
			wantNumResources: 1,
			wantErr:          nil,
		},
		{
			name:         "save already existing bucket error",
			resourceType: model.ResourceTypeBucket,
			newResource: mock.NewFakeBucket(
				"existing-bucket", "bucket/existing-bucket", nil,
			),
			newResourceName: "bucket/existing-bucket",
			setup: func(m *state.Memory) {
				m.Data["bucket/existing-bucket"] = mock.NewFakeBucket(
					"existing-bucket", "bucket/existing-bucket", nil,
				)
			},
			wantName:         "",
			wantKey:          "",
			wantErr:          state.ErrResourceAlreadyExists,
			wantNumResources: 1,
		},
		{
			name:         "bucket with empty key - ok",
			resourceType: model.ResourceTypeBucket,
			newResource: mock.NewFakeBucket(
				"new-bucket",
				"",
				nil,
			),
			newResourceName: "new-bucket",
			setup: func(m *state.Memory) {
				m.Data["bucket/existing-bucket"] = mock.NewFakeBucket("existing-bucket", "bucket/existing-bucket", nil)
			},
			wantName:         "new-bucket",
			wantKey:          "bucket/new-bucket",
			wantNumResources: 2,
		},
		{
			name:             "new bucket with nil memory data",
			newResource:      mock.NewFakeBucket("new-bucket", "bucket/new-bucket", nil),
			newResourceName:  "new-bucket",
			wantName:         "new-bucket",
			wantKey:          "bucket/new-bucket",
			wantNumResources: 1,
		},
		{
			name:             "save different resource than memory resource type",
			resourceType:     model.ResourceTypeNetwork,
			newResource:      mock.NewFakeBucket("new-bucket", "bucket/new-bucket", nil),
			newResourceName:  "new-bucket",
			wantName:         "new-bucket",
			wantKey:          "bucket/new-bucket",
			wantNumResources: 1,
			wantErr:          state.ErrResourceTypeMismatch,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			memory := state.NewMemory()
			memory.ResourceType = tt.resourceType

			if tt.setup != nil {
				tt.setup(memory)
			}
			gotBucket, err := memory.Save(context.Background(), tt.newResource.Key(), tt.newResource)
			correctErr := errors.Is(err, tt.wantErr)

			if tt.wantErr != nil && !correctErr {
				t.Fatalf(
					"got error: %v, expected: %v",
					err,
					tt.wantErr,
				)
			} else if gotBucket != nil {
				if gotBucket.Key() != tt.wantKey {
					t.Fatalf("got bucket key %q, expected %q", gotBucket.Key(), tt.wantKey)
				}
				if gotBucket.Name() != tt.wantName {
					t.Fatalf("bucket key = got: %q, expected: %q", gotBucket.Name(), tt.wantName)
				}
				if len(memory.Data) != tt.wantNumResources {
					t.Fatalf("bucket num = got %d, expected %d", len(memory.Data), tt.wantNumResources)
				}
			}
		})
	}
}

func TestMemory_Delete(t *testing.T) {
	tests := []struct {
		name          string
		setup         func(m *state.Memory)
		deleteKey     string
		wantPreserved bool
		preservedKey  string
		resourceType  string
		wantErr       error
		wantExists    bool
	}{
		{
			name: "delete existing resource",
			setup: func(m *state.Memory) {
				m.Data["bucket/test-bucket"] = mock.NewFakeBucket("test-bucket", "bucket/test-bucket", nil)
			},
			deleteKey:     "bucket/test-bucket",
			wantPreserved: false,
			wantErr:       nil,
			wantExists:    false,
		},
		{
			name: "delete nonexistent resource",

			setup: func(m *state.Memory) {
				m.Data["bucket/test-bucket"] = mock.NewFakeBucket(
					"test-bucket", "bucket/test-bucket", nil,
				)
			},
			deleteKey:     "bucket/bucket-2",
			wantPreserved: false,
			wantErr:       state.ErrResourceNotFound,
			wantExists:    false,
		},
		{
			name: "delete one resource preserves others",
			setup: func(m *state.Memory) {
				m.Data["bucket/bucket-1"] = mock.NewFakeBucket("bucket-1", "bucket/bucket-1", nil)
				m.Data["bucket/bucket-2"] = mock.NewFakeBucket("bucket-2", "bucket/bucket-2", nil)
			},
			deleteKey:     "bucket/bucket-1",
			wantErr:       nil,
			wantExists:    false,
			wantPreserved: true,
			preservedKey:  "bucket/bucket-2",
		},
		{
			name: "delete one resource preserves others",
			setup: func(m *state.Memory) {
				m.Data["bucket/bucket-1"] = mock.NewFakeBucket("bucket-1", "bucket/bucket-1", nil)
				m.Data["bucket/bucket-2"] = mock.NewFakeBucket("bucket-2", "bucket/bucket-2", nil)
			},
			deleteKey:     "bucket/bucket-1",
			wantErr:       nil,
			wantPreserved: true,
			preservedKey:  "bucket/bucket-2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			memory := state.NewMemory()

			if tt.setup != nil {
				tt.setup(memory)
			}

			err := memory.Delete(context.Background(), tt.deleteKey)

			if err != nil && !strings.Contains(err.Error(), tt.wantErr.Error()) {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantPreserved {
				bucket := memory.Data[tt.preservedKey]
				if bucket == nil {
					t.Fatalf("preserved resource %q not found", tt.preservedKey)
				}
			}

			if tt.wantPreserved {
				if _, err := memory.Get(ctx, tt.preservedKey); err != nil {
					t.Fatalf("preserved resource exists = %v, want %v", err, tt.wantErr)
				}
			}
		})
	}
}

func TestMemory_List(t *testing.T) {
	tests := []struct {
		name         string
		setup        func(m *state.Memory)
		wantErr      error
		resourceType string
		wantCount    int
	}{
		{
			name:         "one resource to list",
			resourceType: "bucket",
			setup: func(m *state.Memory) {
				m.Data["bucket/bucket-1"] = mock.NewFakeBucket("bucket-1", "key-1", nil)
			},
			wantCount: 1,
			wantErr:   nil,
		},
		{
			name: "Multiple resources to list",
			setup: func(m *state.Memory) {
				m.Data["bucket/bucket-1"] = mock.NewFakeBucket("bucket-1", "bucket/bucket-1", nil)
				m.Data["bucket/bucket-2"] = mock.NewFakeBucket("bucket-2", "bucket/bucket-2", nil)
			},
			resourceType: "bucket",
			wantCount:    2,
		},
		{
			name:         "empty list",
			resourceType: "bucket",

			setup:     nil,
			wantCount: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			memory := state.NewMemory()

			if tt.setup != nil {
				tt.setup(memory)
			}
			listResource, err := memory.List(ctx)
			if err != nil {
				t.Fatalf("List() error = %v, wantErr %v", err, tt.wantErr)
			}
			if len(listResource) != tt.wantCount {
				t.Fatalf("listResource = %v, want %v", len(listResource), tt.wantCount)
			}
		})
	}
}

func TestMemory_Get(t *testing.T) {
	tests := []struct {
		name             string
		resourceKey      string
		setup            func(m *state.Memory)
		wantResourceType string
		wantBucketName   string
		wantErr          error
	}{
		{
			name:        "load existing resource",
			resourceKey: "bucket/bucket-1",
			setup: func(m *state.Memory) {
				m.Data["bucket/bucket-1"] = mock.NewFakeBucket("bucket-1", "key-1", nil)
			},
			wantResourceType: model.ResourceTypeBucket,
			wantBucketName:   "bucket-1",
			wantErr:          nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			memory := state.NewMemory()
			if tt.setup != nil {
				tt.setup(memory)
			}

			got, _ := memory.Get(ctx, tt.resourceKey)

			if got.Type() != tt.wantResourceType {
				t.Fatalf("Load() = %v, want %v", got.Type(), tt.wantResourceType)
			}
		})
	}
}
