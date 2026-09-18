package handlers

import (
	"context"
	"minicloudstack/internal/service/objectstore"
	"minicloudstack/internal/state"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestBucketHandler_CreateBucketJSONHTTPResponse(t *testing.T) {
	memory := state.NewMemory()
	objectstoreService := objectstore.New(memory)
	handler := NewBucketHandler(objectstoreService)

	var tests = []struct {
		httpMethod string
		body       string
		headers    map[string]string
		path       string
		want       int
	}{
		{
			httpMethod: http.MethodPost,
			path:       "/buckets",
			body:       `{"name": "test", "region": "us-east-1"}`,
			headers: map[string]string{
				"Content-Type": "application/json",
			},
			want: http.StatusCreated,
		},
	}

	for _, tt := range tests {
		req := httptest.NewRequest(
			tt.httpMethod,
			tt.path,
			strings.NewReader(string(tt.body)),
		)

		for key, value := range tt.headers {
			req.Header.Set(key, value)
		}
		rec := httptest.NewRecorder()
		handler.CreateBucket(rec, req)

		if rec.Code != tt.want {
			t.Errorf("want %d, got %d", tt.want, rec.Code)
		}
	}
}

func TestBucketHandler_CreateBucketXMLHTTPResponse(t *testing.T) {
	memory := state.NewMemory()
	objectstoreService := objectstore.New(memory)
	handler := NewBucketHandler(objectstoreService)

	var tests = []struct {
		httpMethod string
		body       string
		headers    map[string]string
		uri        string
		wantStatus int
		wantName   string
		wantKey    string
	}{
		{
			httpMethod: http.MethodPost,
			uri:        "/buckets/test-bucket",
			body: `
				<CreateBucketConfiguration>
					<LocationConstraint>us-west-2</LocationConstraint>
				</CreateBucketConfiguration>
			`,
			headers: map[string]string{
				"Content-Type": "application/xml",
			},
			wantStatus: http.StatusCreated,
			wantName:   "test-bucket",
			wantKey:    "/buckets/test-bucket",
		},
	}

	for _, tt := range tests {
		ctx := context.Background()
		req := httptest.NewRequest(
			tt.httpMethod,
			tt.uri,
			strings.NewReader(string(tt.body)),
		)

		for key, value := range tt.headers {
			req.Header.Set(key, value)
		}
		req.SetPathValue("name", tt.uri)

		rec := httptest.NewRecorder()

		handler.CreateBucket(rec, req)

		gotBucket, err := handler.objectStore.GetBucket(ctx, "test-bucket")
		if err != nil {
			t.Errorf("want no errors, got %v", err)
		}
		if gotBucket.Name() != tt.wantName {
			t.Errorf("want %s, got %s", tt.wantName, gotBucket.Name())
		}

	}
}

func TestBucketHandler_DeleteBucket(t *testing.T) {

}

func TestBucketHandler_GetBucket(t *testing.T) {

}

func TestBucketHandler_ListBuckets(t *testing.T) {

}
