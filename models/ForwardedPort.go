package models

import (
	"errors"
	"strconv"
)

type ForwardedPort struct {
	LocalHostname *string `json:"local_hostname" yaml:"local_hostname,omitempty"`
	LocalPort     int     `json:"local_port" yaml:"local_port"`

	RemoteHostname *string `json:"remote_hostname" yaml:"remote_hostname,omitempty"`
	RemotePort     int     `json:"remote_port" yaml:"remote_port"`

	Reversed bool `json:"reversed" yaml:"reversed"`
}

func (f *ForwardedPort) Validate() error {
	if f.LocalHostname != nil && len(*f.LocalHostname) > 0 {
		f.LocalHostname = nil
	}

	if f.RemoteHostname != nil && len(*f.RemoteHostname) > 0 {
		f.RemoteHostname = nil
	}

	if f.LocalPort < 1 || f.LocalPort > 65535 || f.RemotePort < 1 || f.RemotePort > 65535 {
		return errors.New("Invalid port")
	}

	return nil
}

func (f *ForwardedPort) String() string {
	str := ""

	if f.LocalHostname != nil {
		str += *f.LocalHostname + ":"
	}

	str += strconv.Itoa(f.LocalPort)

	str += " (local)"

	if f.Reversed {
		str += " <- "
	} else {
		str += " -> "
	}

	if f.RemoteHostname != nil {
		str += *f.RemoteHostname + ":"
	}

	local := "(Mapped locally)"
	if f.Reversed {
		local = "(Mapped remotely)"
	}

	str += strconv.Itoa(f.RemotePort) + " (remote) " + local

	return str
}

func (f *ForwardedPort) ToSshString() string {
	local := strconv.Itoa(f.LocalPort)
	remote := strconv.Itoa(f.RemotePort)

	if f.LocalHostname != nil && len(*f.LocalHostname) > 0 {
		local = *f.LocalHostname + ":" + local
	} else {
		local = "localhost" + ":" + local
	}

	if f.RemoteHostname != nil && len(*f.RemoteHostname) > 0 {
		remote = *f.RemoteHostname + ":" + remote
	} else {
		remote = "localhost" + ":" + remote
	}

	if f.Reversed {
		return "RemoteForward " + remote + " " + local
	}

	return "LocalForward " + local + " " + remote
}
