package app

import (
	conf "github.com/mamalmaleki/go-movie/internal/config"
	"log"
)

type config struct {
	*conf.CommonConfig
}

func newConfig() (*config, error) {
	log.Println("beginning the creation of the base config object.")

	commonConf, err := conf.NewCommonConfig()
	if err != nil {
		return nil, err
	}

	return &config{
		CommonConfig: commonConf,
	}, nil
}
