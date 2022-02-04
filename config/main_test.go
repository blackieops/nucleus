package config

import (
	"testing"
)

func TestLoadConfigValid(t *testing.T) {
	c, err := LoadConfig("internal/fixtures/valid.yaml")
	if err != nil {
		t.Errorf("Failed to load valid config: %v", err)
	}
	// Not meant to be exhaustive, just make sure some spot-checked fields have
	// the right values.
	if c.Port != 8888 {
		t.Errorf("Parsed Port incorrectly: %v", c.Port)
	}
	if c.BaseURL != "https://localhost:443" {
		t.Errorf("Parsed BaseURL incorrectly: %v", c.BaseURL)
	}
	if c.DataPath != "/var/db/nucleus_data" {
		t.Errorf("Parsed BaseURL incorrectly: %v", c.DataPath)
	}
}

func TestLoadConfigInvalid(t *testing.T) {
	_, err := LoadConfig("internal/fixtures/notexist.yaml")
	if err == nil {
		t.Errorf("Somehow loaded non-existent config!")
	}
	_, err = LoadConfig("internal/fixtures/invalid.yaml")
	if err == nil {
		t.Errorf("Erroneously loaded invalid config!")
	}
	_, err = LoadConfig("internal/fixtures/useless.yaml")
	if err != nil {
		t.Errorf("Should have loaded useless but syntactically valid config!")
	}
}
