package cmd

import (
	"fmt"
	"os"

	"github.com/oxodao/sshelter_client/models"
	"github.com/oxodao/sshelter_client/services"
	"github.com/spf13/cobra"
)

var (
	editMachine     = &models.Machine{}
	editByName      = ""
	editByShortName = ""
)

var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit a SSH machine",
	Run: func(cmd *cobra.Command, args []string) {
		if len(editByName) == 0 && len(editByShortName) == 0 {
			fmt.Println("You must specify either the name or the shortname")
			os.Exit(1)
		}

		machines, err := services.GetProvider().Client.GetMachines()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		for _, machine := range machines {
			nameMatch := len(editByName) > 0 && machine.Name == editByName
			shortNameMatch := len(editByShortName) > 0 && machine.ShortName == editByShortName

			if nameMatch || shortNameMatch {
				m := &machine
				m.Patch(editMachine)

				err = m.Validate()
				if err != nil {
					fmt.Println(err)

					os.Exit(1)
				}

				err = services.GetProvider().Client.UpdateMachine(m)
				if err != nil {
					fmt.Println(err)

					os.Exit(1)
				}

				fmt.Println("Machine updated")

				break
			}
		}
	},
}

func registerEditCommand() {
	editCmd.Flags().StringVarP(&editByName, "by-name", "", "", "The machine name to edit")
	editCmd.Flags().StringVarP(&editByShortName, "by-shortname", "", "", "The machine shortname to edit")

	editCmd.Flags().StringVarP(&editMachine.Name, "name", "n", "", "The new machine name")
	editCmd.Flags().StringVarP(&editMachine.ShortName, "shortname", "s", "", "The new machine shortname")
	editCmd.Flags().StringVarP(&editMachine.Hostname, "host", "", "", "The new machine hostname / IP")
	editCmd.Flags().IntVarP(&editMachine.Port, "port", "p", -1, "The new SSH port")
	editCmd.Flags().StringVarP(&editMachine.Username, "username", "u", "", "The new username")
	rootCmd.AddCommand(editCmd)
}
