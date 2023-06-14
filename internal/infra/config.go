package infra

import (
	"errors"
	"fmt"
	conf "github.com/mamalmaleki/go-movie/internal/config"
	"log"
	"os"
)

type config struct {
	ServiceDiscoveryUrl string `yaml:"SERVICE_DISCOVERY_URL"`
	JaegerUrl           string `yaml:"JAEGER_URL"`
}

func newConfig() (*config, error) {
	log.Println("beginning the creation of the base config object.")
	viper, err := conf.SetViperConfig(os.Getenv(conf.InfraConfigFile))
	if err != nil {
		return nil, err
	}

	serviceDiscoveryUrl := viper.GetString(conf.VarServiceDiscoveryUrl)
	if serviceDiscoveryUrl == "" {
		return nil, errors.New(fmt.Sprintf("%s is not provided", conf.VarServiceDiscoveryUrl))
	}

	jaegerUrl := viper.GetString(conf.VarJaegerUrl)
	if jaegerUrl == "" {
		return nil, errors.New(fmt.Sprintf("%s is not provided", conf.VarJaegerUrl))
	}

	return &config{
		ServiceDiscoveryUrl: serviceDiscoveryUrl,
		JaegerUrl:           jaegerUrl,
	}, nil
}
