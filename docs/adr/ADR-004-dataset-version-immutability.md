# ADR-004: Dataset Version Immutability

- Status: Accepted
- Date: 2026-07-15

## Context

Orders, licenses, quality reports, evidence records, and delivery packages must refer to an exact set of data. Updating files behind an existing marketplace record would invalidate prior purchases and evidence.

## Decision

`Dataset` is the mutable product identity. `DatasetVersion` is the content identity. A version may be edited only before approval. After publication, content fields are locked permanently; suspension, takedown, and archival change availability but never content.

The following values are version-owned snapshots:

- Manifest root and statistics;
- license text and hash;
- rights declaration;
- verification level;
- future quality, preview, evidence, and delivery package references.

Any content change creates a new monotonically numbered version. Persistence updates use optimistic revision checks.

## Consequences

Historical orders remain reproducible. Storage use increases because snapshots are retained. Administrative corrections use new records or explicit status events rather than overwriting history.
