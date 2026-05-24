package app

import (
	"log/slog"
	"os"

	"github.com/kotafan1rich/geo_logic_api_go/internal/api"
	"github.com/kotafan1rich/geo_logic_api_go/internal/config"
	"github.com/kotafan1rich/geo_logic_api_go/internal/database"
	"github.com/kotafan1rich/geo_logic_api_go/internal/repository"
	"github.com/kotafan1rich/geo_logic_api_go/internal/service"
)

type diContainer struct {
	db database.DB

	userRepo service.UserRepository

	userService service.UserService

	handler api.Handler
}

func newDIContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) DB() database.DB {
	if d.db == nil {
		cfg := config.Get()

		db, err := database.New(cfg.Database.DSN(), cfg.Database.MaxIdleConns, cfg.Database.MaxOpenConns)

		if err != nil {
			slog.Error("failed to connect to db", "err", err)
			os.Exit(1)
		}
		d.db = db
	}
	return d.db
}

func (d *diContainer) UserRepo() service.UserRepository {
	if d.userRepo == nil {
		d.userRepo = repository.NewRepository(d.DB())
	}

	return d.userRepo
}

func (d *diContainer) UserService() service.UserService {
	if d.userService == nil {
		d.userService = service.NewUserService(d.UserRepo())
	}

	return d.userService
}

func (d *diContainer) Handler() api.Handler {
	if d.handler == nil {
		d.handler = api.NewHandler()
	}

	return d.handler
}
