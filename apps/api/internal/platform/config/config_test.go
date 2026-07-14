package config

import "testing"

func TestLoadAPIDefaults(t *testing.T) {
	t.Setenv("APP_ENV", "")
	t.Setenv("API_ADDR", "")
	cfg, err := LoadAPI()
	if err != nil {
		t.Fatalf("LoadAPI() error = %v", err)
	}
	if cfg.Environment != "local" || cfg.Address != ":8080" {
		t.Fatalf("unexpected defaults: %#v", cfg)
	}
}

func TestLoadAPIRejectsMockProvidersInProduction(t *testing.T) {
	t.Setenv("APP_ENV", "production")
	t.Setenv("PAYMENT_PROVIDER", "mock")
	t.Setenv("KYC_PROVIDER", "real")
	t.Setenv("EVIDENCE_PROVIDER", "real")
	if _, err := LoadAPI(); err == nil {
		t.Fatal("LoadAPI() expected an error")
	}
}
