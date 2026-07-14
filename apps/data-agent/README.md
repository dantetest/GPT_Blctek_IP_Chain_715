# BlctekIP Data Agent

The first Agent slice provides a deterministic Manifest CLI.

```bash
go run ./apps/data-agent/cmd/agent manifest \
  --root /path/to/dataset \
  --output /path/to/manifest.json
```

For reproducible tests, pass an explicit RFC3339 timestamp:

```bash
--created-at 2026-07-15T00:00:00Z
```

The output file is written to a temporary file, synchronized, and atomically renamed. A failed build never leaves a partially written target Manifest.

Future slices add SQLite checkpoints, device identity, quality reports, encryption, private Torrent generation, and seeding.
