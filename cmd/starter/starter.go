package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/aovlllo/vue-template/pkg/server"
	"github.com/aovlllo/vue-template/pkg/version"

	//"github.com/sirupsen/logrus"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	configFileFlag = flag.String("config.file", "config.yml", "Path to the configuration file.")
	versionFlag    = flag.Bool("version", false, "Show version information.")
	debugFlag      = flag.Bool("debug", false, "Show debug information.")
)

func init() {
	// Parse command-line flags
	flag.Parse()

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	// Log settings
	if *debugFlag {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
}

func main() {
	// Show version information
	if *versionFlag {
		fmt.Fprintln(os.Stdout, version.Print("starter"))
		os.Exit(0)
	}

	// Create server instance
	instance := server.NewInstance()

	// Interrupt handler
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		log.Info().Msgf("Received %s signal", <-c)
		instance.Shutdown()
	}()

	// Start server
	log.Info().Msgf("Starting starter %s", version.Info())
	log.Info().Msgf("Build context %s", version.BuildContext())
	instance.Start(*configFileFlag)
}
