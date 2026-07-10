package app

import (
	"os"

	"github.com/kotafan1rich/geo_logic_api_go/internal/api"
	"github.com/kotafan1rich/geo_logic_api_go/internal/config"
	"github.com/kotafan1rich/geo_logic_api_go/internal/database"
	"github.com/kotafan1rich/geo_logic_api_go/internal/handler/rent"
	"github.com/kotafan1rich/geo_logic_api_go/internal/handler/user"
	"github.com/kotafan1rich/geo_logic_api_go/internal/logger"
	"github.com/kotafan1rich/geo_logic_api_go/internal/repository"
	rentrepo "github.com/kotafan1rich/geo_logic_api_go/internal/repository/rent"
	userrepo "github.com/kotafan1rich/geo_logic_api_go/internal/repository/user"
	"github.com/kotafan1rich/geo_logic_api_go/internal/service"
	rentservice "github.com/kotafan1rich/geo_logic_api_go/internal/service/rent"
	userservice "github.com/kotafan1rich/geo_logic_api_go/internal/service/user"
)

type diContainer struct {
	db database.DB

	userRepo repository.UserRepository
	rentRepo repository.RentRepository

	userService service.UserService
	rentService service.RentService

	handler api.Handler

	log logger.Logger
}

func newDIContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) DB() database.DB {
	if d.db == nil {
		cfg := config.Get()

		db, err := database.New(cfg.Database.DSN(), cfg.Database.MaxIdleConns, cfg.Database.MaxOpenConns)

		if err != nil {
			d.Logger().Error("failed to connect to db", "err", err)
			os.Exit(1)
		}
		d.db = db
	}
	return d.db
}

func (d *diContainer) UserRepo() repository.UserRepository {
	if d.userRepo == nil {
		d.userRepo = userrepo.NewRepository(d.DB())
	}

	return d.userRepo
}

func (d *diContainer) RentRepo() repository.RentRepository {
	if d.rentRepo == nil {
		d.rentRepo = rentrepo.NewRepository(d.DB())
	}

	return d.rentRepo
}

func (d *diContainer) UserService() service.UserService {
	if d.userService == nil {
		d.userService = userservice.NewUserService(d.Logger(), d.UserRepo())
	}

	return d.userService
}

func (d *diContainer) RentService() service.RentService {
	if d.rentService == nil {
		d.rentService = rentservice.NewRentService(d.Logger(), d.RentRepo())
	}
	return d.rentService
}

func (d *diContainer) Handler() api.Handler {
	if d.handler == nil {
		userHandler := user.NewHandler(d.UserService(), d.Logger())
		rentHandler := rent.NewHandler(d.RentService(), d.Logger())
		d.handler = api.NewMainHandler(d.Logger(), userHandler, rentHandler)
	}

	return d.handler
}

func (d *diContainer) Logger() logger.Logger {
	if d.log == nil {
		cfg := config.Get()
		d.log = logger.New(cfg.Logging.LogLevel,
			cfg.Logging.LogFormat,
			cfg.Logging.LogAddSource,
			os.Stdout,
		)
	}
	return d.log
}
