package commands

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/stavros-k/go-mail-discovery/internal/generators"
	"github.com/stavros-k/go-mail-discovery/internal/utils"
)

var autodiscoverCmd = &cobra.Command{
	Use:   "autodiscover",
	Short: "Generate autodiscover",
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

		config, err := generators.NewAutoDiscoverConfig(generators.AutoDiscoverConfigParams{
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
	autodiscoverCmd.Flags().StringP("email", "e", "", "Email address")
	autodiscoverCmd.MarkFlagRequired("email")
	rootCmd.AddCommand(autodiscoverCmd)
}
