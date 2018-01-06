package main

import (
	"flag"
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

// Config is a struct of expected configuration elements
type Config struct {
	Interface string
	Port      string
}

// ReadConfig can accept a runtime flag of --config with a filename or default to ./config.
// This should probably be changed to somewhere in /etc once we package up chmgt for use.
func ReadConfig() Config {
	var configfile string
	flag.StringVar(&configfile, "config", "./config", "Config file to be used")
	flag.Parse()
	log.Printf("Attempting to use config file: %s", configfile)

	_, err := os.Stat(configfile)
	if err != nil {
		log.Fatal("Config file is missing: ", configfile)
	}

	var config Config
	if _, err := toml.DecodeFile(configfile, &config); err != nil {
		log.Fatal(err)
	}

	return config
}
