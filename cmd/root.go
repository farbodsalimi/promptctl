package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/farbodsalimi/promptctl/cmd/prompt"
	"github.com/farbodsalimi/promptctl/cmd/provider"
	"github.com/farbodsalimi/promptctl/cmd/run"
	"github.com/farbodsalimi/promptctl/cmd/vault"
	"github.com/farbodsalimi/promptctl/internal/db"
)

func NewRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "promptctl",
		Short: "A CLI tool for managing prompt vaults",
		Long:  `promptctl is a CLI tool for storing, versioning, and running prompts with various LLM providers.`,
	}

	rootCmd.AddCommand(prompt.NewRootCmd())
	rootCmd.AddCommand(provider.NewRootCmd())
	rootCmd.AddCommand(run.NewRootCmd())
	rootCmd.AddCommand(vault.NewRootCmd())

	return rootCmd
}

func Execute(version string) error {
	rootCmd := NewRootCommand()
	rootCmd.Version = version

	cobra.OnInitialize(initConfig)
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	return rootCmd.Execute()
}

func initConfig() {
	// Initialize database
	if err := db.InitDB("promptctl.db"); err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
}
