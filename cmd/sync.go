package cmd

import (
	"fmt"
	"os"

	"github.com/oxodao/sshelter_client/services"
)

func Sync(prv *services.Provider) {
	machines, err := prv.Client.GetMachines()
	if err != nil {
		panic(err)
	}

	err = prv.WriteSshConfig(machines)
	if err != nil {
		panic(err)
	}

	fmt.Println("Synced!")

	os.Exit(0)
}
