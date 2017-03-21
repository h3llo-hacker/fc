package main

import (
	"config"
	"handler/monitor"
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
	_, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	go monitorTimeoutChallenges()

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
		log.Errorf("FC Start Error: %v", err)
	}
	log.Info("FC Exit")
}

func monitorTimeoutChallenges() {
	log.Infoln("FC Monitor Started.")
	for {
		time.Sleep(1 * time.Minute)
		monitor.RemoveTimeOutChallenges()
	}
}
