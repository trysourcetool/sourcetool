package models

// Google model definitions.
const GoogleProvider = "google"

// Google creates a Google Gemini model with the specified name and options.
func Google(modelName string, options ...ModelOption) Model {
	return newModel(GoogleProvider, modelName, options...)
}
