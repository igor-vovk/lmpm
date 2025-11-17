package main

import (
	"fmt"
	"os"

	"github.com/hubblew/pim/internal/ui"
)

func demoYesNo() {
	choices := []ui.Choice{
		{Label: "yes", Value: true},
		{Label: "no", Value: false},
	}

	model := ui.NewChoiceDialog("Do you want to continue?", choices)
	model.Cursor = 1 // Default to "no"

	choice, err := model.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return
	}

	if choice == nil {
		fmt.Println("\n❌ Cancelled")
		return
	}

	if choice.Value.(bool) {
		fmt.Println("\n✅ You selected: Yes - Continuing...")
	} else {
		fmt.Println("\n❌ You selected: No - Aborting...")
	}
}

func demoVertical() {
	choices := []ui.Choice{
		{Label: "small", Value: "s"},
		{Label: "medium", Value: "m"},
		{Label: "large", Value: "l"},
		{Label: "extra-large", Value: "xl"},
	}

	model := ui.NewVerticalChoiceDialog("Select your size:", choices)
	model.Cursor = 1 // Default to "medium"

	choice, err := model.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return
	}

	if choice == nil {
		fmt.Println("\n❌ Cancelled")
		return
	}

	fmt.Printf("\n✅ You selected: %s (value: %s)\n", choice.Label, choice.Value)
}

func demoEnvironment() {
	choices := []ui.Choice{
		{Label: "development", Value: "dev"},
		{Label: "staging", Value: "staging"},
		{Label: "production", Value: "prod"},
	}

	model := ui.NewChoiceDialog("Deploy to which environment?", choices)

	choice, err := model.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return
	}

	if choice == nil {
		fmt.Println("\n❌ Deployment cancelled")
		return
	}

	fmt.Printf("\n🚀 Deploying to: %s (environment: %s)\n", choice.Label, choice.Value)
}

func runChoiceDemos() {
	fmt.Println("=== Demo 1: Yes/No Confirmation ===")
	fmt.Println("")
	demoYesNo()

	fmt.Println("\n" + "─────────────────────────────────────────────────")
	fmt.Println("")

	fmt.Println("=== Demo 2: Multiple Options (vertical) ===")
	fmt.Println("")
	demoVertical()

	fmt.Println("\n" + "─────────────────────────────────────────────────")
	fmt.Println("")

	fmt.Println("=== Demo 3: Environment Selection ===")
	fmt.Println("")
	demoEnvironment()
}
