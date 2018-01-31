package config

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"

	"github.com/BurntSushi/toml"
)

// Config binds together configuration items
type Config struct {
	ConfigFile      string
	ListenIP        string `toml:"listenIP"`
	ListenPort      int    `toml:"listenPort"`
	UseProxyHeaders bool   `toml:"useProxyHeaders"`
	DatabaseHost    string `toml:"databaseHost"`
	DatabasePort    int    `toml:"databasePort"`
	DatabaseName    string `toml:"databaseName"`
}

// ReadConfig reads the configuration from a file and sets
// defaults if needed
func ReadConfig(configFile string) (*Config, error) {
	var (
		config *Config
		err    error
	)

	configFiles := []string{
		configFile,
	}
	if configFile == "" {
		usr, err := user.Current()
		if err != nil {
			return nil, err
		}

		configFiles = []string{
			"./config.toml",
			fmt.Sprintf("%s/.chmgt/config.toml", usr.HomeDir),
			"/etc/chmgt/config.toml",
		}
	}

	for _, configFile := range configFiles {
		_, err = toml.DecodeFile(configFile, &config)
		if err != nil {
			if os.IsNotExist(err) {
				// this error is not enough to exit the function
				log.Print(err)
			} else {
				// not sure what happened, but it isn't good
				return nil, err
			}
		} else {
			// Setup any default values here
			if config.ListenPort == 0 {
				config.ListenPort = 8080
			}
			if config.DatabaseHost == "" {
				config.DatabaseHost = "localhost"
			}
			if config.DatabasePort == 0 {
				config.DatabasePort = 27017
			}
			if config.DatabaseName == "" {
				config.DatabaseName = "chmgt"
			}
			return config, nil
		}
	}

	return nil, errors.New("config: no files found")
}

// NewConfig produces a Config item from flags/configs
func NewConfig() (*Config, error) {
	var configFile string
	flag.StringVar(&configFile, "config", "", "Use config file at path specified")
	flag.Parse()

	return ReadConfig(configFile)
}

// ListenAddr produces a string containing the IP and port
func (c *Config) ListenAddr() string {
	return fmt.Sprintf("%s:%d", c.ListenIP, c.ListenPort)
}

// DatabaseConnection produces a string containing the
// database host and port
func (c *Config) DatabaseConnection() string {
	return fmt.Sprintf("%s:%d", c.DatabaseHost, c.DatabasePort)
}
