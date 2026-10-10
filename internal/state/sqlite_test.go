package state_test

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"minicloudstack/internal/model"
	"minicloudstack/internal/state"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func TestSQLiteState_InitDB(t *testing.T) {
	var tableName string

	ctx := context.Background()
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal("failed to open database:", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Error("failed to close database:", err)
		}
	}()

	db.SetMaxOpenConns(1)
	s := state.NewSQLiteState(db)

	err = s.InitDB(context.Background())
	if err != nil {
		t.Fatalf("InitDB err: %v", err)
	}

	err = db.QueryRowContext(ctx, `
		SELECT name
		FROM sqlite_master 
		where type='table' and 
	    name='resources'
		`).Scan(&tableName)

	if err != nil {
		t.Errorf("sqlite init table: %v", err)
	}

}

func TestSqliteState_InitDB_SchemaCreation(t *testing.T) {
	ctx := context.Background()

	db := newTestSqliteDB(t)
	store := state.NewSQLiteState(db)

	// initialize db and create schema
	if err := store.InitDB(ctx); err != nil {
		t.Fatalf("InitDB err: %v", err)
	}

	// scan rows
	rows, err := db.QueryContext(ctx, "PRAGMA table_info('resources')")
	if err != nil {
		t.Errorf("failed to execute query: %v", err)
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			t.Errorf("failed to close rows: %v", err)
		}
	}()

	got := make(map[string]string)
	for rows.Next() {
		var (
			cid        int
			name       string
			dataType   string
			notNull    int
			defaultVal sql.NullString
			primaryKey int
		)

		if err := rows.Scan(&cid, &name, &dataType, &notNull, &defaultVal, &primaryKey); err != nil {
			t.Errorf("failed to scan row: %v", err)
		}
		got[name] = strings.ToUpper(dataType)
	}
	if err := rows.Err(); err != nil {
		t.Errorf("rows err: %v", err)
	}

	want := map[string]string{
		"id":            "TEXT",
		"key":           "TEXT",
		"resource_type": "TEXT",
		"version":       "INTEGER",
		"data":          "BLOB",
		"created_at":    "DATETIME",
		"updated_at":    "DATETIME",
	}

	for column, wantType := range want {
		gotType, ok := want[column]
		if !ok {
			t.Errorf("column %q not found", column)
			continue
		}
		if gotType != wantType {
			t.Fatalf("column %s: got %v, want %v", column, gotType, wantType)
		}
	}
	if len(got) != len(want) {
		t.Fatalf("got %d columns, want %d", len(got), len(want))
	}
}

func TestSqliteState_InitDB_Idempotent(t *testing.T) {
	ctx := context.Background()

	db := newTestSqliteDB(t)
	store := state.NewSQLiteState(db)

	// initdb once
	if err := store.InitDB(ctx); err != nil {
		t.Errorf("first initdb failed to initdb again: %v", err)
	}

	bucket := &model.Bucket{
		BucketName: "bucket",
	}

	key := bucket.Key()

	// test
	if err := store.Save(ctx, key, bucket); err != nil {
		t.Fatalf("failed to save bucket: %v", err)
	}

	// initdb for idempotency test
	if err := store.InitDB(ctx); err != nil {
		t.Errorf("idempotent failed to initdb again: %v", err)
	}

	got, err := store.Get(ctx, key)
	if err != nil {
		t.Fatalf("failed to Get() bucket: %v", err)
	}
	if got == nil {
		t.Fatalf("bucket got nil, want non-nil")
	}

	// assert: expect key
	if got.Key() != key {
		t.Fatalf("bucket key: got %v, want %v", got.Key(), key)
	}
}

func TestSQLiteState_Migrate_V1ToV2(t *testing.T) {
	tests := []struct {
		name           string
		currentVersion int
		migrateVersion int
		setUpQuery     string
		setupArgs      []any
		wantColumns    map[string]bool
		wantVersion    int
		wantLog        string
		wantErr        error
	}{
		{
			name:           "migrate already current",
			wantLog:        state.MsgSchemaAlreadyCurrent,
			currentVersion: 2,
			migrateVersion: 2,
			wantVersion:    2,
		},
		{
			name:           "migrate unsupported version",
			currentVersion: 1,
			migrateVersion: 3,
			wantVersion:    2,
			wantErr:        state.ErrUnsupportedVersion,
		},
		{
			name:           "migrate from v1 to v2",
			currentVersion: 1,
			migrateVersion: 2,
			wantColumns: map[string]bool{
				"description": true,
			},
			wantVersion: 2,
			setUpQuery: `INSERT INTO resources (
					id,
					key,
					resource_type,
					version,
					data
				)
				VALUES (?, ?, ?, ?, ?);`,

			setupArgs: []any{
				"id-1",
				"bucket/foo",
				model.ResourceTypeBucket,
				1,
				[]byte(`{"name":"foo"}`),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			// 1. Open sql
			db := newTestSqliteDB(t)
			store := state.NewSQLiteState(db)

			// 2. create v1 schema
			schemaV1 := `CREATE TABLE IF NOT EXISTS resources (
				id TEXT PRIMARY KEY,
				key TEXT,
				resource_type TEXT NOT NULL,
				version INTEGER NOT NULL,
				data BLOB NOT NULL,
				created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
				updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
			)`

			if _, err := db.ExecContext(ctx, schemaV1); err != nil {
				t.Fatalf("failed to create schemav1: %v", err)
			}

			// 3. set current version
			if _, err := db.ExecContext(ctx, fmt.Sprintf("PRAGMA user_version = %d", tt.currentVersion)); err != nil {
				t.Fatalf("failed set v1 version: %v", err)
			}

			// 4. insert resource using v1 schema
			_, err := db.ExecContext(ctx, tt.setUpQuery, "id-1", "bucket/foo", model.ResourceTypeBucket, tt.currentVersion, []byte(`{"name":"foo"}`))
			if err != nil {
				t.Fatalf("failed to add resource for v1: %v", err)
			}

			// check: make sure the current state is the desired initial version
			var checkCurrentVersion int
			err = db.QueryRowContext(ctx, "PRAGMA user_version").Scan(&checkCurrentVersion)
			if err != nil {
				t.Fatalf("failed to execute PRAGMA user_version to the correct init version: %v", err)
			}

			if checkCurrentVersion != tt.currentVersion {
				t.Errorf("got current version: %d, want %d", checkCurrentVersion, tt.currentVersion)
			}

			// 5. test migration
			err = store.Migrate(ctx, tt.migrateVersion)
			if err != nil {
				// test: expected errors
				if tt.wantErr != nil && !strings.Contains(err.Error(), tt.wantErr.Error()) {
					t.Fatalf("migrate v1 to v2: %v", err)
				}
			}

			// assert: schema version updated
			var gotVersion int
			err = db.QueryRowContext(ctx, `PRAGMA user_version`).Scan(&gotVersion)
			if err != nil {
				t.Fatalf("failed to query PRAGMA user_version: %v", err)
			}

			if gotVersion != tt.wantVersion {
				t.Fatalf("got %d, want %d", gotVersion, tt.migrateVersion)
			}

			// query
			if tt.setupArgs != nil {
				_, err := db.ExecContext(ctx, tt.setUpQuery, tt.setupArgs...)
				if err != nil {
					t.Fatalf("failed to execute setup query: %v", err)
				}
			}

			err = db.QueryRowContext(ctx, "PRAGMA user_version ").Scan(&gotVersion)
			if err != nil {
				t.Fatalf("failed to query PRAGMA user_version: %v", err)
			}

			if gotVersion != 2 {
				t.Fatalf("got %d, want %d", gotVersion, 2)
			}
		})
	}

}

func TestSQLiteState_Save(t *testing.T) {
	tests := []struct {
		name            string
		dataSourceName  string
		dataQuery       string
		resourceKey     string
		resourceType    string
		hasData         bool
		resource        model.Resource
		setupCtx        func() (context.Context, context.CancelFunc)
		wantRecordCount int
		wantErr         error
	}{
		{
			name:            "empty db",
			dataSourceName:  ":memory:",
			resourceType:    model.ResourceTypeBucket,
			resource:        model.NewBucket("new-bucket", model.BucketSpec{}),
			wantRecordCount: 1,
		},
		{
			name:            "save with resource in db",
			dataSourceName:  ":memory:",
			hasData:         true,
			dataQuery:       `INSERT INTO resources (key, resource_type, version, data) VALUES ('bucket/existing-bucket', 'existing-bucket', 1, '{"name":"existing-bucket","key":"bucket/existing-bucket"}')`,
			resourceType:    model.ResourceTypeBucket,
			resourceKey:     "bucket/existing-bucket",
			resource:        model.NewBucket("new-bucket", model.BucketSpec{}),
			wantRecordCount: 2,
		},
		{
			name:            "duplicate resource",
			dataSourceName:  ":memory:",
			dataQuery:       `INSERT INTO resources (key, resource_type, version, data) VALUES ('bucket/existing-bucket', 'existing-bucket', 1, '{"name":"existing-bucket","key":"bucket/existing-bucket"}')`,
			resourceType:    model.ResourceTypeBucket,
			resourceKey:     "bucket/existing-bucket",
			resource:        model.NewBucket("existing-bucket", model.BucketSpec{}),
			wantRecordCount: 1,
		},
		{
			name:           "context cancelled",
			dataSourceName: ":memory:",
			dataQuery:      `INSERT INTO resources (key, resource_type, version, data) VALUES ('bucket/existing-bucket', 'existing-bucket', 1, '{"name":"existing-bucket","key":"bucket/existing-bucket"}')`,
			resourceType:   model.ResourceTypeBucket,
			resource:       model.NewBucket("existing-bucket", model.BucketSpec{}),
			resourceKey:    "bucket/existing-bucket",
			setupCtx: func() (context.Context, context.CancelFunc) {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx, cancel
			},
			wantErr:         context.Canceled,
			wantRecordCount: 0,
		},
		{
			name:     "context deadline exceeded",
			resource: model.NewBucket("existing-bucket", model.BucketSpec{}),
			setupCtx: func() (context.Context, context.CancelFunc) {
				ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(-time.Second))

				return ctx, cancel
			},
			resourceKey:     "bucket/existing-bucket",
			wantErr:         context.DeadlineExceeded,
			wantRecordCount: 0,
		},
		{
			name:     "context timeout 1 second",
			resource: model.NewBucket("existing-bucket", model.BucketSpec{}),
			setupCtx: func() (context.Context, context.CancelFunc) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				return ctx, cancel
			},
			wantRecordCount: 1,
		},
		{
			name:     "context timeout 0 second",
			resource: model.NewBucket("existing-bucket", model.BucketSpec{}),
			setupCtx: func() (context.Context, context.CancelFunc) {
				ctx, cancel := context.WithTimeout(context.Background(), 0)
				return ctx, cancel
			},
			wantErr:         context.DeadlineExceeded,
			wantRecordCount: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var count int
			cancelFunc := func() {}

			ctx := context.Background()
			driverName := "sqlite3"

			// set up context if defined in tests
			if tt.setupCtx != nil {
				ctx, cancelFunc = tt.setupCtx()
			}
			defer cancelFunc()

			db, err := sql.Open(driverName, tt.dataSourceName)
			if err != nil {
				t.Errorf("Open(): %v", err)
			}

			defer func() {
				if err := db.Close(); err != nil {
					t.Errorf("Close(): %v", err)
				}
			}()

			// New State
			db.SetMaxOpenConns(1)
			store := state.NewSQLiteState(db)

			// 1. Create DB
			err = store.InitDB(context.Background())
			if err = store.InitDB(context.Background()); err != nil {
				t.Errorf("Init() error = %v", err)
			}

			// 2. Initialize schema
			if tt.hasData {
				_, err := db.ExecContext(ctx, tt.dataQuery)
				if err != nil {
					t.Fatalf("dataQuery error = %v", err)
				}
			}

			// 3. Test context
			err = store.Save(ctx, tt.resourceKey, tt.resource)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf(
					"Save() error = %v, wantErr = %v",
					err,
					tt.wantErr,
				)
			}

			// 4.Query row count independent of test
			err = db.QueryRowContext(ctx, `SELECT COUNT(*) FROM resources`).Scan(&count)
			if err != nil && tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("QueryRowContext() error = %v, wantErr %v", err, tt.wantErr)
				}
			}

			if count != tt.wantRecordCount {
				t.Fatalf("Save() got = %v, want = %v",
					count,
					tt.wantRecordCount,
				)
			}
		})
	}
}

func testContext(setup func() (context.Context, context.CancelFunc)) (context.Context, context.CancelFunc) {
	if setup == nil {
		return context.Background(), nil
	}
	return setup()
}

func TestSQLiteState_Delete(t *testing.T) {
	tests := []struct {
		name       string
		hasData    bool
		data       string
		deleteKey  string
		expectRows int
		wantErr    error
	}{
		{
			name:    "delete existing resource",
			hasData: true,
			data: `INSERT INTO resources (id, key, resource_type, version, data)
			VALUES
				('id-1', 'bucket/foo', 'bucket', 1, '{"key":"bucket/foo"}'),
				('id-2', 'bucket/bar', 'bucket', 1, '{"key":"bucket/bar"}'),
				('id-3', 'bucket/baz', 'bucket', 1, '{"key":"bucket/baz"}');`,
			deleteKey:  "bucket/baz",
			expectRows: 2,
		},
		{
			name:    "delete non-existing resource",
			hasData: true,
			data: `INSERT INTO resources (id, key, resource_type, version, data)
			VALUES
				('id-2', 'bucket/bar', 'bucket', 1, '{"name":"bar"}'),
				('id-3', 'bucket/baz', 'bucket', 1, '{"name":"baz"}');`,
			deleteKey:  "bucket/foo",
			wantErr:    state.ErrResourceNotFound,
			expectRows: 2,
		},
		{
			name:       "delete empty",
			deleteKey:  "bucket/foo",
			wantErr:    state.ErrResourceNotFound,
			expectRows: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ctx context.Context

			ctx, cancel := context.WithTimeout(
				context.Background(),
				5*time.Second,
			)
			defer cancel()

			//setup: open
			db, _ := sql.Open("sqlite3", ":memory:")
			defer func() {
				if err := db.Close(); err != nil {
					t.Errorf("Close(): %v", err)
				}
			}()

			//new instance of sqlite3
			db.SetMaxOpenConns(1)
			store := state.NewSQLiteState(db)

			// init with context
			if err := store.InitDB(context.Background()); err != nil {
				t.Errorf("Init() error = %v", err)
			}

			// prepopulate db
			if tt.hasData {
				_, err := db.ExecContext(ctx, tt.data)
				if err != nil {
					t.Fatalf("dataQuery error = %v", err)
				}
			}
			err := store.Delete(ctx, tt.deleteKey)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Delete() error = %v, wantErr = %v", err, tt.wantErr)
			}
			_, err = store.Get(ctx, tt.deleteKey)
			if !errors.Is(err, state.ErrResourceNotFound) {
				t.Fatalf("Get() error = %v, want ErrNotFound", err)
			}
		})
	}
}

func TestSQLiteState_List(t *testing.T) {
	tests := []struct {
		name       string
		hasData    bool
		data       string
		setupCtx   func() (context.Context, context.CancelFunc)
		query      string
		expectRows int
		wantErr    error
	}{
		{
			name:    "list existing resource",
			hasData: true,
			data: `INSERT INTO resources (id, key, resource_type, version, data)
			VALUES
				('id-1', 'bucket/foo', 'bucket', 1, '{"name":"foo"}'),
				('id-2', 'bucket/bar', 'bucket', 1, '{"name":"bar"}'),
				('id-3', 'bucket/baz', 'bucket', 1, '{"name":"baz"}');`,
			query:      `DELETE FROM resources WHERE id = 'id-1';`,
			expectRows: 3,
		},
		{
			name:       "List empty resource",
			expectRows: 0,
		},
		{
			name:    "context cancelled",
			hasData: true,
			data:    `INSERT INTO resources (key, resource_type, version, data) VALUES ('bucket/existing-bucket', 'bucket', 1, '{"name":"existing-bucket","key":"bucket/existing-bucket"}')`,
			setupCtx: func() (context.Context, context.CancelFunc) {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx, cancel
			},
			wantErr: context.Canceled,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			cancelFunc := func() {}

			ctx := context.Background()

			if tt.setupCtx != nil {
				ctx, cancelFunc = tt.setupCtx()
			}
			defer cancelFunc()

			//setup: open
			db, _ := sql.Open("sqlite3", ":memory:")
			defer func() {
				if err := db.Close(); err != nil {
					t.Errorf("Close() error = %v", err)
				}
			}()

			//new instance of sqlite3
			db.SetMaxOpenConns(1)
			store := state.NewSQLiteState(db)

			// init with context
			if err = store.InitDB(context.Background()); err != nil {
				t.Errorf("Init() error = %v", err)
			}

			// prepopulate db
			if tt.hasData {
				_, err := db.ExecContext(context.Background(), tt.data)
				if err != nil {
					t.Fatalf("dataQuery error = %v", err)
				}
			}

			resources, err := store.List(ctx)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("List() error = %v, wantErr = %v", err, tt.wantErr)
			}
			if tt.wantErr == nil && len(resources) != tt.expectRows {
				t.Errorf("List() got rows= %v, want rows= %v", len(resources), tt.expectRows)
			}
		})
	}
}

func TestSQLiteState_Get(t *testing.T) {
	tests := []struct {
		name    string
		hasData bool
		data    string
		query   string
		getKey  string
		want    model.Bucket
		wantErr error
	}{
		{
			name:    "get existing resource",
			hasData: true,
			data: `INSERT INTO resources (id, key, resource_type, version, data)
			VALUES
				('id-2', 'bucket/bar', 'bucket', 1, '{"name":"bar"}'),
				('id-3', 'bucket/baz', 'bucket', 1, '{"name":"baz"}');`,
			getKey: "bucket/bar",
			want: model.Bucket{
				BucketName: "bar",
			},
		},
		{
			name:    "get non-existing resource",
			hasData: true,
			data: `INSERT INTO resources (id, key, resource_type, version, data)
			VALUES
				('id-1', 'bucket/bar', 'bucket', 1, '{"name":"bar", "key":"bucket/bar"}'),
				('id-2', 'bucket/baz', 'bucket', 1, '{"name":"baz"}');`,
			getKey:  "bucket/foo",
			wantErr: state.ErrResourceNotFound,
		},
		{
			name: "invalid unmarshal persisted bucket",
			data: `INSERT INTO resources (id, key, resource_type, version, data)
   			VALUES
			    ('id-1', 'bucket/bar', 'bucket', 1, 'invalid json');`,
			hasData: true,
			getKey:  "bucket/bar",
			wantErr: state.ErrUnmarshalResource,
		},
		{
			name: "unsupported resource version",
			data: `INSERT INTO resources (id, key, resource_type, version, data)
   			VALUES
			    ('id-1', 'bucket/bar', 'bucket', 999, 'invalid json');`,
			hasData: true,
			getKey:  "bucket/bar",
			wantErr: state.ErrUnsupportedVersion,
		},
		{
			name:    "supported resource version",
			hasData: true,
			data: `INSERT INTO resources(id, key, resource_type, version, data)
					VALUES (
						'id-1',
						'bucket/bar',
						'bucket',
						1,
						'{"name":"bar"}'
					);`,
			getKey: "bucket/bar",
			want: model.Bucket{
				BucketName: "bar",
			},
		},
		{
			name:    "get empty",
			hasData: false,
			getKey:  "bucket/foo",
			wantErr: state.ErrResourceNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			ctx := context.Background()

			//set up db
			db, err := sql.Open("sqlite3", ":memory:")
			if err != nil {
				t.Fatalf("Open() error = %v", err)
			}
			defer func() {
				if err := db.Close(); err != nil {
					t.Fatalf("Close(): %v", err)
				}
			}()

			db.SetMaxOpenConns(1)
			store := state.NewSQLiteState(db)

			// initial db with context
			if err = store.InitDB(context.Background()); err != nil {
				t.Fatalf("Init() error = %v", err)
			}

			// populate db
			if tt.hasData {
				_, err = db.ExecContext(ctx, tt.data)
				if err != nil {
					t.Fatalf("setup, populating db error = %v", err)
				}
			}

			got, err := store.Get(ctx, tt.getKey)

			if tt.wantErr != nil {
				if !strings.Contains(err.Error(), tt.wantErr.Error()) {
					t.Fatalf(
						"Get() error = %q, want error containing %q",
						err,
						tt.wantErr,
					)
				}

				if got != nil {
					t.Fatalf("Get() got = %v, want nil", got)
				}
				return
			}

			// success
			if got == nil {
				t.Fatalf("Get() got %v, want non-nil", got)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("Get() got = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func newTestSqliteDB(t *testing.T) *sql.DB {
	t.Helper()

	dbPath := filepath.Join(t.TempDir(), "sqlite.db")

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		t.Fatalf("sql.Open() error = %v", err)
	}

	db.SetMaxOpenConns(1)

	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Errorf("Close() error = %v", err)
		}
	})
	return db
}

func hasColumns(ctx context.Context, t *testing.T, db *sql.DB, columnName string) (map[string]bool, error) {
	t.Helper()
	rows, err := db.QueryContext(ctx, fmt.Sprintf("PRAGMA table_info(%s)", columnName))
	if err != nil {
		t.Errorf("table_info hasColumn() error = %v", err)
	}

	if err := rows.Close(); err != nil {
		t.Errorf("rows.Close() error = %v", err)
	}

	columns := make(map[string]bool)
	for rows.Next() {
		var (
			cid        int
			name       string
			dataType   string
			notNull    int
			defaultVal sql.NullString
			primaryKey int
		)
		if err := rows.Scan(&cid, &name, &dataType, &notNull, &defaultVal, &primaryKey); err != nil {
			t.Errorf("table_info hasColumn() error = %v", err)
		}
		columns[name] = true
	}
	return columns, state.ErrColumnDoesNotExist
}
