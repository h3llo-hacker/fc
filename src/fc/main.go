package main

import (
	"config"
	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"net/http"
	// "os"
	"routes"
	"time"
)

func main() {
	log.Info("FC started.")
	log.SetLevel(log.DebugLevel)

	// load config
	conf, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	log.Debug("Swarmkit Endpoints: ", conf.Endpoints)

	router := gin.Default()
	routes.Router(router)

	server := &http.Server{
		Addr:           ":8080",
		Handler:        router,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	err = server.ListenAndServe()
	if err != nil {
		log.Errorf("FC Error: %v", err)
	}
	log.Info("FC Exit 0.")
}
