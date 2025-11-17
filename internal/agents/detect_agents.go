package agents

import (
	"os/exec"
)

// DetectAgentTools checks for known LLM CLI tools in the system
func DetectAgentTools() []AgentTool {
	var tools []AgentTool

	if path, ok := isCommandAvailable("copilot"); ok {
		cmd := exec.Command(path, "--version")
		if err := cmd.Run(); err == nil {
			tools = append(tools, NewGhCopilotAgent(path))
		}
	}
	if path, ok := isCommandAvailable("gemini"); ok {
		cmd := exec.Command(path, "--version")
		if err := cmd.Run(); err == nil {
			tools = append(tools, NewGeminiCLIAgent(path))
		}
	}

	tools = append(tools, NewManualAgent())

	return tools
}

// isCommandAvailable checks if a command is available in PATH
func isCommandAvailable(name string) (string, bool) {
	path, err := exec.LookPath(name)

	return path, err == nil
}
