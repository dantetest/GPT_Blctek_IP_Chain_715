# Current State

## Repository baseline

The repository was empty before this iteration. No legacy production code or database migration exists, so there is no backward-compatibility burden yet.

## Implemented in this iteration

- Deterministic Go workspace with API and Worker modules;
- API health and readiness endpoints;
- strict environment validation, including production Mock Provider rejection;
- Request ID propagation and input sanitization;
- structured response envelope and problem type;
- idempotency-key parser;
- Outbox event and retry policy model;
- MySQL tables for idempotency, Outbox, Provider Callback Inbox, and admin audit logs;
- minimal Next.js landing page;
- Docker Compose infrastructure;
- CI baseline and ADRs.

## Explicitly not implemented yet

- authentication, organizations, KYC;
- Dataset and Dataset Version;
- Manifest and Data Agent;
- orders and payment providers;
- private P2P delivery;
- settlement, ledger, refunds, disputes, and evidence providers.

## Legacy disposition

There is no legacy code in this repository. Future changes must preserve executed migrations and use explicit state transitions.
