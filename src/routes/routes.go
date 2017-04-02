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
	// uG = user group
	uG := router.Group("/user")
	{
		uG.POST("/login", userLogin)
		uG.POST("/create", userCreate)
		uG.POST("/resetpasswd", userResetpasswd)
		uG.POST("/forgetpasswd", userForgetpasswd)
		uG.POST("/sendverifyemail", userSendVerifyEmail)
		uG.POST("/verifyemail", userVerifyEmail)
		uG.POST("/addinvitecodes", userAddInviteCodes)
		uG.GET("/:userURL", userInfo)
		// don't remove user //
		// uG.DELETE("/:userURL/remove", userDelete)
		uG.DELETE("/:userURL/remove", userDeactive)
		uG.GET("/:userURL/info", userInfo)
		uG.POST("/follow/:userURL", userFollow)
		uG.POST("/update/:userURL", userUpdate)
		uG.POST("/active/:userURL", userActive)
		uG.POST("/deactive/:userURL", userDeactive)
		uG.GET("/:userURL/challenges", userChallenges)
		// challenges[?type=0/1/2]
		uG.GET("/:userURL/followers", userFollowers)
		uG.GET("/:userURL/followees", userFollowees)
	}

	// Challenges
	router.GET("/challenges", challenges)
	// cG = challenge group
	cG := router.Group("/challenge")
	{
		cG.GET("/:challengeID", challengeInfo)
		cG.POST("/validate/:challengeID", challengeValidateFlag)
		cG.POST("/create", challengeCreate)
		cG.POST("/remove", challengeRemove)
	}

	// templates
	router.GET("/templates", templates)
	// tG = template group
	tG := router.Group("/template")
	{
		tG.POST("/create", templateCreate)
		tG.GET("/:templateID", templateQuery)
		tG.POST("/enable/:templateID", templateEnable)
		tG.POST("/disable/:templateID", templateDisable)
		tG.POST("/update/:templateID", templateUpdate)
		tG.DELETE("/:templateID/remove", templateRemove)
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
		c.JSON(400, err)
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
