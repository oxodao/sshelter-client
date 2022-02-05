package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/oxodao/sshelter_client/models"
	"github.com/oxodao/sshelter_client/services"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var (
	exportByName      string
	exportByShortname string
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export a machine or a group of machines",
	Run: func(cmd *cobra.Command, args []string) {
		if len(exportByName) == 0 && len(exportByShortname) == 0 {
			fmt.Println("You must specify either --by-name or --by-shortname")
			os.Exit(1)
		}

		if len(exportByName) > 0 && len(exportByShortname) > 0 {
			fmt.Println("You must specify either --by-name or --by-shortname, not both")
			os.Exit(1)
		}

		names := strings.Split(exportByName, ",")
		shortnames := strings.Split(exportByShortname, ",")

		machines, err := services.GetProvider().Client.GetMachines()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		exportedMachines := []*models.Machine{}
		for _, machine := range machines {
			for _, name := range names {
				if machine.Name == name {
					// loving pointers
					m := machine
					exportedMachines = append(exportedMachines, &m)
				}
			}

			for _, shortname := range shortnames {
				if machine.ShortName == shortname {
					m := machine
					exportedMachines = append(exportedMachines, &m)
				}
			}
		}

		for _, machine := range exportedMachines {
			machine.Id = nil
		}

		str, _ := yaml.Marshal(exportedMachines)
		fmt.Println(string(str))
	},
}

func registerExportCommand() {
	exportCmd.Flags().StringVarP(&exportByName, "by-name", "", "", "The name of the machine(s), splitted by comma")
	exportCmd.Flags().StringVarP(&exportByShortname, "by-shortname", "", "", "The shortname of the machine(s) to export, splitted by comma")
	rootCmd.AddCommand(exportCmd)
}
