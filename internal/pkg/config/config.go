package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	Meta *MetaConfig `yaml:"meta"`
}

type MetaConfig struct {
	Port      int32  `yaml:"port"`
	RunMode   string `yaml:"run-mode"`
	LogName   string `yaml:"log-name"`
	BucketCap int64  `yaml:"bucket-cap"`
	BucketQua int64  `yaml:"bucket-qua"`
}

func Load(path string) (*Config, error) {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	cfg := &Config{}
	err = yaml.Unmarshal(f, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
