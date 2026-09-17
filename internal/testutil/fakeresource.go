package testutil

import "time"

type FakeBucketConfig struct {
	Name       string
	Region     string
	Versioning bool
	Encryption string
	Tags       map[string]string
}

type FakeBucket struct {
	config    FakeBucketConfig
	createdAt time.Time
}

func NewFakeBucket(config FakeBucketConfig) *FakeBucket {
	return &FakeBucket{
		config:    config,
		createdAt: time.Now(),
	}
}

func (b *FakeBucket) Key() string {
	return "bucket/" + b.config.Name
}
