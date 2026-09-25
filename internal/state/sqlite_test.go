package state_test

import (
	"context"
	"database/sql"
	"errors"
	"minicloudstack/internal/model"
	"minicloudstack/internal/state"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func TestSQLiteState_Init(t *testing.T) {
	ctx := context.Background()
	db, err := sql.Open("sqlite3", ":memory:")

	if err != nil {
		t.Fatal(err)
	}
	s := state.NewSQLiteState(db)
	err = s.InitDB(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	db.SetMaxOpenConns(1)
	var tableName string
	err = db.QueryRowContext(ctx, `
		SELECT name
		FROM sqlite_master 
		where type='table' and 
	    name='resources'
		`).Scan(&tableName)

	defer db.Close()
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
		setupCtx        func() context.Context
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
			setupCtx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			},
			wantErr:         context.Canceled,
			wantRecordCount: 0,
		},
		{
			name:     "context deadline exceeded",
			resource: model.NewBucket("existing-bucket", model.BucketSpec{}),
			setupCtx: func() context.Context {
				ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(-time.Second))
				cancel()
				return ctx
			},
			wantErr: context.DeadlineExceeded,
		},
		{
			name:     "context timeout 1 second",
			resource: model.NewBucket("existing-bucket", model.BucketSpec{}),
			setupCtx: func() context.Context {
				ctx, _ := context.WithTimeout(context.Background(), time.Second)
				return ctx
			},
			wantRecordCount: 1,
		},
		{
			name:     "context timeout 0 second",
			resource: model.NewBucket("existing-bucket", model.BucketSpec{}),
			setupCtx: func() context.Context {
				ctx, _ := context.WithTimeout(context.Background(), 0)
				return ctx
			},
			wantErr:         context.DeadlineExceeded,
			wantRecordCount: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ctx context.Context
			var count int

			initContext := context.Background()
			driverName := "sqlite3"
			ctx = testContext(tt.setupCtx)

			db, err := sql.Open(driverName, tt.dataSourceName)
			if err != nil {
				t.Fatalf("Open(): %v", err)
			}
			defer db.Close()

			//setup
			db.SetMaxOpenConns(1)
			store := state.NewSQLiteState(db)

			// Initialize schema
			err = store.InitDB(initContext)
			if err = store.InitDB(initContext); err != nil {
				t.Fatalf("Init() error = %v", err)
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

			// Test errors
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

func testContext(setup func() context.Context) context.Context {
	if setup == nil {
		return context.Background()
	}
	return setup()
}
