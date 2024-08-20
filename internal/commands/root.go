package commands

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "mail-discovery",
}

func init() {
	rootCmd.AddCommand(autoconfigCmd)
}

func Execute() {
	rootCmd.Execute()
}
