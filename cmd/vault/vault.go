package vault

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/farbodsalimi/promptctl/internal/db"
)

var vaultCreateCmd = &cobra.Command{
	Use:   "create <name>",
	Short: "Create a new vault",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		if err := db.CreateVault(name); err != nil {
			log.Fatalf("failed to create vault: %v", err)
		}
		fmt.Printf("Created vault: %s\n", name)
	},
}

var vaultListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all vaults",
	Run: func(cmd *cobra.Command, args []string) {
		vaults, err := db.GetVaults()
		if err != nil {
			log.Fatalf("failed to list vaults: %v", err)
		}

		if len(vaults) == 0 {
			fmt.Println("No vaults found")
			return
		}

		fmt.Println("Vaults:")
		for _, vault := range vaults {
			fmt.Printf("  %s (created: %s)\n", vault.Name, vault.Created)
		}
	},
}

var vaultDeleteCmd = &cobra.Command{
	Use:   "delete <name>",
	Short: "Delete a vault",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		rowsAffected, err := db.DeleteVault(name)
		if err != nil {
			log.Fatalf("failed to delete vault: %v", err)
		}
		if rowsAffected == 0 {
			log.Fatalf("vault not found: %s", name)
		}
		fmt.Printf("Deleted vault: %s\n", name)
	},
}

func NewRootCmd() *cobra.Command {
	vaultCmd := &cobra.Command{
		Use:   "vault",
		Short: "Manage prompt vaults",
		Long:  `Create, list, and delete prompt vaults.`,
	}

	vaultCmd.AddCommand(vaultCreateCmd)
	vaultCmd.AddCommand(vaultListCmd)
	vaultCmd.AddCommand(vaultDeleteCmd)

	return vaultCmd
}
