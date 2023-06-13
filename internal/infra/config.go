package infra

import (
	"errors"
	conf "github.com/mamalmaleki/go-movie/internal/config"
	"github.com/spf13/viper"
	"log"
	"os"
)

type config struct {
	ServiceDiscoveryUrl string `yaml:"service_discovery_url"`
	JaegerUrl           string `yaml:"jaeger_url"`
}

func newConfig() (*config, error) {
	log.Println("beginning the creation of the base config object.")
	err := conf.SetViperConfig(os.Getenv(conf.InfraConfigFile))
	if err != nil {
		return nil, err
	}

	serviceDiscoveryUrl := viper.GetString("service_discovery_url")
	if serviceDiscoveryUrl == "" {
		return nil, errors.New("service_discovery_url is not provided")
	}

	jaegerUrl := viper.GetString("jaeger_url")
	if jaegerUrl == "" {
		return nil, errors.New("jaeger_url is not provided")
	}

	return &config{
		ServiceDiscoveryUrl: serviceDiscoveryUrl,
		JaegerUrl:           jaegerUrl,
	}, nil
}
