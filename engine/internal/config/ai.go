package config

type AIConfig struct {
	Mode   AIMode      `yaml:"mode"`
	Local  LocalAI     `yaml:"local"`
	Remote RemoteAI    `yaml:"remote"`
}

type AIMode string

const (
	AIModeLocal  AIMode = "local"
	AIModeRemote AIMode = "remote"
	AIModeHybrid AIMode = "hybrid"
)

type LocalAI struct {
	Provider string `yaml:"provider"`
	Endpoint string `yaml:"endpoint"`
	Model    string `yaml:"model"`
}

type RemoteAI struct {
	Provider string `yaml:"provider"`
	APIKey   string `yaml:"api_key"`
	BaseURL  string `yaml:"base_url"`
	Model    string `yaml:"model"`
}
