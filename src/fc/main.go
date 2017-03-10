package main

import (
	"config"
	"encoding/json"
	"net/http"
	"routes"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

func main() {
	log.Info("FC started.")
	log.SetLevel(log.DebugLevel)

	// load config
	conf, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	ccc, _ := json.Marshal(conf)
	log.Debugln(string(ccc))

	router := gin.Default()

	routes.Router(router)

	server := &http.Server{
		Addr:           ":8080",
		Handler:        router,
		ReadTimeout:    3 * time.Second,
		WriteTimeout:   3 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Errorf("FC Error: %v", err)
	}
	log.Info("FC Exit")
}
