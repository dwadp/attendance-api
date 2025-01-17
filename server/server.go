package server

import (
	"database/sql"
	"github.com/dwadp/attendance-api/config"
	"github.com/dwadp/attendance-api/server/handlers"
	"github.com/dwadp/attendance-api/server/validator"
	"github.com/dwadp/attendance-api/store"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
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
	users := map[string]string{}
	users[s.Config.Auth.User] = s.Config.Auth.Pass

	s.App.Use(basicauth.New(basicauth.Config{
		Users: users,
	}))

	handlers.RegisterEmployee(s.App.Group("employees"), s.Store, s.Validator)
	handlers.RegisterShift(s.App.Group("shifts"), s.Store, s.Validator)
	handlers.RegisterDayOff(s.App.Group("day-offs"), s.Store, s.Validator)
	handlers.RegisterEmployeeShift(s.App.Group("employee-shifts"), s.Store, s.Validator)
	handlers.RegisterAttendance(s.App.Group("attendances"), s.Store, s.Validator)
	handlers.RegisterHolidayHandlers(s.App.Group("holidays"), s.Store, s.Validator)

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
