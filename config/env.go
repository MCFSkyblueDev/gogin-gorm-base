package config

import (
	"log"

	"github.com/spf13/viper"
)

// Env struct holds all the configuration settings
type Env struct {
	Port      string `mapstructure:"PORT"`
	DBUrl     string `mapstructure:"DATABASE_URL"`
	JwtSecret string `mapstructure:"JWT_SECRET"`
	JwtExpire int    `mapstructure:"JWT_EXPIRE"`
}

func NewEnv(filename string, override bool) *Env {
	env := Env{}
	viper.SetConfigFile(filename)

	if override {
		viper.AutomaticEnv()
	}

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Error reading environment file", err)
	}

	err = viper.Unmarshal(&env)
	if err != nil {
		log.Fatal("Error loading environment file", err)
	}
	return &env
}
