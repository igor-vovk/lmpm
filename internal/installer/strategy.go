package installer

import (
	"fmt"

	"github.com/hubblew/pim/internal/config"
	"github.com/spf13/afero"
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
	fs afero.Fs,
	strategyType config.StrategyType,
	outputPath string,
) (Strategy, error) {
	switch strategyType {
	case config.StrategyConcat:
		return NewConcatStrategy(fs, outputPath), nil
	case config.StrategyFlatten:
		return NewFlattenStrategy(fs, outputPath), nil
	case config.StrategyPreserve:
		return NewPreserveStrategy(fs, outputPath), nil
	case "":
		if HasMdExtension(outputPath) {
			return NewStrategy(fs, config.StrategyConcat, outputPath)
		} else {
			return NewStrategy(fs, config.StrategyFlatten, outputPath)
		}
	}

	return nil, fmt.Errorf("unknown strategy type: %s", strategyType)
}
