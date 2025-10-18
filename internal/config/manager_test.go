package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"interactive-feedback-mcp/internal/types"
)

func TestNewConfigManager(t *testing.T) {
	// Test successful creation
	manager, err := NewConfigManager()
	require.NoError(t, err)
	assert.NotNil(t, manager)
}

func TestConfigManager_LoadProjectConfig(t *testing.T) {
	manager, err := NewConfigManager()
	require.NoError(t, err)

	// Test loading non-existent config (should return default)
	config := manager.LoadProjectConfig("/nonexistent/project")
	assert.NotNil(t, config)
	assert.Empty(t, config.RunCommand)
	assert.False(t, config.ExecuteAutomatically)
	assert.False(t, config.CommandSectionVisible)
	assert.Empty(t, config.ConversationHistory)
}

func TestConfigManager_SaveProjectConfig(t *testing.T) {
	manager, err := NewConfigManager()
	require.NoError(t, err)

	// Create test config
	testConfig := &types.ProjectConfig{
		RunCommand:            "npm run test",
		ExecuteAutomatically:  true,
		CommandSectionVisible: false,
		ConversationHistory: []types.ConversationEntry{
			{
				ID:      "test-id",
				Role:    "user",
				Content: "Test message",
			},
		},
	}

	// Create temporary directory for testing
	tempDir := t.TempDir()
	projectPath := filepath.Join(tempDir, "test-project")
	
	// Create project directory
	err = os.MkdirAll(projectPath, 0755)
	require.NoError(t, err)

	// Save config
	err = manager.SaveProjectConfig(projectPath, testConfig)
	require.NoError(t, err)

	// Load config and verify
	loadedConfig := manager.LoadProjectConfig(projectPath)
	assert.Equal(t, testConfig.RunCommand, loadedConfig.RunCommand)
	assert.Equal(t, testConfig.ExecuteAutomatically, loadedConfig.ExecuteAutomatically)
	assert.Equal(t, testConfig.CommandSectionVisible, loadedConfig.CommandSectionVisible)
	assert.Len(t, loadedConfig.ConversationHistory, 1)
	assert.Equal(t, "test-id", loadedConfig.ConversationHistory[0].ID)
}


func TestConfigManager_InvalidJSON(t *testing.T) {
	manager, err := NewConfigManager()
	require.NoError(t, err)

	// Create temporary directory for testing
	tempDir := t.TempDir()
	projectPath := filepath.Join(tempDir, "invalid-project")
	
	// Create project directory
	err = os.MkdirAll(projectPath, 0755)
	require.NoError(t, err)

	// Create invalid JSON file
	configFile := filepath.Join(projectPath, ".interactive-feedback-config.json")
	err = os.WriteFile(configFile, []byte("invalid json"), 0644)
	require.NoError(t, err)

	// Should return default config for invalid JSON
	config := manager.LoadProjectConfig(projectPath)
	assert.NotNil(t, config)
	assert.Empty(t, config.RunCommand)
}
