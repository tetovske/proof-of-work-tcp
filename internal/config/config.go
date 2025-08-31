package config

import (
	"os"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   Server
	Hashcash Hashcash
	Quotes   Quotes `mapstructure:"quotes"`
}

type Hashcash struct {
	Complexity int
	Secret     []byte
	TTL        time.Duration
}

type Quotes struct {
	Data []string `mapstructure:"data"`
}

type Server struct {
	Port string
}

func New() *Config {
	c := &Config{}

	return c
}

func (c *Config) ReadFromConfigAndENV(path *string) error {
	if path != nil && *path != "" {
		_, err := os.Lstat(*path)
		if err != nil {
			return err
		}

		viper.SetConfigFile(*path)

		if err = viper.ReadInConfig(); err != nil {
			return err
		}

		if err = viper.Unmarshal(&c); err != nil {
			return err
		}
	}

	viper.AutomaticEnv()

	if val := viper.GetString("SERVER_PORT"); val != "" {
		c.Server.Port = val
	}

	if val := viper.GetString("HASHCASH_SECRET"); val != "" {
		c.Hashcash.Secret = []byte(val)
	}

	if val := viper.GetInt("HASHCASH_COMPLEXITY"); val != 0 {
		c.Hashcash.Complexity = val
	}

	if val := viper.GetDuration("HASHCASH_TTL"); val != 0 {
		c.Hashcash.TTL = val
	}

	return nil
}
