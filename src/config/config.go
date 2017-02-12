package config

import (
	"encoding/json"
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
)

type MongoDB_Conf struct {
	Host string
	Port string
	User string
	Pass string
	DB   string
}

type Config struct {
	Endpoints   []string
	MongoDB     MongoDB_Conf
	LogLevel    string
	SendGridKey string
	InviteMode  bool
}

var Conf Config

const Salt = "h311oW0rlD"

func PathExist(_path string) bool {
	_, err := os.Stat(_path)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

func ReadConfig() (*Config, error) {
	pwd, err := os.Getwd()
	filePath := pwd + "/config.json"

	if !PathExist(filePath) {
		filePath = os.Getenv("FC_CONFIG")
	}

	if filePath == "" {
		filePath = "/home/mr/Documents/work_space/fc/bin/config.json"
	}

	if !PathExist(filePath) {
		log.Errorf("Config File [%s] Not Found!", filePath)
		os.Exit(1)
	}

	log.Debug("Load config from file \"", filePath, "\"")

	file, _ := os.Open(filePath)
	decoder := json.NewDecoder(file)
	configuration := &Config{}
	err = decoder.Decode(configuration)
	if err != nil {
		log.Error(err)
		panic(err)
	}
	// inject configuration
	Conf = *configuration
	return configuration, nil
}

func setLogLevel(logLevel string) {
	log.Debugf("Setting log level: %s", logLevel)
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

func LoadConfig() (*Config, error) {
	conf, err := ReadConfig()
	if err != nil {
		return nil, err
	}
	setLogLevel(conf.LogLevel)
	return conf, nil
}
