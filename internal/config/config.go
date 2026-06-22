package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config holds the user-level CLI configuration.
type Config struct {
	Endpoint string       `yaml:"endpoint"`
	Defaults Defaults     `yaml:"defaults"`
	Output   OutputConfig `yaml:"output"`
}

// Defaults holds default parameter values.
type Defaults struct {
	ProjectID string `yaml:"project_id"`
	Limit     int    `yaml:"limit"`
}

// OutputConfig holds output preferences.
type OutputConfig struct {
	Format string `yaml:"format"` // table | json | yaml
	Color  bool   `yaml:"color"`
	Pager  string `yaml:"pager"` // auto | never | always
}

// DefaultConfig returns the hardcoded default configuration.
func DefaultConfig() *Config {
	return &Config{
		Endpoint: "https://www.openlibing.com",
		Defaults: Defaults{
			Limit: 10,
		},
		Output: OutputConfig{
			Format: "table",
			Color:  true,
			Pager:  "auto",
		},
	}
}

// LoadConfig reads configuration from ~/.openlibing/config.yaml.
// If the file does not exist, returns DefaultConfig().
func LoadConfig() (*Config, error) {
	cfg := DefaultConfig()

	home, err := os.UserHomeDir()
	if err != nil {
		return cfg, nil // can't determine home, use defaults
	}

	path := filepath.Join(home, ".openlibing", "config.yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil // file doesn't exist, use defaults
		}
		return nil, err
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

// ConfigDir returns the user's openlibing config directory, creating it if needed.
func ConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(home, ".openlibing")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	return dir, nil
}

// SPCUserDir returns the user's custom SPC directory.
func SPCUserDir() (string, error) {
	dir, err := ConfigDir()
	if err != nil {
		return "", err
	}
	spcDir := filepath.Join(dir, "spc")
	if err := os.MkdirAll(spcDir, 0755); err != nil {
		return "", err
	}
	return spcDir, nil
}
