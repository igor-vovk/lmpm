package installer

import (
	"fmt"

	"github.com/hubblew/pim/internal/ui"
)

// UserPrompter handles user interaction for confirmation prompts.
type UserPrompter interface {
	ConfirmOverwrite(path string) (bool, error)
}

// InteractivePrompter prompts the user via stdin for confirmation.
type InteractivePrompter struct{}

var _ UserPrompter = (*InteractivePrompter)(nil)

func NewInteractivePrompter() *InteractivePrompter {
	return &InteractivePrompter{}
}

func (p *InteractivePrompter) ConfirmOverwrite(path string) (bool, error) {
	fmt.Printf("File %s already exists. Overwrite?\n", path)

	choice, err := ui.NewChoiceDialog("Please confirm:", ui.ChoicesYesNo()).Run()

	if err != nil {
		return false, fmt.Errorf("failed to get user input: %w", err)
	}
	if choice == nil {
		return false, nil
	}

	return choice.Value.(bool), nil
}

type AcceptAllPrompter struct{}

var _ UserPrompter = (*AcceptAllPrompter)(nil)

func NewAcceptAllPrompter() *AcceptAllPrompter {
	return &AcceptAllPrompter{}
}

func (p *AcceptAllPrompter) ConfirmOverwrite(_ string) (bool, error) {
	return true, nil
}
