package app

import (
	"os"

	"github.com/kotafan1rich/geo_logic_api_go/internal/api"
	"github.com/kotafan1rich/geo_logic_api_go/internal/config"
	"github.com/kotafan1rich/geo_logic_api_go/internal/database"
	"github.com/kotafan1rich/geo_logic_api_go/internal/handler/event"
	"github.com/kotafan1rich/geo_logic_api_go/internal/handler/infra"
	"github.com/kotafan1rich/geo_logic_api_go/internal/handler/rent"
	trackedlocationhandler "github.com/kotafan1rich/geo_logic_api_go/internal/handler/tracked_location"
	"github.com/kotafan1rich/geo_logic_api_go/internal/handler/user"
	"github.com/kotafan1rich/geo_logic_api_go/internal/logger"
	eventrepo "github.com/kotafan1rich/geo_logic_api_go/internal/repository/event"
	infrarepo "github.com/kotafan1rich/geo_logic_api_go/internal/repository/infra"
	rentrepo "github.com/kotafan1rich/geo_logic_api_go/internal/repository/rent"
	trackedlocationrepo "github.com/kotafan1rich/geo_logic_api_go/internal/repository/tracked_location"
	userrepo "github.com/kotafan1rich/geo_logic_api_go/internal/repository/user"
	eventservice "github.com/kotafan1rich/geo_logic_api_go/internal/service/event"
	infraservice "github.com/kotafan1rich/geo_logic_api_go/internal/service/infra"
	rentservice "github.com/kotafan1rich/geo_logic_api_go/internal/service/rent"
	trackedlocationservice "github.com/kotafan1rich/geo_logic_api_go/internal/service/tracked_location"
	userservice "github.com/kotafan1rich/geo_logic_api_go/internal/service/user"
)

type diContainer struct {
	db database.DB

	userRepo            userrepo.UserRepository
	rentRepo            rentrepo.RentRepository
	eventRepo           eventrepo.EventRepository
	infraTypeRepo       infrarepo.InfraTypeRepository
	infraRepo           infrarepo.InfraRepository
	trackedLocationRepo trackedlocationrepo.TrackedLocationRepository

	userService            userservice.UserService
	rentService            rentservice.RentService
	eventService           eventservice.EventService
	infraTypeService       infraservice.InfraTypeService
	infraService           infraservice.InfraService
	trackedLocationService trackedlocationservice.TrackedLocationService

	handler api.HttpHandler

	log *logger.Logger
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

func (d *diContainer) UserRepo() userrepo.UserRepository {
	if d.userRepo == nil {
		d.userRepo = userrepo.NewRepository(d.DB())
	}

	return d.userRepo
}

func (d *diContainer) RentRepo() rentrepo.RentRepository {
	if d.rentRepo == nil {
		d.rentRepo = rentrepo.NewRepository(d.DB())
	}

	return d.rentRepo
}

func (d *diContainer) EventRepo() eventrepo.EventRepository {
	if d.eventRepo == nil {
		d.eventRepo = eventrepo.NewRepository(d.DB())
	}

	return d.eventRepo
}

func (d *diContainer) InfraTypeRepo() infrarepo.InfraTypeRepository {
	if d.infraTypeRepo == nil {
		d.infraTypeRepo = infrarepo.NewInfraTypeRepository(d.DB())
	}

	return d.infraTypeRepo
}

func (d *diContainer) InfraRepo() infrarepo.InfraRepository {
	if d.infraRepo == nil {
		d.infraRepo = infrarepo.NewInfraRepository(d.DB())
	}

	return d.infraRepo
}

func (d *diContainer) TrackedLocationRepo() trackedlocationrepo.TrackedLocationRepository {
	if d.trackedLocationRepo == nil {
		d.trackedLocationRepo = trackedlocationrepo.NewRepository(d.DB())
	}

	return d.trackedLocationRepo
}

func (d *diContainer) UserService() userservice.UserService {
	if d.userService == nil {
		d.userService = userservice.NewUserService(*d.Logger(), d.UserRepo())
	}

	return d.userService
}

func (d *diContainer) RentService() rentservice.RentService {
	if d.rentService == nil {
		d.rentService = rentservice.NewRentService(*d.Logger(), d.RentRepo())
	}
	return d.rentService
}

func (d *diContainer) EventService() eventservice.EventService {
	if d.eventService == nil {
		d.eventService = eventservice.NewEventService(*d.Logger(), d.EventRepo())
	}
	return d.eventService
}

func (d *diContainer) InfraTypeService() infraservice.InfraTypeService {
	if d.infraTypeService == nil {
		d.infraTypeService = infraservice.NewTypeService(*d.Logger(), d.InfraTypeRepo())
	}
	return d.infraTypeService
}

func (d *diContainer) InfraService() infraservice.InfraService {
	if d.infraService == nil {
		d.infraService = infraservice.NewInfraService(*d.Logger(), d.InfraRepo())
	}
	return d.infraService
}

func (d *diContainer) TrackedLocationService() trackedlocationservice.TrackedLocationService {
	if d.trackedLocationService == nil {
		d.trackedLocationService = trackedlocationservice.NewTrackedLocationService(*d.Logger(), d.TrackedLocationRepo())
	}

	return d.trackedLocationService
}

func (d *diContainer) Handler() api.HttpHandler {
	if d.handler == nil {
		userHandler := user.NewHandler(d.UserService())
		rentHandler := rent.NewHandler(d.RentService())
		eventHandler := event.NewHandler(d.EventService())
		infraTypeHandler := infra.NewTypeHandler(d.InfraTypeService())
		infraHandler := infra.NewInfraHandler(d.InfraService())
		trackedLocationHandler := trackedlocationhandler.NewHandler(d.TrackedLocationService())
		d.handler = api.NewHttpHandler(*d.Logger(), userHandler, rentHandler, eventHandler, infraTypeHandler, infraHandler, trackedLocationHandler)
	}

	return d.handler
}

func (d *diContainer) Logger() *logger.Logger {
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
