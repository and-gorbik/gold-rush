package config

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v3"
)

type Config struct {
	App           App    `yaml:"app"`
	Explorer      Entity `yaml:"explorer"`
	Earner        Entity `yaml:"earner"`
	Exchanger     Entity `yaml:"exchanger"`
	Licenser      Entity `yaml:"licenser"`
	BalanceClient Client `yaml:"balance_client"`
	StatusClient  Client `yaml:"status_client"`
}

type App struct {
	Port     int `yaml:"port"`
	AreaSize int `yaml:"area_size"`
}

type Entity struct {
	Workers int    `yaml:"workers"`
	Client  Client `yaml:"client"`
}

type Client struct {
	MaxIdleConns        int `yaml:"max_idle_conns"`
	MaxConnsPerHost     int `yaml:"max_conns_per_host"`
	MaxIdleConnsPerHost int `yaml:"max_idle_conns_per_host"`
	Timeout             int `yaml:"timeout"`
}

func LoadFrom(path string) *Config {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	c := new(Config)
	if err = yaml.Unmarshal(data, c); err != nil {
		log.Fatal(err)
	}

	return c
}
