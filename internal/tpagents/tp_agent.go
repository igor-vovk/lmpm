package tpagents

type TPAgentTool interface {
	Descriptor() string
	ExecuteCommand(command string) (string, error)
}
