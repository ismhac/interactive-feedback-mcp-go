package types

import "time"

// ProjectConfig represents configuration for a specific project
type ProjectConfig struct {
	RunCommand              string                `json:"run_command"`
	ExecuteAutomatically    bool                  `json:"execute_automatically"`
	CommandSectionVisible   bool                  `json:"command_section_visible"`
	ConversationHistory     []ConversationEntry   `json:"conversation_history"`
}

// ConversationEntry represents a single message in the conversation
type ConversationEntry struct {
	ID        string    `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	Role      string    `json:"role"` // "user" or "assistant"
	Content   string    `json:"content"`
	IsCurrent bool      `json:"is_current"`
}

// CommandHandle represents a running command process
type CommandHandle struct {
	PID       int
	StartTime time.Time
	IsRunning bool
	Output    chan string
	Done      chan error
}

// FeedbackResult represents the final output
type FeedbackResult struct {
	CommandLogs         string                `json:"command_logs"`
	InteractiveFeedback string                `json:"interactive_feedback"`
	ConversationHistory []ConversationEntry   `json:"conversation_history"`
}
