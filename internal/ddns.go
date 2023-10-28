package internal

import (
	"github.com/rs/zerolog/log"
	"strings"
	"time"
)

type RecordAddressMapping struct {
	ID        string
	ARecord   string
	IPAddress string
}

type IPAddressProvider interface {
	// GetIPAddress Get the current ip address from provider
	GetIPAddress() (*string, error)
}

type DNSProvider interface {
	// GetARecordAddresses Get the ip addresses which are currently set for the A records
	GetARecordAddresses() ([]RecordAddressMapping, error)

	// SetARecordAddress Set the current ip address for the provided A record
	SetARecordAddress(string, RecordAddressMapping) error
}

type Retryable func() error

// Retry Repeatedly run the Retryable function and wait between successful and failed attempts
func Retry(retryable Retryable, waitInterval time.Duration, retryInterval time.Duration) {
	for {
		err := retryable()
		if err == nil {
			log.Info().Msgf("Success. Next attempt in %s", waitInterval)
			time.Sleep(waitInterval)
		} else {
			log.Error().Msgf("An error occurred: %s. Retrying in %s", err, retryInterval)
			time.Sleep(retryInterval)
		}
	}
}

// SyncRecords Updates the A records if required
func SyncRecords(i IPAddressProvider, d DNSProvider) func() error {
	return func() error {
		addressToSet, err := i.GetIPAddress()
		if err != nil {
			return err
		}
		log.Info().Msgf("Obtained ip address was %s", *addressToSet)

		setAddresses, err := d.GetARecordAddresses()
		if err != nil {
			return err
		}

		for _, addr := range setAddresses {
			log.Info().Msgf("A record for %s is currently set to %s", addr.ARecord, addr.IPAddress)

			if *addressToSet != addr.IPAddress {
				log.Info().Msg("Ip address of A record did not match obtained address")
				if err := d.SetARecordAddress(*addressToSet, addr); err != nil {
					return err
				}
				DNSARecordInfoGauge.WithLabelValues(*addressToSet, addr.ARecord).Set(1)
			} else {
				log.Info().Msgf("Ip address of A record matched obtained address, no update required")
				DNSARecordInfoGauge.WithLabelValues(addr.IPAddress, addr.ARecord).Set(1)
			}
		}

		return nil
	}
}

// IPAddressProviderFactory Returns an instance of IPAddressProvider based on the passed configuration
func IPAddressProviderFactory(c *Config) IPAddressProvider {
	if c.StaticIPAddressProviderConfig.Enable {
		log.Debug().Msgf("Using StaticIPAddressProvider as IPAddressProvider with ip address %s", c.StaticIPAddressProviderConfig.Address)
		return NewStaticIPAddressProvider(&c.StaticIPAddressProviderConfig)
	} else if c.URLIPAddressProviderConfig.Enable {
		log.Debug().Msgf("Using URLIPAddressProvider as IPAddressProvider with url %s", c.URLIPAddressProviderConfig.URL)
		return NewURLIPAddressProvider(&c.URLIPAddressProviderConfig)
	}

	return nil
}

// DNSProviderFactory Returns an instance of DNSProvider based on the passed configuration
func DNSProviderFactory(c *Config) DNSProvider {
	if c.CloudflareDNSProviderConfig.Enable {
		log.Debug().Msgf("Using CloudflareDNSProvider as DNSProvider with records %s", strings.Join(c.CloudflareDNSProviderConfig.ARecords, ","))
		return NewCloudflareDNSProvider(&c.CloudflareDNSProviderConfig)
	}

	return nil
}
