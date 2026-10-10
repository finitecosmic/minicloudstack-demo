package state

const (
	schemaVersionV1      = 1
	schemaVersionV2      = 2
	currentSchemaVersion = schemaVersionV2

	MsgSchemaAlreadyCurrent = "database schema is already current"

	// Queries
	createResourceTable = `
	CREATE TABLE IF NOT EXISTS resources (
		id TEXT PRIMARY KEY,
		key TEXT,
		resource_type TEXT NOT NULL,
		version INTEGER NOT NULL,
		data BLOB NOT NULL,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	)`

	createResourceVersionTableV2 = `
	CREATE TABLE IF NOT EXISTS resources (
		id TEXT PRIMARY KEY,
		key TEXT,
		description TEXT,
		resource_type TEXT NOT NULL,
		version INTEGER NOT NULL,
		data BLOB NOT NULL,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	)`

	insertResource = `
	INSERT INTO resources (
		id,
	    key,
	    resource_type,
		version,
		data
	)
	VALUES (?, ?, ?, ?, ?);`

	listResources = `SELECT key, resource_type, version, data from resources`

	migrateV1toV2 = `ALTER TABLE resources ADD COLUMN TEXT NOT NULL DEFAULT ''`
)
