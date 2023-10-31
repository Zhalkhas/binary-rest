package app

import "net/http"

type App struct {
	config Config
}

func New(config Config) App {
	return App{config: config}
}

func (a App) Run() error {
	return http.ListenAndServe(":"+a.config.Port, nil)
}
