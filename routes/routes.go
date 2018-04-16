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
		uG.POST("", userCreate)
		userSelfGroup := uG.Group("/self")
		{
			userSelfGroup.GET("/:userURL", userInfo)
			userSelfGroup.GET("/:userURL/info", userInfo)
			userSelfGroup.GET("/:userURL/challenges", userChallenges)
			// challenges[?type=0/1/2]
			userSelfGroup.GET("/:userURL/followers", userFollowers)
			userSelfGroup.GET("/:userURL/followees", userFollowees)
			// don't remove user //
			// uG.DELETE("/:userURL/remove", userDelete)
			userSelfGroup.DELETE("/:userURL", userDeactive)
		}
		userActivationGroup := uG.Group("/activation")
		{
			userActivationGroup.POST("/:userURL", userActive)
			userActivationGroup.DELETE("/:userURL", userDeactive)
		}
		// TODO 这些API怎么抽啊QAQ，全是动作，规范又不能带有动词，哇
		uG.POST("/login", userLogin)
		uG.POST("/passwd", userResetpasswd)
		uG.POST("/passwd/forgotten", userForgetpasswd)
		uG.POST("/email", userSendVerifyEmail)
		uG.POST("/email/verification", userVerifyEmail)
		uG.POST("/invitecodes", userAddInviteCodes)
		uG.POST("/follow/:userID", userFollow)
		uG.POST("/update/:userURL", userUpdate)

	}

	// Challenges
	router.GET("/challenges", challenges)
	// cG = challenge group
	cG := router.Group("/challenge")
	{
		cG.GET("/:challengeID", challengeInfo)
		cG.POST("/:challengeID/validation", challengeValidateFlag)
		cG.POST("", challengeCreate)                // create
		cG.DELETE("/:challengeID", challengeRemove) // delete
	}

	// templates
	router.GET("/templates", templates)
	// tG = template group
	tG := router.Group("/template")
	{
		tG.PUT("", templateCreate)
		tG.GET("/:templateID", templateQuery)
		// TODO 把这俩开关改掉
		tG.POST("/:templateID/enable", templateEnable)
		tG.DELETE("/:templateID/enable", templateDisable)
		tG.POST("/:templateID", templateUpdate)
		tG.DELETE("/:templateID", templateRemove)
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
