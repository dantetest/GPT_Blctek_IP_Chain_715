package domain

import (
	"crypto/sha256"
	"errors"
	"testing"
	"time"
)

func TestDatasetDefaultVersionRequiresPublishedVersionFromSameDataset(t *testing.T) {
	now := time.Date(2026, 7, 15, 0, 0, 0, 0, time.UTC)
	dataset, err := NewDataset("dts_1", OwnerOrganization, "org_1", "Industrial Defects", "industrial-defects", "", now)
	if err != nil {
		t.Fatal(err)
	}
	version, err := NewDatasetVersion("dsv_1", dataset.ID, 1, "v1", now)
	if err != nil {
		t.Fatal(err)
	}
	if err := dataset.SetDefaultVersion(version, now); !errors.Is(err, ErrVersionNotPublished) {
		t.Fatalf("got %v", err)
	}
	readyVersion(t, &version, now)
	if err := dataset.SetDefaultVersion(version, now.Add(time.Hour)); err != nil {
		t.Fatal(err)
	}
	if dataset.DefaultVersionID != version.ID {
		t.Fatalf("default version = %q", dataset.DefaultVersionID)
	}
	other, _ := NewDatasetVersion("dsv_2", "dts_other", 1, "v1", now)
	readyVersion(t, &other, now)
	if err := dataset.SetDefaultVersion(other, now); !errors.Is(err, ErrDatasetVersionMismatch) {
		t.Fatalf("got %v", err)
	}
}

func TestPublishedVersionContentIsImmutable(t *testing.T) {
	now := time.Date(2026, 7, 15, 0, 0, 0, 0, time.UTC)
	version, _ := NewDatasetVersion("dsv_1", "dts_1", 1, "v1", now)
	readyVersion(t, &version, now)
	manifest := ManifestRef{SpecVersion: 1, Root: sha256.Sum256([]byte("other")), TotalFiles: 1, TotalSizeBytes: 5}
	if err := version.AttachManifest(manifest, now.Add(time.Hour)); !errors.Is(err, ErrVersionImmutable) {
		t.Fatalf("AttachManifest() error = %v", err)
	}
	license, _ := NewLicenseSnapshot("other license", now)
	if err := version.AttachLicense(license, now); !errors.Is(err, ErrVersionImmutable) {
		t.Fatalf("AttachLicense() error = %v", err)
	}
}

func TestSubmitReviewRequiresCompleteSnapshots(t *testing.T) {
	now := time.Date(2026, 7, 15, 0, 0, 0, 0, time.UTC)
	version, _ := NewDatasetVersion("dsv_1", "dts_1", 1, "v1", now)
	manifest := ManifestRef{SpecVersion: 1, Root: sha256.Sum256([]byte("root")), TotalFiles: 1, TotalSizeBytes: 5}
	if err := version.AttachManifest(manifest, now); err != nil {
		t.Fatal(err)
	}
	if err := version.SubmitReview(now); !errors.Is(err, ErrLicenseRequired) {
		t.Fatalf("got %v", err)
	}
	license, _ := NewLicenseSnapshot("commercial license", now)
	if err := version.AttachLicense(license, now); err != nil {
		t.Fatal(err)
	}
	if err := version.SubmitReview(now); !errors.Is(err, ErrRightsRequired) {
		t.Fatalf("got %v", err)
	}
}

func TestVersionLifecycle(t *testing.T) {
	now := time.Date(2026, 7, 15, 0, 0, 0, 0, time.UTC)
	version, _ := NewDatasetVersion("dsv_1", "dts_1", 1, "v1", now)
	manifest := ManifestRef{SpecVersion: 1, Root: sha256.Sum256([]byte("root")), TotalFiles: 2, TotalSizeBytes: 8}
	if err := version.AttachManifest(manifest, now); err != nil {
		t.Fatal(err)
	}
	license, _ := NewLicenseSnapshot("commercial license", now)
	if err := version.AttachLicense(license, now); err != nil {
		t.Fatal(err)
	}
	rights := RightsDeclaration{SourceType: "FIRST_PARTY", OwnershipBasis: "SELF_COLLECTED", CommercialUseRight: true, RedistributionRight: true, DeclaredAt: now}
	if err := version.AttachRights(rights, now); err != nil {
		t.Fatal(err)
	}
	if err := version.SetVerificationLevel(VerificationV1, now); err != nil {
		t.Fatal(err)
	}
	if err := version.SubmitReview(now); err != nil {
		t.Fatal(err)
	}
	if err := version.Approve(now); err != nil {
		t.Fatal(err)
	}
	if err := version.Publish(now); err != nil {
		t.Fatal(err)
	}
	if version.Status != VersionStatusPublished || version.PublishedAt == nil {
		t.Fatalf("unexpected version: %#v", version)
	}
	if err := version.Suspend(now); err != nil {
		t.Fatal(err)
	}
	if err := version.Resume(now); err != nil {
		t.Fatal(err)
	}
	if err := version.Takedown(now); err != nil {
		t.Fatal(err)
	}
	if err := version.Archive(now); err != nil {
		t.Fatal(err)
	}
	if version.Status != VersionStatusArchived {
		t.Fatalf("status=%s", version.Status)
	}
}

func readyVersion(t *testing.T, version *DatasetVersion, now time.Time) {
	t.Helper()
	manifest := ManifestRef{SpecVersion: 1, Root: sha256.Sum256([]byte("root")), TotalFiles: 1, TotalSizeBytes: 4}
	if err := version.AttachManifest(manifest, now); err != nil {
		t.Fatal(err)
	}
	license, err := NewLicenseSnapshot("commercial license", now)
	if err != nil {
		t.Fatal(err)
	}
	if err := version.AttachLicense(license, now); err != nil {
		t.Fatal(err)
	}
	rights := RightsDeclaration{SourceType: "FIRST_PARTY", OwnershipBasis: "SELF_COLLECTED", CommercialUseRight: true, RedistributionRight: true, DeclaredAt: now}
	if err := version.AttachRights(rights, now); err != nil {
		t.Fatal(err)
	}
	if err := version.SubmitReview(now); err != nil {
		t.Fatal(err)
	}
	if err := version.Approve(now); err != nil {
		t.Fatal(err)
	}
	if err := version.Publish(now); err != nil {
		t.Fatal(err)
	}
}

func TestApprovedVersionCannotChangeReviewedContent(t *testing.T) {
	now := time.Date(2026, 7, 15, 0, 0, 0, 0, time.UTC)
	version, _ := NewDatasetVersion("dsv_1", "dts_1", 1, "v1", now)
	manifest := ManifestRef{SpecVersion: 1, Root: sha256.Sum256([]byte("root")), TotalFiles: 1, TotalSizeBytes: 4}
	if err := version.AttachManifest(manifest, now); err != nil {
		t.Fatal(err)
	}
	license, _ := NewLicenseSnapshot("commercial license", now)
	if err := version.AttachLicense(license, now); err != nil {
		t.Fatal(err)
	}
	rights := RightsDeclaration{SourceType: "FIRST_PARTY", OwnershipBasis: "SELF_COLLECTED", CommercialUseRight: true, RedistributionRight: true, DeclaredAt: now}
	if err := version.AttachRights(rights, now); err != nil {
		t.Fatal(err)
	}
	if err := version.SubmitReview(now); err != nil {
		t.Fatal(err)
	}
	if err := version.Approve(now); err != nil {
		t.Fatal(err)
	}
	other, _ := NewLicenseSnapshot("changed after approval", now)
	if err := version.AttachLicense(other, now); !errors.Is(err, ErrInvalidTransition) {
		t.Fatalf("AttachLicense() error = %v", err)
	}
}
