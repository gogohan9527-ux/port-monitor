package main

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"
)

type Config struct {
	ListenAddr    string  `json:"listenAddr"`
	IntervalMs    int     `json:"intervalMs"`
	HighRiskPorts []int32 `json:"highRiskPorts"`
}

var defaultConfig = Config{
	ListenAddr: "0.0.0.0:7000",
	IntervalMs: 1000,
	HighRiskPorts: []int32{
		21, 22, 23, 135, 137, 138, 139, 445, 3389,
		3306, 5432, 6379, 27017, 11211, 9200, 2375, 2376,
	},
}

type ConfigStore struct {
	mu   sync.RWMutex
	path string
	cfg  Config
}

func LoadConfig() (*ConfigStore, error) {
	exe, err := os.Executable()
	if err != nil {
		return nil, err
	}
	path := filepath.Join(filepath.Dir(exe), "config.json")

	cs := &ConfigStore{path: path, cfg: defaultConfig}

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			if writeErr := cs.save(); writeErr != nil {
				return nil, writeErr
			}
			return cs, nil
		}
		return nil, err
	}

	var c Config
	if err := json.Unmarshal(data, &c); err != nil {
		return nil, err
	}
	if c.ListenAddr == "" {
		c.ListenAddr = defaultConfig.ListenAddr
	}
	if c.IntervalMs <= 0 {
		c.IntervalMs = defaultConfig.IntervalMs
	}
	if c.HighRiskPorts == nil {
		c.HighRiskPorts = defaultConfig.HighRiskPorts
	}
	cs.cfg = c
	return cs, nil
}

func (cs *ConfigStore) Snapshot() Config {
	cs.mu.RLock()
	defer cs.mu.RUnlock()
	c := cs.cfg
	c.HighRiskPorts = append([]int32(nil), cs.cfg.HighRiskPorts...)
	return c
}

func (cs *ConfigStore) Update(intervalMs int, highRisk []int32) (Config, error) {
	cs.mu.Lock()
	if intervalMs >= 200 {
		cs.cfg.IntervalMs = intervalMs
	}
	if highRisk != nil {
		cs.cfg.HighRiskPorts = highRisk
	}
	cs.mu.Unlock()
	if err := cs.save(); err != nil {
		return Config{}, err
	}
	return cs.Snapshot(), nil
}

func (cs *ConfigStore) save() error {
	cs.mu.RLock()
	data, err := json.MarshalIndent(cs.cfg, "", "  ")
	cs.mu.RUnlock()
	if err != nil {
		return err
	}
	return os.WriteFile(cs.path, data, 0o644)
}
