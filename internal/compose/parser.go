package compose

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// ServiceDefinition holds the declared configuration for a single service.
type ServiceDefinition struct {
	Image       string            `yaml:"image"`
	Environment map[string]string `yaml:"environment"`
	Ports       []string          `yaml:"ports"`
	Volumes     []string          `yaml:"volumes"`
	Command     string            `yaml:"command"`
}

// ComposeFile represents a parsed docker-compose.yml file.
type ComposeFile struct {
	Version  string                       `yaml:"version"`
	Services map[string]ServiceDefinition `yaml:"services"`
}

// ParseFile reads and parses a docker-compose YAML file at the given path.
func ParseFile(path string) (*ComposeFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading compose file %q: %w", path, err)
	}

	var cf ComposeFile
	if err := yaml.Unmarshal(data, &cf); err != nil {
		return nil, fmt.Errorf("parsing compose file %q: %w", path, err)
	}

	if len(cf.Services) == 0 {
		return nil, fmt.Errorf("compose file %q defines no services", path)
	}

	return &cf, nil
}
