package prompt

import (
	"fmt"

	"github.com/farbodsalimi/promptctl/internal/db"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewListCmd() *cobra.Command {
	var (
		vaultName string
	)

	var promptListCmd = &cobra.Command{
		Use:   "list --vault=<vault>",
		Short: "List prompts in a vault",
		Run: func(cmd *cobra.Command, args []string) {
			prompts, err := db.GetPrompts(vaultName)
			if err != nil {
				log.Fatalf("failed to list prompts: %v", err)
			}

			if len(prompts) == 0 {
				fmt.Printf("No prompts found in vault '%s'\n", vaultName)
				return
			}

			fmt.Printf("Prompts in vault '%s':\n", vaultName)
			for _, prompt := range prompts {
				fmt.Printf(
					"  %s (v%d, created: %s)\n",
					prompt.Name,
					prompt.LatestVersion,
					prompt.Created,
				)
			}
		},
	}

	promptListCmd.Flags().StringVarP(&vaultName, "vault", "v", "", "Vault")
	promptListCmd.MarkFlagRequired("vault")

	return promptListCmd
}
