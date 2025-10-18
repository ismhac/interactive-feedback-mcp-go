package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"interactive-feedback-mcp/internal/types"
)

type ConfigManager struct {
}

func NewConfigManager() (*ConfigManager, error) {
	return &ConfigManager{}, nil
}


func (cm *ConfigManager) LoadProjectConfig(projectPath string) *types.ProjectConfig {
	// Load from project directory
	configFile := filepath.Join(projectPath, ".interactive-feedback-config.json")
	if data, err := os.ReadFile(configFile); err == nil {
		var config types.ProjectConfig
		if json.Unmarshal(data, &config) == nil {
			return &config
		}
	}

	// Return default config if not found
	return &types.ProjectConfig{
		RunCommand:            "",
		ExecuteAutomatically:  false,
		CommandSectionVisible: false,
		ConversationHistory:   make([]types.ConversationEntry, 0),
	}
}

func (cm *ConfigManager) SaveProjectConfig(projectPath string, config *types.ProjectConfig) error {
	configFile := filepath.Join(projectPath, ".interactive-feedback-config.json")
	
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}
