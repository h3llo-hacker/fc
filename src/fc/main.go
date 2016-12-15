package main

import (
	// "encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	// "github.com/docker/docker/api/types/swarm"
	"github.com/gin-gonic/gin"
	"handler"
	"net/http"
	"time"
)

func main() {
	log.Info("Main process started.")
	log.SetLevel(log.DebugLevel)
	gin.SetMode(gin.ReleaseMode)

	// // just test
	// endpoint := "127.0.0.1:2374"
	// services, err := handler.ListServices(endpoint)
	// if err == nil {
	// 	for _, service := range services {
	// 		log.Info(service.Spec.Name)
	// 	}
	// } else {
	// 	log.Error(err)
	// }

	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"ping": "pong",
		})
	})

	router.GET("/services/:endpoint", listService)

	server := &http.Server{
		Addr:           ":8001",
		Handler:        router,
		ReadTimeout:    3600 * time.Second,
		WriteTimeout:   3600 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	server.ListenAndServe()
	fmt.Println("Hello world!")
}

func listService(c *gin.Context) {
	endpoint := c.Param("endpoint")
	services, err := handler.ListServices(endpoint)
	if err == nil {
		c.JSON(200, services)
	} else {
		log.Error(err)
	}
}

func createService(c *gin.Context) error {
	return nil
}
