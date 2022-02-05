package cmd

import (
	cmd_port "github.com/oxodao/sshelter_client/cmd/port"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "sshelter",
	Short: "SShelter is a tool for managing your SSH config file",
}

func init() {
	registerListCommand()
	registerCreateCommand()
	registerEditCommand()
	registerDeleteCommand()
	registerSyncCommand()
	registerImportCommand()
	registerExportCommand()
	registerVersionCommand()

	rootCmd.AddCommand(cmd_port.RootCmd)
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
