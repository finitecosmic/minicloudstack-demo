package objectstore

import (
	"context"
	"errors"
	"minicloudstack/internal/model"
	"minicloudstack/internal/state"
	"testing"
)

func TestService_CreateBucketName(t *testing.T) {
	tests := []struct {
		name          string
		newBucketSpec model.BucketSpec
		specData      map[string]model.Resource
		wantName      string
		wantErr       error
	}{
		{
			name: "create bucket",
			newBucketSpec: model.BucketSpec{
				Name: "new-bucket",
			},
			specData: nil,
			wantName: "new-bucket",
			wantErr:  nil,
		},
		{
			name: "create bucket empty name",
			newBucketSpec: model.BucketSpec{
				Name: "new-bucket",
			},
			specData: map[string]model.Resource{},
			wantName: "new-bucket",
			wantErr:  nil,
		},
		{
			name: "create bucket already exists",
			newBucketSpec: model.BucketSpec{
				Key:  "existing/test-bucket",
				Name: "test-bucket",
			},

			specData: map[string]model.Resource{
				"existing-bucket": model.Bucket{
					SpecData: model.BucketSpec{
						Name: "existing-bucket",
					},
				},
			},
			wantName: "",
			wantErr:  ErrBucketAlreadyExists,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			memory := state.NewMemory(tt.specData)
			state.NewMemory()
			service := New(memory)

			newBucket := model.NewBucket(tt.newBucketSpec)

			gotBucket, err := service.CreateBucket(ctx, newBucket.SpecData)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("CreateBucket() error = %v, wantErr %v", err, tt.wantErr)
			}

			gotName := gotBucket.Name()
			if gotBucket != nil && gotBucket.Name() != tt.wantName {
				t.Fatalf("got: %s expected: %s create bucket key err: %v", gotName, tt.wantName, err)
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
					SpecData: model.BucketSpec{
						Name:   "test",
						Region: "us-east-1",
					},
				},
				"bucket/backup": &model.Bucket{
					SpecData: model.BucketSpec{
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
