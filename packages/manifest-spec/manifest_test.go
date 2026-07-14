package manifestspec

import (
	"encoding/base64"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

type vectorFile struct {
	Path          string `json:"path"`
	ContentBase64 string `json:"content_base64"`
	FileHash      string `json:"file_hash"`
}

type vector struct {
	Name     string       `json:"name"`
	RootHash string       `json:"root_hash"`
	Files    []vectorFile `json:"files"`
}

func TestManifestVectors(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("testdata", "vectors.json"))
	if err != nil {
		t.Fatal(err)
	}
	var vectors []vector
	if err := json.Unmarshal(data, &vectors); err != nil {
		t.Fatal(err)
	}
	for _, item := range vectors {
		t.Run(item.Name, func(t *testing.T) {
			root := t.TempDir()
			for _, file := range item.Files {
				content, err := base64.StdEncoding.DecodeString(file.ContentBase64)
				if err != nil {
					t.Fatal(err)
				}
				filename := filepath.Join(root, filepath.FromSlash(file.Path))
				if err := os.MkdirAll(filepath.Dir(filename), 0o755); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(filename, content, 0o600); err != nil {
					t.Fatal(err)
				}
			}
			manifest, err := Build(root, time.Unix(0, 0))
			if err != nil {
				t.Fatal(err)
			}
			if manifest.RootHash.String() != item.RootHash {
				t.Fatalf("root=%s want=%s", manifest.RootHash, item.RootHash)
			}
			if len(manifest.Files) != len(item.Files) {
				t.Fatalf("files=%d", len(manifest.Files))
			}
			for index, file := range manifest.Files {
				if file.Path != item.Files[index].Path || file.FileHash.String() != item.Files[index].FileHash {
					t.Fatalf("file[%d]=%s/%s", index, file.Path, file.FileHash)
				}
			}
			if err := manifest.Validate(); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestManifestDetectsMutation(t *testing.T) {
	root := t.TempDir()
	filename := filepath.Join(root, "file.txt")
	if err := os.WriteFile(filename, []byte("original"), 0o600); err != nil {
		t.Fatal(err)
	}
	first, err := Build(root, time.Unix(0, 0))
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filename, []byte("mutated"), 0o600); err != nil {
		t.Fatal(err)
	}
	second, err := Build(root, time.Unix(0, 0))
	if err != nil {
		t.Fatal(err)
	}
	if first.RootHash == second.RootHash {
		t.Fatal("root hash did not change")
	}
}
