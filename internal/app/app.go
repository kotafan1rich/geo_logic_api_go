package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kotafan1rich/geo_logic_api_go/internal/config"
	eventmodel "github.com/kotafan1rich/geo_logic_api_go/internal/repository/event/dbmodel"
	rentmodel "github.com/kotafan1rich/geo_logic_api_go/internal/repository/rent/dbmodel"
	usermodel "github.com/kotafan1rich/geo_logic_api_go/internal/repository/user/dbmodel"
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
	gin.SetMode(config.Get().HttpServer.GinMode)
	a.httpServer = &http.Server{
		Addr:         fmt.Sprintf(":%s", config.Get().HttpServer.ServerPort),
		Handler:      a.diContainer.Handler().Routes(),
		ReadTimeout:  config.Get().HttpServer.ReadTimeout,
		WriteTimeout: config.Get().HttpServer.WriteTimeout,
		IdleTimeout:  config.Get().HttpServer.IdleTimeout,
	}
}

func (a *App) migrateDB() {
	if err := a.diContainer.DB().GORM().AutoMigrate(&usermodel.User{}, &rentmodel.Rent{}, &eventmodel.Event{}); err != nil {
		slog.Error("Failed to run migrations")
		os.Exit(1)
	}
}

func (a *App) closeDB() error {
	sqlDB, err := a.diContainer.DB().GORM().DB()
	if err != nil {
		slog.Error("failed to close database connection cleanly", "err", err)
		return err
	}

	err = sqlDB.Close()
	if err != nil {
		slog.Error("failed to close database connection cleanly", "err", err)
		return err
	}
	slog.Info("db is closed")
	return nil
}

func (a *App) gracefullShutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := a.httpServer.Shutdown(ctx); err != nil {
		slog.Error("HTTP server graceful shutdown failed", "err", err)
		return fmt.Errorf("server shutdown failed: %w", err)
	}

	err := a.closeDB()
	if err != nil {
		slog.Error("failed to close sql.DB", "err", err)
		return err
	}
	return nil
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

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		slog.Info(fmt.Sprintf("http server is running on %s", a.httpServer.Addr))
		err := a.httpServer.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("HTTP server failed to listen", "err", err)
			os.Exit(1)
		}
	}()

	sig := <-quit
	slog.Info("shutdown signal received, starting graceful shutdown...", "signal", sig.String())
	err := a.gracefullShutdown()
	if err != nil {
		return err
	}
	slog.Info("server stopped cleanly")
	return nil
}
