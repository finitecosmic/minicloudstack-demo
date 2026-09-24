package model

import "time"

type Bucket struct {
	BucketName   string `json:"name" xml:"name""`
	BucketKey    string `json:"key" xml:"key"`
	SpecData     BucketSpec `json:"spec" xml:"spec"`
	DependsOn    []string   `json:"depends_on" xml:"depends_on"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	ResourceType string
}

func NewBucket(name string, spec BucketSpec) Bucket {
	return Bucket{
		BucketName:   name,
		BucketKey:    "bucket/" + name,
		SpecData:     spec,
		DependsOn:    []string{},
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
		ResourceType: ResourceTypeBucket,
	}
}

func (b Bucket) Key() string {
	return "bucket/" + b.BucketName
}

func (b Bucket) Name() string {
	return b.BucketName
}

func (b Bucket) Dependencies() []string {
	return b.DependsOn
}

func (b Bucket) Type() string {
	return "bucket"
}
func (b Bucket) Spec() Spec {
	return b.SpecData.Get()
}

func (b Bucket) SetUpdatedAt(updatedAt time.Time) {
	b.UpdatedAt = updatedAt
	return
}
