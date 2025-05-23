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
		QueriesVectors    string `json:"queries_vectors"`
		TokensVectors     string `json:"tokens_vectors"`
		AnnoyIndex        string `json:"annoy_index"`
	} `json:"artifact"`

	ArtifactRemote struct {
		Queries           string `json:"queries"`
		QueriesCategories string `json:"queries_categories"`
		Nodes             string `json:"nodes"`
		QueriesVectors    string `json:"queries_vectors"`
		TokensVectors     string `json:"tokens_vectors"`
		AnnoyIndex        string `json:"annoy_index"`
	} `json:"artifact_remote"`

	CategoryEngine struct {
		Threshold float64 `json:"threshold"`
	} `json:"category_engine"`

	S3 struct {
		Endpoint   string `json:"endpoint,omitempty"`
		AccessKey  string `json:"access_key,omitempty"`
		SecretKey  string `json:"secret_key,omitempty"`
		BucketName string `json:"bucket_name,omitempty"`
		UseSSL     bool   `json:"use_ssl,omitempty"`
	} `json:"s3"`

	Vector struct {
		Dimension int     `json:"dim"`
		Count     int     `json:"count"`
		MinDist   float32 `json:"min_dist"`
	} `json:"vector"`
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
