package cmd_port

import (
	"fmt"
	"os"

	"github.com/oxodao/sshelter_client/models"
	"github.com/oxodao/sshelter_client/services"
	"github.com/spf13/cobra"
)

var (
	processedMachine *models.Machine
	byName           string
	byShortName      string
)

var RootCmd = &cobra.Command{
	Use:   "port",
	Short: "Manage SSH ports for a machine",
	Long:  "Manage SSH ports for a machine\nUse either --by-name or --by-shortname to select your machine",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if len(byName) == 0 && len(byShortName) == 0 {
			fmt.Println("You must specify either --by-name or --by-shortname")
			os.Exit(1)
		}

		machines, err := services.GetProvider().Client.GetMachines()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		for _, machine := range machines {
			if (len(byName) != 0 && machine.Name == byName) || (len(byShortName) != 0 && machine.ShortName == byShortName) {
				processedMachine = &machine
				return
			}
		}

		fmt.Println("Machine not found")
		os.Exit(1)
	},

	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		err := services.GetProvider().Client.UpdateMachine(processedMachine)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println("Machine updated:")
		fmt.Println(processedMachine.String())
	},
}

func init() {
	RootCmd.PersistentFlags().StringVarP(&byName, "by-name", "", "", "The name of the machine to add/remove port")
	RootCmd.PersistentFlags().StringVarP(&byShortName, "by-shortname", "", "", "The shortname of the machine to add/remove port")

	registerAddCommand()
	registerDeleteCommand()
}
