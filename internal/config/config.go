package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const ConfigFileName = ".vibe-skills.yaml"

type Config struct {
	Skills []string `yaml:"skills"`
}

func Load(dir string) (*Config, error) {
	path := filepath.Join(dir, ConfigFileName)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func Save(dir string, cfg *Config) error {
	path := filepath.Join(dir, ConfigFileName)
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

func Exists(dir string) bool {
	path := filepath.Join(dir, ConfigFileName)
	_, err := os.Stat(path)
	return err == nil
}

func GetDefaultConfig() *Config {
	return &Config{
		Skills: []string{
			"common/commit-convention",
			"common/code-reviewer",
			"common/pull-request",
		},
	}
}
