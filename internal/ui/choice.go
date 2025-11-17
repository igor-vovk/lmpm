package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Choice represents a single option in a choice selector.
type Choice struct {
	Label string
	Value any
}

// ChoiceDialog is a Bubble Tea model for selecting from multiple choices.
type ChoiceDialog struct {
	Prompt      string
	Choices     []Choice
	Cursor      int
	Selected    bool
	Cancelled   bool
	StyleConfig StyleConfig
}

var _ tea.Model = (*ChoiceDialog)(nil)

// StyleConfig holds styling configuration for the choice component.
type StyleConfig struct {
	HighlightStyle lipgloss.Style
	NormalStyle    lipgloss.Style
	PromptStyle    lipgloss.Style
	HelpStyle      lipgloss.Style
}

// DefaultStyleConfig returns sensible default styles.
func DefaultStyleConfig() StyleConfig {
	return StyleConfig{
		HighlightStyle: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("2")),
		NormalStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")),
		PromptStyle: lipgloss.NewStyle(),
		HelpStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Faint(true),
	}
}

// NewChoiceDialog creates a new choice selector model.
func NewChoiceDialog(prompt string, choices []Choice) ChoiceDialog {
	return ChoiceDialog{
		Prompt:      prompt,
		Choices:     choices,
		Cursor:      0,
		StyleConfig: DefaultStyleConfig(),
	}
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

		case "left", "h", "up", "k":
			d.Cursor--
			if d.Cursor < 0 {
				d.Cursor = len(d.Choices) - 1
			}

		case "right", "l", "down", "j":
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
						d.Selected = true
						return d, tea.Quit
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
		b.WriteString(" ")
	}

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
	b.WriteString(d.StyleConfig.HelpStyle.Render("(Use arrow keys to select, Enter to confirm, Esc to cancel)"))
	b.WriteString("\n")

	return b.String()
}

// GetSelectedChoice returns the currently selected choice.
func (d ChoiceDialog) GetSelectedChoice() *Choice {
	if d.Cancelled || !d.Selected || d.Cursor >= len(d.Choices) {
		return nil
	}
	return &d.Choices[d.Cursor]
}

// GetSelectedValue returns the value of the selected choice.
func (d ChoiceDialog) GetSelectedValue() any {
	choice := d.GetSelectedChoice()
	if choice == nil {
		return nil
	}
	return choice.Value
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
