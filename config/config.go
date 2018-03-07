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

type serverConfig struct {
	ListenIP        string `toml:"listenIP"`
	ListenPort      int    `toml:"listenPort"`
	UseProxyHeaders bool   `toml:"useProxyHeaders"`
	SessionSecret   string `toml:"sessionSecret"`
	SessionTimeout  int    `toml:"sessionTimeout"`
}

type databaseConfig struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	Name     string `toml:"name"`
	Username string `toml:"username"`
	Password string `toml:"password"`
	AuthDB   string `toml:"authDB"`
}

// Config binds together configuration items
type Config struct {
	ConfigFile string
	Server     serverConfig
	Database   databaseConfig
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
			if config.Server.ListenPort == 0 {
				config.Server.ListenPort = 8080
			}
			if config.Database.Host == "" {
				config.Database.Host = "localhost"
			}
			if config.Database.Port == 0 {
				config.Database.Port = 27017
			}
			if config.Database.Name == "" {
				config.Database.Name = "chmgt"
			}
			if config.Database.AuthDB == "" {
				config.Database.AuthDB = "admin"
			}
			if config.Database.Username == "" {
				config.Database.Username = "chmgt"
			}
			if config.Database.Password == "" {
				config.Database.Password = "chmgtPass"
			}
			if config.Server.SessionSecret == "" {
				config.Server.SessionSecret = "not_so_secret"
			}
			if config.Server.SessionTimeout == 0 {
				config.Server.SessionTimeout = 30
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
	return fmt.Sprintf("%s:%d", c.Server.ListenIP, c.Server.ListenPort)
}

// DatabaseConnection produces a string containing the
// database host and port
func (c *Config) DatabaseConnection() string {
	return fmt.Sprintf("%s:%d", c.Database.Host, c.Database.Port)
}
