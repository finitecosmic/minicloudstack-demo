# ADR-0002: Use SQLite for Persistent State

## Status
- Accepted

## Date
2026-10-05

## Context
- The initial implementation stored resources in-mem, state is 
  lost whenever the process restart.  The system needed a persistent
  state independent of the service layer. Sqlite can provide that at this current state,
  Postgresql can be implemented later on due to the abstraction set up through the store interface.

```
                    State Interface
                         │
          ┌──────────────┼──────────────┐
          │              │              │
     MemoryState     SQLiteState   PostgresState
          │              │              │
       testing        development     production
```

## Decision

The application will interact with a generic `State` interface:
```
    Save()
    Get()
    List()
    Delete()
    Ready()

```

## Alternatives
- Postgresql
- JSON
- in-memory

```
                         Persistence
                              │
        ┌─────────────────────┼─────────────────────┐
        │                     │                     │
    In-memory              SQLite             PostgreSQL
        │                     │                     │
    ephemeral             single node          multi-instance
                              │                     │
                         current MVP          production scale
```


## Positives
- Simple local development, no infrastructure setup.
- light-weight
- One of the plugins and can be swapped out for another
- embedded
- low operational overhead
- still provides the standard SQL interface `database/sql`

## Drawbacks
- No horizontal scaling
- Replication is limited
- HA model is not matured
- single concurrent user

## Related Artifacts
- 0001-state-interface.md

## Implementation Status
- Implemented 