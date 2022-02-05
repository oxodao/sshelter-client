package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/oxodao/sshelter_client/models"
	"github.com/oxodao/sshelter_client/services"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var toImportFile string

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import a machine or a group of machines",
	Run: func(cmd *cobra.Command, args []string) {
		if len(toImportFile) == 0 {
			fmt.Println("You must specify a file to import")
			os.Exit(1)
		}

		str, err := ioutil.ReadFile(toImportFile)
		if err != nil {
			fmt.Println("Error while reading the file:", err)
			os.Exit(1)
		}

		machines := []models.Machine{}
		err = yaml.Unmarshal(str, &machines)
		if err != nil {
			fmt.Println("Error while reading the file:", err)
			os.Exit(1)
		}

		if len(machines) == 0 {
			fmt.Println("No machine to import")
			return
		}

		for _, machine := range machines {
			err = services.GetProvider().Client.CreateMachine(&machine)
			if err != nil {
				fmt.Println("Error while creating the machine:", err)
				continue
			}

			fmt.Println("Machine created:", machine.Name)
		}
	},
}

func registerImportCommand() {
	importCmd.Flags().StringVarP(&toImportFile, "file", "f", "", "File to import")

	rootCmd.AddCommand(importCmd)
}
