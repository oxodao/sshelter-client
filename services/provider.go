package services

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"

	_ "embed"

	"github.com/oxodao/sshelter_client/config"
	"github.com/oxodao/sshelter_client/models"
	"github.com/oxodao/sshelter_client/sshelter"
)

const VERBOSE = false

type Provider struct {
	Config         *config.Config
	CanReachServer bool
	IsSyncing      bool
	Client         *sshelter.Client
}

func NewProvider(cfg *config.Config) (*Provider, error) {
	client, err := sshelter.New(cfg, Info)
	if err != nil {
		return nil, err
	}

	return &Provider{
		Config:         cfg,
		CanReachServer: false,
		IsSyncing:      false,
		Client:         client,
	}, nil
}

func Info(i interface{}) {
	if VERBOSE {
		fmt.Println(">>> ", i)
	}
}

//go:embed ssh_config_template
var sshelterTemplate string

var sshelterSection = regexp.MustCompile(`(?s)\n# SSHELTER CONFIG(.*?)# SSHELTER END CONFIG`)

func (p *Provider) WriteSshConfig(machines []models.Machine) error {
	file, err := config.GetSshConfigFile()
	if err != nil {
		return err
	}

	stat, err := os.Stat(file)
	if err != nil {
		return err
	}

	perms := stat.Mode().Perm()

	txt, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(file+".sshelterbak", txt, perms)
	if err != nil {
		return errors.New("could not backup ssh config: " + err.Error())
	}

	cleaned := string(sshelterSection.ReplaceAll(txt, []byte("")))
	cleaned += p.buildSshConfig(machines)

	err = ioutil.WriteFile(file, []byte(cleaned), perms)
	if err != nil {
		// @TODO: Restore backup

		return errors.New("restoring the previous ssh config, could not save the file: " + err.Error())
	}

	return nil
}

func (p *Provider) buildSshConfig(machines []models.Machine) string {
	var config string

	for _, machine := range machines {
		if len(machine.ShortName) > 0 {
			config += fmt.Sprintf("Host %v\n", machine.ShortName)
		} else {
			config += fmt.Sprintf("Host %v\n", machine.Hostname)
		}

		if len(machine.Hostname) > 0 {
			config += fmt.Sprintf("\tHostname %v\n", machine.Hostname)
		}

		config += fmt.Sprintf("\tPort %v\n", machine.Port)

		if len(machine.Username) > 0 {
			config += fmt.Sprintf("\tUser %v\n", machine.Username)
		}

		if len(machine.OtherSettings) > 0 {
			lines := strings.Split(machine.OtherSettings, "\n")
			for _, line := range lines {
				config += fmt.Sprintf("\t%v\n", line)
			}
		}

		for _, lp := range machine.ForwardedPorts {
			local := strconv.Itoa(lp.LocalPort)
			remote := strconv.Itoa(lp.RemotePort)

			if lp.LocalHostname != nil && len(*lp.LocalHostname) > 0 {
				local = *lp.LocalHostname + ":" + local
			} else {
				local = "localhost" + ":" + local
			}

			if lp.RemoteHostname != nil && len(*lp.RemoteHostname) > 0 {
				remote = *lp.RemoteHostname + ":" + remote
			} else {
				remote = "localhost" + ":" + remote
			}

			config += fmt.Sprintf("\tLocalForward %v %v\n", local, remote)
		}

		config += "\n"
	}

	return fmt.Sprintf(sshelterTemplate, config)
}
