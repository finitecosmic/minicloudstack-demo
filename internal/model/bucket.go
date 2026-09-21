package model

type Bucket struct {
	SpecData  BucketSpec
	DependsOn []string
}

func NewBucket(config BucketSpec) Bucket {
	return Bucket{
		SpecData: BucketSpec{
			Key:    "bucket/" + config.Name,
			Name:   config.Name,
			Region: config.Region,
		},
		DependsOn: []string{},
	}
}

func (b Bucket) Key() string {

	return "bucket/" + b.SpecData.Name
}

func (b Bucket) Name() string {
	return b.SpecData.Name
}

func (b Bucket) GetSpec() Spec {
	return b.SpecData
}

func (b Bucket) GetDependencies() []string {
	return b.DependsOn
}

func (b Bucket) ResourceType() string {
	return "bucket"
}

type BucketSpec struct {
	Key        string
	Name       string
	Region     string
	Versioning bool
	Encryption string
	Tags       map[string]string
}

func NewBucketSpec() BucketSpec {
	return BucketSpec{
		Name:       "",
		Region:     "",
		Versioning: false,
		Encryption: "",
		Tags:       nil,
	}

}

func (b *BucketSpec) Error() string {
	//TODO implement me
	panic("implement me")
}

type BucketOutput struct {
	Bucket *Bucket
	Error  error
}
