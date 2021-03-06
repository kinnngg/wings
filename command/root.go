package command

import (
	"path/filepath"
	"strconv"

	"github.com/spf13/viper"

	"github.com/pterodactyl/wings/api"
	"github.com/pterodactyl/wings/config"
	"github.com/pterodactyl/wings/constants"
	"github.com/pterodactyl/wings/control"
	"github.com/pterodactyl/wings/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// RootCommand is the root command of wings
var RootCommand = &cobra.Command{
	Use:     "wings",
	Short:   "Wings is the next generation server control daemon for Pterodactyl",
	Long:    "Wings is the next generation server control daemon for Pterodactyl",
	Run:     run,
	Version: constants.Version,
}

var configPath string

func init() {
	RootCommand.Flags().StringVarP(&configPath, "config", "c", "./config.yml", "Allows to set the path of the configuration file.")
}

// Execute registers the RootCommand
func Execute() {
	RootCommand.Execute()
}

func run(cmd *cobra.Command, args []string) {
	utils.InitLogging()
	log.Info("Loading configuration...")
	if err := config.LoadConfiguration(configPath); err != nil {
		log.WithError(err).Fatal("Could not locate a suitable config.yml file. Aborting startup.")
		log.Exit(1)
	}
	utils.ConfigureLogging()

	log.Info(`                     ____`)
	log.Info(`__ Pterodactyl _____/___/_______ _______ ______`)
	log.Info(`\_____\    \/\/    /   /       /  __   /   ___/`)
	log.Info(`   \___\          /   /   /   /  /_/  /___   /`)
	log.Info(`        \___/\___/___/___/___/___    /______/`)
	log.Info(`                            /_______/ v` + constants.Version)
	log.Info()

	log.Info("Configuration loaded successfully.")

	log.Info("Loading configured servers.")
	if err := control.LoadServerConfigurations(filepath.Join(viper.GetString(config.DataPath), constants.ServersPath)); err != nil {
		log.WithError(err).Error("Failed to load configured servers.")
	}
	if amount := len(control.GetServers()); amount > 0 {
		log.Info("Loaded " + strconv.Itoa(amount) + " server(s).")
	}

	log.Info("Starting API server.")
	api := &api.InternalAPI{}
	api.Listen()
}
