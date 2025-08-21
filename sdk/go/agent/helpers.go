package agent

import (
	"context"
	"fmt"
	"reflect"
)

// NewTool creates a new tool with type inference for cleaner syntax.
func NewTool[P any](name, description string, handler func(context.Context, P) (interface{}, error)) Tool {
	return tool[P]{
		Name:        name,
		Description: description,
		Handler:     handler,
	}
}

// StringTool creates a simple tool that accepts a string parameter.
func StringTool(name, description, paramName, paramDesc string, handler func(context.Context, string) (interface{}, error)) Tool {
	type stringParam struct {
		Value string `json:"value" desc:"Parameter value" required:"true"`
	}

	return tool[stringParam]{
		Name:        name,
		Description: description,
		Handler: func(ctx context.Context, param stringParam) (interface{}, error) {
			return handler(ctx, param.Value)
		},
	}
}

// SimpleTool creates a tool without parameters.
func SimpleTool(name, description string, handler func(context.Context) (interface{}, error)) Tool {
	type emptyParam struct{}

	return tool[emptyParam]{
		Name:        name,
		Description: description,
		Handler: func(ctx context.Context, param emptyParam) (interface{}, error) {
			return handler(ctx)
		},
	}
}

// Parameter helper functions for easier tool parameter definition

// StringParam creates a string parameter definition helper.
func StringParam(name, description string) ParamHelper {
	return ParamHelper{
		Name:        name,
		Type:        "string",
		Description: description,
	}
}

// NumberParam creates a number parameter definition helper.
func NumberParam(name, description string) ParamHelper {
	return ParamHelper{
		Name:        name,
		Type:        "number",
		Description: description,
	}
}

// BoolParam creates a boolean parameter definition helper.
func BoolParam(name, description string) ParamHelper {
	return ParamHelper{
		Name:        name,
		Type:        "boolean",
		Description: description,
	}
}

// ArrayParam creates an array parameter definition helper.
func ArrayParam(name, description, itemType string) ParamHelper {
	return ParamHelper{
		Name:        name,
		Type:        "array",
		Description: description,
		Items:       &Property{Type: itemType},
	}
}

// ParamHelper assists in building parameter definitions.
type ParamHelper struct {
	Name        string
	Type        string
	Description string
	Enum        []string
	Default     interface{}
	Min         *float64
	Max         *float64
	Items       *Property
	isRequired  bool
}

// WithEnum adds enum values to the parameter.
func (p ParamHelper) WithEnum(values ...string) ParamHelper {
	p.Enum = values
	return p
}

// WithDefault sets a default value for the parameter.
func (p ParamHelper) WithDefault(value interface{}) ParamHelper {
	p.Default = value
	return p
}

// WithMin sets a minimum value (for numeric parameters).
func (p ParamHelper) WithMin(min float64) ParamHelper {
	p.Min = &min
	return p
}

// WithMax sets a maximum value (for numeric parameters).
func (p ParamHelper) WithMax(max float64) ParamHelper {
	p.Max = &max
	return p
}

// Required marks the parameter as required.
func (p ParamHelper) Required() ParamHelper {
	p.isRequired = true
	return p
}

// ToProperty converts the helper to a Property.
func (p ParamHelper) ToProperty() Property {
	return Property{
		Type:        p.Type,
		Description: p.Description,
		Enum:        p.Enum,
		Default:     p.Default,
		Minimum:     p.Min,
		Maximum:     p.Max,
		Items:       p.Items,
	}
}

// Schema validation helpers

// IsValidSchema checks if a tool schema is valid.
func IsValidSchema(schema *ToolSchema) error {
	if schema == nil {
		return fmt.Errorf("schema cannot be nil")
	}
	if schema.Type != "object" {
		return fmt.Errorf("schema type must be 'object', got '%s'", schema.Type)
	}
	return nil
}

// Type checking utilities

// IsStructType checks if a type is a struct.
func IsStructType(t reflect.Type) bool {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t.Kind() == reflect.Struct
}

// IsSliceType checks if a type is a slice or array.
func IsSliceType(t reflect.Type) bool {
	return t.Kind() == reflect.Slice || t.Kind() == reflect.Array
}

// GetFieldTag extracts a specific tag from a struct field.
func GetFieldTag(field reflect.StructField, tagName string) string {
	return field.Tag.Get(tagName)
}

// Utility functions for common patterns

// MustNewTool creates a tool and panics on error (for initialization).
func MustNewTool[P any](name, description string, handler func(context.Context, P) (interface{}, error)) Tool {
	tool := NewTool(name, description, handler)
	if _, err := tool.GetSchema(); err != nil {
		panic(fmt.Sprintf("failed to create tool %s: %v", name, err))
	}
	return tool
}

// ToolFromFunc converts a simple function to a tool (experimental).
func ToolFromFunc(name, description string, fn interface{}) (Tool, error) {
	fnType := reflect.TypeOf(fn)
	if fnType.Kind() != reflect.Func {
		return nil, fmt.Errorf("expected function, got %v", fnType.Kind())
	}

	// Validate function signature
	if fnType.NumIn() != 2 {
		return nil, fmt.Errorf("function must have exactly 2 parameters (context.Context, params)")
	}

	if fnType.In(0) != reflect.TypeOf((*context.Context)(nil)).Elem() {
		return nil, fmt.Errorf("first parameter must be context.Context")
	}

	if fnType.NumOut() != 2 {
		return nil, fmt.Errorf("function must return (interface{}, error)")
	}

	// This would require more complex reflection to implement fully
	// For now, return an error suggesting to use the typed approach
	return nil, fmt.Errorf("automatic function conversion not implemented, use NewTool[T] instead")
}
