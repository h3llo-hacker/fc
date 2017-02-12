package routes

import (
	"config"
	"handler/docker"
	db "utils/db"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
)

func Router(router *gin.Engine) {

	router.GET("/", func(c *gin.Context) {
		auth, err := c.Cookie("Auth")
		if err != nil || !valiedAuth(auth) {
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

	// Users
	router.GET("/users", users)

	// User
	userGroup := router.Group("/user")
	{
		userGroup.POST("/create", userCreate)
		userGroup.DELETE("/delete", userDelete)
		userGroup.POST("/update/:userURL", userUpdate)
		userGroup.POST("/follow/:userURL", userFollow)
		userGroup.GET("/:userURL", userInfo)
		userGroup.GET("/:userURL/", userInfo)
		userGroup.GET("/:userURL/info", userInfo)
		userGroup.GET("/:userURL/challenges", userChallenges)
		// challenges[?type=0/1/2]
		userGroup.GET("/:userURL/followers", userFollowers)
		userGroup.GET("/:userURL/followees", userFollowees)
	}

	// Challenges
	router.GET("/challenges", challenges)

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
	for _, endpoint := range config.Conf.Endpoints {
		services, err := docker.ListServices(endpoint)
		if err == nil {
			c.JSON(200, services)
		} else {
			log.Error(err)
		}
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

func valiedAuth(auth string) bool {
	C := "auth"
	selector := bson.M{}
	auths, err := db.MongoFind(C, nil, selector)
	if err != nil {
		return false
	}
	for _, Auth := range auths {
		if auth == Auth.(bson.M)["token"] {
			return true
		}
	}
	return false
}
