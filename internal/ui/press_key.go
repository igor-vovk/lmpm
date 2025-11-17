package ui

import (
	tea "github.com/charmbracelet/bubbletea"
)

// PressKeyDialog is a Bubble Tea model for waiting for any key press.
type PressKeyDialog struct {
	Message     string
	Pressed     bool
	StyleConfig StyleConfig
}

var _ tea.Model = (*PressKeyDialog)(nil)

// NewPressKeyDialog creates a new press key dialog.
func NewPressKeyDialog(message string) PressKeyDialog {
	return PressKeyDialog{
		Message:     message,
		StyleConfig: DefaultStyleConfig(),
	}
}

func (d PressKeyDialog) Init() tea.Cmd {
	return nil
}

func (d PressKeyDialog) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() != "" {
			d.Pressed = true
			return d, tea.Quit
		}
	}

	return d, nil
}

func (d PressKeyDialog) View() string {
	if d.Pressed {
		return ""
	}

	return d.StyleConfig.HelpStyle.Render(d.Message) + "\n"
}

// Run is a convenience method to run the press key dialog.
func (d PressKeyDialog) Run() error {
	p := tea.NewProgram(d)
	_, err := p.Run()
	return err
}

// WaitForKey displays a message and waits for any key press.
func WaitForKey(message string) error {
	dialog := NewPressKeyDialog(message)
	return dialog.Run()
}
