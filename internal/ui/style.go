package ui

import "github.com/charmbracelet/lipgloss"

// StyleConfig holds styling configuration for UI components.
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
