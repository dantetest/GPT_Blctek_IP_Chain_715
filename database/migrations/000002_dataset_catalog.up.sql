CREATE TABLE datasets (
    id VARCHAR(64) NOT NULL,
    owner_type VARCHAR(32) NOT NULL,
    owner_id VARCHAR(64) NOT NULL,
    title VARCHAR(255) NOT NULL,
    slug VARCHAR(160) NOT NULL,
    description TEXT NULL,
    status VARCHAR(32) NOT NULL DEFAULT 'ACTIVE',
    default_version_id VARCHAR(64) NULL,
    revision BIGINT UNSIGNED NOT NULL DEFAULT 1,
    created_at DATETIME(6) NOT NULL,
    updated_at DATETIME(6) NOT NULL,
    deleted_at DATETIME(6) NULL,
    PRIMARY KEY (id),
    UNIQUE KEY uq_datasets_owner_slug (owner_type, owner_id, slug),
    KEY idx_datasets_owner (owner_type, owner_id, created_at),
    KEY idx_datasets_status (status, created_at),
    CONSTRAINT chk_datasets_owner_type CHECK (owner_type IN ('USER','ORGANIZATION')),
    CONSTRAINT chk_datasets_status CHECK (status IN ('ACTIVE','SUSPENDED','ARCHIVED'))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE dataset_versions (
    id VARCHAR(64) NOT NULL,
    dataset_id VARCHAR(64) NOT NULL,
    version_number INT UNSIGNED NOT NULL,
    version_label VARCHAR(128) NOT NULL,
    status VARCHAR(32) NOT NULL DEFAULT 'DRAFT',
    manifest_spec_version SMALLINT UNSIGNED NULL,
    manifest_root BINARY(32) NULL,
    manifest_file_count BIGINT UNSIGNED NULL,
    manifest_total_size_bytes BIGINT UNSIGNED NULL,
    verification_level VARCHAR(4) NOT NULL DEFAULT 'V0',
    content_locked_at DATETIME(6) NULL,
    published_at DATETIME(6) NULL,
    revision BIGINT UNSIGNED NOT NULL DEFAULT 1,
    created_at DATETIME(6) NOT NULL,
    updated_at DATETIME(6) NOT NULL,
    PRIMARY KEY (id),
    UNIQUE KEY uq_dataset_versions_number (dataset_id, version_number),
    KEY idx_dataset_versions_status (status, created_at),
    KEY idx_dataset_versions_manifest_root (manifest_root),
    CONSTRAINT fk_dataset_versions_dataset
        FOREIGN KEY (dataset_id) REFERENCES datasets(id)
        ON DELETE RESTRICT ON UPDATE RESTRICT,
    CONSTRAINT chk_dataset_versions_number CHECK (version_number > 0),
    CONSTRAINT chk_dataset_versions_status CHECK (
        status IN ('DRAFT','SCANNING','MANIFEST_READY','REVIEWING','REJECTED','APPROVED','PUBLISHED','SUSPENDED','TAKEDOWN','ARCHIVED')
    ),
    CONSTRAINT chk_dataset_versions_verification CHECK (verification_level IN ('V0','V1','V2','V3','V4')),
    CONSTRAINT chk_dataset_versions_manifest CHECK (
        (manifest_root IS NULL AND manifest_spec_version IS NULL AND manifest_file_count IS NULL AND manifest_total_size_bytes IS NULL)
        OR
        (manifest_root IS NOT NULL AND manifest_spec_version = 1 AND manifest_file_count IS NOT NULL AND manifest_total_size_bytes IS NOT NULL)
    )
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

ALTER TABLE datasets
    ADD CONSTRAINT fk_datasets_default_version
    FOREIGN KEY (default_version_id) REFERENCES dataset_versions(id)
    ON DELETE RESTRICT ON UPDATE RESTRICT;

CREATE TABLE dataset_license_snapshots (
    id VARCHAR(64) NOT NULL,
    dataset_version_id VARCHAR(64) NOT NULL,
    license_text MEDIUMTEXT NOT NULL,
    license_hash BINARY(32) NOT NULL,
    created_at DATETIME(6) NOT NULL,
    PRIMARY KEY (id),
    UNIQUE KEY uq_dataset_license_version (dataset_version_id),
    KEY idx_dataset_license_hash (license_hash),
    CONSTRAINT fk_dataset_license_version
        FOREIGN KEY (dataset_version_id) REFERENCES dataset_versions(id)
        ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE dataset_rights_declarations (
    id VARCHAR(64) NOT NULL,
    dataset_version_id VARCHAR(64) NOT NULL,
    source_type VARCHAR(64) NOT NULL,
    ownership_basis VARCHAR(255) NOT NULL,
    commercial_use_right BOOLEAN NOT NULL,
    redistribution_right BOOLEAN NOT NULL,
    contains_personal_data BOOLEAN NOT NULL DEFAULT FALSE,
    contains_sensitive_data BOOLEAN NOT NULL DEFAULT FALSE,
    contains_biometric_data BOOLEAN NOT NULL DEFAULT FALSE,
    contains_minors_data BOOLEAN NOT NULL DEFAULT FALSE,
    contains_third_party_content BOOLEAN NOT NULL DEFAULT FALSE,
    risk_notes TEXT NULL,
    declared_by VARCHAR(64) NOT NULL,
    declared_at DATETIME(6) NOT NULL,
    created_at DATETIME(6) NOT NULL,
    PRIMARY KEY (id),
    UNIQUE KEY uq_dataset_rights_version (dataset_version_id),
    CONSTRAINT fk_dataset_rights_version
        FOREIGN KEY (dataset_version_id) REFERENCES dataset_versions(id)
        ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE dataset_version_events (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    dataset_version_id VARCHAR(64) NOT NULL,
    from_status VARCHAR(32) NULL,
    to_status VARCHAR(32) NOT NULL,
    actor_type VARCHAR(32) NOT NULL,
    actor_id VARCHAR(64) NOT NULL,
    reason VARCHAR(1024) NULL,
    metadata JSON NULL,
    created_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    PRIMARY KEY (id),
    KEY idx_dataset_version_events_version (dataset_version_id, id),
    CONSTRAINT fk_dataset_version_events_version
        FOREIGN KEY (dataset_version_id) REFERENCES dataset_versions(id)
        ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
