package objectstore

import (
	"context"
	"minicloudstack/internal/model"
)

func (s *Service) CreateBucket(ctx context.Context, key string, spec model.BucketSpec) (model.Resource, error) {

	// validate
	if key == "" {
		return nil, ErrBucketNameRequired
	}
	exists, err := s.bucketExists(ctx, key)
	if exists {
		return nil, ErrBucketAlreadyExists
	}
	newBucket := model.NewBucket(key, spec)

	savedBucket, err := s.state.Save(ctx, key, newBucket)
	if err != nil {
		return nil, err
	}
	return savedBucket, nil
}

func (s *Service) ListBuckets(ctx context.Context) (map[string]model.Resource, error) {
	return s.state.List(ctx)
}

func (s *Service) GetBucket(ctx context.Context, bucketKey string) (model.Resource, error) {
	bucket, err := s.state.Get(ctx, bucketKey)

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

func (s *Service) bucketExists(ctx context.Context, bucketName string) (bool, error) {
	_, err := s.GetBucket(ctx, bucketName)
	if err != nil {
		return false, err
	}
	return true, nil
}
