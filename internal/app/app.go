package app

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/kotafan1rich/geo_logic_api_go/internal/config"
	"github.com/kotafan1rich/geo_logic_api_go/internal/model"
)

type App struct {
	diContainer *diContainer
	httpServer  *http.Server
}

func New() *App {
	app := &App{diContainer: newDIContainer()}

	app.initDeps()

	return app
}

func (a *App) initHTTPServer() {
	a.httpServer = &http.Server{
		Addr:         fmt.Sprintf(":%s", config.Get().HttpServer.ServerPort),
		Handler:      a.diContainer.Handler().Routes(),
		ReadTimeout:  config.Get().HttpServer.ReadTimeout,
		WriteTimeout: config.Get().HttpServer.WriteTimeout,
		IdleTimeout:  config.Get().HttpServer.IdleTimeout,
	}
}

func (a *App) migrateDB() {
	if err := a.diContainer.DB().GORM().AutoMigrate(&model.User{}); err != nil {
		slog.Error("Failed to run migrations")
		os.Exit(1)
	}
}

func (a *App) initDeps() {
	inits := []func(){
		a.initHTTPServer,
		a.migrateDB,
	}

	for _, fn := range inits {
		fn()
	}
}

func (a *App) Run() error {
	slog.Info("server is running")

	return a.httpServer.ListenAndServe()
}
