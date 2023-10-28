package main

import (
	"ddns/main/v2/internal"
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
	// GetIPAddress Get the current ip address from provider
	GetIPAddress() (*string, error)
}

type DNSProvider interface {
	// GetARecordAddresses Get the ip addresses which are currently set for the A records
	GetARecordAddresses() ([]internal.RecordAddressMapping, error)

	// SetARecordAddress Set the current ip address for the provided A record
	SetARecordAddress(string, internal.RecordAddressMapping) error
}

type Retryable func() error

var Version = "development"

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
				internal.DNSARecordInfoGauge.WithLabelValues(*addressToSet, addr.ARecord).Set(1)
			} else {
				log.Info().Msgf("Ip address of A record matched obtained address, no update required")
				internal.DNSARecordInfoGauge.WithLabelValues(addr.IPAddress, addr.ARecord).Set(1)
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
	err := internal.GatherConfig(*configPath)
	if err != nil {
		log.Fatal().Msgf("Error while loading the config file at %s: %s", *configPath, err)
	}

	c := internal.GetConfig()

	// Initialize Metrics
	if c.MetricsServerConfig.Enable {
		internal.VersionGauge.WithLabelValues(Version, runtime.Version()).Set(1)
		router := httprouter.New()
		router.GET("/metrics", internal.Metrics())

		now := time.Now()
		internal.StartTimeGauge.WithLabelValues().Set(float64(now.Unix()))

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

// IPAddressProviderFactory Returns an instance of IPAddressProvider based on the passed configuration
func IPAddressProviderFactory(c *internal.Config) IPAddressProvider {
	if c.StaticIPAddressProviderConfig.Enable {
		log.Debug().Msgf("Using StaticIPAddressProvider as IPAddressProvider")
		return internal.NewStaticIPAddressProvider(&c.StaticIPAddressProviderConfig)
	} else if c.URLIPAddressProviderConfig.Enable {
		log.Debug().Msgf("Using URLIPAddressProvider as IPAddressProvider")
		return internal.NewURLIPAddressProvider(&c.URLIPAddressProviderConfig)
	}

	return nil
}

// DNSProviderFactory Returns an instance of DNSProvider based on the passed configuration
func DNSProviderFactory(c *internal.Config) DNSProvider {
	if c.CloudflareDNSProviderConfig.Enable {
		log.Debug().Msgf("Using CloudflareDNSProvider as DNSProvider")
		return internal.NewCloudflareDNSProvider(&c.CloudflareDNSProviderConfig)
	}

	return nil
}
