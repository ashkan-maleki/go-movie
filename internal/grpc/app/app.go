package app

import "github.com/mamalmaleki/go-movie/internal/infra"

type App struct {
	Infra *infra.Infra
}

func New() (*App, error) {
	appInfra, err := infra.New()
	if err != nil {
		return nil, err
	}
	return &App{
		Infra: appInfra,
	}, nil
}
