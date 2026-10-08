package testutil

import (
	"context"
	"time"
)

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

func (b *FakeBucket) Name(ctx context.Context, name string) (string, error) {
	return b.config.Name, nil
}
