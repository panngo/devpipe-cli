package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type TunnelConfig struct {
	UUID       string `json:"uuid"`
	SecurityKey string `json:"security_key"`
	TunnelID   string `json:"tunnel_id"`
	Port       string `json:"port"`
}

type ConfigManager struct {
	configPath string
}

func NewConfigManager() *ConfigManager {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "."
	}
	
	configDir := filepath.Join(homeDir, ".devpipe")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		// Fallback to current directory
		configDir = "."
	}
	
	return &ConfigManager{
		configPath: filepath.Join(configDir, "tunnel.json"),
	}
}

func (cm *ConfigManager) SaveTunnelConfig(config TunnelConfig) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	
	if err := os.WriteFile(cm.configPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}
	
	return nil
}

func (cm *ConfigManager) LoadTunnelConfig() (*TunnelConfig, error) {
	data, err := os.ReadFile(cm.configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil // No config file exists yet
		}
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}
	
	var config TunnelConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}
	
	return &config, nil
}

func (cm *ConfigManager) ClearTunnelConfig() error {
	if err := os.Remove(cm.configPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove config file: %w", err)
	}
	return nil
} 