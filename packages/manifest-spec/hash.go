package manifestspec

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"hash"
	"io"
)

const (
	SpecVersion      = 1
	HashAlgorithm    = "SHA-256"
	DefaultChunkSize = 4 * 1024 * 1024
)

var emptyRoot = sha256.Sum256([]byte("BLCTEK_EMPTY_V1\x00"))

type FileEntry struct {
	Path        string   `json:"path"`
	SizeBytes   uint64   `json:"size_bytes"`
	ChunkHashes []Digest `json:"chunk_hashes"`
	FileHash    Digest   `json:"file_hash"`
}

func HashReader(reader io.Reader, canonicalPath string, chunkSize int) (FileEntry, error) {
	canonicalPath, err := CanonicalPath(canonicalPath)
	if err != nil {
		return FileEntry{}, err
	}
	if chunkSize <= 0 {
		return FileEntry{}, fmt.Errorf("chunk size must be positive")
	}

	buffer := make([]byte, chunkSize)
	entry := FileEntry{Path: canonicalPath}
	for {
		count, readErr := io.ReadFull(reader, buffer)
		if readErr != nil && readErr != io.EOF && readErr != io.ErrUnexpectedEOF {
			return FileEntry{}, readErr
		}
		if count > 0 {
			entry.SizeBytes += uint64(count)
			entry.ChunkHashes = append(entry.ChunkHashes, sha256.Sum256(buffer[:count]))
		}
		if readErr == io.EOF || readErr == io.ErrUnexpectedEOF {
			break
		}
	}
	entry.FileHash = computeFileHash(entry.Path, entry.SizeBytes, entry.ChunkHashes)
	return entry, nil
}

func computeFileHash(canonicalPath string, size uint64, chunks []Digest) Digest {
	h := sha256.New()
	writeDomain(h, "BLCTEK_FILE_V1")
	writeString(h, canonicalPath)
	writeUint64(h, size)
	writeUint32(h, uint32(len(chunks)))
	for _, chunk := range chunks {
		_, _ = h.Write(chunk[:])
	}
	return digestFromHash(h)
}

func leafHash(entry FileEntry) Digest {
	h := sha256.New()
	writeDomain(h, "BLCTEK_LEAF_V1")
	writeString(h, entry.Path)
	_, _ = h.Write(entry.FileHash[:])
	return digestFromHash(h)
}

func nodeHash(left, right Digest) Digest {
	h := sha256.New()
	writeDomain(h, "BLCTEK_NODE_V1")
	_, _ = h.Write(left[:])
	_, _ = h.Write(right[:])
	return digestFromHash(h)
}

func MerkleRoot(files []FileEntry) Digest {
	if len(files) == 0 {
		return emptyRoot
	}
	level := make([]Digest, len(files))
	for index, file := range files {
		level[index] = leafHash(file)
	}
	for len(level) > 1 {
		next := make([]Digest, 0, (len(level)+1)/2)
		for index := 0; index < len(level); index += 2 {
			right := level[index]
			if index+1 < len(level) {
				right = level[index+1]
			}
			next = append(next, nodeHash(level[index], right))
		}
		level = next
	}
	return level[0]
}

func writeDomain(h hash.Hash, value string) {
	_, _ = h.Write([]byte(value))
	_, _ = h.Write([]byte{0})
}

func writeString(h hash.Hash, value string) {
	writeUint32(h, uint32(len(value)))
	_, _ = h.Write([]byte(value))
}

func writeUint32(h hash.Hash, value uint32) {
	var encoded [4]byte
	binary.BigEndian.PutUint32(encoded[:], value)
	_, _ = h.Write(encoded[:])
}

func writeUint64(h hash.Hash, value uint64) {
	var encoded [8]byte
	binary.BigEndian.PutUint64(encoded[:], value)
	_, _ = h.Write(encoded[:])
}

func digestFromHash(h hash.Hash) Digest {
	var digest Digest
	copy(digest[:], h.Sum(nil))
	return digest
}
