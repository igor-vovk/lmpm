package installer

import (
	"fmt"

	"github.com/hubble-works/pim/internal/config"
)

// Strategy defines the interface for different installation strategies.
//
// Each strategy must implement methods to prepare the output, add files,
// and close any resources when done.
type Strategy interface {
	Prepare() error
	AddFile(srcPath, relativePath string) error
	Close() error
}

func CreateStrategy(
	strategyType config.Strategy,
	outputPath string,
) (Strategy, error) {
	switch strategyType {
	case config.StrategyConcat:
		return NewConcatStrategy(outputPath), nil
	case config.StrategyFlatten:
		return NewFlattenStrategy(outputPath), nil
	case config.StrategyPreserve:
		return NewPreserveStrategy(outputPath), nil
	}

	return nil, fmt.Errorf("unknown strategy type: %s", strategyType)
}
