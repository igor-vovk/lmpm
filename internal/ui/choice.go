package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// Choice represents a single option in a choice selector.
type Choice struct {
	Label string
	Value any
}

func ChoicesYesNo() []Choice {
	return []Choice{
		{Label: "Yes", Value: true},
		{Label: "No", Value: false},
	}
}

// Layout defines how choices are rendered.
type Layout int

const (
	LayoutHorizontal Layout = iota
	LayoutVertical
)

// ChoiceDialog is a Bubble Tea model for selecting from multiple choices.
type ChoiceDialog struct {
	Prompt      string
	Choices     []Choice
	Cursor      int
	Selected    bool
	Cancelled   bool
	StyleConfig StyleConfig
	Layout      Layout
}

var _ tea.Model = (*ChoiceDialog)(nil)

// NewChoiceDialog creates a new choice selector model.
func NewChoiceDialog(prompt string, choices []Choice) ChoiceDialog {
	return ChoiceDialog{
		Prompt:      prompt,
		Choices:     choices,
		Cursor:      0,
		StyleConfig: DefaultStyleConfig(),
		Layout:      LayoutHorizontal,
	}
}

func (d ChoiceDialog) Vertical() ChoiceDialog {
	d.Layout = LayoutVertical
	return d
}

func (d ChoiceDialog) Init() tea.Cmd {
	return nil
}

func (d ChoiceDialog) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			d.Cancelled = true
			return d, tea.Quit

		case "enter":
			d.Selected = true
			return d, tea.Quit

		case "left", "up":
			d.Cursor--
			if d.Cursor < 0 {
				d.Cursor = len(d.Choices) - 1
			}

		case "right", "down":
			d.Cursor++
			if d.Cursor >= len(d.Choices) {
				d.Cursor = 0
			}

		case "home":
			d.Cursor = 0

		case "end":
			d.Cursor = len(d.Choices) - 1

		default:
			// Check if the key matches any choice label's first character
			key := strings.ToLower(msg.String())
			if len(key) == 1 {
				for i, choice := range d.Choices {
					if len(choice.Label) > 0 && strings.ToLower(string(choice.Label[0])) == key {
						d.Cursor = i
						return d, nil
					}
				}
			}
		}
	}

	return d, nil
}

func (d ChoiceDialog) View() string {
	if d.Selected || d.Cancelled {
		return ""
	}

	var b strings.Builder

	if d.Prompt != "" {
		b.WriteString(d.StyleConfig.PromptStyle.Render(d.Prompt))
		if d.Layout == LayoutVertical {
			b.WriteString("\n")
		} else {
			b.WriteString(" ")
		}
	}

	if d.Layout == LayoutVertical {
		for i, choice := range d.Choices {
			cursor := "  "
			if i == d.Cursor {
				cursor = "> "
				b.WriteString(cursor)
				b.WriteString(d.StyleConfig.HighlightStyle.Render(choice.Label))
			} else {
				b.WriteString(cursor)
				b.WriteString(d.StyleConfig.NormalStyle.Render(choice.Label))
			}
			b.WriteString("\n")
		}
	} else {
		for i, choice := range d.Choices {
			if i > 0 {
				b.WriteString(" / ")
			}

			if i == d.Cursor {
				b.WriteString(d.StyleConfig.HighlightStyle.Render(choice.Label))
			} else {
				b.WriteString(d.StyleConfig.NormalStyle.Render(choice.Label))
			}
		}
		b.WriteString("\n")
	}

	b.WriteString(d.StyleConfig.HelpStyle.Render("(Use arrow keys to select, Enter to confirm, Esc to cancel)"))
	b.WriteString("\n")

	return b.String()
}

// GetHighlightedChoice returns the currently highlighted choice.
func (d ChoiceDialog) GetHighlightedChoice() *Choice {
	if d.Cursor < 0 || d.Cursor >= len(d.Choices) {
		return nil
	}
	return &d.Choices[d.Cursor]
}

// GetSelectedChoice returns the currently selected choice.
func (d ChoiceDialog) GetSelectedChoice() *Choice {
	if d.Cancelled || !d.Selected {
		return nil
	}
	return d.GetHighlightedChoice()
}

// Run is a convenience method to run the choice selector and return the result.
func (d ChoiceDialog) Run() (*Choice, error) {
	p := tea.NewProgram(d)
	finalModel, err := p.Run()
	if err != nil {
		return nil, err
	}

	result := finalModel.(ChoiceDialog)
	return result.GetSelectedChoice(), nil
}
