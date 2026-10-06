package state

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"minicloudstack/internal/model"

	"github.com/google/uuid"
)

type SQLiteState struct {
	db *sql.DB
}

func NewSQLiteState(db *sql.DB) *SQLiteState {
	return &SQLiteState{db: db}
}

func (s *SQLiteState) InitDB(ctx context.Context) error {

	_, err := s.db.ExecContext(ctx, createResourceTable)
	if err != nil {
		return err
	}
	return nil
}

func (s *SQLiteState) Save(ctx context.Context, key string, resource model.Resource) error {
	var resourceType string

	if resource == nil {
		return ErrMissingResource
	}

	data, err := json.Marshal(resource)
	if err != nil {
		return err
	}

	id := uuid.NewMD5(uuid.Nil, []byte(key))
	resourceType = resource.Type()
	_, err = s.db.ExecContext(ctx,
		insertResource,
		id,
		key,
		1,
		resourceType,
		data,
	)
	if err != nil {
		return fmt.Errorf("save resource %q: %w", key, err)
	}
	return nil
}

func (s *SQLiteState) Get(ctx context.Context, key string) (model.Resource, error) {
	const query = `
		SELECT resource_type, version, data 
		FROM resources 
		WHERE key=?`

	var (
		resourceType string
		version      int
		data         []byte
	)
	err := s.db.QueryRowContext(ctx, query, key).Scan(&resourceType, &version, &data)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrResourceNotFound
		}
		return nil, err
	}

	resource, err := decodeResource(resourceType, version, data)
	if err != nil {
		return nil, err
	}
	return resource, nil
}

func (s *SQLiteState) List(ctx context.Context) (map[string]model.Resource, error) {

	const query = `
	SELECT resource_type, version, data from resources`

	rows, err := s.db.QueryContext(ctx, query)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	resources := map[string]model.Resource{}
	for rows.Next() {
		var (
			key          string
			resourceType string
			version      int
			data         []byte
		)
		err := rows.Scan(&key, &resourceType, &version, &data)
		if err != nil {
			return nil, err
		}
		resource, err := decodeResource(resourceType, version, data)
		if err != nil {
			return nil, err
		}
		resources[key] = resource
	}

	return resources, nil
}

func (s *SQLiteState) Delete(ctx context.Context, key string) error {
	const query = `DELETE FROM resources WHERE key=?`

	deleted, err := s.db.ExecContext(ctx, query, key)
	if err != nil {
		return err
	}
	rowsAffected, err := deleted.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrResourceNotFound
	}
	return nil
}

func (s *SQLiteState) Ready(ctx context.Context) (bool, error) {
	if err := s.db.PingContext(ctx); err != nil {
		return false, err
	}
	return true, nil
}

func decodeResource(resourceType string, version int, data []byte) (model.Resource, error) {
	switch resourceType {
	case "bucket":
		var bucket model.Bucket
		err := json.Unmarshal(data, &bucket)
		if err != nil {
			panic(err)
		}
		return bucket, nil
	default:
		return nil, fmt.Errorf(`unknown resource type %q`, resourceType)
	}
}
