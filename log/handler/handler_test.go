package handler

import (
	"testing"
)

func TestLevelString(t *testing.T) {
	tests := []struct {
		name  string
		level Level
		want  string
	}{
		{"Debug", DebugLevel, "DEBUG"},
		{"Info", InfoLevel, "INFO"},
		{"Warning", WarningLevel, "WARNING"},
		{"Error", ErrorLevel, "ERROR"},
		{"Fatal", FatalLevel, "FATAL"},
		{"Unknown", Level(99), "UNKNOWN"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.level.String(); got != tt.want {
				t.Errorf("Level.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
