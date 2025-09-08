package providers

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/farbodsalimi/genevieve/pkg/genevieve"
	"github.com/farbodsalimi/genevieve/pkg/providers/anthropic"
	"github.com/farbodsalimi/genevieve/pkg/providers/google"
	"github.com/farbodsalimi/genevieve/pkg/providers/openai"
)

type Config struct {
	OpenAI    OpenAIConfig    `json:"openai"`
	Anthropic AnthropicConfig `json:"anthropic"`
	Google    GoogleConfig    `json:"google"`
}

type OpenAIConfig struct {
	APIKey string `json:"api_key"`
}

type AnthropicConfig struct {
	APIKey string `json:"api_key"`
}

type GoogleConfig struct {
	APIKey string `json:"api_key"`
}

func LoadConfig() (*Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	configPath := filepath.Join(homeDir, ".promptctl", "config.json")

	// Create config directory if it doesn't exist
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, err
	}

	// If config doesn't exist, create default
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		defaultConfig := &Config{}
		data, _ := json.MarshalIndent(defaultConfig, "", "  ")
		os.WriteFile(configPath, data, 0600)
		return defaultConfig, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func GetProvider(ctx context.Context) (*genevieve.Router, error) {
	config, err := LoadConfig()
	if err != nil {
		return nil, err
	}

	openaiClient := openai.NewClient(ctx, config.OpenAI.APIKey)
	anthropicClient := anthropic.NewClient(ctx, config.Anthropic.APIKey)
	geminiClient := google.NewClient(ctx, config.Google.APIKey)

	router := genevieve.NewRouter()
	router.Register(openaiClient)
	router.Register(anthropicClient)
	router.Register(geminiClient)

	return router, nil
}
