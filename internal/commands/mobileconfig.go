package commands

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/stavros-k/go-mail-discovery/internal/generators"
	"github.com/stavros-k/go-mail-discovery/internal/utils"
)

var mobileconfigCmd = &cobra.Command{
	Use:   "mobileconfig",
	Short: "Generate mobileconfig",
	Run: func(cmd *cobra.Command, args []string) {
		emailAddress := cmd.Flag("email").Value.String()
		domain, err := utils.GetDomainFromEmailAddress(emailAddress)
		if err != nil {
			log.Fatal(err)
		}
		provider, err := utils.GetProviderFromMX(domain, false)
		if err != nil {
			log.Fatal(err)
		}

		config, err := generators.NewMobileConfig(generators.MobileConfigParams{
			Domain:      domain,
			DisplayName: emailAddress,
			Username:    emailAddress,
			Provider:    provider,
		})
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(config)
	},
}

func init() {
	mobileconfigCmd.Flags().StringP("email", "e", "", "Email address")
	mobileconfigCmd.MarkFlagRequired("email")
	rootCmd.AddCommand(mobileconfigCmd)
}
