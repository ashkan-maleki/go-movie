package config

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"os"
)


func SetupViper(logger *zap.Logger) error {
	err := readAllConfFiles(logger, EnvVarFiles)
	if err != nil && errors.Is(err, filePathIsEmpty){
		viper.AutomaticEnv()
		logger.Info("Read env variables")
	} else if err != nil {
		return err
	}
	err = validateAllEnvVar(logger)
	if err != nil {
		return err
	}
	return nil
}

func readAllConfFiles(logger *zap.Logger, filenames []string) error{
	for i, filename := range filenames {
		err := readInConf(logger, i, os.Getenv(filename))
		if err != nil {
			return err
		}
	}
	return nil
}

func readInConf(logger *zap.Logger, i int ,filename string) error {
	if filename == "" {
		return filePathIsEmpty
	}
	if _, err := os.Stat(filename); err != nil {
		logger.Info("File does not exists")
		return err
	}
	logger.Info("Filename exists", zap.String("filename:", filename))
	viper.SetConfigType("yaml")
	viper.SetConfigFile(filename)

	if i == 0 {
		if err := viper.ReadInConfig(); err != nil {
			return err
		}
	} else {
		if err := viper.MergeInConfig(); err != nil {
			return err
		}
	}
	return nil
}

func validateAllEnvVar(logger *zap.Logger) error {
	for _, envVar := range EnvVars {
		envValue := viper.Get(envVar)
		switch v := envValue.(type) {
		case int:
			if v == 0 {
				return fmt.Errorf("%s is not provided", envVar)
			}
		case string:
			if v == "" {
				return fmt.Errorf("%s is not provided", envVar)
			}
		default:
			return fmt.Errorf("%s is not provided", envVar)
		}
	}
	logger.Info("All env var are validated")
	return nil
}


