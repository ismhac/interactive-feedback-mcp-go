package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"interactive-feedback-mcp/internal/config"
	"interactive-feedback-mcp/internal/types"
)

type MCPRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
}

type MCPResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   *MCPError   `json:"error,omitempty"`
}

type MCPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Tool struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	InputSchema map[string]interface{} `json:"inputSchema"`
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		var request MCPRequest
		if err := json.Unmarshal([]byte(line), &request); err != nil {
			sendError(nil, -32700, "Parse error", err.Error())
			continue
		}

		response := handleRequest(request)
		sendResponse(response)
	}
}

func handleRequest(request MCPRequest) MCPResponse {
	switch request.Method {
	case "initialize":
		return handleInitialize(request)
	case "tools/list":
		return handleToolsList(request)
	case "tools/call":
		return handleToolsCall(request)
	default:
		return MCPResponse{
			JSONRPC: "2.0",
			ID:      request.ID,
			Error: &MCPError{
				Code:    -32601,
				Message: "Method not found",
			},
		}
	}
}

func handleInitialize(request MCPRequest) MCPResponse {
	return MCPResponse{
		JSONRPC: "2.0",
		ID:      request.ID,
		Result: map[string]interface{}{
			"protocolVersion": "2024-11-05",
			"capabilities": map[string]interface{}{
				"tools": map[string]interface{}{
					"listChanged": true,
				},
			},
			"serverInfo": map[string]string{
				"name":    "interactive-feedback-mcp",
				"version": "1.0.0",
			},
		},
	}
}

func handleToolsList(request MCPRequest) MCPResponse {
	tools := []Tool{
		{
			Name:        "interactive_feedback",
			Description: "Get interactive feedback from user for development tasks",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"projectDirectory": map[string]interface{}{
						"type":        "string",
						"description": "The project directory path",
					},
					"prompt": map[string]interface{}{
						"type":        "string",
						"description": "The prompt to show to the user",
					},
					"previousUserRequest": map[string]interface{}{
						"type":        "string",
						"description": "The previous user request that triggered this interactive feedback",
					},
				},
				"required": []string{"projectDirectory", "prompt", "previousUserRequest"},
			},
		},
	}

	return MCPResponse{
		JSONRPC: "2.0",
		ID:      request.ID,
		Result: map[string]interface{}{
			"tools": tools,
		},
	}
}

func handleToolsCall(request MCPRequest) MCPResponse {
	// Parse the tool call parameters
	paramsBytes, err := json.Marshal(request.Params)
	if err != nil {
		return MCPResponse{
			JSONRPC: "2.0",
			ID:      request.ID,
			Error: &MCPError{
				Code:    -32602,
				Message: "Invalid params",
			},
		}
	}

	var toolCall struct {
		Name      string                 `json:"name"`
		Arguments map[string]interface{} `json:"arguments"`
	}

	if err := json.Unmarshal(paramsBytes, &toolCall); err != nil {
		return MCPResponse{
			JSONRPC: "2.0",
			ID:      request.ID,
			Error: &MCPError{
				Code:    -32602,
				Message: "Invalid params",
			},
		}
	}

	if toolCall.Name != "interactive_feedback" {
		return MCPResponse{
			JSONRPC: "2.0",
			ID:      request.ID,
			Error: &MCPError{
				Code:    -32601,
				Message: "Unknown tool",
			},
		}
	}

	// Extract arguments
	projectDir, _ := toolCall.Arguments["projectDirectory"].(string)
	prompt, _ := toolCall.Arguments["prompt"].(string)
	previousUserRequest, _ := toolCall.Arguments["previousUserRequest"].(string)

	if projectDir == "" {
		projectDir = "."
	}

	// Run interactive feedback with single popup GUI
	result := runInteractiveFeedbackWithSinglePopupGUI(projectDir, prompt, previousUserRequest)

	return MCPResponse{
		JSONRPC: "2.0",
		ID:      request.ID,
		Result: map[string]interface{}{
			"content": []map[string]interface{}{
				{
					"type": "text",
					"text": result,
				},
			},
		},
	}
}

func runInteractiveFeedbackWithSinglePopupGUI(projectDir, prompt, previousUserRequest string) string {
	// Load or create config
	configManager, err := config.NewConfigManager()
	if err != nil {
		return fmt.Sprintf("Error creating config manager: %v", err)
	}
	
	projectConfig := configManager.LoadProjectConfig(projectDir)
	if projectConfig == nil {
		projectConfig = &types.ProjectConfig{
			RunCommand:              "",
			ExecuteAutomatically:    false,
			CommandSectionVisible:   true,
			ConversationHistory:     []types.ConversationEntry{},
		}
	}

	// STEP 1: Add previous user request to conversation history FIRST
	if previousUserRequest != "" {
		userEntry := types.ConversationEntry{
			ID:        uuid.New().String(),
			Timestamp: time.Now(),
			Role:      "user",
			Content:   previousUserRequest,
			IsCurrent: false,
		}
		projectConfig.ConversationHistory = append(projectConfig.ConversationHistory, userEntry)
	}

	// STEP 2: Add agent prompt to conversation history
	assistantEntry := types.ConversationEntry{
		ID:        uuid.New().String(),
		Timestamp: time.Now(),
		Role:      "assistant",
		Content:   prompt,
		IsCurrent: false,
	}
	projectConfig.ConversationHistory = append(projectConfig.ConversationHistory, assistantEntry)

	// STEP 2.5: Trim conversation history to prevent file bloat (keep last 10 entries)
	projectConfig.ConversationHistory = trimConversationHistory(projectConfig.ConversationHistory, 10)

	// STEP 3: Auto-add to .gitignore if not already added
	ensureGitignoreEntry(projectDir)

	// STEP 4: Save config to disk BEFORE calling GUI
	configManager.SaveProjectConfig(projectDir, projectConfig)

	// Get the directory of the current executable
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Sprintf("Error getting executable path: %v", err)
	}
	
	execDir := filepath.Dir(execPath)
	
	// Find the single popup desktop GUI
	desktopGUI := filepath.Join(execDir, "desktop_gui_single.py")
	if _, err := os.Stat(desktopGUI); os.IsNotExist(err) {
		// Try alternative path
		desktopGUI = filepath.Join(execDir, "..", "desktop_gui_single.py")
		if _, err := os.Stat(desktopGUI); os.IsNotExist(err) {
			return "Single popup desktop GUI not found. Please ensure desktop_gui_single.py is in the project directory."
		}
	}

	// STEP 4: Launch single popup desktop GUI AFTER saving config
	cmd := exec.Command("python3", desktopGUI, projectDir, prompt)
	cmd.Dir = filepath.Dir(desktopGUI)
	
	// Capture output
	output, err := cmd.Output()
	if err != nil {
		return fmt.Sprintf("Error running single popup desktop GUI: %v", err)
	}
	
	userFeedback := strings.TrimSpace(string(output))
	// Allow empty feedback - user can choose not to provide feedback
	
	// STEP 5: Add user feedback to conversation only if feedback is provided
	if userFeedback != "" {
		feedbackEntry := types.ConversationEntry{
			ID:        uuid.New().String(),
			Timestamp: time.Now(),
			Role:      "user",
			Content:   userFeedback,
			IsCurrent: false,
		}

		projectConfig.ConversationHistory = append(projectConfig.ConversationHistory, feedbackEntry)

		// Trim conversation history again after adding feedback
		projectConfig.ConversationHistory = trimConversationHistory(projectConfig.ConversationHistory, 10)
	}

	// Save updated config with user feedback
	configManager.SaveProjectConfig(projectDir, projectConfig)

	// Create feedback result
	feedbackResult := types.FeedbackResult{
		CommandLogs:         "",
		InteractiveFeedback: userFeedback,
		ConversationHistory: projectConfig.ConversationHistory,
	}

	// Convert to JSON
	resultBytes, err := json.MarshalIndent(feedbackResult, "", "  ")
	if err != nil {
		return fmt.Sprintf("Error creating feedback result: %v", err)
	}

	return string(resultBytes)
}

func sendResponse(response MCPResponse) {
	responseBytes, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error marshaling response: %v", err)
		return
	}
	fmt.Println(string(responseBytes))
}

func sendError(id interface{}, code int, message, data string) {
	response := MCPResponse{
		JSONRPC: "2.0",
		ID:      id,
		Error: &MCPError{
			Code:    code,
			Message: message,
		},
	}
	sendResponse(response)
}

func trimConversationHistory(history []types.ConversationEntry, maxEntries int) []types.ConversationEntry {
	if len(history) <= maxEntries {
		return history
	}
	
	// Keep the last maxEntries entries
	startIndex := len(history) - maxEntries
	return history[startIndex:]
}

func ensureGitignoreEntry(projectDir string) {
	gitignorePath := filepath.Join(projectDir, ".gitignore")
	configFileName := ".interactive-feedback-config.json"
	
	// Check if .gitignore exists
	if _, err := os.Stat(gitignorePath); os.IsNotExist(err) {
		// Create .gitignore if it doesn't exist
		content := fmt.Sprintf("# Interactive Feedback MCP Configuration\n%s\n", configFileName)
		os.WriteFile(gitignorePath, []byte(content), 0644)
		return
	}
	
	// Read existing .gitignore
	content, err := os.ReadFile(gitignorePath)
	if err != nil {
		return // Skip if can't read
	}
	
	// Check if already contains our entry
	contentStr := string(content)
	if strings.Contains(contentStr, configFileName) {
		return // Already added
	}
	
	// Add our entry to .gitignore
	entry := fmt.Sprintf("\n# Interactive Feedback MCP Configuration\n%s\n", configFileName)
	os.WriteFile(gitignorePath, []byte(contentStr+entry), 0644)
}
