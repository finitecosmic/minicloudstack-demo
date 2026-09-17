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
		bucketConfig BucketConfig
		expectedName string
		expectedErr  error
	}{
		{
			name: "create bucket",
			bucketConfig: BucketConfig{
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
			if bucket.Key() != gotName {
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
				"bucket/test": &Bucket{
					Config: BucketConfig{
						Name:   "test",
						Region: "us-east-1",
					},
				},
				"bucket/backup": &Bucket{
					Config: BucketConfig{
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

//
//func TestService_CreateBucketStateError(t *testing.T) {
//	state := testutil.NewFakeState()
//	service := New(state)
//
//	var tests = []struct {
//		createBucketInput CreateBucketInput
//		region            string
//		currentBuckets    []model.Bucket
//		expectedError     error
//	}{
//		{
//			createBucketInput: CreateBucketInput{
//				Name:   "test-bucket",
//				Region: "us-east-1",
//			},
//			currentBuckets: []model.Bucket{
//				{
//					Name:      "test-bucket",
//					Region:    "us-west-2",
//					CreatedAt: time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC),
//				},
//				{
//					Name:      "logs",
//					Region:    "us-east-1",
//					CreatedAt: time.Date(2026, 1, 2, 12, 0, 0, 0, time.UTC),
//				},
//				{
//					Name:      "backups",
//					Region:    "us-east-1",
//					CreatedAt: time.Date(2026, 1, 3, 12, 0, 0, 0, time.UTC),
//				},
//			},
//			expectedError: errors.New("error creating bucket: test-bucket already exists"),
//		},
//	}
//	for _, tt := range tests {
//		_, err := service.CreateBucket(context.Background(), tt.createBucketInput)
//
//		if !errors.Is(err, tt.expectedError) {
//			t.Fatal("expected", tt.expectedError, "got", err)
//		}
//
//	}
//}
