package config

type NetworkConfig struct {
	Mode         NetworkMode `yaml:"mode"`
	Proxy        ProxyConfig `yaml:"proxy"`
	AllowedHosts []string    `yaml:"allowed_hosts"`
	BlockedHosts []string    `yaml:"blocked_hosts"`
}

type NetworkMode string

const (
	NetworkModeOffline  NetworkMode = "offline"
	NetworkModeProxy    NetworkMode = "proxy"
	NetworkModeDirect   NetworkMode = "direct"
)

type ProxyConfig struct {
	Enabled  bool   `yaml:"enabled"`
	Type     string `yaml:"type"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}
