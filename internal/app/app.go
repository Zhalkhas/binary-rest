package app

import (
	"github.com/go-chi/chi/v5/middleware"
	httpapi "github.com/zhalkhas/binary-rest/internal/api/http"
	"github.com/zhalkhas/binary-rest/internal/repos"
	"log/slog"
	"net/http"
	"os"
)

type App struct {
	config            Config
	indicesController httpapi.IndicesController
}

func New(config Config) (App, error) {
	f, err := os.Open(config.InputFileName)
	if err != nil {
		slog.Error("error while opening indices input file", "err", err, "file_name", config.InputFileName)
		return App{}, err
	}
	// we are not writing to file,
	// and may assume that error is insignificant
	defer func() { _ = f.Close() }()
	indicesRepo, err := repos.NewReaderIndices(f)
	if err != nil {
		slog.Error("error while creating indices repository", "err", err)
		return App{}, err
	}
	indicesController := httpapi.NewIndicesController(
		indicesRepo,
		middleware.RequestLogger(httpapi.SlogFormatter{}),
		middleware.Recoverer,
	)
	return App{
		config:            config,
		indicesController: indicesController,
	}, nil
}

func (a App) RunHTTP() error {
	slog.Info("starting http server", "port", a.config.Port)
	return http.ListenAndServe(":"+a.config.Port, a.indicesController)
}
