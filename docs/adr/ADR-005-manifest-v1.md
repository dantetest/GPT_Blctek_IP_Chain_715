# ADR-005: Manifest v1 Canonical Hashing

- Status: Accepted
- Date: 2026-07-15

## Context

The same directory must produce the same content identity on Windows, macOS, and Linux. Plain concatenation is ambiguous, filesystem paths vary, and large files require streaming.

## Decision

Manifest v1 uses:

- UTF-8 paths normalized to Unicode NFC;
- `/` as the only separator;
- rejection of absolute paths, traversal, empty segments, and symbolic links;
- bytewise ascending canonical-path order;
- SHA-256;
- 4 MiB chunks;
- length-prefixed and big-endian encoded hash inputs;
- domain-separated file, leaf, node, and empty-root hashes;
- duplicate-last Merkle balancing for odd node counts.

Static JSON test vectors are part of the repository and run on Linux, Windows, and macOS CI.

## Consequences

Implementations in other languages must match the byte-level specification exactly. Any future change requires a new manifest version rather than modifying v1.
