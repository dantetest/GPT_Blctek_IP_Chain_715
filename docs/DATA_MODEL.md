# Data Model Baseline

## Dataset aggregate

```text
Dataset
  ├── mutable marketplace metadata
  ├── owner: user or organization
  ├── default published version
  └── DatasetVersion[]
```

## Dataset Version aggregate

```text
DatasetVersion
  ├── version number and label
  ├── Manifest reference
  ├── License Snapshot
  ├── Rights Declaration
  ├── Verification Level
  └── append-only status events
```

## Lifecycle

```text
DRAFT
  → SCANNING
  → MANIFEST_READY
  → REVIEWING
  → APPROVED
  → PUBLISHED
  → SUSPENDED / TAKEDOWN
  → ARCHIVED
```

`REVIEWING` may transition to `REJECTED`; a rejected version may be corrected and resubmitted. Content edits are rejected after approval and permanently locked after publication.

## Persistence rules

- IDs are externally generated prefixed identifiers;
- money is not present in this aggregate;
- Manifest roots and license hashes use `BINARY(32)`;
- version numbers are unique within a Dataset;
- revisions support optimistic concurrency;
- deletion is restricted for version-owned snapshots;
- status changes are appended to `dataset_version_events`.
