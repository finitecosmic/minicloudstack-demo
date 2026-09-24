package mock

import (
	"minicloudstack/internal/model"
	"time"
)

type FakeBucket struct {
	BucketName   string
	BucketKey    string
	SpecData     FakeBucketSpec
	BucketExists bool
	err          error
}

func NewFakeBucket(name string, key string, err error) *FakeBucket {
	return &FakeBucket{
		BucketName:   name,
		BucketKey:    key,
		BucketExists: true,
		SpecData:     FakeBucketSpec{},
		err:          err,
	}
}

func (b *FakeBucket) Key() string {
	return "bucket/" + b.BucketName
}

func (b *FakeBucket) Name() string {
	return b.BucketName
}

func (b *FakeBucket) Dependencies() []string {
	return []string{}
}

func (b *FakeBucket) Type() string {
	return "bucket"
}
func (b *FakeBucket) Spec() model.Spec {
	return &b.SpecData
}

type FakeBucketSpec struct {
	Key    string
	Name   string
	Region string
	err    error
}

func (b *FakeBucketSpec) Validate() error {
	return nil
}

func (b *FakeBucketSpec) Get() model.Spec {
	return b
}

func (b *FakeBucketSpec) Error() string {
	//TODO implement me
	panic("implement me")
}

func (b *FakeBucket) SetUpdatedAt(time time.Time) {}

func (b *FakeBucketSpec) Dependencies() []string {
	return nil
}

type BucketOutput struct {
	Bucket *FakeBucket
	Error  error
}
