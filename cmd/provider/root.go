package provider

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/farbodsalimi/promptctl/internal/providers"
)

var providerAddCmd = &cobra.Command{
	Use:   "add <name> <api_key>",
	Short: "Add or update a provider API key",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		providerName := args[0]
		apiKey := args[1]

		if err := updateProviderConfig(providerName, apiKey); err != nil {
			log.Fatalf("failed to add provider: %v", err)
		}

		fmt.Printf("Added/updated provider: %s\n", providerName)
	},
}

var providerListCmd = &cobra.Command{
	Use:   "list",
	Short: "List configured providers",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := providers.LoadConfig()
		if err != nil {
			log.Fatalf("failed to load config: %v", err)
		}

		fmt.Println("Configured providers:")
		if config.OpenAI.APIKey != "" {
			fmt.Printf("  openai: %s\n", maskAPIKey(config.OpenAI.APIKey))
		}
		if config.Anthropic.APIKey != "" {
			fmt.Printf("  anthropic: %s\n", maskAPIKey(config.Anthropic.APIKey))
		}
		if config.Google.APIKey != "" {
			fmt.Printf("  google: %s\n", maskAPIKey(config.Google.APIKey))
		}
	},
}

var providerDeleteCmd = &cobra.Command{
	Use:   "delete <name>",
	Short: "Delete a provider configuration",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		providerName := args[0]

		if err := updateProviderConfig(providerName, ""); err != nil {
			log.Fatalf("failed to delete provider: %v", err)
		}

		fmt.Printf("Deleted provider: %s\n", providerName)
	},
}

func updateProviderConfig(providerName, apiKey string) error {
	config, err := providers.LoadConfig()
	if err != nil {
		return err
	}

	switch providerName {
	case "openai":
		config.OpenAI.APIKey = apiKey
	case "anthropic":
		config.Anthropic.APIKey = apiKey
	case "google":
		config.Google.APIKey = apiKey
	default:
		return fmt.Errorf(
			"unsupported provider: %s (supported: openai, anthropic, google)",
			providerName,
		)
	}

	return saveConfig(config)
}

func saveConfig(config *providers.Config) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configPath := filepath.Join(homeDir, ".promptctl", "config.json")
	configDir := filepath.Dir(configPath)

	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0600)
}

func maskAPIKey(apiKey string) string {
	if len(apiKey) <= 8 {
		return "***"
	}
	return apiKey[:4] + "..." + apiKey[len(apiKey)-4:]
}

func NewRootCmd() *cobra.Command {
	var providerCmd = &cobra.Command{
		Use:   "provider",
		Short: "Manage LLM providers",
		Long:  `Add, update, list, and delete LLM provider configurations.`,
	}

	providerCmd.AddCommand(providerAddCmd)
	providerCmd.AddCommand(providerListCmd)
	providerCmd.AddCommand(providerDeleteCmd)

	return providerCmd
}
