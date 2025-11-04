package installer

import (
	"testing"

	"github.com/hubblew/pim/internal/config"
)

func TestCreateStrategy(t *testing.T) {
	tests := []struct {
		name         string
		strategyType config.StrategyType
		outputPath   string
		expectError  bool
		expectedType config.StrategyType
	}{
		{
			name:         "explicit concat strategy",
			strategyType: config.StrategyConcat,
			outputPath:   "./output.md",
			expectError:  false,
			expectedType: config.StrategyConcat,
		},
		{
			name:         "explicit flatten strategy",
			strategyType: config.StrategyFlatten,
			outputPath:   "./output",
			expectError:  false,
			expectedType: config.StrategyFlatten,
		},
		{
			name:         "explicit preserve strategy",
			strategyType: config.StrategyPreserve,
			outputPath:   "./output",
			expectError:  false,
			expectedType: config.StrategyPreserve,
		},
		{
			name:         "auto-detect concat from .md extension",
			strategyType: "",
			outputPath:   "./output.md",
			expectError:  false,
			expectedType: config.StrategyConcat,
		},
		{
			name:         "auto-detect concat from .txt extension",
			strategyType: "",
			outputPath:   "./output.txt",
			expectError:  false,
			expectedType: config.StrategyConcat,
		},
		{
			name:         "auto-detect flatten from directory path",
			strategyType: "",
			outputPath:   "./output",
			expectError:  false,
			expectedType: config.StrategyFlatten,
		},
		{
			name:         "auto-detect flatten from .yaml extension",
			strategyType: "",
			outputPath:   "./config.yaml",
			expectError:  false,
			expectedType: config.StrategyFlatten,
		},
		{
			name:         "auto-detect flatten from directory slash",
			strategyType: "",
			outputPath:   "./output/",
			expectError:  false,
			expectedType: config.StrategyFlatten,
		},
		{
			name:         "invalid strategy type",
			strategyType: config.StrategyType("invalid"),
			outputPath:   "./output",
			expectError:  true,
			expectedType: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			strategy, err := NewStrategy(tt.strategyType, tt.outputPath)

			if tt.expectError {
				if err == nil {
					t.Error("expected error but got none")
				}
				if strategy != nil {
					t.Error("expected strategy to be nil when error occurs")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if strategy == nil {
					t.Error("expected strategy to be non-nil")
				}
				if strategy.GetType() != tt.expectedType {
					t.Errorf("expected strategy type %s, got %s", tt.expectedType, strategy.GetType())
				}
			}
		})
	}
}

func TestCreateStrategyRecursion(t *testing.T) {
	// Test that auto-detection recursively calls NewStrategy
	strategy, err := NewStrategy("", "test.md")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, ok := strategy.(*ConcatStrategy)
	if !ok {
		t.Errorf("expected ConcatStrategy, got %T", strategy)
	}
}

func TestCreateStrategyOutputPaths(t *testing.T) {
	tests := []struct {
		name       string
		outputPath string
		strategy   config.StrategyType
	}{
		{"markdown file", "docs.md", config.StrategyConcat},
		{"text file", "output.txt", config.StrategyConcat},
		{"nested markdown", "./nested/path/docs.md", config.StrategyConcat},
		{"directory", "output/", config.StrategyFlatten},
		{"nested directory", "./nested/output", config.StrategyFlatten},
		{"yaml file", "config.yaml", config.StrategyFlatten},
		{"json file", "data.json", config.StrategyFlatten},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			strategy, err := NewStrategy("", tt.outputPath)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if strategy.GetType() != tt.strategy {
				t.Errorf("expected strategy type %s, got %s", tt.strategy, strategy.GetType())
			}
		})
	}
}
