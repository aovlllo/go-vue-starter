package server

import (
	"context"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/aovlllo/vue-template/pkg/api"
	"github.com/aovlllo/vue-template/pkg/app"
	"github.com/aovlllo/vue-template/pkg/db"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
)

// Config represents the server configuration
type Config struct {
	ListenAddress string `yaml:"listen_address"`

	API  *api.Config  `yaml:"api"`
	App  *app.Config  `yaml:"app"`
	DB   *db.Config   `yaml:"database"`
}

// Instance represents an instance of the server
type Instance struct {
	API    *api.API
	App    *app.App
	Config *Config
	DB     db.DB

	httpServer *http.Server
}

// NewInstance returns an new instance of our server
func NewInstance() *Instance {
	return &Instance{}
}

// Start starts the server
func (i *Instance) Start(file string) {
	var err error
	var router = mux.NewRouter()

	// Load configuration file
	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal().Err(err).Msg("Could not load configuration")
	}

	err = yaml.Unmarshal(data, &i.Config)
	if err != nil {
		log.Fatal().Err(err).Msg("Could not load configuration")
	}
	i.Config.DB.MySQL.User = os.Getenv("DB_USER")
	i.Config.DB.MySQL.Password = os.Getenv("DB_PASSWORD")

	// Establish database connection
	i.DB, err = db.NewConnection(i.Config.DB)
	if err != nil {
		log.Fatal().Err(err).Msg("Could not open database connection")
	}
	defer i.DB.CloseConnection()
	log.Debug().Msg("Successfully initiated DB connection")

	i.API, err = api.New(i.Config.API, i.DB, router)
	if err != nil {
		log.Fatal().Err(err).Msg("Could not create API instance")
	}

	i.App, err = app.New(i.Config.App, router)
	if err != nil {
		log.Fatal().Err(err).Msg("Could not create app instance")
	}

	// Startup the HTTP Server in a way that we can gracefully shut it down again
	i.httpServer = &http.Server{
		Addr:    i.Config.ListenAddress,
		Handler: router,
	}

	err = i.httpServer.ListenAndServe()
	if err != http.ErrServerClosed {
		log.Error().Err(err).Msg("HTTP Server stopped unexpected")
		i.Shutdown()
	}
	log.Info().Err(err).Msg("HTTP Server stopped")
}

// Shutdown stops the server
func (i *Instance) Shutdown() {
	// Shutdown all dependencies
	i.DB.CloseConnection()

	// Shutdown HTTP server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := i.httpServer.Shutdown(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to shutdown HTTP server gracefully")
		os.Exit(1)
	}

	log.Info().Msg("Shutdown HTTP server...")
	os.Exit(0)
}
