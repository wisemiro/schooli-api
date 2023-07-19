package commands

import (
	"github.com/spf13/cobra"
)

func Run(args []string) error {
	rootCmd := &cobra.Command{
		Use:   "Schooli",
		Short: "School ecommerce",
		Long:  `School ecommerce`,
	}
	rootCmd.PersistentFlags().StringP("config", "c", "", "Configuration file to use.")
	rootCmd.AddCommand(Server())
	rootCmd.SetArgs(args)
	return rootCmd.Execute()
}
