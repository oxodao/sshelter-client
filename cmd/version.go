package cmd

import (
	"fmt"
	"strings"

	"github.com/oxodao/sshelter_client/utils"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of sshelter",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%v %v (Commit %v) by %v\nhttps://github.com/%v/%v\n", utils.SOFTWARE_NAME, utils.VERSION, utils.COMMIT, utils.AUTHOR, strings.ToLower(utils.AUTHOR), strings.ToLower(utils.SOFTWARE_NAME))
	},
}

func registerVersionCommand() {
	rootCmd.AddCommand(versionCmd)
}
