package server

import (
	"database/sql"
	"github.com/dwadp/attendance-api/config"
	"github.com/dwadp/attendance-api/server/handlers"
	"github.com/dwadp/attendance-api/server/validator"
	"github.com/dwadp/attendance-api/store"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

type Server struct {
	App       *fiber.App
	Config    *config.Config
	Store     store.Store
	DB        *sql.DB
	Validator *validator.Validator
}

func New(cfg *config.Config, s store.Store, db *sql.DB, validator *validator.Validator) *Server {
	return &Server{
		App:       fiber.New(),
		Config:    cfg,
		Store:     s,
		DB:        db,
		Validator: validator,
	}
}

func (s *Server) Start() error {
	handlers.RegisterEmployee(s.App.Group("employees"), s.Store, s.Validator)

	go func() {
		if err := s.App.Listen(net.JoinHostPort(s.Config.Server.Host, s.Config.Server.Port)); err != nil {
			log.Fatal().Err(err).Msg("Failed to start server")
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	_ = <-c
	log.Debug().Msg("Shutting down the server")

	if err := s.App.Shutdown(); err != nil {
		log.Debug().Err(err).Msg("Failed to shutdown server")
		return err
	}

	log.Debug().Msg("Running cleanup tasks")

	if err := s.DB.Close(); err != nil {
		log.Fatal().Err(err).Msg("Failed to close database")
	}

	return nil
}
