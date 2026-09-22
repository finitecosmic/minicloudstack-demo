package model

import "time"

type Bucket struct {
	BucketName   string
	BucketKey    string
	SpecData     BucketSpec
	DependsOn    []string
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
		CreatedAt:    time.Now(),
		ResourceType: "bucket",
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
