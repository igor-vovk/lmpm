package installer

import (
	"fmt"

	"github.com/hubblew/pim/internal/config"
)

// Strategy defines the interface for different installation strategies.
//
// Each strategy must implement methods to Initialize the output, AddFile to add files,
// and Close any resources when done.
type Strategy interface {
	Initialize(prompter UserPrompter) error
	AddFile(srcPath, relativePath string) error
	Close() error
}

func NewStrategy(
	strategyType config.StrategyType,
	outputPath string,
) (Strategy, error) {
	switch strategyType {
	case config.StrategyConcat:
		return NewConcatStrategy(outputPath), nil
	case config.StrategyFlatten:
		return NewFlattenStrategy(outputPath), nil
	case config.StrategyPreserve:
		return NewPreserveStrategy(outputPath), nil
	case "":
		if HasMdExtension(outputPath) {
			return NewStrategy(config.StrategyConcat, outputPath)
		} else {
			return NewStrategy(config.StrategyFlatten, outputPath)
		}
	}

	return nil, fmt.Errorf("unknown strategy type: %s", strategyType)
}
