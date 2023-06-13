package app

import (
	"github.com/mamalmaleki/go-movie/internal/grpc/app"
)

type App struct {
	//Infra  *infra.Infra
	*app.App
	Config *config
}

func New() (*App, error) {
	app, err := app.New()
	if err != nil {
		return nil, err
	}
	conf, err := newConfig()
	if err != nil {
		return nil, err
	}
	return &App{
		App:    app,
		Config: conf,
	}, nil
}
