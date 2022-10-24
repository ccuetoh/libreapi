package config

type NewRelic struct {
	Licence              string `mapstructure:"licence"`
	AppName              string `mapstructure:"app_name,omitempty"`
	LogForwardingEnabled bool   `mapstructure:"log_forwarding_enabled"`
}

type HTTP struct {
	Address      string `mapstructure:"address,omitempty"`
	Port         string `mapstructure:"port,omitempty"`
	DebugEnabled bool   `mapstructure:"debug_enabled"`
}

type Config struct {
	NewRelic NewRelic ``
	HTTP     HTTP
}

func Default() *Config {
	return &Config{
		NewRelic: NewRelic{
			AppName:              "LibreAPI",
			LogForwardingEnabled: true,
		},
		HTTP: HTTP{
			Port:         "443",
			DebugEnabled: false,
		},
	}
}

func Apply(cfg *Config, opts ...Option) *Config {
	for _, op := range opts {
		cfg = op(cfg)
	}

	return cfg
}

func Build(opts ...Option) *Config {
	return Apply(Default(), opts...)
}
