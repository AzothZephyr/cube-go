package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

// NewConfig returns a new decoded Config struct
func NewConfig(configPath string) (*Config, error) {
	// Create config structure
	config := &Config{}

	// Open config file
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Init new YAML decode
	d := yaml.NewDecoder(file)

	// Start YAML decoding from file
	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}

// Config struct for webapp config
type Config struct {
	Client struct {
		// Host is the local machine IP Address to bind the HTTP Server to
		ApiKey string `yaml:"api_key"`

		// Port is the local machine TCP Port to bind the HTTP Server to
		SecretKey string `yaml:"secret_key"`
	} `yaml:"client"`
}

// Run will run the HTTP Server, and the opportunity engine
func (config *Config) Run() {

}
