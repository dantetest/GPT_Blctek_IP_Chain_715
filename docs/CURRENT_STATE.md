# Current State

## Repository baseline

The repository started empty. There is no legacy production code or migration compatibility burden.

## Implemented

### Platform foundation

- deterministic Go workspace with API and Worker modules;
- health and readiness endpoints;
- production Mock Provider rejection;
- Request ID propagation;
- idempotency-key validation;
- Outbox retry model;
- reliability database tables;
- minimal Next.js shell;
- Docker Compose infrastructure and CI.

### Dataset foundation

- mutable Dataset product identity;
- immutable-after-publication Dataset Version aggregate;
- explicit review and publication state transitions;
- license and rights snapshots;
- optimistic revision fields;
- Dataset catalog migration and append-only status event table.

### Manifest v1

- streaming SHA-256 file hashing;
- 4 MiB chunks;
- UTF-8 NFC canonical paths;
- symbolic-link and traversal rejection;
- domain-separated file, leaf, node, and empty hashes;
- deterministic Merkle tree;
- static vectors;
- Linux, Windows, and macOS CI matrix.

## Explicitly not implemented yet

- persistent Dataset repository and HTTP APIs;
- authentication, organizations, and KYC;
- Data Agent CLI and scan checkpoints;
- quality reports and previews;
- orders and payment Providers;
- private P2P delivery;
- settlement, Ledger, refunds, disputes, and evidence Providers.

## Next slice

Implement the persistent Dataset repository and Dataset APIs, then create the first Data Agent CLI around Manifest v1.
