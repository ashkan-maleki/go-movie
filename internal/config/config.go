package config

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path/filepath"
)

type BaseConfig struct {
	Base base `yaml:"base"`
}

type base struct {
	ServiceDiscoveryUrl string `yaml:"service_discovery_url"`
	HttpServerPort      int    `yaml:"httpServerPort"`
}

func loadYamlConfig2(filename string) (*BaseConfig, error) {
	if _, err := os.Stat(filename); err != nil {
		fmt.Printf("File does not exists\n")
		return nil, err
	}
	log.Println(filepath.Abs(filename))
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	//var data []byte = make([]byte, 70)
	//_, err = f.Read(data)
	//if err != nil {
	//	return nil, err
	//}
	//myString := string(data[:])
	//log.Println(myString)
	//log.Println(len(data))
	log.Println("the base config file is read.")
	var cfg BaseConfig
	//err = yaml.Unmarshal(data, cfg)
	if err := yaml.NewDecoder(f).Decode(cfg); err != nil {
		//if err != nil {
		log.Println("yaml decoding: ")
		return nil, err
	}
	return &cfg, nil
}

func loadYamlConfig(filename string) (*BaseConfig, error) {
	if _, err := os.Stat(filename); err != nil {
		fmt.Printf("File does not exists\n")
		return nil, err
	}
	viper.SetConfigType("yaml")
	viper.SetConfigFile(filename)

	log.Println(viper.Get("name"))
	return &BaseConfig{}, nil
}

func loadEnvironmentVariables() (*BaseConfig, error) {
	serviceDiscoveryUrl := os.Getenv("SERVICE_DISCOVERY_URL")
	if serviceDiscoveryUrl == "" {
		return nil, errors.New("SERVICE_DISCOVERY_URL is empty")
	}
	baseC := base{
		ServiceDiscoveryUrl: serviceDiscoveryUrl,
	}
	return &BaseConfig{
		Base: baseC,
	}, nil
}

func New() (*BaseConfig, error) {
	log.Println("beginning the creation of the base config object.")
	filename := os.Getenv(BaseConfigFile)
	if filename == "" {
		return loadEnvironmentVariables()
	} else {
		return loadYamlConfig(filename)
	}
}
