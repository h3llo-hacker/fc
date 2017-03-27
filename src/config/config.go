package config

import (
	"encoding/json"
	"os"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
)

type MongoDB_Conf struct {
	Host string `json:"Host"`
	Port string `json:"Port"`
	User string `json:"User"`
	Pass string `json:"Pass"`
	DB   string `json:"DB"`
}

type Etcd_struct struct {
	Hosts []string `json:"Hosts"`
	User  string   `json:User`
	Pass  string   `json:Pass`
}

type Mail_Templates struct {
	ValidateEmail string `json:"ValidateEmail"`
	ResetPassword string `json:"ResetPassword"`
}

type Mail_config struct {
	SendGridKey string         `json:"SendGridKey"`
	Sender      string         `json:"Sender"`
	Templates   Mail_Templates `json:"Templates"`
}

type Config struct {
	Endpoint          string        `json:"Endpoint"`
	Etcd              Etcd_struct   `json:"Etcd"`
	Mail              Mail_config   `json:"Mail"`
	LogLevel          string        `json:"LogLevel"`
	InviteMode        bool          `json:"InviteMode"`
	InviteCodes       int           `json:"InviteCodes"`
	ComposeFilePath   string        `json:"ComposeFilePath"`
	MongoDB           MongoDB_Conf  `json:"MongoDB"`
	ChallengeDuration time.Duration `json:"ChallengeDuration"`
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

func setTimeZone() {
	time.FixedZone("+8", 8)
	log.Debugf("Setting TimeZone +8, Now: [%v]", time.Now())
}

func LoadConfig() (*Config, error) {
	if Conf.LogLevel != "" {
		return &Conf, nil
	}
	conf, err := ReadConfig()
	if err != nil {
		return nil, err
	}
	setLogLevel(conf.LogLevel)
	setTimeZone()
	return conf, nil
}
