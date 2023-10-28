package run

import (
	"ddns/internal"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	start := &cobra.Command{
		Use:   "run",
		Short: "Run A record synchronization once",
		Long:  `Run A record synchronization once`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return run()
		},
	}

	return start
}

func run() error {
	c := internal.GetConfig()

	// Start Synchronization
	i := internal.IPAddressProviderFactory(c)
	if i == nil {
		log.Fatal().Msgf("no IPAddressProvider was configured and enabled")
	}

	d := internal.DNSProviderFactory(c)
	if d == nil {
		log.Fatal().Msgf("no DNSProvider was configured and enabled")
	}

	if err := internal.SyncRecords(i, d)(); err != nil {
		log.Fatal().Msg(err.Error())
	}
	return nil
}
