package main

import (
	"config"
	"handler/monitor"
	"handler/rank"
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
	go updateUsersRank()

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

func updateUsersRank() {
	log.Infoln("FC UpdateUsersRank Started.")
	for {
		time.Sleep(10 * time.Minute)
		rank.UpdateUsersRank()
	}
}
