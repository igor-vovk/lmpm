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

const yes = "y"
const no = "n"

func (p *InteractivePrompter) ConfirmOverwrite(path string) (bool, error) {
	fmt.Printf("File %s already exists. Overwrite?\n", path)

	choice := ui.NewChoiceDialog(
		"Please confirm:",
		[]ui.Choice{
			{Label: "Yes", Value: yes},
			{Label: "No", Value: no},
		},
	)

	response, err := choice.Run()

	if err != nil {
		return false, fmt.Errorf("failed to get user input: %w", err)
	}
	if response == nil {
		return false, nil
	}

	return response.Value == yes, nil
}

type AcceptAllPrompter struct{}

var _ UserPrompter = (*AcceptAllPrompter)(nil)

func NewAcceptAllPrompter() *AcceptAllPrompter {
	return &AcceptAllPrompter{}
}

func (p *AcceptAllPrompter) ConfirmOverwrite(_ string) (bool, error) {
	return true, nil
}
