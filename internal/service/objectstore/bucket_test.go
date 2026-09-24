package objectstore

import (
	"context"
	"errors"
	"minicloudstack/internal/model"
	"minicloudstack/internal/testutil/mock"
	"testing"
)

func TestService_CreateBucket(t *testing.T) {
	tests := []struct {
		name         string
		bucketSpec   model.BucketSpec
		bucketName   string
		bucketKey    string
		setup        func(m *mock.FakeMemory)
		resourceType string
		wantName     string
		wantKey      string
		wantErr      error
	}{
		{
			name:       "create bucket",
			bucketName: "new-bucket",
			bucketKey:  "bucket/new-bucket",
			bucketSpec: model.BucketSpec{
				SpecName: "new_bucket_spec",
			},
			setup: func(m *mock.FakeMemory) {
				m.Data["bucket/existing-bucket"] = mock.NewFakeBucket("existing-bucket", "bucket/existing-bucket", nil)
			},
			wantName: "new-bucket",
			wantKey:  "bucket/new-bucket",
			wantErr:  nil,
		},
		{
			name: "create bucket empty name",
			bucketSpec: model.BucketSpec{
				SpecName: "new_bucket_spec",
			},
			bucketName:   "",
			resourceType: model.ResourceTypeBucket,
			setup: func(m *mock.FakeMemory) {
				m.Data["bucket/existing-bucket"] = mock.NewFakeBucket("existing-data", "bucket/existing-data", nil)
			},
			wantErr: ErrBucketNameRequired,
		},
		{
			name: "create bucket already exists",
			bucketSpec: model.BucketSpec{
				SpecName: "existing_test-bucket_spec",
			},
			bucketKey:  "bucket/existing-bucket",
			bucketName: "existing_bucket",
			setup: func(m *mock.FakeMemory) {
				m.Data["bucket/existing-bucket"] = mock.NewFakeBucket("existing-bucket", "bucket/existing-bucket", nil)
			},
			wantKey:  "",
			wantName: "",
			wantErr:  ErrBucketAlreadyExists,
		},
		{
			name: "create bucket empty key",
			bucketSpec: model.BucketSpec{
				SpecName: "existing_test-bucket_spec",
			},
			bucketKey:  "",
			bucketName: "existing_bucket",

			setup: func(m *mock.FakeMemory) {
				m.Data["bucket/existing-bucket"] = mock.NewFakeBucket("existing-bucket", "bucket/existing-bucket", nil)
			},
			wantKey:  "bucket/existing_bucket",
			wantName: "existing_bucket",
			wantErr:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			fakeMemory := mock.NewFakeMemory()
			svc := New(fakeMemory)

			if tt.setup != nil {
				tt.setup(fakeMemory)
			}

			gotBucket, err := svc.CreateBucket(ctx, tt.bucketName, tt.bucketKey, tt.bucketSpec)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("CreateBucket() error = %v, wantErr %v", err, tt.wantErr)
			}

			if gotBucket != nil && gotBucket.Name() != tt.wantName {
				t.Fatalf("got: %s expected: %s", gotBucket.Name(), tt.wantName)
			}
			if tt.wantErr == nil && gotBucket.Key() != tt.wantKey {
				t.Fatalf("got: %s expected: %s", gotBucket.Key(), tt.wantKey)
			}

		})
	}

}

func TestService_ListBuckets(t *testing.T) {
	var tests = []struct {
		name            string
		setup           func(m *mock.FakeMemory)
		resourceType    string
		expectedBuckets int
	}{
		{
			name: "list buckets",
			setup: func(m *mock.FakeMemory) {
				m.Data["bucket/test"] = mock.NewFakeBucket("test", "bucket/test", nil)
				m.Data["bucket/backup"] = mock.NewFakeBucket("backup", "bucket/backup", nil)
			},
			expectedBuckets: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			memory := mock.NewFakeMemory()
			service := New(memory)
			if tt.setup != nil {
				tt.setup(memory)
			}

			buckets, _ := service.ListBuckets(ctx)
			if len(buckets) != tt.expectedBuckets {
				t.Fatalf("got: %d, expected: %d", len(buckets), tt.expectedBuckets)
			}
		})
	}
}

func TestService_DeleteBucket(t *testing.T) {
	tests := []struct {
		name          string
		setup         func(m *mock.FakeMemory)
		memoryData    map[string]model.Resource
		deleteKey     string
		wantErr       error
		wantPreserved []string
	}{
		{
			name: "delete bucket",
			memoryData: map[string]model.Resource{
				"bucket/test": &model.Bucket{
					BucketName: "bucket_test",
				},
			},
			deleteKey:     "bucket/test",
			wantErr:       nil,
			wantPreserved: nil,
		},
		{
			name: "delete bucket preserve one",
			memoryData: map[string]model.Resource{
				"bucket/test": &model.Bucket{
					BucketName: "bucket_test",
				},
				"bucket/test2": &model.Bucket{
					BucketName: "bucket_test",
				},
			},
			deleteKey:     "bucket/test",
			wantErr:       nil,
			wantPreserved: []string{"bucket/test2"},
		},
		{
			name: "delete bucket preserve multiple",
			setup: func(m *mock.FakeMemory) {
				m.Data["bucket/test1"] = mock.NewFakeBucket(
					"test1",
					"bucket/test1",
					nil,
				)

				m.Data["bucket/test2"] = mock.NewFakeBucket(
					"test2",
					"bucket/test2",
					nil,
				)

				m.Data["bucket/test3"] = mock.NewFakeBucket(
					"test2",
					"bucket/test3",
					nil,
				)
			},
			memoryData: map[string]model.Resource{
				"bucket/test": &model.Bucket{
					BucketName: "bucket_test",
				},
				"bucket/test2": &model.Bucket{
					BucketName: "bucket_test",
				},
				"bucket/test3": &model.Bucket{
					BucketName: "bucket_test",
				},
			},
			deleteKey:     "bucket/test",
			wantErr:       nil,
			wantPreserved: []string{"bucket/test2", "bucket/test3"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			memory := mock.NewFakeMemory()
			if tt.setup != nil {
				tt.setup(memory)
			}
			svc := New(memory)
			err := svc.DeleteBucket(ctx, tt.deleteKey)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("DeleteBucket() error = %v, wantErr %v", err, tt.wantErr)
			}

			for _, key := range tt.wantPreserved {
				if _, err := memory.Get(ctx, key); err != nil {
					t.Errorf("preserved resource %q does not exist: %v", key, err)
				}
			}
		})
	}
}
func TestService_GetBucket(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(m *mock.FakeMemory)
		key     string
		wantErr error
	}{
		{
			name: "get existing bucket",
			key:  "bucket/bucket_test",
			setup: func(m *mock.FakeMemory) {
				m.Data["bucket/bucket_test"] = mock.NewFakeBucket("bucket/bucket_test", "bucket_test", nil)
			},

			wantErr: nil,
		},
		{
			name: "get non existing bucket",
			key:  "bucket/test1",
			setup: func(m *mock.FakeMemory) {
				m.Data["bucket/bucket_test"] = mock.NewFakeBucket("bucket/bucket_test", "bucket_test", nil)
			},

			wantErr: ErrBucketNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			fakeMemory := mock.NewFakeMemory()
			svc := New(fakeMemory)
			if tt.setup != nil {
				tt.setup(fakeMemory)
			}

			bucket, err := svc.GetBucket(ctx, tt.key)

			// Assert: error
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("GetBucket() error = %v, wantErr %v", err, tt.wantErr)
			}
			// Assert: result
			if bucket != nil && bucket.Name() != tt.key {
				t.Errorf("got: %s, want: %s", bucket.Name(), tt.key)
			}
		})
	}
}
