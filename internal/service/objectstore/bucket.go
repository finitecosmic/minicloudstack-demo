package objectstore

import (
	"context"
	"minicloudstack/internal/model"
	"minicloudstack/internal/state"
	"regexp"
	"time"
)

type BucketService struct {
	state        state.State
	resourceType string
	ctx          context.Context //todo remove
}

func NewBucketService(ctx context.Context, state state.State) *BucketService {
	return &BucketService{
		state:        state,
		resourceType: model.ResourceTypeBucket,
		ctx:          ctx,
	}
}

func (s *Service) CreateBucket(ctx context.Context, name string, key string,
	spec model.BucketSpec) (model.Resource, error) {
	var err error
	var exists bool

	if name == "" {
		return nil, ErrBucketNameRequired
	}
	if key == "" {
		key = "bucket/" + name
	}

	err = validateBucketCreation(ctx, s.state, name, key)
	if err != nil {
		return nil, err
	}
	exists, err = bucketExists(ctx, s.state, key)
	if exists {
		return nil, ErrBucketAlreadyExists
	}

	newBucket := model.NewBucket(name, spec)
	newBucket.CreatedAt = time.Now().UTC()
	newBucket.UpdatedAt = time.Now().UTC()

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
	if bucket == nil {
		return nil, ErrBucketNotFound
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

func bucketExists(ctx context.Context, state state.State, key string) (bool, error) {
	bucket, err := state.Get(ctx, key)
	if err != nil || bucket == nil {
		return false, err
	}

	return true, nil
}

func validateBucketCreation(ctx context.Context, state state.State, name string, key string) error {
	if name == "" {
		return ErrBucketNameRequired
	}

	// correct key format
	regex := regexp.MustCompile("^bucket/[^/]+$")
	if !regex.MatchString(key) {
		return ErrInvalidKey
	}

	return nil
}
