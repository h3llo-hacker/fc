package config

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"os"
	"strings"
)

type MongoDB_Conf struct {
	Host string
	Port string
	User string
	Pass string
	DB   string
}

type Config struct {
	Endpoints []string
	MongoDB   MongoDB_Conf
	LogLevel  string
}

var C Config

func readConfig() *Config {
	pwd, err := os.Getwd()
	filePath := pwd + "/config.json"
	log.Debug("Load config from file \"", filePath, "\"")
	file, err := os.Open(filePath)
	if err != nil {
		log.Error(err)
		panic(err)
	}
	decoder := json.NewDecoder(file)
	configuration := &Config{}
	err = decoder.Decode(configuration)
	if err != nil {
		log.Error(err)
		panic(err)
	}
	// inject configuration
	C = *configuration
	return configuration
}

func setLogLevel(logLevel string) {
	log.Debug("Setting log level ", logLevel)
	switch strings.ToLower(logLevel) {
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "debug":
		log.SetLevel(log.DebugLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}
}

func LoadConfig() *Config {
	conf := readConfig()
	setLogLevel(conf.LogLevel)
	return conf
}
