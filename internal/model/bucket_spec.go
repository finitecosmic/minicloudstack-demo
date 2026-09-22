package model

type BucketSpec struct {
	SpecKey      string
	SpecName     string
	Region       string
	Versioning   bool
	Encryption   string
	ResourceType string
	Tags         map[string]string
}

func NewBucketSpec(name string, key string) BucketSpec {
	return BucketSpec{
		SpecKey:      key,
		SpecName:     name,
		Region:       "",
		Versioning:   false,
		Encryption:   "",
		Tags:         nil,
		ResourceType: "bucketspec",
	}

}

func (b *BucketSpec) Validate() error {
	if b.Name() == "" {
		return ErrMissingBucketSpecName
	}
	if b.SpecKey == "" {
		return ErrMissingBucketKey
	}
	return nil
}

func (b *BucketSpec) Get() Spec {
	return b
}

func (b *BucketSpec) Name() string {
	return b.SpecName
}

func (b *BucketSpec) Error() string {
	//TODO implement me
	panic("implement me")
}

type BucketOutput struct {
	Bucket *Bucket
	Error  error
}
