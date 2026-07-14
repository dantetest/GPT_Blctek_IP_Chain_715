package manifestspec

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

type Digest [sha256.Size]byte

func ParseDigest(value string) (Digest, error) {
	var digest Digest
	decoded, err := hex.DecodeString(value)
	if err != nil || len(decoded) != sha256.Size {
		return digest, fmt.Errorf("invalid SHA-256 digest")
	}
	copy(digest[:], decoded)
	return digest, nil
}

func (d Digest) String() string { return hex.EncodeToString(d[:]) }
func (d Digest) IsZero() bool   { return d == Digest{} }

func (d Digest) MarshalJSON() ([]byte, error) { return json.Marshal(d.String()) }

func (d *Digest) UnmarshalJSON(data []byte) error {
	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	parsed, err := ParseDigest(value)
	if err != nil {
		return err
	}
	*d = parsed
	return nil
}
