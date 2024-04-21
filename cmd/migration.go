package cmd

import (
	"context"
	"github.com/dwadp/attendance-api/store/db"
	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog/log"

	"github.com/spf13/cobra"
)

var migrationCmd = &cobra.Command{
	Use:   "migration",
	Short: "Run migration tools",
	Long:  `A list of tools that you can use to provision your database`,
	Run: func(cmd *cobra.Command, args []string) {
		migrationDirFlag, err := cmd.Flags().GetString("migration-dir")
		if err != nil {
			log.Fatal().Err(err).Send()
		}

		database, err := db.New(&cfg.Database)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to connect to database")
		}

		defer func() {
			if err := database.Close(); err != nil {
				log.Fatal().Err(err).Msg("Failed to close the database connection")
			}
		}()

		if len(args) < 1 {
			log.Fatal().Msg("Not enough arguments")
		}

		if err := goose.RunContext(context.TODO(), args[0], database, migrationDirFlag, args[1:]...); err != nil {
			log.Fatal().Err(err).Msg("Failed to run migration")
		}
	},
}

func init() {
	migrationCmd.PersistentFlags().StringP("migration-dir", "", "./store/migrations", "Set the migration dir")

	rootCmd.AddCommand(migrationCmd)
}
