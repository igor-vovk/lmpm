package installer

import (
	"reflect"
	"testing"

	"github.com/hubblew/pim/internal/config"
)

func TestCreateStrategy(t *testing.T) {
	tests := []struct {
		name         string
		strategyType config.StrategyType
		outputPath   string
		expectError  bool
		expectedType reflect.Type
	}{
		{
			name:         "explicit concat strategy",
			strategyType: config.StrategyConcat,
			outputPath:   "./output.md",
			expectError:  false,
			expectedType: reflect.TypeOf(&ConcatStrategy{}),
		},
		{
			name:         "explicit flatten strategy",
			strategyType: config.StrategyFlatten,
			outputPath:   "./output",
			expectError:  false,
			expectedType: reflect.TypeOf(&FlattenStrategy{}),
		},
		{
			name:         "explicit preserve strategy",
			strategyType: config.StrategyPreserve,
			outputPath:   "./output",
			expectError:  false,
			expectedType: reflect.TypeOf(&PreserveStrategy{}),
		},
		{
			name:         "auto-detect concat from .md extension",
			strategyType: "",
			outputPath:   "./output.md",
			expectError:  false,
			expectedType: reflect.TypeOf(&ConcatStrategy{}),
		},
		{
			name:         "auto-detect concat from .txt extension",
			strategyType: "",
			outputPath:   "./output.txt",
			expectError:  false,
			expectedType: reflect.TypeOf(&ConcatStrategy{}),
		},
		{
			name:         "auto-detect flatten from directory path",
			strategyType: "",
			outputPath:   "./output",
			expectError:  false,
			expectedType: reflect.TypeOf(&FlattenStrategy{}),
		},
		{
			name:         "invalid strategy type",
			strategyType: config.StrategyType("invalid"),
			outputPath:   "./output",
			expectError:  true,
			expectedType: nil,
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

				actualType := reflect.TypeOf(strategy)

				if actualType != tt.expectedType {
					t.Errorf("expected strategy type %v, got %v", tt.expectedType, actualType)
				}
			}
		})
	}
}

func TestCreateStrategyOutputPaths(t *testing.T) {
	tests := []struct {
		name         string
		outputPath   string
		expectedType reflect.Type
	}{
		{"markdown file", "docs.md", reflect.TypeOf(&ConcatStrategy{})},
		{"text file", "output.txt", reflect.TypeOf(&ConcatStrategy{})},
		{"nested markdown", "./nested/path/docs.md", reflect.TypeOf(&ConcatStrategy{})},
		{"directory", "output/", reflect.TypeOf(&FlattenStrategy{})},
		{"nested directory", "./nested/output", reflect.TypeOf(&FlattenStrategy{})},
		{"yaml file", "config.yaml", reflect.TypeOf(&FlattenStrategy{})},
		{"json file", "data.json", reflect.TypeOf(&FlattenStrategy{})},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			strategy, err := NewStrategy("", tt.outputPath)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			actualType := reflect.TypeOf(strategy)
			if actualType != tt.expectedType {
				t.Errorf("expected strategy type %s, got %s", tt.expectedType, actualType)
			}
		})
	}
}
