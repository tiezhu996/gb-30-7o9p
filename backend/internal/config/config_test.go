package config

import "testing"

func TestLoadDefaults(t *testing.T) {
	cfg := Load()
	if cfg.ServerPort != "8080" {
		t.Errorf("ServerPort = %s", cfg.ServerPort)
	}
	if cfg.MinIOBucket != "gbadopt" {
		t.Errorf("MinIOBucket = %s", cfg.MinIOBucket)
	}
	if cfg.RateLimitReq <= 0 {
		t.Errorf("RateLimitReq invalid: %d", cfg.RateLimitReq)
	}
	if cfg.DSN() == "" {
		t.Error("DSN should not be empty")
	}
}
