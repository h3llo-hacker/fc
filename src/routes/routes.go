package routes

import (
	"config"
	"handler/docker"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

func Router(router *gin.Engine) {

	router.GET("/", func(c *gin.Context) {
		auth, err := c.Cookie("Auth")
		if err != nil || !validateAuth(auth) {
			c.JSON(401, gin.H{
				"err": "Auth Failed.",
			})
		} else {
			c.JSON(200, gin.H{
				"apis": router.Routes(),
			})
		}
	})

	// Ping
	router.GET("/ping", Ping)

	// User
	router.GET("/users", users)
	userGroup := router.Group("/user")
	{
		userGroup.POST("/login", userLogin)
		userGroup.POST("/create", userCreate)
		userGroup.DELETE("/delete", userDelete)
		userGroup.GET("/:userURL", userInfo)
		userGroup.GET("/:userURL/info", userInfo)
		userGroup.POST("/follow/:userURL", userFollow)
		userGroup.POST("/update/:userURL", userUpdate)
		userGroup.GET("/:userURL/challenges", userChallenges)
		// challenges[?type=0/1/2]
		userGroup.GET("/:userURL/followers", userFollowers)
		userGroup.GET("/:userURL/followees", userFollowees)
	}

	// Challenges
	router.GET("/challenges", challenges)
	challengeGroup := router.Group("/challenge")
	{
		challengeGroup.GET("/:challengeID", challengeInfo)
		challengeGroup.POST("/create", challengeCreate)
		challengeGroup.DELETE("/remove", challengeRemove)
	}

	// templates
	router.GET("/templates", templates)
	templateGroup := router.Group("/template")
	{
		templateGroup.GET("/:templateID", templateQuery)
		templateGroup.POST("/create", templateCreate)
		templateGroup.DELETE("/remove", templateRemove)
	}

	// list all services
	router.GET("/services", listServices)

	// actions
	serviceGroup := router.Group("/service")
	{
		serviceGroup.GET("/:serviceID", inspectService)
		serviceGroup.GET("/:serviceID/status", getServiceStatus)
	}
}

func listServices(c *gin.Context) {
	endpoint := config.Conf.Endpoint
	services, err := docker.ListServices(endpoint)
	if err == nil {
		c.JSON(200, services)
	} else {
		log.Error(err)
		c.JSON(500, err)
	}
}

func inspectService(c *gin.Context) {
	serviceID := c.Param("serviceID")
	service, err := docker.InspectService(serviceID)
	if err == nil {
		if service.ID != "" {
			c.JSON(200, service)
		} else {
			c.JSON(404, gin.H{
				"error": "not found",
			})
		}
	} else {
		log.Error(err)
	}
}

func getServiceStatus(c *gin.Context) {
	serviceID := c.Param("serviceID")
	service, err := docker.InspectServiceTasks(serviceID)
	if err == nil {
		if service.ID != "" {
			c.JSON(200, service)
		} else {
			c.JSON(404, gin.H{
				"error": "not found",
			})
		}
	} else {
		log.Error(err)
	}
}
