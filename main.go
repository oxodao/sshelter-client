package main

import (
	"github.com/oxodao/sshelter_client/cmd"
	"github.com/oxodao/sshelter_client/config"
	"github.com/oxodao/sshelter_client/services"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	err = services.NewProvider(cfg)
	if err != nil {
		panic(err)
	}

	cmd.Execute()
}
