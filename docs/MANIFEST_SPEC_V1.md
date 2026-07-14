# BlctekIP Manifest Specification v1

## 1. Constants

```text
manifest_version = 1
hash_algorithm = SHA-256
chunk_size_bytes = 4,194,304
path_encoding = UTF-8
unicode_normalization = NFC
path_separator = /
symlink_policy = REJECT
file_order = canonical UTF-8 byte order ascending
```

## 2. Canonical paths

A path is valid only when it is relative and contains no:

- absolute Unix prefix;
- Windows drive prefix;
- UNC prefix;
- NUL;
- empty segment;
- `.` segment;
- `..` segment;
- invalid UTF-8.

Backslashes are converted to `/`, then the complete string is normalized to NFC. Two source paths that normalize to the same canonical path make the manifest invalid.

## 3. Integer encoding

All integers used inside hash inputs are unsigned and encoded in network byte order:

```text
uint32 = 4-byte big-endian
uint64 = 8-byte big-endian
```

A string is encoded as:

```text
uint32(byte_length) || UTF-8 bytes
```

A domain separator is encoded as:

```text
ASCII domain bytes || 0x00
```

## 4. Chunk hash

```text
chunk_hash = SHA256(chunk_bytes)
```

An empty file has zero chunks.

## 5. File hash

```text
file_hash = SHA256(
    domain("BLCTEK_FILE_V1") ||
    string(canonical_path) ||
    uint64(file_size_bytes) ||
    uint32(chunk_count) ||
    chunk_hash[0] || ... || chunk_hash[n-1]
)
```

## 6. Leaf hash

```text
leaf_hash = SHA256(
    domain("BLCTEK_LEAF_V1") ||
    string(canonical_path) ||
    file_hash
)
```

## 7. Merkle tree

Files are sorted by canonical path. Parent nodes are:

```text
node_hash = SHA256(
    domain("BLCTEK_NODE_V1") ||
    left_hash ||
    right_hash
)
```

If a level has an odd number of nodes, the final node is duplicated as both left and right. An empty dataset uses:

```text
SHA256("BLCTEK_EMPTY_V1" || 0x00)
```

## 8. Compatibility

Published Manifest v1 bytes and roots are immutable. A change to path policy, encoding, chunk size, hash algorithm, or Merkle balancing requires Manifest v2.
