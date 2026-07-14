ALTER TABLE datasets DROP FOREIGN KEY fk_datasets_default_version;
DROP TABLE IF EXISTS dataset_version_events;
DROP TABLE IF EXISTS dataset_rights_declarations;
DROP TABLE IF EXISTS dataset_license_snapshots;
DROP TABLE IF EXISTS dataset_versions;
DROP TABLE IF EXISTS datasets;
