package main

import (
	"fmt"
	"os"

	"github.com/hubblew/pim/internal/ui"
)

type Demo struct {
	Name string
	Run  func()
}

func main() {
	demos := []Demo{
		{Name: "Choice Component Demos", Run: runChoiceDemos},
		{Name: "Press Key Demo", Run: runPressKeyDemos},
	}

	choices := make([]ui.Choice, len(demos))
	for i, demo := range demos {
		choices[i] = ui.Choice{
			Label: demo.Name,
			Value: demo,
		}
	}

	fmt.Println("\n" + "═════════════════════════════════════════════════")
	fmt.Println(" PIM UI Component Demos")
	fmt.Println("═════════════════════════════════════════════════")

	model := ui.NewVerticalChoiceDialog("Select a demo to run:", choices)
	choice, err := model.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if choice == nil {
		fmt.Println("\n❌ Cancelled")
		os.Exit(0)
	}

	selectedDemo := choice.Value.(Demo)
	fmt.Println("")
	selectedDemo.Run()
	fmt.Println("")
	fmt.Println("✅ Demo completed!")

	main()
}
