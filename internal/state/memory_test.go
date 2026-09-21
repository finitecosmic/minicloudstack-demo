// internal/state/memory_test.go
package state

import (
	"context"
	"errors"
	"minicloudstack/internal/model"
	"testing"
)

func TestMemory_Save(t *testing.T) {
	tests := []struct {
		name           string
		newBucketKey   string
		newBucketName  string
		memoryData     []model.Resource
		wantName       string
		wantKey        string
		wantNumBuckets int
		wantErr        error
	}{
		{
			name:          "valid save bucket",
			newBucketKey:  "bucket/new-bucket",
			newBucketName: "new-bucket",
			memoryData: []model.Resource{
				model.Bucket{
					SpecData: model.BucketSpec{
						Key:  "bucket/existing-bucket",
						Name: "existing-bucket",
					},
				},
			},
			wantName:       "new-bucket",
			wantKey:        "bucket/new-bucket",
			wantNumBuckets: 2,
			wantErr:        nil,
		},
		{
			name:           "valid save bucket",
			newBucketKey:   "bucket/new-bucket",
			newBucketName:  "new-bucket",
			memoryData:     nil,
			wantName:       "new-bucket",
			wantKey:        "bucket/new-bucket",
			wantNumBuckets: 1,
			wantErr:        nil,
		},
		{
			name:          "save already existing bucket",
			newBucketKey:  "bucket/existing-bucket",
			newBucketName: "existing-bucket",
			memoryData: []model.Resource{
				model.Bucket{
					SpecData: model.BucketSpec{
						Key:  "bucket/existing-bucket",
						Name: "existing-bucket",
					},
				},
			},
			wantName:       "",
			wantKey:        "",
			wantErr:        ErrResourceAlreadyExists,
			wantNumBuckets: 1,
		},
		{
			name:          "invalid bucket with empty key",
			newBucketKey:  "",
			newBucketName: "existing-bucket",
			memoryData: []model.Resource{
				model.Bucket{
					SpecData: model.BucketSpec{
						Key:  "bucket/existing-bucket",
						Name: "existing-bucket",
					},
				},
			},
			wantName:       "existing-bucket",
			wantKey:        "bucket/existing-bucket",
			wantNumBuckets: 1,
			wantErr:        nil,
		},
		{
			name:          "new bucket with nil memory data",
			newBucketKey:  "buckey/new-bucket",
			newBucketName: "new-bucket",
			memoryData: []model.Resource{
				model.Bucket{
					SpecData: model.BucketSpec{
						Key:  "bucket/existing-bucket",
						Name: "existing-bucket",
					},
				},
			},
			wantName:       "new-bucket",
			wantKey:        "bucket/new-bucket",
			wantNumBuckets: 2,
			wantErr:        nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			memory := NewMemory()

			for _, v := range tt.memoryData {
				memory.data[v.Key()] = v
			}

			newBucket := model.Bucket{
				SpecData: model.BucketSpec{
					Key:  tt.newBucketKey,
					Name: tt.newBucketName,
				},
			}
			gotBucket, err := memory.Save(context.Background(), tt.newBucketKey, newBucket)
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
					t.Fatalf("got bucket name %q, expected: %q", gotBucket.Name(), tt.wantName)
				}
				if len(memory.data) != tt.wantNumBuckets {
					t.Fatalf("got %d buckets, expected %d", len(memory.data), tt.wantNumBuckets)
				}
			}
		})
	}
}

//type testResource struct {
//	key string
//}
//
//func (r *testResource) Key() string {
//	return r.key
//}
//
//func TestMemoryEmpty(t *testing.T) {
//	m := NewMemory()
//
//	resources, err := m.List(context.Background())
//	if err != nil {
//		t.Fatalf("List() error = %v", err)
//	}
//
//	if resources == nil {
//		t.Fatal("List() returned nil map")
//	}
//
//	if len(resources) != 0 {
//		t.Fatalf("List() returned %d resources, want 0", len(resources))
//	}
//}
//
//func TestMemory_List(t *testing.T) {
//	m := NewMemory()
//
//	bucket := &testResource{key: "bucket/photos"}
//	blob := &testResource{key: "blob/photos/image.jpg"}
//
//	if err := m.Save(context.Background(), bucket); err != nil {
//		t.Fatalf("Save() error = %v", err)
//	}
//
//	if err := m.Save(context.Background(), blob); err != nil {
//		t.Fatalf("Save() error = %v", err)
//	}
//
//	resources, err := m.List(context.Background())
//	if err != nil {
//		t.Fatalf("List() error = %v", err)
//	}
//
//	if len(resources) != 2 {
//		t.Fatalf("List() returned %d resources, want 2", len(resources))
//	}
//
//	if got, ok := resources["bucket/photos"]; !ok || got != bucket {
//		t.Errorf("List() missing bucket/photos")
//	}
//
//	if got, ok := resources["blob/photos/image.jpg"]; !ok || got != blob {
//		t.Errorf("List() missing blob/photos/image.jpg")
//	}
//}
//
//func TestMemory_List_ReturnsCopy(t *testing.T) {
//	m := NewMemory()
//
//	resource := &testResource{key: "bucket/photos"}
//
//	if err := m.Save(context.Background(), resource); err != nil {
//		t.Fatalf("Save() error = %v", err)
//	}
//
//	resources, err := m.List(context.Background())
//	if err != nil {
//		t.Fatalf("List() error = %v", err)
//	}
//
//	delete(resources, "bucket/photos")
//
//	resourcesAgain, err := m.List(context.Background())
//	if err != nil {
//		t.Fatalf("second List() error = %v", err)
//	}
//
//	if len(resourcesAgain) != 1 {
//		t.Fatalf(
//			"modifying returned map changed internal state: got %d resources, want 1",
//			len(resourcesAgain),
//		)
//	}
//
//	if _, ok := resourcesAgain["bucket/photos"]; !ok {
//		t.Error("resource was removed from internal state")
//	}
//}
//
//func TestMemory_List_Concurrent(t *testing.T) {
//	m := NewMemory()
//
//	const count = 100
//
//	for i := 0; i < count; i++ {
//		resource := &testResource{
//			key: "resource/" + string(rune(i)),
//		}
//
//		if err, _ := m.Save(context.Background(), _, resource); err != nil {
//			t.Fatalf("Save() error = %v", err)
//		}
//	}
//
//	var wg sync.WaitGroup
//
//	for i := 0; i < 10; i++ {
//		wg.Add(1)
//
//		go func() {
//			defer wg.Done()
//
//			resources, err := m.List(context.Background())
//			if err != nil {
//				t.Errorf("List() error = %v", err)
//				return
//			}
//
//			if len(resources) != count {
//				t.Errorf(
//					"List() returned %d resources, want %d",
//					len(resources),
//					count,
//				)
//			}
//		}()
//	}
//
//	wg.Wait()
//}
