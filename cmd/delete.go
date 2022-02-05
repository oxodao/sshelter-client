package cmd

import (
	"fmt"
	"os"

	"github.com/oxodao/sshelter_client/models"
	"github.com/oxodao/sshelter_client/services"
	"github.com/spf13/cobra"
)

var deleteMachine = &models.Machine{}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a machine",
	Long:  "Delete a machine from SShelter.\n You can specify either the name or the shortname of the machine",
	Run: func(cmd *cobra.Command, args []string) {
		if len(deleteMachine.Name) == 0 && len(deleteMachine.ShortName) == 0 {
			fmt.Println("You must specify either the name or the shortname")
			os.Exit(1)
		}

		machines, err := services.GetProvider().Client.GetMachines()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		for _, machine := range machines {
			nameMatch := len(deleteMachine.Name) > 0 && machine.Name == deleteMachine.Name
			shortNameMatch := len(deleteMachine.ShortName) > 0 && machine.ShortName == deleteMachine.ShortName

			if nameMatch || shortNameMatch {
				err = services.GetProvider().Client.DeleteMachine(&machine)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				fmt.Println("Machine deleted")
				break
			}
		}
	},
}

func registerDeleteCommand() {
	deleteCmd.Flags().StringVarP(&deleteMachine.Name, "name", "n", "", "Machine name in SSHelter")
	deleteCmd.Flags().StringVarP(&deleteMachine.ShortName, "shortname", "s", "", "Machine shortname in SSHelter")
	rootCmd.AddCommand(deleteCmd)
}
