package main

import (
	"github.com/zhalkhas/binary-rest/internal/app"
	"log/slog"
)

func main() {
	config, err := app.NewConfig()
	if err != nil {
		// using logger before init???
		slog.Error("error initializing config", "err", err)
		return
	}
	a, err := app.New(config)
	if err != nil {
		slog.Error("error initializing app", "err", err, "config", config)
		return
	}
	if err := a.RunHTTP(); err != nil {
		slog.Error("error during app running", "err", err)
		return
	}
}
