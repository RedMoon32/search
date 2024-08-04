package utils

import (
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Kafka struct {
		Brokers string `yaml:"brokers"`
		Topic   string `yaml:"topic"`
	} `yaml:"kafka"`

	KVStore struct {
		Path string `yaml:"path"`
	} `yaml:"kvstore"`

	Storage struct {
		Path string `yaml:"path"`
	} `yaml:"storage"`
}

func LoadConfig() *Config {
	data, err := os.ReadFile("/home/red/projects/search/web-crawler/config/config.yaml")
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("Failed to parse config file: %v", err)
	}

	return &config
}

func FetchURL(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func CurrentTimestamp() int64 {
	return time.Now().Unix()
}

func ReadFile(filename string) (string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
