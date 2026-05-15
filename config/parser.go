package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

const Delimiter = "L/DD77"

type Dependency struct {
	URL string `json:"url"`
}

type TaskConfig struct {
	Command string   `json:"command"`
	Inputs  []string `json:"inputs,omitempty"`
	Outputs []string `json:"outputs,omitempty"`
}

type ProjectConfig struct {
	Name         string                `json:"name"`
	Version      string                `json:"version"`
	Dependencies map[string]string     `json:"dependencies"`
	Tasks        map[string]TaskConfig `json:"tasks"`
}

type DependencyLock struct {
	Hash      string `json:"hash"`
	UpdatedAt string `json:"updated_at"`
	Size      int64  `json:"size"`
}

type TaskLock struct {
	InputHash string `json:"input_hash"`
	Success   bool   `json:"success"`
}

type LockData struct {
	Dependencies map[string]DependencyLock `json:"dependencies"`
	Tasks        map[string]TaskLock       `json:"tasks"`
}

type CrowFile struct {
	Config ProjectConfig
	Locks  LockData
	Path   string
}

// Load reads and parses the .crw file
func Load(path string) (*CrowFile, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	parts := strings.Split(string(content), Delimiter)
	if len(parts) < 1 {
		return nil, fmt.Errorf("invalid .crw format: missing configuration")
	}

	var cf CrowFile
	cf.Path = path

	// Parse Part 1: Config
	if err := json.Unmarshal([]byte(parts[0]), &cf.Config); err != nil {
		return nil, fmt.Errorf("error parsing config: %v", err)
	}

	// Parse Part 2: Locks (if exists)
	cf.Locks.Dependencies = make(map[string]DependencyLock)
	cf.Locks.Tasks = make(map[string]TaskLock)
	if len(parts) > 1 && strings.TrimSpace(parts[1]) != "" {
		if err := json.Unmarshal([]byte(parts[1]), &cf.Locks); err != nil {
			return nil, fmt.Errorf("error parsing locks: %v", err)
		}
	}

	return &cf, nil
}

// Save writes the project back to the .crw file
func (cf *CrowFile) Save() error {
	configJSON, err := json.MarshalIndent(cf.Config, "", "  ")
	if err != nil {
		return err
	}

	locksJSON, err := json.MarshalIndent(cf.Locks, "", "  ")
	if err != nil {
		return err
	}

	content := string(configJSON) + "\n" + Delimiter + "\n" + string(locksJSON)
	return os.WriteFile(cf.Path, []byte(content), 0644)
}
