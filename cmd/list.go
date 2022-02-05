package cmd

import (
	"fmt"
	"os"

	"github.com/oxodao/sshelter_client/services"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all machines",
	Run: func(cmd *cobra.Command, args []string) {
		machines, err := services.GetProvider().Client.GetMachines()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		for _, machine := range machines {
			fmt.Println(machine.String())
		}
	},
}

func registerListCommand() {
	rootCmd.AddCommand(listCmd)
}
