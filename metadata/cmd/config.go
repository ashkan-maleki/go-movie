package main

import (
	"errors"
	"fmt"
	config2 "github.com/mamalmaleki/go-movie/internal/config"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"strconv"
)

type jaegerConfig struct {
	URL string `yaml:"url"`
}

type apiConfig struct {
	Port int `yaml:"port"`
}

type prometheusConfig struct {
	MetricsPort int `yaml:"metricsPort"`
}

type config struct {
	API        apiConfig        `yaml:"api"`
	Jaeger     jaegerConfig     `yaml:"jaeger"`
	Prometheus prometheusConfig `yaml:"prometheus"`
}
type Config struct {
	config2.BaseConfig
}

func newConfig() (*Config, error) {
	log.Println("beginning the creation of the config object.")
	var cfg *Config

	baseCfg, err := config2.New()
	if err != nil {
		log.Println(fmt.Sprintf("base config creation failed: %s", err))
		return nil, err
	}
	cfg.BaseConfig = *baseCfg

	log.Println("the base config object is created.")

	filename := os.Getenv(config2.AppConfigFile)
	if filename == "" {
		log.Println("app config yaml file is not provided, then use environment variables")
		err := cfg.loadEnvironmentVariables()
		if err != nil {
			return nil, err
		}
	} else {
		log.Println("app config yaml file is provided")
		err := cfg.loadYamlConfig(filename)
		if err != nil {
			return nil, err
		}
	}
	log.Println("ending the creation of the config object.")
	return cfg, nil
}

func (cfg *Config) loadYamlConfig(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
		return err
	}
	return nil
}

func (cfg *Config) loadEnvironmentVariables() error {
	httpServerPortStr := os.Getenv("HTTP_SERVER_PORT")
	if httpServerPortStr == "" {
		return errors.New("HTTP_SERVER_PORT is empty")
	}
	httpServerPort, err := strconv.Atoi(httpServerPortStr)
	if err != nil {
		return errors.New(fmt.Sprintf("HTTP_SERVER_PORT is not correct: %s", err))
	}
	cfg.Base.HttpServerPort = httpServerPort
	return nil
}
