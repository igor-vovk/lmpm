package agents

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"reflect"

	"github.com/hubblew/pim/internal/utils"
)

type GeminiCLIAgent struct {
	exec string
}

var GeminiCLIAgentType = reflect.TypeOf(new(GeminiCLIAgent))

var _ AgentTool = (*GeminiCLIAgent)(nil)

func NewGeminiCLIAgent(path string) *GeminiCLIAgent {
	return &GeminiCLIAgent{
		exec: path,
	}
}

func (a *GeminiCLIAgent) Descriptor() string {
	return "Gemini CLI (" + a.exec + ")"
}

func (a *GeminiCLIAgent) ExecuteCommand(command string) (string, error) {
	cmd := exec.Command(a.exec, "--approval-mode=yolo", fmt.Sprintf(`"%s"`, command))

	var buf bytes.Buffer
	prefix := "  Gemini> "
	cmd.Stdout = utils.NewPrefixWriter(os.Stdout, prefix)
	cmd.Stderr = utils.NewPrefixWriter(os.Stderr, prefix)

	err := cmd.Run()

	return buf.String(), err
}
