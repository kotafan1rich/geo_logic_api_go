package main

import (
	"log/slog"
	"os"

	"github.com/kotafan1rich/geo_logic_api_go/internal/app"
	"github.com/kotafan1rich/geo_logic_api_go/internal/config"
)

func main() {
	config.MustLoad()

	app := app.New()

	if err := app.Run(); err != nil {
		slog.Error("app error", "err", err)
		os.Exit(1)
	}
}
