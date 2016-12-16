package main

import (
	"config"
	// "fmt"
	log "github.com/Sirupsen/logrus"
	// "github.com/docker/docker/api/types/swarm"
	"github.com/gin-gonic/gin"
	"handler"
	"net/http"
	"os"
	"time"
)

func main() {
	log.Info("Main process started.")
	log.SetLevel(log.DebugLevel)
	if os.Getenv("release") != "" {
		gin.SetMode(gin.ReleaseMode)
	}

	// load config
	conf, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	log.Debug("Swarmkit Endpoints: ", conf.Endpoints)

	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"apis": router.Routes(),
		})
	})
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"ping": "pong",
		})
	})
	router.GET("/pong", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"pong": "ping",
		})
	})

	// list all services
	router.GET("/services", listServices)

	// service actions
	serviceGroup := router.Group("/service")
	{
		serviceGroup.GET("/:serviceID", inspectAservice)
		serviceGroup.GET("/:serviceID/status", getServiceStatus)
	}

	server := &http.Server{
		Addr:           ":8083",
		Handler:        router,
		ReadTimeout:    3600 * time.Second,
		WriteTimeout:   3600 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	server.ListenAndServe()
	log.Info("Exit 0")
}

func listServices(c *gin.Context) {
	for _, endpoint := range config.Conf.Endpoints {
		services, err := handler.ListServices(endpoint)
		if err == nil {
			c.JSON(200, services)
		} else {
			log.Error(err)
		}
	}
}

func inspectAservice(c *gin.Context) {
	serviceID := c.Param("serviceID")
	service, err := handler.InspectService(serviceID)
	if err == nil {
		if service.ID != "" {
			c.JSON(200, service)
		} else {
			c.JSON(404, gin.H{
				"error": "service not found",
			})
		}
	} else {
		log.Error(err)
	}
}

func getServiceStatus(c *gin.Context) {
	serviceID := c.Param("serviceID")
	service, err := handler.InspectServiceTasks(serviceID)
	if err == nil {
		if service.ID != "" {
			c.JSON(200, service)
		} else {
			c.JSON(404, gin.H{
				"error": "service not found",
			})
		}
	} else {
		log.Error(err)
	}
}
