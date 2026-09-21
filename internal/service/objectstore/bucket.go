package objectstore

import (
	"context"
	"errors"
	"minicloudstack/internal/model"
)

func (s *Service) CreateBucket(ctx context.Context, bucketSpec model.BucketSpec) (model.Resource, error) {
	newBucket := model.NewBucket(bucketSpec)

	name := bucketSpec.Name
	savedBucket, err := s.state.Save(ctx, name, newBucket)
	if err != nil {
		if errors.Is(err, ErrBucketAlreadyExists) {
			return savedBucket, err
		}

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
