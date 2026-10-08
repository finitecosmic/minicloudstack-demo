# ADR-0001: State Interface

- Status: Accepted
- Date: 2026-10-05 

## Context
Introduces a generic `State` interface to abstract resource persistence. 

The interface provides operations for:

```
Save — persist a resource
Get — retrieve a resource by key
List — retrieve resources
Delete — remove a resource
Ready — determine whether the state backend is available
```

```
type State interface {
    Save(ctx context.Context, key string, resource model.Resource) error
    Get(ctx context.Context, key string) (model.Resource, error)
    List(ctx context.Context) (map[string]model.Resource, error)
    Delete(ctx context.Context, key string) error
    Ready(ctx context.Context) error
}

```
## Rationale 
This allows for expanding state stores. Consistent contract for implementing and testing
different storage systems. For instance, `MemoryState`, `SqliteState`, and `PostgresState`
can be developed simultaneously without depending on one another.

## Positives
- Independent development
- Pluggable
- storage is decoupled from service
- scalable. Other state/stores can be added later
- Services can be tested independent of state


## Negatives
- Abstraction and complexity
- Some backend-specific capabilities may not fit cleanly into the generic interface.
- Additional implementations require conformance testing to ensure that each implementation satisfies 
  the behavioral contract defined by the `State` interface.
```
               State Contract
                  Tests
                    │
          ┌─────────┼─────────┐
          ▼         ▼         ▼
      MemoryState SQLiteState PostgresState
          │         │         │
       passes?   passes?   passes?
```

## Status of Implementation
- Implemented