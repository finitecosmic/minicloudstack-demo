package state

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"minicloudstack/internal/model"

	"github.com/google/uuid"
)

type SQLiteState struct {
	db     *sql.DB
	logger *slog.Logger
}

func NewSQLiteState(db *sql.DB) *SQLiteState {
	return &SQLiteState{
		db:     db,
		logger: slog.Default(),
	}
}

func (s *SQLiteState) InitDB(ctx context.Context) error {

	if _, err := s.db.ExecContext(ctx, createResourceTable); err != nil {
		return fmt.Errorf("initialize database: %w", err)
	}

	return nil
}

func (s *SQLiteState) Save(ctx context.Context, key string, resource model.Resource) error {

	if resource == nil {
		return ErrMissingResource
	}

	data, err := json.Marshal(resource)
	if err != nil {
		return err
	}

	id := uuid.NewMD5(uuid.Nil, []byte(key))
	resourceType := resource.Type()
	_, err = s.db.ExecContext(ctx,
		insertResource,
		id,
		key,
		resourceType,
		model.BucketVersion,
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
	rows, err := s.db.QueryContext(ctx, listResources)
	if err != nil {
		return nil, fmt.Errorf("list resources: %w", err)
	}

	defer rows.Close()

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
			return nil, fmt.Errorf("decode resource: %w", err)
		}
		resources[key] = resource
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate resources: %w", err)
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

func (s *SQLiteState) Migrate(ctx context.Context) error {
	var version int

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("%w %w", ErrBeginTransaction, err)
	}
	defer tx.Rollback() //rollback if not commit, atomic

	if err := tx.QueryRowContext(ctx, "PRAGMA user_version").Scan(&version); err != nil {
		return fmt.Errorf("%w %w", ErrQueryRowContext, err)
	}

	log.Printf("migrated sqlite state version %d", version)

	switch version {
	case schemaVersionV1:
		if _, err := tx.ExecContext(ctx, migrateV1toV2); err != nil {
			return fmt.Errorf("%w %w", ErrMigrateV1toV2, err)
		}

		if _, err := tx.ExecContext(ctx, "PRAGMA user_version = ?", currentSchemaVersion); err != nil {
			return fmt.Errorf("%w %w", ErrSetSchemaVersion, err)
		}
	case schemaVersionV2:
		s.logger.InfoContext(ctx, "%d", MsgSchemaAlreadyCurrent, currentSchemaVersion)
	default:
		return fmt.Errorf("%w %d", ErrUnsupportedVersion, version)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("%w %w", ErrTxCommit, err)
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
		return decodeBucket(data, version)

	default:
		return nil, fmt.Errorf(`unknown resource type %q`, resourceType)
	}
}

func decodeBucket(data []byte, version int) (model.Resource, error) {
	switch version {
	case 1:
		var bucket model.Bucket
		err := json.Unmarshal(data, &bucket)
		if err != nil {
			return nil, fmt.Errorf("%w: bucket: %v", ErrUnmarshalResource, err)
		}
		return bucket, nil
	default:
		return nil, fmt.Errorf(`%s %d`, ErrUnsupportedVersion, version)
	}
}
