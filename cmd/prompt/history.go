package prompt

import (
	"fmt"

	"github.com/farbodsalimi/promptctl/internal/db"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewHistoryCmd() *cobra.Command {
	var (
		vaultName  string
		promptName string
	)

	var promptHistoryCmd = &cobra.Command{
		Use:   "history <vault> <name>",
		Short: "Show version history of a prompt",
		Run: func(cmd *cobra.Command, args []string) {
			prompt, err := db.GetPromptByName(vaultName, promptName)
			if err != nil {
				log.Fatalf("prompt not found: %s/%s", vaultName, promptName)
			}

			history, err := db.GetPromptHistory(prompt.ID)
			if err != nil {
				log.Fatalf("failed to get prompt history: %v", err)
			}

			fmt.Printf("History for prompt '%s' in vault '%s':\n", promptName, vaultName)
			for _, entry := range history {
				fmt.Printf("  %s\n", entry)
			}
		},
	}

	promptHistoryCmd.Flags().StringVarP(&vaultName, "vault", "v", "", "Vault")
	promptHistoryCmd.Flags().StringVarP(&promptName, "name", "n", "", "Name")

	promptHistoryCmd.MarkFlagRequired("vault")
	promptHistoryCmd.MarkFlagRequired("name")

	return promptHistoryCmd
}
