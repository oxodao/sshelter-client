package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"

	"github.com/oxodao/sshelter_client/config"
	"github.com/oxodao/sshelter_client/services"
	"github.com/spf13/cobra"
)

var sshelterSection = regexp.MustCompile(`(?s)\n# SSHELTER CONFIG(.*?)# SSHELTER END CONFIG`)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync your SSH config file",
	Long:  "Fetch your machines from the SSHelter instance and apply it to your sshconfig.",
	Run: func(cmd *cobra.Command, args []string) {
		machines, err := services.GetProvider().Client.GetMachines()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		sshConfig, err := config.GetSshConfigFile()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		stat, err := os.Stat(sshConfig)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		perms := stat.Mode().Perm()

		txt, err := ioutil.ReadFile(sshConfig)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		err = ioutil.WriteFile(sshConfig+".sshelterbak", txt, perms)
		if err != nil {
			fmt.Println("could not backup ssh config: " + err.Error())
			os.Exit(1)
		}

		cleaned := string(sshelterSection.ReplaceAll(txt, []byte("")))

		cleaned += "# SSHELTER CONFIG\n"
		cleaned += "# Please do not edit this section. Any changes here will be overwritten or might break sshelter.\n\n"

		for _, m := range machines {
			cleaned += m.ToSshString()
		}

		cleaned += "# SSHELTER END CONFIG"

		err = ioutil.WriteFile(sshConfig, []byte(cleaned), perms)
		if err != nil {
			// @TODO: Restore backup

			fmt.Println("Please restore your .ssh/config manually from .ssh/sshelterbak, could not save the file: " + err.Error())
		}

		fmt.Println("Synced!")
	},
}

func registerSyncCommand() {
	rootCmd.AddCommand(syncCmd)
}
