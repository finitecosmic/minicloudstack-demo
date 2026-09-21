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

func TestBucketHandler_CreateBucket(t *testing.T) {
	var tests = []struct {
		name             string
		httpMethod       string
		body             string
		headers          map[string]string
		uri              string
		wantName         string
		wantCode         int
		wantErrorMessage error
	}{
		{
			name:       "json correct bucket name",
			httpMethod: http.MethodPost,
			uri:        "/buckets",
			body:       `{"name": "test", "region": "us-east-1"}`,
			headers: map[string]string{
				"Content-Type": "application/json",
			},
			wantName:         "test",
			wantCode:         http.StatusCreated,
			wantErrorMessage: nil,
		},
		{
			name:       "xml correct bucket name",
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
			wantCode:         http.StatusCreated,
			wantName:         "test-bucket",
			wantErrorMessage: nil,
		},
		{
			name:       "xml incorrect correct bucket name",
			httpMethod: http.MethodPost,
			uri:        "/buckets/",
			body: `
				<CreateBucketConfiguration>
					<LocationConstraint>us-west-2</LocationConstraint>
				</CreateBucketConfiguration>
			`,
			headers: map[string]string{
				"Content-Type": "application/xml",
			},
			wantCode:         http.StatusBadRequest,
			wantName:         "",
			wantErrorMessage: ErrBucketNameRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.httpMethod, func(t *testing.T) {
			memory := state.NewMemory()
			objectstoreService := objectstore.New(memory)
			handler := NewBucketHandler(objectstoreService)

			req := httptest.NewRequest(
				tt.httpMethod,
				tt.uri,
				strings.NewReader(tt.body),
			)

			for key, value := range tt.headers {
				req.Header.Set(key, value)
			}
			rec := httptest.NewRecorder()
			handler.Create(rec, req)

			if rec.Code != tt.wantCode {
				t.Errorf("want %d, got %d", tt.wantCode, rec.Code)
			}
		})
	}
}

func TestBucketHandler_DeleteBucket(t *testing.T) {
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
			strings.NewReader(tt.body),
		)

		for key, value := range tt.headers {
			req.Header.Set(key, value)
		}
		req.SetPathValue("name", tt.uri)

		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		gotBucket, err := handler.service.GetBucket(ctx, "test-bucket")
		if err != nil {
			t.Errorf("want no errors, got %v", err)
		}
		if gotBucket.Name() != tt.wantName {
			t.Errorf("want %s, got %s", tt.wantName, gotBucket.Name())
		}

	}
}

func TestBucketHandler_GetBucket(t *testing.T) {

}

func TestBucketHandler_ListBuckets(t *testing.T) {

}
