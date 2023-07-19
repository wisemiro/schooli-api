package config

import (
	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
)

type Config struct {
	Database struct {
		Server      string
		Port        uint16
		Database    string
		User        string
		Password    string
		Secure      bool
		Timeout     int
		Automigrate bool
		Log         bool
	}
	Minio struct {
		Address    string
		BucketName string
		Region     string
		User       string
		Password   string
		Timeout    int
	}
	Secrets struct {
		TTL          int
		Jwt          string
		FireBaseFile string
	}
	Redis struct {
		Address string
		Port    string
		DB      int
		Timeout int
	}
	Server struct {
		Address string
		Timeout int
		Prod    bool
	}
}

func Load(loc string) (Config, error) {
	var config Config
	if _, err := toml.DecodeFile(loc, &config); err != nil {
		return config, errors.Wrap(err, "Unable to decode config")
	}
	return config, nil
}
