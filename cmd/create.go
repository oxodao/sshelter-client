package cmd

import (
	"fmt"
	"os"

	"github.com/oxodao/sshelter_client/models"
	"github.com/oxodao/sshelter_client/services"
	"github.com/spf13/cobra"
)

var (
	createMachine = &models.Machine{}
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new SSH machine",
	Run: func(cmd *cobra.Command, args []string) {
		err := createMachine.Validate()
		if err != nil {
			fmt.Println(err.Error()) // @TODO: Better error handling
			os.Exit(1)
		}

		err = services.GetProvider().Client.CreateMachine(createMachine)
		if err != nil {
			fmt.Println(err.Error()) // @TODO: Better error handling
			os.Exit(1)
		}

		fmt.Println("Machine created")
	},
}

func registerCreateCommand() {
	createCmd.Flags().StringVarP(&createMachine.Name, "name", "n", "", "Machine name in SSHelter")
	createCmd.Flags().StringVarP(&createMachine.ShortName, "shortname", "s", "", "The \"shortcut name\" (SSH name to use to connect to)")
	createCmd.Flags().StringVarP(&createMachine.Hostname, "host", "", "", "Hostname / IP of the machine")
	createCmd.Flags().IntVarP(&createMachine.Port, "port", "p", 22, "SSH port")
	createCmd.Flags().StringVarP(&createMachine.Username, "username", "u", "", "Username used to connect to the machine")
	rootCmd.AddCommand(createCmd)
}
