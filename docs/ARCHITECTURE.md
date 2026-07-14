# Architecture Baseline

## Runtime components

- `apps/api`: modular Go HTTP API;
- `apps/worker`: asynchronous Outbox and scheduled-job processor;
- `apps/web`: Next.js user, seller, enterprise, and admin surfaces;
- future `apps/data-agent`: local scanning, hashing, quality, encryption, and seeding;
- future `apps/private-tracker`: authenticated private BitTorrent tracker.

## Reliability boundary

Business transactions write domain state and an Outbox event in the same database transaction. External Providers are invoked asynchronously. Provider callbacks are persisted to an Inbox table before domain processing.

## Initial dependency policy

The first bootstrap keeps the Go runtime on the standard library to ensure the repository can compile before external dependency lockfiles are introduced. HTTP framework adoption will be isolated behind transport packages and reviewed in a later ADR without changing domain packages.
