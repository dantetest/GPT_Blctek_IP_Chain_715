package domain

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

type OwnerType string

const (
	OwnerUser         OwnerType = "USER"
	OwnerOrganization OwnerType = "ORGANIZATION"
)

type DatasetStatus string

const (
	DatasetStatusActive    DatasetStatus = "ACTIVE"
	DatasetStatusSuspended DatasetStatus = "SUSPENDED"
	DatasetStatusArchived  DatasetStatus = "ARCHIVED"
)

type VersionStatus string

const (
	VersionStatusDraft         VersionStatus = "DRAFT"
	VersionStatusScanning      VersionStatus = "SCANNING"
	VersionStatusManifestReady VersionStatus = "MANIFEST_READY"
	VersionStatusReviewing     VersionStatus = "REVIEWING"
	VersionStatusRejected      VersionStatus = "REJECTED"
	VersionStatusApproved      VersionStatus = "APPROVED"
	VersionStatusPublished     VersionStatus = "PUBLISHED"
	VersionStatusSuspended     VersionStatus = "SUSPENDED"
	VersionStatusTakedown      VersionStatus = "TAKEDOWN"
	VersionStatusArchived      VersionStatus = "ARCHIVED"
)

type VerificationLevel string

const (
	VerificationV0 VerificationLevel = "V0"
	VerificationV1 VerificationLevel = "V1"
	VerificationV2 VerificationLevel = "V2"
	VerificationV3 VerificationLevel = "V3"
	VerificationV4 VerificationLevel = "V4"
)

type Digest [sha256.Size]byte

func ParseDigest(value string) (Digest, error) {
	var result Digest
	decoded, err := hex.DecodeString(strings.TrimSpace(value))
	if err != nil || len(decoded) != sha256.Size {
		return result, fmt.Errorf("parse digest: %w", ErrInvalidVersion)
	}
	copy(result[:], decoded)
	return result, nil
}

func (d Digest) String() string { return hex.EncodeToString(d[:]) }
func (d Digest) IsZero() bool   { return d == Digest{} }

type ManifestRef struct {
	SpecVersion    int
	Root           Digest
	TotalFiles     uint64
	TotalSizeBytes uint64
}

func (m ManifestRef) Valid() bool {
	return m.SpecVersion == 1 && !m.Root.IsZero()
}

type LicenseSnapshot struct {
	Text      string
	Hash      Digest
	CreatedAt time.Time
}

func NewLicenseSnapshot(text string, now time.Time) (LicenseSnapshot, error) {
	text = strings.TrimSpace(text)
	if text == "" {
		return LicenseSnapshot{}, ErrLicenseRequired
	}
	return LicenseSnapshot{
		Text:      text,
		Hash:      sha256.Sum256([]byte(text)),
		CreatedAt: now.UTC(),
	}, nil
}

type RightsDeclaration struct {
	SourceType                string
	OwnershipBasis            string
	CommercialUseRight        bool
	RedistributionRight       bool
	ContainsPersonalData      bool
	ContainsSensitiveData     bool
	ContainsBiometricData     bool
	ContainsMinorsData        bool
	ContainsThirdPartyContent bool
	RiskNotes                 string
	DeclaredAt                time.Time
}

func (r RightsDeclaration) Valid() bool {
	return strings.TrimSpace(r.SourceType) != "" &&
		strings.TrimSpace(r.OwnershipBasis) != "" &&
		!r.DeclaredAt.IsZero()
}
