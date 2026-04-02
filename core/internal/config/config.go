package config

import (
	"os"
	"path/filepath"
)

type Config struct {
	DatabasePath string
	LogPath    string
	ListenAddr string
}

func DefaultConfig() *Config {
	wd, _ := os.Getwd()
	lieyingDir := filepath.Join(wd, "data")

	os.MkdirAll(lieyingDir, 0755)
	os.MkdirAll(filepath.Join(lieyingDir, "logs"), 0755)

	return &Config{
		DatabasePath: filepath.Join(lieyingDir, "lieying.db"),
		LogPath:    filepath.Join(lieyingDir, "logs"),
		ListenAddr: ":8081",
	}
}
