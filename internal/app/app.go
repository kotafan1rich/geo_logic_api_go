package app

import "log/slog"

type App struct {
	diContainer *diContainer
}

func New() *App {
	app := &App{diContainer: newDIContainer()}
	return  app
}

func (a *App) Run() error {
	slog.Info("server is running")

	_ = a.diContainer.DB()

	return nil
}
