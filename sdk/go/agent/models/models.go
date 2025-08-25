package models

import "fmt"

// Model represents an AI model with provider and configuration.
type Model interface {
	Provider() string
	Name() string
	ID() string // Full model identifier (provider/model)
	Config() *Config
}

// Config holds model configuration options.
type Config struct {
	Temperature *float64 `json:"temperature,omitempty"`
	MaxTokens   int      `json:"max_tokens,omitempty"`
	TopP        *float64 `json:"top_p,omitempty"`
	TopK        *int     `json:"top_k,omitempty"`
	Stop        []string `json:"stop,omitempty"`
	Seed        *int     `json:"seed,omitempty"`
}

// baseModel provides common functionality for all models.
type baseModel struct {
	provider string
	name     string
	config   *Config
}

func (m *baseModel) Provider() string {
	return m.provider
}

func (m *baseModel) Name() string {
	return m.name
}

func (m *baseModel) ID() string {
	return fmt.Sprintf("%s/%s", m.provider, m.name)
}

func (m *baseModel) Config() *Config {
	if m.config == nil {
		return &Config{}
	}
	return m.config
}

// ModelOption configures a model.
type ModelOption func(*Config)

// WithTemperature sets the temperature for the model.
func WithTemperature(temp float64) ModelOption {
	return func(c *Config) {
		c.Temperature = &temp
	}
}

// WithMaxTokens sets the maximum tokens for the model.
func WithMaxTokens(tokens int) ModelOption {
	return func(c *Config) {
		c.MaxTokens = tokens
	}
}

// WithTopP sets the top-p sampling parameter.
func WithTopP(p float64) ModelOption {
	return func(c *Config) {
		c.TopP = &p
	}
}

// WithTopK sets the top-k sampling parameter.
func WithTopK(k int) ModelOption {
	return func(c *Config) {
		c.TopK = &k
	}
}

// WithStop sets stop sequences for the model.
func WithStop(stop ...string) ModelOption {
	return func(c *Config) {
		c.Stop = stop
	}
}

// WithSeed sets a seed for deterministic generation.
func WithSeed(seed int) ModelOption {
	return func(c *Config) {
		c.Seed = &seed
	}
}

// newModel creates a new model with the given provider, name, and options.
func newModel(provider, name string, options ...ModelOption) Model {
	config := &Config{}
	for _, opt := range options {
		opt(config)
	}

	return &baseModel{
		provider: provider,
		name:     name,
		config:   config,
	}
}
