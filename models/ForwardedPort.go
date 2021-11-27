package models

import "errors"

type ForwardedPort struct {
	LocalHostname *string `json:"local_hostname" yaml:"local_hostname,omitempty"`
	LocalPort     int     `json:"local_port" yaml:"local_port"`

	RemoteHostname *string `json:"remote_hostname" yaml:"remote_hostname,omitempty"`
	RemotePort     int     `json:"remote_port" yaml:"remote_port"`
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
