package state_test

import (
	"context"
	"database/sql"
	"errors"
	"minicloudstack/internal/model"
	"minicloudstack/internal/state"
	"reflect"
	"strings"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func TestSQLiteState_Init(t *testing.T) {
	t.Parallel()
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
		t.Errorf("InitDB err: %v", err)
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
			resource:        model.NewBucket("new-bucket", model.BucketSpec{}),
			wantRecordCount: 2,
		},
		{
			name:            "duplicate resource",
			dataSourceName:  ":memory:",
			dataQuery:       `INSERT INTO resources (key, resource_type, version, data) VALUES ('bucket/existing-bucket', 'existing-bucket', 1, '{"name":"existing-bucket","key":"bucket/existing-bucket"}')`,
			resourceType:    model.ResourceTypeBucket,
			resource:        model.NewBucket("existing-bucket", model.BucketSpec{}),
			wantRecordCount: 1,
		},
		{
			name:           "context cancelled",
			dataSourceName: ":memory:",
			dataQuery:      `INSERT INTO resources (key, resource_type, version, data) VALUES ('bucket/existing-bucket', 'existing-bucket', 1, '{"name":"existing-bucket","key":"bucket/existing-bucket"}')`,
			resourceType:   model.ResourceTypeBucket,
			resource:       model.NewBucket("existing-bucket", model.BucketSpec{}),
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

			// open the db
			db, err := sql.Open(driverName, tt.dataSourceName)
			if err != nil {
				t.Errorf("Open(): %v", err)
			}

			// defer close db and handle error
			defer func() {
				if err := db.Close(); err != nil {
					t.Errorf("Close(): %v", err)
				}
			}()

			// New State
			db.SetMaxOpenConns(1)
			store := state.NewSQLiteState(db)

			// Initialize schema
			err = store.InitDB(context.Background())
			if err = store.InitDB(context.Background()); err != nil {
				t.Errorf("Init() error = %v", err)
			}

			// Populate data
			if tt.hasData {
				_, err := db.ExecContext(ctx, tt.dataQuery)
				if err != nil {
					t.Fatalf("dataQuery error = %v", err)
				}
			}

			// Test Case
			err = store.Save(ctx, tt.resourceKey, tt.resource)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf(
					"Save() error = %v, wantErr = %v",
					err,
					tt.wantErr,
				)
			}

			// Expected number of row entries
			db.QueryRowContext(ctx, `SELECT COUNT(*) FROM resources`).Scan(&count)
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
		query      string
		expectRows int
		wantErr    error
	}{
		{
			name:    "delete existing resource",
			hasData: true,
			data: `INSERT INTO resources (id, key, resource_type, version, data)
			VALUES
				('id-1', 'bucket/foo', 'bucket', 1, '{"name":"foo"}'),
				('id-2', 'bucket/bar', 'bucket', 1, '{"name":"bar"}'),
				('id-3', 'bucket/baz', 'bucket', 1, '{"name":"baz"}');`,
			query:      `DELETE FROM resources WHERE id = 'id-1';`,
			deleteKey:  "id-1",
			wantErr:    state.ErrResourceNotFound,
			expectRows: 2,
		},
		{
			name:    "delete non-existing resource",
			hasData: true,
			data: `INSERT INTO resources (id, key, resource_type, version, data)
			VALUES
				('id-2', 'bucket/bar', 'bucket', 1, '{"name":"bar"}'),
				('id-3', 'bucket/baz', 'bucket', 1, '{"name":"baz"}');`,
			query:      `DELETE FROM resources WHERE id = 'id-1';`,
			deleteKey:  "id-1",
			wantErr:    state.ErrResourceNotFound,
			expectRows: 2,
		},
		{
			name:       "delete empty",
			query:      `DELETE FROM resources WHERE id = 'id-1';`,
			deleteKey:  "id-1",
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
			err := store.Delete(ctx, tt.query)
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
			name: "unknown resource version",
			data: `INSERT INTO resources (id, key, resource_type, version, data)
   			VALUES
			    ('id-1', 'bucket/bar', 'bucket', 1, 'invalid json');`,
			hasData: true,
			getKey:  "bucket/bar",
			wantErr: state.ErrUnmarshalResource,
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

			// expected an error
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Get() error = %v, want ErrResourceNotFound", err)
			}

			if tt.wantErr != nil {
				if got != nil {
					t.Fatalf("Get() got = %v, want nil", got)
				}
				if !strings.Contains(err.Error(), tt.wantErr.Error()) {
					t.Fatalf(
						"Get() error = %q, want error containing %q",
						err,
						tt.wantErr,
					)
				}
				return
			}

			t.Logf("Get() got = %v", got.Key())
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("Get() got = %#v, want %#v", got, tt.want)
			}
		})
	}
}
