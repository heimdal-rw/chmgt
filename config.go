package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/user"

	"github.com/BurntSushi/toml"
)

// Config is a struct of expected configuration elements
type Config struct {
	ServerListen string
	Database     string
}

// ReadConfig takes the configfile string and attempts to open it and parse toml
// if the commandline flag for config file is not specified, it tries a few other locations
func ReadConfig(configFileFlag string) (*Config, error) {
	// define paths of where we might find the config file
	configFiles := []string{
		configFileFlag,
	}

	if configFileFlag == "" {
		// only use predefined values if the config flag was not used
		usr, err := user.Current()
		if err != nil {
			return nil, err
		}
		// replace the configFiles slice with preconfigured paths
		configFiles = []string{
			"./config",
			fmt.Sprintf("%s/.chmgt/config", usr.HomeDir),
			"/etc/chmgt/config",
		}
	}

	var config *Config
	for _, configFile := range configFiles {
		log.Printf("Attempting to use config file: %s", configFile)
		if _, err := toml.DecodeFile(configFile, &config); err != nil {
			if os.IsNotExist(err) {
				// this error is not enough to exit the function
				log.Print("Config file not found.")
			} else {
				// not sure what happened, but it isn't good
				return nil, err
			}
		} else {
			// found and parsed a valid config, return it
			return config, nil
		}
	}

	return nil, errors.New("no config files found")
}
