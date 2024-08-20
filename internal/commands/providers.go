package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stavros-k/go-mail-discovery/internal/providers"
)

var providersCmd = &cobra.Command{
	Use:   "providers",
	Short: "List available providers",
	Run: func(cmd *cobra.Command, args []string) {
		for _, provider := range providers.ListProvidersWithInfo() {
			fmt.Println(provider)
		}
	},
}
var providerCmd = &cobra.Command{
	Use:   "provider",
	Short: "Get provider info",
	Run: func(cmd *cobra.Command, args []string) {
		id, _ := cmd.Flags().GetString("id")
		if id == "" {
			fmt.Println("Error: provider ID is required")
			return
		}
		info, err := providers.GetProviderInfo(id)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Println(info)

	},
}

func init() {
	providerCmd.Flags().StringP("id", "i", "", "Provider ID")
	providerCmd.MarkFlagRequired("id")

	rootCmd.AddCommand(providerCmd)
	rootCmd.AddCommand(providersCmd)
}
