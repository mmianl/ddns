package serve

import (
	"ddns/internal"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"net/http"
)

func New() *cobra.Command {
	start := &cobra.Command{
		Use:   "serve",
		Short: "Serve daemon that periodically performs A record synchronization",
		Long:  `Serve daemon that periodically performs A record synchronization`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return serve()
		},
	}

	return start
}

func serve() error {
	c := internal.GetConfig()

	// Initialize Metrics Handler
	if c.MetricsServerConfig.Enable {
		router := httprouter.New()
		router.GET("/metrics", internal.Metrics())

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
	i := internal.IPAddressProviderFactory(c)
	if i == nil {
		log.Fatal().Msgf("no IPAddressProvider was configured and enabled")
	}

	d := internal.DNSProviderFactory(c)
	if d == nil {
		log.Fatal().Msgf("no DNSProvider was configured and enabled")
	}

	internal.Retry(internal.SyncRecords(i, d), c.WaitInterval, c.RetryInterval)
	return nil
}
