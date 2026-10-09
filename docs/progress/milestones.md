# MiniCloudStack Milestones

## Project Goal

Build a Go-based, CloudStack-inspired infrastructure control plane that demonstrates practical engineering in resource management, durable state, reconciliation, distributed systems, infrastructure provisioning, container execution, and observability.

**Primary objective:** 
- Build a cohesive project demonstrating various infrastructure and platform concepts.
- Serves as a lab to 

**Reference:** Use Apache CloudStack for control-plane concepts and MiniStack as the primary reference for architecture, API scope, configuration-driven provisioning, and practical integrations. This is an educational implementation, not a full CloudStack clone.

## Engineering Principles

- Design before implementation: maintain architecture diagrams, design documents, ADRs, and RFCs.
- Test behavior and failure modes, not just successful execution.
- Keep infrastructure providers replaceable behind stable interfaces.
- Prefer incremental vertical slices that deliver working functionality.
- Document tradeoffs, limitations, and operational behavior.
- Prioritize production engineering skills over feature count.


---

## Milestone 0 — Project Foundation

**Status:** In progress / foundational work started

### Deliverables

- [ ] Define the MVP scope and non-goals.
- [ ] Establish the Go module and package conventions.
- [ ] Document system context and component architecture.
- [ ] Define initial API and resource model conventions.
- [ ] Create a testing strategy covering unit, integration, and failure tests.
- [ ] Establish CI with formatting, linting, and automated tests.

### Exit criteria

- A new contributor can understand the architecture and run the tests.
- Major components and their responsibilities are documented.
- Initial technical decisions are recorded in ADRs.

## Milestone 1 — Resource State and Persistence

**Status:** In progress

### Deliverables

- [x] Define a generic state interface.
- [x] Implement an in-memory state backend.
- [x] Introduce SQLite-backed persistence.
- [x] Complete SQLite `Save`, `Get`, `List`, and `Delete`.
- [x] Implement resource serialization and deserialization.
- [x] Define resource type and schema-version handling.
- [ ] Standardize not-found, invalid-resource, and persistence errors.
- [ ] Test duplicate keys, malformed data, empty results, and database failures.
- [ ] Test context cancellation and propagation.
- [ ] Test persistence across database connections or store reinitialization.
- [ ] Document state interface and persistence decisions.

### Exit criteria

- CRUD behavior is covered by deterministic tests.
- The implementation correctly handles invalid inputs and persistence failures.
- Cancellation is tested on the operation under test, with fixture setup performed independently.
- Resource type and version metadata are validated during decoding.
- State storage can be replaced without rewriting the service layer.

### Key engineering concepts

Go interfaces, `database/sql`, context propagation, serialization, schema evolution, error wrapping, SQL constraints, test fixtures, and resource identity.

## Milestone 2 — Resource API and Lifecycle

### Deliverables

- [ ] Implement resource creation and retrieval through the service layer.
- [ ] Add resource listing and deletion.
- [ ] Validate incoming resource configurations.
- [ ] Define resource lifecycle states and permitted transitions.
- [ ] Separate API handlers, business logic, persistence, and models.
- [ ] Return consistent API errors and status codes.
- [ ] Add unit and HTTP integration tests.
- [ ] Document the resource lifecycle and API contract.

### Exit criteria

A client can create, inspect, list, and delete a resource through the API, with validation and persistence handled by distinct layers.

### Key engineering concepts

API design, separation of concerns, domain modeling, lifecycle state machines, and contract testing.

## Milestone 3 — Desired State and Reconciliation

### Deliverables

- [ ] Distinguish desired state from observed state.
- [ ] Define resource status, conditions, and lifecycle transitions.
- [ ] Implement a reconciliation loop.
- [ ] Make reconciliation idempotent.
- [ ] Handle retries, transient failures, and partial completion.
- [ ] Introduce bounded retry and backoff policies.
- [ ] Track reconciliation attempts and failure reasons.
- [ ] Test restart recovery and repeated reconciliation.
- [ ] Document reconciliation invariants and failure modes.

### Exit criteria

The control plane can repeatedly reconcile a resource without producing unintended duplicate effects. Recoverable failures can be retried, and persistent failures remain observable.

### Key engineering concepts

Control loops, idempotency, convergence, retries, state machines, eventual consistency, and recovery.

## Milestone 4 — Infrastructure Provider Abstraction

### Deliverables

- [ ] Define a provider interface for infrastructure operations.
- [ ] Implement a deterministic fake provider for tests.
- [ ] Implement an initial real provider integration.
- [ ] Map resource lifecycle operations to provider operations.
- [ ] Handle provider timeouts, errors, and partial failures.
- [ ] Track external resource identifiers.
- [ ] Define cleanup and orphan-resource behavior.
- [ ] Document provider boundaries and extension points.

### Exit criteria

A resource can be provisioned and reconciled through the provider abstraction without coupling core control-plane logic to a specific provider.

### Key engineering concepts

Ports and adapters, dependency inversion, external system integration, timeouts, and distributed failure handling.

## Milestone 5 — Compute and Container Integration

### Deliverables

- [ ] Define compute resource specifications.
- [ ] Model placement requirements and resource capacity.
- [ ] Implement a minimal scheduler or placement decision component.
- [ ] Integrate `conman` as a container execution component.
- [ ] Define container create, inspect, stop, and delete operations.
- [ ] Track desired and observed container state.
- [ ] Handle runtime failures and unavailable execution targets.
- [ ] Document the boundary between scheduling and execution.

### Exit criteria

A compute resource can move through a documented lifecycle, with placement decisions separated from runtime execution and failures reflected in observed state.

### Key engineering concepts

Scheduling, resource allocation, runtime abstraction, control-plane/data-plane separation, and lifecycle coordination.

## Milestone 6 — Distributed State and Concurrency

### Deliverables

- [ ] Document the consistency requirements of the state layer.
- [ ] Define concurrent update and resource ownership semantics.
- [ ] Add resource versions and optimistic concurrency control.
- [ ] Test conflicting updates and stale writes.
- [ ] Define leader election or work-ownership requirements where needed.
- [ ] Evaluate a distributed database appropriate to the requirements.
- [ ] Implement a separate distributed state backend if justified.
- [ ] Test recovery, network failures, and concurrent reconciliation.
- [ ] Document consistency, availability, and operational tradeoffs.

### Exit criteria

The system has explicit concurrency semantics, tested stale-write behavior, and a documented strategy for coordinating work across multiple control-plane instances.

A distributed database integration is complete only when its failure and consistency behavior has been tested—not merely when the application can connect to it.

### Key engineering concepts

Consistency models, optimistic concurrency, transactions, coordination, replication, partition tolerance, and fault recovery.

## Milestone 7 — Telemetry and Operability

### Deliverables

- [ ] Define service-level indicators for API availability and latency.
- [ ] Instrument API requests and reconciliation operations.
- [ ] Add structured logging with resource and operation identifiers.
- [ ] Add metrics for reconciliation outcomes, retries, and failures.
- [ ] Add distributed tracing across service boundaries.
- [ ] Define health, readiness, and liveness behavior.
- [ ] Build dashboards and actionable alerts.
- [ ] Document common operational failure scenarios.
- [ ] Define initial SLOs and an error-budget approach.

### Exit criteria

An operator can identify failed resources, locate the failing component, follow a request across component boundaries, and measure whether the system meets its reliability objectives.

### Key engineering concepts

OpenTelemetry, metrics, logs, traces, SLOs, alert design, and incident investigation.

## Milestone 8 — Reliability and Production Simulation

### Deliverables

- [ ] Run multiple control-plane instances.
- [ ] Inject database, provider, and runtime failures.
- [ ] Simulate process crashes and restarts.
- [ ] Test duplicate requests and repeated reconciliation.
- [ ] Test timeouts, cancellation, and dependency unavailability.
- [ ] Measure recovery time and reconciliation throughput.
- [ ] Run load tests and identify bottlenecks.
- [ ] Document incident scenarios and postmortems.
- [ ] Create a repeatable local or cloud-based test environment.

### Exit criteria

The system demonstrates documented recovery behavior under realistic failures. Performance claims are supported by repeatable measurements.

## Milestone 9 — Configuration, Intermediate Representation, and Extensibility

### Deliverables

- [ ] Define a versioned, provider-neutral intermediate representation (IR).
- [ ] Define validation and normalization rules.
- [ ] Translate user-facing configuration into the IR.
- [ ] Map supported IR resources to provider operations.
- [ ] Define resource dependencies and provisioning order.
- [ ] Document unsupported features and provider-specific differences.
- [ ] Add end-to-end configuration-to-provisioning tests.

### Exit criteria

A configuration can be parsed, validated, represented independently of a specific provider, and executed through the control plane with clear dependency and error semantics.

### Key engineering concepts

Intermediate representations, compiler pipelines, dependency graphs, declarative configuration, and provider portability.

## Milestone 10 — Integrated Demonstration and Technical Portfolio

### Deliverables

- [ ] Publish the end-to-end architecture diagram.
- [ ] Document key ADRs and rejected alternatives.
- [ ] Document API contracts and lifecycle state machines.
- [ ] Publish a failure-mode analysis.
- [ ] Create a repeatable end-to-end demo.
- [ ] Demonstrate provisioning, reconciliation, telemetry, and recovery.
- [ ] Document measured performance and known limitations.
- [ ] Write a technical walkthrough explaining key engineering decisions.

### Exit criteria

A reviewer can run the system, reproduce the demo, inspect the tests and design documents, and understand the reasoning behind its major architectural decisions.

---

## Cross-Cutting Requirements

These requirements apply throughout every milestone.

- **Testing:** Unit tests, integration tests, failure tests, and end-to-end tests as appropriate.
- **Documentation:** Keep architecture diagrams, ADRs, API contracts, and operational notes current.
- **Reliability:** Define expected behavior for cancellation, timeouts, retries, and partial failure.
- **Security:** Validate inputs, avoid leaking secrets, and define authentication and authorization before exposing the API.
- **Observability:** Add instrumentation as components emerge rather than postponing all telemetry until the end.
- **Maintainability:** Keep package responsibilities and interfaces explicit.
- **Evidence:** Support performance and reliability claims with reproducible tests or measurements.

## Scope Guardrails

The goal is not to reimplement all of Apache CloudStack or Kubernetes.

Prioritize a coherent, demonstrable control plane with:

1. Durable resource state.
2. Correct resource lifecycle semantics.
3. Idempotent reconciliation.
4. Provider abstraction and real infrastructure integration.
5. Failure recovery and concurrency handling.
6. Useful telemetry and operational evidence.

Defer broad UI development, numerous provider integrations, elaborate scheduling, and advanced distributed coordination until the core system requires them.

## Progress Tracking

Update this file whenever a milestone changes:

- Use `[ ]` for incomplete work and `[x]` for completed work.
- Mark a milestone complete only when its exit criteria are satisfied.
- Record significant scope changes in an ADR.
- Link relevant tests, design documents, and implementation notes.
- Prefer demonstrable working behavior over lines of code or feature count.