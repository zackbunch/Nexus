package main

import (
	"github.com/BurntSushi/toml"
)

type Config struct {
	Repositories []Repository `toml:"repositories"`
}

type Repository struct {
	Name        string   `toml:"name"`
	Directories []string `toml:"directories"`
	LocalPath   string   `toml:"local_path"`
}

func LoadConfig(path string) (*Config, error) {
	var config Config
	_, err := toml.DecodeFile(path, &config)
	return &config, err
}
