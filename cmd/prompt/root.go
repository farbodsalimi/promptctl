package prompt

import (
	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	promptCmd := &cobra.Command{
		Use:   "prompt",
		Short: "Manage prompts",
		Long:  `Add, update, list, and view prompts in vaults.`,
	}

	promptCmd.AddCommand(NewAddCmd())
	promptCmd.AddCommand(NewUpdateCmd())
	promptCmd.AddCommand(NewListCmd())
	promptCmd.AddCommand(NewHistoryCmd())
	promptCmd.AddCommand(NewShowCmd())

	return promptCmd
}
