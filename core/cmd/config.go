package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kunlun-sec/lunying/pkg/logger"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	configShow   bool
	configSet    string
	configValue  string
	configPath   string
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "配置管理",
	Long: `管理猎影CLI配置文件，包括：
  - 查看当前配置
  - 设置配置项
  - 初始化配置`,
	Run: func(cmd *cobra.Command, args []string) {
		if configShow {
			showConfig()
			return
		}

		if configSet != "" && configValue != "" {
			setConfig(configSet, configValue)
			return
		}

		// 默认显示帮助
		cmd.Help()
	},
}

var configInitCmd = &cobra.Command{
	Use:   "init",
	Short: "初始化配置文件",
	Run: func(cmd *cobra.Command, args []string) {
		initConfig()
	},
}

func init() {
	configCmd.Flags().BoolVarP(&configShow, "show", "s", false, "显示当前配置")
	configCmd.Flags().StringVarP(&configSet, "set", "", "", "设置配置项 (如: api.key)")
	configCmd.Flags().StringVarP(&configValue, "value", "v", "", "配置值")
	configCmd.Flags().StringVarP(&configPath, "path", "p", "", "配置文件路径")

	configCmd.AddCommand(configInitCmd)
}

type CLIConfig struct {
	API      APIConfig      `yaml:"api"`
	Network  NetworkConfig  `yaml:"network"`
	Logging  LoggingConfig  `yaml:"logging"`
	Database DatabaseConfig `yaml:"database"`
}

type APIConfig struct {
	Key     string `yaml:"key"`
	BaseURL string `yaml:"base_url"`
}

type NetworkConfig struct {
	Timeout         int    `yaml:"timeout"`
	Concurrency     int    `yaml:"concurrency"`
	RateLimit       int    `yaml:"rate_limit"`
	UserAgent       string `yaml:"user_agent"`
	FollowRedirects bool   `yaml:"follow_redirects"`
	Proxy           string `yaml:"proxy"`
}

type LoggingConfig struct {
	Level  string `yaml:"level"`
	Path   string `yaml:"path"`
	Output string `yaml:"output"`
}

type DatabaseConfig struct {
	Path string `yaml:"path"`
	Type string `yaml:"type"`
}

func getConfigFilePath() string {
	if configPath != "" {
		return configPath
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return ".lieying/config.yaml"
	}
	return filepath.Join(home, ".lieying", "config.yaml")
}

func loadConfig() (*CLIConfig, error) {
	path := getConfigFilePath()

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config CLIConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func saveConfig(config *CLIConfig) error {
	path := getConfigFilePath()

	// 确保目录存在
	dir := filepath.Dir(path)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

func getDefaultConfig() *CLIConfig {
	return &CLIConfig{
		API: APIConfig{
			Key:     "",
			BaseURL: "https://api.deepseek.com",
		},
		Network: NetworkConfig{
			Timeout:         30,
			Concurrency:     10,
			RateLimit:       50,
			UserAgent:       "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
			FollowRedirects: true,
			Proxy:           "",
		},
		Logging: LoggingConfig{
			Level:  "info",
			Path:   "logs",
			Output: "",
		},
		Database: DatabaseConfig{
			Path: "data/lieying.db",
			Type: "sqlite",
		},
	}
}

func showConfig() {
	config, err := loadConfig()
	if err != nil {
		fmt.Printf("无法加载配置: %v\n", err)
		fmt.Println("请运行 'lieying config init' 初始化配置")
		return
	}

	fmt.Println("=====================================")
	fmt.Println("  猎影CLI配置")
	fmt.Println("=====================================")
	fmt.Printf("配置文件: %s\n\n", getConfigFilePath())

	fmt.Println("API配置:")
	fmt.Printf("  Base URL: %s\n", config.API.BaseURL)
	if config.API.Key != "" {
		fmt.Printf("  API Key: (已设置)\n")
	} else {
		fmt.Printf("  API Key: (未设置)\n")
	}

	fmt.Println("\n网络配置:")
	fmt.Printf("  Timeout: %d秒\n", config.Network.Timeout)
	fmt.Printf("  Concurrency: %d\n", config.Network.Concurrency)
	fmt.Printf("  Rate Limit: %d QPS\n", config.Network.RateLimit)
	fmt.Printf("  User Agent: %s\n", config.Network.UserAgent)
	fmt.Printf("  Follow Redirects: %v\n", config.Network.FollowRedirects)
	if config.Network.Proxy != "" {
		fmt.Printf("  Proxy: %s\n", config.Network.Proxy)
	}

	fmt.Println("\n日志配置:")
	fmt.Printf("  Level: %s\n", config.Logging.Level)
	fmt.Printf("  Path: %s\n", config.Logging.Path)

	fmt.Println("\n数据库配置:")
	fmt.Printf("  Type: %s\n", config.Database.Type)
	fmt.Printf("  Path: %s\n", config.Database.Path)

	fmt.Println("=====================================")
}

func setConfig(key, value string) {
	config, err := loadConfig()
	if err != nil {
		config = getDefaultConfig()
	}

	// 根据key设置对应的配置项
	switch key {
	case "api.key":
		config.API.Key = value
	case "api.base_url":
		config.API.BaseURL = value
	case "network.timeout":
		config.Network.Timeout = parseInt(value)
	case "network.concurrency":
		config.Network.Concurrency = parseInt(value)
	case "network.rate_limit":
		config.Network.RateLimit = parseInt(value)
	case "network.proxy":
		config.Network.Proxy = value
	case "network.user_agent":
		config.Network.UserAgent = value
	case "logging.level":
		config.Logging.Level = value
	case "logging.path":
		config.Logging.Path = value
	case "database.path":
		config.Database.Path = value
	default:
		fmt.Printf("未知的配置项: %s\n", key)
		return
	}

	err = saveConfig(config)
	if err != nil {
		fmt.Printf("保存配置失败: %v\n", err)
		return
	}

	fmt.Printf("已设置 %s = %s\n", key, value)
}

func initConfig() {
	path := getConfigFilePath()

	// 检查是否已存在
	if _, err := os.Stat(path); err == nil {
		fmt.Printf("配置文件已存在: %s\n", path)
		fmt.Println("如需重新初始化，请先删除现有配置文件")
		return
	}

	config := getDefaultConfig()
	err := saveConfig(config)
	if err != nil {
		fmt.Printf("初始化配置失败: %v\n", err)
		return
	}

	// 初始化日志
	if err := logger.Init(config.Logging.Level, config.Logging.Path, config.Logging.Output); err != nil {
		fmt.Printf("日志初始化失败: %v\n", err)
	}

	fmt.Println("=====================================")
	fmt.Println("  配置初始化完成")
	fmt.Println("=====================================")
	fmt.Printf("配置文件: %s\n", path)
	fmt.Println("\n可以使用以下命令修改配置:")
	fmt.Println("  lieying config --set api.key --value YOUR_KEY")
	fmt.Println("  lieying config --show")
	fmt.Println("=====================================")
}

func parseInt(s string) int {
	var result int
	fmt.Sscanf(s, "%d", &result)
	return result
}