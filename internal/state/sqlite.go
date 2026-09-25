package state

import (
	"context"
	"database/sql"
	"encoding/json"
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
	return nil, nil
}
func (s *SQLiteState) List(ctx context.Context) (map[string]model.Resource, error) {
	return make(map[string]model.Resource), nil
}

func (s *SQLiteState) Delete(ctx context.Context, key string) error {
	return nil
}

func (s *SQLiteState) Ready(ctx context.Context) (bool, error) {
	return true, nil
}
