package services

import (
	"fmt"

	_ "embed"

	"github.com/oxodao/sshelter_client/config"
	"github.com/oxodao/sshelter_client/sshelter"
)

const VERBOSE = false

var provider *Provider

type Provider struct {
	Config         *config.Config
	CanReachServer bool
	IsSyncing      bool
	Client         *sshelter.Client
}

func NewProvider(cfg *config.Config) error {
	client, err := sshelter.New(cfg, Info)
	if err != nil {
		return err
	}

	provider = &Provider{
		Config:         cfg,
		CanReachServer: false,
		IsSyncing:      false,
		Client:         client,
	}

	return nil
}

func GetProvider() *Provider {
	return provider
}

func Info(i interface{}) {
	if VERBOSE {
		fmt.Println(">>> ", i)
	}
}
