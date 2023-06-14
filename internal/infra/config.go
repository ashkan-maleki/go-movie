package infra

import (
	"errors"
	"fmt"
	validator "github.com/go-playground/validator/v10"
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
	viper, err := conf.SetViperConfig(os.Getenv(conf.InfraConfigFile))
	if err != nil {
		return nil, err
	}

	var conf config
	if err := viper.Unmarshal(&conf); err != nil {
		return nil, errors.New(fmt.Sprintf("unable to unmarshall the config %v", err))
	}
	validate := validator.New()
	if err := validate.Struct(&conf); err != nil {
		return nil, errors.New(fmt.Sprintf("Missing required attributes %v\n", err))
	}
	log.Println(conf)

	return &conf, nil
}
