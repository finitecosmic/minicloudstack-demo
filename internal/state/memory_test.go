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
		name           string
		newBucket      model.Resource
		newBucketName  string
		err            error
		memoryData     map[string]model.Resource
		resourceType   string
		wantName       string
		wantKey        string
		wantNumBuckets int
		wantErr        error
	}{
		{
			name:      "valid save bucket",
			newBucket: mock.NewFakeBucket("new-bucket", "bucket/new-bucket", nil),
			memoryData: map[string]model.Resource{
				"bucket/existing-bucket": model.Bucket{
					SpecData: model.BucketSpec{},
				},
			},
			wantName:       "new-bucket",
			wantKey:        "bucket/new-bucket",
			wantNumBuckets: 2,
			wantErr:        nil,
		},
		{
			name: "valid save bucket",
			newBucket: mock.NewFakeBucket(
				"new-bucket", "bucket/new-bucket", nil,
			),

			newBucketName:  "new-bucket",
			memoryData:     nil,
			wantName:       "new-bucket",
			wantKey:        "bucket/new-bucket",
			wantNumBuckets: 1,
			wantErr:        nil,
		},
		{
			name: "save already existing bucket error",
			newBucket: mock.NewFakeBucket(
				"existing-bucket", "bucket/existing-bucket", nil,
			),
			resourceType:  model.ResourceTypeBucket,
			newBucketName: "bucket/existing-bucket",
			memoryData: map[string]model.Resource{
				"bucket/existing-bucket": model.Bucket{
					ResourceType: model.ResourceTypeBucket,
					BucketName:   "existing-bucket",
					SpecData:     model.BucketSpec{},
				},
			},
			wantName:       "",
			wantKey:        "",
			wantErr:        state.ErrResourceAlreadyExists,
			wantNumBuckets: 1,
		},
		{
			name: "bucket with empty key - ok",
			newBucket: mock.NewFakeBucket(
				"new-bucket",
				"",
				nil,
			),
			newBucketName: "new-bucket",
			memoryData: map[string]model.Resource{
				"bucket/existing-bucket": model.Bucket{
					BucketName: "existing-bucket",
					SpecData: model.BucketSpec{
						Region: "us-east-1",
					},
				},
			},
			wantName:       "new-bucket",
			wantKey:        "bucket/new-bucket",
			wantNumBuckets: 2,
			wantErr:        nil, // replace with expected error
		},
		{
			name:           "new bucket with nil memory data",
			newBucket:      mock.NewFakeBucket("new-bucket", "bucket/new-bucket", nil),
			newBucketName:  "new-bucket",
			memoryData:     nil,
			wantName:       "new-bucket",
			wantKey:        "bucket/new-bucket",
			wantNumBuckets: 1,
			wantErr:        nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			memory := state.NewMemory("bucket", tt.memoryData)

			gotBucket, err := memory.Save(context.Background(), tt.newBucket.Key(), tt.newBucket)
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
				if len(memory.Data) != tt.wantNumBuckets {
					t.Fatalf("bucket num = got %d, expected %d", len(memory.Data), tt.wantNumBuckets)
				}
			}
		})
	}
}

func TestMemory_Delete(t *testing.T) {
	tests := []struct {
		name          string
		memoryData    map[string]model.Resource
		deleteKey     string
		wantPreserved bool
		preservedKey  string
		resourceType  string
		wantErr       error
		wantExists    bool
	}{
		{
			name: "delete existing resource",
			memoryData: map[string]model.Resource{
				"bucket/test-bucket": &mock.FakeBucket{
					BucketName:   "test-bucket",
					BucketKey:    "bucket/test-bucket",
					SpecData:     mock.FakeBucketSpec{},
					BucketExists: false,
				},
			},
			deleteKey:     "bucket/test-bucket",
			wantPreserved: false,
			wantErr:       nil,
			wantExists:    false,
		},
		{
			name: "delete nonexistent resource",

			memoryData: map[string]model.Resource{
				"bucket/test-bucket": mock.NewFakeBucket("test-bucket", "bucket-1", nil),
			},
			deleteKey:     "bucket/bucket-2",
			wantPreserved: false,
			wantErr:       state.ErrResourceNotFound,
			wantExists:    false,
		},
		{
			name: "delete one resource preserves others",
			memoryData: map[string]model.Resource{
				"bucket/bucket-1": mock.NewFakeBucket("bucket-1", "key-1", nil),
				"bucket/bucket-2": mock.NewFakeBucket("bucket-2", "key-2", nil),
			},
			deleteKey:     "bucket/bucket-1",
			wantErr:       nil,
			wantExists:    false,
			wantPreserved: true,
			preservedKey:  "bucket/bucket-2",
		},
		{
			name: "delete one resource preserves others",
			memoryData: map[string]model.Resource{
				"bucket/bucket-1": mock.NewFakeBucket("bucket-1", "key-1", nil),
				"bucket/bucket-2": mock.NewFakeBucket("bucket-2", "key-2", nil),
			},
			deleteKey:     "bucket/bucket-1",
			wantErr:       nil,
			wantExists:    false,
			wantPreserved: true,
			preservedKey:  "bucket/bucket-2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			memory := state.NewMemory("", tt.memoryData)

			for _, v := range tt.memoryData {
				memory.Data[v.Key()] = v
			}

			err := memory.Delete(context.Background(), tt.deleteKey)

			if err != nil && !strings.Contains(err.Error(), tt.wantErr.Error()) {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}

			_, err = memory.Get(ctx, tt.deleteKey)

			exists := err == nil
			if exists != tt.wantExists {
				t.Fatalf("resource exists = %v, want %v", exists, tt.wantExists)
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
		memoryData   map[string]model.Resource
		wantErr      error
		resourceType string
		wantCount    int
	}{
		{
			name:         "one resource to list",
			resourceType: "bucket",
			memoryData: map[string]model.Resource{
				"bucket/bucket-1": mock.NewFakeBucket("bucket-1", "key-1", nil),
			},
			wantCount: 1,
			wantErr:   nil,
		},
		{
			name: "Multiple resources to list",
			memoryData: map[string]model.Resource{
				"bucket/bucket-1": mock.NewFakeBucket("bucket-1", "key-1", nil),
				"bucket/bucket-2": mock.NewFakeBucket("bucket-2", "key-2", nil),
			},
			resourceType: "bucket",
			wantCount:    2,
		},
		{
			name:         "empty list",
			resourceType: "bucket",

			memoryData: nil,
			wantCount:  0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			memory := state.NewMemory(tt.resourceType, tt.memoryData)
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

func TestMemory_Load(t *testing.T) {
	tests := []struct {
		name           string
		getBucket      string
		memoryData     map[string]model.Resource
		wantBucketName string
		wantErr        error
	}{
		{
			name:      "load existing resource",
			getBucket: "bucket/bucket-1",
			memoryData: map[string]model.Resource{
				"bucket/bucket-1": mock.NewFakeBucket("bucket-1", "key-1", nil),
			},
			wantBucketName: "bucket-1",
			wantErr:        nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			memory := state.NewMemory(tt.getBucket, tt.memoryData)
		})
	}
}
