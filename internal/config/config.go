package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

const (
	configPath = "config.yaml"
)

//Config is the app config struct
type Config struct {
	AppName  string `yaml:"app_name"`
	Postgres struct {
		DBName   string `yaml:"db_name"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Port     string `yaml:"port"`
		DSN      string `yaml:"dsn"`
	}
	Kafka struct {
		Brokers []string `yaml:"brokers"`
	} `yaml:"kafka"`
}

//NewConfig loads the config file
// it's just a temporal config handling strategy
func NewConfig() (*Config, error) {
	conf := &Config{}
	file, err := os.Open(configPath)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	if err = yaml.NewDecoder(file).Decode(conf); err != nil {
		return nil, err
	}

	return conf, nil
}
