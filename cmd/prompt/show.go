package prompt

import (
	"fmt"
	"strconv"

	"github.com/farbodsalimi/promptctl/internal/db"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewShowCmd() *cobra.Command {
	var (
		vaultName  string
		promptName string
		revision   int
	)

	var promptShowCmd = &cobra.Command{
		Use:   "show --vault=<vault> --prompt=<prompt>",
		Short: "Show prompt content",
		Run: func(cmd *cobra.Command, args []string) {
			var content string
			var err error

			if revision > 0 {
				content, err = db.GetPromptVersionContentByVersion(vaultName, promptName, revision)
			} else {
				prompt, err2 := db.GetPromptByName(vaultName, promptName)
				if err2 != nil {
					log.Fatalf("prompt not found: %s/%s", vaultName, promptName)
				}
				content, err = db.GetPromptContent(prompt.ID)
			}

			if err != nil {
				log.Fatalf("failed to get prompt content: %v", err)
			}

			versionStr := "latest"
			if revision > 0 {
				versionStr = "v" + strconv.Itoa(revision)
			}

			fmt.Printf("Prompt '%s' in vault '%s' (%s):\n", promptName, vaultName, versionStr)
			fmt.Printf("---\n%s\n---\n", content)
		},
	}

	promptShowCmd.Flags().
		StringVarP(&vaultName, "vault", "v", "", "Name of the vault containing the prompt")
	promptShowCmd.Flags().
		StringVarP(&promptName, "prompt", "p", "", "Name of the prompt to display")
	promptShowCmd.Flags().
		IntVarP(&revision, "revision", "r", 0, "Specific revision number to show (default: latest revision)")

	promptShowCmd.MarkFlagRequired("vault")
	promptShowCmd.MarkFlagRequired("prompt")

	return promptShowCmd
}
