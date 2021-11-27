package main

import (
	"flag"
	"fmt"
	"strconv"

	"github.com/oxodao/sshelter_client/cmd"
	"github.com/oxodao/sshelter_client/config"
	"github.com/oxodao/sshelter_client/services"
)

// @TODO forwarding port both ways

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	prv, err := services.NewProvider(cfg)
	if err != nil {
		panic(err)
	}

	execCommands(prv)

	machines, err := prv.Client.GetMachines()
	if err != nil {
		panic(err)
	}

	for _, machine := range machines {
		fmt.Printf("%v", machine.Name)
		if len(machine.ShortName) > 0 {
			fmt.Printf(" (%v)", machine.ShortName)
		}
		fmt.Println()

		dsn := fmt.Sprintf("%v:%v", machine.Hostname, machine.Port)
		if len(machine.Username) > 0 {
			dsn = machine.Username + "@" + dsn
		}

		fmt.Println("\t" + dsn)

		if len(machine.ForwardedPorts) > 0 {
			fmt.Println("\tForwarded ports:")
			for _, port := range machine.ForwardedPorts {
				// Should not happen but it seems like the API is accepting them anyway
				if port.LocalPort == 0 || port.RemotePort == 0 {
					continue
				}

				localPort := strconv.Itoa(port.LocalPort)
				remotePort := strconv.Itoa(port.LocalPort)

				if port.LocalHostname != nil && len(*port.LocalHostname) > 0 {
					localPort = *port.LocalHostname + ":" + localPort
				}

				if port.RemoteHostname != nil && len(*port.RemoteHostname) > 0 {
					remotePort = *port.RemoteHostname + ":" + remotePort
				}

				fmt.Printf("\t\t%v -> %v\n", localPort, remotePort)
			}
		}

		if len(machine.OtherSettings) > 0 {
			fmt.Println("\tOther settings:")
			fmt.Printf("\t\t%v\n", machine.OtherSettings)
		}

		fmt.Println()
	}
}

func execCommands(prv *services.Provider) {
	genSkeleton := flag.String("gen-skel", "", "Generate the skeleton to create machines")
	createMachine := flag.String("create", "", "Create new machine(s)")
	deleteMachine := flag.String("delete", "", "Delete a machine")
	exportMachine := flag.String("export", "", "Export one or multiple machines, name separated by commas (or 'all' to export everything)")

	syncMachines := flag.Bool("sync", false, "Sync the machines with the server")
	version := flag.Bool("version", false, "Displays the version")

	flag.Parse()

	if *version {
		cmd.Version(prv)
	}

	if len(*genSkeleton) > 0 {
		cmd.GenerateSkeleton(prv, *genSkeleton)
	}

	if len(*createMachine) > 0 {
		cmd.CreateMachine(prv, *createMachine)
	}

	if len(*deleteMachine) > 0 {
		cmd.DeleteMachine(prv, *deleteMachine)
	}

	if len(*exportMachine) > 0 {
		cmd.ExportMachine(prv, *exportMachine)
	}

	if *syncMachines {
		cmd.Sync(prv)
	}
}
