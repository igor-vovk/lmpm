package main

import (
	"fmt"

	"github.com/hubblew/pim/internal/ui"
)

func demoPressKey() {
	fmt.Println("=== Press Key Demo ===")
	fmt.Println()
	fmt.Println("This demo shows how to wait for user input.")
	fmt.Println()

	err := ui.WaitForKey("Press any key to continue...")
	if err != nil {
		fmt.Printf("\nError: %v\n", err)
		return
	}

	fmt.Println("\n✅ Key pressed! Continuing...")
	fmt.Println()
	fmt.Println("Here's another one with a custom message:")
	fmt.Println()

	err = ui.WaitForKey("Hit any key when you're ready to proceed →")
	if err != nil {
		fmt.Printf("\nError: %v\n", err)
		return
	}

	fmt.Println("\n✅ All done!")
}

func runPressKeyDemos() {
	demoPressKey()
}
