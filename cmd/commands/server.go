package commands

import (
	"schooli-api/cmd/app"

	"github.com/spf13/cobra"
)

func Server() *cobra.Command {
	serverCmd := &cobra.Command{
		Use:          "server",
		Short:        "Run the Schooli server",
		SilenceUsage: true,
	}
	serverCmd.RunE = func(cmd *cobra.Command, _ []string) error {
		cfg, err := cmd.Flags().GetString("config")
		if err != nil {
			return err
		}
		server, errx := app.StartApplication(cfg)
		if errx != nil {
			return errx
		}

		return server.Run()
	}
	return serverCmd
}
