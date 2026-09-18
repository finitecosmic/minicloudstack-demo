package objectstore

import (
	"context"
	"errors"
	"minicloudstack/internal/model"
)

func (s *Service) CreateBucket(ctx context.Context, bucketConfig model.BucketConfig) (model.Resource, error) {
	bucket := model.NewBucket(bucketConfig)

	name := bucketConfig.Name
	savedBucket, err := s.state.Save(ctx, name, bucket)
	if err != nil {
		if errors.Is(err, ErrBucketAlreadyExists) {
			return savedBucket, ErrBucketAlreadyExists
		}

		return savedBucket, err
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
