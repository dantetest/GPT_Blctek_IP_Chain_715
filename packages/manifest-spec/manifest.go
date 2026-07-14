package manifestspec

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"time"
)

var ErrInvalidManifest = errors.New("manifest is invalid")

type Manifest struct {
	ManifestVersion int         `json:"manifest_version"`
	HashAlgorithm   string      `json:"hash_algorithm"`
	ChunkSizeBytes  int         `json:"chunk_size_bytes"`
	PathEncoding    string      `json:"path_encoding"`
	PathSeparator   string      `json:"path_separator"`
	UnicodeForm     string      `json:"unicode_normalization"`
	SymlinkPolicy   string      `json:"symlink_policy"`
	CreatedAt       time.Time   `json:"created_at"`
	Files           []FileEntry `json:"files"`
	RootHash        Digest      `json:"root_hash"`
}

func Build(root string, createdAt time.Time) (Manifest, error) {
	rootInfo, err := os.Stat(root)
	if err != nil {
		return Manifest{}, err
	}
	if !rootInfo.IsDir() {
		return Manifest{}, fmt.Errorf("manifest root must be a directory")
	}

	files := make([]FileEntry, 0)
	seen := make(map[string]struct{})
	err = filepath.WalkDir(root, func(filename string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if filename == root {
			return nil
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return ErrSymbolicLink
		}
		if entry.IsDir() {
			return nil
		}
		if !info.Mode().IsRegular() {
			return fmt.Errorf("unsupported non-regular file: %s", filename)
		}
		canonicalPath, err := RelativeCanonicalPath(root, filename)
		if err != nil {
			return err
		}
		if _, exists := seen[canonicalPath]; exists {
			return ErrDuplicatePath
		}
		seen[canonicalPath] = struct{}{}
		file, err := os.Open(filename)
		if err != nil {
			return err
		}
		hashed, hashErr := HashReader(file, canonicalPath, DefaultChunkSize)
		closeErr := file.Close()
		if hashErr != nil {
			return hashErr
		}
		if closeErr != nil {
			return closeErr
		}
		files = append(files, hashed)
		return nil
	})
	if err != nil {
		return Manifest{}, err
	}

	sort.Slice(files, func(left, right int) bool { return files[left].Path < files[right].Path })
	manifest := Manifest{
		ManifestVersion: SpecVersion,
		HashAlgorithm:   HashAlgorithm,
		ChunkSizeBytes:  DefaultChunkSize,
		PathEncoding:    "UTF-8",
		PathSeparator:   "/",
		UnicodeForm:     "NFC",
		SymlinkPolicy:   "REJECT",
		CreatedAt:       createdAt.UTC(),
		Files:           files,
		RootHash:        MerkleRoot(files),
	}
	if err := manifest.Validate(); err != nil {
		return Manifest{}, err
	}
	return manifest, nil
}

func (m Manifest) Validate() error {
	if m.ManifestVersion != SpecVersion || m.HashAlgorithm != HashAlgorithm || m.ChunkSizeBytes != DefaultChunkSize || m.RootHash.IsZero() {
		return ErrInvalidManifest
	}
	seen := make(map[string]struct{}, len(m.Files))
	for index, entry := range m.Files {
		canonical, err := CanonicalPath(entry.Path)
		if err != nil || canonical != entry.Path || entry.FileHash.IsZero() {
			return ErrInvalidManifest
		}
		if index > 0 && m.Files[index-1].Path >= entry.Path {
			return ErrInvalidManifest
		}
		if _, exists := seen[entry.Path]; exists {
			return ErrDuplicatePath
		}
		seen[entry.Path] = struct{}{}
		if computeFileHash(entry.Path, entry.SizeBytes, entry.ChunkHashes) != entry.FileHash {
			return ErrInvalidManifest
		}
	}
	if MerkleRoot(m.Files) != m.RootHash {
		return ErrInvalidManifest
	}
	return nil
}
