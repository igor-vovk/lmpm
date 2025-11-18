package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/hubblew/pim/internal/config"
	"github.com/hubblew/pim/internal/templates"
	"github.com/hubblew/pim/internal/tpagents"
	"github.com/hubblew/pim/internal/ui"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new PIM configuration",
	Long:  "Detect LLM tools, discover existing instruction files, and create a pim.yaml configuration.",
	RunE:  runInit,
}

func init() {
	alphaCmd.AddCommand(initCmd)
}

func runInit(_ *cobra.Command, _ []string) error {
	reader := bufio.NewReader(os.Stdin)

	// Step 1: Detect LLM agents
	var tools []tpagents.TPAgentTool
	err := ui.RunWithSpinner("Detecting CLI agents in your system...", func() error {
		tools = tpagents.DetectTPAgentTools()
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to detect agents: %w", err)
	}

	if len(tools) == 0 {
		fmt.Println("No CLI agents detected in your system.")
		fmt.Println("Currently supported tools: GitHub Copilot, Google Gemini CLI.")
		return nil
	}

	// Step 2: Ask user to choose a tool
	var selectedTool tpagents.TPAgentTool
	if len(tools) == 1 {
		selectedTool = tools[0]
		fmt.Printf("\nUsing detected tool: %s\n", selectedTool.Descriptor())
	} else {
		choices := make([]ui.Choice, len(tools))
		for i, tool := range tools {
			choices[i] = ui.Choice{Label: tool.Descriptor(), Value: tool}
		}

		choice, err := ui.NewChoiceDialog("\nSelect an agent:", choices).Vertical().Run()
		if err != nil {
			return fmt.Errorf("failed to run selection dialog: %w", err)
		}
		if choice == nil {
			return fmt.Errorf("no agent selected")
		}
		selectedTool = choice.Value.(tpagents.TPAgentTool)
	}

	// Step 3: Ask for configuration file name
	fmt.Print("\nEnter configuration file name (default: pim.yaml): ")
	configName, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}
	configName = strings.TrimSpace(configName)
	if configName == "" {
		configName = "pim.yaml"
	}

	// Step 4: Check if config file already exists
	if _, err := os.Stat(configName); err == nil {
		return fmt.Errorf("configuration file '%s' already exists", configName)
	}

	// Step 5: Ask for instructions folder name
	fmt.Print("\nEnter instructions folder name (default: ./instructions): ")
	instructionsDir, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}
	instructionsDir = strings.TrimSpace(instructionsDir)
	if instructionsDir == "" {
		instructionsDir = "./instructions"
	}

	// Step 6: Create instructions directory if it doesn't exist
	if err := os.MkdirAll(instructionsDir, 0755); err != nil {
		return fmt.Errorf("failed to create instructions directory: %w", err)
	}
	fmt.Printf("Instructions directory: %s\n", instructionsDir)

	// Step 7: Ask if user wants to generate instructions using AI
	fmt.Printf("\nDo you want to generate instruction files using %s?\n", selectedTool.Descriptor())

	choice, err := ui.NewChoiceDialog("", ui.ChoicesYesNo()).Run()
	if err != nil {
		return fmt.Errorf("failed to run generation choice dialog: %w", err)
	}
	if choice != nil && choice.Value.(bool) {
		if err := generateInstructions(selectedTool, instructionsDir); err != nil {
			fmt.Printf("Warning: failed to generate instructions: %v\n", err)
		}
	}

	// Step 8: Look for existing instruction files
	existingFiles := discoverInstructionFiles(instructionsDir)

	// Step 9: Generate pim.yaml based on detected tool
	cfg, err := generateConfig(selectedTool, instructionsDir, existingFiles)
	if err != nil {
		return fmt.Errorf("failed to generate config: %w", err)
	}

	// Step 10: Write config to file
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configName, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	fmt.Printf("\n✓ Configuration file '%s' created successfully!\n", configName)
	fmt.Printf("✓ Instructions directory: %s\n", instructionsDir)
	if len(existingFiles) > 0 {
		fmt.Printf("✓ Found %d existing instruction file(s)\n", len(existingFiles))
	}
	fmt.Println("\nNext steps:")
	fmt.Printf("  1. Review and edit %s\n", configName)
	fmt.Printf("  2. Add your instruction files to %s\n", instructionsDir)
	fmt.Printf("  3. Run 'pim install' to apply the configuration\n")

	return nil
}

// discoverInstructionFiles looks for existing instruction files
func discoverInstructionFiles(instructionsDir string) []string {
	var files []string

	// Common instruction file locations
	candidates := []string{
		filepath.Join(instructionsDir, "*.md"),
		"AGENTS.md",
		".github/copilot-instructions.md",
	}

	for _, pattern := range candidates {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			continue
		}
		files = append(files, matches...)
	}

	// List files in instructions directory if it exists
	if entries, err := os.ReadDir(instructionsDir); err == nil {
		for _, entry := range entries {
			if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".md") {
				filePath := filepath.Join(instructionsDir, entry.Name())
				// Avoid duplicates
				found := false
				for _, f := range files {
					if f == filePath {
						found = true
						break
					}
				}
				if !found {
					files = append(files, filePath)
				}
			}
		}
	}

	return files
}

// generateConfig creates a config based on the selected tool and discovered files
func generateConfig(tool tpagents.TPAgentTool, instructionsDir string, existingFiles []string) (*config.Config, error) {
	cfg := config.NewConfig()

	switch reflect.TypeOf(tool) {
	case tpagents.GhCopilotAgentType:
		target := config.Target{
			Name:    "copilot-instructions",
			Output:  ".github/copilot-instructions.md",
			Include: []string{},
		}

		// Add existing files to include list
		for _, file := range existingFiles {
			// Skip the output file itself
			if file == ".github/copilot-instructions.md" {
				continue
			}
			target.Include = append(target.Include, file)
		}

		cfg.Targets = []config.Target{target}

	case tpagents.GeminiCLIAgentType:
		target := config.Target{
			Name:    "gemini-instructions",
			Output:  "GEMINI.md",
			Include: []string{},
		}

		// Add existing files to include list
		for _, file := range existingFiles {
			// Skip the output file itself
			if file == "GEMINI.md" {
				continue
			}
			target.Include = append(target.Include, file)
		}

		cfg.Targets = []config.Target{target}

	case tpagents.ManualAgentType:
		target := config.Target{
			Name:    "manual-instructions",
			Output:  "AGENTS.md",
			Include: []string{},
		}

		// Add existing files to include list
		for _, file := range existingFiles {
			// Skip the output file itself
			if file == "AGENTS.md" {
				continue
			}
			target.Include = append(target.Include, file)
		}

		cfg.Targets = []config.Target{target}

	default:
		return nil, fmt.Errorf("unsupported tool: %s", tool.Descriptor())
	}

	return cfg, nil
}

// generateInstructions uses the agent tool to generate instruction files
func generateInstructions(tool tpagents.TPAgentTool, instructionsDir string) error {
	fmt.Printf("\nGenerating instruction files using %s (this may take a while)...\n\n", tool.Descriptor())

	prompt, err := templates.RenderGenerateInstructionsPrompt(instructionsDir)
	if err != nil {
		return fmt.Errorf("failed to render instructions template: %w", err)
	}

	_, err = tool.ExecuteCommand(prompt)
	if err != nil {
		return fmt.Errorf("failed to execute agent command: %w", err)
	}

	fmt.Println("\nYou can now manually create these instruction files in the", instructionsDir, "directory.")

	return nil
}
