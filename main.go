package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/h3llo-hacker/fc/config"
	"github.com/h3llo-hacker/fc/handler/monitor"
	"github.com/h3llo-hacker/fc/handler/rank"
	"github.com/h3llo-hacker/fc/routes"
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
		monitor.RemoveTimeOutChallenges()
		time.Sleep(3 * time.Minute)
	}
}

func updateUsersRank() {
	log.Infoln("FC UpdateUsersRank Started.")
	for {
		rank.UpdateUsersRank()
		time.Sleep(10 * time.Minute)
	}
}
