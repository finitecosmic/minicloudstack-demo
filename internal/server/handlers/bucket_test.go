package handlers

import (
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
		httpMethod   string
		body         string
		headers      map[string]string
		path         string
		wantStatus   int
		wantLocation string
	}{
		{
			httpMethod: http.MethodPost,
			path:       "/buckets/test-bucket",
			body: `
				<CreateBucketConfiguration>
					<LocationConstraint>us-west-2</LocationConstraint>
				</CreateBucketConfiguration>
			`,
			headers: map[string]string{
				"Content-Type": "application/xml",
			},
			wantStatus:   http.StatusCreated,
			wantLocation: "/buckets/test-bucket",
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
		req.SetPathValue("name", tt.path)

		rec := httptest.NewRecorder()
		handler.CreateBucket(rec, req)

		if rec.Code != tt.wantStatus {
			t.Errorf("want %d, got %d", tt.wantStatus, rec.Code)
		}
		if got := rec.Header().Get("Location"); got != tt.wantLocation {
			t.Errorf("want %s, got %s", tt.wantLocation, got)
		}
	}
}

func TestBucketHandler_DeleteBucket(t *testing.T) {

}

func TestBucketHandler_GetBucket(t *testing.T) {

}

func TestBucketHandler_ListBuckets(t *testing.T) {

}
