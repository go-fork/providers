// Package model defines the core data structures used in the configuration system.
//
// This package contains the fundamental types for storing configuration values and
// their types in a type-safe manner. The main types are ConfigValue which represents
// a single configuration value with its type information, ValueType which is an enum
// of possible value types, and ConfigMap which is a flattened map of configuration keys
// to their values.
package model

// ValueType represents the type of a configuration value.
// This is used to provide type-safety when retrieving values from the configuration.
type ValueType int

const (
	// TypeUnknown represents an unknown or uninitialized value type.
	TypeUnknown ValueType = iota
	// TypeString represents a string value.
	TypeString
	// TypeInt represents an integer value.
	TypeInt
	// TypeFloat represents a floating-point value.
	TypeFloat
	// TypeBool represents a boolean value.
	TypeBool
	// TypeSlice represents a slice/array value.
	TypeSlice
	// TypeMap represents a map/object value.
	TypeMap
	// TypeNil represents a nil/null value.
	TypeNil
)

// String returns a string representation of the ValueType.
func (vt ValueType) String() string {
	switch vt {
	case TypeString:
		return "string"
	case TypeInt:
		return "int"
	case TypeFloat:
		return "float"
	case TypeBool:
		return "bool"
	case TypeSlice:
		return "slice"
	case TypeMap:
		return "map"
	case TypeNil:
		return "nil"
	default:
		return "unknown"
	}
}

// ConfigValue represents a single configuration value with its type information.
// It stores both the original value and its type to allow for type-safe access.
type ConfigValue struct {
	// Value is the original value (string, int, bool, slice, map)
	Value interface{}
	// Type is the type of the value, used for type-safe access
	Type ValueType
}

// ConfigMap is a flattened map of configuration keys to their values.
// The keys use dot notation to represent hierarchical structures.
type ConfigMap map[string]ConfigValue
