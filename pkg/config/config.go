package config

type Config struct {
	NewRelic NewRelic `mapstructure:"new_relic"`
	HTTP     HTTP     `mapstructure:"http"`
}

type NewRelic struct {
	Licence              string `mapstructure:"licence"`
	AppName              string `mapstructure:"app_name"`
	LogForwardingEnabled bool   `mapstructure:"log_forwarding_enabled"`
}

type HTTP struct {
	Address             string `mapstructure:"address"`
	Port                string `mapstructure:"port"`
	DebugEnabled        bool   `mapstructure:"debug_enabled"`
	ProxyClientIPHeader string `mapstructure:"proxy_client_ip_header"`
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
