package infra

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	conf "github.com/mamalmaleki/go-movie/internal/config"
	"log"
	"os"
)

type config struct {
	ServiceDiscoveryUrl string `mapstructure:"SERVICE_DISCOVERY_URL"`
	JaegerUrl           string `mapstructure:"JAEGER_URL"`
}

func newConfig() (*config, error) {
	log.Println("beginning the creation of the base config object.")
	viper, err := conf.SetViperConfig(os.Getenv(conf.InfraConfigFile))
	if err != nil {
		return nil, err
	}

	var cf config
	if err := viper.Unmarshal(&cf); err != nil {
		return nil, errors.New(fmt.Sprintf("unable to unmarshall the config %v", err))
	}
	log.Println(cf)
	log.Println(os.Getenv(conf.VarHttpServerPort))
	log.Println(viper.GetInt(conf.VarHttpServerPort))
	validate := validator.New()
	if err := validate.Struct(&cf); err != nil {
		return nil, errors.New(fmt.Sprintf("Missing required attributes %v\n", err))
	}

	return &cf, nil
}
