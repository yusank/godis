package config

import (
	"io"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	DefaultTimeout time.Duration `yaml:"default_timeout"`
	Persistence    string        `yaml:"persistence"`
}

func InitConfig(r io.Reader) error {
	d := yaml.NewDecoder(r)

	tmp := new(Config)
	if err := d.Decode(tmp); err != nil {
		return err
	}

	cfg = tmp
	return nil
}

var cfg *Config

func GetDefaultTimeout() time.Duration {
	if cfg == nil || cfg.DefaultTimeout == 0 {
		return time.Second * 5
	}

	return cfg.DefaultTimeout
}
