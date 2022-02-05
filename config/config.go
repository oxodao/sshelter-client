package config

import (
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Server        Server        `yaml:"server"`
	LocalOverride LocalOverride `yaml:"local_override"`
}

type Server struct {
	Url      string `yaml:"url"`
	Username string `yaml:"username"`

	Token        string `yaml:"token"`
	RefreshToken string `yaml:"refresh_token"`
}

type LocalOverride struct {
}

func Load() (*Config, error) {
	path, err := getConfigPath()
	if err != nil {
		return nil, err
	}

	cfg := Config{}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, &cfg)
	validate(&cfg)

	return &cfg, err
}

func validate(cfg *Config) {
	if len(cfg.Server.Url) == 0 {
		panic("No server URL has been set!")
	}

	if len(cfg.Server.Username) == 0 {
		panic("No username has been set!")
	}

	if !strings.HasSuffix(cfg.Server.Url, "/") {
		cfg.Server.Url += "/"
	}
}

func (cfg *Config) Save() error {
	path, err := getConfigPath()
	if err != nil {
		return err
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, data, 0644)
}
