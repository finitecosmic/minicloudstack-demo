package state

const (
	createResourceTable = `
	CREATE TABLE IF NOT EXISTS resources (
		id TEXT PRIMARY KEY,
		key VARCHAR(255),
		resource_type TEXT NOT NULL,
		version INTEGER NOT NULL,
		data BLOB NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
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
)
