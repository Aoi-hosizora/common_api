package config

import (
	"github.com/Aoi-hosizora/ahlib-web/xvalidator"
	"github.com/Aoi-hosizora/ahlib/xreflect"
	"gopkg.in/yaml.v2"
	"os"
)

type Config struct {
	Meta *MetaConfig `yaml:"meta"  validate:"required"`
}

type MetaConfig struct {
	Port    int32  `yaml:"port"     validate:"required"`
	RunMode string `yaml:"run-mode" default:"debug"`
	LogName string `yaml:"log-name" default:"./logs/console"`
	Swagger bool   `yaml:"swagger"`
	Host    string `yaml:"host"`

	BucketCap int64 `yaml:"bucket-cap" default:"200" validate:"gt=0"`
	BucketQua int64 `yaml:"bucket-qua" default:"100" validate:"gt=0"`
}

var _debugMode = true

func IsDebugMode() bool {
	return _debugMode
}

func Load(path string) (*Config, error) {
	f, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	cfg := &Config{}
	if err = yaml.Unmarshal(f, cfg); err != nil {
		return nil, err
	}
	if _, err = xreflect.FillDefaultFields(cfg); err != nil {
		return nil, err
	}
	if err = validateConfig(cfg); err != nil {
		return nil, err
	}

	_debugMode = cfg.Meta.RunMode == "debug"
	return cfg, nil
}

func validateConfig(cfg *Config) error {
	val := xvalidator.NewCustomStructValidator()
	val.SetValidatorTagName("validate")
	val.SetMessageTagName("message")
	xvalidator.UseTagAsFieldName(val.ValidateEngine(), "yaml")
	err := val.ValidateStruct(cfg)
	if err != nil {
		ut, _ := xvalidator.ApplyTranslator(val.ValidateEngine(), xvalidator.EnLocaleTranslator(), xvalidator.EnTranslationRegisterFunc())
		return xvalidator.FlattedMapToError(err.(*xvalidator.ValidateFieldsError).Translate(ut, false))
	}
	return nil
}
