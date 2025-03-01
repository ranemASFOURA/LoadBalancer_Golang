package config

import (
	"encoding/json"
	"os"
)

// Server struct represents an individual server
type Server struct {
	Name              string `json:"Name"`
	URL               string `json:"URL"`
	ActiveConnections int    `json:"ActiveConnections"`
	Healthy           bool   `json:"Healthy"`
}

// Config struct for JSON settings
type Config struct {
	HealthCheckInterval string   `json:"healthCheckInterval"`
	Servers             []Server `json:"servers"`
	ListenPort          string   `json:"listenPort"`
}

// LoadConfig loads configuration from JSON file
func LoadConfig(filename string) (Config, error) {
	var config Config

	file, err := os.ReadFile(filename)
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(file, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}
