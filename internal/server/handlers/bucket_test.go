package handlers

import (
	"errors"
	"minicloudstack/internal/model"
	"minicloudstack/internal/service/objectstore"
	"minicloudstack/internal/state"
	"minicloudstack/internal/testutil/mock"
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
		setup            func(m state.Memory)
		wantName         string
		wantCode         int
		wantErrorMessage error
	}{
		{
			name:       "json correct bucket name",
			httpMethod: http.MethodPost,
			uri:        "/bucket",
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
			uri:        "/bucket/test-bucket",
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
			uri:        "/bucket/",
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
			wantErrorMessage: ErrInvalidUri,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			memory := state.NewMemory()
			memory.ResourceType = model.ResourceTypeBucket
			objectstoreService := objectstore.New(memory)
			handler := NewBucketHandler(objectstoreService)

			reader := strings.NewReader(tt.body)
			req := httptest.NewRequest(tt.httpMethod, tt.uri, reader)

			for key, value := range tt.headers {
				req.Header.Set(key, value)
			}
			rec := httptest.NewRecorder()

			// handler Create
			handler.Create(rec, req)

			if rec.Code != tt.wantCode {
				t.Errorf("want %d, got %d", tt.wantCode, rec.Code)
			}
		})
	}
}

func TestBucketHandler_DeleteBucket(t *testing.T) {
	var tests = []struct {
		name             string
		httpMethod       string
		uri              string
		setup            func(*state.Memory)
		wantCode         int
		wantBucketCount  int
		wantErrorMessage string
	}{
		{
			name:       "delete existing bucket",
			httpMethod: http.MethodDelete,
			setup: func(m *state.Memory) {
				m.Data["bucket/testbucket"] = mock.NewFakeBucket(
					"testbucket",
					"bucket/testbucket",
					nil,
				)
			},
			uri:             "/bucket/testbucket",
			wantBucketCount: 0,
			wantCode:        http.StatusNoContent,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			memory := state.NewMemory()
			if tt.setup != nil {
				tt.setup(memory)
			}

			service := objectstore.New(memory)
			handler := NewBucketHandler(service)

			req := httptest.NewRequest(
				tt.httpMethod,
				tt.uri,
				strings.NewReader(""),
			)
			rec := httptest.NewRecorder()
			handler.Delete(rec, req)

			if rec.Code != tt.wantCode {
				t.Errorf("want %d, got %d", tt.wantCode, rec.Code)
			}
		})

	}
}

func TestBucketHandler_GetBucket(t *testing.T) {

}

func TestBucketHandler_ListBuckets(t *testing.T) {

}

func TestBucketHandler_validateXmlPath(t *testing.T) {
	var tests = []struct {
		name string
		uri  string

		wantError error
	}{
		{
			name:      "xml correct uri",
			uri:       "/bucket/test-bucket",
			wantError: nil,
		},
		{
			name:      " xml incorrect uri",
			uri:       "/buckets/",
			wantError: ErrInvalidUri,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateXmlPath(tt.uri)
			if !errors.Is(err, tt.wantError) {
				t.Errorf("want no errors, got %v", err)
			}

		})
	}
}
