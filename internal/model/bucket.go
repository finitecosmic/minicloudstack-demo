package model

type Bucket struct {
	Config     BucketConfig
	Versioning bool
}

func NewBucket(config BucketConfig) Bucket {
	return Bucket{
		Config: BucketConfig{
			Key:    "bucket/" + config.Name,
			Name:   config.Name,
			Region: config.Region,
		},
		Versioning: false,
	}
}

func (b Bucket) Key() string {
	return b.Config.Key
}

func (b Bucket) Name() string {
	return b.Config.Name
}

type BucketConfig struct {
	Key        string
	Name       string
	Region     string
	Versioning bool
	Encryption string
	Tags       map[string]string
}

func NewBucketConfig() BucketConfig {
	return BucketConfig{
		Name:       "",
		Region:     "",
		Versioning: false,
		Encryption: "",
		Tags:       nil,
	}

}

func (b *BucketConfig) Error() string {
	//TODO implement me
	panic("implement me")
}

type BucketOutput struct {
	Bucket *Bucket
	Error  error
}
