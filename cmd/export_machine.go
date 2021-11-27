package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/oxodao/sshelter_client/models"
	"github.com/oxodao/sshelter_client/services"
	"gopkg.in/yaml.v2"
)

func ExportMachine(prv *services.Provider, machinesList string) {
	remoteMachines, err := prv.Client.GetMachines()
	if err != nil {
		fmt.Println("Could not get remote machines: ", err)
		os.Exit(1)
	}

	machines := strings.Split(machinesList, ",")
	filename := "sshelter-export.yml"
	if len(machines) == 1 && machines[0] != "all" {
		filename = fmt.Sprintf("%s.yml", machines[0])
	}

	exportedMachines := struct {
		Machines []models.Machine `yaml:"machines"`
	}{
		Machines: []models.Machine{},
	}

	if machinesList != "all" {
		for _, machine := range machines {
			if currentMachine := getMachineInArray(machine, remoteMachines); currentMachine != nil {
				exportedMachines.Machines = append(exportedMachines.Machines, *currentMachine)
			}
		}
	} else {
		exportedMachines.Machines = remoteMachines
	}

	if len(exportedMachines.Machines) == 0 {
		fmt.Println("No machines found")
		os.Exit(1)
	}

	str, err := yaml.Marshal(exportedMachines)
	if err != nil {
		fmt.Println("Could not export machines: ", err)
		os.Exit(1)
	}

	err = ioutil.WriteFile(filename, str, 0644)
	if err != nil {
		fmt.Println("Could not write exported machines: ", err)
		os.Exit(1)
	}

	fmt.Println("Machines exported as " + filename)
	os.Exit(0)
}

func getMachineInArray(machineName string, machines []models.Machine) *models.Machine {
	for _, machine := range machines {
		if machine.Name == machineName {
			return &machine
		}
	}

	return nil
}
