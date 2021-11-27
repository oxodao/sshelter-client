package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/oxodao/sshelter_client/models"
	"github.com/oxodao/sshelter_client/services"
	"gopkg.in/yaml.v2"
)

type machineFile struct {
	Machines []models.Machine `yaml:"machines"`
}

func CreateMachine(prv *services.Provider, filename string) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Error reading file: " + err.Error())
		os.Exit(1)
	}

	var machines machineFile
	err = yaml.Unmarshal(file, &machines)
	if err != nil {
		fmt.Println("Error reading file: " + err.Error())
		os.Exit(1)
	}

	for _, machine := range machines.Machines {
		err := machine.Validate()
		if err != nil {
			fmt.Println("Invalid machine: " + err.Error())
			os.Exit(1)
		}

		err = prv.Client.CreateMachine(&machine)
		if err != nil {
			fmt.Println("Error creating machine: " + err.Error())
			os.Exit(1)
		}
	}

	if len(machines.Machines) > 1 {
		fmt.Println("Machines created")
	} else if len(machines.Machines) == 1 {
		fmt.Println("Machine created")
	} else {
		fmt.Println("No machine created")
	}
	os.Exit(0)
}
