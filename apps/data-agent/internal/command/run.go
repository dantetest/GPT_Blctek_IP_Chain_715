package command

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	manifestspec "github.com/dantetest/GPT_Blctek_IP_Chain_715/packages/manifest-spec"
)

var ErrUsage = errors.New("usage: blctek-agent manifest --root <directory> --output <manifest.json>")

type summary struct {
	ManifestVersion int    `json:"manifest_version"`
	RootHash        string `json:"root_hash"`
	Files           int    `json:"files"`
	TotalSizeBytes  uint64 `json:"total_size_bytes"`
	Output          string `json:"output"`
}

func Run(ctx context.Context, args []string, stdout, stderr io.Writer) error {
	if len(args) == 0 {
		return ErrUsage
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	switch args[0] {
	case "manifest":
		return runManifest(ctx, args[1:], stdout, stderr)
	case "help", "-h", "--help":
		_, _ = fmt.Fprintln(stdout, ErrUsage)
		return nil
	default:
		return fmt.Errorf("unknown command %q: %w", args[0], ErrUsage)
	}
}

func runManifest(ctx context.Context, args []string, stdout, stderr io.Writer) error {
	flags := flag.NewFlagSet("manifest", flag.ContinueOnError)
	flags.SetOutput(stderr)
	root := flags.String("root", "", "dataset root directory")
	output := flags.String("output", "", "output manifest JSON path")
	createdAtValue := flags.String("created-at", "", "RFC3339 creation time; defaults to current UTC time")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() != 0 || strings.TrimSpace(*root) == "" || strings.TrimSpace(*output) == "" {
		return ErrUsage
	}

	createdAt := time.Now().UTC()
	if strings.TrimSpace(*createdAtValue) != "" {
		parsed, err := time.Parse(time.RFC3339Nano, strings.TrimSpace(*createdAtValue))
		if err != nil {
			return fmt.Errorf("parse --created-at: %w", err)
		}
		createdAt = parsed.UTC()
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	manifest, err := manifestspec.Build(*root, createdAt)
	if err != nil {
		return fmt.Errorf("build manifest: %w", err)
	}
	encoded, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return fmt.Errorf("encode manifest: %w", err)
	}
	encoded = append(encoded, '\n')
	if err := writeAtomic(*output, encoded, 0o600); err != nil {
		return fmt.Errorf("write manifest: %w", err)
	}

	var totalSize uint64
	for _, file := range manifest.Files {
		totalSize += file.SizeBytes
	}
	result := summary{
		ManifestVersion: manifest.ManifestVersion,
		RootHash:        manifest.RootHash.String(),
		Files:           len(manifest.Files),
		TotalSizeBytes:  totalSize,
		Output:          filepath.Clean(*output),
	}
	if err := json.NewEncoder(stdout).Encode(result); err != nil {
		return fmt.Errorf("write summary: %w", err)
	}
	return nil
}

func writeAtomic(filename string, content []byte, mode os.FileMode) error {
	directory := filepath.Dir(filename)
	if err := os.MkdirAll(directory, 0o755); err != nil {
		return err
	}
	temporary, err := os.CreateTemp(directory, ".manifest-*.tmp")
	if err != nil {
		return err
	}
	temporaryName := temporary.Name()
	defer func() { _ = os.Remove(temporaryName) }()

	if err := temporary.Chmod(mode); err != nil {
		_ = temporary.Close()
		return err
	}
	if _, err := temporary.Write(content); err != nil {
		_ = temporary.Close()
		return err
	}
	if err := temporary.Sync(); err != nil {
		_ = temporary.Close()
		return err
	}
	if err := temporary.Close(); err != nil {
		return err
	}
	return os.Rename(temporaryName, filename)
}
