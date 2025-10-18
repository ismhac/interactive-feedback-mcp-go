package types

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestProjectConfig_JSONSerialization(t *testing.T) {
	config := &ProjectConfig{
		RunCommand:            "npm run dev",
		ExecuteAutomatically:  true,
		CommandSectionVisible: false,
		ConversationHistory: []ConversationEntry{
			{
				ID:        "test-id",
				Timestamp: time.Now(),
				Role:      "user",
				Content:   "Test message",
				IsCurrent: false,
			},
		},
	}

	// Test JSON marshaling
	data, err := json.Marshal(config)
	assert.NoError(t, err)
	assert.NotEmpty(t, data)

	// Test JSON unmarshaling
	var decoded ProjectConfig
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)
	assert.Equal(t, config.RunCommand, decoded.RunCommand)
	assert.Equal(t, config.ExecuteAutomatically, decoded.ExecuteAutomatically)
	assert.Equal(t, config.CommandSectionVisible, decoded.CommandSectionVisible)
	assert.Len(t, decoded.ConversationHistory, 1)
	assert.Equal(t, config.ConversationHistory[0].ID, decoded.ConversationHistory[0].ID)
}

func TestConversationEntry_JSONSerialization(t *testing.T) {
	entry := ConversationEntry{
		ID:        "test-id",
		Timestamp: time.Now(),
		Role:      "assistant",
		Content:   "Test response",
		IsCurrent: true,
	}

	// Test JSON marshaling
	data, err := json.Marshal(entry)
	assert.NoError(t, err)
	assert.NotEmpty(t, data)

	// Test JSON unmarshaling
	var decoded ConversationEntry
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)
	assert.Equal(t, entry.ID, decoded.ID)
	assert.Equal(t, entry.Role, decoded.Role)
	assert.Equal(t, entry.Content, decoded.Content)
	assert.Equal(t, entry.IsCurrent, decoded.IsCurrent)
}

func TestFeedbackResult_JSONSerialization(t *testing.T) {
	result := &FeedbackResult{
		CommandLogs:         "Command output here",
		InteractiveFeedback: "User feedback here",
		ConversationHistory: []ConversationEntry{
			{
				ID:        "test-id",
				Timestamp: time.Now(),
				Role:      "user",
				Content:   "Test message",
				IsCurrent: false,
			},
		},
	}

	// Test JSON marshaling
	data, err := json.Marshal(result)
	assert.NoError(t, err)
	assert.NotEmpty(t, data)

	// Test JSON unmarshaling
	var decoded FeedbackResult
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)
	assert.Equal(t, result.CommandLogs, decoded.CommandLogs)
	assert.Equal(t, result.InteractiveFeedback, decoded.InteractiveFeedback)
	assert.Len(t, decoded.ConversationHistory, 1)
}
