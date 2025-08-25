package models

// Anthropic model definitions.
const AnthropicProvider = "anthropic"

// Anthropic creates an Anthropic model with the specified name and options.
func Anthropic(modelName string, options ...ModelOption) Model {
	return newModel(AnthropicProvider, modelName, options...)
}
