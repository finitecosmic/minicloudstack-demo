CREATE TABLE IF NOT EXISTS resources (
     id TEXT PRIMARY KEY,
     key TEXT NOT NULL UNIQUE,
     resource_type TEXT NOT NULL,
     version INTEGER NOT NULL DEFAULT 1 CHECK (version > 0),
     data BLOB NOT NULL,
     created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
     updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);


INSERT INTO resources (
    id,
    key,
    resource_type,
    version,
    data
)
VALUES (
           '550e8400-e29b-41d4-a716-446655440000',
           'c',
           'bucket',
           1,
           '{"name":"existing-bucket","key":"bucket/existing-bucket"}'
       );

INSERT INTO resources (
    id,
    key,
    resource_type,
    data
)
VALUES (
           'random',
           'bucket/existing-bucket',
           'bucket',
1,
'{"name":"existing-bucket","key":"bucket/existing-bucket"}'
       );
