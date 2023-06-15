package infra

import (
	conf "github.com/mamalmaleki/go-movie/internal/config"
	"log"
	"os"
)

type config struct {
	ServiceDiscoveryUrl string `mapstructure:"SERVICE_DISCOVERY_URL" validate:"required"`
	JaegerUrl           string `mapstructure:"JAEGER_URL" validate:"required"`
}

func newConfig() (*config, error) {
	log.Println("beginning the creation of the base config object.")
	var cf config
	err := conf.Configure(os.Getenv(conf.InfraConfigFile), &cf)
	if err != nil {
		return nil, err
	}
	log.Println("ending the creation of the base config object.")
	return &cf, nil
}
