package objectstore

import (
	"context"
	"minicloudstack/internal/model"
	"minicloudstack/internal/state"
	"testing"
)

func TestService_CreateBucket(t *testing.T) {
	tests := []struct {
		name         string
		bucketConfig model.BucketConfig
		expectedName string
		expectedErr  error
	}{
		{
			name: "create bucket",
			bucketConfig: model.BucketConfig{
				Name:   "test",
				Region: "us-east-1",
			},
			expectedName: "test",
			expectedErr:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			memory := state.NewMemory()
			service := New(memory)

			bucket, err := service.CreateBucket(ctx, tt.bucketConfig)
			if err != nil {
				t.Fatalf("create bucket err: %v", err)
			}
			gotName := bucket.Name()
			if bucket.Name() != tt.expectedName {
				t.Fatalf("got: %s expected: %s create bucket key err: %v", gotName, tt.expectedName, err)
			}
		})
	}

}

func TestService_ListBuckets(t *testing.T) {
	var tests = []struct {
		name            string
		buckets         map[string]model.Resource
		expectedBuckets int
	}{
		{
			name: "list buckets",
			buckets: map[string]model.Resource{
				"bucket/test": &model.Bucket{
					Config: model.BucketConfig{
						Name:   "test",
						Region: "us-east-1",
					},
				},
				"bucket/backup": &model.Bucket{
					Config: model.BucketConfig{
						Name:   "backup",
						Region: "us-west-2",
					},
				},
			},
			expectedBuckets: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			memory := state.NewMemory(tt.buckets)

			service := New(memory)

			buckets, _ := service.ListBuckets(ctx)
			if len(buckets) != tt.expectedBuckets {
				t.Fatalf("got: %d, expected: %d", len(buckets), tt.expectedBuckets)
			}
		})
	}
}
