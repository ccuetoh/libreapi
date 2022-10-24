package config

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"strconv"
)

type Option func(cfg *Config) *Config

func FromViper(v *viper.Viper) Option {
	return func(cfg *Config) *Config {
		err := v.Unmarshal(cfg)
		if err != nil {
			panic(errors.Wrap(err, "unable to unmarshall configuration"))
		}

		return cfg
	}
}

func SetNewRelicLicence(licence string) Option {
	return func(cfg *Config) *Config {
		cfg.NewRelic.Licence = licence
		return cfg
	}
}

func SetPort(port int) Option {
	return func(cfg *Config) *Config {
		cfg.HTTP.Port = strconv.Itoa(port)
		return cfg
	}
}
