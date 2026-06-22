package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Auth holds authentication credentials.
type Auth struct {
	OpenLibing OpenLibingAuth `yaml:"openlibing"`
	LLM        LLMAuth        `yaml:"llm,omitempty"`
}

// OpenLibingAuth holds OpenLibing platform credentials.
type OpenLibingAuth struct {
	Token     string `yaml:"token"`
	TokenType string `yaml:"token_type"` // Bearer
}

// LLMAuth holds LLM API credentials.
type LLMAuth struct {
	Provider string `yaml:"provider"`
	APIKey   string `yaml:"api_key"`
}

// LoadAuth reads authentication from ~/.openlibing/auth.yaml.
// Returns empty Auth if file doesn't exist.
func LoadAuth() (*Auth, error) {
	auth := &Auth{
		OpenLibing: OpenLibingAuth{
			TokenType: "Bearer",
		},
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return auth, nil
	}

	path := filepath.Join(home, ".openlibing", "auth.yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return auth, nil
		}
		return nil, err
	}

	if err := yaml.Unmarshal(data, auth); err != nil {
		return nil, err
	}

	return auth, nil
}
