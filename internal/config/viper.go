package config

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
	"reflect"
)

// https://benchkram.de/blog/dev/ultimate-config-for-golang-apps

func BindEnvAll(vip *viper.Viper, config any) error {
	t := reflect.TypeOf(config)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("mapstructure")
		err := vip.BindEnv(tag)
		if err != nil {
			return err
		}
	}
	return nil
}

func Configure(filename string, config any) error {
	vip := viper.New()
	if filename == "" {
		log.Println("filename is empty")
		vip.AutomaticEnv()
		if err := BindEnvAll(vip, config); err != nil {
			return err
		}
	} else {
		if _, err := os.Stat(filename); err != nil {
			log.Printf("File does not exists\n")
			return err
		}
		log.Println("filename exists")
		vip.SetConfigType("yaml")
		vip.SetConfigFile(filename)

		if err := vip.ReadInConfig(); err != nil {
			log.Printf("Reading config failed\n")
			return err
		}
	}
	if err := vip.Unmarshal(config); err != nil {
		return err
	}
	log.Println(config)
	if err := IsValid(config); err != nil {
		return errors.New(fmt.Sprintf("Missing required attributes %v\n", err))
	}
	return nil
}
