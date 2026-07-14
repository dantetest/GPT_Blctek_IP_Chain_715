# Manifest v1 Go implementation

This module is the canonical Go implementation of `docs/MANIFEST_SPEC_V1.md`.

```go
manifest, err := manifestspec.Build("/data/dataset", time.Now())
```

Compatibility is guarded by static vectors in `testdata/vectors.json`. Do not change the byte-level hashing rules under Manifest v1; introduce a new specification version instead.
