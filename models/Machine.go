package models

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
