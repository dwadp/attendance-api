package cmd

import (
	"github.com/dwadp/attendance-api/config"
	"github.com/dwadp/attendance-api/server"
	"github.com/dwadp/attendance-api/server/validator"
	"github.com/dwadp/attendance-api/store/db"
	"github.com/dwadp/attendance-api/store/postgres"
	"github.com/rs/zerolog/log"
	"os"
	"path"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	cfg     config.Config
)

var rootCmd = &cobra.Command{
	Use:   "attendance-api",
	Short: "Attendance API",
	Long:  `An API to manage employee attendance`,
	Run: func(cmd *cobra.Command, args []string) {
		database, err := db.New(&cfg.Database)
		if err != nil {
			log.Fatal().Err(err).Msg("Unable to connect to the database")
		}

		s := postgres.NewPostgres(database)
		v, err := validator.New()
		if err != nil {
			log.Fatal().Err(err).Msg("Unable to create validator")
		}

		srv := server.New(&cfg, s, database, v)
		if err := srv.Start(); err != nil {
			log.Fatal().Err(err).Msg("Failed to start the server")
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "Set the config file (default is $HOME/attendance-api/config.yml)")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(path.Join(home, "attendance-api"))
		viper.SetConfigType("yml")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal().Err(err).Str("file_name", viper.ConfigFileUsed()).Msg("Unable to read config file")
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatal().Err(err).Msg("Unable to decode the configuration")
	}
}
