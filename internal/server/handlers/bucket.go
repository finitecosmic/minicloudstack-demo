package handlers

import (
	"encoding/json"
	"encoding/xml"
	"minicloudstack/internal/model"
	"minicloudstack/internal/service/objectstore"
	"net/http"
	"path"
	"strings"
)

type BucketHandler struct {
	service *objectstore.Service
}

func NewBucketHandler(objectStore *objectstore.Service) *BucketHandler {
	return &BucketHandler{
		service: objectStore,
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
	Buckets []model.Bucket `json:"buckets"`
}

type ListBucketsXMLResponse struct {
	XMLName xml.Name       `xml:"ListAllMyBucketsResult"`
	Buckets []model.Bucket `xml:"Buckets>Bucket"`
}

func (ListBucketsJSONResponse) isListBucketsResponse() {}
func (ListBucketsXMLResponse) isListBucketsResponse()  {}

type BucketResponse struct {
	Name     string `json:"name" xml:"Name"`
	Region   string `json:"region" xml:"Region"`
	CreateAt string `json:"create_at" xml:"CreateAt"`
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

func (h *BucketHandler) Create(w http.ResponseWriter, r *http.Request) {
	var contentType string
	var name string
	var region string
	var err error
	var bucket model.Resource

	ctx := r.Context()

	if contentType, err = parseContentType(r); err != nil {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(err.Error()))
		return
	}
	switch contentType {
	case "application/json":
		var req CreateBucketJSONRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, ErrInvalidJSON.Error(), http.StatusBadRequest)
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

		name = path.Base(r.URL.Path)
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

	bucketSpec := model.NewBucketSpec(name, region)

	bucket, err = h.service.CreateBucket(ctx, bucket.Key(), bucketSpec)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Location", "/buckets/"+bucket.Key())

}

func (h *BucketHandler) List(w http.ResponseWriter, r *http.Request) {
	var buckets []model.Bucket
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

func (h *BucketHandler) Get(w http.ResponseWriter, r *http.Request) {
	var req GetBucketRequest
	ctx := r.Context()

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
	}

	bucket, err := h.service.GetBucket(ctx, req.Name)

	if bucket == nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bucket)
}

func (h *BucketHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	name := r.PathValue("name")

	err := h.service.DeleteBucket(ctx, name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
}
