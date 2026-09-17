package objectstore

import (
	"context"
	"errors"
	"minicloudstack/internal/model"
)

type Bucket struct {
	Config     BucketConfig
	Versioning bool
}

func NewBucket(config BucketConfig) Bucket {
	return Bucket{
		Config:     BucketConfig{},
		Versioning: false,
	}
}

type BucketConfig struct {
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

func (s *Service) CreateBucket(ctx context.Context, bucketConfig BucketConfig) (model.Resource, error) {
	bucket := NewBucket(bucketConfig)

	name := bucketConfig.Name
	savedBucket, err := s.state.Save(ctx, name, bucket)
	if err != nil {
		if errors.Is(err, ErrBucketAlreadyExists) {
			return nil, ErrBucketAlreadyExists
		}

		return nil, err
	}

	return savedBucket, nil
}

func (s *Service) ListBuckets(ctx context.Context) (map[string]model.Resource, error) {
	return s.state.List(ctx)
}

func (s *Service) GetBucket(ctx context.Context, bucketName string) (model.Resource, error) {
	bucket, err := s.state.Get(ctx, bucketName)

	if err != nil {
		return nil, err
	}
	return bucket, nil
}

func (s *Service) DeleteBucket(ctx context.Context, bucketName string) error {
	err := s.state.Delete(ctx, bucketName)
	if err != nil {
		return err
	}
	return nil
}

func (b Bucket) Key() string {
	return b.Config.Name
}

func (b Bucket) Name() string {
	return b.Config.Name
}
