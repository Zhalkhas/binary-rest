package app

import (
	"github.com/go-chi/chi/v5/middleware"
	httpapi "github.com/zhalkhas/binary-rest/internal/api/http"
	"github.com/zhalkhas/binary-rest/internal/repos"
	"log/slog"
	"net/http"
	"os"
)

// App is core of the application,
// which is responsible for collecting, orchestrating and launching
// all the components of the application.
type App struct {
	config Config
	// rootHandler is HTTP handler that is responsible for handling all the requests.
	rootHandler http.Handler
}

func New(config Config) (App, error) {
	// initialize repositories
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

	// initialize controllers
	indicesController := httpapi.NewIndicesController(
		indicesRepo,
		middleware.RequestLogger(httpapi.SlogFormatter{}),
		middleware.Recoverer,
	)

	return App{
		config: config,
		// since we have only one handler, we can use it as root handler,
		// if we have more handlers, we can mount subrouters to the root handler
		rootHandler: indicesController,
	}, nil
}

// RunHTTP starts http server at port specified in [Config.Port] during initialization.
func (a App) RunHTTP() error {
	slog.Info("starting http server", "port", a.config.Port)
	return http.ListenAndServe(":"+a.config.Port, a.rootHandler)
}
