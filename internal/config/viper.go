package config

import (
	"github.com/spf13/viper"
	"log"
	"os"
)

func SetViperConfig(filename string) error {
	if filename == "" {
		log.Println("filename is empty")
		viper.AutomaticEnv()
	} else {
		if _, err := os.Stat(filename); err != nil {
			log.Printf("File does not exists\n")
			return err
		}
		log.Println("filename exists")
		viper.SetConfigType("yaml")
		viper.SetConfigFile(filename)

	}

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	return nil
}
