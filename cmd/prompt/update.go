package prompt

import (
	"fmt"

	"github.com/farbodsalimi/promptctl/internal/db"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewUpdateCmd() *cobra.Command {
	var (
		vaultName  string
		promptName string
		content    string
	)

	var promptUpdateCmd = &cobra.Command{
		Use:   "update <vault> <name> <prompt>",
		Short: "Update an existing prompt (creates new version)",

		Run: func(cmd *cobra.Command, args []string) {
			if err := db.UpdatePrompt(vaultName, promptName, content); err != nil {
				log.Fatalf("failed to update prompt: %v", err)
			}

			fmt.Printf("Updated prompt '%s' in vault '%s'\n", promptName, vaultName)
		},
	}

	promptUpdateCmd.Flags().StringVarP(&vaultName, "vault", "v", "", "Vault")
	promptUpdateCmd.Flags().StringVarP(&promptName, "name", "n", "", "Name")
	promptUpdateCmd.Flags().StringVarP(&content, "prompt", "p", "", "Prompt")

	promptUpdateCmd.MarkFlagRequired("vault")
	promptUpdateCmd.MarkFlagRequired("name")
	promptUpdateCmd.MarkFlagRequired("prompt")

	return promptUpdateCmd
}
