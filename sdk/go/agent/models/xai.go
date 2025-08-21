package models

// xAI model definitions.
const XAIProvider = "xai"

// XAI creates an xAI model with the specified name and options.
func XAI(modelName string, options ...ModelOption) Model {
	return newModel(XAIProvider, modelName, options...)
}
