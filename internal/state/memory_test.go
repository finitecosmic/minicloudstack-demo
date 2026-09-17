// internal/state/memory_test.go
package state_test

import (
	"fmt"
	"minicloudstack/internal/state"
	"minicloudstack/internal/testutil"
	"testing"
)

func TestMemory_List_Bucket(t *testing.T) {
	memory := state.NewMemory()

	_ = testutil.NewFakeBucket(
		testutil.FakeBucketConfig{
			Name:       "",
			Region:     "",
			Versioning: false,
			Encryption: "",
			Tags:       nil,
		},
	)
	fmt.Print(memory)
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
