package cmd

import (
	"fmt"
	"os"

	"github.com/oxodao/sshelter_client/config"
	"github.com/oxodao/sshelter_client/services"
)

func Version(prv *services.Provider) {
	fmt.Printf("SSHelter client version %v by %v\n", config.VERSION, config.AUTHOR)
	os.Exit(0)
}
