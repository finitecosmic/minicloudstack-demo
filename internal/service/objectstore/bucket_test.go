package objectstore

import (
	"context"
	"errors"
	"minicloudstack/internal/model"
	"minicloudstack/internal/state"
	"minicloudstack/internal/testutil/mock"
	"testing"
)

func TestService_CreateBucketName(t *testing.T) {
	tests := []struct {
		name         string
		bucketSpec   model.BucketSpec
		bucketName   string
		bucketKey    string
		resourceType string
		existingData map[string]model.Resource
		wantName     string
		wantErr      error
	}{
		{
			name:       "create bucket",
			bucketName: "new-bucket",
			bucketKey:  "bucket/new-bucket",
			bucketSpec: model.BucketSpec{
				SpecName: "new_bucket_spec",
			},
			existingData: nil,
			wantName:     "new-bucket",
			wantErr:      nil,
		},
		{
			name: "create bucket empty name",
			bucketSpec: model.BucketSpec{
				SpecName: "new_bucket_spec",
			},
			bucketName:   "",
			resourceType: model.BucketResourceType,
			existingData: map[string]model.Resource{},
			wantName:     "new-bucket",
			wantErr:      nil,
		},
		{
			name: "create bucket already exists",
			bucketSpec: model.BucketSpec{
				SpecName: "existing_test-bucket_spec",
			},
			bucketName: "existing_bucket",
			existingData: map[string]model.Resource{
				"existing-bucket": model.Bucket{
					BucketName: "existing_bucket",
				},
			},
			wantName: "",
			wantErr:  ErrBucketAlreadyExists,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			fakeMemory := mock.NewFakeMemory(tt.existingData)
			svc := New(fakeMemory)
			gotBucket, err := svc.CreateBucket(ctx, tt.bucketKey, tt.bucketSpec)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("CreateBucket() error = %v, wantErr %v", err, tt.wantErr)
			}

			if gotBucket != nil && gotBucket.Name() != tt.wantName {
				t.Fatalf("got: %s expected: %s", gotBucket.Name(), tt.wantName)
			}

		})
	}

}

func TestService_ListBuckets(t *testing.T) {
	var tests = []struct {
		name            string
		buckets         map[string]model.Resource
		resourceType    string
		expectedBuckets int
	}{
		{
			name: "list buckets",
			buckets: map[string]model.Resource{
				"bucket/test": &model.Bucket{
					SpecData: model.BucketSpec{
						SpecName: "spec_test",
						Region:   "us-east-1",
					},
				},
				"bucket/backup": &model.Bucket{
					SpecData: model.BucketSpec{
						SpecName: "spec_backup",
						Region:   "us-west-2",
					},
				},
			},
			expectedBuckets: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			memory := state.NewMemory(model.BucketResourceType, tt.buckets)

			service := New(memory)

			buckets, _ := service.ListBuckets(ctx)
			if len(buckets) != tt.expectedBuckets {
				t.Fatalf("got: %d, expected: %d", len(buckets), tt.expectedBuckets)
			}
		})
	}
}
