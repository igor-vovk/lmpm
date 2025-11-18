package tpagents

import (
	"fmt"
	"reflect"

	"github.com/charmbracelet/glamour"
	"github.com/hubblew/pim/internal/ui"
)

type ManualAgent struct{}

var ManualAgentType = reflect.TypeOf(new(ManualAgent))

var _ TPAgentTool = (*ManualAgent)(nil)

func NewManualAgent() *ManualAgent {
	return &ManualAgent{}
}

func (a *ManualAgent) Descriptor() string {
	return "Manual (just output prompts, let me do it myself)"
}

func (a *ManualAgent) ExecuteCommand(command string) (string, error) {
	r, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(100),
	)
	if err != nil {
		return "", fmt.Errorf("failed to create markdown renderer: %w", err)
	}

	rendered, err := r.Render(command)
	if err != nil {
		return "", fmt.Errorf("failed to render markdown: %w", err)
	}

	fmt.Println()
	fmt.Println("=== Manual Task ===")
	fmt.Print(rendered)
	fmt.Println("===================")
	fmt.Println()

	err = ui.WaitForKey("Press any key to continue...")
	if err != nil {
		return "", fmt.Errorf("failed to wait for key press: %w", err)
	}

	return "", nil
}
