package domain

import "errors"

var (
	ErrInvalidDataset         = errors.New("dataset is invalid")
	ErrInvalidVersion         = errors.New("dataset version is invalid")
	ErrInvalidTransition      = errors.New("dataset version transition is invalid")
	ErrVersionImmutable       = errors.New("published dataset version content is immutable")
	ErrManifestRequired       = errors.New("manifest is required")
	ErrLicenseRequired        = errors.New("license snapshot is required")
	ErrRightsRequired         = errors.New("rights declaration is required")
	ErrVersionNotPublished    = errors.New("dataset version is not published")
	ErrDatasetVersionMismatch = errors.New("dataset version belongs to another dataset")
)
