package commands

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/stavros-k/go-mail-discovery/internal/generators"
	"github.com/stavros-k/go-mail-discovery/internal/utils"
)

var autoconfigCmd = &cobra.Command{
	Use:   "autoconfig",
	Short: "Generate autoconfig",
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

		config, err := generators.NewConfigV1_1(generators.ConfigV1_1Params{
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
	autoconfigCmd.Flags().StringP("email", "e", "", "Email address")
	autoconfigCmd.MarkFlagRequired("email")
}
