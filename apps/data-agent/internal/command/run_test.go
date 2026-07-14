package command

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestRunManifestWritesValidatedOutput(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "nested"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "nested", "data.txt"), []byte("payload"), 0o600); err != nil {
		t.Fatal(err)
	}
	output := filepath.Join(t.TempDir(), "manifest.json")
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	err := Run(context.Background(), []string{
		"manifest",
		"--root", root,
		"--output", output,
		"--created-at", "2026-07-15T00:00:00Z",
	}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("Run() error = %v, stderr = %s", err, stderr.String())
	}
	data, err := os.ReadFile(output)
	if err != nil {
		t.Fatal(err)
	}
	var manifest map[string]any
	if err := json.Unmarshal(data, &manifest); err != nil {
		t.Fatal(err)
	}
	if manifest["manifest_version"] != float64(1) || manifest["root_hash"] == "" {
		t.Fatalf("unexpected manifest: %#v", manifest)
	}
	var got summary
	if err := json.Unmarshal(stdout.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if got.Files != 1 || got.TotalSizeBytes != 7 || got.Output != output {
		t.Fatalf("unexpected summary: %#v", got)
	}
}

func TestRunRejectsMissingArguments(t *testing.T) {
	err := Run(context.Background(), []string{"manifest"}, &bytes.Buffer{}, &bytes.Buffer{})
	if !errors.Is(err, ErrUsage) {
		t.Fatalf("Run() error = %v", err)
	}
}

func TestWriteAtomicReplacesExistingFile(t *testing.T) {
	filename := filepath.Join(t.TempDir(), "manifest.json")
	if err := os.WriteFile(filename, []byte("old"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := writeAtomic(filename, []byte("new"), 0o600); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "new" {
		t.Fatalf("content = %q", data)
	}
}
