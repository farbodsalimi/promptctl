package prompt

import (
	"fmt"

	"github.com/farbodsalimi/promptctl/internal/db"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewAddCmd() *cobra.Command {
	var (
		vaultName  string
		promptName string
		content    string
	)

	promptAddCmd := &cobra.Command{
		Use:   "add --vault=<vault> --name=<name> --prompt=<prompt>",
		Short: "Add a new prompt to a vault",
		Run: func(cmd *cobra.Command, args []string) {
			// Get vault
			vault, err := db.GetVaultByName(vaultName)
			if err != nil {
				log.Fatalf("vault not found: %s", vaultName)
			}

			if err := db.CreatePrompt(vault.ID, promptName, content); err != nil {
				log.Fatalf("failed to create prompt: %v", err)
			}

			fmt.Printf("Added prompt '%s' to vault '%s'\n", promptName, vaultName)
		},
	}

	promptAddCmd.Flags().
		StringVarP(&vaultName, "vault", "v", "", "Name of the vault to add the prompt to")
	promptAddCmd.Flags().
		StringVarP(&promptName, "name", "n", "", "Name for the new prompt (must be unique within vault)")
	promptAddCmd.Flags().
		StringVarP(&content, "prompt", "p", "", "Prompt content (supports Go template syntax)")

	promptAddCmd.MarkFlagRequired("vault")
	promptAddCmd.MarkFlagRequired("name")
	promptAddCmd.MarkFlagRequired("prompt")

	return promptAddCmd
}
