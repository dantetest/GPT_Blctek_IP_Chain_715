package config

import "testing"

func TestLoadAPIDefaults(t *testing.T) {
	t.Setenv("APP_ENV", "")
	t.Setenv("API_ADDR", "")
	t.Setenv("DATASET_REPOSITORY", "")
	t.Setenv("DATABASE_DSN", "")
	cfg, err := LoadAPI()
	if err != nil {
		t.Fatalf("LoadAPI() error = %v", err)
	}
	if cfg.Environment != "local" || cfg.Address != ":8080" || cfg.DatasetRepository != DatasetRepositoryMemory {
		t.Fatalf("unexpected defaults: %#v", cfg)
	}
}

func TestLoadAPIRequiresDSNForMySQL(t *testing.T) {
	t.Setenv("APP_ENV", "staging")
	t.Setenv("DATASET_REPOSITORY", DatasetRepositoryMySQL)
	t.Setenv("DATABASE_DSN", "")
	if _, err := LoadAPI(); err == nil {
		t.Fatal("LoadAPI() expected a missing DSN error")
	}
}

func TestLoadAPIRejectsMemoryRepositoryInProduction(t *testing.T) {
	t.Setenv("APP_ENV", "production")
	t.Setenv("DATASET_REPOSITORY", DatasetRepositoryMemory)
	t.Setenv("PAYMENT_PROVIDER", "real")
	t.Setenv("KYC_PROVIDER", "real")
	t.Setenv("EVIDENCE_PROVIDER", "real")
	if _, err := LoadAPI(); err == nil {
		t.Fatal("LoadAPI() expected a production repository error")
	}
}

func TestLoadAPIRejectsMockProvidersInProduction(t *testing.T) {
	t.Setenv("APP_ENV", "production")
	t.Setenv("DATASET_REPOSITORY", DatasetRepositoryMySQL)
	t.Setenv("DATABASE_DSN", "user:pass@tcp(mysql:3306)/blctekip")
	t.Setenv("PAYMENT_PROVIDER", "mock")
	t.Setenv("KYC_PROVIDER", "real")
	t.Setenv("EVIDENCE_PROVIDER", "real")
	if _, err := LoadAPI(); err == nil {
		t.Fatal("LoadAPI() expected an error")
	}
}
