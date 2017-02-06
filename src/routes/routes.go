package routes

import (
	"config"
	// network "handler/network"
	service "handler/service"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

func Router(router *gin.Engine) {

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"apis": router.Routes(),
		})
	})

	// Ping
	router.GET("/ping", Ping)

	// Users
	router.GET("/users", users)
	userGroup := router.Group("/user")
	{
		userGroup.POST("/create", userCreate)
		userGroup.POST("/update/:userURL", userUpdate)
		userGroup.POST("/follow/:userURL", userFollow)
		userGroup.GET("/:userURL/info", userInfo)
		userGroup.GET("/:userURL/challenges", userChallenges)
		userGroup.GET("/:userURL/followers", userFollowers)
		userGroup.GET("/:userURL/followees", userFollowees)
	}

	// Challenges
	router.GET("/challenges", challenges)

	// list all services
	router.GET("/services", listServices)

	// service actions
	serviceGroup := router.Group("/service")
	{
		serviceGroup.GET("/:serviceID", inspectService)
		serviceGroup.GET("/:serviceID/status", getServiceStatus)
	}
}

func listServices(c *gin.Context) {
	for _, endpoint := range config.Conf.Endpoints {
		services, err := service.ListServices(endpoint)
		if err == nil {
			c.JSON(200, services)
		} else {
			log.Error(err)
		}
	}
}

func inspectService(c *gin.Context) {
	serviceID := c.Param("serviceID")
	service, err := service.InspectService(serviceID)
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
	service, err := service.InspectServiceTasks(serviceID)
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
