package handlers

import (
	"encoding/json"
	"encoding/xml"
	"minicloudstack/internal/service/objectstore"
	"net/http"
	"strings"
)

type BucketHandler struct {
	objectStore *objectstore.Service
}

func NewBucketHandler(objectStore *objectstore.Service) *BucketHandler {
	return &BucketHandler{
		objectStore: objectStore,
	}
}

// CreateBucketRequest Create Bucket Request interface
type CreateBucketRequest interface {
	isCreateBucketRequest()
}
type CreateBucketXMLRequest struct {
	XMLName            xml.Name `xml:"CreateBucketConfiguration"`
	LocationConstraint string   `xml:"LocationConstraint"`
}

type CreateBucketJSONRequest struct {
	Name   string `json:"name"`
	Region string `json:"region"`
}

func (req *CreateBucketXMLRequest) isCreateBucketRequest()  {}
func (req *CreateBucketJSONRequest) isCreateBucketRequest() {}

// CreateBucketResponse Create Bucket Response Interface
type CreateBucketResponse interface {
	isCreateBucketResponse()
}

type CreateBucketJSONResponse struct {
	Name      string `json:"name"`
	Region    string `json:"region"`
	CreatedAt string `json:"created_at"`
}
type CreateBucketXMLResponse struct {
	XMLName xml.Name `xml:"CreateBucketConfiguration"`
}

func (CreateBucketJSONResponse) isCreateBucketResponse() {}
func (CreateBucketXMLResponse) isCreateBucketResponse()  {}

// ListBucketsResponse List Bucket Response
type ListBucketsResponse interface {
	isListBucketsResponse()
}
type ListBucketsJSONResponse struct {
	Buckets []objectstore.Bucket `json:"buckets"`
}

type ListBucketsXMLResponse struct {
	XMLName xml.Name             `xml:"ListAllMyBucketsResult"`
	Buckets []objectstore.Bucket `xml:"Buckets>Bucket"`
}

func (ListBucketsJSONResponse) isListBucketsResponse() {}
func (ListBucketsXMLResponse) isListBucketsResponse()  {}

type BucketResponse struct {
	Name     string `"json":"name" "xml":"name"`
	Region   string `"json":"region" "xml":"region"`
	CreateAt string `"json":"createAt" "xml":"createAt"`
}

type ListBucketRequest struct {
	Region string `json:"region"`
}

type GetBucketRequest struct {
	Name   string `json:"name"`
	Region string `json:"region"`
}

type DeleteBucketRequest struct {
	Name   string `json:"name"`
	Region string `json:"region"`
}

func (h *BucketHandler) CreateBucket(w http.ResponseWriter, r *http.Request) {
	var contentType string
	var err error
	var name string
	var region string

	ctx := r.Context()

	if contentType, err = parseContentType(r); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	switch contentType {
	case "application/json":
		var req CreateBucketJSONRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}
		if req.Name != "" {
			name = req.Name
		}
		if req.Region != "" {
			region = req.Region
		}

	case "application/xml", "text/xml":
		var req CreateBucketXMLRequest
		if err := xml.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid xml", http.StatusBadRequest)
			return
		}

		name = r.PathValue("name")
		region = req.LocationConstraint
	}

	if name == "" {
		http.Error(w, "missing bucket name", http.StatusBadRequest)
		return
	}
	if region == "" {
		http.Error(w, "missing bucket region", http.StatusBadRequest)
		return
	}

	bucket, err := h.objectStore.CreateBucket(ctx, objectstore.BucketConfig{})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Location", "/bucket"+bucket.Key())

}

func (h *BucketHandler) ListBuckets(w http.ResponseWriter, r *http.Request) {
	var buckets []objectstore.Bucket
	var response ListBucketsResponse

	accept := r.Header.Get("Accept")

	// XML Response
	if strings.Contains(accept, "application/xml") {
		response = ListBucketsXMLResponse{
			Buckets: buckets,
		}
	} else {
		response = ListBucketsJSONResponse{
			Buckets: buckets,
		}
	}

	writeResponse(w, r, response)
}

func (h *BucketHandler) GetBucket(w http.ResponseWriter, r *http.Request) {
	var req GetBucketRequest
	ctx := r.Context()

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
	}

	bucket, err := h.objectStore.GetBucket(ctx, req.Name)

	if bucket == nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bucket)
}

func (h *BucketHandler) DeleteBucket(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	name := r.PathValue("name")

	err := h.objectStore.DeleteBucket(ctx, name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
}
