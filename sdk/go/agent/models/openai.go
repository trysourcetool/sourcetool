package models

// OpenAI model definitions.
const OpenAIProvider = "openai"

// OpenAI creates an OpenAI model with the specified name and options.
func OpenAI(modelName string, options ...ModelOption) Model {
	return newModel(OpenAIProvider, modelName, options...)
}
