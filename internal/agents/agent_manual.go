package agents

import (
	"fmt"
	"reflect"

	"github.com/hubblew/pim/internal/ui"
)

type ManualAgent struct{}

var ManualAgentType = reflect.TypeOf(new(ManualAgent))

var _ AgentTool = (*ManualAgent)(nil)

func NewManualAgent() *ManualAgent {
	return &ManualAgent{}
}

func (a *ManualAgent) Descriptor() string {
	return "Manual (just output prompts, let me do it myself)"
}

func (a *ManualAgent) ExecuteCommand(command string) (string, error) {
	fmt.Println()
	fmt.Println("=== Manual Task ===")
	fmt.Println(command)
	fmt.Println("===================")
	fmt.Println()

	err := ui.WaitForKey("Press any key to continue...")
	if err != nil {
		return "", fmt.Errorf("failed to wait for key press: %w", err)
	}

	return "", nil
}
