package main

import (
	"github.com/rs/zerolog/log"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

// MetricsServerConfig Config section governing the metrics http server
type MetricsServerConfig struct {
	// Switch to turn on or off the http server that serves the metrics endpoint at /metrics
	Enable bool `yaml:"enable"`

	// Host for the metrics server to listen on
	Host string `yaml:"host"`

	// Port for the metrics server to listen on
	Port string `yaml:"port"`
}

// Config section governing the general configuration of the application
type Config struct {
	// Go duration to wait after a successful update attempt
	WaitInterval time.Duration `yaml:"waitInterval"`

	// Go duration to wait after a failed update attempt
	RetryInterval time.Duration `yaml:"retryInterval"`

	// Config section governing the metrics http server
	MetricsServerConfig MetricsServerConfig `yaml:"metricsServer"`

	// Config section governing the static ip address provider
	StaticIPAddressProviderConfig StaticIPAddressProviderConfig `yaml:"staticIPAddressProvider"`

	// Config section governing the url ip address provider
	URLIPAddressProviderConfig URLIPAddressProviderConfig `yaml:"urlIPAddressProvider"`

	// Config section governing the url ip address provider
	CloudflareDNSProviderConfig CloudflareDNSProviderConfig `yaml:"cloudflareDNSProvider"`
}

var defaultMetricsServerConfig = &MetricsServerConfig{
	Enable: true,
	Host:   "0.0.0.0",
	Port:   "9097",
}

var defaultConfig = &Config{
	WaitInterval:                  1 * time.Minute,
	RetryInterval:                 5 * time.Second,
	MetricsServerConfig:           *defaultMetricsServerConfig,
	URLIPAddressProviderConfig:    *defaultURLIPAddressProviderConfig,
	StaticIPAddressProviderConfig: *defaultStaticIPAddressProviderConfig,
	CloudflareDNSProviderConfig:   *defaultCloudflareDNSProviderConfig,
}

// NewConfig returns an instance of type *Config with values read from the passed config file
func NewConfig(configPath string) (*Config, error) {
	config := defaultConfig

	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Error().Msgf("error %s occurred while closing file %s", err, file.Name())
		}
	}(file)

	d := yaml.NewDecoder(file)

	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}
