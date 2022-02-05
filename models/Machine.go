package models

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Machine struct {
	Id             *string         `json:"@id" yaml:"-,omitempty"`
	Name           string          `json:"name" yaml:"name"`
	ShortName      string          `json:"shortName" yaml:"short_name,omitempty"`
	Hostname       string          `json:"hostname" yaml:"hostname"`
	Port           int             `json:"port" yaml:"port"`
	Username       string          `json:"username" yaml:"username,omitempty"`
	OtherSettings  string          `json:"otherSettings" yaml:"other_settings,omitempty"`
	ForwardedPorts []ForwardedPort `json:"forwardedPorts" yaml:"forwarded_ports,omitempty"`
}

func (m *Machine) Validate() error {
	if len(m.Name) == 0 {
		return errors.New("The name is required")
	}

	if len(m.Hostname) == 0 {
		return errors.New("The hostname is required")
	}

	if m.Port < 1 || m.Port > 65535 {
		m.Port = 22
	}

	for _, fp := range m.ForwardedPorts {
		err := fp.Validate()
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *Machine) Patch(nm *Machine) {
	if len(nm.Name) > 0 {
		m.Name = nm.Name
	}

	if len(nm.ShortName) > 0 {
		m.ShortName = nm.ShortName
	}

	if len(nm.Hostname) > 0 {
		m.Hostname = nm.Hostname
	}

	if nm.Port > 0 {
		m.Port = nm.Port
	}

	if len(nm.Username) > 0 {
		m.Username = nm.Username
	}

	if len(nm.OtherSettings) > 0 {
		m.OtherSettings = nm.OtherSettings
	}

	if len(nm.ForwardedPorts) > 0 {
		m.ForwardedPorts = nm.ForwardedPorts
	}
}

func (m *Machine) String() string {
	str := ""

	connName := m.ShortName
	if len(connName) > 0 {
		connName = " (" + connName + ")"
	}

	str += "--- " + m.Name + connName + " ---\n"
	user := ""
	if len(m.Username) > 0 {
		user = m.Username + "@"
	}

	str += user + m.Hostname + ":" + strconv.Itoa(m.Port) + "\n"
	if len(m.OtherSettings) > 0 {
		str += "Other settings: " + m.OtherSettings + "\n"
	}

	if len(m.ForwardedPorts) > 0 {
		str += "Remote ports mapped on local machine:\n"
		for i, fp := range m.ForwardedPorts {
			str += fmt.Sprintf("\t - [%v] %v\n", i, fp.String())
		}
	}

	return str
}

func (m *Machine) ToSshString() string {
	config := ""
	if len(m.ShortName) > 0 {
		config += fmt.Sprintf("Host %v\n", m.ShortName)
	} else {
		config += fmt.Sprintf("Host %v\n", m.Hostname)
	}

	if len(m.Hostname) > 0 {
		config += fmt.Sprintf("\tHostname %v\n", m.Hostname)
	}

	config += fmt.Sprintf("\tPort %v\n", m.Port)

	if len(m.Username) > 0 {
		config += fmt.Sprintf("\tUser %v\n", m.Username)
	}

	if len(m.OtherSettings) > 0 {
		lines := strings.Split(m.OtherSettings, "\n")
		for _, line := range lines {
			config += fmt.Sprintf("\t%v\n", line)
		}
	}

	for _, lp := range m.ForwardedPorts {
		config += "\t" + lp.ToSshString() + "\n"
	}

	config += "\n"

	return config
}
