package internal

import (
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v3"
	"os"
	"time"
)

var global Config

// MetricsServerConfig Config section governing the metrics http server
type MetricsServerConfig struct {
	// Switch to turn on or off the http server that serves the metrics endpoint at /metrics
	Enable bool `yaml:"enable" envconfig:"DDNS_METRICS_ENABLE" required:"false"`

	// Host for the metrics server to listen on
	Host string `yaml:"host" envconfig:"DDNS_METRICS_HOST" required:"false"`

	// Port for the metrics server to listen on
	Port string `yaml:"port" envconfig:"DDNS_METRICS_PORT" required:"false"`
}

// Config section governing the general configuration of the application
type Config struct {
	// Go duration to wait after a successful update attempt
	WaitInterval time.Duration `yaml:"waitInterval" envconfig:"DDNS_WAIT_INTERVAL" required:"false"`

	// Go duration to wait after a failed update attempt
	RetryInterval time.Duration `yaml:"retryInterval" envconfig:"DDNS_RETRY_INTERVAL" required:"false"`

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

var defaultConfig = Config{
	WaitInterval:                  1 * time.Minute,
	RetryInterval:                 5 * time.Second,
	MetricsServerConfig:           *defaultMetricsServerConfig,
	URLIPAddressProviderConfig:    *defaultURLIPAddressProviderConfig,
	StaticIPAddressProviderConfig: *defaultStaticIPAddressProviderConfig,
	CloudflareDNSProviderConfig:   *defaultCloudflareDNSProviderConfig,
}

// GatherConfig sets the globalConfig with values read from the passed config file and the environment
func GatherConfig(configPath string) error {
	global = defaultConfig
	if err := gatherFromFile(configPath); err != nil {
		return err
	}
	return gatherFromEnv()
}

func gatherFromFile(configPath string) error {
	file, err := os.Open(configPath)
	if err != nil {
		return err
	}

	defer file.Close()
	d := yaml.NewDecoder(file)

	return d.Decode(&global)
}

func gatherFromEnv() error {
	return envconfig.Process("", &global)
}

// GetConfig returns the globalConfig
func GetConfig() *Config {
	return &global
}
