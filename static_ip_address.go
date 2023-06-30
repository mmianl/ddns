package main

type StaticIPAddressProviderConfig struct {
	// Switch to enable or disable this provider
	Enable bool `yaml:"enable"`

	// Static ip address returned by this provider
	Address string `yaml:"address"`
}

var defaultStaticIPAddressProviderConfig = &StaticIPAddressProviderConfig{
	Enable:  false,
	Address: "127.0.0.1",
}

type StaticIPAddressProvider struct {
	address string
}

// NewStaticIPAddressProvider Returns an instance of StaticIPAddressProvider based on the passed configuration
func NewStaticIPAddressProvider(config *StaticIPAddressProviderConfig) *StaticIPAddressProvider {
	return &StaticIPAddressProvider{
		address: config.Address,
	}
}

// GetIPAddress Returns the static ip address passed via configuration
func (s *StaticIPAddressProvider) GetIPAddress() (*string, error) {
	return &s.address, nil
}
