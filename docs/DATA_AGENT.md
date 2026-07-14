# Data Agent Design Baseline

## Current command

```text
blctek-agent manifest --root <directory> --output <manifest.json>
```

The command:

1. validates arguments;
2. scans the directory with Manifest v1 rules;
3. hashes files as streams in 4 MiB chunks;
4. validates the resulting Manifest;
5. writes formatted JSON atomically;
6. emits a one-line JSON summary to stdout.

## Atomic output

The Agent writes into the destination directory using a temporary file, calls `fsync`, closes it, and renames it over the target. Errors remove the temporary file. The target is never truncated before a complete Manifest is available.

## Exit contract

- exit `0`: Manifest and summary were written;
- non-zero: no successful result should be consumed;
- operational errors are written to stderr;
- stdout is reserved for machine-readable success output.

## Next implementation slice

- SQLite scan job and checkpoint schema;
- incremental file metadata cache;
- pause and resume;
- file mutation detection;
- Agent device key and API registration;
- signed report upload.
