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

func TestLevelComparison(t *testing.T) {
	// Kiểm tra các level có thứ tự đúng
	if !(DebugLevel < InfoLevel) {
		t.Error("DebugLevel nên nhỏ hơn InfoLevel")
	}
	if !(InfoLevel < WarningLevel) {
		t.Error("InfoLevel nên nhỏ hơn WarningLevel")
	}
	if !(WarningLevel < ErrorLevel) {
		t.Error("WarningLevel nên nhỏ hơn ErrorLevel")
	}
	if !(ErrorLevel < FatalLevel) {
		t.Error("ErrorLevel nên nhỏ hơn FatalLevel")
	}

	// Kiểm tra cấp độ tùy chỉnh
	customLevel := Level(42)
	if customLevel.String() != "UNKNOWN" {
		t.Errorf("Level tùy chỉnh nên trả về 'UNKNOWN', got %v", customLevel.String())
	}
}

// TestLevelStringWithValues kiểm tra giá trị số của các level
func TestLevelStringWithValues(t *testing.T) {
	tests := []struct {
		name  string
		level Level
		value int
	}{
		{"Debug", DebugLevel, 0},
		{"Info", InfoLevel, 1},
		{"Warning", WarningLevel, 2},
		{"Error", ErrorLevel, 3},
		{"Fatal", FatalLevel, 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if int(tt.level) != tt.value {
				t.Errorf("Level %v nên có giá trị = %v, got = %v",
					tt.name, tt.value, int(tt.level))
			}
		})
	}
}
