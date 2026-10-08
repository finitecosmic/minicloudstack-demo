# ADR-0004: Resource Model Abstraction

- Status: Accepted
- Date: 2026-10-06

## Context
This project will manage multiple resources, such as networks, vm, buckets, and other
infrastructure resources. As a result, the control plan needs a way to manage resources without coupling higher-level services to
individual resource implementation. The object storage service should be able to persist 
a bucket without needing to know whether the state is backed by memory, SQLite, or PostgreSQL.

## Rationale
The state layer should be able to operate on different resources through a common abstraction.
Without a common interface `Resource`, each state would have to implement duplicate logic for each resource.
This is better know as the `Strategy` design pattern.

```
             model.Resource 
                  │ 
   ┌──────────────┼──────────────┐ 
   ▼              ▼              ▼ 
 Bucket        Network          VM 
   │              │              │              
   └──────────────┼──────────────┘ 
                  │ 
                  ▼ 
            State interface 
                  │ 
       ┌──────────┼──────────┐ 
       ▼          ▼          ▼ 
     Memory    SQLite    PostgreSQL
```

## Implementation

The Resource Interface
```
type Resource interface {
	Key() string
	Name() string
	Dependencies() []string
	Type() string
	Spec() Spec
	SetUpdatedAt(time.Time)
}
```
The Bucket logic adhering to `Resource`. In `Bucket`,as long as the following functions are 
implemented, it will be recognized a `Resource` type at the higher level, but it will also be
recognized as a `Bucket` type

```
func (b Bucket) Key() string {}
func (b Bucket) Name() string {}
func (b Bucket) Dependencies() {}
func (b Bucket) Type() string {}
func (b Bucket) Spec() Spec {}
func (b Bucket) SetUpdatedAt(updatedAt time.Time) {}
```

The state takes in any type `Resource` to persist in state.
For instance, in `state/memory.go`, the `Save()` will accept any resource type. As more types
get developed, it can be passed into memory for save.

```
func (m *Memory) Save(ctx context.Context, key string, resource model.Resource) (model.Resource, error) {
```
## Alternatives

### 1. Separate interfaces for every resource type

```
type BucketStore interface  { 
    SaveBucket(...) 
    GetBucket(...) 
}
```
#### Rejected
- higher-level components will be tightly-coupled
- adding new resources will require additional storage interface

### 2. Store untyped data
- representing data as raw json

#### Rejected
- Removes compile-time type safety.
- Must interpret un-typed data

### 3. Resource-based persistent implementation
- Give each resource its own database implementation.

#### Rejected
- Couples resource models directly to storage.
- Duplicate common persistence behavior

## Consequences
### Positive
- Reduce duplicated code; each state will not have to implement different logic for every
  individual resources.
- scalable

### Negative
- Abstraction introduces complexities
- Resources have to follow strict compliance to `Resource` interface. Some resource cannot
  be expressed through generic abstraction.
- Increased likelihood of breaking changes
- Changes to the Resource interface can affect every resource implementation.

## Future
```
              Resource 
                 │ 
    ┌────────────┼────────────┐ 
    ▼            ▼            ▼ 
Desired        State      Observed 
Config         Store       State 
    │            │            │ 
    └────────────┼────────────┘ 
                 ▼ 
             Reconciler 
                 │ 
                 ▼ 
              Provider 
                 │ 
       ┌─────────┼─────────┐ 
       ▼         ▼         ▼ 
      AWS       GCP       Azure
```

## Related Decisions
- 0001-state-interface.md
- 0002-sqlite-persistent-state.md
- 0003-resource-model-abstraction.md