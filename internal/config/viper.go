package config

import (
	"github.com/spf13/viper"
	"log"
	"os"
)

func SetViperConfig(filename string) (*viper.Viper, error) {
	vip := viper.New()
	if filename == "" {
		log.Println("filename is empty")

		vip.AutomaticEnv()
	} else {
		if _, err := os.Stat(filename); err != nil {
			log.Printf("File does not exists\n")
			return nil, err
		}
		log.Println("filename exists")
		vip.SetConfigType("yaml")
		vip.SetConfigFile(filename)
		err := vip.ReadInConfig()
		if err != nil {
			log.Printf("Reading config failed\n")
			return nil, err
		}

	}

	return vip, nil
}
