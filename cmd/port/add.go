package cmd_port

import (
	"fmt"
	"os"

	"github.com/oxodao/sshelter_client/models"
	"github.com/spf13/cobra"
)

var (
	toAddLocalHost string
	toAddLocalPort int

	toAddRemoteHost string
	toAddRemotePort int

	toAddReversed bool
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a port to a machine",
	Run: func(cmd *cobra.Command, args []string) {
		invalid := false

		if toAddLocalPort < 0 || toAddLocalPort > 65535 {
			fmt.Printf("Invalid port: %v\n", toAddLocalPort)
			invalid = true
		}

		if toAddRemotePort < 0 || toAddRemotePort > 65535 {
			fmt.Printf("Invalid port: %v\n", toAddRemotePort)
			invalid = true
		}

		if invalid {
			os.Exit(1)
		}

		processedMachine.ForwardedPorts = append(processedMachine.ForwardedPorts, models.ForwardedPort{
			LocalHostname:  &toAddLocalHost,
			LocalPort:      toAddLocalPort,
			RemoteHostname: &toAddRemoteHost,
			RemotePort:     toAddRemotePort,
			Reversed:       toAddReversed,
		})

	},
}

func registerAddCommand() {
	addCmd.Flags().StringVarP(&toAddLocalHost, "lh", "", "", "The local host to forward")
	addCmd.Flags().IntVarP(&toAddLocalPort, "lp", "", -1, "The port on the local machine")
	addCmd.Flags().StringVarP(&toAddRemoteHost, "rh", "", "", "The remote host to forward")
	addCmd.Flags().IntVarP(&toAddRemotePort, "rp", "", -1, "The port on the remote machine")
	addCmd.Flags().BoolVarP(&toAddReversed, "reversed", "", false, "The local port should be mapped on the remote machine instead of the other way around")

	RootCmd.AddCommand(addCmd)
}
