# ADR-0003: Resource Versioning

- Status: Accepted
- Date: 2026-10-06

## Context
Resources managed may evolve over time.For example, a Bucket resource may initially contain only a name, and later 
gain additional configuration such as encryption, versioning, or tags.  Without versioning, 
the system cannot reliably determine which schema to use in the database.

- Database schema versioning tracks changes to the database structure.
- Resource versioning tracks changes to the serialized representation of an individual resource.

The state layer need to determine the version of stored resources data and select the
appropriate decoder.

## Decision
Persistent Resources are versioned independently of the database schema.

Each persisted resource will contain a version in the `resources.version` column.
For example:
```
const BucketVersion = "1.0.0"
```
The persistence layer will store the resource's current version when saving it.

When loading a resource, the decoder will use both the resource type and the version:

```
resource_type + version
          │ 
          ▼
     decodeResource()
          │ 
          ▼
     decodeBucket()
          │ 
          ▼
   version-specific decoder
```
Each version-specific decoder will serialize data according to its version.

```
func decodeBucket(data []byte, version int) (model.Resource, error) {
    switch version { 
    case model.BucketVersion: 
        // Decode current bucket representation. 
    default: 
        return nil, fmt.Errorf("unknown bucket version %d", version) } }
```

As resource schemas evolve, additional versions may be supported:

Bucket v1
Bucket v2
Bucket v3

## Alternative
### 1. No Version Control

#### Reject:
- Schema changes will be difficult to decode
- Future migration becomes harder to implement safely

### 2. Database Schema Versioning Only
Use only database versioning to present resource versioning

#### Rejected:
Database migrations describe the structure of the database itself. 
They do not necessarily describe the serialized schema of individual 
resources stored in the data column.

### 3. Store Version inside the serialized resource
The goal is to have portable resources or self-contained serialization.
With a database, a column will be decided to `version`

## Consequences
### Positive
- Resource versioning and database versioning can evolve independent of one another.

### Negative
- Old versions will have to be supported
- Each version will require separate decoding logic
- Resource migration adds to testing complexities.
- Version management becomes part of the resource lifecycle

## Implementation Notes
```
const BucketVersion=1
```
resource table will store
```
resource_type = "bucket"
version = 1
data = serialized data
```
Future version may introduce migration functions

```
       Bucket v1
          │ 
          ▼
     Migration v1 -> v2
          │ 
          ▼
       Bucket v2

```
Resource versioning should remain independent from database migration tooling.

## Related Decisions
- ADR-0001: State Interface
- ADR-0002: SQLite Persistent State