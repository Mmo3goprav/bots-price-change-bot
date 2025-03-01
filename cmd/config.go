package main

import (
	"errors"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Config struct {
	Debug   bool   `mapstructure:"debug"`
	AdminID int    `mapstructure:"admin-id"`
	Token   string `mapstructure:"token"`
}

//nolint:lll,funlen // flags are easier to read in single line, cannot make it shorter
func ReadConfig() *Config {
	err := godotenv.Load()
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		panic(err)
	}

	pflag.Bool("debug", false, "enable debug")
	pflag.Int("admin-id", 0, "admin id")
	pflag.String("token", "", "bot token")

	pflag.Parse()

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	err = viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		panic(err)
	}

	var config Config

	if err := viper.Unmarshal(&config); err != nil {
		panic(err)
	}

	return &config
}
