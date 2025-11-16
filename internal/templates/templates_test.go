package templates

import (
	"strings"
	"testing"
)

func TestRenderGenerateInstructionsPrompt(t *testing.T) {
	tpl, err := RenderGenerateInstructionsPrompt("/path/to/instructions")

	if (strings.Contains(tpl, "/path/to/instructions")) == false {
		t.Fatalf("Expected instructionsDir to be included in the template output")
	}
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}
