package main

import (
	"flag"
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type IPAddressProvider interface {
	// Get the current ip address from provider
	GetIPAddress() (*string, error)
}

type DNSProvider interface {
	// Get the ip addresses which are currently set for the A records
	GetARecordAddresses() ([]RecordAddressMapping, error)

	// Set the current ip address for the provided A record
	SetARecordAddress(string, RecordAddressMapping) error
}

type Retryable func() error

type RecordAddressMapping struct {
	id        string
	aRecord   string
	ipAddress string
}

var Version = "development"

// Repeatedly run the Retryable function and wait between successful and failed attempts
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

// Updates the A records if required
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
			log.Info().Msgf("A record for %s is currently set to %s", addr.aRecord, addr.ipAddress)

			if *addressToSet != addr.ipAddress {
				log.Info().Msg("Ip address of A record did not match obtained address")
				ddnsDNSARecordUpdateTimeGauge.WithLabelValues(*addressToSet, addr.aRecord).SetToCurrentTime()

				if err := d.SetARecordAddress(*addressToSet, addr); err != nil {
					return err
				}
			} else {
				log.Info().Msgf("Ip address of A record matched obtained address, no update required")
			}
		}

		return nil
	}
}

func main() {
	// Initialize Logging
	logLevel := flag.String("logLevel", "info", "Log level, possible values: trace, debug, info, warn, error, fatal, panic")
	configPath := flag.String("config", "./config.yml", "Relative or absolute path to the config file")
	flag.Parse()
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	if l, err := zerolog.ParseLevel(*logLevel); err == nil {
		zerolog.SetGlobalLevel(l)
	} else {
		log.Error().Msgf("Could not parse log level '%s', defaulting to 'info'", *logLevel)
	}
	log.Info().Msgf("Running ddns version %s", Version)

	// Initialize Config
	c, err := NewConfig(*configPath)
	if err != nil {
		log.Fatal().Msgf("Error while loading the config file at %s: %s", *configPath, err)
	}

	// Initialize Metrics
	if c.MetricsServerConfig.Enable {
		ddnsVersionGauge.WithLabelValues(Version, runtime.Version()).Set(1)
		router := httprouter.New()
		router.GET("/metrics", Metrics())

		now := time.Now()
		ddnsStartTimeGauge.WithLabelValues().Set(float64(now.Unix()))

		ch := make(chan bool)
		go func() {
			listen := fmt.Sprintf("%s:%s", c.MetricsServerConfig.Host, c.MetricsServerConfig.Port)
			log.Info().Msgf("Metrics endpoint listening on %s...", listen)
			ch <- true
			if err := http.ListenAndServe(listen, router); err != nil {
				log.Fatal().Msgf("Could not listen on %s: %s", listen, err)
			}
		}()
		<-ch
	}

	// Start Main Loop
	i := IPAddressProviderFactory(c)
	if i == nil {
		log.Fatal().Msgf("No IPAddressProvider was configured and enabled")
	}

	d := DNSProviderFactory(c)
	if d == nil {
		log.Fatal().Msgf("No DNSProvider was configured and enabled")
	}

	Retry(SyncRecords(i, d), c.WaitInterval, c.RetryInterval)
}

// Returns an instrance of IPAddressProvider based on the passed configuration
func IPAddressProviderFactory(c *Config) IPAddressProvider {
	if c.StaticIPAddressProviderConfig.Enable {
		log.Debug().Msgf("Using StaticIPAddressProvider as IPAddressProvider")
		return NewStaticIPAddressProvider(&c.StaticIPAddressProviderConfig)
	} else if c.URLIPAddressProviderConfig.Enable {
		log.Debug().Msgf("Using URLIPAddressProvider as IPAddressProvider")
		return NewURLIPAddressProvider(&c.URLIPAddressProviderConfig)
	}

	return nil
}

// Returns an instrance of DNSProvider based on the passed configuration
func DNSProviderFactory(c *Config) DNSProvider {
	if c.CloudflareDNSProviderConfig.Enable {
		log.Debug().Msgf("Using CloudflareDNSProvider as DNSProvider")
		return NewCloudflareDNSProvider(&c.CloudflareDNSProviderConfig)
	}

	return nil
}
