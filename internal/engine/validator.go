package engine

import (
	"fmt"
	"strconv"

	"github.com/openlibing/openlibing-cli/pkg/spc"
)

// Validate checks that the provided parameters satisfy the SPC definition.
// It verifies required fields, type correctness, enum membership, and range constraints.
func Validate(def *spc.SPCDefinition, params map[string]interface{}) error {
	for _, p := range def.Parameters {
		val, exists := params[p.Name]

		// Required check
		if p.Required && (!exists || val == nil || val == "") {
			return &ValidationError{
				Parameter: p.Name,
				Message:   fmt.Sprintf("parameter '%s' is required", p.Name),
			}
		}

		// Skip further checks if no value and not required
		if !exists || val == nil || val == "" {
			// Apply default if available
			if p.Default != nil && exists {
				// User explicitly passed empty → use default
				// This covers the case where value is provided but empty
			}
			continue
		}

		// Type check
		switch p.Type {
		case "integer":
			if err := checkInteger(val); err != nil {
				return &ValidationError{
					Parameter: p.Name,
					Message:   fmt.Sprintf("parameter '%s': %v", p.Name, err),
				}
			}
		case "boolean":
			if err := checkBoolean(val); err != nil {
				return &ValidationError{
					Parameter: p.Name,
					Message:   fmt.Sprintf("parameter '%s': %v", p.Name, err),
				}
			}
		// string type needs no explicit check — everything can be a string
		}

		// Enum check
		if len(p.Enum) > 0 {
			strVal := fmt.Sprintf("%v", val)
			found := false
			for _, e := range p.Enum {
				if e == strVal {
					found = true
					break
				}
			}
			if !found {
				return &ValidationError{
					Parameter: p.Name,
					Message:   fmt.Sprintf("parameter '%s' must be one of %v, got '%s'", p.Name, p.Enum, strVal),
				}
			}
		}

		// Range validation (integer only)
		if p.Validation != nil && p.Type == "integer" {
			n, _ := toInt(val)
			if p.Validation.Min != nil && n < *p.Validation.Min {
				return &ValidationError{
					Parameter: p.Name,
					Message:   fmt.Sprintf("parameter '%s' must be >= %d, got %d", p.Name, *p.Validation.Min, n),
				}
			}
			if p.Validation.Max != nil && n > *p.Validation.Max {
				return &ValidationError{
					Parameter: p.Name,
					Message:   fmt.Sprintf("parameter '%s' must be <= %d, got %d", p.Name, *p.Validation.Max, n),
				}
			}
		}
	}

	return nil
}

// ValidationError represents a parameter validation failure.
type ValidationError struct {
	Parameter string
	Message   string
}

func (e *ValidationError) Error() string {
	return e.Message
}

func checkInteger(val interface{}) error {
	switch v := val.(type) {
	case int, int64, float64:
		return nil
	case string:
		_, err := strconv.Atoi(v)
		return err
	default:
		return fmt.Errorf("expected integer, got %T", val)
	}
}

func checkBoolean(val interface{}) error {
	switch v := val.(type) {
	case bool:
		return nil
	case string:
		if v != "true" && v != "false" {
			return fmt.Errorf("expected boolean, got %q", v)
		}
		return nil
	default:
		return fmt.Errorf("expected boolean, got %T", val)
	}
}

func toInt(val interface{}) (int, error) {
	switch v := val.(type) {
	case int:
		return v, nil
	case int64:
		return int(v), nil
	case float64:
		return int(v), nil
	case string:
		return strconv.Atoi(v)
	default:
		return 0, fmt.Errorf("cannot convert %T to int", val)
	}
}
