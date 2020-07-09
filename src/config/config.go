package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

var Configs *Config

type MetaConfig struct {
	RunMode string `yaml:"run-mode"`
	Port    int32  `yaml:"port"`
	LogPath string `yaml:"log-path"`
}

type Config struct {
	Meta *MetaConfig `yaml:"meta"`
}

func Load(path string) error {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	cfg := &Config{}
	err = yaml.Unmarshal(f, &cfg)
	if err != nil {
		return err
	}

	Configs = cfg
	return nil
}
