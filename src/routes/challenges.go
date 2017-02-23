package routes

import (
	"fmt"
	Challenge "handler/challenge"

	"github.com/gin-gonic/gin"
	"github.com/nu7hatch/gouuid"
)

func challenges(c *gin.Context) {
	challenges, err := Challenge.QueryChallenges("all")
	if err != nil {
		c.JSON(500, gin.H{
			"err": err.Error(),
		})
	} else {
		c.JSON(200, challenges)
	}
}

func challengeInfo(c *gin.Context) {
	cid := c.Param("challengeID")
	challenges, err := Challenge.QueryChallenges(cid)
	if err != nil {
		c.JSON(500, gin.H{
			"err": err.Error(),
		})
	} else {
		c.JSON(200, challenges[0])
	}
}

func challengeCreate(c *gin.Context) {
	userID := c.Request.PostFormValue("uid")
	templateID := c.Request.PostFormValue("templateID")

	if !validateUser(userID) {
		c.JSON(500, gin.H{
			"err": "user [" + userID + "] not found",
		})
		return
	}

	if !validateTemplate(templateID) {
		c.JSON(500, gin.H{
			"err": "template [" + templateID + "] not found",
		})
		return
	}

	uid, _ := uuid.NewV4()
	challengeID := fmt.Sprintf("%v", uid)

	go func() {
		_, err := Challenge.CreateChallenge(userID,
			templateID, challengeID)
		if err != nil {
			Challenge.UpdateChallengeState(challengeID, "failed")
		} else {
			Challenge.UpdateChallengeState(challengeID, "running")
		}
	}()
	c.JSON(201, gin.H{
		"id": challengeID,
	})
}
