package cmd

import (
	"ddns/cmd/run"
	"ddns/cmd/serve"
	"ddns/internal"
	"errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"io"
	"runtime"
	"time"
)

func New(_ io.Writer, _ io.Reader, _ []string, version string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ddns",
		Short: "The DDNS CLI lets you interact with the DDNS service",
		Long:  `The DDNS CLI lets you interact with the DDNS service`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return errors.New("no additional command provided")
		},
		Version: version,
	}

	var logLevel, configPath string
	cobra.OnInitialize(initConfig(&version, &logLevel, &configPath))
	cmd.PersistentFlags().StringVar(&logLevel, "loglevel", "info", "log level, possible values: trace, debug, info, warn, error, fatal, panic")
	cmd.PersistentFlags().StringVar(&configPath, "config", "./config.yml", "relative or absolute path to the config file")
	cmd.InitDefaultVersionFlag()

	cmd.AddCommand(
		serve.New(),
		run.New(),
	)

	return cmd
}

func initConfig(version, logLevel, configPath *string) func() {

	return func() {
		// Initialize Logging
		zerolog.SetGlobalLevel(zerolog.InfoLevel)

		if l, err := zerolog.ParseLevel(*logLevel); err == nil {
			zerolog.SetGlobalLevel(l)
		} else {
			log.Error().Msgf("Could not parse log level '%s', defaulting to 'info'", *logLevel)
		}
		log.Info().Msgf("Running DDNS version %s", *version)

		// Initialize Config
		err := internal.GatherConfig(*configPath)
		if err != nil {
			log.Fatal().Msgf("Error while loading the config file at %s: %s", *configPath, err)
		}

		// Initialize Metric
		internal.VersionGauge.WithLabelValues(*version, runtime.Version()).Set(1)
		now := time.Now()
		internal.StartTimeGauge.WithLabelValues().Set(float64(now.Unix()))
	}
}
