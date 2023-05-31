package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	APIURL     string `mapstructure:"api_url"`
	APIPort    string `mapstructure:"api_port"`
	DBType     string `mapstructure:"db_type"`
	DBHost     string `mapstructure:"db_host"`
	DBPort     int    `mapstructure:"db_port"`
	DBUser     string `mapstructure:"db_user"`
	DBPassword string `mapstructure:"db_password"`
	DBName     string `mapstructure:"db_name"`
	Env        string `mapstructure:"env"`
}

func Get() *Config {
	viper.SetConfigFile("app/config/local.yml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}

	conf := &Config{}
	if err := viper.Unmarshal(conf); err != nil {
		log.Fatal(err)
	}

	return conf
}
