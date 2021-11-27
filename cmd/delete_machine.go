package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/oxodao/sshelter_client/services"
)

func DeleteMachine(prv *services.Provider, machinesNames string) {
	remoteMachines, err := prv.Client.GetMachines()
	if err != nil {
		fmt.Println("Could not get remote machines: ", err)
		os.Exit(1)
	}

	machines := strings.Split(machinesNames, ",")

	removedMachine := 0
	for _, machine := range machines {
		if currentMachine := getMachineInArray(machine, remoteMachines); currentMachine != nil {
			err := prv.Client.DeleteMachine(currentMachine)
			if err != nil {
				fmt.Println("Could not delete machine: ", err)
				os.Exit(1)
			}

			fmt.Println("Machine removed: ", currentMachine.Name)

			removedMachine++
		}
	}

	if removedMachine == 0 {
		fmt.Println("No machines found")
		os.Exit(1)
	}

	os.Exit(0)
}
