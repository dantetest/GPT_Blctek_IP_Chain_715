package mysqlrepo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/dantetest/GPT_Blctek_IP_Chain_715/apps/api/internal/dataset/domain"
)

type execer interface {
	ExecContext(context.Context, string, ...any) (sql.Result, error)
}

func persistCompliance(ctx context.Context, tx *sql.Tx, version domain.DatasetVersion, actorID string) error {
	if version.License != nil {
		if _, err := tx.ExecContext(ctx, `INSERT INTO dataset_license_snapshots
			(id, dataset_version_id, license_text, license_hash, created_at)
			VALUES (?, ?, ?, ?, ?)
			ON DUPLICATE KEY UPDATE license_text = VALUES(license_text), license_hash = VALUES(license_hash), created_at = VALUES(created_at)`,
			"lic_"+version.ID, version.ID, version.License.Text, version.License.Hash[:], version.License.CreatedAt,
		); err != nil {
			return fmt.Errorf("persist license snapshot: %w", err)
		}
	}

	if version.Rights != nil {
		if actorID == "" {
			actorID = "system"
		}
		rights := version.Rights
		if _, err := tx.ExecContext(ctx, `INSERT INTO dataset_rights_declarations
			(id, dataset_version_id, source_type, ownership_basis, commercial_use_right, redistribution_right,
			 contains_personal_data, contains_sensitive_data, contains_biometric_data, contains_minors_data,
			 contains_third_party_content, risk_notes, declared_by, declared_at, created_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
			ON DUPLICATE KEY UPDATE source_type = VALUES(source_type), ownership_basis = VALUES(ownership_basis),
			 commercial_use_right = VALUES(commercial_use_right), redistribution_right = VALUES(redistribution_right),
			 contains_personal_data = VALUES(contains_personal_data), contains_sensitive_data = VALUES(contains_sensitive_data),
			 contains_biometric_data = VALUES(contains_biometric_data), contains_minors_data = VALUES(contains_minors_data),
			 contains_third_party_content = VALUES(contains_third_party_content), risk_notes = VALUES(risk_notes),
			 declared_by = VALUES(declared_by), declared_at = VALUES(declared_at)`,
			"rgt_"+version.ID, version.ID, rights.SourceType, rights.OwnershipBasis,
			rights.CommercialUseRight, rights.RedistributionRight, rights.ContainsPersonalData,
			rights.ContainsSensitiveData, rights.ContainsBiometricData, rights.ContainsMinorsData,
			rights.ContainsThirdPartyContent, rights.RiskNotes, actorID, rights.DeclaredAt, version.UpdatedAt,
		); err != nil {
			return fmt.Errorf("persist rights declaration: %w", err)
		}
	}
	return nil
}

func (r *Repository) loadCompliance(ctx context.Context, version *domain.DatasetVersion) error {
	var license domain.LicenseSnapshot
	var licenseHash []byte
	err := r.db.QueryRowContext(ctx, `SELECT license_text, license_hash, created_at
		FROM dataset_license_snapshots WHERE dataset_version_id = ?`, version.ID).
		Scan(&license.Text, &licenseHash, &license.CreatedAt)
	switch {
	case err == nil:
		if len(licenseHash) != len(license.Hash) {
			return fmt.Errorf("invalid persisted license hash for version %s", version.ID)
		}
		copy(license.Hash[:], licenseHash)
		version.License = &license
	case errors.Is(err, sql.ErrNoRows):
	default:
		return fmt.Errorf("load license snapshot: %w", err)
	}

	var rights domain.RightsDeclaration
	err = r.db.QueryRowContext(ctx, `SELECT source_type, ownership_basis, commercial_use_right, redistribution_right,
		contains_personal_data, contains_sensitive_data, contains_biometric_data, contains_minors_data,
		contains_third_party_content, COALESCE(risk_notes, ''), declared_at
		FROM dataset_rights_declarations WHERE dataset_version_id = ?`, version.ID).
		Scan(&rights.SourceType, &rights.OwnershipBasis, &rights.CommercialUseRight, &rights.RedistributionRight,
			&rights.ContainsPersonalData, &rights.ContainsSensitiveData, &rights.ContainsBiometricData,
			&rights.ContainsMinorsData, &rights.ContainsThirdPartyContent, &rights.RiskNotes, &rights.DeclaredAt)
	switch {
	case err == nil:
		version.Rights = &rights
	case errors.Is(err, sql.ErrNoRows):
	default:
		return fmt.Errorf("load rights declaration: %w", err)
	}
	return nil
}
