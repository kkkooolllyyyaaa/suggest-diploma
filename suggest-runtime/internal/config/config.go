package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Server struct {
		Host string `json:"host"`
		Port int    `json:"port"`
	} `json:"server"`

	Pprof struct {
		Enable bool `json:"enable"`
	} `json:"pprof"`

	Redis struct {
		Host string `json:"host"`
	} `json:"redis"`

	Artifact struct {
		Queries           string `json:"queries"`
		QueriesCategories string `json:"queries_categories"`
		Nodes             string `json:"nodes"`
	} `json:"artifact"`
}

func LoadConfig(path string) (*Config, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(file, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
