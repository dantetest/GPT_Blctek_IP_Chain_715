package config

import (
	"fmt"
	"os"
	"strings"
)

const defaultAPIAddress = ":8080"

type API struct {
	Environment string
	Address     string
}

func LoadAPI() (API, error) {
	cfg := API{
		Environment: valueOrDefault("APP_ENV", "local"),
		Address:     valueOrDefault("API_ADDR", defaultAPIAddress),
	}

	if err := validateEnvironment(cfg.Environment); err != nil {
		return API{}, err
	}
	if !strings.HasPrefix(cfg.Address, ":") && !strings.Contains(cfg.Address, ":") {
		return API{}, fmt.Errorf("API_ADDR must include a port: %q", cfg.Address)
	}
	if cfg.Environment == "production" {
		for _, key := range []string{"PAYMENT_PROVIDER", "KYC_PROVIDER", "EVIDENCE_PROVIDER"} {
			if strings.EqualFold(strings.TrimSpace(os.Getenv(key)), "mock") || strings.TrimSpace(os.Getenv(key)) == "" {
				return API{}, fmt.Errorf("%s must use a non-mock provider in production", key)
			}
		}
	}
	return cfg, nil
}

func valueOrDefault(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}

func validateEnvironment(value string) error {
	switch value {
	case "local", "test", "staging", "production":
		return nil
	default:
		return fmt.Errorf("APP_ENV must be local, test, staging, or production: %q", value)
	}
}
