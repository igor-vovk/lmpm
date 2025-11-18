package ui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestNewChoiceDialog(t *testing.T) {
	model := NewChoiceDialog("Confirm?", ChoicesYesNo())

	if model.Prompt != "Confirm?" {
		t.Errorf("expected prompt 'Confirm?', got '%s'", model.Prompt)
	}

	if len(model.Choices) != 2 {
		t.Errorf("expected 2 choices, got %d", len(model.Choices))
	}

	if model.Cursor != 0 {
		t.Errorf("expected cursor at 0, got %d", model.Cursor)
	}
}

func TestChoiceDialog_Navigation(t *testing.T) {
	choices := []Choice{
		{Label: "One", Value: 1},
		{Label: "Two", Value: 2},
		{Label: "Three", Value: 3},
	}

	model := NewChoiceDialog("Select:", choices)

	// Test right navigation
	result, _ := model.Update(tea.KeyMsg{Type: tea.KeyRight})
	model = result.(ChoiceDialog)
	if model.Cursor != 1 {
		t.Errorf("expected cursor at 1, got %d", model.Cursor)
	}

	// Test left navigation
	result, _ = model.Update(tea.KeyMsg{Type: tea.KeyLeft})
	model = result.(ChoiceDialog)
	if model.Cursor != 0 {
		t.Errorf("expected cursor at 0, got %d", model.Cursor)
	}

	// Test end key
	result, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnd})
	model = result.(ChoiceDialog)
	if model.Cursor != 2 {
		t.Errorf("expected cursor at 2, got %d", model.Cursor)
	}

	// Test home key
	result, _ = model.Update(tea.KeyMsg{Type: tea.KeyHome})
	model = result.(ChoiceDialog)
	if model.Cursor != 0 {
		t.Errorf("expected cursor at 0, got %d", model.Cursor)
	}
}

func TestChoiceDialog_Selection(t *testing.T) {
	choices := ChoicesYesNo()
	if choices[1].Value != false {
		t.Fatalf("expected second choice value to be false, got %v", choices[1].Value)
	}

	model := NewChoiceDialog("Confirm?", choices)
	model.Cursor = 1

	result, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	model = result.(ChoiceDialog)

	if !model.Selected {
		t.Error("expected model to be selected")
	}

	if cmd == nil {
		t.Error("expected quit command")
	}

	choice := model.GetSelectedChoice()
	if choice == nil {
		t.Fatal("expected selected choice")
	}

	if choice.Value != false {
		t.Errorf("expected value false, got %v", choice.Value)
	}
}

func TestChoiceDialog_Cancel(t *testing.T) {
	model := NewChoiceDialog("Confirm?", ChoicesYesNo())
	result, _ := model.Update(tea.KeyMsg{Type: tea.KeyEsc})
	model = result.(ChoiceDialog)

	if !model.Cancelled {
		t.Error("expected model to be cancelled")
	}

	choice := model.GetSelectedChoice()
	if choice != nil {
		t.Error("expected no selected choice when cancelled")
	}
}

func TestChoiceDialog_QuickSelect(t *testing.T) {
	model := NewChoiceDialog("Confirm?", ChoicesYesNo())

	// Press 'n' to quickly highlight "No"
	result, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
	model = result.(ChoiceDialog)

	if model.Selected {
		t.Error("expected model not to be selected")
	}

	if model.Cursor != 1 {
		t.Errorf("expected cursor at 1 (No), got %d", model.Cursor)
	}

	choice := model.GetHighlightedChoice()
	if choice == nil || choice.Value != false {
		t.Error("expected 'No' to be highlighted")
	}

	result, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	model = result.(ChoiceDialog)

	if !model.Selected {
		t.Error("expected model to be selected")
	}
	choice = model.GetSelectedChoice()
	if choice == nil || choice.Value != false {
		t.Error("expected 'No' to be selected")
	}
}

func TestChoiceDialog_CircularNavigation(t *testing.T) {
	choices := []Choice{
		{Label: "One", Value: 1},
		{Label: "Two", Value: 2},
		{Label: "Three", Value: 3},
	}

	model := NewChoiceDialog("Select:", choices)
	model.Cursor = 0

	// Going left from first wraps to last
	result, _ := model.Update(tea.KeyMsg{Type: tea.KeyLeft})
	model = result.(ChoiceDialog)
	if model.Cursor != 2 {
		t.Errorf("expected cursor at 2 (last), got %d", model.Cursor)
	}

	// Going right from last wraps to first
	model.Cursor = 2
	result, _ = model.Update(tea.KeyMsg{Type: tea.KeyRight})
	model = result.(ChoiceDialog)
	if model.Cursor != 0 {
		t.Errorf("expected cursor at 0 (first), got %d", model.Cursor)
	}
}

func TestChoiceDialog_View(t *testing.T) {
	model := NewChoiceDialog("Confirm?", ChoicesYesNo())
	view := model.View()

	if view == "" {
		t.Error("expected non-empty view")
	}

	// After selection, view should be empty
	model.Selected = true
	view = model.View()
	if view != "" {
		t.Error("expected empty view after selection")
	}
}
