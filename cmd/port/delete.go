package cmd_port

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	toDeletePort = -1
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Remove a port from a machine",
	Run: func(cmd *cobra.Command, args []string) {
		if toDeletePort < 0 || toDeletePort >= len(processedMachine.ForwardedPorts) {
			fmt.Println("Invalid port index specified")
			os.Exit(1)
		}

		processedMachine.ForwardedPorts = append(processedMachine.ForwardedPorts[:toDeletePort], processedMachine.ForwardedPorts[toDeletePort+1:]...)
	},
}

func registerDeleteCommand() {
	deleteCmd.Flags().IntVarP(&toDeletePort, "index", "i", -1, "The port to delete")
	RootCmd.AddCommand(deleteCmd)
}
