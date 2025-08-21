package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// Tool is the interface that all tools must implement.
type Tool interface {
	GetName() string
	GetDescription() string
	GetSchema() (*ToolSchema, error)
	Execute(ctx context.Context, params json.RawMessage) (interface{}, error)
}

// ToolInfo contains basic information about a tool.
type ToolInfo struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Schema      *ToolSchema `json:"schema,omitempty"`
}

// tool is a type-safe tool definition using generics.
type tool[P any] struct {
	Name        string
	Description string
	Handler     func(context.Context, P) (interface{}, error)

	// Cache for reflection results
	schema *ToolSchema
}

// ToolSchema represents the JSON schema for tool parameters (OpenAI/Anthropic compatible).
type ToolSchema struct {
	Type       string              `json:"type"`
	Properties map[string]Property `json:"properties"`
	Required   []string            `json:"required,omitempty"`
}

// Property represents a parameter property in the schema.
type Property struct {
	Type        string      `json:"type"`
	Description string      `json:"description,omitempty"`
	Enum        []string    `json:"enum,omitempty"`
	Default     interface{} `json:"default,omitempty"`
	Minimum     *float64    `json:"minimum,omitempty"`
	Maximum     *float64    `json:"maximum,omitempty"`
	Items       *Property   `json:"items,omitempty"` // For array types
}

// Validator is an optional interface for parameter validation.
type Validator interface {
	Validate() error
}

// GetName implements Tool interface.
func (t tool[P]) GetName() string {
	return t.Name
}

// GetDescription implements Tool interface.
func (t tool[P]) GetDescription() string {
	return t.Description
}

// GetSchema generates JSON Schema from Go struct tags.
func (t tool[P]) GetSchema() (*ToolSchema, error) {
	if t.schema != nil {
		return t.schema, nil
	}

	var p P
	typ := reflect.TypeOf(p)

	// Handle pointer types
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	if typ.Kind() != reflect.Struct {
		return nil, fmt.Errorf("tool parameters must be a struct, got %v", typ.Kind())
	}

	schema := &ToolSchema{
		Type:       "object",
		Properties: make(map[string]Property),
		Required:   []string{},
	}

	// Parse struct fields using reflection
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)

		// Skip unexported fields
		if !field.IsExported() {
			continue
		}

		// Parse JSON tag
		jsonTag := field.Tag.Get("json")
		if jsonTag == "-" {
			continue
		}

		fieldName, omitempty := parseJSONTag(jsonTag)
		if fieldName == "" {
			fieldName = field.Name
		}

		prop := Property{
			Type:        goTypeToJSONType(field.Type),
			Description: field.Tag.Get("desc"),
		}

		// Parse enum tag
		if enumTag := field.Tag.Get("enum"); enumTag != "" {
			prop.Enum = strings.Split(enumTag, ",")
		}

		// Parse default tag
		if defaultTag := field.Tag.Get("default"); defaultTag != "" {
			prop.Default = parseDefaultValue(defaultTag, field.Type)
		}

		// Parse min tag
		if minTag := field.Tag.Get("min"); minTag != "" {
			if v, err := strconv.ParseFloat(minTag, 64); err == nil {
				prop.Minimum = &v
			}
		}

		// Parse max tag
		if maxTag := field.Tag.Get("max"); maxTag != "" {
			if v, err := strconv.ParseFloat(maxTag, 64); err == nil {
				prop.Maximum = &v
			}
		}

		// Handle array/slice types
		if field.Type.Kind() == reflect.Slice || field.Type.Kind() == reflect.Array {
			elemType := field.Type.Elem()
			itemProp := &Property{
				Type: goTypeToJSONType(elemType),
			}
			prop.Items = itemProp
		}

		schema.Properties[fieldName] = prop

		// Determine if field is required
		requiredTag := field.Tag.Get("required")
		if requiredTag == "true" || (requiredTag != "false" && !omitempty) {
			schema.Required = append(schema.Required, fieldName)
		}
	}

	t.schema = schema
	return t.schema, nil
}

// Execute implements Tool interface with type safety.
func (t tool[P]) Execute(ctx context.Context, params json.RawMessage) (interface{}, error) {
	var p P
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, fmt.Errorf("invalid parameters for tool %s: %w", t.Name, err)
	}

	// Optional validation
	if validator, ok := interface{}(p).(Validator); ok {
		if err := validator.Validate(); err != nil {
			return nil, fmt.Errorf("parameter validation failed: %w", err)
		}
	}

	return t.Handler(ctx, p)
}

// parseJSONTag parses the json struct tag.
func parseJSONTag(tag string) (name string, omitempty bool) {
	if tag == "" {
		return "", false
	}
	parts := strings.Split(tag, ",")
	name = parts[0]
	for i := 1; i < len(parts); i++ {
		if parts[i] == "omitempty" {
			omitempty = true
			break
		}
	}
	return name, omitempty
}

// goTypeToJSONType converts Go types to JSON schema types.
func goTypeToJSONType(t reflect.Type) string {
	switch t.Kind() {
	case reflect.String:
		return "string"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return "integer"
	case reflect.Float32, reflect.Float64:
		return "number"
	case reflect.Bool:
		return "boolean"
	case reflect.Slice, reflect.Array:
		return "array"
	case reflect.Struct:
		// Check for time.Time
		if t.String() == "time.Time" {
			return "string" // ISO 8601 date-time string
		}
		return "object"
	case reflect.Map:
		return "object"
	case reflect.Ptr:
		return goTypeToJSONType(t.Elem())
	default:
		return "string" // Default fallback
	}
}

// parseDefaultValue parses the default value based on type.
func parseDefaultValue(value string, t reflect.Type) interface{} {
	switch t.Kind() {
	case reflect.String:
		return value
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if v, err := strconv.ParseInt(value, 10, 64); err == nil {
			return v
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if v, err := strconv.ParseUint(value, 10, 64); err == nil {
			return v
		}
	case reflect.Float32, reflect.Float64:
		if v, err := strconv.ParseFloat(value, 64); err == nil {
			return v
		}
	case reflect.Bool:
		if v, err := strconv.ParseBool(value); err == nil {
			return v
		}
	}
	return value
}
