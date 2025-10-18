package ui

import (
	"fmt"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/google/uuid"
	"interactive-feedback-mcp/internal/types"
)

type ConversationSection struct {
	container   *fyne.Container
	historyList *widget.List
	copyButton  *widget.Button
	clearButton *widget.Button
	entries     []types.ConversationEntry
	onCopy      func(string)
	onClear     func()
}

func NewConversationSection() *ConversationSection {
	cs := &ConversationSection{
		entries: make([]types.ConversationEntry, 0),
	}

	cs.createUI()
	return cs
}

func (cs *ConversationSection) createUI() {
	// Create conversation list
	cs.historyList = widget.NewList(
		func() int {
			return len(cs.entries)
		},
		func() fyne.CanvasObject {
			return container.NewVBox(
				widget.NewLabel(""), // Role label
				widget.NewLabel(""), // Content label
				widget.NewLabel(""), // Timestamp label
			)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			if id >= len(cs.entries) {
				return
			}

			entry := cs.entries[id]
			container := obj.(*fyne.Container)

			// Role label with styling
			roleLabel := container.Objects[0].(*widget.Label)
			roleLabel.SetText(fmt.Sprintf("%s:", strings.Title(entry.Role)))

			// Style based on role
			if entry.Role == "user" {
				roleLabel.Importance = widget.MediumImportance
			} else {
				roleLabel.Importance = widget.HighImportance
			}

			// Content label
			contentLabel := container.Objects[1].(*widget.Label)
			contentLabel.SetText(entry.Content)
			contentLabel.Wrapping = fyne.TextWrapWord

			// Timestamp label
			timeLabel := container.Objects[2].(*widget.Label)
			timeLabel.SetText(entry.Timestamp.Format("15:04:05"))
			timeLabel.Importance = widget.LowImportance
		},
	)

	// Set minimum height for conversation list
	cs.historyList.Resize(fyne.NewSize(0, 200))

	// Create buttons
	cs.copyButton = widget.NewButton("Copy All", cs.copyAllConversation)
	cs.clearButton = widget.NewButton("Clear History", cs.clearHistory)

	// Button container
	buttonContainer := container.NewHBox(
		cs.copyButton,
		widget.NewSeparator(),
		cs.clearButton,
	)

	// Main container
	cs.container = container.NewVBox(
		widget.NewLabel("Conversation History"),
		cs.historyList,
		buttonContainer,
	)
}

func (cs *ConversationSection) AddEntry(role, content string) {
	entry := types.ConversationEntry{
		ID:        uuid.New().String(),
		Timestamp: time.Now(),
		Role:      role,
		Content:   content,
		IsCurrent: false,
	}

	// Mark previous entries as not current
	for i := range cs.entries {
		cs.entries[i].IsCurrent = false
	}

	// Add new entry
	cs.entries = append(cs.entries, entry)

	// Limit history size (keep last 50 entries)
	if len(cs.entries) > 50 {
		cs.entries = cs.entries[len(cs.entries)-50:]
	}

	// Refresh the list
	cs.historyList.Refresh()

	// Scroll to bottom
	cs.historyList.ScrollToBottom()
}

func (cs *ConversationSection) copyAllConversation() {
	var conversation strings.Builder

	for _, entry := range cs.entries {
		conversation.WriteString(fmt.Sprintf("%s: %s\n",
			strings.Title(entry.Role),
			entry.Content))
	}

	if cs.onCopy != nil {
		cs.onCopy(conversation.String())
	}
}

func (cs *ConversationSection) clearHistory() {
	cs.entries = make([]types.ConversationEntry, 0)
	cs.historyList.Refresh()

	if cs.onClear != nil {
		cs.onClear()
	}
}

func (cs *ConversationSection) GetContainer() *fyne.Container {
	return cs.container
}

func (cs *ConversationSection) SetOnCopy(callback func(string)) {
	cs.onCopy = callback
}

func (cs *ConversationSection) SetOnClear(callback func()) {
	cs.onClear = callback
}
