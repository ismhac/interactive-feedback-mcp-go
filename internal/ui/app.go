package ui

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"interactive-feedback-mcp/internal/config"
	"interactive-feedback-mcp/internal/executor"
	"interactive-feedback-mcp/internal/types"
)

type FeedbackApp struct {
	app                fyne.App
	window             fyne.Window
	projectDirectory   string
	prompt             string
	configManager      *config.ConfigManager
	commandExecutor    *executor.CommandExecutor
	currentHandle      *types.CommandHandle

	// UI Components
	commandEntry       *widget.Entry
	runButton          *widget.Button
	consoleText        *widget.Entry
	feedbackText       *widget.Entry
	submitButton       *widget.Button
	commandSection     *widget.Card
	conversationSection *ConversationSection
}

func NewFeedbackApp(projectDirectory, prompt string) (*FeedbackApp, error) {
	myApp := app.NewWithID("com.interactivefeedback.mcp")

	configManager, err := config.NewConfigManager()
	if err != nil {
		return nil, fmt.Errorf("failed to create config manager: %w", err)
	}

	commandExecutor := executor.NewCommandExecutor()

	window := myApp.NewWindow("Interactive Feedback MCP")
	window.Resize(fyne.NewSize(800, 600))
	window.CenterOnScreen()

	feedbackApp := &FeedbackApp{
		app:              myApp,
		window:           window,
		projectDirectory: projectDirectory,
		prompt:           prompt,
		configManager:    configManager,
		commandExecutor:  commandExecutor,
	}

	feedbackApp.createUI()
	feedbackApp.loadConfig()

	return feedbackApp, nil
}

func (fa *FeedbackApp) createUI() {
	// Command Section
	fa.commandEntry = widget.NewEntry()
	fa.commandEntry.SetPlaceHolder("Enter command to run")
	fa.commandEntry.OnSubmitted = func(text string) {
		fa.runCommand()
	}

	fa.runButton = widget.NewButton("Run", fa.runCommand)
	fa.runButton.Importance = widget.HighImportance

	commandContainer := container.NewBorder(nil, nil, nil, fa.runButton, fa.commandEntry)

	fa.commandSection = widget.NewCard("Command", "", commandContainer)

	// Console Section
	fa.consoleText = widget.NewMultiLineEntry()
	fa.consoleText.SetText("")
	fa.consoleText.Disable()
	fa.consoleText.Wrapping = fyne.TextWrapWord

	clearButton := widget.NewButton("Clear", func() {
		fa.consoleText.SetText("")
	})

	consoleContainer := container.NewBorder(nil, clearButton, nil, nil, fa.consoleText)
	consoleCard := widget.NewCard("Console", "", consoleContainer)

	// Conversation History Section (NEW)
	fa.conversationSection = NewConversationSection()
	fa.conversationSection.SetOnCopy(fa.copyToClipboard)
	fa.conversationSection.SetOnClear(fa.clearConversationHistory)

	// Add initial conversation entry
	fa.conversationSection.AddEntry("assistant", fa.prompt)

	// Feedback Section
	promptLabel := widget.NewLabel(fa.prompt)
	promptLabel.Wrapping = fyne.TextWrapWord

	fa.feedbackText = widget.NewMultiLineEntry()
	fa.feedbackText.SetPlaceHolder("Enter your feedback here...")
	fa.feedbackText.Wrapping = fyne.TextWrapWord

	fa.submitButton = widget.NewButton("Submit Feedback", fa.submitFeedback)
	fa.submitButton.Importance = widget.HighImportance

	feedbackContainer := container.NewVBox(
		promptLabel,
		fa.feedbackText,
		fa.submitButton,
	)
	feedbackCard := widget.NewCard("Feedback", "", feedbackContainer)

	// Main layout
	content := container.NewVBox(
		fa.commandSection,
		consoleCard,
		fa.conversationSection.GetContainer(), // NEW: Conversation section
		feedbackCard,
	)

	fa.window.SetContent(content)
}

func (fa *FeedbackApp) loadConfig() {
	config := fa.configManager.LoadProjectConfig(fa.projectDirectory)
	fa.commandEntry.SetText(config.RunCommand)

	if config.ExecuteAutomatically && config.RunCommand != "" {
		fa.runCommand()
	}
}

func (fa *FeedbackApp) runCommand() {
	command := strings.TrimSpace(fa.commandEntry.Text)
	if command == "" {
		return
	}

	if fa.currentHandle != nil && fa.currentHandle.IsRunning {
		// Stop current command
		fa.commandExecutor.KillProcessTree(fa.currentHandle.PID)
		fa.runButton.SetText("Run")
		fa.currentHandle = nil
		return
	}

	// Start new command
	fa.runButton.SetText("Stop")
	fa.appendToConsole(fmt.Sprintf("$ %s\n", command))

	handle, err := fa.commandExecutor.ExecuteCommand(command, fa.projectDirectory)
	if err != nil {
		fa.appendToConsole(fmt.Sprintf("Error: %v\n", err))
		fa.runButton.SetText("Run")
		return
	}

	fa.currentHandle = handle

	// Start goroutine to read output
	go func() {
		for output := range handle.Output {
			fa.appendToConsole(output)
		}

		// Command finished
		fa.runButton.SetText("Run")
		fa.currentHandle = nil
	}()
}

func (fa *FeedbackApp) appendToConsole(text string) {
	currentText := fa.consoleText.Text
	fa.consoleText.SetText(currentText + text)
	fa.consoleText.CursorRow = len(strings.Split(fa.consoleText.Text, "\n")) - 1
}

func (fa *FeedbackApp) copyToClipboard(text string) {
	// Copy to system clipboard
	fa.window.Clipboard().SetContent(text)

	// Show notification
	dialog.ShowInformation("Copied", "Conversation copied to clipboard!", fa.window)
}

func (fa *FeedbackApp) clearConversationHistory() {
	// Clear conversation history
	fa.conversationSection.clearHistory()

	// Re-add the current prompt
	fa.conversationSection.AddEntry("assistant", fa.prompt)
}

func (fa *FeedbackApp) submitFeedback() {
	// Add user feedback to conversation history
	feedback := fa.feedbackText.Text
	if strings.TrimSpace(feedback) != "" {
		fa.conversationSection.AddEntry("user", feedback)
	}

	// Save configuration
	config := fa.configManager.LoadProjectConfig(fa.projectDirectory)
	config.RunCommand = fa.commandEntry.Text

	if err := fa.configManager.SaveProjectConfig(fa.projectDirectory, config); err != nil {
		dialog.ShowError(fmt.Errorf("failed to save config: %w", err), fa.window)
		return
	}

	// Create response
	response := &types.FeedbackResult{
		CommandLogs:         fa.consoleText.Text,
		InteractiveFeedback: fa.feedbackText.Text,
		ConversationHistory: fa.conversationSection.entries,
	}

	// Save response to file (for MCP server to read)
	responseData, _ := json.Marshal(response)
	os.WriteFile("feedback_result.json", responseData, 0644)

	fa.window.Close()
}

func (fa *FeedbackApp) Run() {
	fa.window.ShowAndRun()
}
